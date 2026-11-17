import { Check, CircleX, Flag } from 'lucide-react'

export type RulerState = 'empty' | 'answered' | 'flagged' | 'ok' | 'bad'

export interface RulerCell {
  pos: number
  state: RulerState
  /** Spoken name, e.g. "Soal 4, terjawab". */
  label: string
}

// A row of numbered cells, one per question. With onPick the cells are buttons; without it they only show state.
export function Ruler({
  cells,
  current,
  label,
  onPick,
}: {
  cells: RulerCell[]
  /** pos of the question on screen; 0 when none. */
  current: number
  label: string
  onPick?: (pos: number) => void
}) {
  return (
    <div className="grid ruler" role="group" aria-label={label}>
      {cells.map((c) => {
        const cur = c.pos === current
        const cls = `cell ${c.state}${cur ? ' cur' : ''}`
        const mark =
          c.state === 'flagged' ? (
            <Flag className="ic" aria-hidden="true" />
          ) : c.state === 'bad' ? (
            <CircleX className="ic" aria-hidden="true" />
          ) : c.state === 'answered' || c.state === 'ok' ? (
            <Check className="ic" aria-hidden="true" />
          ) : null
        return onPick ? (
          <button key={c.pos} className={cls} aria-label={c.label} aria-current={cur} onClick={() => onPick(c.pos)}>
            {c.pos}
            {mark}
          </button>
        ) : (
          <span key={c.pos} className={cls} role="img" aria-label={c.label} aria-current={cur}>
            {c.pos}
            {mark}
          </span>
        )
      })}
    </div>
  )
}
