package game

import (
	"testing"

	"rakshyak-98/pokemon-backend/models"
)

func TestBuildPowerDeck(t *testing.T) {
	deck := buildPowerDeck("player1")
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

	for i := 0; i < MaxPowerHandSlots; i++ {
		p1 = e.getPlayer("player1")
		p1.HasDrawn = false
		if err := e.DrawCard("player1"); err != nil {
			t.Fatalf("draw %d: %v", i+1, err)
		}
		// Advance and come back so draw flag resets without needing full attack loop.
		_ = e.advanceTurn("player1")
		_ = e.advanceTurn("player2")
	}
	p1 = e.getPlayer("player1")
	if len(p1.Hand) != MaxPowerHandSlots {
		t.Fatalf("expected %d cards in hand, got %d", MaxPowerHandSlots, len(p1.Hand))
	}
	p1.HasDrawn = false
	if err := e.DrawCard("player1"); err == nil {
		t.Fatal("expected draw to fail when hand is full")
	}
	// Keeping cards across turns: hand unchanged if we skip PlayPower.
	if len(p1.Hand) != MaxPowerHandSlots {
		t.Fatalf("kept cards should remain for next turn")
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
	before := len(p1.PowerDeck)
	if err := e.DrawCard("player1"); err != nil {
		t.Fatal(err)
	}
	p1 = e.getPlayer("player1")
	if len(p1.Hand) != 1 || len(p1.PowerDeck) != before-1 {
		t.Fatalf("draw should move one card to hand")
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
