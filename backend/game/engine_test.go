package game

import (
	"testing"
)

func TestGOMatchFlow(t *testing.T) {
	e := NewEngine(nil)
	if err := e.StartGame("player1", "player2"); err != nil {
		t.Fatal(err)
	}
	if e.State.Phase != "TeamPreview" {
		t.Fatalf("phase %s", e.State.Phase)
	}
	p1 := e.getPlayer("player1")
	p2 := e.getPlayer("player2")
	if len(p1.BattleTeam) < 3 || len(p2.BattleTeam) < 3 {
		t.Fatalf("teams too small: %d / %d", len(p1.BattleTeam), len(p2.BattleTeam))
	}

	ids1 := []string{p1.BattleTeam[0].ID, p1.BattleTeam[1].ID, p1.BattleTeam[2].ID}
	ids2 := []string{p2.BattleTeam[0].ID, p2.BattleTeam[1].ID, p2.BattleTeam[2].ID}
	if err := e.SelectParty("player1", ids1); err != nil {
		t.Fatal(err)
	}
	if err := e.SelectParty("player2", ids2); err != nil {
		t.Fatal(err)
	}
	if e.State.Phase != "InBattle" {
		t.Fatalf("expected InBattle, got %s", e.State.Phase)
	}
	if p1.ActivePokemon == nil || p2.ActivePokemon == nil {
		t.Fatal("missing actives")
	}

	if err := e.AttachEnergy("player1", ""); err != nil {
		t.Fatal(err)
	}
	// Ensure enough energy for first attack
	for p1.ActivePokemon.EnergyAttached < p1.ActivePokemon.Attacks[0].Cost {
		p1.HasAttached = false
		_ = e.AttachEnergy("player1", "")
	}
	if err := e.Attack("player1", 0); err != nil {
		t.Fatal(err)
	}
}
