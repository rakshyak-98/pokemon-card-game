import React, { useEffect } from 'react';
import { Zap, Droplets, Flame, Leaf, Eye, Shield, X } from 'lucide-react';
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
            return <Shield className={`type-icon trainer ${className}`} size={size} />;
        case 'Energy':
            return <Zap className={`type-icon energy ${className}`} size={size} />;
        default:
            return <div className={`type-icon neutral ${className}`} />;
    }
};

const STAT_KEYS = [
    { key: 'hp', label: 'HP' },
    { key: 'attack', label: 'ATK' },
    { key: 'defense', label: 'DEF' },
    { key: 'spAttack', label: 'SP.ATK' },
    { key: 'spDefense', label: 'SP.DEF' },
    { key: 'speed', label: 'SPD' },
];

const statBarWidth = (value) => `${Math.min(100, Math.round((value / 200) * 100))}%`;

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
    const maxHp = card.maxHp || card.hp;
    const hpPct = maxHp ? Math.max(0, Math.min(100, Math.round((card.hp / maxHp) * 100))) : 100;

    return (
        <div className="card-detail-backdrop" onClick={onClose} role="presentation">
            <div
                className="card-detail-modal pixel-panel animate-slam-in"
                onClick={(e) => e.stopPropagation()}
                role="dialog"
                aria-modal="true"
                aria-labelledby="card-detail-title"
            >
                <div className="card-detail-marquee">POKÉDEX ENTRY</div>

                <button type="button" className="card-detail-close pixel-btn" onClick={onClose} aria-label="Close">
                    <X size={14} />
                </button>

                <div className="card-detail-layout">
                    <div className="card-detail-art pixel-screen">
                        {card.imageUrl ? (
                            <img src={card.imageUrl} alt={card.name} draggable={false} />
                        ) : (
                            <TypeIcon type={elementType} size={64} className="card-art-icon" />
                        )}
                    </div>

                    <div className="card-detail-info">
                        {ownerLabel && <p className="card-detail-owner">{ownerLabel}</p>}
                        <div className="card-detail-header">
                            <h2 id="card-detail-title">{card.name}</h2>
                            <div className="card-detail-type">
                                <TypeIcon type={elementType} size={18} />
                                <span>{elementType}</span>
                            </div>
                        </div>

                        {card.hp != null && (
                            <div className="card-detail-hp-block">
                                <div className="card-detail-hp-row">
                                    <span className="hp-label">HP</span>
                                    <span className="hp-value">
                                        {card.hp}
                                        {maxHp ? ` / ${maxHp}` : ''}
                                    </span>
                                </div>
                                <div className="hp-bar-track">
                                    <div
                                        className={`hp-bar-fill ${hpPct <= 25 ? 'low' : hpPct <= 50 ? 'mid' : ''}`}
                                        style={{ width: `${hpPct}%` }}
                                    />
                                </div>
                            </div>
                        )}

                        <div className="card-detail-meta">
                            <span className="meta-chip">ENERGY ×{card.energyAttached || 0}</span>
                            {card.pokeApiId > 0 && (
                                <span className="meta-chip">#{String(card.pokeApiId).padStart(3, '0')}</span>
                            )}
                        </div>
                    </div>
                </div>

                {card.stats && (
                    <section className="card-detail-section">
                        <h3>■ BASE STATS</h3>
                        <div className="stat-grid">
                            {STAT_KEYS.map(({ key, label }) => (
                                <div key={key} className="stat-row">
                                    <span className="stat-label">{label}</span>
                                    <div className="stat-bar-track">
                                        <div
                                            className="stat-bar-fill"
                                            style={{ width: statBarWidth(card.stats[key] || 0) }}
                                        />
                                    </div>
                                    <span className="stat-value">{card.stats[key] ?? 0}</span>
                                </div>
                            ))}
                        </div>
                    </section>
                )}

                {card.attacks?.length > 0 && (
                    <section className="card-detail-section">
                        <h3>■ ATTACKS</h3>
                        <ul className="detail-attacks">
                            {card.attacks.map((att, i) => (
                                <li key={i}>
                                    <div className="detail-attack-cost">
                                        {Array(att.cost)
                                            .fill(0)
                                            .map((_, j) => (
                                                <TypeIcon key={j} type="Energy" size={14} />
                                            ))}
                                        {att.cost === 0 && <span className="cost-free">FREE</span>}
                                    </div>
                                    <span className="detail-attack-name">{att.name}</span>
                                    <span className="detail-attack-dmg">{att.damage} DMG</span>
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
