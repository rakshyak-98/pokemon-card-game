import mongoose from "mongoose";
import { Repository } from "./repository/db";

async function main() {
	const repo = new Repository();
	try {
		const data = await repo.getCards();
		console.log(data);
		mongoose.connection.close();
	} finally {
	}
	// const board = new Board(deck);
	// const player = new Player();
	// player.drawFromDeck(deck);
}

main();
