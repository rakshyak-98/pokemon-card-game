interface IPlayer {
	inHandCards: TCGCard.PlayerInHandCardLength<TCGCard.CardType>[];
	shuffleDeck(): void;
	drawFromDeck(): TCGCard.PlayerInHandCardLength<TCGCard.CardType[]>;
	placeCardsOnBoard(): void;
	putCardOnActiveSport(card: TCGCard.BasicPokemon): void;
	putCardsInBench(card: TCGCard.Bench): void;
	drawFromPrizeCards(): void;
	discardCard(card: TCGCard.CardType): void;
	moveCardFromBenchToActiveSpot(index: number): void;
}

class Player implements IPlayer {
	inHandCards: TCGCard.PlayerInHandCardLength<TCGCard.CardType>[];
	constructor(
		inHandCards: TCGCard.PlayerInHandCardLength<TCGCard.CardType>[],
		deck: TCGCard.Deck
	) {
		this.inHandCards = inHandCards;
	}
	putCardOnActiveSport(card: TCGCard.BasicPokemon): void {
		throw new Error("Method not implemented.");
	}
	putCardsInBench(card: TCGCard.Bench): void {
		throw new Error("Method not implemented.");
	}
	moveCardFromBenchToActiveSpot(index: number): void {
		throw new Error("Method not implemented.");
	}
	shuffleDeck(): void {
		throw new Error("Method not implemented.");
	}
	drawFromDeck(): TCGCard.PlayerInHandCardLength<TCGCard.CardType[]> {
		throw new Error("Method not implemented.");
	}
	placeCardsOnBoard(): void {
		throw new Error("Method not implemented.");
	}
	drawFromPrizeCards(): void {
		throw new Error("Method not implemented.");
	}
	discardCard(card: TCGCard.CardType): void {
		throw new Error("Method not implemented.");
	}
}
