package game

import (
	"testing"
)

func TestFullTurnLoop(t *testing.T) {
	e := NewEngine(nil)
	if err := e.StartGame("player1", "player2"); err != nil {
		t.Fatal(err)
	}
	p1 := e.getPlayer("player1")
	if len(p1.Hand) != 7 {
		t.Fatalf("hand size %d", len(p1.Hand))
	}
	if len(p1.PrizeCards) != 6 {
		t.Fatalf("prizes %d", len(p1.PrizeCards))
	}

	// Find a pokemon in hand to set active
	var pokeID, energyID string
	for _, c := range p1.Hand {
		if c.Type == "Pokemon" && pokeID == "" {
			pokeID = c.ID
		}
		if c.Type == "Energy" && energyID == "" {
			energyID = c.ID
		}
	}
	if pokeID == "" {
		// draw until we get one
		for i := 0; i < 10 && pokeID == ""; i++ {
			_ = e.DrawCard("player1")
			p1 = e.getPlayer("player1")
			for _, c := range p1.Hand {
				if c.Type == "Pokemon" {
					pokeID = c.ID
					break
				}
			}
			p1.HasDrawn = false
		}
	}
	if pokeID == "" {
		t.Fatal("no pokemon in hand")
	}
	if err := e.SetActive("player1", pokeID); err != nil {
		t.Fatal(err)
	}

	// Ensure energy available
	p1 = e.getPlayer("player1")
	for _, c := range p1.Hand {
		if c.Type == "Energy" {
			energyID = c.ID
			break
		}
	}
	if energyID == "" {
		t.Skip("no energy in opening hand")
	}
	if err := e.AttachEnergy("player1", energyID); err != nil {
		t.Fatal(err)
	}

	// Opponent needs active too
	p2 := e.getPlayer("player2")
	var oppPoke string
	for _, c := range p2.Hand {
		if c.Type == "Pokemon" {
			oppPoke = c.ID
			break
		}
	}
	e.State.CurrentTurn = "player2"
	if err := e.SetActive("player2", oppPoke); err != nil {
		t.Fatal(err)
	}
	e.State.CurrentTurn = "player1"

	if err := e.Attack("player1", 0); err != nil {
		t.Fatal(err)
	}
	if e.State.CurrentTurn != "player2" {
		t.Fatalf("expected turn to pass, got %s", e.State.CurrentTurn)
	}
}
