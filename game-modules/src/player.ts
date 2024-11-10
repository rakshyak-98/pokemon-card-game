interface IPlayer {
	inHandCards: TCGCard.CardType[];
	shuffleDeck(): void;
	drawFromDeck(deck: TCGCard.Deck): void;
	placeCardsOnBoard(): void;
	putCardOnActiveSport(card: TCGCard.BasicPokemon): void;
	putCardsInBench(card: TCGCard.Bench): void;
	drawFromPrizeCards(): void;
	discardCard(card: TCGCard.CardType): void;
	moveCardFromBenchToActiveSpot(index: number): void;
	getInHandCards(): TCGCard.CardType[];
}

class Player implements IPlayer {
	getInHandCards(): TCGCard.CardType[] {
		return this.inHandCards;
	}
	inHandCards: TCGCard.CardType[] = [];
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
	drawFromDeck(deck: TCGCard.Deck) {
		this.inHandCards = deck.shuffle().draw(7);
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
