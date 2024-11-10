import type { DescriptionType } from "./card-description";
import type { RulesType } from "./rules";

type SymbolType = {
	typeName: string;
	symbolType: keyof DescriptionType;
};

type PokemonType = {};
type TrainerType = {};
type CardType = TCGCard.CardType;

interface IPokemonCard {
	pokemonType: SymbolType;
	hp: number;
	cardName: string;
	stage: string;
	evolvesFromPokemon: PokemonType;
	expansionCode: number;
	collectorCardNumber: number;
}

interface ITrainerCard {
	cardName: string;
	trainerType: TrainerType;
	cardType: CardType;
	textBox: string;
	trainerRule: keyof RulesType;
}

interface IEnergyCard {}

export class PokemonCard implements IPokemonCard {
	pokemonType: SymbolType;
	hp: number;
	cardName: string;
	stage: string;
	evolvesFromPokemon: PokemonType;
	expansionCode: number;
	collectorCardNumber: number;
	constructor(
		pokemonType: SymbolType,
		hp: number,
		cardName: string,
		stage: string,
		evolvesFromPokemon: PokemonType,
		expansionCode: number,
		collectorCardNumber: number
	) {
		this.cardName = cardName;
		this.pokemonType = pokemonType;
		this.hp = hp;
		this.stage = stage;
		this.evolvesFromPokemon = evolvesFromPokemon;
		this.expansionCode = expansionCode;
		this.collectorCardNumber = collectorCardNumber;
	}
}

export class TrainerCard implements ITrainerCard {
	cardName: string;
	trainerType: TrainerType;
	cardType: CardType;
	textBox: string;
	trainerRule: "describe";
	constructor(
		cardName: string,
		trainerType: TrainerType,
		cardType: CardType,
		textBox: string,
		trainerRule: "describe"
	) {
		this.cardName = cardName;
		this.trainerRule = trainerRule;
		this.trainerType = trainerType;
		this.cardType = cardType;
		this.textBox = textBox;
	}
}

export class EnergyCard implements IEnergyCard {}
