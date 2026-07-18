package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"rakshyak-98/pokemon-backend/models"

	_ "modernc.org/sqlite"
)

// GameStore is the DIP abstraction for persistence (Repository).
type GameStore interface {
	SaveGame(state *models.GameState) error
	LoadGame(id string) (*models.GameState, error)
	LoadLatestGame() (*models.GameState, error)
	AppendAction(log models.ActionLog) (models.ActionLog, error)
	ListActions(gameID string, limit int) ([]models.ActionLog, error)
	CountPokemon() (int, error)
	UpsertPokemon(p models.Pokemon) error
	ListPokemon() ([]models.Pokemon, error)
	GetPokemon(pokeAPIID int) (*models.Pokemon, error)
	CountPowerCards() (int, error)
	UpsertPowerCard(p models.PowerCard) error
	ListPowerCards() ([]models.PowerCard, error)
	GetPowerCard(pokeAPIID int) (*models.PowerCard, error)
	Close() error
}

var (
	once     sync.Once
	instance *SQLiteStore
	initErr  error
)

// GetSQLite is a Singleton access point for the local DB connection.
func GetSQLite(dsn string) (*SQLiteStore, error) {
	once.Do(func() {
		instance, initErr = openSQLite(dsn)
	})
	return instance, initErr
}

// OpenForTest opens a fresh SQLite DB (bypasses Singleton for isolated tests).
func OpenForTest(dsn string) (*SQLiteStore, error) {
	return openSQLite(dsn)
}

type SQLiteStore struct {
	db *sql.DB
	mu sync.Mutex
}

func openSQLite(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *SQLiteStore) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS games (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    current_turn TEXT,
    winner TEXT,
    turn_number INTEGER NOT NULL DEFAULT 1,
    state_json TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS action_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    action_type TEXT NOT NULL,
    payload_json TEXT,
    success INTEGER NOT NULL,
    error_message TEXT,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games(id)
);

CREATE INDEX IF NOT EXISTS idx_action_logs_game ON action_logs(game_id, created_at);
`
	if _, err := s.db.Exec(schema); err != nil {
		return err
	}
	if err := s.migratePokemon(); err != nil {
		return err
	}
	return s.migratePowerCards()
}

func (s *SQLiteStore) SaveGame(state *models.GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	state.UpdatedAt = now
	raw, err := json.Marshal(state)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
INSERT INTO games (id, status, current_turn, winner, turn_number, state_json, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    status = excluded.status,
    current_turn = excluded.current_turn,
    winner = excluded.winner,
    turn_number = excluded.turn_number,
    state_json = excluded.state_json,
    updated_at = excluded.updated_at
`, state.ID, string(state.Status), state.CurrentTurn, state.Winner, state.TurnNumber, string(raw), now, now)
	return err
}

func (s *SQLiteStore) LoadGame(id string) (*models.GameState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var raw string
	err := s.db.QueryRow(`SELECT state_json FROM games WHERE id = ?`, id).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var state models.GameState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *SQLiteStore) LoadLatestGame() (*models.GameState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var raw string
	err := s.db.QueryRow(`SELECT state_json FROM games ORDER BY updated_at DESC LIMIT 1`).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var state models.GameState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *SQLiteStore) AppendAction(log models.ActionLog) (models.ActionLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	success := 0
	if log.Success {
		success = 1
	}
	res, err := s.db.Exec(`
INSERT INTO action_logs (game_id, player_id, action_type, payload_json, success, error_message, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
`, log.GameID, log.PlayerID, log.ActionType, log.PayloadJSON, success, log.ErrorMessage, log.CreatedAt)
	if err != nil {
		return log, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return log, err
	}
	log.ID = id
	return log, nil
}

func (s *SQLiteStore) ListActions(gameID string, limit int) ([]models.ActionLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(`
SELECT id, game_id, player_id, action_type, COALESCE(payload_json, ''), success, COALESCE(error_message, ''), created_at
FROM action_logs
WHERE game_id = ?
ORDER BY id DESC
LIMIT ?
`, gameID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.ActionLog
	for rows.Next() {
		var a models.ActionLog
		var success int
		if err := rows.Scan(&a.ID, &a.GameID, &a.PlayerID, &a.ActionType, &a.PayloadJSON, &success, &a.ErrorMessage, &a.CreatedAt); err != nil {
			return nil, err
		}
		a.Success = success == 1
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

var ErrNotFound = errors.New("not found")

// MemoryStateStore holds the live game in process memory (local state storage).
type MemoryStateStore struct {
	mu    sync.RWMutex
	state *models.GameState
}

func NewMemoryStateStore() *MemoryStateStore {
	return &MemoryStateStore{}
}

func (m *MemoryStateStore) Get() *models.GameState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *MemoryStateStore) Set(state *models.GameState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = state
}
