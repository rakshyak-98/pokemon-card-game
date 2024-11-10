interface IGame {
	activePokemon: TCGCard.BasicPokemon;
	deck: TCGCard.BoardDeck[];
	bench: TCGCard.Bench[];
	prizeCards: TCGCard.PriceCards;
	discardPile: TCGCard.DiscardPile;
}
