// seed.js

const mongoose = require("mongoose");
const pokemonSeedData = require("./data");

const uri = "mongodb://localhost:27017/pokemonDB"; // Replace with your MongoDB URI

// Define the schema
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

// Define the model
const PokemonCard = mongoose.model("PokemonCard", pokemonSchema);

// Seed function
async function seedDatabase() {
	try {
		await mongoose.connect(uri);
		console.log("Connected to MongoDB");

		// Clear existing data
		await PokemonCard.deleteMany({});
		console.log("Existing data cleared");

		// Insert sample data
		await PokemonCard.insertMany(pokemonSeedData);
		console.log("Database seeded with Pokémon cards");

		mongoose.connection.close();
	} catch (error) {
		console.error("Error seeding database:", error);
		mongoose.connection.close();
	}
}

// Run the seed function
seedDatabase();
