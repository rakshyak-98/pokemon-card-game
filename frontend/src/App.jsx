import React, { useState } from 'react'
import { GameBoard } from './components/GameBoard'
import { Rules } from './components/Rules'

function App() {
  const [view, setView] = useState('game')

  if (view === 'rules') {
    return <Rules onBack={() => setView('game')} />
  }

  return <GameBoard onShowRules={() => setView('rules')} />
}

export default App
