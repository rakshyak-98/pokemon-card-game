"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const card_description_1 = __importDefault(require("./card-description"));
class GrassCard {
    constructor() {
        this.description = card_description_1.default.grass;
    }
}
class FireCard {
    constructor() {
        this.description = card_description_1.default.fire;
    }
}
class WaterCard {
    constructor() {
        this.description = card_description_1.default.water;
    }
}
class LightningCard {
    constructor() {
        this.description = card_description_1.default.lightning;
    }
}
class PsychicCard {
    constructor() {
        this.description = card_description_1.default.psychic;
    }
}
class FightingCard {
    constructor() {
        this.description = card_description_1.default.fighting;
    }
}
class DarknessCard {
    constructor() {
        this.description = card_description_1.default.darkness;
    }
}
class MetalCard {
    constructor() {
        this.description = card_description_1.default.metal;
    }
}
class FairyCard {
    constructor() {
        this.description = card_description_1.default.fairy;
    }
}
class DragonCard {
    constructor() {
        this.description = card_description_1.default.dragon;
    }
}
class ColorlessCard {
    constructor() {
        this.description = card_description_1.default.colorless;
    }
}
