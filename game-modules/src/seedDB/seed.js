// seed.js

const mongoose = require("mongoose");
const pokemonCards = require("./data");

const uri = "mongodb://localhost:27017/pokemonDB"; // Replace with your MongoDB URI

// Define the schema
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
		await PokemonCard.insertMany(pokemonCards);
		console.log("Database seeded with Pokémon cards");

		mongoose.connection.close();
	} catch (error) {
		console.error("Error seeding database:", error);
		mongoose.connection.close();
	}
}

// Run the seed function
seedDatabase();
