import { EnergyCard, PokemonCard, TrainerCard } from "./cards/card-types";

declare global {
	namespace TCGCard {
		type PlayerInHandCardLength<T> = T[] & { length: 7 };
		type CardType = PokemonCard | TrainerCard | EnergyCard;
	}
}

export { };
