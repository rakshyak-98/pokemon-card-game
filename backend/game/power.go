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

// PowerDeckSize is how many special power cards each competitor starts with.
const PowerDeckSize = 12

// MaxPowerHandSlots — each player may hold at most this many power cards (empty slots to fill).
const MaxPowerHandSlots = rules.MaxPowerHandSlots

// Power card effect magnitudes (balanced for GO-style HP / damage ranges).
const (
	PowerAttackBonus  = 20
	PowerDefenseBonus = 15
	PowerHealAmount   = 30
)

type powerTemplate struct {
	Name        string
	Effect      string
	EffectValue int
	ImageURL    string
}

var powerTemplates = []powerTemplate{
	{
		Name: "Power Strike", Effect: EffectBoostAttack, EffectValue: PowerAttackBonus,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/x-attack.png",
	},
	{
		Name: "Iron Guard", Effect: EffectBoostDefense, EffectValue: PowerDefenseBonus,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/x-defense.png",
	},
	{
		Name: "Potion Heal", Effect: EffectHeal, EffectValue: PowerHealAmount,
		ImageURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/potion.png",
	},
}

// buildPowerDeck creates a shuffled special-power deck for a player (attack / defense / heal).
func buildPowerDeck(playerID string) []models.Card {
	deck := make([]models.Card, 0, PowerDeckSize)
	for i := 0; i < PowerDeckSize; i++ {
		tpl := powerTemplates[i%len(powerTemplates)]
		deck = append(deck, models.Card{
			ID:          fmt.Sprintf("%s-power-%d", playerID, i),
			Name:        tpl.Name,
			Type:        models.TypePower,
			Effect:      tpl.Effect,
			EffectValue: tpl.EffectValue,
			ImageURL:    tpl.ImageURL,
		})
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// drawPowerCard moves the top card from PowerDeck into Hand. Returns the drawn card name.
func drawPowerCard(player *models.PlayerState) (string, error) {
	if len(player.Hand) >= MaxPowerHandSlots {
		return "", fmt.Errorf("power hand is full (%d slots)", MaxPowerHandSlots)
	}
	if len(player.PowerDeck) == 0 {
		return "", fmt.Errorf("power deck is empty")
	}
	card := player.PowerDeck[0]
	player.PowerDeck = player.PowerDeck[1:]
	player.Hand = append(player.Hand, card)
	return card.Name, nil
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
