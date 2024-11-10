"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Board = void 0;
class Board {
    constructor(deck) {
        this.activePokemon = null;
        this.bench = null;
        this.discardPile = [];
        this.prizeCards = [];
        this.deck = deck;
    }
}
exports.Board = Board;
