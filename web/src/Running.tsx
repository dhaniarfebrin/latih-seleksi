import { useCallback, useEffect, useRef, useState } from 'react'
import {
  ApiError,
  getExam,
  getExamResult,
  getServerOffset,
  putExamResponse,
  submitExam,
  type Exam,
  type ExamItem,
  type ExamResult,
} from './api'
import { getToken } from './auth'
import { fmtClock, itemStatus, nextSeq, remainingSeconds, summarizeSections } from './exam'

const STATUS_LABEL = { empty: 'belum dijawab', answered: 'terjawab', flagged: 'ragu-ragu' }

export function Running(props: {
  exam: Exam
  onAnswered: (pos: number, r: NonNullable<ExamItem['response']>) => void
  onFinished: (r: ExamResult, exam?: Exam) => void
}) {
  const { exam, onAnswered, onFinished } = props
  const [idx, setIdx] = useState(0)
  const [left, setLeft] = useState(() => remainingSeconds(exam.deadline_at, Date.now(), getServerOffset()))
  const [error, setError] = useState('')
  const [confirming, setConfirming] = useState(false)
  const [busy, setBusy] = useState(false)
  const submitted = useRef(false)
  const seq = useRef(nextSeq(exam.items))
  const spent = useRef(new Map<number, number>()) // ms already spent per pos
  const shownAt = useRef(performance.now())
  const item = exam.items[idx]

  const finish = useCallback(async () => {
    if (submitted.current) return
    submitted.current = true
    setBusy(true)
    setError('')
    try {
      const token = await getToken()
      let result: ExamResult
      try {
        result = await submitExam(token, exam.id)
      } catch (e) {
        // 409: already closed (deadline passed or submitted elsewhere); just read the result.
        if (e instanceof ApiError && e.status === 409) result = await getExamResult(token, exam.id)
        else throw e
      }
      onFinished(result, await getExam(token, exam.id))
    } catch (e) {
      submitted.current = false
      setBusy(false)
      setError((e as Error).message)
    }
  }, [exam.id, onFinished])

  useEffect(() => {
    const t = setInterval(() => setLeft(remainingSeconds(exam.deadline_at, Date.now(), getServerOffset())), 500)
    return () => clearInterval(t)
  }, [exam.deadline_at])

  useEffect(() => {
    if (left === 0) void finish()
  }, [left, finish])

  const goto = (i: number) => {
    const now = performance.now()
    spent.current.set(item.pos, (spent.current.get(item.pos) ?? 0) + now - shownAt.current)
    shownAt.current = now
    setIdx(i)
  }

  const send = async (answer: string, flagged: boolean) => {
    setError('')
    const ms = Math.round((spent.current.get(item.pos) ?? 0) + performance.now() - shownAt.current)
    try {
      const resp = await putExamResponse(await getToken(), exam.id, item.pos, {
        answer,
        time_ms: ms,
        flagged,
        seq: seq.current++,
      })
      onAnswered(item.pos, resp)
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) void finish()
      else setError((e as Error).message)
    }
  }

  const r = item.response
  const sections = summarizeSections(exam.items)
  const unanswered = exam.items.filter((it) => itemStatus(it) === 'empty').length
  const flagged = sections.reduce((n, x) => n + x.flagged, 0)
  const timer = (
    <span className={`num timer${left <= 60 ? ' low' : ''}`} role="timer" aria-label="Sisa waktu">
      {fmtClock(left)}
    </span>
  )
  const grid = (
    <div className="grid" role="group" aria-label="Nomor soal">
      {exam.items.map((it, i) => (
        <button
          key={it.pos}
          className={`cell ${itemStatus(it)}${i === idx ? ' cur' : ''}`}
          aria-label={`Soal ${it.pos}, ${STATUS_LABEL[itemStatus(it)]}`}
          aria-current={i === idx}
          onClick={() => goto(i)}
        >
          {it.pos}
        </button>
      ))}
    </div>
  )

  return (
    <div className="wide exam">
      <div className="bar mo">
        <p className="eyebrow">
          Soal {item.pos} / {exam.items.length}
        </p>
        {timer}
      </div>
      <div className="qmain">
        <p className="eyebrow dk">
          Soal {item.pos} / {exam.items.length}
        </p>
        <details className="mo">
          <summary>
            Nomor soal ({exam.items.length - unanswered}/{exam.items.length} terjawab)
          </summary>
          {grid}
        </details>
        <p className="stem">{item.stem}</p>
        <ul className="opts" role="list">
          {item.options.map((o) => (
            <li key={o.key}>
              <button
                className={`opt${r?.answer === o.key ? ' picked' : ''}`}
                aria-pressed={r?.answer === o.key}
                disabled={busy}
                onClick={() => void send(o.key, r?.flagged ?? false)}
              >
                <span className="k num">{o.key}</span>
                <span>{o.text}</span>
              </button>
            </li>
          ))}
        </ul>
        {error && (
          <p role="alert" className="err">
            {error}
          </p>
        )}
        <div className="row">
          <button className="btn" disabled={idx === 0} onClick={() => goto(idx - 1)}>
            Sebelumnya
          </button>
          <button className="btn" disabled={!r || busy} onClick={() => r && void send(r.answer, !r.flagged)}>
            {r?.flagged ? 'Batal ragu' : 'Tandai ragu'}
          </button>
          <button className="btn" disabled={idx === exam.items.length - 1} onClick={() => goto(idx + 1)}>
            Berikutnya
          </button>
        </div>
      </div>
      <aside className="side">
        <div className="dk side-timer">
          <span className="eyebrow">Sisa waktu</span>
          {timer}
        </div>
        <p className="muted">
          {sections.map((x) => `${x.id} ${x.answered}/${x.total}`).join(' | ')}
          {flagged > 0 ? ` | ragu ${flagged}` : ''}
        </p>
        <div className="dk">{grid}</div>
        {confirming ? (
          <section className="panel" role="alertdialog" aria-label="Konfirmasi selesai">
            <p>{unanswered > 0 ? `${unanswered} soal belum dijawab. Akhiri simulasi sekarang?` : 'Akhiri simulasi sekarang?'}</p>
            <div className="row">
              <button className="btn hi" disabled={busy} onClick={() => void finish()}>
                {busy ? 'Mengirim...' : 'Ya, selesai'}
              </button>
              <button className="btn" disabled={busy} onClick={() => setConfirming(false)}>
                Lanjut mengerjakan
              </button>
            </div>
          </section>
        ) : (
          <button className="btn hi" onClick={() => setConfirming(true)}>
            Selesai
          </button>
        )}
      </aside>
    </div>
  )
}
