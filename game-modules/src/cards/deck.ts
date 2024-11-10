import { Rules } from "./rules";

export interface IDeck {
	shuffle(): this;
	draw(n: number): TCGCard.CardType[];
}

export class Deck implements IDeck {
	cards: TCGCard.CardType[];
	constructor() {
		this.cards = [];
		if (this.cards.length < 60) {
			throw new Error(Rules.deckLength);
		}
	}
	shuffle() {
		this.cards.sort(() => Math.random() - 0.5);
		return this;
	}

	draw(n: number) {
		return this.cards.splice(0, n);
	}
}
