import React, { useState, useEffect } from 'react';
import { useGameState } from '../hooks/useGameState';
import { Card } from './Card';
import { CardDetail } from './CardDetail';
import { PartySelect } from './PartySelect';
import { PowerReplace } from './PowerReplace';
import { MAX_POWER_HAND_SLOTS } from '../rules/handbook';
import './GameBoard.css';

export const GameBoard = ({ onShowRules }) => {
    const {
        gameState, actionLog, loading, error, isMyTurn, me, opponent,
        actions, setPlayerId, playerId, needsPromote, needsPartySelect,
        vsCPU, setVsCPU, isPractice, cpuThinking, pendingDraw
    } = useGameState();
    const [detailView, setDetailView] = useState(null);
    const [confirmingParty, setConfirmingParty] = useState(false);
    const [menuOpen, setMenuOpen] = useState(false);
    const [resolvingReplace, setResolvingReplace] = useState(false);
    const [switchMode, setSwitchMode] = useState(false);

    const isMatchOver = gameState?.status === 'GameOver' || gameState?.phase === 'MatchOver';
    const needsPowerReplace = Boolean(pendingDraw?.length);
    const canSwitch =
        !isMatchOver &&
        isMyTurn &&
        !needsPromote &&
        !needsPowerReplace &&
        !cpuThinking &&
        !!me?.activePokemon &&
        !me?.hasSwitched &&
        (me?.benchedPokemon?.length ?? 0) > 0;

    // Hooks must run unconditionally (before early returns).
    useEffect(() => {
        if (switchMode && !canSwitch) {
            setSwitchMode(false);
        }
    }, [switchMode, canSwitch]);

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

    const drawnPower = pendingDraw?.[0] || null;

    const resultHeadline =
        gameState.winner === gameState.cpuPlayerId
            ? 'CPU WINS!'
            : gameState.winner === playerId
              ? 'YOU WIN!'
              : gameState.winner
                ? `${gameState.winner} WINS!`
                : 'MATCH OVER';

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

    const handleBenchCardClick = (card, ownerLabel) => {
        if (!card) return;

        if (needsPromote && !isMatchOver && !needsPowerReplace) {
            actions.promote(card.id);
            return;
        }

        if (switchMode && canSwitch) {
            actions.switchActive(card.id);
            setSwitchMode(false);
            return;
        }

        openDetails(card, ownerLabel);
    };

    const renderBench = (player, isMe) => {
        const benchSlots = Array(2).fill(null);
        player?.benchedPokemon?.forEach((card, i) => {
            if (i < 2) benchSlots[i] = card;
        });

        const pickingPromote = isMe && needsPromote && !isMatchOver && !needsPowerReplace;
        const pickingSwitch = isMe && switchMode && canSwitch;
        const selectable = pickingPromote || pickingSwitch;
        const ownerLabel = isMe ? `YOU (${me?.id})` : opponentLabel;

        return (
            <div className={`bench-area ${isMe ? 'my-bench' : 'opponent-bench'} ${pickingSwitch ? 'switch-picking' : ''}`}>
                <span className="zone-label">
                    {pickingSwitch
                        ? 'TAP A CARD TO SWITCH IN'
                        : pickingPromote
                          ? 'TAP A CARD TO PROMOTE'
                          : isMe
                            ? 'YOUR BACK LINE'
                            : 'OPP BACK LINE'}
                </span>
                <div className="bench-slots">
                    {benchSlots.map((card, idx) => (
                        <div key={idx} className="bench-slot-wrapper">
                            {card ? (
                                <Card
                                    card={
                                        isMe
                                            ? card
                                            : {
                                                ...card,
                                                energyAttached: undefined,
                                                hp: undefined,
                                                maxHp: undefined,
                                                stats: undefined,
                                            }
                                    }
                                    size="sm"
                                    compact
                                    className="tap-details"
                                    isPlayable={selectable}
                                    onClick={() => handleBenchCardClick(card, ownerLabel)}
                                />
                            ) : (
                                <div className="card empty-slot size-sm compact">BACK</div>
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

    const statusHint = isMatchOver
        ? 'MATCH COMPLETE'
        : needsPromote
          ? 'PICK BACK-LINE TO PROMOTE (§6.3)'
          : needsPowerReplace
            ? 'POWER SLOTS FULL — SWAP OR KEEP HAND'
            : switchMode
              ? 'TAP A BACK-LINE CARD TO SWITCH IN'
              : isMyTurn
                ? 'USE POWER · CHARGE · SWITCH · ATTACK'
                : cpuThinking
                  ? 'CPU IS PLAYING…'
                  : 'WAITING FOR OPPONENT';

    const canPlayPower =
        !isMatchOver &&
        isMyTurn &&
        !needsPromote &&
        !needsPowerReplace &&
        !cpuThinking &&
        !!me?.activePokemon;

    const isPowerCardPlayable = (card) => {
        if (!canPlayPower || !card) return false;
        if (card.effect === 'heal' && me?.activePokemon?.hp >= me?.activePokemon?.maxHp) {
            return false;
        }
        return true;
    };

    const handlePowerClick = (card) => {
        if (!isPowerCardPlayable(card)) return;
        actions.playPower(card.id);
    };

    const resolveReplace = async (cardId) => {
        setResolvingReplace(true);
        try {
            await actions.selectDraw(cardId);
        } finally {
            setResolvingReplace(false);
        }
    };

    const renderPowerStrip = () => {
        const hand = me?.hand || [];
        const deckCount = me?.powerDeck?.length ?? 0;
        const slots = Array.from({ length: MAX_POWER_HAND_SLOTS }, (_, i) => hand[i] || null);

        return (
            <div className="power-strip">
                <div className="power-strip-header">
                    <span className="zone-label">POWER SLOTS ({hand.length}/{MAX_POWER_HAND_SLOTS})</span>
                    <span className="stat-chip">AUTO DRAW</span>
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
                            const playable = isPowerCardPlayable(card);
                            return (
                                <div key={card.id} className="hand-card-wrapper">
                                    <Card
                                        card={card}
                                        size="sm"
                                        isPlayable={playable}
                                        onClick={playable ? handlePowerClick : undefined}
                                    />
                                    <span className="cp-chip power-chip">
                                        {playable ? 'TAP TO USE' : 'HOLD'}
                                    </span>
                                </div>
                            );
                        }

                        return (
                            <div key={`empty-power-${idx}`} className="hand-card-wrapper">
                                <div
                                    className="card empty-slot size-sm power-draw-slot"
                                    title="Empty slot — filled automatically at turn start"
                                >
                                    {me?.hasDrawn ? 'EMPTY' : 'AUTO'}
                                </div>
                                <span className="cp-chip power-chip">SLOT {idx + 1}</span>
                            </div>
                        );
                    })}
                </div>
            </div>
        );
    };

    const renderActiveOptions = () => {
        const active = me?.activePokemon;
        const myTurnLive =
            !isMatchOver && isMyTurn && !needsPromote && !needsPowerReplace && !cpuThinking;

        if (!active) {
            return (
                <div className="rail-options">
                    <p className="rail-options-hint">
                        {isMatchOver ? 'MATCH OVER' : needsPromote ? 'PROMOTE FROM BACK LINE' : 'WAITING…'}
                    </p>
                </div>
            );
        }

        return (
            <div className="rail-options">
                <div className="rail-options-actions">
                    {myTurnLive ? (
                        <>
                            <button
                                type="button"
                                className="pixel-btn primary attach-btn"
                                onClick={() => {
                                    setSwitchMode(false);
                                    actions.attachEnergy('');
                                }}
                                disabled={me.hasAttached}
                            >
                                CHARGE
                            </button>
                            <button
                                type="button"
                                className={`pixel-btn switch-btn ${switchMode ? 'is-active' : ''}`}
                                onClick={() => setSwitchMode((on) => !on)}
                                disabled={!canSwitch && !switchMode}
                                title={
                                    me.hasSwitched
                                        ? 'Already switched this turn'
                                        : !(me?.benchedPokemon?.length > 0)
                                          ? 'No back-line Pokémon to switch in'
                                          : switchMode
                                            ? 'Cancel switch'
                                            : 'Switch active with a back-line Pokémon'
                                }
                            >
                                {switchMode ? 'CANCEL' : me.hasSwitched ? 'SWITCHED' : 'SWITCH'}
                            </button>
                            {(active.attacks || []).map((att, i) => (
                                <button
                                    key={i}
                                    type="button"
                                    className="pixel-btn danger attack-button"
                                    onClick={() => {
                                        setSwitchMode(false);
                                        actions.attack(i);
                                    }}
                                    disabled={
                                        !opponent?.activePokemon ||
                                        (active.energyAttached || 0) < att.cost
                                    }
                                >
                                    {att.name} · {att.damage}
                                </button>
                            ))}
                            <button
                                type="button"
                                className="pixel-btn end-turn-btn"
                                onClick={() => {
                                    setSwitchMode(false);
                                    actions.endTurn();
                                }}
                            >
                                END TURN
                            </button>
                        </>
                    ) : (
                        <p className="rail-options-hint">
                            {needsPowerReplace
                                ? 'RESOLVE POWER DRAW…'
                                : cpuThinking
                                  ? 'CPU TURN…'
                                  : isMyTurn
                                    ? 'PROMOTE FIRST'
                                    : 'WAIT FOR TURN'}
                        </p>
                    )}
                </div>
            </div>
        );
    };

    const turnLabel = isMatchOver
        ? 'GAME OVER'
        : needsPromote
          ? 'PROMOTE!'
          : needsPowerReplace
            ? 'SWAP POWER'
            : switchMode
              ? 'PICK SWITCH'
              : cpuThinking
                ? 'CPU TURN'
                : isMyTurn
                  ? 'YOUR TURN'
                  : 'WAIT';

    return (
        <div className={`game-board ${isMatchOver ? 'match-over' : ''}`}>
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
                                    <Card
                                        card={oppActiveDisplay}
                                        size="lg"
                                        isActive={true}
                                        className="tap-details"
                                        onClick={() => openDetails(oppActiveDisplay, opponentLabel)}
                                    />
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
                                    <Card
                                        card={me.activePokemon}
                                        size="lg"
                                        isActive={true}
                                        className="tap-details"
                                        onClick={() => openDetails(me.activePokemon, `YOU (${me.id})`)}
                                    />
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

            {needsPowerReplace && drawnPower && (
                <PowerReplace
                    drawn={drawnPower}
                    hand={me?.hand || []}
                    resolving={resolvingReplace}
                    onReplace={(cardId) => resolveReplace(cardId)}
                    onKeep={() => resolveReplace('_keep')}
                />
            )}

            {isMatchOver && (
                <div
                    className="game-result-backdrop"
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby="game-result-title"
                >
                    <div className="game-result-modal pixel-panel animate-slam-in">
                        <div className="arcade-marquee game-result-marquee">
                            {isPractice ? 'PRACTICE OVER' : 'MATCH OVER · §6.4'}
                        </div>
                        <h2 id="game-result-title" className="game-result-title">
                            GAME OVER
                        </h2>
                        <p className="winner-banner">{resultHeadline}</p>
                        {gameState.lastAction && (
                            <p className="game-result-detail">{gameState.lastAction}</p>
                        )}
                        <div className="game-result-score">
                            <span className="stat-chip">YOU {me?.gamesWon || 0}</span>
                            <span className="stat-chip">
                                {isPractice ? 'CPU' : 'OPP'} {opponent?.gamesWon || 0}
                            </span>
                        </div>
                        <div className="game-result-actions">
                            <button
                                type="button"
                                className="pixel-btn primary start-btn"
                                onClick={() => actions.startGame({ vsCPU: isPractice })}
                            >
                                {isPractice ? 'PRACTICE AGAIN' : 'PLAY AGAIN'}
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
            )}
        </div>
    );
};
