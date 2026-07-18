import { useState, useCallback, useEffect } from 'react';
import axios from 'axios';
import { ACTIONS } from '../rules/handbook';
import { validateAction } from '../rules/validateAction';

const API_BASE = '/api/game';

export function useGameState() {
    const [gameState, setGameState] = useState(null);
    const [actionLog, setActionLog] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [playerId, setPlayerId] = useState('player1');
    const [vsCPU, setVsCPU] = useState(true);

    const fetchActions = useCallback(async () => {
        try {
            const res = await axios.get(`${API_BASE}/actions`, { params: { limit: 30 } });
            setActionLog(Array.isArray(res.data) ? res.data : []);
        } catch {
            // no active game yet
        }
    }, []);

    const fetchState = useCallback(async () => {
        try {
            const res = await axios.get(API_BASE);
            setGameState(res.data);
            setError(null);
            if (res.data?.id) {
                await fetchActions();
            }
        } catch (err) {
            if (err.response?.status !== 404) {
                setError(err.message || 'Failed to fetch game state');
            }
        } finally {
            setLoading(false);
        }
    }, [fetchActions]);

    useEffect(() => {
        fetchState();
        const interval = setInterval(fetchState, 2000);
        return () => clearInterval(interval);
    }, [fetchState]);

    const applyResult = (data) => {
        setGameState(data);
        setError(null);
        fetchActions();
    };

    const guard = (action, payload = {}) => {
        const result = validateAction({
            gameState,
            playerId,
            action,
            payload,
        });
        if (!result.ok) {
            setError(result.error);
            return false;
        }
        return true;
    };

    const startGame = async (opts = {}) => {
        const practice = opts.vsCPU ?? vsCPU;
        if (!guard(ACTIONS.START_GAME)) return;
        try {
            if (practice) {
                setPlayerId('player1');
            }
            const res = await axios.post(`${API_BASE}/start`, { vsCPU: !!practice });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to start game');
        }
    };

    const selectParty = async (cardIds) => {
        if (!guard(ACTIONS.SELECT_PARTY, { cardIds })) return;
        try {
            const res = await axios.post(`${API_BASE}/select-party`, { playerId, cardIds });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to select party');
        }
    };

    const drawCard = async () => {
        if (!guard(ACTIONS.DRAW_CARD)) return;
        try {
            const res = await axios.post(`${API_BASE}/draw`, { playerId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to draw card');
        }
    };

    const selectDraw = async (cardId) => {
        if (!guard(ACTIONS.SELECT_DRAW, { cardId })) return;
        try {
            const res = await axios.post(`${API_BASE}/draw/select`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to select draw card');
        }
    };

    const playBench = async (cardId) => {
        if (!guard(ACTIONS.PLAY_BENCH, { cardId })) return;
        try {
            const res = await axios.post(`${API_BASE}/play-bench`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to play to bench');
        }
    };

    const setActive = async (cardId) => {
        if (!guard(ACTIONS.SET_ACTIVE, { cardId })) return;
        try {
            const res = await axios.post(`${API_BASE}/set-active`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to set active Pokemon');
        }
    };

    const attachEnergy = async (cardId = '') => {
        if (!guard(ACTIONS.ATTACH_ENERGY, { cardId })) return;
        try {
            const res = await axios.post(`${API_BASE}/attach-energy`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to charge energy');
        }
    };

    const attack = async (attackIndex) => {
        if (!guard(ACTIONS.ATTACK, { attackIndex })) return;
        try {
            const res = await axios.post(`${API_BASE}/attack`, { playerId, attackIndex });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to attack');
        }
    };

    const endTurn = async () => {
        if (!guard(ACTIONS.END_TURN)) return;
        try {
            const res = await axios.post(`${API_BASE}/end-turn`, { playerId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to end turn');
        }
    };

    const promote = async (cardId) => {
        if (!guard(ACTIONS.PROMOTE, { cardId })) return;
        try {
            const res = await axios.post(`${API_BASE}/promote`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to promote');
        }
    };

    const isMyTurn = gameState?.currentTurn === playerId;
    const me = gameState?.players?.find((p) => p.id === playerId);
    const opponent = gameState?.players?.find((p) => p.id !== playerId);
    const needsPromote = Boolean(me && !me.activePokemon && me.benchedPokemon?.length > 0);
    const pendingDraw = me?.pendingDraw?.length > 0 ? me.pendingDraw : null;
    const needsPartySelect =
        ['TeamPreview', 'PartySelect', 'BetweenGames'].includes(gameState?.phase) && !me?.partyReady;
    const isPractice = Boolean(gameState?.vsCPU ?? vsCPU);
    const cpuThinking =
        isPractice &&
        gameState?.phase === 'InBattle' &&
        gameState?.currentTurn === gameState?.cpuPlayerId;

    return {
        gameState,
        actionLog,
        loading,
        error,
        playerId,
        setPlayerId,
        vsCPU,
        setVsCPU,
        isPractice,
        cpuThinking,
        me,
        opponent,
        isMyTurn,
        needsPromote,
        pendingDraw,
        needsPartySelect,
        actions: {
            fetchState,
            startGame,
            selectParty,
            drawCard,
            selectDraw,
            playBench,
            setActive,
            attachEnergy,
            attack,
            endTurn,
            promote,
        },
    };
}
