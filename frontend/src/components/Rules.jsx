import React from 'react';
import './Rules.css';

const SECTIONS = [
  {
    id: 'objective',
    marquee: '01 · WIN',
    title: 'How to win',
    body: (
      <>
        <p>
          Two players duel until one side claims victory. You win by either:
        </p>
        <ul>
          <li>Taking all <strong>6 prize cards</strong></li>
          <li>Knocking out the opponent&apos;s last Pokémon (no Active and empty Bench)</li>
        </ul>
      </>
    ),
  },
  {
    id: 'setup',
    marquee: '02 · DEAL',
    title: 'Match setup',
    body: (
      <>
        <p>When the match starts, each player receives:</p>
        <ul>
          <li>A shuffled 60-card deck (20 Pokémon + 40 Energy)</li>
          <li><strong>6 prize cards</strong> set aside face-down</li>
          <li>A <strong>7-card opening hand</strong></li>
        </ul>
        <p>
          Player 1 goes first. Before battling, each side must set an{' '}
          <strong>Active Pokémon</strong> from their hand.
        </p>
      </>
    ),
  },
  {
    id: 'zones',
    marquee: '03 · FIELD',
    title: 'The battlefield',
    body: (
      <>
        <ul>
          <li>
            <strong>Active</strong> — your fighter in play. Only this Pokémon attacks and
            receives Energy.
          </li>
          <li>
            <strong>Bench</strong> — up to 5 backup Pokémon. They wait until promoted.
          </li>
          <li>
            <strong>Hand / Deck / Discard</strong> — cards you hold, draw from, or send after a KO.
          </li>
          <li>
            <strong>Prizes</strong> — take one into your hand each time you knock out an
            opposing Active Pokémon.
          </li>
        </ul>
      </>
    ),
  },
  {
    id: 'turn',
    marquee: '04 · TURN',
    title: 'Your turn',
    body: (
      <>
        <p>On your turn you may:</p>
        <ol>
          <li>
            <strong>Draw</strong> one card (once per turn).
          </li>
          <li>
            <strong>Set Active</strong> if you have none — from hand or Bench.
          </li>
          <li>
            <strong>Bench</strong> Pokémon from your hand (max 5 on the Bench).
          </li>
          <li>
            <strong>Attach Energy</strong> to your Active Pokémon (once per turn).
          </li>
          <li>
            <strong>Attack</strong> if your Active has enough Energy — this ends your turn —
            or press <strong>End Turn</strong> to pass.
          </li>
        </ol>
        <p>
          If your Active was knocked out last turn, you must{' '}
          <strong>Promote</strong> a Benched Pokémon before doing anything else.
        </p>
      </>
    ),
  },
  {
    id: 'battle',
    marquee: '05 · FIGHT',
    title: 'Attacking & knockouts',
    body: (
      <>
        <ul>
          <li>
            Each attack lists a damage value and an Energy <strong>cost</strong>. Your Active
            needs at least that many Energy attached.
          </li>
          <li>
            Damage reduces the opponent&apos;s Active HP. At 0 HP they are knocked out and
            discarded.
          </li>
          <li>
            The attacker takes the top prize card into their hand. If that was the 6th prize,
            they win immediately.
          </li>
          <li>
            The defender must promote from the Bench. If the Bench is empty, the attacker wins.
          </li>
        </ul>
      </>
    ),
  },
  {
    id: 'controls',
    marquee: '06 · CONTROLS',
    title: 'How to play in the UI',
    body: (
      <>
        <ul>
          <li>Select a card in your hand, then click an empty Bench slot or Active zone.</li>
          <li>Select an Energy, then attach it to your Active (or use Attach Energy).</li>
          <li>Use the attack buttons on your Active when Energy cost is met.</li>
          <li>
            After a KO, click a Benched Pokémon to Promote it.
          </li>
          <li>
            Use <strong>Swap P1/P2</strong> to play both sides on one screen.
          </li>
        </ul>
      </>
    ),
  },
];

export const Rules = ({ onBack }) => {
  return (
    <div className="rules-page">
      <header className="rules-hero pixel-panel animate-slam-in">
        <div className="arcade-marquee">INSTRUCTION MANUAL</div>
        <h1 className="game-title">HOW TO PLAY</h1>
        <p className="rules-lead">
          A simplified Pokémon TCG duel — draw, energy up, attack, take prizes.
        </p>
        <button type="button" className="pixel-btn primary" onClick={onBack}>
          Back to game
        </button>
      </header>

      <div className="rules-sections">
        {SECTIONS.map((section) => (
          <section key={section.id} className="rules-section pixel-panel" id={section.id}>
            <div className="rules-section-marquee">{section.marquee}</div>
            <h2>{section.title}</h2>
            <div className="rules-section-body">{section.body}</div>
          </section>
        ))}
      </div>

      <footer className="rules-footer">
        <p className="insert-coin animate-insert-coin">READY PLAYER ONE?</p>
        <button type="button" className="pixel-btn primary" onClick={onBack}>
          Press start
        </button>
      </footer>
    </div>
  );
};
