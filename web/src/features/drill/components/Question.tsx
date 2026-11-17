import { CircleCheck, CircleX } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { fmtClock } from '@/shared/lib/time'
import { misconceptionLabel } from '../session'
import type { Item } from '../types'
import { errorMessage } from '@/shared/lib/error'

export function Question(props: {
  item: Item
  total: number
  targetSec: number | null
  labels: Record<string, string>
  onAnswer: (pos: number, key: string, timeMs: number) => Promise<void>
  onNext: () => void
}) {
  const { item, total, targetSec, labels, onAnswer, onNext } = props
  const startedAt = useRef(performance.now())
  const [elapsed, setElapsed] = useState(0)
  const [sending, setSending] = useState(false)
  const [error, setError] = useState('')
  const r = item.response
  const nextRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    if (r) return
    const t = setInterval(() => setElapsed(Math.floor((performance.now() - startedAt.current) / 1000)), 500)
    return () => clearInterval(t)
  }, [r])

  useEffect(() => {
    if (r) nextRef.current?.focus()
  }, [r])

  const answer = async (key: string) => {
    if (r || sending) return
    setSending(true)
    setError('')
    try {
      await onAnswer(item.pos, key, Math.round(performance.now() - startedAt.current))
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setSending(false)
    }
  }

  return (
    <>
      <div className="qhead">
        <p className="eyebrow">
          Soal {item.pos} / {total}
        </p>
        <span className="num timer" aria-label="Waktu berjalan">
          {fmtClock(r ? r.time_ms / 1000 : elapsed)}
        </span>
        {targetSec !== null && (
          <span className={`num timer${(r ? r.time_ms / 1000 : elapsed) > targetSec ? ' low' : ''}`}>
            Target {targetSec} d{(r ? r.time_ms / 1000 : elapsed) > targetSec ? ' (lewat)' : ''}
          </span>
        )}
      </div>
      <progress max={total} value={item.pos - (r ? 0 : 1)} aria-label="Kemajuan" />
      <p className="stem">{item.stem}</p>
      <ul className="opts" role="list">
        {item.options.map((o) => {
          const chosen = r?.answer === o.key
          const isKey = r?.answer_key === o.key
          const cls = r ? (isKey ? 'ok' : chosen ? 'bad' : '') : ''
          return (
            <li key={o.key}>
              <button
                className={`opt ${cls}`}
                disabled={!!r || sending}
                aria-pressed={chosen}
                onClick={() => answer(o.key)}
              >
                <span className="k num">{o.key}</span>
                <span>{o.text}</span>
                {r && isKey && (
                  <span className="tag">
                    <CircleCheck className="ic" aria-hidden="true" />
                    Kunci
                  </span>
                )}
                {r && chosen && !isKey && (
                  <span className="tag">
                    <CircleX className="ic" aria-hidden="true" />
                    Pilihanmu
                  </span>
                )}
              </button>
            </li>
          )
        })}
      </ul>
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      {r && (
        <section className="panel" aria-live="polite">
          <p className={`verdict ${r.correct ? 'ok' : 'bad'}`}>
            {r.correct ? <CircleCheck className="ic" aria-hidden="true" /> : <CircleX className="ic" aria-hidden="true" />}
            {r.correct ? 'Benar' : `Salah. Kunci: ${r.answer_key}`}
          </p>
          {r.misconception_id && <p className="chip warn">{misconceptionLabel(labels, r.misconception_id)}</p>}
          <h2 className="sub">Langkah</h2>
          <ol className="steps">
            {r.steps.map((st, i) => (
              <li key={i}>{st}</li>
            ))}
          </ol>
          <button ref={nextRef} className="btn" onClick={onNext}>
            {item.pos < total ? 'Lanjut' : 'Lihat ringkasan'}
          </button>
        </section>
      )}
    </>
  )
}
