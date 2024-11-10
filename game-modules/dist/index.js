"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const deck_1 = require("./cards/deck");
const gameBoard_1 = require("./gameBoard");
function main() {
    const deck = new deck_1.Deck();
    const board = new gameBoard_1.Board(deck);
    const player = new Player();
    player.drawFromDeck(deck);
    console.log(player.getInHandCards());
}
main();
