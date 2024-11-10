import type { DescriptionType } from "./card-description";
import type { RulesType } from "./rules";

type SymbolType = {
	typeName: string;
	symbolType: keyof DescriptionType;
};

type PokemonType = {};
type TrainerType = {};
type CardType = {};

interface PokemonCard {
	pokemonType: SymbolType;
	hp: number;
	cardName: string;
	stage: string;
	evolvesFromPokemon: PokemonType;
	expansionCode: number;
	collectorCardNumber: number;
}

interface TrainerCard {
	cardName: string;
	trainerType: TrainerType;
	cardType: CardType;
	textBox: string;
	trainerRule: keyof RulesType;
}

export class Card implements PokemonCard, TrainerCard {
	pokemonType: SymbolType;
	hp: number;
	cardName: string;
	stage: string;
	evolvesFromPokemon: PokemonType;
	expansionCode: number;
	collectorCardNumber: number;
	trainerType: TrainerType;
	cardType: CardType;
	textBox: string;
	trainerRule: keyof RulesType;
}
