import { useEffect, useReducer, useRef, useState } from 'react'
import { createAttempt, listPatterns, putResponse, submitAttempt, type Item, type Pattern } from './api'
import { initialSession, misconceptionLabel, reduce, summarize } from './session'

const N = 8

export function App() {
  const [s, dispatch] = useReducer(reduce, initialSession)
  const [patterns, setPatterns] = useState<Pattern[]>([])
  const [patternsError, setPatternsError] = useState('')

  useEffect(() => {
    listPatterns()
      .then(setPatterns)
      .catch((e: Error) => setPatternsError(e.message))
  }, [])

  const labelsOf = (patternId: string) => patterns.find((p) => p.id === patternId)?.misconceptions ?? {}
  return (
    <div className="app">
      <main className="screen">
        {s.screen === 'start' && (
          <Start
            patterns={patterns}
            loadError={patternsError}
            onStarted={(attempt) => dispatch({ type: 'started', attempt })}
          />
        )}
        {s.screen === 'question' && s.attempt && (
          <Question
            key={s.index}
            attemptId={s.attempt.id}
            item={s.attempt.items[s.index]}
            total={s.attempt.items.length}
            labels={labelsOf(s.attempt.pattern_id)}
            onAnswered={(pos, response) => dispatch({ type: 'answered', pos, response })}
            onNext={() => dispatch({ type: 'next' })}
          />
        )}
        {s.screen === 'summary' && s.attempt && (
          <Summary
            attemptId={s.attempt.id}
            items={s.attempt.items}
            labels={labelsOf(s.attempt.pattern_id)}
            onAgain={() => dispatch({ type: 'reset' })}
          />
        )}
      </main>
      <Footer />
    </div>
  )
}

function Start(props: {
  patterns: Pattern[]
  loadError: string
  onStarted: (a: Awaited<ReturnType<typeof createAttempt>>) => void
}) {
  const { patterns, loadError, onStarted } = props
  const [patternId, setPatternId] = useState('')
  const [band, setBand] = useState(2)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  const pattern = patterns.find((p) => p.id === patternId) ?? patterns[0]
  const start = async () => {
    if (!pattern) return
    setBusy(true)
    setError('')
    try {
      onStarted(await createAttempt(pattern.id, N, band))
    } catch (e) {
      setError((e as Error).message)
      setBusy(false)
    }
  }

  return (
    <>
      <p className="eyebrow">Drill {N} soal</p>
      <h1 className="d">
        Latih <mark>Seleksi</mark>
      </h1>
      <p className="muted">{pattern ? pattern.label : 'Memuat pola soal...'} (TIU). Jawab, lalu lihat langkah dan jenis salah hitungnya.</p>
      <fieldset className="field">
        <legend>Pola soal</legend>
        <div className="seg">
          {patterns.map((p) => (
            <label key={p.id} className={pattern?.id === p.id ? 'on' : ''}>
              <input
                type="radio"
                name="pattern"
                checked={pattern?.id === p.id}
                onChange={() => setPatternId(p.id)}
              />
              {p.label}
            </label>
          ))}
        </div>
      </fieldset>
      <fieldset className="field">
        <legend>Tingkat kesulitan</legend>
        <div className="seg">
          {[
            [1, 'Mudah'],
            [2, 'Sedang'],
            [3, 'Sulit'],
          ].map(([v, label]) => (
            <label key={v} className={band === v ? 'on' : ''}>
              <input type="radio" name="band" checked={band === v} onChange={() => setBand(v as number)} />
              {label}
            </label>
          ))}
        </div>
      </fieldset>
      {(loadError || error) && (
        <p role="alert" className="err">
          {error || loadError}
        </p>
      )}
      <button className="btn hi" disabled={!pattern || busy} onClick={start}>
        {busy ? 'Menyiapkan...' : 'Mulai latihan'}
      </button>
    </>
  )
}

function Question(props: {
  attemptId: string
  item: Item
  total: number
  labels: Record<string, string>
  onAnswered: (pos: number, r: NonNullable<Item['response']>) => void
  onNext: () => void
}) {
  const { attemptId, item, total, labels, onAnswered, onNext } = props
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
      const resp = await putResponse(attemptId, item.pos, key, Math.round(performance.now() - startedAt.current))
      onAnswered(item.pos, resp)
    } catch (e) {
      setError((e as Error).message)
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
                {r && isKey && <span className="tag">Kunci</span>}
                {r && chosen && !isKey && <span className="tag">Pilihanmu</span>}
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
          <p className={`verdict ${r.correct ? 'ok' : 'bad'}`}>{r.correct ? 'Benar' : `Salah. Kunci: ${r.answer_key}`}</p>
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

function Summary(props: { attemptId: string; items: Item[]; labels: Record<string, string>; onAgain: () => void }) {
  const { attemptId, items, labels, onAgain } = props
  const sum = summarize(items)
  const [submitError, setSubmitError] = useState('')
  const sent = useRef(false)

  useEffect(() => {
    if (sent.current) return
    sent.current = true
    submitAttempt(attemptId).catch((e: Error) => setSubmitError(e.message))
  }, [attemptId])

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

function Footer() {
  return (
    <footer className="foot">
      Aplikasi latihan independen, bukan produk resmi penyelenggara seleksi. Soal buatan sendiri dan tidak menjamin kelulusan.{' '}
      <a href="https://cat.bkn.go.id" target="_blank" rel="noopener noreferrer">
        Simulasi resmi (cat.bkn.go.id)
      </a>
    </footer>
  )
}

function fmtClock(sec: number) {
  const t = Math.floor(sec)
  return `${Math.floor(t / 60)}:${String(t % 60).padStart(2, '0')}`
}
