# Pokémon Card Game — Backend

Local Go server with SQLite persistence, in-memory live state, and a full action audit log.

## Patterns used

| Pattern | Where |
|---------|--------|
| **Facade** | `service.GameFacade` — single entry for play + persist + audit |
| **Command** | `command/*` — each player action is an object (`Execute`) |
| **Repository / DIP** | `store.GameStore` interface; SQLite implements it |
| **Singleton** | `store.GetSQLite` — one DB connection |
| **Memento** | full `GameState` JSON snapshot in `games.state_json` |
| **State** | turn/status rules in `game.Engine` (`Waiting` → `InProgress` → `GameOver`) |

## Database (local SQLite)

Default path: `backend/data/pokemon.db` (override with `DATABASE_PATH`).

- **games** — live game row + JSON state snapshot
- **action_logs** — every user action (success and failure)
- **pokemons** — catalog seeded from [PokeAPI](https://pokeapi.co) (official artwork, types, base stats, card HP/attacks)

On first boot the server fetches Gen 1 (ids 1–151) into `pokemons` and builds decks from that catalog. Override with:

| Env | Default | Meaning |
|-----|---------|---------|
| `POKEAPI_SEED_FROM` | `1` | first national dex id |
| `POKEAPI_SEED_TO` | `151` | last national dex id |
| `POKEAPI_SEED_WORKERS` | `6` | concurrent fetch workers |
| `POKEAPI_SEED_FORCE` | — | set `1` to re-fetch and overwrite |

## API

| Method | Path | Body |
|--------|------|------|
| GET | `/api/pokemon` | — (full catalog) |
| GET | `/api/pokemon/{id}` | — (one entry by PokeAPI id) |
| GET | `/api/game` | — |
| POST | `/api/game/start` | — |
| POST | `/api/game/draw` | `{ "playerId" }` |
| POST | `/api/game/play-bench` | `{ "playerId", "cardId" }` |
| POST | `/api/game/set-active` | `{ "playerId", "cardId" }` |
| POST | `/api/game/attach-energy` | `{ "playerId", "cardId" }` |
| POST | `/api/game/attack` | `{ "playerId", "attackIndex" }` |
| POST | `/api/game/end-turn` | `{ "playerId" }` |
| POST | `/api/game/promote` | `{ "playerId", "cardId" }` |
| GET | `/api/game/actions?limit=50` | — |

## Run

```bash
cd backend
go run .
# or: ./pokemon-backend
```

Frontend (Vite proxies `/api` → `:8080`):

```bash
cd frontend && npm run dev
```

## Play loop

1. Start game → each player gets 7-card hand + 6 prizes  
2. Set Active Pokémon (from hand)  
3. Optionally Bench Pokémon / Attach Energy (1 per turn)  
4. Attack (needs enough energy) or End Turn  
5. On KO → take a prize; opponent Promotes from bench  
6. Win by 6 prizes or when opponent has no Pokémon left  
