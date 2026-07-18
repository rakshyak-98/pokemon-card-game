package models

import "time"

type CardType string

const (
	TypePokemon CardType = "Pokemon"
	TypeEnergy  CardType = "Energy"
	TypeTrainer CardType = "Trainer"
)

type Attack struct {
	Name   string `json:"name"`
	Damage int    `json:"damage"`
	Cost   int    `json:"cost"` // Energy cost
}

// PokemonStats are base stats from PokeAPI (stored on cards for UI).
type PokemonStats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"spAttack"`
	SpDefense int `json:"spDefense"`
	Speed     int `json:"speed"`
}

// Pokemon is a catalog entry seeded from PokeAPI.
type Pokemon struct {
	PokeAPIID   int          `json:"pokeApiId"`
	Name        string       `json:"name"`
	ImageURL    string       `json:"imageUrl"`
	PrimaryType string       `json:"primaryType"`
	Types       []string     `json:"types"`
	Stats       PokemonStats `json:"stats"`
	CardHP      int          `json:"cardHp"`
	Attacks     []Attack     `json:"attacks"`
}

type Card struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Type           CardType      `json:"type"`
	HP             int           `json:"hp,omitempty"`
	MaxHP          int           `json:"maxHp,omitempty"`
	Attacks        []Attack      `json:"attacks,omitempty"`
	EnergyAttached int           `json:"energyAttached,omitempty"`
	EnergyType     string        `json:"energyType,omitempty"`
	ImageURL       string        `json:"imageUrl,omitempty"`
	ElementType    string        `json:"elementType,omitempty"`
	PokeAPIID      int           `json:"pokeApiId,omitempty"`
	Stats          *PokemonStats `json:"stats,omitempty"`
}

type PlayerState struct {
	ID             string `json:"id"`
	Deck           []Card `json:"deck"`
	Hand           []Card `json:"hand"`
	ActivePokemon  *Card  `json:"activePokemon"`
	BenchedPokemon []Card `json:"benchedPokemon"`
	PrizeCards     []Card `json:"prizeCards"`
	DiscardPile    []Card `json:"discardPile"`
	HasDrawn       bool   `json:"hasDrawn"`
	HasAttached    bool   `json:"hasAttached"`
	PrizesTaken    int    `json:"prizesTaken"`
}

type GameStatus string

const (
	StatusWaiting    GameStatus = "Waiting"
	StatusSetup      GameStatus = "Setup"
	StatusInProgress GameStatus = "InProgress"
	StatusGameOver   GameStatus = "GameOver"
)

type GameState struct {
	ID          string        `json:"id"`
	Status      GameStatus    `json:"status"`
	Players     []PlayerState `json:"players"`
	CurrentTurn string        `json:"currentTurn"`
	Winner      string        `json:"winner,omitempty"`
	TurnNumber  int           `json:"turnNumber"`
	LastAction  string        `json:"lastAction,omitempty"`
	UpdatedAt   time.Time     `json:"updatedAt,omitempty"`
}

// ActionLog is an audit record of a player action (Command + history).
type ActionLog struct {
	ID           int64     `json:"id"`
	GameID       string    `json:"gameId"`
	PlayerID     string    `json:"playerId"`
	ActionType   string    `json:"actionType"`
	PayloadJSON  string    `json:"payloadJson,omitempty"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}
