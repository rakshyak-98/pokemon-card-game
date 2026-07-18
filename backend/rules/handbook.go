// Package rules encodes Play! Pokémon Pokémon GO Tournament Handbook (rev. May 21, 2026)
// constraints for team construction, match structure, and in-game information classes.
package rules

// Handbook revision tracked for UI / docs.
const HandbookRevision = "May 21, 2026"

// --- §3.1 Battle Team Setup ---

const (
	// MinBattleTeamSize / MaxBattleTeamSize — team may consist of up to six Pokémon, minimum three (§3.1).
	MinBattleTeamSize = 3
	MaxBattleTeamSize = 6

	// BattlePartySize — competitors choose any three of their total Pokémon to bring to each battle (§3.1).
	BattlePartySize = 3

	// GreatLeagueCPCap — tournaments use Great League; Pokémon must be 1,500 CP or less (§3.1).
	GreatLeagueCPCap = 1500

	// MaxBestBuddyBoosts — no more than one Pokémon benefitting from a Best Buddy CP boost (§3.1.2).
	MaxBestBuddyBoosts = 1
)

// HardBannedSpecies cannot be included due to in-game restrictions (§3.1.1).
// Names are lowercase for case-insensitive matching.
var HardBannedSpecies = map[string]struct{}{
	"ditto":    {},
	"shedinja": {},
	"xerneas":  {},
	"yveltal":  {},
}

// HardBannedDexIDs mirrors HardBannedSpecies by National Pokédex number.
var HardBannedDexIDs = map[int]struct{}{
	132: {}, // Ditto
	292: {}, // Shedinja
	716: {}, // Xerneas
	717: {}, // Yveltal
}

// --- §6 Match / Game timing ---

const (
	// MaxBattleSeconds — Pokémon GO maximum battle length (§6.1): 270 seconds.
	MaxBattleSeconds = 270

	// TeamPreviewMaxSeconds — maximum time for team preview phase (§6.1): 2 minutes.
	TeamPreviewMaxSeconds = 120

	// BetweenGamesMaxSeconds — time between battles should not exceed 2 minutes (§6.1).
	BetweenGamesMaxSeconds = 120

	// BestOfThreeWins — most matches are best-of-three (§6.4).
	BestOfThreeWins = 2

	// BestOfFiveWins — Winners/Losers/Grand Finals are best-of-five (§6.4).
	BestOfFiveWins = 3

	// DefaultMatchFormatWins required to take the match (best-of-three → 2).
	DefaultMatchFormatWins = BestOfThreeWins
)

// ProtectShieldsPerGame — each competitor starts a game with 2 Protect Shields (public info §6.5.1).
const ProtectShieldsPerGame = 2

// --- Game phases (adapted for this app) ---

type Phase string

const (
	PhaseWaiting      Phase = "Waiting"
	PhaseTeamPreview  Phase = "TeamPreview"
	PhasePartySelect  Phase = "PartySelect"
	PhaseInBattle     Phase = "InBattle"
	PhaseBetweenGames Phase = "BetweenGames"
	PhaseMatchOver    Phase = "MatchOver"
)

// --- §6.5 Information classes ---

type InfoClass string

const (
	InfoPublic   InfoClass = "public"
	InfoInferred InfoClass = "inferred"
	InfoPrivate  InfoClass = "private"
)

// PublicFields are always shareable with the opponent (§6.5.1).
var PublicFields = []string{
	"species", "form", "cp", "bestBuddy", "moves", "shadowOrPurified",
	"protectShields", "pokemonRemaining", "types",
}

// PrivateFields must not be exposed to the opponent (§6.5.3).
var PrivateFields = []string{
	"energy", "hp", "switchTimer", "unrevealedBackLine",
}
