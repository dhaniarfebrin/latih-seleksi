import { CircleCheck, CircleX, Minus, Scale } from 'lucide-react'
import { Link } from 'react-router-dom'
import { ReportForm } from '@/shared/components/ReportForm'
import { BUMN_NOTICE, isBumnSpec, timeLine, WEIGHTED_SCORE_NOTE, WEIGHTED_SECTIONS, weightedSectionIds } from '../exam'
import { reportItem } from '../api'
import type { Exam, ExamItem, ExamResult } from '../types'

export function Result(props: { result: ExamResult; exam: Exam | null }) {
  const { result, exam } = props
  const weightedIds = weightedSectionIds(result.sections.map((x) => x.id))
  return (
    <div className="wide split">
      <div className="col">
        <h1 className="d">
          Skor <mark className="num">{result.total}</mark> / {result.max_total}
        </h1>
        {result.status === 'expired' && <p className="muted">Waktu habis.</p>}
        {isBumnSpec(result.spec_id) && <p className="chip warn">{BUMN_NOTICE}</p>}
        {result.partial && (
          <p className="chip warn">Simulasi parsial: tiap subtes dibandingkan dengan ambangnya, bukan status lulus.</p>
        )}
        {result.duration_basis === 'prorated' && (
          <p className="muted">Durasi dihitung prorata dari durasi resmi, jadi bukan durasi resmi ujian.</p>
        )}
        {weightedIds.length > 0 && (
          <p className="chip warn">
            Bobot {weightedIds.join(' dan ')} adalah penilaian penyusun, bukan kunci resmi.
          </p>
        )}
        {result.weighted !== undefined && (
          <p>
            Nilai tertimbang <mark className="num">{result.weighted.toFixed(1)}</mark> / 100
            <span className="muted"> ({WEIGHTED_SCORE_NOTE})</span>
          </p>
        )}
        <section className="panel">
          <h2 className="sub">Per subtes</h2>
          <ul className="sheet-rows" role="list">
            {result.sections.map((x) => (
              <li key={x.id} className="sheet-row">
                <span className="sheet-name">{x.id}</span>
                <span className="sheet-bar" aria-hidden="true">
                  <i style={{ width: `${x.max > 0 ? Math.round((x.score * 100) / x.max) : 0}%` }} />
                </span>
                <span className="num">
                  {x.score}/{x.max}
                </span>
                <div className="sheet-why">
                  <span>
                    {x.answered}/{x.items} dijawab
                    {x.min !== undefined && (
                      <>
                        {', ambang '}
                        {x.min}{' '}
                        <span className={x.met ? 'okt' : 'badt'}>
                          {x.met ? <CircleCheck className="ic" aria-hidden="true" /> : <CircleX className="ic" aria-hidden="true" />}{' '}
                          {x.met ? 'tercapai' : 'belum tercapai'}
                        </span>
                      </>
                    )}
                  </span>
                </div>
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
        const weighted = WEIGHTED_SECTIONS.has(it.section)
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
