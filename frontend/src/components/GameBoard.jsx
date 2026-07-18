import React, { useState } from 'react';
import { useGameState } from '../hooks/useGameState';
import { Card } from './Card';
import './GameBoard.css';

export const GameBoard = () => {
    const {
        gameState, actionLog, loading, error, isMyTurn, me, opponent,
        actions, setPlayerId, playerId, needsPromote
    } = useGameState();
    const [selectedHandCard, setSelectedHandCard] = useState(null);
    const [selectedBenchedCard, setSelectedBenchedCard] = useState(null);

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

    if (!gameState || gameState.status === 'Waiting') {
        return (
            <div className="waiting-screen">
                <div className="arcade-cabinet pixel-panel">
                    <div className="arcade-marquee">1 PLAYER · COIN OP</div>
                    <h1 className="game-title">POKÉMON<br />CARD BATTLE</h1>
                    <p className="game-subtitle">8-BIT DUEL MODE</p>
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
                    <p className="insert-coin animate-insert-coin">INSERT COIN</p>
                </div>
            </div>
        );
    }

    if (gameState.status === 'GameOver') {
        return (
            <div className="waiting-screen">
                <div className="arcade-cabinet pixel-panel animate-slam-in">
                    <div className="arcade-marquee">ROUND OVER</div>
                    <h1 className="game-title">GAME OVER</h1>
                    <p className="winner-banner">{gameState.winner} WINS!</p>
                    <p className="last-action">{gameState.lastAction}</p>
                    <button className="pixel-btn primary start-btn" onClick={actions.startGame}>
                        CONTINUE?
                    </button>
                    <p className="insert-coin animate-insert-coin">PRESS START</p>
                </div>
            </div>
        );
    }

    const handleHandCardClick = (card) => {
        if (!isMyTurn) return;
        if (selectedHandCard?.id === card.id) {
            setSelectedHandCard(null);
        } else {
            setSelectedHandCard(card);
            setSelectedBenchedCard(null);
        }
    };

    const handleBenchClick = (card) => {
        if (!isMyTurn && !needsPromote) return;

        if (needsPromote && card) {
            actions.promote(card.id);
            setSelectedBenchedCard(null);
            return;
        }

        if (!isMyTurn) return;

        if (card) {
            setSelectedBenchedCard(selectedBenchedCard?.id === card.id ? null : card);
            setSelectedHandCard(null);
        } else if (selectedHandCard && selectedHandCard.type === 'Pokemon') {
            actions.playBench(selectedHandCard.id);
            setSelectedHandCard(null);
        }
    };

    const handleActiveClick = () => {
        if (!isMyTurn) return;
        if (!me?.activePokemon) {
            if (selectedBenchedCard) {
                actions.setActive(selectedBenchedCard.id);
                setSelectedBenchedCard(null);
            } else if (selectedHandCard && selectedHandCard.type === 'Pokemon') {
                actions.setActive(selectedHandCard.id);
                setSelectedHandCard(null);
            }
        } else if (selectedHandCard && selectedHandCard.type === 'Energy') {
            actions.attachEnergy(selectedHandCard.id);
            setSelectedHandCard(null);
        }
    };

    const renderBench = (player, isMe) => {
        const benchSlots = Array(5).fill(null);
        player?.benchedPokemon?.forEach((card, i) => {
            benchSlots[i] = card;
        });

        return (
            <div className={`bench-area ${isMe ? 'my-bench' : 'opponent-bench'}`}>
                <div className="bench-slots">
                    {benchSlots.map((card, idx) => (
                        <div key={idx} className="bench-slot-wrapper" onClick={() => isMe ? handleBenchClick(card) : null}>
                            {card ? (
                                <Card
                                    card={card}
                                    isPlayable={isMe && ((isMyTurn && !me?.activePokemon) || needsPromote)}
                                />
                            ) : (
                                <div className={`card empty-slot ${isMe && isMyTurn && selectedHandCard?.type === 'Pokemon' ? 'playable' : ''}`}>
                                    BENCH
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            </div>
        );
    };

    return (
        <div className="game-board">
            <div className="hud pixel-panel">
                <div className="player-info opponent-info">
                    <span className="player-name">VS {opponent?.id}</span>
                    <div>
                        <span className="stat-chip">DECK {opponent?.deck?.length}</span>
                        <span className="stat-chip">HAND {opponent?.hand?.length}</span>
                        <span className="stat-chip">PRIZE {opponent?.prizeCards?.length}</span>
                    </div>
                </div>
                <div className="turn-indicator">
                    <h2 className={isMyTurn ? 'glow' : ''}>
                        {needsPromote ? 'PROMOTE!' : isMyTurn ? 'YOUR TURN' : "CPU WAIT"}
                    </h2>
                    <p className="turn-meta">TURN {gameState.turnNumber} · {gameState.lastAction}</p>
                </div>
                <div className="player-info my-info">
                    <span className="player-name">YOU ({me?.id})</span>
                    <div>
                        <span className="stat-chip">DECK {me?.deck?.length}</span>
                        <span className="stat-chip">PRIZE {me?.prizeCards?.length}</span>
                        <span className="stat-chip">DISC {me?.discardPile?.length}</span>
                    </div>
                </div>
            </div>

            {error && <div className="banner-error">{error}</div>}

            <div className="arena">
                <div className="player-side top-side">
                    {renderBench(opponent, false)}
                    <div className="active-area">
                        {opponent?.activePokemon ? (
                            <Card card={opponent.activePokemon} isActive={true} />
                        ) : (
                            <div className="card empty-slot">ACTIVE</div>
                        )}
                    </div>
                </div>

                <div className="battle-line">
                    <span className="vs-badge">VS</span>
                </div>

                <div className="player-side bottom-side">
                    <div className="active-area" onClick={handleActiveClick}>
                        {me?.activePokemon ? (
                            <div className="active-card-container">
                                <Card card={me.activePokemon} isActive={true} />
                                {isMyTurn && !needsPromote && me.activePokemon.attacks?.map((att, i) => (
                                    <button
                                        key={i}
                                        className="pixel-btn danger attack-button"
                                        onClick={(e) => { e.stopPropagation(); actions.attack(i); }}
                                        disabled={!opponent?.activePokemon || (me.activePokemon.energyAttached || 0) < att.cost}
                                    >
                                        {att.name} · {att.damage} DMG
                                    </button>
                                ))}
                                {isMyTurn && selectedHandCard?.type === 'Energy' && (
                                    <button
                                        className="pixel-btn primary attach-btn"
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            actions.attachEnergy(selectedHandCard.id);
                                            setSelectedHandCard(null);
                                        }}
                                        disabled={me.hasAttached}
                                    >
                                        ATTACH ENERGY
                                    </button>
                                )}
                            </div>
                        ) : (
                            <div className={`card empty-slot ${(selectedBenchedCard || selectedHandCard) ? 'playable' : ''}`}>
                                SET ACTIVE
                            </div>
                        )}
                    </div>
                    {renderBench(me, true)}
                </div>

                <aside className="action-log pixel-screen">
                    <h3>■ LOG</h3>
                    <ul>
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

            <div className="hand-container pixel-panel">
                <div className="hand-actions">
                    <button
                        className="pixel-btn primary draw-btn"
                        onClick={actions.drawCard}
                        disabled={!isMyTurn || needsPromote || me?.hasDrawn || me?.deck?.length === 0}
                    >
                        DRAW
                    </button>
                    <div className="selected-info">
                        {selectedHandCard && (
                            <span>
                                SEL: {selectedHandCard.name}
                                {selectedHandCard.type === 'Energy' ? ' → ACTIVE' : ' → ZONE'}
                            </span>
                        )}
                        {selectedBenchedCard && <span>SEL BENCH: {selectedBenchedCard.name}</span>}
                        {needsPromote && <span>PICK BENCH TO PROMOTE</span>}
                        {!selectedHandCard && !selectedBenchedCard && !needsPromote && (
                            <span className="animate-insert-coin">SELECT A CARD</span>
                        )}
                    </div>
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
                        SWAP P1/P2
                    </button>
                </div>

                <div className="hand-cards">
                    {me?.hand?.map((card, i) => (
                        <div key={card.id || i} className={`hand-card-wrapper ${selectedHandCard?.id === card.id ? 'selected-card' : ''}`}>
                            <Card
                                card={card}
                                isPlayable={isMyTurn && !needsPromote}
                                onClick={() => handleHandCardClick(card)}
                            />
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};
