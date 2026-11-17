import { CircleCheck, CircleX, Clock } from 'lucide-react'
import { Link } from 'react-router-dom'
import type { ExamSummary, Spec } from '../types'

const STATUS: Record<ExamSummary['status'], string> = {
  in_progress: 'Sedang berjalan',
  submitted: 'Selesai',
  expired: 'Waktu habis',
}

export const fmtDate = (iso: string) =>
  new Date(iso).toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })

// One exam in a list: date, spec and track, status, and each section's score
// against its threshold. The whole row links to the exam.
export function ExamRow({ e, specs }: { e: ExamSummary; specs: Spec[] }) {
  const spec = specs.find((s) => s.id === e.spec_id)
  const track = spec?.tracks.find((t) => t.id === e.track)
  return (
    <li>
      <Link className="row-link" to={`/simulation/${e.id}`}>
        <span className="exam-head">
          <span className="num">{fmtDate(e.started_at)}</span>
          <span>
            {spec?.label ?? e.spec_id} | {track?.label ?? e.track}
          </span>
          <span className={e.status === 'in_progress' ? 'chip warn' : 'muted'}>
            {e.status === 'in_progress' && <Clock className="ic" aria-hidden="true" />}
            {STATUS[e.status]}
          </span>
        </span>
        {e.result && (
          <span className="exam-scores">
            {e.result.sections.map((x) => (
              <span key={x.id} className="exam-score">
                <span className="num">
                  {x.id} {x.score}/{x.max}
                </span>
                {x.min !== undefined && (
                  <span className={x.met ? 'met' : 'unmet'}>
                    {x.met ? <CircleCheck className="ic" aria-hidden="true" /> : <CircleX className="ic" aria-hidden="true" />}
                    ambang {x.min}, {x.met ? 'tercapai' : 'belum tercapai'}
                  </span>
                )}
              </span>
            ))}
          </span>
        )}
      </Link>
    </li>
  )
}
