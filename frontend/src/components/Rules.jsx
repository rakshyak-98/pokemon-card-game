import React from 'react';
import './Rules.css';
import { HANDBOOK_REVISION, GREAT_LEAGUE_CP_CAP, BATTLE_PARTY_SIZE, MAX_BATTLE_TEAM_SIZE, MIN_BATTLE_TEAM_SIZE } from '../rules/handbook';

const SECTIONS = [
  {
    id: 'source',
    marquee: '00 · SOURCE',
    title: 'Official handbook',
    body: (
      <>
        <p>
          Rules follow the <strong>Play! Pokémon Pokémon GO Tournament Handbook</strong> (revision{' '}
          {HANDBOOK_REVISION}). This app adapts Great League match play into a turn-based duel.
        </p>
      </>
    ),
  },
  {
    id: 'objective',
    marquee: '01 · WIN',
    title: 'How to win a match',
    body: (
      <>
        <p>
          Matches are <strong>best-of-three</strong> (§6.4). Win a game by knocking out the opponent&apos;s
          last Pokémon in their battle party of three (§6.3). Win the match by taking two games.
        </p>
      </>
    ),
  },
  {
    id: 'team',
    marquee: '02 · TEAM',
    title: 'Battle Team construction (§3.1)',
    body: (
      <>
        <ul>
          <li>
            Register a Battle Team of <strong>{MIN_BATTLE_TEAM_SIZE}–{MAX_BATTLE_TEAM_SIZE}</strong> Pokémon.
          </li>
          <li>
            Great League: each Pokémon must be <strong>≤ {GREAT_LEAGUE_CP_CAP} CP</strong>.
          </li>
          <li>No two Pokémon with the same National Pokédex number.</li>
          <li>Banned: Ditto, Shedinja, Xerneas, Yveltal (plus any Ban List updates).</li>
          <li>No Mega Evolved / Primal forms. At most one Best Buddy CP boost.</li>
          <li>
            Each game, choose any <strong>{BATTLE_PARTY_SIZE}</strong> from your team to battle (§3.1 / §6.1).
          </li>
        </ul>
      </>
    ),
  },
  {
    id: 'preview',
    marquee: '03 · PREVIEW',
    title: 'Team preview (§6.1 / §6.5.1)',
    body: (
      <>
        <p>
          Before each game, exchange team preview lists (species, CP, moves, form flags). Preview time is
          limited to two minutes in live events. Public info also includes Protect Shields remaining, Pokémon
          left, and types.
        </p>
        <p>
          <strong>Private</strong> (§6.5.3): current energy, switch timer, and unrevealed back-line details —
          do not demand these from your opponent.
        </p>
      </>
    ),
  },
  {
    id: 'battle',
    marquee: '04 · BATTLE',
    title: 'During a game',
    body: (
      <>
        <ol>
          <li>Lead with your first selected Pokémon; the other two form your back line.</li>
          <li>
            <strong>Power cards</strong> draw automatically once per turn into an empty slot (max 4)
            and keep dealing until the game ends (the special-power deck refills when empty).
            If slots are full, choose a card to swap out—or keep your hand and discard the draw.
            Tap any held cards to use them on your Active this turn (Attack, Defense, or Heal)—you may play
            multiple—or leave them for later.
          </li>
          <li>
            <strong>Charge Energy</strong> once per turn on your Active (approximates Fast Attack energy gain).
          </li>
          <li>
            <strong>Attack</strong> when you have enough energy, or <strong>End Turn</strong>.
          </li>
          <li>On your turn, tap a back-line Pokémon to switch it with your active (once per turn).</li>
          <li>After a KO, promote from the back line. If none remain, you lose the game (§6.3).</li>
        </ol>
      </>
    ),
  },
  {
    id: 'controls',
    marquee: '05 · CONTROLS',
    title: 'How to play in this UI',
    body: (
      <>
        <ul>
          <li>Start a match → review opponent preview → lock a party of 3.</li>
          <li>Tap an Active Pokémon in the arena to open its card details.</li>
          <li>
            Drawn power cards fill up to 4 slots on the bottom-right automatically each turn — tap any to use on Active (you may play several), or keep for next turn. When slots are full, a popup lets you swap or discard the new draw.
          </li>
          <li>SWAP switches which seat you control on this local client (2-player only).</li>
          <li>Between games, both players select a new party of 3 from the same Battle Team.</li>
        </ul>
      </>
    ),
  },
  {
    id: 'practice',
    marquee: '06 · PRACTICE',
    title: 'Practice vs CPU',
    body: (
      <>
        <p>
          Enable <strong>Practice vs CPU</strong> on the start screen to learn the flow solo. The CPU
          picks a legal party, draws and plays power cards, charges energy, attacks when it can, and
          promotes after knockouts — so you can focus on handbook timing without a second player.
        </p>
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
          Pokémon GO Championship Series rules — Great League, best-of-three, team of six, party of three.
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
