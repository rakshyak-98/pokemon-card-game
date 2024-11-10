"use strict";
class Player {
    constructor() {
        this.inHandCards = [];
    }
    getInHandCards() {
        return this.inHandCards;
    }
    putCardOnActiveSport(card) {
        throw new Error("Method not implemented.");
    }
    putCardsInBench(card) {
        throw new Error("Method not implemented.");
    }
    moveCardFromBenchToActiveSpot(index) {
        throw new Error("Method not implemented.");
    }
    shuffleDeck() {
        throw new Error("Method not implemented.");
    }
    drawFromDeck(deck) {
        this.inHandCards = deck.shuffle().draw(7);
    }
    placeCardsOnBoard() {
        throw new Error("Method not implemented.");
    }
    drawFromPrizeCards() {
        throw new Error("Method not implemented.");
    }
    discardCard(card) {
        throw new Error("Method not implemented.");
    }
}
