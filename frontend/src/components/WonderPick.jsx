import React, { useEffect, useState } from 'react';
import { Card } from './Card';
import './WonderPick.css';

function shuffle(list) {
    const next = [...list];
    for (let i = next.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [next[i], next[j]] = [next[j], next[i]];
    }
    return next;
}

/**
 * Pokémon GO / TCG Pocket–style Wonder Pick:
 * preview face-up → shuffle face-down → pick one card into hand.
 */
export const WonderPick = ({ cards, onSelect, selecting }) => {
    const [phase, setPhase] = useState('preview'); // preview | shuffle | pick | reveal
    const [slots, setSlots] = useState(cards);
    const [pickedId, setPickedId] = useState(null);
    const pickKey = cards.map((c) => c.id).join('|');

    useEffect(() => {
        setSlots(cards);
        setPhase('preview');
        setPickedId(null);

        const toShuffle = setTimeout(() => setPhase('shuffle'), 1800);
        const toPick = setTimeout(() => {
            setSlots((prev) => shuffle(prev));
            setPhase('pick');
        }, 2800);

        return () => {
            clearTimeout(toShuffle);
            clearTimeout(toPick);
        };
        // Restart only when the pending card set changes (not on poll refreshes).
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [pickKey]);

    const handlePick = (card) => {
        if (phase !== 'pick' || selecting || pickedId) return;
        setPickedId(card.id);
        setPhase('reveal');
        setTimeout(() => onSelect(card.id), 700);
    };

    const faceUp = phase === 'preview' || phase === 'reveal';

    return (
        <div className="wonder-pick-backdrop" role="dialog" aria-modal="true" aria-labelledby="wonder-pick-title">
            <div className="wonder-pick-modal pixel-panel animate-slam-in">
                <div className="wonder-pick-marquee">WONDER PICK</div>
                <h2 id="wonder-pick-title" className="wonder-pick-title">
                    {phase === 'preview' && 'MEMORIZE THE CARDS'}
                    {phase === 'shuffle' && 'SHUFFLING…'}
                    {phase === 'pick' && 'PICK ONE CARD'}
                    {phase === 'reveal' && 'GOT IT!'}
                </h2>
                <p className="wonder-pick-hint">
                    {phase === 'preview' && 'Cards flip face-down after a moment — choose carefully.'}
                    {phase === 'shuffle' && 'Equal chance for each card.'}
                    {phase === 'pick' && 'Tap a card to add it to your hand. The rest return to your deck.'}
                    {phase === 'reveal' && 'Adding to hand…'}
                </p>

                <div className={`wonder-pick-row ${phase === 'shuffle' ? 'is-shuffling' : ''}`}>
                    {slots.map((card, i) => {
                        const showFace = faceUp || (phase === 'reveal' && card.id === pickedId);
                        const isPicked = card.id === pickedId;
                        return (
                            <button
                                key={`${card.id}-${i}`}
                                type="button"
                                className={`wonder-slot ${showFace ? 'face-up' : 'face-down'} ${
                                    isPicked ? 'picked' : ''
                                } ${phase === 'pick' ? 'selectable' : ''}`}
                                onClick={() => handlePick(card)}
                                disabled={phase !== 'pick' || selecting}
                                aria-label={showFace ? card.name : `Face-down card ${i + 1}`}
                            >
                                {showFace ? (
                                    <Card card={card} size="md" />
                                ) : (
                                    <div className="card-back pixel-screen">
                                        <span className="card-back-mark">?</span>
                                    </div>
                                )}
                            </button>
                        );
                    })}
                </div>
            </div>
        </div>
    );
};
