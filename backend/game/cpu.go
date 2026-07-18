package game

import (
	"sort"

	"rakshyak-98/pokemon-backend/models"
	"rakshyak-98/pokemon-backend/rules"
)

const cpuMaxSteps = 12

// RunCPUIfNeeded advances the practice AI until it is no longer the CPU's turn
// to act (party select or battle), or the match ends.
func (e *Engine) RunCPUIfNeeded() {
	if e.State == nil || !e.State.VsCPU || e.State.CPUPlayerID == "" {
		return
	}
	for i := 0; i < cpuMaxSteps; i++ {
		if !e.cpuShouldAct() {
			break
		}
		if err := e.cpuActOnce(); err != nil {
			e.State.LastAction = "CPU pause: " + err.Error()
			break
		}
	}
	// Ensure the player who receives the turn after the CPU still gets their power draw.
	if e.State.CurrentTurn != "" && e.State.CurrentTurn != e.State.CPUPlayerID {
		e.tryAutoDrawPower(e.State.CurrentTurn)
	}
}

func (e *Engine) cpuShouldAct() bool {
	if e.State.Status == models.StatusGameOver || e.State.Phase == string(rules.PhaseMatchOver) {
		return false
	}
	cpu := e.getPlayer(e.State.CPUPlayerID)
	if cpu == nil {
		return false
	}

	phase := rules.Phase(e.State.Phase)
	switch phase {
	case rules.PhaseTeamPreview, rules.PhasePartySelect, rules.PhaseBetweenGames:
		return !cpu.PartyReady
	case rules.PhaseInBattle:
		if cpu.ActivePokemon == nil && len(cpu.BenchedPokemon) > 0 {
			return true
		}
		return e.State.CurrentTurn == e.State.CPUPlayerID
	default:
		return false
	}
}

func (e *Engine) cpuActOnce() error {
	cpuID := e.State.CPUPlayerID
	cpu := e.getPlayer(cpuID)
	if cpu == nil {
		return nil
	}

	phase := rules.Phase(e.State.Phase)
	if phase == rules.PhaseTeamPreview || phase == rules.PhasePartySelect || phase == rules.PhaseBetweenGames {
		if cpu.PartyReady {
			return nil
		}
		ids := e.cpuPickParty(cpu)
		return e.SelectParty(cpuID, ids)
	}

	// Must promote after KO before anything else.
	if cpu.ActivePokemon == nil && len(cpu.BenchedPokemon) > 0 {
		return e.Promote(cpuID, cpu.BenchedPokemon[0].ID)
	}

	if e.State.CurrentTurn != cpuID {
		return nil
	}

	// Resolve a full-hand power draw before other actions.
	if len(cpu.PendingDraw) > 0 {
		replaceID := cpuPickReplace(cpu)
		return e.SelectDraw(cpuID, replaceID)
	}

	// Charge once if possible.
	if !cpu.HasAttached && cpu.ActivePokemon != nil {
		if err := e.AttachEnergy(cpuID, ""); err != nil {
			return err
		}
		cpu = e.getPlayer(cpuID)
	}

	// Draw is auto-triggered on turn start; only draw here if somehow skipped.
	if cpu != nil && !cpu.HasDrawn && len(cpu.PowerDeck) > 0 {
		if err := e.DrawCard(cpuID); err != nil {
			return err
		}
		cpu = e.getPlayer(cpuID)
		if cpu != nil && len(cpu.PendingDraw) > 0 {
			return e.SelectDraw(cpuID, cpuPickReplace(cpu))
		}
	}

	// Play beneficial power cards from slots onto the active Pokémon (any number per turn).
	for cpu != nil && len(cpu.Hand) > 0 && cpu.ActivePokemon != nil {
		cardID := cpuPickPower(cpu)
		if cardID == "" {
			break
		}
		if err := e.PlayPower(cpuID, cardID); err != nil {
			return err
		}
		cpu = e.getPlayer(cpuID)
	}

	// Attack with the strongest affordable move if opponent is active.
	opp := e.getOpponent(cpuID)
	if cpu.ActivePokemon != nil && opp != nil && opp.ActivePokemon != nil {
		bestIdx := -1
		bestDmg := -1
		for i, atk := range cpu.ActivePokemon.Attacks {
			if cpu.ActivePokemon.EnergyAttached >= atk.Cost && atk.Damage > bestDmg {
				bestDmg = atk.Damage
				bestIdx = i
			}
		}
		if bestIdx >= 0 {
			return e.Attack(cpuID, bestIdx)
		}
	}

	return e.EndTurn(cpuID)
}

func (e *Engine) cpuPickParty(cpu *models.PlayerState) []string {
	team := append([]models.Card(nil), cpu.BattleTeam...)
	sort.SliceStable(team, func(i, j int) bool {
		// Prefer higher HP, then CP, for a simple practice opponent.
		if team[i].MaxHP != team[j].MaxHP {
			return team[i].MaxHP > team[j].MaxHP
		}
		return team[i].CombatPower > team[j].CombatPower
	})
	n := rules.BattlePartySize
	if len(team) < n {
		n = len(team)
	}
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, team[i].ID)
	}
	return ids
}

// cpuPickPower chooses a beneficial power card from hand, or "" to skip.
func cpuPickPower(cpu *models.PlayerState) string {
	var healID, attackID, defenseID string
	for _, c := range cpu.Hand {
		if c.Type != models.TypePower {
			continue
		}
		switch c.Effect {
		case EffectHeal:
			if cpu.ActivePokemon != nil && cpu.ActivePokemon.HP < cpu.ActivePokemon.MaxHP {
				healID = c.ID
			}
		case EffectBoostAttack:
			attackID = c.ID
		case EffectBoostDefense:
			defenseID = c.ID
		}
	}
	// Prefer heal when below half HP, otherwise lean into offense.
	if healID != "" && cpu.ActivePokemon != nil && cpu.ActivePokemon.HP*2 <= cpu.ActivePokemon.MaxHP {
		return healID
	}
	if attackID != "" {
		return attackID
	}
	if healID != "" {
		return healID
	}
	return defenseID
}

// cpuPickReplace chooses a hand card to discard for the pending draw, or "_keep".
func cpuPickReplace(cpu *models.PlayerState) string {
	if len(cpu.PendingDraw) == 0 || len(cpu.Hand) == 0 {
		return "_keep"
	}
	pending := cpu.PendingDraw[0]
	pendingRank := powerRank(pending.Effect, cpu)

	worstID := ""
	worstRank := 999
	for _, c := range cpu.Hand {
		r := powerRank(c.Effect, cpu)
		if r < worstRank {
			worstRank = r
			worstID = c.ID
		}
	}
	if pendingRank > worstRank {
		return worstID
	}
	return "_keep"
}

func powerRank(effect string, cpu *models.PlayerState) int {
	switch effect {
	case EffectHeal:
		if cpu.ActivePokemon != nil && cpu.ActivePokemon.HP*2 <= cpu.ActivePokemon.MaxHP {
			return 3
		}
		if cpu.ActivePokemon != nil && cpu.ActivePokemon.HP < cpu.ActivePokemon.MaxHP {
			return 2
		}
		return 0
	case EffectBoostAttack:
		return 2
	case EffectBoostDefense:
		return 1
	default:
		return 0
	}
}
