import React, { useState } from 'react';
import { useGameState } from '../hooks/useGameState';
import { Card } from './Card';
import { CardDetail } from './CardDetail';
import { PartySelect } from './PartySelect';
import './GameBoard.css';

export const GameBoard = ({ onShowRules }) => {
    const {
        gameState, actionLog, loading, error, isMyTurn, me, opponent,
        actions, setPlayerId, playerId, needsPromote, needsPartySelect
    } = useGameState();
    const [selectedBenchedCard, setSelectedBenchedCard] = useState(null);
    const [detailView, setDetailView] = useState(null);
    const [confirmingParty, setConfirmingParty] = useState(false);

    const openDetails = (card, ownerLabel) => {
        if (!card) return;
        setDetailView({ card, ownerLabel });
    };

    if (loading) {
        return (
            <div className="loading-screen">
                <div className="arcade-cabinet pixel-panel">
                    <div className="arcade-marquee">SYSTEM BOOT</div>
                    <h1 className="game-title">LOADING...</h1>
                    <p className="insert-coin animate-insert-coin">PLEASE WAIT</p>
                </div>
            </div>
        );
    }

    if (!gameState || gameState.status === 'Waiting' || !gameState.phase || gameState.phase === 'Waiting') {
        return (
            <div className="waiting-screen">
                <div className="arcade-cabinet pixel-panel">
                    <div className="arcade-marquee">GO TOURNAMENT · GREAT LEAGUE</div>
                    <h1 className="game-title">POKÉMON GO<br />BATTLE</h1>
                    <p className="game-subtitle">HANDBOOK MATCH MODE</p>
                    <div className="player-select">
                        <label>PLAYER SELECT</label>
                        <select value={playerId} onChange={(e) => setPlayerId(e.target.value)}>
                            <option value="player1">PLAYER 1</option>
                            <option value="player2">PLAYER 2</option>
                        </select>
                    </div>
                    {error && <div className="inline-error">{error}</div>}
                    <button className="pixel-btn primary start-btn" onClick={actions.startGame}>
                        PRESS START
                    </button>
                    {onShowRules && (
                        <button type="button" className="pixel-btn rules-link-btn" onClick={onShowRules}>
                            How to play
                        </button>
                    )}
                    <p className="insert-coin animate-insert-coin">INSERT COIN</p>
                </div>
            </div>
        );
    }

    if (gameState.status === 'GameOver' || gameState.phase === 'MatchOver') {
        return (
            <div className="waiting-screen">
                <div className="arcade-cabinet pixel-panel animate-slam-in">
                    <div className="arcade-marquee">MATCH OVER · §6.4</div>
                    <h1 className="game-title">GAME OVER</h1>
                    <p className="winner-banner">{gameState.winner} WINS!</p>
                    <p className="last-action">{gameState.lastAction}</p>
                    <button className="pixel-btn primary start-btn" onClick={actions.startGame}>
                        CONTINUE?
                    </button>
                    {onShowRules && (
                        <button type="button" className="pixel-btn rules-link-btn" onClick={onShowRules}>
                            How to play
                        </button>
                    )}
                    <p className="insert-coin animate-insert-coin">PRESS START</p>
                </div>
            </div>
        );
    }

    if (needsPartySelect || ['TeamPreview', 'PartySelect', 'BetweenGames'].includes(gameState.phase)) {
        if (!me?.partyReady || !opponent?.partyReady) {
            return (
                <PartySelect
                    me={me}
                    opponent={opponent}
                    gameNumber={gameState.gameNumber}
                    winsNeeded={gameState.winsNeeded}
                    confirming={confirmingParty}
                    onConfirm={async (cardIds) => {
                        setConfirmingParty(true);
                        try {
                            await actions.selectParty(cardIds);
                        } finally {
                            setConfirmingParty(false);
                        }
                    }}
                />
            );
        }
    }

    const handleBenchClick = (card) => {
        if (!isMyTurn && !needsPromote) return;
        if (needsPromote && card) {
            actions.promote(card.id);
            setSelectedBenchedCard(null);
            return;
        }
        if (!isMyTurn || !card) return;
        setSelectedBenchedCard(selectedBenchedCard?.id === card.id ? null : card);
    };

    const renderBench = (player, isMe) => {
        const benchSlots = Array(2).fill(null);
        player?.benchedPokemon?.forEach((card, i) => {
            if (i < 2) benchSlots[i] = card;
        });

        return (
            <div className={`bench-area ${isMe ? 'my-bench' : 'opponent-bench'}`}>
                <span className="zone-label">{isMe ? 'YOUR BACK LINE' : 'OPP BACK LINE'}</span>
                <div className="bench-slots">
                    {benchSlots.map((card, idx) => (
                        <div
                            key={idx}
                            className="bench-slot-wrapper"
                            onClick={() => (isMe ? handleBenchClick(card) : null)}
                        >
                            {card ? (
                                <Card
                                    card={isMe ? card : { ...card, energyAttached: undefined, hp: undefined, maxHp: undefined, stats: undefined }}
                                    size="sm"
                                    isPlayable={isMe && needsPromote}
                                />
                            ) : (
                                <div className="card empty-slot size-sm">BACK</div>
                            )}
                        </div>
                    ))}
                </div>
            </div>
        );
    };

    const oppActiveDisplay = opponent?.activePokemon
        ? {
            ...opponent.activePokemon,
            // Private info (§6.5.3): hide energy from opponent view
            energyAttached: undefined,
        }
        : null;

    return (
        <div className="game-board">
            <header className="hud pixel-panel">
                <div className="player-info opponent-info">
                    <span className="player-name">VS {opponent?.id}</span>
                    <div className="stat-row">
                        <span className="stat-chip">WINS {opponent?.gamesWon || 0}</span>
                        <span className="stat-chip">LEFT {(opponent?.activePokemon ? 1 : 0) + (opponent?.benchedPokemon?.length || 0)}</span>
                        <span className="stat-chip">SHIELD {opponent?.protectShields ?? 0}</span>
                    </div>
                </div>
                <div className="turn-indicator">
                    <h2 className={isMyTurn ? 'glow' : ''}>
                        {needsPromote ? 'PROMOTE!' : isMyTurn ? 'YOUR TURN' : 'WAIT'}
                    </h2>
                    <p className="turn-meta">
                        G{gameState.gameNumber} · T{gameState.turnNumber}
                        {gameState.lastAction ? ` · ${gameState.lastAction}` : ''}
                    </p>
                </div>
                <div className="player-info my-info">
                    <span className="player-name">YOU ({me?.id})</span>
                    <div className="stat-row">
                        <span className="stat-chip">WINS {me?.gamesWon || 0}</span>
                        <span className="stat-chip">SHIELD {me?.protectShields ?? 0}</span>
                        <span className="stat-chip">CP TEAM</span>
                    </div>
                </div>
            </header>

            {error && <div className="banner-error">{error}</div>}

            <div className="battlefield">
                <section className="field-column">
                    <div className="player-side top-side">
                        {renderBench(opponent, false)}
                        <div className="active-zone">
                            <span className="zone-label">OPP ACTIVE</span>
                            <div className="active-area">
                                {oppActiveDisplay ? (
                                    <Card card={oppActiveDisplay} size="lg" isActive={true} />
                                ) : (
                                    <div className="card empty-slot size-lg">ACTIVE</div>
                                )}
                            </div>
                        </div>
                    </div>

                    <div className="battle-line" aria-hidden="true">
                        <span className="vs-badge">VS</span>
                    </div>

                    <div className="player-side bottom-side">
                        <div className="active-zone">
                            <span className="zone-label">YOUR ACTIVE</span>
                            <div className="active-row">
                                <div className="active-area">
                                    {me?.activePokemon ? (
                                        <Card card={me.activePokemon} size="lg" isActive={true} />
                                    ) : (
                                        <div className="card empty-slot size-lg">PROMOTE</div>
                                    )}
                                </div>
                                {me?.activePokemon && (
                                    <div className="combat-panel pixel-panel">
                                        <span className="combat-label">
                                            {isMyTurn && !needsPromote ? 'ACTIONS' : 'INFO'}
                                        </span>
                                        <button
                                            type="button"
                                            className="pixel-btn details-btn"
                                            onClick={() => openDetails(me.activePokemon, `YOU (${me.id})`)}
                                        >
                                            DETAILS
                                        </button>
                                        {isMyTurn && !needsPromote && (
                                            <button
                                                type="button"
                                                className="pixel-btn primary attach-btn"
                                                onClick={() => actions.attachEnergy('')}
                                                disabled={me.hasAttached}
                                            >
                                                CHARGE ENERGY
                                            </button>
                                        )}
                                        {isMyTurn &&
                                            !needsPromote &&
                                            me.activePokemon.attacks?.map((att, i) => (
                                                <button
                                                    key={i}
                                                    className="pixel-btn danger attack-button"
                                                    onClick={() => actions.attack(i)}
                                                    disabled={
                                                        !opponent?.activePokemon ||
                                                        (me.activePokemon.energyAttached || 0) < att.cost
                                                    }
                                                >
                                                    {att.name} · {att.damage}
                                                </button>
                                            ))}
                                    </div>
                                )}
                            </div>
                        </div>
                        {renderBench(me, true)}
                    </div>
                </section>

                <aside className="action-log pixel-screen">
                    <h3>■ LOG</h3>
                    <ul>
                        {actionLog.length === 0 && (
                            <li className="ok">
                                <span className="log-type">READY</span>
                                <span className="log-player">Waiting for actions…</span>
                            </li>
                        )}
                        {actionLog.map((a) => (
                            <li key={a.id} className={a.success ? 'ok' : 'fail'}>
                                <span className="log-type">{a.actionType}</span>
                                <span className="log-player">{a.playerId}</span>
                                {!a.success && <span className="log-err">{a.errorMessage}</span>}
                            </li>
                        ))}
                    </ul>
                </aside>
            </div>

            <footer className="hand-container pixel-panel">
                <div className="hand-actions">
                    <div className="selected-info">
                        {needsPromote && <span>PICK BACK-LINE TO PROMOTE (§6.3)</span>}
                        {!needsPromote && isMyTurn && (
                            <span className="animate-insert-coin">CHARGE · ATTACK · OR END TURN</span>
                        )}
                        {!needsPromote && !isMyTurn && <span>WAITING FOR OPPONENT</span>}
                    </div>
                    <div className="hand-action-btns">
                        <button
                            className="pixel-btn end-turn-btn"
                            disabled={!isMyTurn || needsPromote}
                            onClick={actions.endTurn}
                        >
                            END TURN
                        </button>
                        <button
                            className="pixel-btn"
                            onClick={() => setPlayerId(playerId === 'player1' ? 'player2' : 'player1')}
                        >
                            SWAP
                        </button>
                        {onShowRules && (
                            <button type="button" className="pixel-btn" onClick={onShowRules}>
                                RULES
                            </button>
                        )}
                    </div>
                </div>
                <div className="hand-cards party-strip">
                    {(me?.battleTeam || []).map((card) => (
                        <div key={card.id} className="hand-card-wrapper">
                            <Card card={card} size="sm" />
                            <span className="cp-chip">CP {card.combatPower}</span>
                        </div>
                    ))}
                </div>
            </footer>

            {detailView && (
                <CardDetail
                    card={detailView.card}
                    ownerLabel={detailView.ownerLabel}
                    onClose={() => setDetailView(null)}
                />
            )}
        </div>
    );
};
