import React, { useState } from 'react';
import { Card } from './Card';
import { BATTLE_PARTY_SIZE, GREAT_LEAGUE_CP_CAP } from '../rules/handbook';
import { validateBattleParty } from '../rules/validateAction';
import './PartySelect.css';

export const PartySelect = ({ me, gameNumber, winsNeeded, onConfirm, confirming, isPractice, onNewGame }) => {
  const [selected, setSelected] = useState([]);

  const toggle = (cardId) => {
    setSelected((prev) => {
      if (prev.includes(cardId)) return prev.filter((id) => id !== cardId);
      if (prev.length >= BATTLE_PARTY_SIZE) return prev;
      return [...prev, cardId];
    });
  };

  const check = validateBattleParty(me?.battleTeam || [], selected);

  return (
    <div className="party-select-screen">
      <div className="party-select-panel pixel-panel animate-slam-in">
        <div className="arcade-marquee">
          {isPractice ? 'PRACTICE · VS CPU' : 'TEAM PREVIEW · §6.1'}
        </div>
        <h1 className="game-title">SELECT BATTLE PARTY</h1>
        <p className="party-select-meta">
          Game {gameNumber || 1} · First to {winsNeeded || 2} · Great League ≤{GREAT_LEAGUE_CP_CAP} CP
        </p>
        <p className="party-select-hint">
          Choose exactly {BATTLE_PARTY_SIZE} Pokémon from your Battle Team. First pick leads as Active.
        </p>

        {me?.partyReady ? (
          <p className="party-waiting animate-insert-coin">PARTY LOCKED — WAITING FOR OPPONENT</p>
        ) : (
          <>
            <div className="party-team-grid">
              {(me?.battleTeam || []).map((card) => {
                const on = selected.includes(card.id);
                const order = selected.indexOf(card.id);
                return (
                  <button
                    key={card.id}
                    type="button"
                    className={`party-pick ${on ? 'picked' : ''}`}
                    onClick={() => toggle(card.id)}
                  >
                    {on && <span className="pick-order">{order === 0 ? 'LEAD' : order + 1}</span>}
                    <Card card={card} size="sm" />
                    <span className="cp-chip">CP {card.combatPower}</span>
                  </button>
                );
              })}
            </div>

            {!check.ok && selected.length > 0 && (
              <div className="inline-error">{check.error}</div>
            )}

            <button
              type="button"
              className="pixel-btn primary start-btn"
              disabled={!check.ok || confirming}
              onClick={() => onConfirm(selected)}
            >
              LOCK PARTY
            </button>
          </>
        )}

        {onNewGame && (
          <button type="button" className="pixel-btn" onClick={onNewGame}>
            NEW GAME
          </button>
        )}
      </div>
    </div>
  );
};
