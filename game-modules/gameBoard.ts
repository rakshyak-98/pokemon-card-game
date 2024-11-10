interface IBoard {
	activePokemon: TCGCard.BasicPokemon | null;
	bench: TCGCard.Bench[] | null;
	deck: TCGCard.Deck;
	prizeCards: TCGCard.PriceCards | [];
	discardPile: TCGCard.DiscardPile | [];
}

export class Board implements IBoard {
	activePokemon = null;
	bench = null;
	deck;
	discardPile = [];
	prizeCards = [];
	constructor(deck: TCGCard.Deck) {
		this.deck = deck;
	}
}
