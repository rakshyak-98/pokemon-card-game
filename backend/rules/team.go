package rules

import (
	"fmt"
	"math"
	"strings"

	"rakshyak-98/pokemon-backend/models"
)

// TeamMember is the minimum data needed to validate a Battle Team entry (§3.2.2).
type TeamMember struct {
	PokeAPIID   int
	Name        string
	CP          int
	BestBuddy   bool
	MegaOrPrimal bool
	Attacks     []models.Attack
}

// EstimateCombatPower derives a Great-League-style CP from PokeAPI base stats.
// Competitors power down to ≤1500 for Great League; callers should clamp with ClampCP.
func EstimateCombatPower(stats models.PokemonStats) int {
	atk := float64(maxInt(stats.Attack, stats.SpAttack))
	def := float64(maxInt(stats.Defense, stats.SpDefense))
	sta := float64(stats.HP)
	if atk < 1 {
		atk = 1
	}
	if def < 1 {
		def = 1
	}
	if sta < 1 {
		sta = 1
	}
	cp := int(math.Floor((atk * math.Sqrt(def) * math.Sqrt(sta)) / 10))
	if cp < 10 {
		return 10
	}
	return cp
}

// ClampCP enforces Great League cap (§3.1).
func ClampCP(cp int) int {
	if cp < 10 {
		return 10
	}
	if cp > GreatLeagueCPCap {
		return GreatLeagueCPCap
	}
	return cp
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// ValidateBattleTeam checks §3.1 / §3.1.2 team construction rules.
func ValidateBattleTeam(members []TeamMember) error {
	n := len(members)
	if n < MinBattleTeamSize {
		return fmt.Errorf("battle team must have at least %d Pokémon (§3.1)", MinBattleTeamSize)
	}
	if n > MaxBattleTeamSize {
		return fmt.Errorf("battle team may have at most %d Pokémon (§3.1)", MaxBattleTeamSize)
	}

	seenDex := map[int]struct{}{}
	bestBuddyCount := 0

	for i, m := range members {
		if m.PokeAPIID <= 0 && m.Name == "" {
			return fmt.Errorf("team slot %d is missing species identity (§3.2.2)", i+1)
		}
		if _, banned := HardBannedDexIDs[m.PokeAPIID]; banned {
			return fmt.Errorf("%s is banned from tournament play (§3.1.1)", displayName(m))
		}
		if _, banned := HardBannedSpecies[normalizeName(m.Name)]; banned {
			return fmt.Errorf("%s is banned from tournament play (§3.1.1)", displayName(m))
		}
		if m.MegaOrPrimal {
			return fmt.Errorf("Mega Evolved / Primal Pokémon are not allowed (§3.1.2)")
		}
		if m.CP <= 0 || m.CP > GreatLeagueCPCap {
			return fmt.Errorf("%s CP must be 1–%d for Great League (§3.1)", displayName(m), GreatLeagueCPCap)
		}
		if m.PokeAPIID > 0 {
			if _, dup := seenDex[m.PokeAPIID]; dup {
				return fmt.Errorf("team cannot contain two Pokémon with the same National Pokédex number (§3.1.2)")
			}
			seenDex[m.PokeAPIID] = struct{}{}
		}
		if m.BestBuddy {
			bestBuddyCount++
		}
		if len(m.Attacks) == 0 {
			return fmt.Errorf("%s must list known attacks (§3.2.2)", displayName(m))
		}
	}

	if bestBuddyCount > MaxBestBuddyBoosts {
		return fmt.Errorf("team may contain at most %d Best Buddy CP boost (§3.1.2)", MaxBestBuddyBoosts)
	}
	return nil
}

func displayName(m TeamMember) string {
	if m.Name != "" {
		return m.Name
	}
	return fmt.Sprintf("#%d", m.PokeAPIID)
}

// ValidateBattleParty ensures exactly three team members are selected for a game (§3.1).
func ValidateBattleParty(team []TeamMember, partyIDs []int) error {
	if len(partyIDs) != BattlePartySize {
		return fmt.Errorf("must select exactly %d Pokémon for battle (§3.1)", BattlePartySize)
	}
	seen := map[int]struct{}{}
	teamIndex := map[int]TeamMember{}
	for _, m := range team {
		teamIndex[m.PokeAPIID] = m
	}
	for _, id := range partyIDs {
		if _, ok := seen[id]; ok {
			return fmt.Errorf("battle party cannot repeat the same Pokémon")
		}
		seen[id] = struct{}{}
		if _, ok := teamIndex[id]; !ok {
			return fmt.Errorf("battle party Pokémon must come from the registered Battle Team (§3.1)")
		}
	}
	return nil
}

// CardsToTeamMembers adapts game cards for team validation.
func CardsToTeamMembers(cards []models.Card) []TeamMember {
	out := make([]TeamMember, 0, len(cards))
	for _, c := range cards {
		out = append(out, TeamMember{
			PokeAPIID: c.PokeAPIID,
			Name:      c.Name,
			CP:        c.CombatPower,
			BestBuddy: c.BestBuddy,
			Attacks:   c.Attacks,
		})
	}
	return out
}
