package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"rakshyak-98/pokemon-backend/service"
	"rakshyak-98/pokemon-backend/store"
)

func setupHandler(t *testing.T) *GameHandler {
	t.Helper()
	dir := t.TempDir()
	dsn := filepath.Join(dir, "test.db")
	db, err := store.OpenForTest(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	facade := service.NewGameFacade(db, store.NewMemoryStateStore())
	return NewGameHandler(facade)
}

func TestGameFlow(t *testing.T) {
	h := setupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/game/start", nil)
	rr := httptest.NewRecorder()
	h.StartGame(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("start: status %d body %s", rr.Code, rr.Body.String())
	}

	var state map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	if state["status"] != "InProgress" {
		t.Fatalf("expected InProgress, got %v", state["status"])
	}

	body, _ := json.Marshal(ActionRequest{PlayerID: "player1"})
	req = httptest.NewRequest(http.MethodPost, "/api/game/draw", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	h.DrawCard(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("draw: status %d body %s", rr.Code, rr.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/game/actions", nil)
	rr = httptest.NewRecorder()
	h.ListActions(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("actions: status %d", rr.Code)
	}
	var actions []map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&actions); err != nil {
		t.Fatal(err)
	}
	if len(actions) < 2 {
		t.Fatalf("expected at least 2 action logs, got %d", len(actions))
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
