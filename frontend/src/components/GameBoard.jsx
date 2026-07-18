import React, { useState } from 'react';
import { useGameState } from '../hooks/useGameState';
import { Card } from './Card';
import { CardDetail } from './CardDetail';
import { PartySelect } from './PartySelect';
import { MAX_POWER_HAND_SLOTS } from '../rules/handbook';
import './GameBoard.css';

export const GameBoard = ({ onShowRules }) => {
    const {
        gameState, actionLog, loading, error, isMyTurn, me, opponent,
        actions, setPlayerId, playerId, needsPromote, needsPartySelect,
        vsCPU, setVsCPU, isPractice, cpuThinking
    } = useGameState();
    const [selectedBenchedCard, setSelectedBenchedCard] = useState(null);
    const [detailView, setDetailView] = useState(null);
    const [confirmingParty, setConfirmingParty] = useState(false);
    const [menuOpen, setMenuOpen] = useState(false);

    const openDetails = (card, ownerLabel) => {
        if (!card) return;
        setDetailView({ card, ownerLabel });
    };

    const opponentLabel = isPractice ? 'CPU' : opponent?.id;

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

                    <label className="mode-toggle">
                        <input
                            type="checkbox"
                            checked={vsCPU}
                            onChange={(e) => setVsCPU(e.target.checked)}
                        />
                        <span>Practice vs CPU (learn the rules)</span>
                    </label>

                    {!vsCPU && (
                        <div className="player-select">
                            <label>PLAYER SELECT</label>
                            <select value={playerId} onChange={(e) => setPlayerId(e.target.value)}>
                                <option value="player1">PLAYER 1</option>
                                <option value="player2">PLAYER 2</option>
                            </select>
                        </div>
                    )}
                    {error && <div className="inline-error">{error}</div>}
                    <div className="cabinet-actions">
                        <button
                            className="pixel-btn primary start-btn"
                            onClick={() => actions.startGame({ vsCPU })}
                        >
                            {vsCPU ? 'PRACTICE START' : 'PRESS START'}
                        </button>
                        {onShowRules && (
                            <button type="button" className="pixel-btn rules-link-btn" onClick={onShowRules}>
                                How to play
                            </button>
                        )}
                    </div>
                    <p className="insert-coin animate-insert-coin">INSERT COIN</p>
                </div>
            </div>
        );
    }

    if (gameState.status === 'GameOver' || gameState.phase === 'MatchOver') {
        return (
            <div className="waiting-screen">
                <div className="arcade-cabinet pixel-panel animate-slam-in">
                    <div className="arcade-marquee">
                        {isPractice ? 'PRACTICE OVER' : 'MATCH OVER · §6.4'}
                    </div>
                    <h1 className="game-title">GAME OVER</h1>
                    <p className="winner-banner">
                        {gameState.winner === gameState.cpuPlayerId
                            ? 'CPU WINS!'
                            : gameState.winner === playerId
                              ? 'YOU WIN!'
                              : `${gameState.winner} WINS!`}
                    </p>
                    <p className="last-action">{gameState.lastAction}</p>
                    <div className="cabinet-actions">
                        <button
                            className="pixel-btn primary start-btn"
                            onClick={() => actions.startGame({ vsCPU: isPractice })}
                        >
                            {isPractice ? 'PRACTICE AGAIN' : 'CONTINUE?'}
                        </button>
                        {onShowRules && (
                            <button type="button" className="pixel-btn rules-link-btn" onClick={onShowRules}>
                                How to play
                            </button>
                        )}
                    </div>
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
                    gameNumber={gameState.gameNumber}
                    winsNeeded={gameState.winsNeeded}
                    confirming={confirmingParty}
                    isPractice={isPractice}
                    onNewGame={() => actions.startGame({ vsCPU: isPractice })}
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
        if (!isMyTurn || !card || cpuThinking) return;
        if (me?.activePokemon && !me.hasSwitched) {
            actions.switchActive(card.id);
            setSelectedBenchedCard(null);
            return;
        }
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
                                    isPlayable={isMe && (needsPromote || (isMyTurn && !cpuThinking && !!me?.activePokemon && !me?.hasSwitched))}
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

    const renderActiveHp = (card) => {
        if (!card || card.hp == null) return null;
        const maxHp = card.maxHp || card.hp;
        const pct = maxHp ? Math.max(0, Math.min(100, Math.round((card.hp / maxHp) * 100))) : 100;
        return (
            <div className="active-hp-meter" aria-label={`HP ${card.hp} of ${maxHp}`}>
                <div className="active-hp-row">
                    <span className="hp-label">HP</span>
                    <span className="hp-value">{card.hp} / {maxHp}</span>
                </div>
                <div className="active-hp-track">
                    <div
                        className={`active-hp-fill ${pct <= 25 ? 'low' : pct <= 50 ? 'mid' : ''}`}
                        style={{ width: `${pct}%` }}
                    />
                </div>
            </div>
        );
    };

    const runMenuAction = (fn) => {
        setMenuOpen(false);
        fn?.();
    };

    const statusHint = needsPromote
        ? 'PICK BACK-LINE TO PROMOTE (§6.3)'
        : isMyTurn
          ? 'DRAW INTO EMPTY SLOT · USE OR KEEP · CHARGE · ATTACK'
          : cpuThinking
            ? 'CPU IS PLAYING…'
            : 'WAITING FOR OPPONENT';

    const canDrawPower =
        isMyTurn &&
        !needsPromote &&
        !cpuThinking &&
        !!me?.activePokemon &&
        !me?.hasDrawn &&
        (me?.hand?.length ?? 0) < MAX_POWER_HAND_SLOTS &&
        (me?.powerDeck?.length ?? 0) > 0;

    const canPlayPower = isMyTurn && !needsPromote && !cpuThinking && !me?.hasPlayedPower && !!me?.activePokemon;

    const handlePowerClick = (card) => {
        if (!canPlayPower || !card) return;
        actions.playPower(card.id);
    };

    const renderPowerStrip = () => {
        const hand = me?.hand || [];
        const deckCount = me?.powerDeck?.length ?? 0;
        const slots = Array.from({ length: MAX_POWER_HAND_SLOTS }, (_, i) => hand[i] || null);
        const nextEmptyIndex = slots.findIndex((c) => !c);

        return (
            <div className="power-strip">
                <div className="power-strip-header">
                    <span className="zone-label">POWER SLOTS ({hand.length}/{MAX_POWER_HAND_SLOTS})</span>
                    <span className="stat-chip">USE OR KEEP</span>
                    <span className="stat-chip">DECK {deckCount}</span>
                    {(me?.attackBonus > 0 || me?.defenseBonus > 0) && (
                        <span className="stat-chip boost-chip">
                            {me.attackBonus > 0 ? `ATK +${me.attackBonus}` : ''}
                            {me.attackBonus > 0 && me.defenseBonus > 0 ? ' · ' : ''}
                            {me.defenseBonus > 0 ? `DEF +${me.defenseBonus}` : ''}
                        </span>
                    )}
                </div>
                <div className="hand-cards power-hand-strip">
                    {slots.map((card, idx) => {
                        if (card) {
                            return (
                                <div key={card.id} className="hand-card-wrapper">
                                    <Card
                                        card={card}
                                        size="sm"
                                        isPlayable={canPlayPower}
                                        onClick={canPlayPower ? handlePowerClick : undefined}
                                    />
                                    <span className="cp-chip power-chip">
                                        {canPlayPower ? 'TAP TO USE' : 'KEEP'}
                                    </span>
                                </div>
                            );
                        }

                        const isNextDrawSlot = idx === nextEmptyIndex;
                        const canFill = canDrawPower && isNextDrawSlot;
                        return (
                            <div key={`empty-power-${idx}`} className="hand-card-wrapper">
                                <button
                                    type="button"
                                    className={`card empty-slot size-sm power-draw-slot ${canFill ? 'playable' : ''}`}
                                    onClick={() => canFill && actions.drawCard()}
                                    disabled={!canFill}
                                    title={
                                        me?.hasDrawn
                                            ? 'Already drew this turn — keep cards for next turn'
                                            : hand.length >= MAX_POWER_HAND_SLOTS
                                              ? 'All 4 power slots are full'
                                              : deckCount === 0
                                                ? 'Power deck is empty'
                                                : isNextDrawSlot
                                                  ? 'Draw 1 power card into this slot (once per turn)'
                                                  : 'Fill earlier empty slots first'
                                    }
                                >
                                    {me?.hasDrawn ? 'EMPTY' : canFill ? 'DRAW' : 'SLOT'}
                                </button>
                                <span className="cp-chip power-chip">
                                    {canFill ? '1/TURN' : `SLOT ${idx + 1}`}
                                </span>
                            </div>
                        );
                    })}
                </div>
            </div>
        );
    };

    const renderActiveOptions = () => {
        if (!me?.activePokemon) {
            return (
                <div className="rail-options">
                    <span className="zone-label">YOUR OPTIONS</span>
                    <p className="rail-options-hint">
                        {needsPromote ? 'PROMOTE FROM BACK LINE' : 'WAITING…'}
                    </p>
                </div>
            );
        }

        return (
            <div className="rail-options">
                <div className="rail-options-header">
                    <span className="zone-label">YOUR OPTIONS</span>
                    <button
                        type="button"
                        className="pixel-btn details-btn"
                        onClick={() => openDetails(me.activePokemon, `YOU (${me.id})`)}
                    >
                        DETAILS
                    </button>
                </div>
                <div className="rail-options-actions">
                    {isMyTurn && !needsPromote && !cpuThinking ? (
                        <>
                            <button
                                type="button"
                                className="pixel-btn primary draw-btn"
                                onClick={actions.drawCard}
                                disabled={!canDrawPower}
                            >
                                {me.hasDrawn
                                    ? 'KEEP FOR NEXT'
                                    : (me.hand?.length || 0) >= MAX_POWER_HAND_SLOTS
                                      ? 'SLOTS FULL'
                                      : 'DRAW POWER'}
                            </button>
                            <button
                                type="button"
                                className="pixel-btn primary attach-btn"
                                onClick={() => actions.attachEnergy('')}
                                disabled={me.hasAttached}
                            >
                                CHARGE
                            </button>
                            {me.activePokemon.attacks?.map((att, i) => (
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
                            <button
                                type="button"
                                className="pixel-btn end-turn-btn"
                                onClick={actions.endTurn}
                            >
                                END TURN
                            </button>
                        </>
                    ) : (
                        <p className="rail-options-hint">
                            {cpuThinking ? 'CPU TURN…' : isMyTurn ? 'PROMOTE FIRST' : 'WAIT FOR TURN'}
                        </p>
                    )}
                </div>
            </div>
        );
    };

    const turnLabel = needsPromote
        ? 'PROMOTE!'
        : cpuThinking
          ? 'CPU TURN'
          : isMyTurn
            ? 'YOUR TURN'
            : 'WAIT';

    return (
        <div className="game-board">
            <div className={`board-menu ${menuOpen ? 'open' : ''}`}>
                <button
                    type="button"
                    className="pixel-btn board-menu-toggle"
                    aria-expanded={menuOpen}
                    aria-haspopup="menu"
                    onClick={() => setMenuOpen((open) => !open)}
                >
                    MENU {menuOpen ? '▴' : '▾'}
                </button>
                {menuOpen && (
                    <div className="board-menu-dropdown pixel-panel" role="menu">
                        {!isPractice && (
                            <button
                                type="button"
                                className="board-menu-item"
                                role="menuitem"
                                onClick={() => runMenuAction(() =>
                                    setPlayerId(playerId === 'player1' ? 'player2' : 'player1')
                                )}
                            >
                                SWAP
                            </button>
                        )}
                        <button
                            type="button"
                            className="board-menu-item"
                            role="menuitem"
                            onClick={() => runMenuAction(() =>
                                actions.startGame({ vsCPU: isPractice })
                            )}
                        >
                            NEW GAME
                        </button>
                        {onShowRules && (
                            <button
                                type="button"
                                className="board-menu-item"
                                role="menuitem"
                                onClick={() => runMenuAction(onShowRules)}
                            >
                                RULES
                            </button>
                        )}
                    </div>
                )}
            </div>

            <header className="card-rail top-rail pixel-panel">
                <div className="rail-meta">
                    <span className="player-name">VS {opponentLabel}</span>
                        <div className="stat-row">
                            <span className="stat-chip">WINS {opponent?.gamesWon || 0}</span>
                            <span className="stat-chip">LEFT {(opponent?.activePokemon ? 1 : 0) + (opponent?.benchedPokemon?.length || 0)}</span>
                            <span className="stat-chip">SHIELD {opponent?.protectShields ?? 0}</span>
                            <span className="stat-chip">PWR {(opponent?.hand?.length ?? 0)}</span>
                            {(opponent?.attackBonus > 0 || opponent?.defenseBonus > 0) && (
                                <span className="stat-chip boost-chip">
                                    {opponent.attackBonus > 0 ? `ATK +${opponent.attackBonus}` : ''}
                                    {opponent.attackBonus > 0 && opponent.defenseBonus > 0 ? ' · ' : ''}
                                    {opponent.defenseBonus > 0 ? `DEF +${opponent.defenseBonus}` : ''}
                                </span>
                            )}
                        </div>
                </div>
            </header>

            <div className="mid-row">
                <aside className="message-panel pixel-screen">
                    <h3>■ STATUS</h3>
                    <p className={`msg-turn ${isMyTurn ? 'glow' : ''}`}>{turnLabel}</p>
                    <p className="msg-meta">
                        G{gameState.gameNumber} · T{gameState.turnNumber}
                    </p>
                    <p className={`msg-hint ${isMyTurn && !needsPromote ? 'animate-insert-coin' : ''}`}>
                        {statusHint}
                    </p>
                    {gameState.lastAction && (
                        <p className="msg-last">{gameState.lastAction}</p>
                    )}
                    {error && <p className="msg-error">{error}</p>}
                    <div className="msg-you">
                        <span className="player-name">YOU ({me?.id})</span>
                        <div className="stat-row">
                            <span className="stat-chip">WINS {me?.gamesWon || 0}</span>
                            <span className="stat-chip">SHIELD {me?.protectShields ?? 0}</span>
                            {(me?.attackBonus > 0 || me?.defenseBonus > 0) && (
                                <span className="stat-chip boost-chip">
                                    {me.attackBonus > 0 ? `ATK +${me.attackBonus}` : ''}
                                    {me.attackBonus > 0 && me.defenseBonus > 0 ? ' · ' : ''}
                                    {me.defenseBonus > 0 ? `DEF +${me.defenseBonus}` : ''}
                                </span>
                            )}
                        </div>
                    </div>
                </aside>

                <section className="active-arena">
                    <div className="arena-faceoff">
                        <div className="arena-card">
                            {oppActiveDisplay ? (
                                <>
                                    <Card card={oppActiveDisplay} size="lg" isActive={true} />
                                    {renderActiveHp(oppActiveDisplay)}
                                </>
                            ) : (
                                <div className="card empty-slot size-lg">OPP</div>
                            )}
                        </div>
                        <span className="vs-badge">VS</span>
                        <div className="arena-card">
                            {me?.activePokemon ? (
                                <>
                                    <Card card={me.activePokemon} size="lg" isActive={true} />
                                    {renderActiveHp(me.activePokemon)}
                                </>
                            ) : (
                                <div className="card empty-slot size-lg">PROMOTE</div>
                            )}
                        </div>
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

            <footer className="card-rail bottom-rail pixel-panel">
                <div className="rail-cards">
                    {renderBench(me, true)}
                    {renderActiveOptions()}
                    {renderPowerStrip()}
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
