import { Rules } from "./rules";

interface IDeck {
	shuffle(): TCGCard.CardType[];
	draw(): TCGCard.PlayerInHandCardLength<TCGCard.CardType[]>;
}

export class Deck implements IDeck {
	cards: TCGCard.CardType[];
	constructor(cards: TCGCard.CardType[]) {
		if (cards.length < 60) {
			throw new Error(Rules.deckLength);
		}
		this.cards = cards;
	}
	draw(): TCGCard.PlayerInHandCardLength<TCGCard.CardType[]> {
		throw new Error("Method not implemented.");
	}
	shuffle(): TCGCard.CardType[] {
		throw new Error("Method not implemented.");
	}
}
