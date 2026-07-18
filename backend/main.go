package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"

	"rakshyak-98/pokemon-backend/handlers"
	"rakshyak-98/pokemon-backend/pokeapi"
	"rakshyak-98/pokemon-backend/service"
	"rakshyak-98/pokemon-backend/store"
)

func main() {
	_ = godotenv.Load()

	const addr = ":8080"
	dsn := os.Getenv("DATABASE_PATH")
	if dsn == "" {
		dsn = filepath.Join("data", "pokemon.db")
	}
	if err := os.MkdirAll(filepath.Dir(dsn), 0o755); err != nil {
		log.Fatalf("create data dir: %v", err)
	}

	db, err := store.GetSQLite(dsn)
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	seedFrom := envInt("POKEAPI_SEED_FROM", 1)
	seedTo := envInt("POKEAPI_SEED_TO", 151)
	if err := pokeapi.SeedIfEmpty(db, pokeapi.NewClient(), pokeapi.SeedOptions{
		FromID:  seedFrom,
		ToID:    seedTo,
		Workers: envInt("POKEAPI_SEED_WORKERS", 6),
		Force:   os.Getenv("POKEAPI_SEED_FORCE") == "1",
	}); err != nil {
		log.Printf("warning: pokeapi seed failed (using fallback catalog): %v", err)
	}

	catalog, err := db.ListPokemon()
	if err != nil {
		log.Fatalf("load pokemon catalog: %v", err)
	}
	log.Printf("loaded %d pokemon into deck catalog\n", len(catalog))

	memory := store.NewMemoryStateStore()
	facade := service.NewGameFacade(db, memory, catalog)
	gameHandler := handlers.NewGameHandler(facade)
	pokemonHandler := handlers.NewPokemonHandler(db)

	mux := http.NewServeMux()
	gameHandler.Register(mux)
	pokemonHandler.Register(mux)
	mux.Handle("/", http.FileServer(http.Dir("public")))

	log.Printf("Server starting on %s (sqlite=%s)\n", addr, dsn)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
