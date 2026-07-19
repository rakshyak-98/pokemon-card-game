package game

import (
	"testing"

	"rakshyak-98/pokemon-backend/models"
)

func TestBuildPowerDeck(t *testing.T) {
	e := NewEngine(nil)
	deck := e.buildPowerDeck("player1")
	if len(deck) != PowerDeckSize {
		t.Fatalf("expected deck size %d, got %d", PowerDeckSize, len(deck))
	}
	counts := map[string]int{}
	for _, c := range deck {
		if c.Type != models.TypePower {
			t.Fatalf("expected TypePower, got %s", c.Type)
		}
		counts[c.Effect]++
	}
	if counts[EffectBoostAttack] == 0 || counts[EffectBoostDefense] == 0 || counts[EffectHeal] == 0 {
		t.Fatalf("deck missing effect kinds: %+v", counts)
	}
}

func TestBuildPowerDeckUsesCatalog(t *testing.T) {
	e := NewEngine(nil)
	e.SetPowerCatalog([]models.PowerCard{
		{PokeAPIID: 1, Name: "Hyper Potion", Effect: EffectHeal, EffectValue: 50, ImageURL: "https://example/heal.png"},
		{PokeAPIID: 2, Name: "X Attack 3", Effect: EffectBoostAttack, EffectValue: 30, ImageURL: "https://example/atk.png"},
		{PokeAPIID: 3, Name: "X Defense 3", Effect: EffectBoostDefense, EffectValue: 25, ImageURL: "https://example/def.png"},
	})
	deck := e.buildPowerDeck("p1")
	if len(deck) != PowerDeckSize {
		t.Fatalf("expected deck size %d, got %d", PowerDeckSize, len(deck))
	}
	names := map[string]bool{}
	for _, c := range deck {
		names[c.Name] = true
		if c.PokeAPIID == 0 {
			t.Fatalf("expected pokeApiId from catalog on %s", c.Name)
		}
	}
	if !names["Hyper Potion"] || !names["X Attack 3"] || !names["X Defense 3"] {
		t.Fatalf("deck should sample seeded catalog names, got %+v", names)
	}
}

func TestDrawRespectsMaxHandSlots(t *testing.T) {
	e := NewEngine(fallbackCatalog())
	if err := e.StartGame("player1", "player2", false); err != nil {
		t.Fatal(err)
	}
	p1 := e.getPlayer("player1")
	p2 := e.getPlayer("player2")
	ids1 := []string{p1.BattleTeam[0].ID, p1.BattleTeam[1].ID, p1.BattleTeam[2].ID}
	ids2 := []string{p2.BattleTeam[0].ID, p2.BattleTeam[1].ID, p2.BattleTeam[2].ID}
	if err := e.SelectParty("player1", ids1); err != nil {
		t.Fatal(err)
	}
	if err := e.SelectParty("player2", ids2); err != nil {
		t.Fatal(err)
	}

	// beginBattle auto-draws once for the starting player.
	p1 = e.getPlayer("player1")
	if !p1.HasDrawn || len(p1.Hand) != 1 {
		t.Fatalf("expected auto-draw on battle start, hand=%d hasDrawn=%v", len(p1.Hand), p1.HasDrawn)
	}

	for len(p1.Hand) < MaxPowerHandSlots {
		p1.HasDrawn = false
		if err := e.DrawCard("player1"); err != nil {
			t.Fatalf("fill draw: %v", err)
		}
		p1 = e.getPlayer("player1")
	}
	if len(p1.Hand) != MaxPowerHandSlots {
		t.Fatalf("expected %d cards in hand, got %d", MaxPowerHandSlots, len(p1.Hand))
	}

	p1.HasDrawn = false
	beforeDeck := len(p1.PowerDeck)
	if err := e.DrawCard("player1"); err != nil {
		t.Fatalf("full-hand draw should go pending: %v", err)
	}
	p1 = e.getPlayer("player1")
	if len(p1.PendingDraw) != 1 {
		t.Fatalf("expected pending draw when hand full, got %d", len(p1.PendingDraw))
	}
	if len(p1.Hand) != MaxPowerHandSlots {
		t.Fatalf("hand should stay full until replace, got %d", len(p1.Hand))
	}
	if len(p1.PowerDeck) != beforeDeck-1 {
		t.Fatalf("deck should lose the pending card")
	}

	replaceID := p1.Hand[0].ID
	pendingName := p1.PendingDraw[0].Name
	if err := e.SelectDraw("player1", replaceID); err != nil {
		t.Fatal(err)
	}
	p1 = e.getPlayer("player1")
	if len(p1.PendingDraw) != 0 {
		t.Fatal("pending should clear after replace")
	}
	if len(p1.Hand) != MaxPowerHandSlots {
		t.Fatalf("hand should remain full after replace")
	}
	if p1.Hand[0].Name != pendingName {
		t.Fatalf("expected slot 0 to become %s, got %s", pendingName, p1.Hand[0].Name)
	}
}

func TestDrawAndPlayPower(t *testing.T) {
	e := NewEngine(fallbackCatalog())
	if err := e.StartGame("player1", "player2", false); err != nil {
		t.Fatal(err)
	}
	p1 := e.getPlayer("player1")
	p2 := e.getPlayer("player2")
	if len(p1.PowerDeck) != PowerDeckSize || len(p2.PowerDeck) != PowerDeckSize {
		t.Fatalf("both players should start with power decks")
	}

	// Lock parties so battle begins.
	ids1 := []string{p1.BattleTeam[0].ID, p1.BattleTeam[1].ID, p1.BattleTeam[2].ID}
	ids2 := []string{p2.BattleTeam[0].ID, p2.BattleTeam[1].ID, p2.BattleTeam[2].ID}
	if err := e.SelectParty("player1", ids1); err != nil {
		t.Fatal(err)
	}
	if err := e.SelectParty("player2", ids2); err != nil {
		t.Fatal(err)
	}

	p1 = e.getPlayer("player1")
	if len(p1.Hand) != 1 || !p1.HasDrawn {
		t.Fatalf("battle start should auto-draw one power card")
	}

	card := p1.Hand[0]
	// Force a heal-friendly state if needed.
	if card.Effect == EffectHeal {
		p1.ActivePokemon.HP = p1.ActivePokemon.MaxHP - 10
	}
	if err := e.PlayPower("player1", card.ID); err != nil {
		t.Fatal(err)
	}
	p1 = e.getPlayer("player1")
	if len(p1.Hand) != 0 || !p1.HasPlayedPower {
		t.Fatalf("power card should leave hand and set HasPlayedPower")
	}
	switch card.Effect {
	case EffectBoostAttack:
		if p1.AttackBonus != card.EffectValue {
			t.Fatalf("expected attack bonus %d, got %d", card.EffectValue, p1.AttackBonus)
		}
	case EffectBoostDefense:
		if p1.DefenseBonus != card.EffectValue {
			t.Fatalf("expected defense bonus %d, got %d", card.EffectValue, p1.DefenseBonus)
		}
	case EffectHeal:
		if p1.ActivePokemon.HP <= p1.ActivePokemon.MaxHP-10 {
			t.Fatalf("heal should restore HP")
		}
	}
}

func TestPlayMultiplePowersSameTurn(t *testing.T) {
	e := NewEngine(fallbackCatalog())
	if err := e.StartGame("player1", "player2", false); err != nil {
		t.Fatal(err)
	}
	p1 := e.getPlayer("player1")
	ids1 := []string{p1.BattleTeam[0].ID, p1.BattleTeam[1].ID, p1.BattleTeam[2].ID}
	p2 := e.getPlayer("player2")
	ids2 := []string{p2.BattleTeam[0].ID, p2.BattleTeam[1].ID, p2.BattleTeam[2].ID}
	if err := e.SelectParty("player1", ids1); err != nil {
		t.Fatal(err)
	}
	if err := e.SelectParty("player2", ids2); err != nil {
		t.Fatal(err)
	}

	p1 = e.getPlayer("player1")
	// Seed two boost cards into hand (beyond the auto-draw).
	p1.Hand = []models.Card{
		{ID: "atk-1", Name: "Power Strike", Type: models.TypePower, Effect: EffectBoostAttack, EffectValue: PowerAttackBonus},
		{ID: "def-1", Name: "Iron Guard", Type: models.TypePower, Effect: EffectBoostDefense, EffectValue: PowerDefenseBonus},
	}

	if err := e.PlayPower("player1", "atk-1"); err != nil {
		t.Fatal(err)
	}
	if err := e.PlayPower("player1", "def-1"); err != nil {
		t.Fatalf("should allow a second power card same turn: %v", err)
	}
	p1 = e.getPlayer("player1")
	if len(p1.Hand) != 0 {
		t.Fatalf("both power cards should leave hand, got %d", len(p1.Hand))
	}
	if p1.AttackBonus != PowerAttackBonus || p1.DefenseBonus != PowerDefenseBonus {
		t.Fatalf("expected stacked bonuses atk=%d def=%d, got atk=%d def=%d",
			PowerAttackBonus, PowerDefenseBonus, p1.AttackBonus, p1.DefenseBonus)
	}
}

func TestComputeDamageUsesBonuses(t *testing.T) {
	attacker := &models.PlayerState{AttackBonus: 20}
	defender := &models.PlayerState{DefenseBonus: 15}
	if got := computeDamage(30, attacker, defender); got != 35 {
		t.Fatalf("expected 35 damage, got %d", got)
	}
	if got := computeDamage(10, attacker, defender); got != 15 {
		t.Fatalf("expected 15 damage, got %d", got)
	}
	defender.DefenseBonus = 100
	if got := computeDamage(10, attacker, defender); got != 0 {
		t.Fatalf("damage should floor at 0, got %d", got)
	}
}

func TestPowerDeckRefillsUntilGameEnds(t *testing.T) {
	e := NewEngine(fallbackCatalog())
	if err := e.StartGame("player1", "player2", false); err != nil {
		t.Fatal(err)
	}
	p1 := e.getPlayer("player1")
	p2 := e.getPlayer("player2")
	ids1 := []string{p1.BattleTeam[0].ID, p1.BattleTeam[1].ID, p1.BattleTeam[2].ID}
	ids2 := []string{p2.BattleTeam[0].ID, p2.BattleTeam[1].ID, p2.BattleTeam[2].ID}
	if err := e.SelectParty("player1", ids1); err != nil {
		t.Fatal(err)
	}
	if err := e.SelectParty("player2", ids2); err != nil {
		t.Fatal(err)
	}

	p1 = e.getPlayer("player1")
	// Exhaust the starting deck (one card already auto-drawn into hand).
	p1.PowerDeck = nil
	p1.HasDrawn = false
	beforeHand := len(p1.Hand)

	if err := e.DrawCard("player1"); err != nil {
		t.Fatalf("draw should refill empty power deck: %v", err)
	}
	p1 = e.getPlayer("player1")
	if len(p1.Hand) != beforeHand+1 {
		t.Fatalf("expected hand to grow after refill draw, hand=%d", len(p1.Hand))
	}
	if len(p1.PowerDeck) != PowerDeckSize-1 {
		t.Fatalf("expected refilled deck with %d remaining, got %d", PowerDeckSize-1, len(p1.PowerDeck))
	}

	// Simulate many turns: auto-draw must keep working after further empties.
	for turn := 0; turn < PowerDeckSize+3; turn++ {
		p1 = e.getPlayer("player1")
		p1.HasDrawn = false
		p1.PendingDraw = nil
		if len(p1.Hand) >= MaxPowerHandSlots {
			p1.Hand = p1.Hand[:MaxPowerHandSlots-1]
		}
		e.tryAutoDrawPower("player1")
		p1 = e.getPlayer("player1")
		if !p1.HasDrawn {
			t.Fatalf("turn %d: expected auto-draw after refill", turn)
		}
	}
}

func TestFallbackCatalogHasASCTrainerVariety(t *testing.T) {
	names := map[string]bool{}
	effects := map[string]int{}
	for _, c := range fallbackPowerCatalog {
		names[c.Name] = true
		effects[c.Effect]++
	}
	for _, want := range []string{"Potion", "Hyper Potion", "X Attack", "Guard Spec", "Oran Berry", "Muscle Band"} {
		if !names[want] {
			t.Fatalf("fallback catalog missing ASC-themed card %q", want)
		}
	}
	if effects[EffectBoostAttack] == 0 || effects[EffectBoostDefense] == 0 || effects[EffectHeal] == 0 {
		t.Fatalf("fallback catalog missing an effect kind: %+v", effects)
	}
}
