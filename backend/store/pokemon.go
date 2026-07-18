package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"rakshyak-98/pokemon-backend/models"
)

func (s *SQLiteStore) migratePokemon() error {
	schema := `
CREATE TABLE IF NOT EXISTS pokemons (
    poke_api_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    primary_type TEXT NOT NULL,
    types_json TEXT NOT NULL,
    hp INTEGER NOT NULL,
    attack INTEGER NOT NULL,
    defense INTEGER NOT NULL,
    sp_attack INTEGER NOT NULL,
    sp_defense INTEGER NOT NULL,
    speed INTEGER NOT NULL,
    card_hp INTEGER NOT NULL,
    attacks_json TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pokemons_name ON pokemons(name);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *SQLiteStore) CountPokemon() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM pokemons`).Scan(&n)
	return n, err
}

func (s *SQLiteStore) UpsertPokemon(p models.Pokemon) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	typesJSON, err := json.Marshal(p.Types)
	if err != nil {
		return err
	}
	attacksJSON, err := json.Marshal(p.Attacks)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
INSERT INTO pokemons (
    poke_api_id, name, image_url, primary_type, types_json,
    hp, attack, defense, sp_attack, sp_defense, speed,
    card_hp, attacks_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(poke_api_id) DO UPDATE SET
    name = excluded.name,
    image_url = excluded.image_url,
    primary_type = excluded.primary_type,
    types_json = excluded.types_json,
    hp = excluded.hp,
    attack = excluded.attack,
    defense = excluded.defense,
    sp_attack = excluded.sp_attack,
    sp_defense = excluded.sp_defense,
    speed = excluded.speed,
    card_hp = excluded.card_hp,
    attacks_json = excluded.attacks_json
`, p.PokeAPIID, p.Name, p.ImageURL, p.PrimaryType, string(typesJSON),
		p.Stats.HP, p.Stats.Attack, p.Stats.Defense, p.Stats.SpAttack, p.Stats.SpDefense, p.Stats.Speed,
		p.CardHP, string(attacksJSON))
	return err
}

func (s *SQLiteStore) ListPokemon() ([]models.Pokemon, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`
SELECT poke_api_id, name, image_url, primary_type, types_json,
       hp, attack, defense, sp_attack, sp_defense, speed,
       card_hp, attacks_json
FROM pokemons
ORDER BY poke_api_id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Pokemon
	for rows.Next() {
		p, err := scanPokemon(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetPokemon(pokeAPIID int) (*models.Pokemon, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow(`
SELECT poke_api_id, name, image_url, primary_type, types_json,
       hp, attack, defense, sp_attack, sp_defense, speed,
       card_hp, attacks_json
FROM pokemons
WHERE poke_api_id = ?
`, pokeAPIID)
	p, err := scanPokemon(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPokemon(row rowScanner) (models.Pokemon, error) {
	var p models.Pokemon
	var typesJSON, attacksJSON string
	err := row.Scan(
		&p.PokeAPIID, &p.Name, &p.ImageURL, &p.PrimaryType, &typesJSON,
		&p.Stats.HP, &p.Stats.Attack, &p.Stats.Defense, &p.Stats.SpAttack, &p.Stats.SpDefense, &p.Stats.Speed,
		&p.CardHP, &attacksJSON,
	)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal([]byte(typesJSON), &p.Types); err != nil {
		return p, err
	}
	if err := json.Unmarshal([]byte(attacksJSON), &p.Attacks); err != nil {
		return p, err
	}
	if p.Name != "" {
		p.Name = titleCaseName(p.Name)
	}
	return p, nil
}

func titleCaseName(name string) string {
	if name == "" {
		return name
	}
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
