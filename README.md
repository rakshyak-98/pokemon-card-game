# Pokémon Card Game

A turn-based Pokémon trading-card style game you can play in the browser. Two players take turns drawing cards, setting Active and Bench Pokémon, attaching energy, and attacking until one side wins by collecting all prize cards or knocking out the opponent’s last Pokémon.

## Goal

Recreate the core loop of a Pokémon TCG match as a local full-stack app:

- **Start a match** — each player gets a 7-card hand and 6 prize cards
- **Set up the field** — choose an Active Pokémon; optionally bench others and attach energy (one per turn)
- **Battle** — attack when you have enough energy, or end your turn
- **Win** — take a prize on each KO; promote from the bench after a KO; win with 6 prizes or when the opponent has no Pokémon left

The backend seeds Gen 1 Pokémon (ids 1–151) from [PokeAPI](https://pokeapi.co) into a local SQLite database and drives game state, persistence, and an action audit log. The React frontend talks to the API through Vite’s proxy.

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

1. Use the UI to **start a game**
2. Set each player’s **Active Pokémon** from their hand
3. Optionally **bench** Pokémon and **attach energy**
4. **Attack** or **end turn**, then continue until someone wins

## More detail

- Backend API, patterns, and env options: [backend/README.md](backend/README.md)
