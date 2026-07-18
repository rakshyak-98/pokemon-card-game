package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"rakshyak-98/pokemon-backend/models"
)

const defaultBaseURL = "https://pokeapi.co/api/v2"

// Client fetches Pokémon data from the public PokeAPI.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

type apiPokemon struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
		Other        struct {
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
			} `json:"official-artwork"`
		} `json:"other"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
	} `json:"moves"`
}

// FetchPokemon loads one Pokémon by national dex id.
func (c *Client) FetchPokemon(id int) (models.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%d", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.Pokemon{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "pokemon-card-game/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return models.Pokemon{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return models.Pokemon{}, fmt.Errorf("pokeapi %d: %s (%s)", id, resp.Status, strings.TrimSpace(string(body)))
	}

	var raw apiPokemon
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return models.Pokemon{}, err
	}
	return mapPokemon(raw), nil
}

func mapPokemon(raw apiPokemon) models.Pokemon {
	stats := models.PokemonStats{}
	for _, s := range raw.Stats {
		switch s.Stat.Name {
		case "hp":
			stats.HP = s.BaseStat
		case "attack":
			stats.Attack = s.BaseStat
		case "defense":
			stats.Defense = s.BaseStat
		case "special-attack":
			stats.SpAttack = s.BaseStat
		case "special-defense":
			stats.SpDefense = s.BaseStat
		case "speed":
			stats.Speed = s.BaseStat
		}
	}

	types := make([]string, 0, len(raw.Types))
	primary := ""
	for _, t := range raw.Types {
		name := titleCase(t.Type.Name)
		types = append(types, name)
		if t.Slot == 1 || primary == "" {
			primary = name
		}
	}

	imageURL := raw.Sprites.Other.OfficialArtwork.FrontDefault
	if imageURL == "" {
		imageURL = raw.Sprites.FrontDefault
	}

	cardHP := stats.HP * 2
	if cardHP < 40 {
		cardHP = 40
	}
	if cardHP > 200 {
		cardHP = 200
	}

	offensive := stats.Attack
	if stats.SpAttack > offensive {
		offensive = stats.SpAttack
	}
	damage := offensive / 2
	if damage < 10 {
		damage = 10
	}
	if damage > 80 {
		damage = 80
	}
	cost := 1
	if damage >= 40 {
		cost = 2
	}

	attackName := defaultAttackName(primary)
	if len(raw.Moves) > 0 && raw.Moves[0].Move.Name != "" {
		attackName = titleCase(strings.ReplaceAll(raw.Moves[0].Move.Name, "-", " "))
	}

	return models.Pokemon{
		PokeAPIID:   raw.ID,
		Name:        titleCase(raw.Name),
		ImageURL:    imageURL,
		PrimaryType: primary,
		Types:       types,
		Stats:       stats,
		CardHP:      cardHP,
		Attacks: []models.Attack{
			{Name: attackName, Damage: damage, Cost: cost},
		},
	}
}

func defaultAttackName(element string) string {
	switch strings.ToLower(element) {
	case "electric":
		return "Thunder Shock"
	case "fire":
		return "Ember"
	case "water":
		return "Water Gun"
	case "grass":
		return "Vine Whip"
	case "psychic":
		return "Confusion"
	case "fighting":
		return "Low Kick"
	case "poison":
		return "Poison Sting"
	case "ground":
		return "Mud-Slap"
	case "rock":
		return "Rock Throw"
	case "bug":
		return "Bug Bite"
	case "ghost":
		return "Lick"
	case "ice":
		return "Powder Snow"
	case "dragon":
		return "Dragon Rage"
	case "dark":
		return "Bite"
	case "steel":
		return "Metal Claw"
	case "fairy":
		return "Fairy Wind"
	case "flying":
		return "Gust"
	case "normal":
		return "Tackle"
	default:
		return "Strike"
	}
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Fields(strings.ReplaceAll(s, "-", " "))
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, " ")
}

// PokemonWriter persists catalog rows (implemented by store.SQLiteStore).
type PokemonWriter interface {
	CountPokemon() (int, error)
	UpsertPokemon(p models.Pokemon) error
}

// SeedOptions controls how many Pokémon are loaded from PokeAPI.
type SeedOptions struct {
	// FromID / ToID are inclusive national dex ids (default Gen 1: 1–151).
	FromID     int
	ToID       int
	Workers    int
	Force      bool // re-fetch even when the table already has rows
	OnProgress func(done, total int, name string)
}

// SeedIfEmpty populates the local DB from PokeAPI when empty (or Force).
func SeedIfEmpty(w PokemonWriter, client *Client, opts SeedOptions) error {
	if opts.FromID <= 0 {
		opts.FromID = 1
	}
	if opts.ToID <= 0 {
		opts.ToID = 151
	}
	if opts.Workers <= 0 {
		opts.Workers = 6
	}
	if opts.ToID < opts.FromID {
		return fmt.Errorf("invalid seed range %d–%d", opts.FromID, opts.ToID)
	}

	count, err := w.CountPokemon()
	if err != nil {
		return err
	}
	expected := opts.ToID - opts.FromID + 1
	if !opts.Force && count >= expected {
		log.Printf("pokemon catalog already seeded (%d rows)", count)
		return nil
	}
	if !opts.Force && count > 0 && count < expected {
		log.Printf("pokemon catalog incomplete (%d/%d); filling gaps", count, expected)
	}

	if client == nil {
		client = NewClient()
	}

	ids := make([]int, 0, expected)
	for id := opts.FromID; id <= opts.ToID; id++ {
		ids = append(ids, id)
	}

	type result struct {
		p   models.Pokemon
		err error
		id  int
	}

	jobs := make(chan int)
	results := make(chan result)
	var wg sync.WaitGroup
	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				var p models.Pokemon
				var err error
				for attempt := 0; attempt < 3; attempt++ {
					p, err = client.FetchPokemon(id)
					if err == nil {
						break
					}
					time.Sleep(time.Duration(attempt+1) * 300 * time.Millisecond)
				}
				results <- result{p: p, err: err, id: id}
			}
		}()
	}

	go func() {
		for _, id := range ids {
			jobs <- id
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	done := 0
	var firstErr error
	for res := range results {
		done++
		if res.err != nil {
			log.Printf("seed pokemon %d failed: %v", res.id, res.err)
			if firstErr == nil {
				firstErr = res.err
			}
			continue
		}
		if err := w.UpsertPokemon(res.p); err != nil {
			log.Printf("upsert pokemon %d failed: %v", res.id, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if opts.OnProgress != nil {
			opts.OnProgress(done, expected, res.p.Name)
		} else if done%25 == 0 || done == expected {
			log.Printf("seeded pokemon %d/%d (%s)", done, expected, res.p.Name)
		}
	}

	finalCount, _ := w.CountPokemon()
	log.Printf("pokemon catalog ready: %d entries", finalCount)
	if finalCount == 0 && firstErr != nil {
		return firstErr
	}
	return nil
}
