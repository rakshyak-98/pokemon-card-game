import mongoose from "mongoose";

try {
	const uri = "mongodb://localhost:27017/pokemonDB";
	mongoose.connect(uri);
} catch (error) {
	console.log(error);
}

const pokemonSchema = new mongoose.Schema(
	{
		pokemonType: {
			type: String,
			enum: [
				"Fire",
				"Water",
				"Grass",
				"Electric",
				"Fighting",
				"Psychic",
				"Ghost",
				"Dragon",
			],
			required: true,
		},
		hp: {
			type: Number,
			required: true,
			min: 1,
		},
		cardName: {
			type: String,
			required: true,
		},
		stage: {
			type: String,
			enum: ["Basic", "Stage 1", "Stage 2", "Legendary"],
			required: true,
		},
		evolvesFromPokemon: {
			type: String,
			required: true,
		},
		expansionCode: {
			type: Number,
			required: true,
		},
		collectorCardNumber: {
			type: Number,
			required: true,
		},
	},
	{
		timestamps: true, // Automatically add createdAt and updatedAt timestamps
	}
);

export class Repository {
	PokemonCard = mongoose.model("PokemonCard", pokemonSchema);

	async getCards() {
		const data = await this.PokemonCard.find().select("-_id -__v");
		data.map((d) => {
			const data = d.toObject();
			const {} = data;
		});
	}
}
