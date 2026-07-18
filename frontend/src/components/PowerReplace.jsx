import React from 'react';
import { Card } from './Card';
import './PowerReplace.css';

/**
 * When all power slots are full, the new draw sits in pendingDraw.
 * Player picks a hand card to discard and replace, or keeps the current hand.
 */
export const PowerReplace = ({ drawn, hand, onReplace, onKeep, resolving }) => {
    if (!drawn) return null;

    return (
        <div
            className="power-replace-backdrop"
            role="dialog"
            aria-modal="true"
            aria-labelledby="power-replace-title"
        >
            <div className="power-replace-modal pixel-panel animate-slam-in">
                <div className="arcade-marquee power-replace-marquee">POWER SLOTS FULL</div>
                <h2 id="power-replace-title" className="power-replace-title">
                    SWITCH A POWER CARD?
                </h2>
                <p className="power-replace-hint">
                    You drew a new card but all 4 slots are filled. Tap a card in your hand to replace,
                    or keep your hand and discard the draw.
                </p>

                <div className="power-replace-drawn">
                    <span className="zone-label">NEW DRAW</span>
                    <Card card={drawn} size="md" className="power-replace-new" />
                </div>

                <div className="power-replace-hand">
                    <span className="zone-label">YOUR HAND — TAP TO REPLACE</span>
                    <div className="power-replace-row">
                        {(hand || []).map((card) => (
                            <button
                                key={card.id}
                                type="button"
                                className="power-replace-slot"
                                onClick={() => !resolving && onReplace(card.id)}
                                disabled={resolving}
                                aria-label={`Replace ${card.name}`}
                            >
                                <Card card={card} size="sm" isPlayable={!resolving} />
                                <span className="cp-chip power-chip">SWAP</span>
                            </button>
                        ))}
                    </div>
                </div>

                <button
                    type="button"
                    className="pixel-btn power-replace-keep"
                    onClick={() => !resolving && onKeep()}
                    disabled={resolving}
                >
                    KEEP HAND · DISCARD DRAW
                </button>
            </div>
        </div>
    );
};
