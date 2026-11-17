import { ChevronLeft, ChevronRight, Clock, Flag } from 'lucide-react'
import { useCallback, useEffect, useRef, useState } from 'react'
import { ApiError, getServerOffset } from '@/shared/api/client'
import { getExam, getExamResult, putExamResponse, submitExam } from '../api'
import type { Exam, ExamItem, ExamResult } from '../types'
import { fmtClock } from '@/shared/lib/time'
import { itemStatus, nextSeq, remainingSeconds, summarizeSections } from '../exam'
import { errorMessage } from '@/shared/lib/error'
import { Ruler } from '@/shared/components/Ruler'

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
      let result: ExamResult
      try {
        result = await submitExam(exam.id)
      } catch (e) {
        // 409: already closed (deadline passed or submitted elsewhere); just read the result.
        if (e instanceof ApiError && e.status === 409) result = await getExamResult(exam.id)
        else throw e
      }
      onFinished(result, await getExam(exam.id))
    } catch (e) {
      submitted.current = false
      setBusy(false)
      setError(errorMessage(e))
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
      const resp = await putExamResponse(exam.id, item.pos, {
        answer,
        time_ms: ms,
        flagged,
        seq: seq.current++,
      })
      onAnswered(item.pos, resp)
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) void finish()
      else setError(errorMessage(e))
    }
  }

  const r = item.response
  const sections = summarizeSections(exam.items)
  const unanswered = exam.items.filter((it) => itemStatus(it) === 'empty').length
  const flagged = sections.reduce((n, x) => n + x.flagged, 0)
  const timer = (
    <span className={`num timer${left <= 60 ? ' low' : ''}`} role="timer" aria-label="Sisa waktu">
      <Clock className="ic" aria-hidden="true" /> {fmtClock(left)}
    </span>
  )
  // The ruler covers the subtest on screen; the tabs move between subtests.
  const ruler = (
    <Ruler
      label="Nomor soal"
      current={item.pos}
      cells={exam.items
        .filter((it) => it.section === item.section)
        .map((it) => ({ pos: it.pos, state: itemStatus(it), label: `Soal ${it.pos}, ${STATUS_LABEL[itemStatus(it)]}` }))}
      onPick={(pos) => goto(exam.items.findIndex((it) => it.pos === pos))}
    />
  )
  const tabs = (
    <div className="tabs" role="group" aria-label="Subtes">
      {sections.map((x) => (
        <button
          key={x.id}
          className={`tab${x.id === item.section ? ' on' : ''}`}
          aria-pressed={x.id === item.section}
          onClick={() => goto(exam.items.findIndex((it) => it.section === x.id))}
        >
          {x.id} {x.answered}/{x.total}
        </button>
      ))}
    </div>
  )

  return (
    <div className="wide exam">
      <header className="pin">
        <div className="band">
          <div className="band-top">
            {timer}
            <div className="dk tabs-slot">{tabs}</div>
            {!confirming && (
              <button className="btn hi band-end" onClick={() => setConfirming(true)}>
                Selesai
              </button>
            )}
          </div>
          <progress max={exam.items.length} value={exam.items.length - unanswered} aria-label="Soal terjawab" />
          <div className="dk">{ruler}</div>
          <details className="mo">
            <summary>
              Nomor soal ({exam.items.length - unanswered}/{exam.items.length} terjawab
              {flagged > 0 ? `, ragu ${flagged}` : ''})
            </summary>
            {tabs}
            {ruler}
          </details>
        </div>
      </header>
      {confirming && (
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
      )}
      <section className="qcard" aria-label={`Soal ${item.pos}`}>
        <p className="muted">
          Soal {item.pos} dari {exam.items.length}, {item.section}
        </p>
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
            <ChevronLeft className="ic" aria-hidden="true" />
            Sebelumnya
          </button>
          <button className="btn" disabled={!r || busy} onClick={() => r && void send(r.answer, !r.flagged)}>
            <Flag className="ic" aria-hidden="true" />&nbsp;{r?.flagged ? 'Batal ragu' : 'Tandai ragu'}
          </button>
          <button className="btn" disabled={idx === exam.items.length - 1} onClick={() => goto(idx + 1)}>
            Berikutnya
            <ChevronRight className="ic" aria-hidden="true" />
          </button>
        </div>
      </section>
    </div>
  )
}
