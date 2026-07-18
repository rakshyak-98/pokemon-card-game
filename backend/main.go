package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"rakshyak-98/pokemon-backend/handlers"
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

	memory := store.NewMemoryStateStore()
	facade := service.NewGameFacade(db, memory)
	handler := handlers.NewGameHandler(facade)

	mux := http.NewServeMux()
	handler.Register(mux)
	mux.Handle("/", http.FileServer(http.Dir("public")))

	log.Printf("Server starting on %s (sqlite=%s)\n", addr, dsn)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
