package rules

import (
	"testing"

	"rakshyak-98/pokemon-backend/models"
)

func TestValidateBattleTeam_OK(t *testing.T) {
	team := []TeamMember{
		{PokeAPIID: 25, Name: "Pikachu", CP: 900, Attacks: []models.Attack{{Name: "Thunder Shock", Damage: 20, Cost: 1}}},
		{PokeAPIID: 4, Name: "Charmander", CP: 800, Attacks: []models.Attack{{Name: "Ember", Damage: 30, Cost: 1}}},
		{PokeAPIID: 7, Name: "Squirtle", CP: 850, Attacks: []models.Attack{{Name: "Water Gun", Damage: 20, Cost: 1}}},
	}
	if err := ValidateBattleTeam(team); err != nil {
		t.Fatal(err)
	}
}

func TestValidateBattleTeam_DuplicateDex(t *testing.T) {
	team := []TeamMember{
		{PokeAPIID: 25, Name: "Pikachu", CP: 900, Attacks: []models.Attack{{Name: "A"}}},
		{PokeAPIID: 25, Name: "Pikachu", CP: 800, Attacks: []models.Attack{{Name: "B"}}},
		{PokeAPIID: 7, Name: "Squirtle", CP: 850, Attacks: []models.Attack{{Name: "C"}}},
	}
	if err := ValidateBattleTeam(team); err == nil {
		t.Fatal("expected duplicate dex error")
	}
}

func TestValidateBattleTeam_Banned(t *testing.T) {
	team := []TeamMember{
		{PokeAPIID: 132, Name: "Ditto", CP: 500, Attacks: []models.Attack{{Name: "Transform"}}},
		{PokeAPIID: 4, Name: "Charmander", CP: 800, Attacks: []models.Attack{{Name: "Ember"}}},
		{PokeAPIID: 7, Name: "Squirtle", CP: 850, Attacks: []models.Attack{{Name: "Water Gun"}}},
	}
	if err := ValidateBattleTeam(team); err == nil {
		t.Fatal("expected ban error")
	}
}

func TestValidateBattleTeam_CPOver(t *testing.T) {
	team := []TeamMember{
		{PokeAPIID: 25, Name: "Pikachu", CP: 1501, Attacks: []models.Attack{{Name: "A"}}},
		{PokeAPIID: 4, Name: "Charmander", CP: 800, Attacks: []models.Attack{{Name: "B"}}},
		{PokeAPIID: 7, Name: "Squirtle", CP: 850, Attacks: []models.Attack{{Name: "C"}}},
	}
	if err := ValidateBattleTeam(team); err == nil {
		t.Fatal("expected CP error")
	}
}

func TestValidateBattleParty(t *testing.T) {
	team := []TeamMember{
		{PokeAPIID: 25, Name: "Pikachu", CP: 900, Attacks: []models.Attack{{Name: "A"}}},
		{PokeAPIID: 4, Name: "Charmander", CP: 800, Attacks: []models.Attack{{Name: "B"}}},
		{PokeAPIID: 7, Name: "Squirtle", CP: 850, Attacks: []models.Attack{{Name: "C"}}},
		{PokeAPIID: 1, Name: "Bulbasaur", CP: 700, Attacks: []models.Attack{{Name: "D"}}},
	}
	if err := ValidateBattleParty(team, []int{25, 4, 7}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateBattleParty(team, []int{25, 4}); err == nil {
		t.Fatal("expected size error")
	}
}

func TestClampCP(t *testing.T) {
	if ClampCP(2000) != GreatLeagueCPCap {
		t.Fatalf("got %d", ClampCP(2000))
	}
}
