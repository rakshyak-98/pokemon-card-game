import { Rules } from "./rules";

export interface IDeck {
	shuffle(): TCGCard.CardType[];
}

export class Deck implements IDeck {
	cards: TCGCard.CardType[];
	constructor(cards: TCGCard.CardType[]) {
		if (cards.length < 60) {
			throw new Error(Rules.deckLength);
		}
		this.cards = cards;
	}
	shuffle(): TCGCard.CardType[] {
		throw new Error("Method not implemented.");
	}
}
