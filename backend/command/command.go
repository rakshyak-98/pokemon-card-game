package command

import (
	"encoding/json"

	"rakshyak-98/pokemon-backend/models"
)

// Action names for the audit log.
const (
	ActionStartGame   = "start_game"
	ActionDrawCard    = "draw_card"
	ActionPlayBench   = "play_bench"
	ActionSetActive   = "set_active"
	ActionAttachEnergy = "attach_energy"
	ActionAttack      = "attack"
	ActionEndTurn     = "end_turn"
	ActionPromote     = "promote"
)

// Command encapsulates a request as an object (Command pattern).
type Command interface {
	Name() string
	PlayerID() string
	Payload() any
	Execute() error
}

// GameActions is the receiver interface Command objects call into.
type GameActions interface {
	StartGame(player1ID, player2ID string) error
	DrawCard(playerID string) error
	PlayBench(playerID, cardID string) error
	SetActive(playerID, cardID string) error
	AttachEnergy(playerID, energyCardID string) error
	Attack(playerID string, attackIndex int) error
	EndTurn(playerID string) error
	Promote(playerID, cardID string) error
	State() *models.GameState
}

type StartGameCommand struct {
	Receiver             GameActions
	Player1ID, Player2ID string
}

func (c *StartGameCommand) Name() string     { return ActionStartGame }
func (c *StartGameCommand) PlayerID() string { return c.Player1ID }
func (c *StartGameCommand) Payload() any {
	return map[string]string{"player1Id": c.Player1ID, "player2Id": c.Player2ID}
}
func (c *StartGameCommand) Execute() error {
	return c.Receiver.StartGame(c.Player1ID, c.Player2ID)
}

type DrawCardCommand struct {
	Receiver GameActions
	PID      string
}

func (c *DrawCardCommand) Name() string     { return ActionDrawCard }
func (c *DrawCardCommand) PlayerID() string { return c.PID }
func (c *DrawCardCommand) Payload() any     { return map[string]string{"playerId": c.PID} }
func (c *DrawCardCommand) Execute() error   { return c.Receiver.DrawCard(c.PID) }

type PlayBenchCommand struct {
	Receiver GameActions
	PID      string
	CardID   string
}

func (c *PlayBenchCommand) Name() string     { return ActionPlayBench }
func (c *PlayBenchCommand) PlayerID() string { return c.PID }
func (c *PlayBenchCommand) Payload() any {
	return map[string]string{"playerId": c.PID, "cardId": c.CardID}
}
func (c *PlayBenchCommand) Execute() error { return c.Receiver.PlayBench(c.PID, c.CardID) }

type SetActiveCommand struct {
	Receiver GameActions
	PID      string
	CardID   string
}

func (c *SetActiveCommand) Name() string     { return ActionSetActive }
func (c *SetActiveCommand) PlayerID() string { return c.PID }
func (c *SetActiveCommand) Payload() any {
	return map[string]string{"playerId": c.PID, "cardId": c.CardID}
}
func (c *SetActiveCommand) Execute() error { return c.Receiver.SetActive(c.PID, c.CardID) }

type AttachEnergyCommand struct {
	Receiver     GameActions
	PID          string
	EnergyCardID string
}

func (c *AttachEnergyCommand) Name() string     { return ActionAttachEnergy }
func (c *AttachEnergyCommand) PlayerID() string { return c.PID }
func (c *AttachEnergyCommand) Payload() any {
	return map[string]string{"playerId": c.PID, "cardId": c.EnergyCardID}
}
func (c *AttachEnergyCommand) Execute() error {
	return c.Receiver.AttachEnergy(c.PID, c.EnergyCardID)
}

type AttackCommand struct {
	Receiver    GameActions
	PID         string
	AttackIndex int
}

func (c *AttackCommand) Name() string     { return ActionAttack }
func (c *AttackCommand) PlayerID() string { return c.PID }
func (c *AttackCommand) Payload() any {
	return map[string]any{"playerId": c.PID, "attackIndex": c.AttackIndex}
}
func (c *AttackCommand) Execute() error { return c.Receiver.Attack(c.PID, c.AttackIndex) }

type EndTurnCommand struct {
	Receiver GameActions
	PID      string
}

func (c *EndTurnCommand) Name() string     { return ActionEndTurn }
func (c *EndTurnCommand) PlayerID() string { return c.PID }
func (c *EndTurnCommand) Payload() any     { return map[string]string{"playerId": c.PID} }
func (c *EndTurnCommand) Execute() error   { return c.Receiver.EndTurn(c.PID) }

type PromoteCommand struct {
	Receiver GameActions
	PID      string
	CardID   string
}

func (c *PromoteCommand) Name() string     { return ActionPromote }
func (c *PromoteCommand) PlayerID() string { return c.PID }
func (c *PromoteCommand) Payload() any {
	return map[string]string{"playerId": c.PID, "cardId": c.CardID}
}
func (c *PromoteCommand) Execute() error { return c.Receiver.Promote(c.PID, c.CardID) }

// MarshalPayload encodes command payload for the audit log.
func MarshalPayload(c Command) string {
	b, err := json.Marshal(c.Payload())
	if err != nil {
		return "{}"
	}
	return string(b)
}
