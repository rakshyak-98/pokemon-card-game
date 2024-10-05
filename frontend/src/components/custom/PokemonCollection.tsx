import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import React, { useState } from 'react';

// Define the structure of a Pokémon card
interface PokemonCard {
  id: number;
  name: string;
  type: string;
  hp: number;
  image: string;
}

// Mock data for Pokémon cards
const pokemonCards: PokemonCard[] = [
  { id: 1, name: "Pikachu", type: "Electric", hp: 60, image: "/placeholder.svg?height=100&width=100" },
  { id: 2, name: "Charizard", type: "Fire", hp: 120, image: "/placeholder.svg?height=100&width=100" },
  { id: 3, name: "Bulbasaur", type: "Grass", hp: 45, image: "/placeholder.svg?height=100&width=100" },
  { id: 4, name: "Squirtle", type: "Water", hp: 44, image: "/placeholder.svg?height=100&width=100" },
  { id: 5, name: "Mewtwo", type: "Psychic", hp: 150, image: "/placeholder.svg?height=100&width=100" },
  { id: 6, name: "Gengar", type: "Ghost", hp: 100, image: "/placeholder.svg?height=100&width=100" },
]

const Component: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState<string>('')

  const filteredCards = pokemonCards.filter(card =>
    card.name.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value)
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6 text-center">My Pokémon Collection</h1>

      <Input
        type="text"
        placeholder="Search Pokémon..."
        value={searchTerm}
        onChange={handleSearchChange}
        className="mb-6"
      />

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {filteredCards.map((card: PokemonCard) => (
          <Card key={card.id} className="overflow-hidden">
            <CardHeader className="p-4">
              <CardTitle className="text-lg">{card.name}</CardTitle>
            </CardHeader>
            <CardContent>
              <img src={card.image} alt={card.name} className="w-full h-40 object-contain mb-2" />
              <p className="text-sm text-gray-600">Type: {card.type}</p>
              <p className="text-sm text-gray-600">HP: {card.hp}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      {filteredCards.length === 0 && (
        <p className="text-center text-gray-500 mt-4">No Pokémon found. Try a different search term.</p>
      )}
    </div>
  )
}

export default Component