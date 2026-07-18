package service

import (
	"errors"
	"sync"
	"time"

	"rakshyak-98/pokemon-backend/command"
	"rakshyak-98/pokemon-backend/game"
	"rakshyak-98/pokemon-backend/models"
	"rakshyak-98/pokemon-backend/store"
)

// GameFacade is a thin orchestrator over engine + memory + DB + audit (Facade pattern).
type GameFacade struct {
	mu     sync.Mutex
	engine *game.Engine
	memory *store.MemoryStateStore
	db     store.GameStore
}

func NewGameFacade(db store.GameStore, memory *store.MemoryStateStore) *GameFacade {
	f := &GameFacade{
		engine: game.NewEngine(),
		memory: memory,
		db:     db,
	}
	f.hydrate()
	return f
}

func (f *GameFacade) hydrate() {
	state, err := f.db.LoadLatestGame()
	if err != nil {
		f.memory.Set(f.engine.State)
		return
	}
	f.engine.Restore(state)
	f.memory.Set(state)
}

func (f *GameFacade) State() *models.GameState {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.engine.State
}

// Execute runs a Command, persists state (Memento), and writes the action audit log.
func (f *GameFacade) Execute(cmd command.Command) (*models.GameState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	err := cmd.Execute()
	gameID := ""
	if f.engine.State != nil {
		gameID = f.engine.State.ID
	}

	logEntry := models.ActionLog{
		GameID:       gameID,
		PlayerID:     cmd.PlayerID(),
		ActionType:   cmd.Name(),
		PayloadJSON:  command.MarshalPayload(cmd),
		Success:      err == nil,
		CreatedAt:    time.Now().UTC(),
	}
	if err != nil {
		logEntry.ErrorMessage = err.Error()
		if gameID != "" {
			_, _ = f.db.AppendAction(logEntry)
		}
		return f.engine.State, err
	}

	f.engine.State.UpdatedAt = time.Now().UTC()
	f.memory.Set(f.engine.State)
	if saveErr := f.db.SaveGame(f.engine.State); saveErr != nil {
		return f.engine.State, saveErr
	}
	if _, logErr := f.db.AppendAction(logEntry); logErr != nil {
		return f.engine.State, logErr
	}
	return f.engine.State, nil
}

func (f *GameFacade) ListActions(limit int) ([]models.ActionLog, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.engine.State == nil || f.engine.State.ID == "" {
		return []models.ActionLog{}, nil
	}
	return f.db.ListActions(f.engine.State.ID, limit)
}

// --- command.GameActions receiver ---

func (f *GameFacade) StartGame(player1ID, player2ID string) error {
	return f.engine.StartGame(player1ID, player2ID)
}

func (f *GameFacade) DrawCard(playerID string) error {
	return f.engine.DrawCard(playerID)
}

func (f *GameFacade) PlayBench(playerID, cardID string) error {
	return f.engine.PlayBench(playerID, cardID)
}

func (f *GameFacade) SetActive(playerID, cardID string) error {
	return f.engine.SetActive(playerID, cardID)
}

func (f *GameFacade) AttachEnergy(playerID, energyCardID string) error {
	return f.engine.AttachEnergy(playerID, energyCardID)
}

func (f *GameFacade) Attack(playerID string, attackIndex int) error {
	return f.engine.Attack(playerID, attackIndex)
}

func (f *GameFacade) EndTurn(playerID string) error {
	return f.engine.EndTurn(playerID)
}

func (f *GameFacade) Promote(playerID, cardID string) error {
	return f.engine.Promote(playerID, cardID)
}

// Ensure GameFacade satisfies command.GameActions at compile time.
var _ command.GameActions = (*GameFacade)(nil)

var ErrNoGame = errors.New("no active game")
