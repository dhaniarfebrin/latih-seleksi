import { CircleCheck, CircleX, Minus } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { reportItem } from '../api'
import { misconceptionLabel, summarize } from '../session'
import type { Item } from '../types'
import { ReportForm } from '@/shared/components/ReportForm'
import { errorMessage } from '@/shared/lib/error'
import { Ruler } from '@/shared/components/Ruler'
import { fmtSeconds } from '@/shared/lib/time'

export function Summary(props: {
  attemptId: string
  items: Item[]
  targetSec: number | null
  labels: Record<string, string>
  onSubmit: () => Promise<unknown>
  onAgain: () => void
}) {
  const { attemptId, items, targetSec, labels, onSubmit, onAgain } = props
  const sum = summarize(items, targetSec)
  const [submitError, setSubmitError] = useState('')
  const sent = useRef(false)

  useEffect(() => {
    if (sent.current) return
    sent.current = true
    onSubmit().catch((e) => setSubmitError(errorMessage(e)))
  }, [onSubmit])

  const tally = Object.entries(sum.misconceptions).sort((a, b) => b[1] - a[1])
  return (
    <>
      <h1 className="d">
        Skor <mark className="num">{sum.score}</mark> / {sum.n}
      </h1>
      <section className="band">
        <Ruler
          label="Hasil tiap soal"
          current={0}
          cells={items.map((it) => ({
            pos: it.pos,
            state: it.response ? (it.response.correct ? 'ok' : 'bad') : 'empty',
            label: `Soal ${it.pos}, ${it.response ? (it.response.correct ? 'benar' : 'salah') : 'tidak dijawab'}`,
          }))}
        />
      </section>
      <p className="muted">
        Rata-rata <span className="num">{fmtSeconds(sum.avgMs)}</span> detik per soal.
      </p>
      {sum.withinTarget !== null && (
        <p className="muted">
          <span className="num">{sum.withinTarget}</span> dari {sum.n} soal dalam target {targetSec} detik.
        </p>
      )}
      {submitError && (
        <p role="alert" className="err">
          Hasil belum tersimpan di server: {submitError}
        </p>
      )}
      <section className="panel">
        <h2 className="sub">Jenis salah hitung</h2>
        {tally.length === 0 ? (
          <p className="muted">Tidak ada salah hitung yang terdeteksi.</p>
        ) : (
          <ul className="tally" role="list">
            {tally.map(([id, c]) => (
              <li key={id}>
                <span>{misconceptionLabel(labels, id)}</span>
                <span className="num">{c}x</span>
              </li>
            ))}
          </ul>
        )}
      </section>
      <section className="panel">
        <h2 className="sub">Tinjauan soal</h2>
        {items.map((it) => {
          const r = it.response
          return (
            <details key={it.pos}>
              <summary>
                <span className="num">{it.pos}.</span>{' '}
                <span className={r ? (r.correct ? 'okt' : 'badt') : 'muted'}>
                  {r ? (
                    r.correct ? (
                      <CircleCheck className="ic" aria-hidden="true" />
                    ) : (
                      <CircleX className="ic" aria-hidden="true" />
                    )
                  ) : (
                    <Minus className="ic" aria-hidden="true" />
                  )}{' '}
                  {r ? (r.correct ? 'Benar' : 'Salah') : 'Tidak dijawab'}
                </span>
              </summary>
              <p className="stem">{it.stem}</p>
              <p>
                Jawabanmu: <span className="num">{r?.answer || '-'}</span> | Kunci: <span className="num">{r?.answer_key ?? '-'}</span>
              </p>
              {r && (
                <p className="muted">
                  Waktu <span className="num">{fmtSeconds(r.time_ms)}</span> detik
                  {r.misconception_id && <> | {misconceptionLabel(labels, r.misconception_id)}</>}
                </p>
              )}
              {r && r.steps.length > 0 && (
                <ol className="steps">
                  {r.steps.map((st, i) => (
                    <li key={i}>{st}</li>
                  ))}
                </ol>
              )}
              {r && <ReportForm onSend={(reason, note) => reportItem(attemptId, it.pos, reason, note)} />}
            </details>
          )
        })}
      </section>
      <button className="btn hi" onClick={onAgain}>
        Ulangi
      </button>
    </>
  )
}
