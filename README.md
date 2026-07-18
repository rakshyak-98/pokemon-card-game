# Pokémon Card Game

A turn-based Pokémon GO Championship–style battle you can play in the browser. Rules are adapted from the Play! Pokémon Pokémon GO Tournament Handbook (Great League, Battle Teams, best-of-three).

## Goal

Recreate the core loop of a Pokémon GO tournament match as a local full-stack app:

- **Start a match** — each player gets a legal Great League Battle Team (up to 6 Pokémon, ≤1500 CP)
- **Team preview** — exchange public team lists, then each locks a party of 3
- **Battle** — charge energy, attack, promote from the back line after a KO
- **Win** — take games by KOing the opponent’s last Pokémon; match is best-of-three

The backend seeds Gen 1 Pokémon (ids 1–151) from [PokeAPI](https://pokeapi.co) into a local SQLite database and drives game state, persistence, handbook rule validation, and an action audit log. The React frontend validates actions client-side before API calls and talks through Vite’s proxy.

## Stack

| Layer | Tech |
|-------|------|
| Frontend | React + Vite |
| Backend | Go HTTP API |
| Database | SQLite (`backend/data/pokemon.db`) |

## How to start the game

You need **Go** (1.25+) and **Node.js** (npm) installed.

### 1. Start the backend

```bash
cd backend
go run .
```

The API listens on [http://localhost:8080](http://localhost:8080). On first boot it may take a moment to seed the Pokémon catalog from PokeAPI.

### 2. Start the frontend

In a second terminal:

```bash
cd frontend
npm install
npm run dev
```

Open the URL Vite prints (usually [http://localhost:5173](http://localhost:5173)). API calls to `/api` are proxied to the backend on port 8080.

### 3. Play

1. Use the UI to **start a match**
2. Review the opponent **team preview**, then **lock a party of 3**
3. **Charge energy**, **attack**, or **end turn**; promote after a KO
4. Win games to take the **best-of-three** match

Handbook PDF: `play-pokemon-pokemon-go-tournament-handbook-en.pdf`

## More detail

- Backend API, patterns, and env options: [backend/README.md](backend/README.md)
