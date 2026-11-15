import { useEffect, useRef, useState } from 'react'
import { misconceptionLabel, summarize } from '../session'
import type { Item } from '../types'
import { errorMessage } from '@/shared/lib/error'

export function Summary(props: {
  items: Item[]
  labels: Record<string, string>
  onSubmit: () => Promise<unknown>
  onAgain: () => void
}) {
  const { items, labels, onSubmit, onAgain } = props
  const sum = summarize(items)
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
      <p className="eyebrow">Ringkasan</p>
      <h1 className="d">
        Skor <mark className="num">{sum.score}</mark> / {sum.n}
      </h1>
      <p className="muted">
        Rata-rata <span className="num">{(sum.avgMs / 1000).toFixed(1)}</span> detik per soal.
      </p>
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
      <button className="btn hi" onClick={onAgain}>
        Ulangi
      </button>
    </>
  )
}
