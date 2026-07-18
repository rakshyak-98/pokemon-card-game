package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"rakshyak-98/pokemon-backend/models"
	"rakshyak-98/pokemon-backend/store"
)

type PokemonHandler struct {
	DB store.GameStore
}

func NewPokemonHandler(db store.GameStore) *PokemonHandler {
	return &PokemonHandler{DB: db}
}

func (h *PokemonHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/pokemon", withCORS(h.ListPokemon))
	mux.HandleFunc("/api/pokemon/", withCORS(h.GetPokemon))
}

func (h *PokemonHandler) ListPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, err := h.DB.ListPokemon()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if list == nil {
		list = []models.Pokemon{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/pokemon/")
	if idStr == "" {
		h.ListPokemon(w, r)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid pokemon id"})
		return
	}
	p, err := h.DB.GetPokemon(id)
	if err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "pokemon not found"})
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}
