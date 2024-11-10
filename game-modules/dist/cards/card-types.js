"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EnergyCard = exports.TrainerCard = exports.PokemonCard = void 0;
class PokemonCard {
    constructor(pokemonType, hp, cardName, stage, evolvesFromPokemon, expansionCode, collectorCardNumber) {
        this.cardName = cardName;
        this.pokemonType = pokemonType;
        this.hp = hp;
        this.stage = stage;
        this.evolvesFromPokemon = evolvesFromPokemon;
        this.expansionCode = expansionCode;
        this.collectorCardNumber = collectorCardNumber;
    }
}
exports.PokemonCard = PokemonCard;
class TrainerCard {
    constructor(cardName, trainerType, cardType, textBox, trainerRule) {
        this.cardName = cardName;
        this.trainerRule = trainerRule;
        this.trainerType = trainerType;
        this.cardType = cardType;
        this.textBox = textBox;
    }
}
exports.TrainerCard = TrainerCard;
class EnergyCard {
}
exports.EnergyCard = EnergyCard;
