"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Deck = void 0;
const rules_1 = require("./rules");
class Deck {
    constructor() {
        this.cards = [];
        if (this.cards.length < 60) {
            throw new Error(rules_1.Rules.deckLength);
        }
    }
    shuffle() {
        this.cards.sort(() => Math.random() - 0.5);
        return this;
    }
    draw(n) {
        return this.cards.splice(0, n);
    }
}
exports.Deck = Deck;
