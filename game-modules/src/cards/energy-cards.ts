import DESCRIPTION from "./card-description";

type DescriptionType = {
	grass: string;
	fire: string;
	water: string;
	lightning: string;
	psychic: string;
	fighting: string;
	darkness: string;
	metal: string;
	fairy: string;
	dragon: string;
	colorless: string;
};

interface CardType {
	description: string;
}

class GrassCard implements CardType {
	public description = DESCRIPTION.grass;
}

class FireCard implements CardType {
	public description = DESCRIPTION.fire;
}
class WaterCard implements CardType {
	public description = DESCRIPTION.water;
}
class LightningCard implements CardType {
	public description = DESCRIPTION.lightning;
}
class PsychicCard implements CardType {
	public description = DESCRIPTION.psychic;
}
class FightingCard implements CardType {
	public description = DESCRIPTION.fighting;
}
class DarknessCard implements CardType {
	public description = DESCRIPTION.darkness;
}
class MetalCard implements CardType {
	public description = DESCRIPTION.metal;
}
class FairyCard implements CardType {
	public description = DESCRIPTION.fairy;
}
class DragonCard implements CardType {
	public description = DESCRIPTION.dragon;
}
class ColorlessCard implements CardType {
	public description = DESCRIPTION.colorless;
}
