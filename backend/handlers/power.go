package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"rakshyak-98/pokemon-backend/models"
	"rakshyak-98/pokemon-backend/store"
)

type PowerHandler struct {
	DB store.GameStore
}

func NewPowerHandler(db store.GameStore) *PowerHandler {
	return &PowerHandler{DB: db}
}

func (h *PowerHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/power-cards", withCORS(h.ListPowerCards))
	mux.HandleFunc("/api/power-cards/", withCORS(h.GetPowerCard))
}

func (h *PowerHandler) ListPowerCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, err := h.DB.ListPowerCards()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if list == nil {
		list = []models.PowerCard{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *PowerHandler) GetPowerCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/power-cards/")
	if idStr == "" {
		h.ListPowerCards(w, r)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid power card id"})
		return
	}
	p, err := h.DB.GetPowerCard(id)
	if err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "power card not found"})
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}
