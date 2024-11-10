interface IBoard {
	activePokemon: TCGCard.BasicPokemon;
	bench: TCGCard.Bench[];
	deck: TCGCard.Deck;
	prizeCards: TCGCard.PriceCards;
	discardPile: TCGCard.DiscardPile;
}

export class Board implements IBoard {
	activePokemon: TCGCard.BasicPokemon;
	bench: TCGCard.Bench[];
	deck: TCGCard.Deck;
	prizeCards: TCGCard.PriceCards;
	discardPile: TCGCard.DiscardPile;
	constructor(
		activePokemon: TCGCard.BasicPokemon,
		bench: TCGCard.Bench[],
		prizeCards: TCGCard.PriceCards,
		discardPile: TCGCard.DiscardPile,
		deck: TCGCard.Deck
	) {
		this.activePokemon = activePokemon;
		this.prizeCards = prizeCards;
		this.bench = bench;
		this.discardPile = discardPile;
		this.deck = deck;
	}
}
