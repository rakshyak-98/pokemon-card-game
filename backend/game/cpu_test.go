package game

import "testing"

func TestPracticeVsCPU(t *testing.T) {
	e := NewEngine(nil)
	if err := e.StartGame("player1", "player2", true); err != nil {
		t.Fatal(err)
	}
	if !e.State.VsCPU || e.State.CPUPlayerID != "player2" {
		t.Fatalf("cpu flags: vs=%v id=%s", e.State.VsCPU, e.State.CPUPlayerID)
	}

	// CPU should lock a party once asked to act.
	e.RunCPUIfNeeded()
	cpu := e.getPlayer("player2")
	if !cpu.PartyReady || cpu.ActivePokemon == nil {
		t.Fatal("CPU should have selected a party")
	}

	p1 := e.getPlayer("player1")
	ids := []string{p1.BattleTeam[0].ID, p1.BattleTeam[1].ID, p1.BattleTeam[2].ID}
	if err := e.SelectParty("player1", ids); err != nil {
		t.Fatal(err)
	}
	if e.State.Phase != "InBattle" {
		t.Fatalf("phase %s", e.State.Phase)
	}

	// Human turn: charge + end; CPU should respond.
	if err := e.AttachEnergy("player1", ""); err != nil {
		t.Fatal(err)
	}
	if err := e.EndTurn("player1"); err != nil {
		t.Fatal(err)
	}
	e.RunCPUIfNeeded()
	if e.State.CurrentTurn != "player1" && e.State.Phase == "InBattle" {
		// CPU may have ended on human, or game advanced; either is fine if match continues.
		if e.State.CurrentTurn == "player2" {
			t.Fatal("CPU should have finished its turn")
		}
	}
}
