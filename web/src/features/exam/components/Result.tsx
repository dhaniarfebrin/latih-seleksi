import { CircleCheck, CircleX, Minus, Scale } from 'lucide-react'
import { Link } from 'react-router-dom'
import { ReportForm } from '@/shared/components/ReportForm'
import { timeLine, WEIGHTED_SECTION } from '../exam'
import { reportItem } from '../api'
import type { Exam, ExamItem, ExamResult } from '../types'

export function Result(props: { result: ExamResult; exam: Exam | null }) {
  const { result, exam } = props
  return (
    <div className="wide split">
      <div className="col">
        <p className="eyebrow">{result.status === 'expired' ? 'Waktu habis' : 'Selesai'}</p>
        <h1 className="d">
          Skor <mark className="num">{result.total}</mark> / {result.max_total}
        </h1>
        {result.partial && (
          <p className="chip warn">Simulasi parsial: tiap subtes dibandingkan dengan ambangnya, bukan status lulus.</p>
        )}
        {result.duration_basis === 'prorated' && (
          <p className="muted">Durasi dihitung prorata dari durasi resmi, jadi bukan durasi resmi ujian.</p>
        )}
        {result.sections.some((x) => x.id === WEIGHTED_SECTION) && (
          <p className="chip warn">Bobot TKP adalah penilaian penyusun, bukan kunci resmi.</p>
        )}
        <section className="panel">
          <h2 className="sub">Per subtes</h2>
          <ul className="tally" role="list">
            {result.sections.map((x) => (
              <li key={x.id}>
                <span>
                  {x.id}: {x.answered}/{x.items} dijawab
                </span>
                <span className="num">
                  {x.score}/{x.max}
                  {x.min !== undefined && (
                    <>
                      {' '}
                      (ambang {x.min}){' '}
                      {x.met ? <CircleCheck className="ic" aria-hidden="true" /> : <CircleX className="ic" aria-hidden="true" />}{' '}
                      {x.met ? 'tercapai' : 'belum tercapai'}
                    </>
                  )}
                </span>
              </li>
            ))}
          </ul>
        </section>
        <Link className="btn hi" to="/simulation">
          Simulasi lagi
        </Link>
        <Link className="btn" to="/history">
          Riwayat
        </Link>
      </div>
      <div className="col">{exam && <Review examId={exam.id} items={exam.items} />}</div>
    </div>
  )
}

function Review({ examId, items }: { examId: string; items: ExamItem[] }) {
  return (
    <section className="panel">
      <h2 className="sub">Tinjauan soal</h2>
      {items.map((it) => {
        const r = it.response
        const weighted = it.section === WEIGHTED_SECTION
        return (
          <details key={it.pos}>
            <summary>
              <span className="num">{it.pos}.</span>{' '}
              {weighted && r ? (
                <span>
                  <Scale className="ic" aria-hidden="true" /> Bobot {r.weight ?? '-'} dari 5
                </span>
              ) : (
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
              )}
            </summary>
            <p className="stem">{it.stem}</p>
            <p>
              Jawabanmu: <span className="num">{r?.answer || '-'}</span> | {weighted ? 'Bobot tertinggi' : 'Kunci'}: <span className="num">{it.answer_key ?? '-'}</span>
            </p>
            {r && <p className="muted">{timeLine(r.time_ms, it.median_ms)}</p>}
            {it.steps && it.steps.length > 0 && (
              <ol className="steps">
                {it.steps.map((st, i) => (
                  <li key={i}>{st}</li>
                ))}
              </ol>
            )}
            <ReportForm onSend={(reason, note) => reportItem(examId, it.pos, reason, note)} />
          </details>
        )
      })}
    </section>
  )
}
