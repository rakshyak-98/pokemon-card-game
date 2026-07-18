package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rakshyak-98/pokemon-backend/command"
	"rakshyak-98/pokemon-backend/service"
)

type GameHandler struct {
	Facade *service.GameFacade
}

func NewGameHandler(facade *service.GameFacade) *GameHandler {
	return &GameHandler{Facade: facade}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func (h *GameHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/game", withCORS(h.GetGameState))
	mux.HandleFunc("/api/game/start", withCORS(h.StartGame))
	mux.HandleFunc("/api/game/draw", withCORS(h.DrawCard))
	mux.HandleFunc("/api/game/draw/select", withCORS(h.SelectDraw))
	mux.HandleFunc("/api/game/select-party", withCORS(h.SelectParty))
	mux.HandleFunc("/api/game/play-bench", withCORS(h.PlayBench))
	mux.HandleFunc("/api/game/set-active", withCORS(h.SetActive))
	mux.HandleFunc("/api/game/attach-energy", withCORS(h.AttachEnergy))
	mux.HandleFunc("/api/game/attack", withCORS(h.Attack))
	mux.HandleFunc("/api/game/end-turn", withCORS(h.EndTurn))
	mux.HandleFunc("/api/game/promote", withCORS(h.Promote))
	mux.HandleFunc("/api/game/actions", withCORS(h.ListActions))
}

func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, h.Facade.State())
}

type ActionRequest struct {
	PlayerID    string   `json:"playerId"`
	CardID      string   `json:"cardId,omitempty"`
	CardIDs     []string `json:"cardIds,omitempty"`
	AttackIndex int      `json:"attackIndex"`
}

func (h *GameHandler) run(w http.ResponseWriter, cmd command.Command) {
	state, err := h.Facade.Execute(cmd)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		VsCPU bool `json:"vsCPU"`
	}
	if r.Body != nil && r.ContentLength != 0 {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	h.run(w, &command.StartGameCommand{
		Receiver:  h.Facade,
		Player1ID: "player1",
		Player2ID: "player2",
		VsCPU:     req.VsCPU,
	})
}

func (h *GameHandler) DrawCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.DrawCardCommand{Receiver: h.Facade, PID: req.PlayerID})
}

func (h *GameHandler) SelectDraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.SelectDrawCommand{Receiver: h.Facade, PID: req.PlayerID, CardID: req.CardID})
}

func (h *GameHandler) SelectParty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.SelectPartyCommand{Receiver: h.Facade, PID: req.PlayerID, CardIDs: req.CardIDs})
}

func (h *GameHandler) PlayBench(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.PlayBenchCommand{Receiver: h.Facade, PID: req.PlayerID, CardID: req.CardID})
}

func (h *GameHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.SetActiveCommand{Receiver: h.Facade, PID: req.PlayerID, CardID: req.CardID})
}

func (h *GameHandler) AttachEnergy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.AttachEnergyCommand{Receiver: h.Facade, PID: req.PlayerID, EnergyCardID: req.CardID})
}

func (h *GameHandler) Attack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.AttackCommand{Receiver: h.Facade, PID: req.PlayerID, AttackIndex: req.AttackIndex})
}

func (h *GameHandler) EndTurn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.EndTurnCommand{Receiver: h.Facade, PID: req.PlayerID})
}

func (h *GameHandler) Promote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.run(w, &command.PromoteCommand{Receiver: h.Facade, PID: req.PlayerID, CardID: req.CardID})
}

func (h *GameHandler) ListActions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	limit := 50
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			limit = n
		}
	}
	actions, err := h.Facade.ListActions(limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, actions)
}
