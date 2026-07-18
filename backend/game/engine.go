package game

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	"rakshyak-98/pokemon-backend/models"
)

// Engine owns mutable game state and turn rules (State-machine style transitions).
type Engine struct {
	State   *models.GameState
	catalog []models.Pokemon
}

func NewEngine(catalog []models.Pokemon) *Engine {
	return &Engine{
		catalog: catalog,
		State: &models.GameState{
			Status:     models.StatusWaiting,
			Players:    []models.PlayerState{},
			TurnNumber: 0,
		},
	}
}

// SetCatalog updates the Pokémon pool used when building decks.
func (e *Engine) SetCatalog(catalog []models.Pokemon) {
	e.catalog = catalog
}

func (e *Engine) StateSnapshot() *models.GameState {
	return e.State
}

// StartGame initializes a new game with two players, decks, prizes, and opening hands.
func (e *Engine) StartGame(player1ID, player2ID string) error {
	if player1ID == "" || player2ID == "" {
		return errors.New("both player ids are required")
	}
	if player1ID == player2ID {
		return errors.New("player ids must be different")
	}

	e.State = &models.GameState{
		ID:          uuid.NewString(),
		Status:      models.StatusInProgress,
		Players:     []models.PlayerState{e.createPlayer(player1ID), e.createPlayer(player2ID)},
		CurrentTurn: player1ID,
		TurnNumber:  1,
		LastAction:  "game_started",
		UpdatedAt:   time.Now().UTC(),
	}

	for i := range e.State.Players {
		p := &e.State.Players[i]
		for j := 0; j < 6; j++ {
			if len(p.Deck) == 0 {
				return errors.New("deck too small for prizes")
			}
			card := p.Deck[0]
			p.Deck = p.Deck[1:]
			p.PrizeCards = append(p.PrizeCards, card)
		}
		for j := 0; j < 7; j++ {
			if err := e.drawCard(p); err != nil {
				return err
			}
		}
		p.HasDrawn = false
		p.HasAttached = false
	}

	return nil
}

func (e *Engine) createPlayer(id string) models.PlayerState {
	deck := e.buildStarterDeck(id)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return models.PlayerState{
		ID:             id,
		Deck:           deck,
		Hand:           []models.Card{},
		BenchedPokemon: []models.Card{},
		PrizeCards:     []models.Card{},
		DiscardPile:    []models.Card{},
	}
}

func fallbackCatalog() []models.Pokemon {
	return []models.Pokemon{
		{
			PokeAPIID: 25, Name: "Pikachu", PrimaryType: "Electric", Types: []string{"Electric"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/25.png",
			Stats:    models.PokemonStats{HP: 35, Attack: 55, Defense: 40, SpAttack: 50, SpDefense: 50, Speed: 90},
			CardHP:   70, Attacks: []models.Attack{{Name: "Thunder Shock", Damage: 20, Cost: 1}},
		},
		{
			PokeAPIID: 4, Name: "Charmander", PrimaryType: "Fire", Types: []string{"Fire"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/4.png",
			Stats:    models.PokemonStats{HP: 39, Attack: 52, Defense: 43, SpAttack: 60, SpDefense: 50, Speed: 65},
			CardHP:   78, Attacks: []models.Attack{{Name: "Ember", Damage: 30, Cost: 1}},
		},
		{
			PokeAPIID: 7, Name: "Squirtle", PrimaryType: "Water", Types: []string{"Water"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/7.png",
			Stats:    models.PokemonStats{HP: 44, Attack: 48, Defense: 65, SpAttack: 50, SpDefense: 64, Speed: 43},
			CardHP:   88, Attacks: []models.Attack{{Name: "Water Gun", Damage: 20, Cost: 1}},
		},
		{
			PokeAPIID: 1, Name: "Bulbasaur", PrimaryType: "Grass", Types: []string{"Grass", "Poison"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/1.png",
			Stats:    models.PokemonStats{HP: 45, Attack: 49, Defense: 49, SpAttack: 65, SpDefense: 65, Speed: 45},
			CardHP:   90, Attacks: []models.Attack{{Name: "Vine Whip", Damage: 20, Cost: 1}},
		},
	}
}

func (e *Engine) catalogPool() []models.Pokemon {
	if len(e.catalog) > 0 {
		return e.catalog
	}
	return fallbackCatalog()
}

func cardFromPokemon(idPrefix string, index int, p models.Pokemon) models.Card {
	stats := p.Stats
	attacks := append([]models.Attack(nil), p.Attacks...)
	return models.Card{
		ID:          fmt.Sprintf("%s-poke-%d", idPrefix, index),
		Name:        p.Name,
		Type:        models.TypePokemon,
		HP:          p.CardHP,
		MaxHP:       p.CardHP,
		Attacks:     attacks,
		ImageURL:    p.ImageURL,
		ElementType: p.PrimaryType,
		PokeAPIID:   p.PokeAPIID,
		Stats:       &stats,
	}
}

func energyForType(element string) string {
	switch strings.ToLower(element) {
	case "electric":
		return "Electric"
	case "fire":
		return "Fire"
	case "water", "ice":
		return "Water"
	case "grass", "bug", "poison":
		return "Grass"
	case "psychic", "ghost", "fairy", "dark":
		return "Psychic"
	default:
		return "Normal"
	}
}

func (e *Engine) buildStarterDeck(id string) []models.Card {
	pool := e.catalogPool()
	deck := make([]models.Card, 0, 60)

	// 20 Pokémon cards drawn from the seeded catalog.
	for i := 0; i < 20; i++ {
		p := pool[i%len(pool)]
		deck = append(deck, cardFromPokemon(id, i, p))
	}

	energyTypes := make([]string, 0, 4)
	seen := map[string]bool{}
	for _, p := range pool {
		et := energyForType(p.PrimaryType)
		if !seen[et] {
			seen[et] = true
			energyTypes = append(energyTypes, et)
		}
		if len(energyTypes) >= 4 {
			break
		}
	}
	if len(energyTypes) == 0 {
		energyTypes = []string{"Electric", "Fire", "Water", "Grass"}
	}

	for i := 0; i < 40; i++ {
		et := energyTypes[i%len(energyTypes)]
		deck = append(deck, models.Card{
			ID:          fmt.Sprintf("%s-energy-%d", id, i),
			Name:        et + " Energy",
			Type:        models.TypeEnergy,
			EnergyType:  et,
			ElementType: et,
		})
	}
	return deck
}

func (e *Engine) drawCard(player *models.PlayerState) error {
	if len(player.Deck) == 0 {
		return errors.New("deck is empty")
	}
	card := player.Deck[0]
	player.Deck = player.Deck[1:]
	player.Hand = append(player.Hand, card)
	return nil
}

func (e *Engine) getPlayer(playerID string) *models.PlayerState {
	for i := range e.State.Players {
		if e.State.Players[i].ID == playerID {
			return &e.State.Players[i]
		}
	}
	return nil
}

func (e *Engine) getOpponent(playerID string) *models.PlayerState {
	for i := range e.State.Players {
		if e.State.Players[i].ID != playerID {
			return &e.State.Players[i]
		}
	}
	return nil
}

func (e *Engine) requireTurn(playerID string) (*models.PlayerState, error) {
	if e.State.Status == models.StatusGameOver {
		return nil, errors.New("game is over")
	}
	if e.State.Status != models.StatusInProgress {
		return nil, errors.New("game not in progress")
	}
	player := e.getPlayer(playerID)
	if player == nil {
		return nil, errors.New("player not found")
	}
	if e.State.CurrentTurn != playerID {
		return nil, errors.New("not your turn")
	}
	return player, nil
}

func (e *Engine) requirePlayable(playerID string) (*models.PlayerState, error) {
	player, err := e.requireTurn(playerID)
	if err != nil {
		return nil, err
	}
	if player.ActivePokemon == nil && len(player.BenchedPokemon) > 0 {
		return nil, errors.New("must promote a benched pokemon first")
	}
	return player, nil
}

func (e *Engine) DrawCard(playerID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.HasDrawn {
		return errors.New("already drew this turn")
	}
	if err := e.drawCard(player); err != nil {
		return err
	}
	player.HasDrawn = true
	e.State.LastAction = fmt.Sprintf("%s drew a card", playerID)
	return nil
}

func (e *Engine) PlayBench(playerID, cardID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	idx := findCard(player.Hand, cardID)
	if idx < 0 {
		return errors.New("card not found in hand")
	}
	card := player.Hand[idx]
	if card.Type != models.TypePokemon {
		return errors.New("can only bench pokemon")
	}
	if len(player.BenchedPokemon) >= 5 {
		return errors.New("bench is full")
	}
	player.Hand = append(player.Hand[:idx], player.Hand[idx+1:]...)
	player.BenchedPokemon = append(player.BenchedPokemon, card)
	e.State.LastAction = fmt.Sprintf("%s benched %s", playerID, card.Name)
	return nil
}

func (e *Engine) SetActive(playerID, cardID string) error {
	player, err := e.requireTurn(playerID)
	if err != nil {
		return err
	}
	if player.ActivePokemon != nil {
		return errors.New("already have an active pokemon")
	}

	if idx := findCard(player.BenchedPokemon, cardID); idx >= 0 {
		card := player.BenchedPokemon[idx]
		player.BenchedPokemon = append(player.BenchedPokemon[:idx], player.BenchedPokemon[idx+1:]...)
		player.ActivePokemon = &card
		e.State.LastAction = fmt.Sprintf("%s set active %s", playerID, card.Name)
		return nil
	}

	if idx := findCard(player.Hand, cardID); idx >= 0 {
		card := player.Hand[idx]
		if card.Type != models.TypePokemon {
			return errors.New("can only set pokemon as active")
		}
		player.Hand = append(player.Hand[:idx], player.Hand[idx+1:]...)
		player.ActivePokemon = &card
		e.State.LastAction = fmt.Sprintf("%s set active %s", playerID, card.Name)
		return nil
	}
	return errors.New("card not found in hand or bench")
}

// Promote moves a benched Pokemon to active after a knockout.
func (e *Engine) Promote(playerID, cardID string) error {
	if e.State.Status != models.StatusInProgress {
		return errors.New("game not in progress")
	}
	player := e.getPlayer(playerID)
	if player == nil {
		return errors.New("player not found")
	}
	if player.ActivePokemon != nil {
		return errors.New("already have an active pokemon")
	}
	idx := findCard(player.BenchedPokemon, cardID)
	if idx < 0 {
		return errors.New("card not found on bench")
	}
	card := player.BenchedPokemon[idx]
	player.BenchedPokemon = append(player.BenchedPokemon[:idx], player.BenchedPokemon[idx+1:]...)
	player.ActivePokemon = &card
	e.State.LastAction = fmt.Sprintf("%s promoted %s", playerID, card.Name)
	return nil
}

func (e *Engine) AttachEnergy(playerID, energyCardID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.HasAttached {
		return errors.New("already attached energy this turn")
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon to attach energy to")
	}
	idx := findCard(player.Hand, energyCardID)
	if idx < 0 {
		return errors.New("energy card not found in hand")
	}
	card := player.Hand[idx]
	if card.Type != models.TypeEnergy {
		return errors.New("card is not an energy")
	}
	player.Hand = append(player.Hand[:idx], player.Hand[idx+1:]...)
	player.ActivePokemon.EnergyAttached++
	player.HasAttached = true
	e.State.LastAction = fmt.Sprintf("%s attached %s to %s", playerID, card.Name, player.ActivePokemon.Name)
	return nil
}

func (e *Engine) Attack(playerID string, attackIndex int) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon")
	}
	if attackIndex < 0 || attackIndex >= len(player.ActivePokemon.Attacks) {
		return errors.New("invalid attack index")
	}
	attack := player.ActivePokemon.Attacks[attackIndex]
	if player.ActivePokemon.EnergyAttached < attack.Cost {
		return fmt.Errorf("not enough energy (need %d)", attack.Cost)
	}

	opponent := e.getOpponent(playerID)
	if opponent == nil || opponent.ActivePokemon == nil {
		return errors.New("opponent has no active pokemon")
	}

	opponent.ActivePokemon.HP -= attack.Damage
	e.State.LastAction = fmt.Sprintf("%s used %s for %d damage", playerID, attack.Name, attack.Damage)

	if opponent.ActivePokemon.HP <= 0 {
		e.resolveKnockout(player, opponent)
	}

	if e.State.Status == models.StatusGameOver {
		return nil
	}

	return e.advanceTurn(playerID)
}

func (e *Engine) resolveKnockout(attacker, defender *models.PlayerState) {
	knocked := *defender.ActivePokemon
	defender.DiscardPile = append(defender.DiscardPile, knocked)
	defender.ActivePokemon = nil

	if len(attacker.PrizeCards) > 0 {
		prize := attacker.PrizeCards[0]
		attacker.PrizeCards = attacker.PrizeCards[1:]
		attacker.Hand = append(attacker.Hand, prize)
		attacker.PrizesTaken++
	}

	if attacker.PrizesTaken >= 6 {
		e.State.Status = models.StatusGameOver
		e.State.Winner = attacker.ID
		e.State.LastAction = fmt.Sprintf("%s wins by taking all prizes", attacker.ID)
		return
	}

	if len(defender.BenchedPokemon) == 0 {
		e.State.Status = models.StatusGameOver
		e.State.Winner = attacker.ID
		e.State.LastAction = fmt.Sprintf("%s wins — opponent has no pokemon left", attacker.ID)
		return
	}

	e.State.LastAction += fmt.Sprintf("; %s must promote", defender.ID)
}

func (e *Engine) EndTurn(playerID string) error {
	if _, err := e.requirePlayable(playerID); err != nil {
		return err
	}
	e.State.LastAction = fmt.Sprintf("%s ended turn", playerID)
	return e.advanceTurn(playerID)
}

func (e *Engine) advanceTurn(currentPlayerID string) error {
	opponent := e.getOpponent(currentPlayerID)
	if opponent == nil {
		return errors.New("opponent not found")
	}
	e.State.CurrentTurn = opponent.ID
	e.State.TurnNumber++
	opponent.HasDrawn = false
	opponent.HasAttached = false
	return nil
}

func findCard(cards []models.Card, id string) int {
	for i, c := range cards {
		if c.ID == id {
			return i
		}
	}
	return -1
}

// Restore loads a persisted state into the engine (Memento restore).
func (e *Engine) Restore(state *models.GameState) {
	e.State = state
}
