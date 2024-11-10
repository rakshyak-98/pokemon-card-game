import { EnergyCard, PokemonCard, TrainerCard } from "./cards/card-types";

declare global {
	namespace TCGCard {
		type PlayerInHandCardLength<T> = T[] & { length: 7 };
		type CardType = PokemonCard | TrainerCard | EnergyCard;
		type BasicPokemon = PokemonCard;
		type BoardDeck = CardType[] & { length: 60 };
		type Bench = CardType[] & { length: 6 };
		type PriceCards = CardType[] & { length: 6 };
		type DiscardPile = CardType[];
	}
}

export { };
