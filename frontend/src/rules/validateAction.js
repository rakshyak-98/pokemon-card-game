import {
  BATTLE_PARTY_SIZE,
  GREAT_LEAGUE_CP_CAP,
  HARD_BANNED_DEX_IDS,
  HARD_BANNED_SPECIES,
  MAX_BATTLE_TEAM_SIZE,
  MAX_BEST_BUDDY_BOOSTS,
  MIN_BATTLE_TEAM_SIZE,
  PHASE,
  ACTIONS,
} from './handbook';

function normalizeName(name) {
  return (name || '').trim().toLowerCase();
}

export function validateBattleTeam(members) {
  if (!Array.isArray(members)) {
    return { ok: false, error: 'Battle Team is required (§3.1)' };
  }
  if (members.length < MIN_BATTLE_TEAM_SIZE) {
    return { ok: false, error: `Battle Team must have at least ${MIN_BATTLE_TEAM_SIZE} Pokémon (§3.1)` };
  }
  if (members.length > MAX_BATTLE_TEAM_SIZE) {
    return { ok: false, error: `Battle Team may have at most ${MAX_BATTLE_TEAM_SIZE} Pokémon (§3.1)` };
  }

  const seenDex = new Set();
  let bestBuddy = 0;

  for (const m of members) {
    const name = normalizeName(m.name);
    const dex = m.pokeApiId || 0;
    if (HARD_BANNED_DEX_IDS.has(dex) || HARD_BANNED_SPECIES.has(name)) {
      return { ok: false, error: `${m.name || dex} is banned from tournament play (§3.1.1)` };
    }
    if (m.megaOrPrimal) {
      return { ok: false, error: 'Mega Evolved / Primal Pokémon are not allowed (§3.1.2)' };
    }
    if (!m.combatPower || m.combatPower < 1 || m.combatPower > GREAT_LEAGUE_CP_CAP) {
      return {
        ok: false,
        error: `${m.name || 'Pokémon'} CP must be 1–${GREAT_LEAGUE_CP_CAP} for Great League (§3.1)`,
      };
    }
    if (dex) {
      if (seenDex.has(dex)) {
        return {
          ok: false,
          error: 'Team cannot contain two Pokémon with the same National Pokédex number (§3.1.2)',
        };
      }
      seenDex.add(dex);
    }
    if (m.bestBuddy) bestBuddy += 1;
    if (!m.attacks?.length) {
      return { ok: false, error: `${m.name || 'Pokémon'} must list known attacks (§3.2.2)` };
    }
  }

  if (bestBuddy > MAX_BEST_BUDDY_BOOSTS) {
    return {
      ok: false,
      error: `Team may contain at most ${MAX_BEST_BUDDY_BOOSTS} Best Buddy CP boost (§3.1.2)`,
    };
  }
  return { ok: true };
}

export function validateBattleParty(team, cardIds) {
  if (!Array.isArray(cardIds) || cardIds.length !== BATTLE_PARTY_SIZE) {
    return { ok: false, error: `Must select exactly ${BATTLE_PARTY_SIZE} Pokémon for battle (§3.1)` };
  }
  const teamIds = new Set((team || []).map((c) => c.id));
  const seen = new Set();
  for (const id of cardIds) {
    if (seen.has(id)) return { ok: false, error: 'Battle party cannot repeat the same Pokémon' };
    seen.add(id);
    if (!teamIds.has(id)) {
      return { ok: false, error: 'Battle party Pokémon must come from the registered Battle Team (§3.1)' };
    }
  }
  return { ok: true };
}

/**
 * Client-side gate before API calls. Backend still re-validates.
 */
export function validateAction({ gameState, playerId, action, payload = {} }) {
  if (!gameState && action !== ACTIONS.START_GAME) {
    return { ok: false, error: 'No active game' };
  }

  const phase = gameState?.phase;
  const me = gameState?.players?.find((p) => p.id === playerId);
  const isMyTurn = gameState?.currentTurn === playerId;

  switch (action) {
    case ACTIONS.START_GAME:
      return { ok: true };

    case ACTIONS.SELECT_PARTY: {
      if (![PHASE.TEAM_PREVIEW, PHASE.PARTY_SELECT, PHASE.BETWEEN_GAMES].includes(phase)) {
        return { ok: false, error: 'Can only select a battle party during preview / between games (§6.1)' };
      }
      if (!me) return { ok: false, error: 'Player not found' };
      const teamCheck = validateBattleTeam(me.battleTeam || []);
      if (!teamCheck.ok) return teamCheck;
      return validateBattleParty(me.battleTeam, payload.cardIds);
    }

    case ACTIONS.DRAW_CARD: {
      if (phase !== PHASE.IN_BATTLE) return { ok: false, error: 'Draw only during battle' };
      if (!isMyTurn) return { ok: false, error: 'Not your turn' };
      if (!me?.activePokemon) return { ok: false, error: 'No active Pokémon' };
      if (me.hasDrawn) return { ok: false, error: 'Already drawn this turn' };
      return { ok: true };
    }

    case ACTIONS.SELECT_DRAW:
    case ACTIONS.PLAY_BENCH:
      return { ok: false, error: `${action} is not used in Pokémon GO tournament battles (handbook §6)` };

    case ACTIONS.ATTACH_ENERGY: {
      if (phase !== PHASE.IN_BATTLE) return { ok: false, error: 'Charge energy only during battle' };
      if (!isMyTurn) return { ok: false, error: 'Not your turn' };
      if (!me?.activePokemon) return { ok: false, error: 'No active Pokémon to charge' };
      if (me.hasAttached) return { ok: false, error: 'Already charged energy this turn' };
      return { ok: true };
    }

    case ACTIONS.ATTACK: {
      if (phase !== PHASE.IN_BATTLE) return { ok: false, error: 'Attack only during battle (§6)' };
      if (!isMyTurn) return { ok: false, error: 'Not your turn' };
      if (!me?.activePokemon) return { ok: false, error: 'No active Pokémon' };
      const opp = gameState.players?.find((p) => p.id !== playerId);
      if (!opp?.activePokemon) return { ok: false, error: 'Opponent has no active Pokémon' };
      return { ok: true };
    }

    case ACTIONS.END_TURN: {
      if (phase !== PHASE.IN_BATTLE) return { ok: false, error: 'End turn only during battle' };
      if (!isMyTurn) return { ok: false, error: 'Not your turn' };
      return { ok: true };
    }

    case ACTIONS.PROMOTE: {
      if (phase !== PHASE.IN_BATTLE) return { ok: false, error: 'Promote only during battle' };
      if (me?.activePokemon) return { ok: false, error: 'Active Pokémon still in play' };
      if (!me?.benchedPokemon?.length) return { ok: false, error: 'No benched Pokémon to promote' };
      return { ok: true };
    }

    case ACTIONS.SET_ACTIVE:
      if (phase !== PHASE.IN_BATTLE) return { ok: false, error: 'Set active only during battle' };
      if (!isMyTurn) return { ok: false, error: 'Not your turn' };
      return { ok: true };

    default:
      return { ok: false, error: `Unknown action ${action}` };
  }
}

/** Opponent-safe public preview (§6.5.1). */
export function publicTeamPreview(team = []) {
  return team.map((c) => ({
    name: c.name,
    pokeApiId: c.pokeApiId,
    elementType: c.elementType,
    combatPower: c.combatPower,
    bestBuddy: !!c.bestBuddy,
    moves: (c.attacks || []).map((a) => a.name),
  }));
}
