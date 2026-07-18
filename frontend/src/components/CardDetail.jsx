import React, { useEffect } from 'react';
import {
    Zap, Droplets, Flame, Leaf, Eye, Shield, X,
    Heart, Swords, ShieldBan, Sparkles, Wind,
} from 'lucide-react';
import './CardDetail.css';

const TypeIcon = ({ type, size = 16, className = '' }) => {
    switch (type) {
        case 'Electric':
            return <Zap className={`type-icon electric ${className}`} size={size} />;
        case 'Fire':
            return <Flame className={`type-icon fire ${className}`} size={size} />;
        case 'Water':
            return <Droplets className={`type-icon water ${className}`} size={size} />;
        case 'Grass':
            return <Leaf className={`type-icon grass ${className}`} size={size} />;
        case 'Psychic':
            return <Eye className={`type-icon psychic ${className}`} size={size} />;
        case 'Trainer':
        case 'Power':
            return <Shield className={`type-icon trainer ${className}`} size={size} />;
        case 'Energy':
            return <Zap className={`type-icon energy ${className}`} size={size} />;
        default:
            return <div className={`type-icon neutral ${className}`} />;
    }
};

const STAT_KEYS = [
    { key: 'hp', label: 'HP', Icon: Heart, tone: 'hp' },
    { key: 'attack', label: 'Attack', Icon: Swords, tone: 'atk' },
    { key: 'defense', label: 'Defense', Icon: Shield, tone: 'def' },
    { key: 'spAttack', label: 'Sp. Atk', Icon: Sparkles, tone: 'spatk' },
    { key: 'spDefense', label: 'Sp. Def', Icon: ShieldBan, tone: 'spdef' },
    { key: 'speed', label: 'Speed', Icon: Wind, tone: 'spd' },
];

const typeTone = (type = '') => {
    const t = String(type).toLowerCase();
    if (t.includes('fire')) return 'fire';
    if (t.includes('water')) return 'water';
    if (t.includes('grass') || t.includes('bug')) return 'grass';
    if (t.includes('electric')) return 'electric';
    if (t.includes('psychic') || t.includes('fairy') || t.includes('ghost')) return 'psychic';
    if (t.includes('fighting') || t.includes('rock') || t.includes('ground')) return 'fighting';
    if (t.includes('ice') || t.includes('dragon') || t.includes('steel')) return 'steel';
    if (t.includes('poison') || t.includes('dark')) return 'poison';
    return 'neutral';
};

const statBarWidth = (value) => `${Math.min(100, Math.round((Number(value) / 180) * 100))}%`;

export const CardDetail = ({ card, ownerLabel, onClose }) => {
    useEffect(() => {
        const onKey = (e) => {
            if (e.key === 'Escape') onClose();
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [onClose]);

    if (!card) return null;

    const elementType = card.elementType || card.type;
    const tone = typeTone(elementType);
    const maxHp = card.maxHp || card.hp;
    const hpPct = maxHp ? Math.max(0, Math.min(100, Math.round((card.hp / maxHp) * 100))) : 100;

    return (
        <div className="card-detail-backdrop" onClick={onClose} role="presentation">
            <div
                className={`card-detail-modal pixel-panel animate-slam-in tone-${tone}`}
                onClick={(e) => e.stopPropagation()}
                role="dialog"
                aria-modal="true"
                aria-labelledby="card-detail-title"
            >
                <div className="card-detail-ribbon">
                    <span>POKÉMON CARD</span>
                    {ownerLabel && <span className="card-detail-owner">{ownerLabel}</span>}
                </div>

                <button type="button" className="card-detail-close pixel-btn" onClick={onClose} aria-label="Close">
                    <X size={14} />
                </button>

                <header className="tcg-header">
                    <div className="tcg-identity">
                        <h2 id="card-detail-title">{card.name}</h2>
                        <div className={`tcg-type-badge tone-${tone}`}>
                            <TypeIcon type={elementType} size={16} />
                            <span>{elementType}</span>
                        </div>
                    </div>
                    {card.hp != null && (
                        <div className="tcg-hp-badge" aria-label={`HP ${card.hp} of ${maxHp}`}>
                            <span className="tcg-hp-label">HP</span>
                            <span className="tcg-hp-value">{card.hp}</span>
                            {maxHp != null && <span className="tcg-hp-max">/{maxHp}</span>}
                        </div>
                    )}
                </header>

                <div className={`tcg-art-frame tone-${tone}`}>
                    <div className="tcg-art pixel-screen">
                        {card.imageUrl ? (
                            <img src={card.imageUrl} alt={card.name} draggable={false} />
                        ) : (
                            <TypeIcon type={elementType} size={72} className="card-art-icon" />
                        )}
                    </div>
                    {card.hp != null && (
                        <div className="tcg-hp-bar" aria-hidden="true">
                            <div
                                className={`hp-bar-fill ${hpPct <= 25 ? 'low' : hpPct <= 50 ? 'mid' : ''}`}
                                style={{ width: `${hpPct}%` }}
                            />
                        </div>
                    )}
                </div>

                <div className="tcg-meta-row">
                    <span className="meta-chip">Energy ×{card.energyAttached || 0}</span>
                    {card.combatPower > 0 && (
                        <span className="meta-chip">CP {card.combatPower}</span>
                    )}
                    {card.pokeApiId > 0 && (
                        <span className="meta-chip">#{String(card.pokeApiId).padStart(3, '0')}</span>
                    )}
                </div>

                {card.stats && (
                    <section className="card-detail-section">
                        <h3>BASE STATS</h3>
                        <div className="stat-tiles">
                            {STAT_KEYS.map(({ key, label, Icon, tone: statTone }) => {
                                const value = card.stats[key] ?? 0;
                                return (
                                    <div key={key} className={`stat-tile tone-${statTone}`}>
                                        <div className="stat-tile-top">
                                            <span className="stat-tile-icon" aria-hidden="true">
                                                <Icon size={14} />
                                            </span>
                                            <span className="stat-tile-label">{label}</span>
                                            <span className="stat-tile-value">{value}</span>
                                        </div>
                                        <div className="stat-bar-track">
                                            <div
                                                className="stat-bar-fill"
                                                style={{ width: statBarWidth(value) }}
                                            />
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </section>
                )}

                {card.attacks?.length > 0 && (
                    <section className="card-detail-section">
                        <h3>ATTACKS</h3>
                        <ul className="detail-attacks">
                            {card.attacks.map((att, i) => (
                                <li key={i} className="tcg-move">
                                    <div className="detail-attack-cost">
                                        {att.cost > 0
                                            ? Array(att.cost)
                                                .fill(0)
                                                .map((_, j) => (
                                                    <span key={j} className={`energy-orb tone-${tone}`}>
                                                        <TypeIcon type={elementType} size={12} />
                                                    </span>
                                                ))
                                            : <span className="cost-free">FREE</span>}
                                    </div>
                                    <div className="tcg-move-body">
                                        <span className="detail-attack-name">{att.name}</span>
                                        <span className="detail-attack-hint">
                                            Needs {att.cost} energy
                                        </span>
                                    </div>
                                    <span className="detail-attack-dmg">
                                        <span className="dmg-num">{att.damage}</span>
                                        <span className="dmg-unit">DMG</span>
                                    </span>
                                </li>
                            ))}
                        </ul>
                    </section>
                )}

                <button type="button" className="pixel-btn primary card-detail-dismiss" onClick={onClose}>
                    CLOSE
                </button>
            </div>
        </div>
    );
};
