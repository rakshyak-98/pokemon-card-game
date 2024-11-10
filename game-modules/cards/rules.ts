const Rules = {
	describe: {
		zonesOfTCG: {
			prizeCards:
				"Each player has their own Prize cards. Prize cards are 6 cards that each player sets aside, face down, from the top of their own deck while setting up to play. When you Knock Out an opposing Pokémon, you take one of your Prize cards and put it into your hand. If you’re the first one to take your last Prize card, you win!",
			deck: "Each player starts with their own deck of 60 cards to play the game.  While both players know how many cards are in each deck, no one can look at or change the order of the cards in either player’s deck unless a card says so.",
			inPlay: {
				activeSpot:
					"The top row of a player’s in-play section is the Active Spot. Each player starts with (and must always have) one Pokémon in their Active Spot— this is the Active Pokémon. Each player may have only one Active Pokémon at a time. If your opponent doesn’t have any more Pokémon in play, you win the game!",
				bench:
					"The bottom row of a player’s in-play section is for the Benched Pokémon. Each player may have up to 5 Pokémon on the Bench at any one time. Any Pokémon in play other than the Active Pokémon must be put on the Bench",
			},
			discardPile:
				"The bottom row of a player’s in-play section is for the Benched Pokémon. Each player may have up to 5 Pokémon on the Bench at any one time. Any Pokémon in play other than the Active Pokémon must be put on the Bench",
		},
	},
};
export type RulesType = typeof Rules;
