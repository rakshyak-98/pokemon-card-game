import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import React, { useState } from "react";
import PokemonDialog from "./PokemonDialog";

// Define the structure of a Pokémon card
interface PokemonCard {
	id: number;
	name: string;
	type: string;
	hp: number;
	image: string;
	description?: string;
	height?: number;
	weight?: number;
	abilities?: string[];
}

// Mock data for Pokémon cards
const pokemonCards: PokemonCard[] = [
	{
		id: 1,
		name: "Pikachu",
		type: "Electric",
		hp: 60,
		image: "/placeholder.svg?height=100&width=100",
		description: "A mouse-like PokÃ©mon that can generate electricity.",
		height: 0.4,
		weight: 6.0,
		abilities: ["Static", "Lightning Rod"],
	},
	{
		id: 2,
		name: "Charizard",
		type: "Fire/Flying",
		hp: 120,
		image: "/placeholder.svg?height=100&width=100",
		description:
			"A dragon-like PokÃ©mon that spits fire hot enough to melt boulders.",
		height: 1.7,
		weight: 90.5,
		abilities: ["Blaze", "Solar Power"],
	},
	{
		id: 3,
		name: "Bulbasaur",
		type: "Grass/Poison",
		hp: 45,
		image: "/placeholder.svg?height=100&width=100",
		description:
			"A strange seed was planted on its back at birth. The plant sprouts and grows with this PokÃ©mon.",
		height: 0.7,
		weight: 6.9,
		abilities: ["Overgrow", "Chlorophyll"],
	},
	{
		id: 4,
		name: "Squirtle",
		type: "Water",
		hp: 44,
		image: "/placeholder.svg?height=100&width=100",
		description: "A turtle-like PokÃ©mon that can squirt water from its mouth.",
		height: 0.5,
		weight: 9.0,
		abilities: ["Torrent", "Rain Dish"],
	},
	{
		id: 5,
		name: "Mewtwo",
		type: "Psychic",
		hp: 150,
		image: "/placeholder.svg?height=100&width=100",
		description:
			"A genetically engineered PokÃ©mon with powerful psychic abilities.",
		height: 2.0,
		weight: 122.0,
		abilities: ["Pressure", "Unnerve"],
	},
	{
		id: 6,
		name: "Gengar",
		type: "Ghost/Poison",
		hp: 100,
		image: "/placeholder.svg?height=100&width=100",
		description:
			"A shadowy PokÃ©mon that is said to be the animated soul of a deceased being.",
		height: 1.5,
		weight: 40.5,
		abilities: ["Cursed Body"],
	},
];

const Component: React.FC = () => {
	const [searchTerm, setSearchTerm] = useState<string>("");
	const [selectedPokemon, setSelectedPokemon] = useState<PokemonCard | null>(
		null
	);

	const filteredCards = pokemonCards.filter((card) =>
		card.name.toLowerCase().includes(searchTerm.toLowerCase())
	);

	const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setSearchTerm(event.target.value);
	};

	const openModal = (pokemon: PokemonCard) => {
		setSelectedPokemon(pokemon);
	};

	const closeModal = () => {
		setSelectedPokemon(null);
	};

	return (
		<div className="container mx-auto p-4">
			<h1 className="text-3xl font-bold mb-6 text-center">
				My Pokémon Collection
			</h1>

			<Input
				type="text"
				placeholder="Search Pokémon..."
				value={searchTerm}
				onChange={handleSearchChange}
				className="mb-6"
			/>

			<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
				{filteredCards.map((card: PokemonCard) => (
					<Card
						key={card.id}
						className="overflow-hidden"
						onClick={() => openModal(card)}
					>
						<CardHeader className="p-4">
							<CardTitle className="text-lg">{card.name}</CardTitle>
						</CardHeader>
						<CardContent>
							<img
								src={card.image}
								alt={card.name}
								className="w-full h-40 object-contain mb-2"
							/>
							<p className="text-sm text-gray-600">Type: {card.type}</p>
							<p className="text-sm text-gray-600">HP: {card.hp}</p>
						</CardContent>
					</Card>
				))}
			</div>

			{filteredCards.length === 0 && (
				<p className="text-center text-gray-500 mt-4">
					No Pokémon found. Try a different search term.
				</p>
			)}
			<PokemonDialog
				closeModal={closeModal}
				selectedPokemon={selectedPokemon}
			/>
		</div>
	);
};

export default Component;
