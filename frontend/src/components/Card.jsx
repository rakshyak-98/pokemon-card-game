import React from 'react';
import { Zap, Droplets, Flame, Leaf, Eye, Shield } from 'lucide-react';
import './Card.css';

const TypeIcon = ({ type, className }) => {
    switch (type) {
        case 'Pokemon': // Generic for now, we'll map Electric/Fire/Water later if backend adds them
        case 'Electric': return <Zap className={`type-icon electric ${className}`} size={16} />;
        case 'Fire': return <Flame className={`type-icon fire ${className}`} size={16} />;
        case 'Water': return <Droplets className={`type-icon water ${className}`} size={16} />;
        case 'Grass': return <Leaf className={`type-icon grass ${className}`} size={16} />;
        case 'Psychic': return <Eye className={`type-icon psychic ${className}`} size={16} />;
        case 'Trainer': return <Shield className={`type-icon trainer ${className}`} size={16} />;
        case 'Energy': return <Zap className={`type-icon energy ${className}`} size={16} />;
        default: return <div className={`type-icon neutral ${className}`} />;
    }
};

export const Card = ({ card, className = '', onClick, isPlayable = false, isActive = false, size = 'md' }) => {
    if (!card) return <div className={`card empty-slot size-${size}`}>EMPTY</div>;

    const elementType = card.elementType || (card.name?.includes('Electric') ? 'Electric' : card.type);

    return (
        <div
            className={`card size-${size} ${card.type?.toLowerCase() || ''} ${isPlayable ? 'playable' : ''} ${isActive ? 'active-card' : ''} ${className}`}
            onClick={() => onClick && onClick(card)}
        >
            <div className="card-header">
                <span className="card-name">{card.name}</span>
                {card.hp != null && (
                    <div className="card-hp">
                        <span className="hp-label">HP</span>
                        <span className="hp-value">
                            {card.hp}
                            {card.maxHp != null ? `/${card.maxHp}` : ''}
                        </span>
                        <TypeIcon type={elementType} />
                    </div>
                )}
            </div>

            <div className="card-art-box">
                {card.imageUrl ? (
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
                    {card.elementType || card.type}
                    {card.stats && (
                        <span className="card-stats-inline">
                            {' '}ATK {card.stats.attack} / DEF {card.stats.defense}
                        </span>
                    )}
                </div>
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
