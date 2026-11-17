import { Link } from 'react-router-dom'
import { queueReason } from './format'
import type { QueueItem } from './types'

// The first queue item is the practice to do next; the rest follow as plain rows.
// `plain` drops the highlight when another card on the screen already holds it.
export function NextPractice({
  queue,
  label,
  plain = false,
  firstOnly = false,
}: {
  queue: QueueItem[]
  label: (patternId: string) => string
  plain?: boolean
  firstOnly?: boolean
}) {
  if (queue.length === 0) return null
  const [first, ...more] = queue
  const rest = firstOnly ? [] : more
  const href = (q: QueueItem) => `/drill?pattern=${encodeURIComponent(q.pattern_id)}`
  return (
    <div className="next-stack">
      <section className={plain ? 'panel next' : 'panel next hi'} aria-labelledby="next-h">
        <h2 id="next-h">Latih berikutnya</h2>
        <p className="next-name">{label(first.pattern_id)}</p>
        <p className={plain ? 'muted' : undefined}>{queueReason(first)}</p>
        <Link className={plain ? 'btn' : 'btn ink'} to={href(first)}>
          Latih sekarang
        </Link>
      </section>
      {rest.length > 0 && (
        <ul className="list" role="list" aria-label="Antrean latihan berikutnya">
          {rest.map((q) => (
            <li key={q.pattern_id} className="panel next-more">
              <span>
                <strong>{label(q.pattern_id)}</strong>
                <span className="muted"> {queueReason(q)}</span>
              </span>
              <Link className="btn ghost" to={href(q)}>
                Latih
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
