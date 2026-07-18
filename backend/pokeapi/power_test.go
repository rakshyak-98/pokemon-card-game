package pokeapi

import (
	"testing"

	"rakshyak-98/pokemon-backend/models"
)

func TestMapItemToPowerCard(t *testing.T) {
	raw := apiItem{
		ID:   57,
		Name: "x-attack",
		Category: apiNamedResource{
			Name: "stat-boosts",
		},
		EffectEntries: []struct {
			Effect      string `json:"effect"`
			ShortEffect string `json:"short_effect"`
			Language    struct {
				Name string `json:"name"`
			} `json:"language"`
		}{
			{
				ShortEffect: "Raises Attack by one stage.",
				Language: struct {
					Name string `json:"name"`
				}{Name: "en"},
			},
		},
		Sprites: struct {
			Default string `json:"default"`
		}{Default: "https://example/x-attack.png"},
	}

	card, ok := MapItemToPowerCard(raw)
	if !ok {
		t.Fatal("expected x-attack to map")
	}
	if card.Effect != "boost_attack" || card.EffectValue != 20 {
		t.Fatalf("unexpected mapping: %+v", card)
	}
	if card.Name != "X Attack" {
		t.Fatalf("expected titled name, got %q", card.Name)
	}
}

func TestMapItemToPowerCardUnknown(t *testing.T) {
	_, ok := MapItemToPowerCard(apiItem{Name: "master-ball"})
	if ok {
		t.Fatal("master-ball should not map to a power card")
	}
}

type memPowerWriter struct {
	cards map[int]models.PowerCard
}

func (m *memPowerWriter) CountPowerCards() (int, error) {
	return len(m.cards), nil
}

func (m *memPowerWriter) UpsertPowerCard(p models.PowerCard) error {
	if m.cards == nil {
		m.cards = map[int]models.PowerCard{}
	}
	m.cards[p.PokeAPIID] = p
	return nil
}

func TestSeedPowerIfEmptySkipsWhenFull(t *testing.T) {
	w := &memPowerWriter{cards: map[int]models.PowerCard{}}
	for i := 0; i < len(powerItemSpecs); i++ {
		w.cards[i+1] = models.PowerCard{PokeAPIID: i + 1, Name: "x", Effect: "heal", EffectValue: 1}
	}
	if err := SeedPowerIfEmpty(w, NewClient(), PowerSeedOptions{}); err != nil {
		t.Fatal(err)
	}
	if len(w.cards) != len(powerItemSpecs) {
		t.Fatalf("should not mutate a complete catalog")
	}
}
