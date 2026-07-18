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

	client := pokeapi.NewClient()
	seedForce := os.Getenv("POKEAPI_SEED_FORCE") == "1"
	powerForce := seedForce || os.Getenv("POWER_SEED_FORCE") == "1"

	seedFrom := envInt("POKEAPI_SEED_FROM", 1)
	seedTo := envInt("POKEAPI_SEED_TO", 151)
	if err := pokeapi.SeedIfEmpty(db, client, pokeapi.SeedOptions{
		FromID:  seedFrom,
		ToID:    seedTo,
		Workers: envInt("POKEAPI_SEED_WORKERS", 6),
		Force:   seedForce,
	}); err != nil {
		log.Printf("warning: pokeapi seed failed (using fallback catalog): %v", err)
	}

	if err := pokeapi.SeedPowerIfEmpty(db, client, pokeapi.PowerSeedOptions{
		Workers: envInt("POKEAPI_SEED_WORKERS", 6),
		Force:   powerForce,
	}); err != nil {
		log.Printf("warning: power card seed failed (using fallback templates): %v", err)
	}

	catalog, err := db.ListPokemon()
	if err != nil {
		log.Fatalf("load pokemon catalog: %v", err)
	}
	log.Printf("loaded %d pokemon into deck catalog\n", len(catalog))

	powers, err := db.ListPowerCards()
	if err != nil {
		log.Fatalf("load power card catalog: %v", err)
	}
	log.Printf("loaded %d power cards into special-card catalog\n", len(powers))

	memory := store.NewMemoryStateStore()
	facade := service.NewGameFacade(db, memory, catalog, powers)
	gameHandler := handlers.NewGameHandler(facade)
	pokemonHandler := handlers.NewPokemonHandler(db)
	powerHandler := handlers.NewPowerHandler(db)

	mux := http.NewServeMux()
	gameHandler.Register(mux)
	pokemonHandler.Register(mux)
	powerHandler.Register(mux)
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
