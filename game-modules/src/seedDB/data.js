// sampleData.js

const pokemonCards = [
  {
    name: "Pikachu",
    type: "Electric",
    hp: 60,
    attacks: [
      { name: "Thunder Shock", damage: 30 },
      { name: "Electro Ball", damage: 50 },
    ],
    rarity: "Common",
    weakness: "Ground",
    resistance: "Flying",
    retreatCost: 1,
  },
  {
    name: "Charmander",
    type: "Fire",
    hp: 50,
    attacks: [
      { name: "Scratch", damage: 10 },
      { name: "Flame Tail", damage: 30 },
    ],
    rarity: "Common",
    weakness: "Water",
    resistance: "Grass",
    retreatCost: 1,
  },
  {
    name: "Squirtle",
    type: "Water",
    hp: 60,
    attacks: [
      { name: "Bubble", damage: 10 },
      { name: "Aqua Tail", damage: 20 },
    ],
    rarity: "Common",
    weakness: "Electric",
    resistance: "Fire",
    retreatCost: 1,
  },
  {
    name: "Bulbasaur",
    type: "Grass",
    hp: 70,
    attacks: [
      { name: "Vine Whip", damage: 20 },
      { name: "Seed Bomb", damage: 40 },
    ],
    rarity: "Common",
    weakness: "Fire",
    resistance: "Water",
    retreatCost: 1,
  },
  {
    name: "Jigglypuff",
    type: "Fairy",
    hp: 80,
    attacks: [
      { name: "Sing", damage: 0 },
      { name: "Double Slap", damage: 40 },
    ],
    rarity: "Uncommon",
    weakness: "Steel",
    resistance: "Dark",
    retreatCost: 1,
  },
  {
    name: "Gengar",
    type: "Ghost",
    hp: 120,
    attacks: [
      { name: "Shadow Punch", damage: 60 },
      { name: "Dark Corridor", damage: 90 },
    ],
    rarity: "Rare",
    weakness: "Psychic",
    resistance: "Normal",
    retreatCost: 2,
  },
  {
    name: "Dragonite",
    type: "Dragon",
    hp: 160,
    attacks: [
      { name: "Dragon Claw", damage: 80 },
      { name: "Hyper Beam", damage: 150 },
    ],
    rarity: "Legendary",
    weakness: "Fairy",
    resistance: "Grass",
    retreatCost: 3,
  },
  {
    name: "Magikarp",
    type: "Water",
    hp: 30,
    attacks: [
      { name: "Splash", damage: 0 },
    ],
    rarity: "Common",
    weakness: "Electric",
    resistance: "Fire",
    retreatCost: 1,
  },
  {
    name: "Snorlax",
    type: "Normal",
    hp: 150,
    attacks: [
      { name: "Body Slam", damage: 70 },
      { name: "Heavy Impact", damage: 100 },
    ],
    rarity: "Rare",
    weakness: "Fighting",
    resistance: "Psychic",
    retreatCost: 4,
  },
  {
    name: "Eevee",
    type: "Normal",
    hp: 50,
    attacks: [
      { name: "Tackle", damage: 10 },
      { name: "Quick Attack", damage: 20 },
    ],
    rarity: "Common",
    weakness: "Fighting",
    resistance: "Ghost",
    retreatCost: 1,
  },
];

// Extend this list to 60 cards by duplicating and modifying attributes.
module.exports = pokemonCards;
