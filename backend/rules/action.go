package rules

import (
	"fmt"
	"strings"

	"rakshyak-98/pokemon-backend/models"
)

// Action names mirrored from command package (avoid import cycle).
const (
	ActionStartGame        = "start_game"
	ActionSelectParty      = "select_party"
	ActionDrawCard         = "draw_card"
	ActionSelectDraw       = "select_draw"
	ActionPlayBench        = "play_bench"
	ActionSetActive        = "set_active"
	ActionAttachEnergy     = "attach_energy"
	ActionAttack           = "attack"
	ActionEndTurn          = "end_turn"
	ActionPromote          = "promote"
	ActionUseShield        = "use_shield"
)

// ValidateAction checks whether an action is legal for the current handbook-adapted game state.
func ValidateAction(state *models.GameState, playerID, action string, payload map[string]any) error {
	if state == nil {
		return fmt.Errorf("no active game")
	}

	switch action {
	case ActionStartGame:
		return nil

	case ActionSelectParty:
		return validateSelectParty(state, playerID, payload)

	case ActionDrawCard, ActionSelectDraw, ActionPlayBench:
		return fmt.Errorf("%s is not used in Pokémon GO tournament battles (handbook §6)", action)

	case ActionAttachEnergy:
		return validateChargeEnergy(state, playerID)

	case ActionSetActive:
		return validateSetActive(state, playerID)

	case ActionAttack:
		return validateAttack(state, playerID)

	case ActionEndTurn:
		return validateEndTurn(state, playerID)

	case ActionPromote:
		return validatePromote(state, playerID)

	case ActionUseShield:
		return validateUseShield(state, playerID)

	default:
		return fmt.Errorf("unknown action %q", action)
	}
}

func playerByID(state *models.GameState, id string) *models.PlayerState {
	for i := range state.Players {
		if state.Players[i].ID == id {
			return &state.Players[i]
		}
	}
	return nil
}

func requirePlayerTurn(state *models.GameState, playerID string) (*models.PlayerState, error) {
	if state.Status == models.StatusGameOver || state.Phase == string(PhaseMatchOver) {
		return nil, fmt.Errorf("match is over")
	}
	p := playerByID(state, playerID)
	if p == nil {
		return nil, fmt.Errorf("player not found")
	}
	if state.CurrentTurn != "" && state.CurrentTurn != playerID {
		return nil, fmt.Errorf("not your turn")
	}
	return p, nil
}

func validateSelectParty(state *models.GameState, playerID string, payload map[string]any) error {
	phase := Phase(state.Phase)
	if phase != PhaseTeamPreview && phase != PhasePartySelect && phase != PhaseBetweenGames {
		return fmt.Errorf("can only select a battle party during team preview / between games (§6.1)")
	}
	p := playerByID(state, playerID)
	if p == nil {
		return fmt.Errorf("player not found")
	}
	if err := ValidateBattleTeam(CardsToTeamMembers(p.BattleTeam)); err != nil {
		return err
	}

	ids, err := payloadIntSlice(payload, "pokeApiIds")
	if err != nil {
		// also accept cardIds → resolve via team
		cardIDs, err2 := payloadStringSlice(payload, "cardIds")
		if err2 != nil {
			return fmt.Errorf("select_party requires pokeApiIds or cardIds")
		}
		ids = make([]int, 0, len(cardIDs))
		for _, cid := range cardIDs {
			found := false
			for _, c := range p.BattleTeam {
				if c.ID == cid {
					ids = append(ids, c.PokeAPIID)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("card %s is not on your Battle Team", cid)
			}
		}
	}
	return ValidateBattleParty(CardsToTeamMembers(p.BattleTeam), ids)
}

func validateChargeEnergy(state *models.GameState, playerID string) error {
	if Phase(state.Phase) != PhaseInBattle {
		return fmt.Errorf("charge energy only during battle")
	}
	p, err := requirePlayerTurn(state, playerID)
	if err != nil {
		return err
	}
	if p.ActivePokemon == nil {
		return fmt.Errorf("no active Pokémon to charge")
	}
	if p.HasAttached {
		return fmt.Errorf("already charged energy this turn")
	}
	return nil
}

func validateSetActive(state *models.GameState, playerID string) error {
	if Phase(state.Phase) != PhaseInBattle {
		return fmt.Errorf("set active only during battle")
	}
	_, err := requirePlayerTurn(state, playerID)
	return err
}

func validateAttack(state *models.GameState, playerID string) error {
	if Phase(state.Phase) != PhaseInBattle {
		return fmt.Errorf("attack only during battle (§6)")
	}
	p, err := requirePlayerTurn(state, playerID)
	if err != nil {
		return err
	}
	if p.ActivePokemon == nil {
		return fmt.Errorf("no active Pokémon to attack with")
	}
	opp := opponentOf(state, playerID)
	if opp == nil || opp.ActivePokemon == nil {
		return fmt.Errorf("opponent has no active Pokémon")
	}
	return nil
}

func validateEndTurn(state *models.GameState, playerID string) error {
	if Phase(state.Phase) != PhaseInBattle {
		return fmt.Errorf("end turn only during battle")
	}
	_, err := requirePlayerTurn(state, playerID)
	return err
}

func validatePromote(state *models.GameState, playerID string) error {
	if Phase(state.Phase) != PhaseInBattle {
		return fmt.Errorf("promote only during battle")
	}
	p := playerByID(state, playerID)
	if p == nil {
		return fmt.Errorf("player not found")
	}
	if p.ActivePokemon != nil {
		return fmt.Errorf("active Pokémon still in play")
	}
	if len(p.BenchedPokemon) == 0 {
		return fmt.Errorf("no benched Pokémon to promote")
	}
	return nil
}

func validateUseShield(state *models.GameState, playerID string) error {
	if Phase(state.Phase) != PhaseInBattle {
		return fmt.Errorf("shields only during battle")
	}
	p, err := requirePlayerTurn(state, playerID)
	if err != nil {
		return err
	}
	if p.ProtectShields <= 0 {
		return fmt.Errorf("no Protect Shields remaining (§6.5.1)")
	}
	return nil
}

func opponentOf(state *models.GameState, playerID string) *models.PlayerState {
	for i := range state.Players {
		if state.Players[i].ID != playerID {
			return &state.Players[i]
		}
	}
	return nil
}

func payloadIntSlice(payload map[string]any, key string) ([]int, error) {
	if payload == nil {
		return nil, fmt.Errorf("missing %s", key)
	}
	raw, ok := payload[key]
	if !ok {
		return nil, fmt.Errorf("missing %s", key)
	}
	switch v := raw.(type) {
	case []int:
		return v, nil
	case []any:
		out := make([]int, 0, len(v))
		for _, item := range v {
			switch n := item.(type) {
			case float64:
				out = append(out, int(n))
			case int:
				out = append(out, n)
			default:
				return nil, fmt.Errorf("invalid %s entry", key)
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid %s", key)
	}
}

func payloadStringSlice(payload map[string]any, key string) ([]string, error) {
	if payload == nil {
		return nil, fmt.Errorf("missing %s", key)
	}
	raw, ok := payload[key]
	if !ok {
		return nil, fmt.Errorf("missing %s", key)
	}
	switch v := raw.(type) {
	case []string:
		return v, nil
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("invalid %s entry", key)
			}
			out = append(out, s)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid %s", key)
	}
}

// MatchWinnerID returns the player who has reached the wins needed for the match format (§6.4).
func MatchWinnerID(state *models.GameState) string {
	need := state.WinsNeeded
	if need <= 0 {
		need = DefaultMatchFormatWins
	}
	for _, p := range state.Players {
		if p.GamesWon >= need {
			return p.ID
		}
	}
	return ""
}

// PublicTeamPreview builds opponent-visible team preview data (§3.2.2.1 / §6.5.1).
func PublicTeamPreview(cards []models.Card) []map[string]any {
	out := make([]map[string]any, 0, len(cards))
	for _, c := range cards {
		moves := make([]string, 0, len(c.Attacks))
		for _, a := range c.Attacks {
			moves = append(moves, a.Name)
		}
		out = append(out, map[string]any{
			"name":        c.Name,
			"pokeApiId":   c.PokeAPIID,
			"elementType": c.ElementType,
			"combatPower": c.CombatPower,
			"bestBuddy":   c.BestBuddy,
			"moves":       moves,
			"types":       strings.TrimSpace(c.ElementType),
		})
	}
	return out
}
