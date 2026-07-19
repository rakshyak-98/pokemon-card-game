package game

import (
	"fmt"
	"math/rand"
	"time"

	"rakshyak-98/pokemon-backend/models"
	"rakshyak-98/pokemon-backend/rules"
)

// Power card effect kinds — enhance attack, defense, or heal.
const (
	EffectBoostAttack  = "boost_attack"
	EffectBoostDefense = "boost_defense"
	EffectHeal         = "heal"
)

// PowerDeckSize is how many special power cards each competitor starts with
// (and how many are added on each mid-game refill).
const PowerDeckSize = 12

// MaxPowerHandSlots — each player may hold at most this many power cards (empty slots to fill).
const MaxPowerHandSlots = rules.MaxPowerHandSlots

// Power card effect magnitudes used by the built-in fallback catalog.
const (
	PowerAttackBonus  = 20
	PowerDefenseBonus = 15
	PowerHealAmount   = 30
)

// fallbackPowerCatalog is used when the DB / PokeAPI power seed is empty.
// Names lean on ASC / MEW Trainer Item and healing themes (Potion, X Attack, berries, etc.).
var fallbackPowerCatalog = []models.PowerCard{
	{
		PokeAPIID: 57, Name: "X Attack", Effect: EffectBoostAttack, EffectValue: PowerAttackBonus,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/x-attack.png",
		Category: "stat-boosts", Description: "Raises Attack for the next hit.",
	},
	{
		PokeAPIID: 58, Name: "X Defense", Effect: EffectBoostDefense, EffectValue: PowerDefenseBonus,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/x-defense.png",
		Category: "stat-boosts", Description: "Raises Defense against the next hit.",
	},
	{
		PokeAPIID: 17, Name: "Potion", Effect: EffectHeal, EffectValue: PowerHealAmount,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/potion.png",
		Category: "healing", Description: "Restores HP to the Active Pokémon.",
	},
	{
		PokeAPIID: 18, Name: "Super Potion", Effect: EffectHeal, EffectValue: 35,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/super-potion.png",
		Category: "healing", Description: "A stronger Potion that restores more HP.",
	},
	{
		PokeAPIID: 19, Name: "Hyper Potion", Effect: EffectHeal, EffectValue: 50,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/hyper-potion.png",
		Category: "healing", Description: "Greatly restores HP to the Active Pokémon.",
	},
	{
		PokeAPIID: 55, Name: "Guard Spec", Effect: EffectBoostDefense, EffectValue: 10,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/guard-spec.png",
		Category: "stat-boosts", Description: "Temporarily hardens defenses.",
	},
	{
		PokeAPIID: 56, Name: "Dire Hit", Effect: EffectBoostAttack, EffectValue: 15,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/dire-hit.png",
		Category: "stat-boosts", Description: "Sharpens focus for a stronger hit.",
	},
	{
		PokeAPIID: 126, Name: "Oran Berry", Effect: EffectHeal, EffectValue: 20,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/oran-berry.png",
		Category: "medicine", Description: "A restorative berry that heals a little HP.",
	},
	{
		PokeAPIID: 158, Name: "Sitrus Berry", Effect: EffectHeal, EffectValue: 40,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/sitrus-berry.png",
		Category: "medicine", Description: "A restorative berry that heals a solid amount of HP.",
	},
	{
		PokeAPIID: 266, Name: "Muscle Band", Effect: EffectBoostAttack, EffectValue: 25,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/muscle-band.png",
		Category: "held-items", Description: "A Trainer tool that boosts physical striking power.",
	},
	{
		PokeAPIID: 640, Name: "Assault Vest", Effect: EffectBoostDefense, EffectValue: 25,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/assault-vest.png",
		Category: "held-items", Description: "A Trainer tool that hardens the Active Pokémon.",
	},
	{
		PokeAPIID: 234, Name: "Leftovers", Effect: EffectHeal, EffectValue: 25,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/leftovers.png",
		Category: "held-items", Description: "Gradually restores HP to the Active Pokémon.",
	},
}

// SetPowerCatalog updates the special power-card pool used when building decks.
func (e *Engine) SetPowerCatalog(catalog []models.PowerCard) {
	e.powerCatalog = catalog
}

func (e *Engine) powerPool() []models.PowerCard {
	if len(e.powerCatalog) > 0 {
		return e.powerCatalog
	}
	return fallbackPowerCatalog
}

// buildPowerDeck creates a shuffled special-power deck for a player from the catalog.
// Cards are drawn with roughly equal representation of attack / defense / heal effects.
func (e *Engine) buildPowerDeck(playerID string) []models.Card {
	pool := e.powerPool()
	byEffect := map[string][]models.PowerCard{}
	for _, p := range pool {
		byEffect[p.Effect] = append(byEffect[p.Effect], p)
	}
	effects := []string{EffectBoostAttack, EffectBoostDefense, EffectHeal}
	// Prefer effect kinds that actually exist in the catalog.
	available := make([]string, 0, len(effects))
	for _, eff := range effects {
		if len(byEffect[eff]) > 0 {
			available = append(available, eff)
		}
	}
	if len(available) == 0 {
		available = effects
		byEffect = map[string][]models.PowerCard{
			EffectBoostAttack:  {fallbackPowerCatalog[0]},
			EffectBoostDefense: {fallbackPowerCatalog[1]},
			EffectHeal:         {fallbackPowerCatalog[2]},
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	deck := make([]models.Card, 0, PowerDeckSize)
	for i := 0; i < PowerDeckSize; i++ {
		eff := available[i%len(available)]
		choices := byEffect[eff]
		tpl := choices[r.Intn(len(choices))]
		e.powerSeq++
		deck = append(deck, models.Card{
			ID:          fmt.Sprintf("%s-power-%d", playerID, e.powerSeq),
			Name:        tpl.Name,
			Type:        models.TypePower,
			Effect:      tpl.Effect,
			EffectValue: tpl.EffectValue,
			ImageURL:    tpl.ImageURL,
			PokeAPIID:   tpl.PokeAPIID,
		})
	}
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// ensurePowerDeck refills the player's special-power deck when empty so draws
// continue every turn until the game ends.
func (e *Engine) ensurePowerDeck(player *models.PlayerState) {
	if player == nil || len(player.PowerDeck) > 0 {
		return
	}
	player.PowerDeck = e.buildPowerDeck(player.ID)
}

// drawPowerCard moves the top card from PowerDeck into Hand, or into PendingDraw
// when all power slots are full (player must swap or discard).
func drawPowerCard(player *models.PlayerState) (string, bool, error) {
	if len(player.PowerDeck) == 0 {
		return "", false, fmt.Errorf("power deck is empty")
	}
	card := player.PowerDeck[0]
	player.PowerDeck = player.PowerDeck[1:]
	if len(player.Hand) >= MaxPowerHandSlots {
		player.PendingDraw = []models.Card{card}
		return card.Name, true, nil
	}
	player.Hand = append(player.Hand, card)
	return card.Name, false, nil
}

// applyPowerEffect resolves a played power card onto the active Pokémon / combat bonuses.
func applyPowerEffect(player *models.PlayerState, card models.Card) (string, error) {
	if player.ActivePokemon == nil {
		return "", fmt.Errorf("no active Pokémon to apply power to")
	}

	switch card.Effect {
	case EffectBoostAttack:
		player.AttackBonus += card.EffectValue
		return fmt.Sprintf("attack +%d (next hit)", card.EffectValue), nil

	case EffectBoostDefense:
		player.DefenseBonus += card.EffectValue
		return fmt.Sprintf("defense +%d (incoming)", card.EffectValue), nil

	case EffectHeal:
		if player.ActivePokemon.HP >= player.ActivePokemon.MaxHP {
			return "", fmt.Errorf("%s is already at full HP", player.ActivePokemon.Name)
		}
		before := player.ActivePokemon.HP
		player.ActivePokemon.HP += card.EffectValue
		if player.ActivePokemon.HP > player.ActivePokemon.MaxHP {
			player.ActivePokemon.HP = player.ActivePokemon.MaxHP
		}
		healed := player.ActivePokemon.HP - before
		return fmt.Sprintf("healed %s for %d HP", player.ActivePokemon.Name, healed), nil

	default:
		return "", fmt.Errorf("unknown power effect %q", card.Effect)
	}
}

// computeDamage applies attack/defense power bonuses for a combat hit.
func computeDamage(baseDamage int, attacker, defender *models.PlayerState) int {
	dmg := baseDamage + attacker.AttackBonus
	dmg -= defender.DefenseBonus
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}
