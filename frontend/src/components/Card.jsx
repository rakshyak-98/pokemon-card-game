import React from 'react';
import { Zap, Droplets, Flame, Leaf, Eye, Shield, Swords, Heart } from 'lucide-react';
import './Card.css';

const powerMeta = (effect, value) => {
    switch (effect) {
        case 'boost_attack':
            return { label: 'ATTACK', detail: `+${value}`, iconType: 'Attack' };
        case 'boost_defense':
            return { label: 'DEFENSE', detail: `+${value}`, iconType: 'Defense' };
        case 'heal':
            return { label: 'HEAL', detail: `+${value}`, iconType: 'Heal' };
        default:
            return { label: 'POWER', detail: value ? `+${value}` : '', iconType: 'Power' };
    }
};

const TypeIcon = ({ type, className, size = 16 }) => {
    switch (type) {
        case 'Pokemon':
        case 'Electric': return <Zap className={`type-icon electric ${className}`} size={size} />;
        case 'Fire': return <Flame className={`type-icon fire ${className}`} size={size} />;
        case 'Water': return <Droplets className={`type-icon water ${className}`} size={size} />;
        case 'Grass': return <Leaf className={`type-icon grass ${className}`} size={size} />;
        case 'Psychic': return <Eye className={`type-icon psychic ${className}`} size={size} />;
        case 'Trainer': return <Shield className={`type-icon trainer ${className}`} size={size} />;
        case 'Power': return <Shield className={`type-icon trainer ${className}`} size={size} />;
        case 'Attack': return <Swords className={`type-icon attack ${className}`} size={size} />;
        case 'Defense': return <Shield className={`type-icon defense ${className}`} size={size} />;
        case 'Heal': return <Heart className={`type-icon heal ${className}`} size={size} />;
        case 'Energy': return <Zap className={`type-icon energy ${className}`} size={size} />;
        default: return <div className={`type-icon neutral ${className}`} />;
    }
};

export const Card = ({ card, className = '', onClick, isPlayable = false, isActive = false, size = 'md' }) => {
    if (!card) return <div className={`card empty-slot size-${size}`}>EMPTY</div>;

    const isPower = card.type === 'Power';
    const power = isPower ? powerMeta(card.effect, card.effectValue) : null;
    const elementType = card.elementType || (card.name?.includes('Electric') ? 'Electric' : card.type);

    return (
        <div
            className={`card size-${size} ${card.type?.toLowerCase() || ''} ${isPower ? `power-${card.effect || 'generic'}` : ''} ${isPlayable ? 'playable' : ''} ${isActive ? 'active-card' : ''} ${className}`}
            onClick={() => onClick && onClick(card)}
            role={onClick ? 'button' : undefined}
            tabIndex={onClick ? 0 : undefined}
            onKeyDown={onClick ? (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    onClick(card);
                }
            } : undefined}
        >
            <div className="card-header">
                <span className="card-name">{card.name}</span>
                {isPower ? (
                    <div className="card-hp power-badge">
                        <span className="hp-value">{power.detail}</span>
                        <TypeIcon type={power.iconType} size={14} />
                    </div>
                ) : card.hp != null ? (
                    <div className="card-hp">
                        <span className="hp-label">HP</span>
                        <span className="hp-value">
                            {card.hp}
                            {card.maxHp != null ? `/${card.maxHp}` : ''}
                        </span>
                        <TypeIcon type={elementType} />
                    </div>
                ) : null}
            </div>

            <div className={`card-art-box ${isPower ? 'power-art' : ''}`}>
                {isPower ? (
                    <TypeIcon type={power.iconType} className="card-art-icon power-art-icon" size={36} />
                ) : card.imageUrl ? (
                    <img
                        className="card-art-image"
                        src={card.imageUrl}
                        alt={card.name}
                        loading="lazy"
                        draggable={false}
                    />
                ) : (
                    <TypeIcon type={elementType} className="card-art-icon" />
                )}
            </div>

            <div className="card-body">
                <div className="card-type-label">
                    {isPower ? power.label : (card.elementType || card.type)}
                    {!isPower && card.stats && (
                        <span className="card-stats-inline">
                            {' '}ATK {card.stats.attack} / DEF {card.stats.defense}
                        </span>
                    )}
                </div>
                {isPower && (
                    <div className="power-effect-label">
                        USE ON ACTIVE
                    </div>
                )}
                {card.energyAttached > 0 && (
                    <div className="energy-attached">⚡ ×{card.energyAttached}</div>
                )}
                {card.attacks && card.attacks.length > 0 && (
                    <div className="attacks-list">
                        {card.attacks.map((attack, idx) => (
                            <div key={idx} className="attack-row">
                                <div className="attack-cost">
                                    {Array(attack.cost).fill(0).map((_, i) => (
                                        <TypeIcon key={i} type="Energy" size={12} />
                                    ))}
                                </div>
                                <span className="attack-name">{attack.name}</span>
                                <span className="attack-damage">{attack.damage}</span>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};
