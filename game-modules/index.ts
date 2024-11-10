import { Deck } from "./cards/deck";
import { Board } from "./gameBoard";

function main() {
	const deck = new Deck();
	const board = new Board(deck);
	const player = new Player();
	player.drawFromDeck(deck);
	console.log(player.getInHandCards());
}

main();
