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
- **power_cards** — special power cards seeded from PokeAPI items (X Attack, potions, pinch berries, etc.)

On first boot the server fetches Gen 1 (ids 1–151) into `pokemons` and builds decks from that catalog. It also seeds battle items into `power_cards` for the special-power deck. Override with:

| Env | Default | Meaning |
|-----|---------|---------|
| `POKEAPI_SEED_FROM` | `1` | first national dex id |
| `POKEAPI_SEED_TO` | `151` | last national dex id |
| `POKEAPI_SEED_WORKERS` | `6` | concurrent fetch workers |
| `POKEAPI_SEED_FORCE` | — | set `1` to re-fetch Pokémon **and** power cards |
| `POWER_SEED_FORCE` | — | set `1` to re-fetch only power cards |

## API

| Method | Path | Body |
|--------|------|------|
| GET | `/api/pokemon` | — (full catalog) |
| GET | `/api/pokemon/{id}` | — (one entry by PokeAPI id) |
| GET | `/api/power-cards` | — (special power-card catalog) |
| GET | `/api/power-cards/{id}` | — (one power card by PokeAPI item id) |
| GET | `/api/game` | — |
| POST | `/api/game/start` | `{ "vsCPU": true }` optional — practice match vs AI |
| POST | `/api/game/select-party` | `{ "playerId", "cardIds": [3 ids] }` |
| POST | `/api/game/draw` | disabled under GO handbook rules |
| POST | `/api/game/draw/select` | disabled under GO handbook rules |
| POST | `/api/game/play-bench` | disabled — party select handles lineup |
| POST | `/api/game/set-active` | `{ "playerId", "cardId" }` |
| POST | `/api/game/attach-energy` | `{ "playerId" }` — charge energy on Active (once/turn) |
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

## Play loop (Pokémon GO handbook)

1. Start match → each player gets a legal Great League Battle Team (up to 6)  
2. Team preview → each selects a party of 3  
3. Charge energy / attack / promote after KO  
4. Win the game by KOing the last opposing Pokémon; match is best-of-three (§6.4)  

