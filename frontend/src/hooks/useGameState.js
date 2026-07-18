import { useState, useCallback, useEffect } from 'react';
import axios from 'axios';

const API_BASE = '/api/game';

export function useGameState() {
    const [gameState, setGameState] = useState(null);
    const [actionLog, setActionLog] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [playerId, setPlayerId] = useState('player1');

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

    const startGame = async () => {
        try {
            const res = await axios.post(`${API_BASE}/start`);
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to start game');
        }
    };

    const drawCard = async () => {
        try {
            const res = await axios.post(`${API_BASE}/draw`, { playerId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to draw card');
        }
    };

    const playBench = async (cardId) => {
        try {
            const res = await axios.post(`${API_BASE}/play-bench`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to play to bench');
        }
    };

    const setActive = async (cardId) => {
        try {
            const res = await axios.post(`${API_BASE}/set-active`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to set active Pokemon');
        }
    };

    const attachEnergy = async (cardId) => {
        try {
            const res = await axios.post(`${API_BASE}/attach-energy`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to attach energy');
        }
    };

    const attack = async (attackIndex) => {
        try {
            const res = await axios.post(`${API_BASE}/attack`, { playerId, attackIndex });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to attack');
        }
    };

    const endTurn = async () => {
        try {
            const res = await axios.post(`${API_BASE}/end-turn`, { playerId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to end turn');
        }
    };

    const promote = async (cardId) => {
        try {
            const res = await axios.post(`${API_BASE}/promote`, { playerId, cardId });
            applyResult(res.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to promote');
        }
    };

    const isMyTurn = gameState?.currentTurn === playerId;
    const me = gameState?.players?.find(p => p.id === playerId);
    const opponent = gameState?.players?.find(p => p.id !== playerId);
    const needsPromote = Boolean(me && !me.activePokemon && me.benchedPokemon?.length > 0);

    return {
        gameState,
        actionLog,
        loading,
        error,
        playerId,
        setPlayerId,
        me,
        opponent,
        isMyTurn,
        needsPromote,
        actions: {
            fetchState,
            startGame,
            drawCard,
            playBench,
            setActive,
            attachEnergy,
            attack,
            endTurn,
            promote
        }
    };
}
