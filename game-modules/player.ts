interface IPlayer {
	shuffleDeck(): void;
	drawFromDeck(): void;
	inHandCards: TCGCard.PlayerInHandCardLength<TCGCard.CardType>[];
	placeCardsOnBoard: () => void;
}

class Player {}
