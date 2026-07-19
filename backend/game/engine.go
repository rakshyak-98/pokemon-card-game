package game

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	"rakshyak-98/pokemon-backend/models"
	"rakshyak-98/pokemon-backend/rules"
)

// Engine owns mutable game state and turn rules (State-machine style transitions).
type Engine struct {
	State        *models.GameState
	catalog      []models.Pokemon
	powerCatalog []models.PowerCard
	// powerSeq ensures unique power-card IDs across mid-game deck refills.
	powerSeq int
}

func NewEngine(catalog []models.Pokemon) *Engine {
	return &Engine{
		catalog: catalog,
		State: &models.GameState{
			Status:     models.StatusWaiting,
			Players:    []models.PlayerState{},
			TurnNumber: 0,
		},
	}
}

// SetCatalog updates the Pokémon pool used when building decks.
func (e *Engine) SetCatalog(catalog []models.Pokemon) {
	e.catalog = catalog
}

func (e *Engine) StateSnapshot() *models.GameState {
	return e.State
}

// StartGame initializes a Pokémon GO–style match: each player gets a legal Great League
// Battle Team (6 unique Pokémon), then Team Preview / party select begins (§3.1, §6.1).
// If vsCPU is true, player2 is controlled by the practice AI.
func (e *Engine) StartGame(player1ID, player2ID string, vsCPU bool) error {
	if player1ID == "" || player2ID == "" {
		return errors.New("both player ids are required")
	}
	if player1ID == player2ID {
		return errors.New("player ids must be different")
	}

	p1, err := e.createBattlePlayer(player1ID)
	if err != nil {
		return err
	}
	p2, err := e.createBattlePlayer(player2ID)
	if err != nil {
		return err
	}

	last := "match_started — exchange team preview lists (§6.1)"
	if vsCPU {
		last = "practice match vs CPU — pick your battle party to learn the rules"
	}

	e.State = &models.GameState{
		ID:          uuid.NewString(),
		Status:      models.StatusSetup,
		Phase:       string(rules.PhaseTeamPreview),
		Players:     []models.PlayerState{p1, p2},
		CurrentTurn: "",
		TurnNumber:  0,
		GameNumber:  1,
		WinsNeeded:  rules.DefaultMatchFormatWins,
		VsCPU:       vsCPU,
		CPUPlayerID: "",
		LastAction:  last,
		UpdatedAt:   time.Now().UTC(),
	}
	if vsCPU {
		e.State.CPUPlayerID = player2ID
	}
	return nil
}

func (e *Engine) createBattlePlayer(id string) (models.PlayerState, error) {
	team, err := e.buildBattleTeam(id)
	if err != nil {
		return models.PlayerState{}, err
	}
	if err := rules.ValidateBattleTeam(rules.CardsToTeamMembers(team)); err != nil {
		return models.PlayerState{}, err
	}
	return models.PlayerState{
		ID:             id,
		BattleTeam:     team,
		PowerDeck:      e.buildPowerDeck(id),
		Hand:           []models.Card{},
		BenchedPokemon: []models.Card{},
		DiscardPile:    []models.Card{},
		ProtectShields: rules.ProtectShieldsPerGame,
		GamesWon:       0,
		PartyReady:     false,
	}, nil
}

func (e *Engine) buildBattleTeam(id string) ([]models.Card, error) {
	pool := e.catalogPool()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	order := r.Perm(len(pool))

	team := make([]models.Card, 0, rules.MaxBattleTeamSize)
	seen := map[int]struct{}{}

	for _, oi := range order {
		if len(team) >= rules.MaxBattleTeamSize {
			break
		}
		p := pool[oi]
		if _, banned := rules.HardBannedDexIDs[p.PokeAPIID]; banned {
			continue
		}
		if _, banned := rules.HardBannedSpecies[strings.ToLower(p.Name)]; banned {
			continue
		}
		if _, dup := seen[p.PokeAPIID]; dup {
			continue
		}
		if len(p.Attacks) == 0 {
			continue
		}
		seen[p.PokeAPIID] = struct{}{}
		team = append(team, cardFromPokemon(id, len(team), p))
	}

	if len(team) < rules.MinBattleTeamSize {
		return nil, errors.New("catalog too small to build a legal Battle Team (§3.1)")
	}
	return team, nil
}

func (e *Engine) createPlayer(id string) models.PlayerState {
	// Legacy helper retained for older tests; prefer createBattlePlayer.
	p, err := e.createBattlePlayer(id)
	if err != nil {
		deck := e.buildStarterDeck(id)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})
		return models.PlayerState{ID: id, Deck: deck, Hand: []models.Card{}, BenchedPokemon: []models.Card{}, PrizeCards: []models.Card{}, DiscardPile: []models.Card{}}
	}
	return p
}

func fallbackCatalog() []models.Pokemon {
	return []models.Pokemon{
		{
			PokeAPIID: 25, Name: "Pikachu", PrimaryType: "Electric", Types: []string{"Electric"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/25.png",
			Stats:    models.PokemonStats{HP: 35, Attack: 55, Defense: 40, SpAttack: 50, SpDefense: 50, Speed: 90},
			CardHP:   70, Attacks: []models.Attack{{Name: "Thunder Shock", Damage: 20, Cost: 1}},
		},
		{
			PokeAPIID: 4, Name: "Charmander", PrimaryType: "Fire", Types: []string{"Fire"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/4.png",
			Stats:    models.PokemonStats{HP: 39, Attack: 52, Defense: 43, SpAttack: 60, SpDefense: 50, Speed: 65},
			CardHP:   78, Attacks: []models.Attack{{Name: "Ember", Damage: 30, Cost: 1}},
		},
		{
			PokeAPIID: 7, Name: "Squirtle", PrimaryType: "Water", Types: []string{"Water"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/7.png",
			Stats:    models.PokemonStats{HP: 44, Attack: 48, Defense: 65, SpAttack: 50, SpDefense: 64, Speed: 43},
			CardHP:   88, Attacks: []models.Attack{{Name: "Water Gun", Damage: 20, Cost: 1}},
		},
		{
			PokeAPIID: 1, Name: "Bulbasaur", PrimaryType: "Grass", Types: []string{"Grass", "Poison"},
			ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/1.png",
			Stats:    models.PokemonStats{HP: 45, Attack: 49, Defense: 49, SpAttack: 65, SpDefense: 65, Speed: 45},
			CardHP:   90, Attacks: []models.Attack{{Name: "Vine Whip", Damage: 20, Cost: 1}},
		},
	}
}

func (e *Engine) catalogPool() []models.Pokemon {
	if len(e.catalog) > 0 {
		return e.catalog
	}
	return fallbackCatalog()
}

func cardFromPokemon(idPrefix string, index int, p models.Pokemon) models.Card {
	stats := p.Stats
	attacks := append([]models.Attack(nil), p.Attacks...)
	cp := rules.ClampCP(rules.EstimateCombatPower(p.Stats))
	return models.Card{
		ID:          fmt.Sprintf("%s-poke-%d", idPrefix, index),
		Name:        p.Name,
		Type:        models.TypePokemon,
		HP:          p.CardHP,
		MaxHP:       p.CardHP,
		Attacks:     attacks,
		ImageURL:    p.ImageURL,
		ElementType: p.PrimaryType,
		PokeAPIID:   p.PokeAPIID,
		Stats:       &stats,
		CombatPower: cp,
	}
}

func energyForType(element string) string {
	switch strings.ToLower(element) {
	case "electric":
		return "Electric"
	case "fire":
		return "Fire"
	case "water", "ice":
		return "Water"
	case "grass", "bug", "poison":
		return "Grass"
	case "psychic", "ghost", "fairy", "dark":
		return "Psychic"
	default:
		return "Normal"
	}
}

func (e *Engine) buildStarterDeck(id string) []models.Card {
	pool := e.catalogPool()
	deck := make([]models.Card, 0, 60)

	// 20 Pokémon cards drawn from the seeded catalog.
	for i := 0; i < 20; i++ {
		p := pool[i%len(pool)]
		deck = append(deck, cardFromPokemon(id, i, p))
	}

	energyTypes := make([]string, 0, 4)
	seen := map[string]bool{}
	for _, p := range pool {
		et := energyForType(p.PrimaryType)
		if !seen[et] {
			seen[et] = true
			energyTypes = append(energyTypes, et)
		}
		if len(energyTypes) >= 4 {
			break
		}
	}
	if len(energyTypes) == 0 {
		energyTypes = []string{"Electric", "Fire", "Water", "Grass"}
	}

	for i := 0; i < 40; i++ {
		et := energyTypes[i%len(energyTypes)]
		deck = append(deck, models.Card{
			ID:          fmt.Sprintf("%s-energy-%d", id, i),
			Name:        et + " Energy",
			Type:        models.TypeEnergy,
			EnergyType:  et,
			ElementType: et,
		})
	}
	return deck
}

func (e *Engine) drawCard(player *models.PlayerState) error {
	if len(player.Deck) == 0 {
		return errors.New("deck is empty")
	}
	card := player.Deck[0]
	player.Deck = player.Deck[1:]
	player.Hand = append(player.Hand, card)
	return nil
}

func (e *Engine) getPlayer(playerID string) *models.PlayerState {
	for i := range e.State.Players {
		if e.State.Players[i].ID == playerID {
			return &e.State.Players[i]
		}
	}
	return nil
}

func (e *Engine) getOpponent(playerID string) *models.PlayerState {
	for i := range e.State.Players {
		if e.State.Players[i].ID != playerID {
			return &e.State.Players[i]
		}
	}
	return nil
}

func (e *Engine) requireTurn(playerID string) (*models.PlayerState, error) {
	if e.State.Status == models.StatusGameOver || e.State.Phase == string(rules.PhaseMatchOver) {
		return nil, errors.New("game is over")
	}
	if e.State.Status != models.StatusInProgress && e.State.Phase != string(rules.PhaseInBattle) {
		return nil, errors.New("game not in progress")
	}
	player := e.getPlayer(playerID)
	if player == nil {
		return nil, errors.New("player not found")
	}
	if e.State.CurrentTurn != playerID {
		return nil, errors.New("not your turn")
	}
	return player, nil
}

func (e *Engine) requirePlayable(playerID string) (*models.PlayerState, error) {
	player, err := e.requireTurn(playerID)
	if err != nil {
		return nil, err
	}
	if len(player.PendingDraw) > 0 {
		return nil, errors.New("choose a power card to replace (or keep hand)")
	}
	if player.ActivePokemon == nil && len(player.BenchedPokemon) > 0 {
		return nil, errors.New("must promote a benched pokemon first")
	}
	return player, nil
}

// DrawCard pulls one special power card from the player's power deck into hand (once per turn).
// If all power slots are full, the card goes to PendingDraw for a swap/keep choice.
func (e *Engine) DrawCard(playerID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.HasDrawn {
		return errors.New("already drawn this turn")
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon")
	}
	e.ensurePowerDeck(player)
	name, pending, err := drawPowerCard(player)
	if err != nil {
		return err
	}
	player.HasDrawn = true
	if pending {
		e.State.LastAction = fmt.Sprintf("%s drew %s — power slots full, choose a card to replace", playerID, name)
	} else {
		e.State.LastAction = fmt.Sprintf("%s drew power card %s", playerID, name)
	}
	return nil
}

// tryAutoDrawPower draws the turn's power card automatically when legal.
// Uses direct deck→hand logic (not DrawCard/requirePlayable) so turn-start
// draws cannot fail silently after EndTurn / Attack advance.
func (e *Engine) tryAutoDrawPower(playerID string) {
	if e.State == nil || e.State.Phase != string(rules.PhaseInBattle) {
		return
	}
	if e.State.CurrentTurn != playerID {
		return
	}
	player := e.getPlayer(playerID)
	if player == nil || player.HasDrawn || player.ActivePokemon == nil {
		return
	}
	if len(player.PendingDraw) > 0 {
		return
	}
	e.ensurePowerDeck(player)
	if len(player.PowerDeck) == 0 {
		return
	}
	name, pending, err := drawPowerCard(player)
	if err != nil {
		return
	}
	player.HasDrawn = true
	if pending {
		e.State.LastAction = fmt.Sprintf("%s drew %s — power slots full, choose a card to replace", playerID, name)
	} else {
		e.State.LastAction = fmt.Sprintf("%s drew power card %s", playerID, name)
	}
}

// PlayPower plays a special power card from hand onto the active Pokémon.
// Players may use any number of available power-slot cards in a turn.
func (e *Engine) PlayPower(playerID, cardID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon to apply power to")
	}
	idx := findCard(player.Hand, cardID)
	if idx < 0 {
		return errors.New("power card not found in hand")
	}
	card := player.Hand[idx]
	if card.Type != models.TypePower {
		return errors.New("card is not a special power card")
	}

	summary, err := applyPowerEffect(player, card)
	if err != nil {
		return err
	}

	player.Hand = append(player.Hand[:idx], player.Hand[idx+1:]...)
	player.DiscardPile = append(player.DiscardPile, card)
	player.HasPlayedPower = true
	e.State.LastAction = fmt.Sprintf("%s played %s — %s", playerID, card.Name, summary)
	return nil
}

// SelectDraw resolves a full-hand power draw: cardID is the hand card to discard
// and replace with the pending draw. Empty cardID keeps the hand and discards the draw.
func (e *Engine) SelectDraw(playerID, cardID string) error {
	player, err := e.requireTurn(playerID)
	if err != nil {
		return err
	}
	if len(player.PendingDraw) == 0 {
		return errors.New("no pending power card to place")
	}
	pending := player.PendingDraw[0]

	if cardID == "" || cardID == "_keep" {
		player.DiscardPile = append(player.DiscardPile, pending)
		player.PendingDraw = nil
		e.State.LastAction = fmt.Sprintf("%s kept hand and discarded drawn %s", playerID, pending.Name)
		return nil
	}

	idx := findCard(player.Hand, cardID)
	if idx < 0 {
		return errors.New("hand card not found to replace")
	}
	replaced := player.Hand[idx]
	player.Hand[idx] = pending
	player.DiscardPile = append(player.DiscardPile, replaced)
	player.PendingDraw = nil
	e.State.LastAction = fmt.Sprintf("%s replaced %s with %s", playerID, replaced.Name, pending.Name)
	return nil
}

// SelectParty locks in the three Pokémon brought to the current game (§3.1 / §6.1).
func (e *Engine) SelectParty(playerID string, cardIDs []string) error {
	phase := rules.Phase(e.State.Phase)
	if phase != rules.PhaseTeamPreview && phase != rules.PhasePartySelect && phase != rules.PhaseBetweenGames {
		return errors.New("can only select party during preview or between games")
	}
	player := e.getPlayer(playerID)
	if player == nil {
		return errors.New("player not found")
	}
	if player.PartyReady {
		return errors.New("party already selected for this game")
	}
	if len(cardIDs) != rules.BattlePartySize {
		return fmt.Errorf("must select exactly %d Pokémon (§3.1)", rules.BattlePartySize)
	}

	chosen := make([]models.Card, 0, rules.BattlePartySize)
	seen := map[string]struct{}{}
	for _, cid := range cardIDs {
		if _, dup := seen[cid]; dup {
			return errors.New("duplicate party selection")
		}
		seen[cid] = struct{}{}
		idx := findCard(player.BattleTeam, cid)
		if idx < 0 {
			return errors.New("card not on Battle Team")
		}
		c := player.BattleTeam[idx]
		// Fresh copy for this game (full HP / no energy).
		fresh := c
		fresh.HP = c.MaxHP
		if fresh.HP <= 0 {
			fresh.HP = c.HP
			fresh.MaxHP = c.HP
		}
		fresh.EnergyAttached = 0
		chosen = append(chosen, fresh)
	}

	dexIDs := make([]int, 0, len(chosen))
	for _, c := range chosen {
		dexIDs = append(dexIDs, c.PokeAPIID)
	}
	if err := rules.ValidateBattleParty(rules.CardsToTeamMembers(player.BattleTeam), dexIDs); err != nil {
		return err
	}

	active := chosen[0]
	player.ActivePokemon = &active
	player.BenchedPokemon = append([]models.Card(nil), chosen[1:]...)
	player.ProtectShields = rules.ProtectShieldsPerGame
	player.HasAttached = false
	player.HasDrawn = false
	player.HasSwitched = false
	player.HasPlayedPower = false
	player.AttackBonus = 0
	player.DefenseBonus = 0
	player.PartyReady = true
	player.DiscardPile = nil
	// Fresh power deck each game so both sides keep leverage.
	player.PowerDeck = e.buildPowerDeck(playerID)
	player.Hand = nil
	player.PendingDraw = nil
	e.State.LastAction = fmt.Sprintf("%s selected battle party", playerID)

	if e.bothPartiesReady() {
		e.beginBattle()
	} else {
		e.State.Phase = string(rules.PhasePartySelect)
		e.State.Status = models.StatusSetup
	}
	return nil
}

func (e *Engine) bothPartiesReady() bool {
	for i := range e.State.Players {
		if !e.State.Players[i].PartyReady {
			return false
		}
	}
	return len(e.State.Players) >= 2
}

func (e *Engine) beginBattle() {
	e.State.Status = models.StatusInProgress
	e.State.Phase = string(rules.PhaseInBattle)
	e.State.CurrentTurn = e.State.Players[0].ID
	e.State.TurnNumber = 1
	e.State.LastAction = fmt.Sprintf("game %d begin — %s to move", e.State.GameNumber, e.State.CurrentTurn)
	e.tryAutoDrawPower(e.State.CurrentTurn)
}

func (e *Engine) PlayBench(playerID, cardID string) error {
	return errors.New("benching from hand is not used — select your battle party of 3 (§3.1)")
}

func (e *Engine) SetActive(playerID, cardID string) error {
	player, err := e.requireTurn(playerID)
	if err != nil {
		return err
	}
	if player.ActivePokemon != nil {
		return errors.New("already have an active pokemon")
	}

	if idx := findCard(player.BenchedPokemon, cardID); idx >= 0 {
		card := player.BenchedPokemon[idx]
		player.BenchedPokemon = append(player.BenchedPokemon[:idx], player.BenchedPokemon[idx+1:]...)
		player.ActivePokemon = &card
		e.State.LastAction = fmt.Sprintf("%s set active %s", playerID, card.Name)
		return nil
	}

	if idx := findCard(player.Hand, cardID); idx >= 0 {
		card := player.Hand[idx]
		if card.Type != models.TypePokemon {
			return errors.New("can only set pokemon as active")
		}
		player.Hand = append(player.Hand[:idx], player.Hand[idx+1:]...)
		player.ActivePokemon = &card
		e.State.LastAction = fmt.Sprintf("%s set active %s", playerID, card.Name)
		return nil
	}
	return errors.New("card not found in hand or bench")
}

// Promote moves a benched Pokemon to active after a knockout.
func (e *Engine) Promote(playerID, cardID string) error {
	if e.State.Status != models.StatusInProgress {
		return errors.New("game not in progress")
	}
	player := e.getPlayer(playerID)
	if player == nil {
		return errors.New("player not found")
	}
	if player.ActivePokemon != nil {
		return errors.New("already have an active pokemon")
	}
	idx := findCard(player.BenchedPokemon, cardID)
	if idx < 0 {
		return errors.New("card not found on bench")
	}
	card := player.BenchedPokemon[idx]
	player.BenchedPokemon = append(player.BenchedPokemon[:idx], player.BenchedPokemon[idx+1:]...)
	player.ActivePokemon = &card
	e.State.LastAction = fmt.Sprintf("%s promoted %s", playerID, card.Name)
	e.tryAutoDrawPower(playerID)
	return nil
}

// Switch swaps the active Pokémon with a back-line Pokémon (once per turn).
func (e *Engine) Switch(playerID, cardID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon to switch out — promote instead")
	}
	if player.HasSwitched {
		return errors.New("already switched this turn")
	}
	idx := findCard(player.BenchedPokemon, cardID)
	if idx < 0 {
		return errors.New("card not found on back line")
	}

	incoming := player.BenchedPokemon[idx]
	outgoing := *player.ActivePokemon
	player.BenchedPokemon[idx] = outgoing
	player.ActivePokemon = &incoming
	player.HasSwitched = true
	e.State.LastAction = fmt.Sprintf("%s switched to %s", playerID, incoming.Name)
	return nil
}

func (e *Engine) AttachEnergy(playerID, energyCardID string) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.HasAttached {
		return errors.New("already charged energy this turn")
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon to charge")
	}
	// GO-style: one energy charge per turn on the active (approximates Fast Attack energy).
	// Optional energyCardID ignored — energy is not drawn from a deck in this format.
	_ = energyCardID
	player.ActivePokemon.EnergyAttached++
	player.HasAttached = true
	e.State.LastAction = fmt.Sprintf("%s charged energy on %s", playerID, player.ActivePokemon.Name)
	return nil
}

func (e *Engine) Attack(playerID string, attackIndex int) error {
	player, err := e.requirePlayable(playerID)
	if err != nil {
		return err
	}
	if player.ActivePokemon == nil {
		return errors.New("no active pokemon")
	}
	if attackIndex < 0 || attackIndex >= len(player.ActivePokemon.Attacks) {
		return errors.New("invalid attack index")
	}
	attack := player.ActivePokemon.Attacks[attackIndex]
	if player.ActivePokemon.EnergyAttached < attack.Cost {
		return fmt.Errorf("not enough energy (need %d)", attack.Cost)
	}

	opponent := e.getOpponent(playerID)
	if opponent == nil || opponent.ActivePokemon == nil {
		return errors.New("opponent has no active pokemon")
	}

	dmg := computeDamage(attack.Damage, player, opponent)
	opponent.ActivePokemon.HP -= dmg
	// Attack bonus is one-shot; defense persists until consumed by a hit.
	player.AttackBonus = 0
	if opponent.DefenseBonus > 0 {
		opponent.DefenseBonus = 0
	}
	e.State.LastAction = fmt.Sprintf("%s used %s for %d damage", playerID, attack.Name, dmg)

	if opponent.ActivePokemon.HP <= 0 {
		e.resolveKnockout(player, opponent)
	}

	if e.State.Status == models.StatusGameOver {
		return nil
	}

	return e.advanceTurn(playerID)
}

func (e *Engine) resolveKnockout(attacker, defender *models.PlayerState) {
	knocked := *defender.ActivePokemon
	defender.DiscardPile = append(defender.DiscardPile, knocked)
	defender.ActivePokemon = nil

	// Game ends when a competitor knocks out their opponent's final Pokémon (§6.3).
	if len(defender.BenchedPokemon) == 0 {
		e.finishGame(attacker.ID, fmt.Sprintf("%s wins game %d — opponent has no Pokémon left (§6.3)", attacker.ID, e.State.GameNumber))
		return
	}

	e.State.LastAction += fmt.Sprintf("; %s must promote", defender.ID)
}

func (e *Engine) finishGame(winnerID, message string) {
	winner := e.getPlayer(winnerID)
	if winner != nil {
		winner.GamesWon++
	}
	e.State.LastAction = message

	if matchWinner := rules.MatchWinnerID(e.State); matchWinner != "" {
		e.State.Status = models.StatusGameOver
		e.State.Phase = string(rules.PhaseMatchOver)
		e.State.Winner = matchWinner
		e.State.LastAction = fmt.Sprintf("%s wins the match (§6.4 best-of-%d)", matchWinner, e.State.WinsNeeded*2-1)
		return
	}

	// Prepare next game — team preview / party select (§6.1 between games).
	e.State.GameNumber++
	e.State.Status = models.StatusSetup
	e.State.Phase = string(rules.PhaseBetweenGames)
	e.State.CurrentTurn = ""
	e.State.TurnNumber = 0
	for i := range e.State.Players {
		p := &e.State.Players[i]
		p.ActivePokemon = nil
		p.BenchedPokemon = nil
		p.DiscardPile = nil
		p.Hand = nil
		p.PowerDeck = e.buildPowerDeck(p.ID)
		p.PartyReady = false
		p.ProtectShields = rules.ProtectShieldsPerGame
		p.HasAttached = false
		p.HasDrawn = false
		p.HasSwitched = false
		p.HasPlayedPower = false
		p.PendingDraw = nil
		p.AttackBonus = 0
		p.DefenseBonus = 0
	}
	e.State.LastAction += fmt.Sprintf("; select parties for game %d", e.State.GameNumber)
}

func (e *Engine) EndTurn(playerID string) error {
	if _, err := e.requirePlayable(playerID); err != nil {
		return err
	}
	e.State.LastAction = fmt.Sprintf("%s ended turn", playerID)
	return e.advanceTurn(playerID)
}

func (e *Engine) advanceTurn(currentPlayerID string) error {
	opponent := e.getOpponent(currentPlayerID)
	if opponent == nil {
		return errors.New("opponent not found")
	}
	e.State.CurrentTurn = opponent.ID
	e.State.TurnNumber++
	opponent.HasDrawn = false
	opponent.HasAttached = false
	opponent.HasSwitched = false
	opponent.HasPlayedPower = false
	e.tryAutoDrawPower(opponent.ID)
	return nil
}

func findCard(cards []models.Card, id string) int {
	for i, c := range cards {
		if c.ID == id {
			return i
		}
	}
	return -1
}

// Restore loads a persisted state into the engine (Memento restore).
func (e *Engine) Restore(state *models.GameState) {
	e.State = state
}
