package store

import (
	"database/sql"
	"errors"

	"rakshyak-98/pokemon-backend/models"
)

func (s *SQLiteStore) migratePowerCards() error {
	schema := `
CREATE TABLE IF NOT EXISTS power_cards (
    poke_api_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    effect TEXT NOT NULL,
    effect_value INTEGER NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_power_cards_effect ON power_cards(effect);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *SQLiteStore) CountPowerCards() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM power_cards`).Scan(&n)
	return n, err
}

func (s *SQLiteStore) UpsertPowerCard(p models.PowerCard) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
INSERT INTO power_cards (
    poke_api_id, name, image_url, effect, effect_value, category, description
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(poke_api_id) DO UPDATE SET
    name = excluded.name,
    image_url = excluded.image_url,
    effect = excluded.effect,
    effect_value = excluded.effect_value,
    category = excluded.category,
    description = excluded.description
`, p.PokeAPIID, p.Name, p.ImageURL, p.Effect, p.EffectValue, p.Category, p.Description)
	return err
}

func (s *SQLiteStore) ListPowerCards() ([]models.PowerCard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`
SELECT poke_api_id, name, image_url, effect, effect_value, category, description
FROM power_cards
ORDER BY effect, effect_value, poke_api_id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PowerCard
	for rows.Next() {
		p, err := scanPowerCard(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetPowerCard(pokeAPIID int) (*models.PowerCard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow(`
SELECT poke_api_id, name, image_url, effect, effect_value, category, description
FROM power_cards
WHERE poke_api_id = ?
`, pokeAPIID)
	p, err := scanPowerCard(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func scanPowerCard(row rowScanner) (models.PowerCard, error) {
	var p models.PowerCard
	err := row.Scan(
		&p.PokeAPIID, &p.Name, &p.ImageURL, &p.Effect, &p.EffectValue, &p.Category, &p.Description,
	)
	if err != nil {
		return p, err
	}
	if p.Name != "" {
		p.Name = titleCaseName(p.Name)
	}
	return p, nil
}
