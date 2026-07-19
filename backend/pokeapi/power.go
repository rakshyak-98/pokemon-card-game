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

// Battle-relevant PokeAPI item categories used to discover power cards.
// Categories mirror ASC / MEW Trainer Item themes: healing, battle boosts, berries, held tools.
var powerItemCategories = []string{
	"stat-boosts",
	"healing",
	"status-cures",
	"revival",
	"medicine",
	"in-a-pinch",
	"miracle-shooter",
	"held-items",
	"type-enhancement",
	"vitamins",
}

// powerItemSpec maps a PokeAPI item slug onto a game power effect.
// Values are balanced for GO-style HP / damage ranges (not literal main-series numbers).
// Catalog expands ASC rulebook Trainer Item ideas (Potion, X Attack, Guard Spec, berries, tools).
type powerItemSpec struct {
	Effect string
	Value  int
}

var powerItemSpecs = map[string]powerItemSpec{
	// Attack boosts — X Attack / Dire Hit family + vitamins & held tools
	"x-attack":       {Effect: "boost_attack", Value: 20},
	"x-sp-atk":       {Effect: "boost_attack", Value: 20},
	"x-attack-2":     {Effect: "boost_attack", Value: 25},
	"x-sp-atk-2":     {Effect: "boost_attack", Value: 25},
	"x-attack-3":     {Effect: "boost_attack", Value: 30},
	"x-sp-atk-3":     {Effect: "boost_attack", Value: 30},
	"x-attack-6":     {Effect: "boost_attack", Value: 40},
	"x-sp-atk-6":     {Effect: "boost_attack", Value: 40},
	"dire-hit":       {Effect: "boost_attack", Value: 15},
	"dire-hit-2":     {Effect: "boost_attack", Value: 20},
	"dire-hit-3":     {Effect: "boost_attack", Value: 25},
	"liechi-berry":   {Effect: "boost_attack", Value: 25},
	"petaya-berry":   {Effect: "boost_attack", Value: 25},
	"protein":        {Effect: "boost_attack", Value: 20},
	"calcium":        {Effect: "boost_attack", Value: 20},
	"choice-band":    {Effect: "boost_attack", Value: 30},
	"choice-specs":   {Effect: "boost_attack", Value: 30},
	"muscle-band":    {Effect: "boost_attack", Value: 25},
	"wise-glasses":   {Effect: "boost_attack", Value: 25},
	"life-orb":       {Effect: "boost_attack", Value: 35},
	"expert-belt":    {Effect: "boost_attack", Value: 28},
	"black-belt":     {Effect: "boost_attack", Value: 22},
	"charcoal":       {Effect: "boost_attack", Value: 18},
	"mystic-water":   {Effect: "boost_attack", Value: 18},
	"miracle-seed":   {Effect: "boost_attack", Value: 18},
	"magnet":         {Effect: "boost_attack", Value: 18},
	"twisted-spoon":  {Effect: "boost_attack", Value: 18},
	"never-melt-ice": {Effect: "boost_attack", Value: 18},
	"black-glasses":  {Effect: "boost_attack", Value: 18},
	"spell-tag":      {Effect: "boost_attack", Value: 18},
	"sharp-beak":     {Effect: "boost_attack", Value: 18},
	"silk-scarf":     {Effect: "boost_attack", Value: 18},
	"dragon-fang":    {Effect: "boost_attack", Value: 18},
	"metal-coat":     {Effect: "boost_attack", Value: 18},
	"soft-sand":      {Effect: "boost_attack", Value: 18},
	"poison-barb":    {Effect: "boost_attack", Value: 18},
	"silver-powder":  {Effect: "boost_attack", Value: 18},

	// Defense boosts — X Defense / Guard Spec family + held tools
	"x-defense":    {Effect: "boost_defense", Value: 15},
	"x-sp-def":     {Effect: "boost_defense", Value: 15},
	"x-defense-2":  {Effect: "boost_defense", Value: 20},
	"x-sp-def-2":   {Effect: "boost_defense", Value: 20},
	"x-defense-3":  {Effect: "boost_defense", Value: 25},
	"x-sp-def-3":   {Effect: "boost_defense", Value: 25},
	"x-defense-6":  {Effect: "boost_defense", Value: 35},
	"x-sp-def-6":   {Effect: "boost_defense", Value: 35},
	"guard-spec":   {Effect: "boost_defense", Value: 10},
	"ganlon-berry": {Effect: "boost_defense", Value: 20},
	"apicot-berry": {Effect: "boost_defense", Value: 20},
	"iron":         {Effect: "boost_defense", Value: 18},
	"zinc":         {Effect: "boost_defense", Value: 18},
	"assault-vest": {Effect: "boost_defense", Value: 25},
	"eviolite":     {Effect: "boost_defense", Value: 30},
	"focus-band":   {Effect: "boost_defense", Value: 15},
	"focus-sash":   {Effect: "boost_defense", Value: 20},
	"rocky-helmet": {Effect: "boost_defense", Value: 22},

	// Healing — Potion family, drinks, berries, revival & status cures (ASC Item theme)
	"potion":          {Effect: "heal", Value: 20},
	"super-potion":    {Effect: "heal", Value: 35},
	"hyper-potion":    {Effect: "heal", Value: 50},
	"max-potion":      {Effect: "heal", Value: 70},
	"full-restore":    {Effect: "heal", Value: 80},
	"fresh-water":     {Effect: "heal", Value: 25},
	"soda-pop":        {Effect: "heal", Value: 30},
	"lemonade":        {Effect: "heal", Value: 35},
	"moomoo-milk":     {Effect: "heal", Value: 40},
	"energy-powder":   {Effect: "heal", Value: 25},
	"energy-root":     {Effect: "heal", Value: 45},
	"berry-juice":     {Effect: "heal", Value: 20},
	"sweet-heart":     {Effect: "heal", Value: 20},
	"full-heal":       {Effect: "heal", Value: 15},
	"heal-powder":     {Effect: "heal", Value: 20},
	"revival-herb":    {Effect: "heal", Value: 55},
	"revive":          {Effect: "heal", Value: 45},
	"max-revive":      {Effect: "heal", Value: 70},
	"sacred-ash":      {Effect: "heal", Value: 80},
	"oran-berry":      {Effect: "heal", Value: 20},
	"sitrus-berry":    {Effect: "heal", Value: 40},
	"figy-berry":      {Effect: "heal", Value: 30},
	"wiki-berry":      {Effect: "heal", Value: 30},
	"mago-berry":      {Effect: "heal", Value: 30},
	"aguav-berry":     {Effect: "heal", Value: 30},
	"iapapa-berry":    {Effect: "heal", Value: 30},
	"lum-berry":       {Effect: "heal", Value: 25},
	"leftovers":       {Effect: "heal", Value: 25},
	"shell-bell":      {Effect: "heal", Value: 22},
	"lava-cookie":     {Effect: "heal", Value: 20},
	"old-gateau":      {Effect: "heal", Value: 20},
	"casteliacone":    {Effect: "heal", Value: 20},
	"lumiose-galette": {Effect: "heal", Value: 20},
	"shalour-sable":   {Effect: "heal", Value: 20},
	"rage-candy-bar":  {Effect: "heal", Value: 20},
	"antidote":        {Effect: "heal", Value: 10},
	"burn-heal":       {Effect: "heal", Value: 10},
	"ice-heal":        {Effect: "heal", Value: 10},
	"awakening":       {Effect: "heal", Value: 10},
	"paralyze-heal":   {Effect: "heal", Value: 10},
	"hp-up":           {Effect: "heal", Value: 35},
}

type apiNamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type apiItemCategory struct {
	ID    int                `json:"id"`
	Name  string             `json:"name"`
	Items []apiNamedResource `json:"items"`
}

type apiItem struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Category      apiNamedResource `json:"category"`
	EffectEntries []struct {
		Effect      string `json:"effect"`
		ShortEffect string `json:"short_effect"`
		Language    struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"effect_entries"`
	Sprites struct {
		Default string `json:"default"`
	} `json:"sprites"`
}

// PowerWriter persists power-card catalog rows (implemented by store.SQLiteStore).
type PowerWriter interface {
	CountPowerCards() (int, error)
	UpsertPowerCard(p models.PowerCard) error
}

// PowerSeedOptions controls how power cards are loaded from PokeAPI items.
type PowerSeedOptions struct {
	Workers    int
	Force      bool // re-fetch even when the table already has rows
	OnProgress func(done, total int, name string)
}

// FetchItem loads one item by PokeAPI id or slug name.
func (c *Client) FetchItem(idOrName string) (apiItem, error) {
	url := fmt.Sprintf("%s/item/%s", c.BaseURL, idOrName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return apiItem{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "pokemon-card-game/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return apiItem{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return apiItem{}, fmt.Errorf("pokeapi item %s: %s (%s)", idOrName, resp.Status, strings.TrimSpace(string(body)))
	}

	var raw apiItem
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return apiItem{}, err
	}
	return raw, nil
}

// FetchItemCategory lists items in a PokeAPI item category.
func (c *Client) FetchItemCategory(name string) (apiItemCategory, error) {
	url := fmt.Sprintf("%s/item-category/%s", c.BaseURL, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return apiItemCategory{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "pokemon-card-game/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return apiItemCategory{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return apiItemCategory{}, fmt.Errorf("pokeapi item-category %s: %s (%s)", name, resp.Status, strings.TrimSpace(string(body)))
	}

	var raw apiItemCategory
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return apiItemCategory{}, err
	}
	return raw, nil
}

// MapItemToPowerCard converts a PokeAPI item into a power card when the slug is known.
func MapItemToPowerCard(raw apiItem) (models.PowerCard, bool) {
	spec, ok := powerItemSpecs[raw.Name]
	if !ok {
		return models.PowerCard{}, false
	}
	desc := ""
	for _, e := range raw.EffectEntries {
		if e.Language.Name == "en" {
			desc = e.ShortEffect
			if desc == "" {
				desc = e.Effect
			}
			break
		}
	}
	imageURL := raw.Sprites.Default
	if imageURL == "" {
		imageURL = fmt.Sprintf(
			"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/%s.png",
			raw.Name,
		)
	}
	return models.PowerCard{
		PokeAPIID:   raw.ID,
		Name:        titleCase(raw.Name),
		ImageURL:    imageURL,
		Effect:      spec.Effect,
		EffectValue: spec.Value,
		Category:    raw.Category.Name,
		Description: desc,
	}, true
}

// DiscoverPowerItemSlugs walks battle-relevant item categories and returns slugs
// that have a known power-card mapping.
func (c *Client) DiscoverPowerItemSlugs() ([]string, error) {
	seen := map[string]struct{}{}
	var slugs []string
	for _, cat := range powerItemCategories {
		detail, err := c.FetchItemCategory(cat)
		if err != nil {
			return nil, err
		}
		for _, item := range detail.Items {
			if _, ok := powerItemSpecs[item.Name]; !ok {
				continue
			}
			if _, dup := seen[item.Name]; dup {
				continue
			}
			seen[item.Name] = struct{}{}
			slugs = append(slugs, item.Name)
		}
	}
	// Always include curated specs even if a category fetch missed them.
	for slug := range powerItemSpecs {
		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}
		slugs = append(slugs, slug)
	}
	return slugs, nil
}

// SeedPowerIfEmpty populates the local power_cards table from PokeAPI items.
func SeedPowerIfEmpty(w PowerWriter, client *Client, opts PowerSeedOptions) error {
	if opts.Workers <= 0 {
		opts.Workers = 6
	}

	count, err := w.CountPowerCards()
	if err != nil {
		return err
	}
	expected := len(powerItemSpecs)
	if !opts.Force && count >= expected {
		log.Printf("power card catalog already seeded (%d rows)", count)
		return nil
	}
	if !opts.Force && count > 0 && count < expected {
		log.Printf("power card catalog incomplete (%d/%d); filling gaps", count, expected)
	}

	if client == nil {
		client = NewClient()
	}

	slugs, err := client.DiscoverPowerItemSlugs()
	if err != nil {
		log.Printf("power item category discovery failed (%v); using curated slug list", err)
		slugs = make([]string, 0, len(powerItemSpecs))
		for slug := range powerItemSpecs {
			slugs = append(slugs, slug)
		}
	}

	type result struct {
		p    models.PowerCard
		ok   bool
		err  error
		slug string
	}

	jobs := make(chan string)
	results := make(chan result)
	var wg sync.WaitGroup
	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for slug := range jobs {
				var raw apiItem
				var fetchErr error
				for attempt := 0; attempt < 3; attempt++ {
					raw, fetchErr = client.FetchItem(slug)
					if fetchErr == nil {
						break
					}
					time.Sleep(time.Duration(attempt+1) * 300 * time.Millisecond)
				}
				if fetchErr != nil {
					results <- result{err: fetchErr, slug: slug}
					continue
				}
				card, ok := MapItemToPowerCard(raw)
				results <- result{p: card, ok: ok, slug: slug}
			}
		}()
	}

	go func() {
		for _, slug := range slugs {
			jobs <- slug
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	done := 0
	total := len(slugs)
	var firstErr error
	for res := range results {
		done++
		if res.err != nil {
			log.Printf("seed power item %s failed: %v", res.slug, res.err)
			if firstErr == nil {
				firstErr = res.err
			}
			continue
		}
		if !res.ok {
			continue
		}
		if err := w.UpsertPowerCard(res.p); err != nil {
			log.Printf("upsert power card %s failed: %v", res.slug, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if opts.OnProgress != nil {
			opts.OnProgress(done, total, res.p.Name)
		} else if done%10 == 0 || done == total {
			log.Printf("seeded power card %d/%d (%s)", done, total, res.p.Name)
		}
	}

	finalCount, _ := w.CountPowerCards()
	log.Printf("power card catalog ready: %d entries", finalCount)
	if finalCount == 0 && firstErr != nil {
		return firstErr
	}
	return nil
}
