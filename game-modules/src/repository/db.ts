import mongoose from "mongoose";

try {
	const uri = "mongodb://localhost:27017/pokemonDB";
	mongoose.connect(uri);
} catch (error) {
	console.log(error);
}

const pokemonSchema = new mongoose.Schema({
	name: String,
	type: String,
	hp: Number,
	attacks: [
		{
			name: String,
			damage: Number,
		},
	],
	rarity: String,
	weakness: String,
	resistance: String,
	retreatCost: Number,
});

export class Repository {
	PokemonCard = mongoose.model("PokemonCard", pokemonSchema);

	async getCards() {
		const data = await this.PokemonCard.find().select("-_id -__v");
		data.map((d) => {
			new Card();
		});
	}
}
