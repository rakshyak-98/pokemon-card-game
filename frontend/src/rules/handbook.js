/**
 * Play! Pokémon Pokémon GO Tournament Handbook (rev. May 21, 2026)
 * Shared constants for frontend rule validation.
 */
export const HANDBOOK_REVISION = 'May 21, 2026';

export const MIN_BATTLE_TEAM_SIZE = 3;
export const MAX_BATTLE_TEAM_SIZE = 6;
export const BATTLE_PARTY_SIZE = 3;
export const GREAT_LEAGUE_CP_CAP = 1500;
export const MAX_BEST_BUDDY_BOOSTS = 1;
export const PROTECT_SHIELDS_PER_GAME = 2;
export const MAX_POWER_HAND_SLOTS = 4;
export const DEFAULT_MATCH_WINS_NEEDED = 2;
export const MAX_BATTLE_SECONDS = 270;
export const TEAM_PREVIEW_MAX_SECONDS = 120;

export const HARD_BANNED_SPECIES = new Set(['ditto', 'shedinja', 'xerneas', 'yveltal']);
export const HARD_BANNED_DEX_IDS = new Set([132, 292, 716, 717]);

export const PHASE = {
  WAITING: 'Waiting',
  TEAM_PREVIEW: 'TeamPreview',
  PARTY_SELECT: 'PartySelect',
  IN_BATTLE: 'InBattle',
  BETWEEN_GAMES: 'BetweenGames',
  MATCH_OVER: 'MatchOver',
};

export const ACTIONS = {
  START_GAME: 'start_game',
  SELECT_PARTY: 'select_party',
  DRAW_CARD: 'draw_card',
  SELECT_DRAW: 'select_draw',
  PLAY_BENCH: 'play_bench',
  SET_ACTIVE: 'set_active',
  ATTACH_ENERGY: 'attach_energy',
  ATTACK: 'attack',
  END_TURN: 'end_turn',
  PROMOTE: 'promote',
  SWITCH: 'switch',
  PLAY_POWER: 'play_power',
  USE_SHIELD: 'use_shield',
};
