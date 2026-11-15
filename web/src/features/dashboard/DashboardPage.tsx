import { Link } from 'react-router-dom'
import { listPatterns } from '@/features/drill'
import { useLoad } from '@/shared/hooks/useLoad'
import { getDashboard } from './api'
import { causeText, fmtDuration, fmtPct, queueReason } from './format'
import type { CauseCounts } from './types'

const CAUSE_KEYS: (keyof CauseCounts)[] = ['tekanan_waktu', 'ceroboh', 'kurang_paham', 'belum_jelas']

export function DashboardPage() {
  const { data, error, loading } = useLoad(async () => {
    const [dash, patterns] = await Promise.all([getDashboard(), listPatterns().catch(() => [])])
    return { dash, patterns }
  })
  const label = (id: string) => data?.patterns.find((p) => p.id === id)?.label ?? id
  const misLabel = (patternId: string, misId: string) =>
    data?.patterns.find((p) => p.id === patternId)?.misconceptions[misId] ?? misId

  return (
    <section className="screen wide">
      <p className="eyebrow">Dasbor</p>
      <h1 className="d">
        Titik <mark>lemah</mark>
      </h1>
      {loading && <p className="muted">Memuat...</p>}
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      {data && data.dash.patterns.length === 0 && (
        <div className="panel">
          <p className="muted">Belum ada jawaban. Kerjakan satu drill untuk melihat statistik.</p>
          <Link className="btn hi" to="/drill">
            Mulai drill
          </Link>
        </div>
      )}
      {data && data.dash.patterns.length > 0 && (
        <>
          <p className="muted">Berdasarkan 2000 jawaban terakhir dari drill dan simulasi (sesi tertua bisa terpotong). Penyebab salah hanyalah perkiraan dari waktu menjawab.</p>
          <h2 className="sub">Per subtes</h2>
          <table className="tbl">
            <thead>
              <tr>
                <th>Subtes</th>
                <th>Dijawab</th>
                <th>Akurasi</th>
                <th>Rata-rata waktu</th>
              </tr>
            </thead>
            <tbody>
              {data.dash.sections.map((s) => (
                <tr key={s.section}>
                  <td>{s.section}</td>
                  <td>{s.answered}</td>
                  <td>{fmtPct(s.correct, s.answered)}</td>
                  <td>{fmtDuration(s.avg_time_ms)}</td>
                </tr>
              ))}
            </tbody>
          </table>

          <h2 className="sub">Per pola soal</h2>
          <ul className="list" role="list">
            {data.dash.patterns.map((p) => (
              <li key={p.pattern_id} className="panel">
                <h3 className="sub">{label(p.pattern_id)}</h3>
                <p>
                  Akurasi {fmtPct(p.correct, p.answered)} ({p.correct}/{p.answered}), rata-rata{' '}
                  {fmtDuration(p.avg_time_ms)}
                </p>
                {p.answered > p.correct && (
                  <>
                    <p className="muted">Soal yang salah, dilihat dari waktu menjawab:</p>
                    <ul className="muted">
                      {CAUSE_KEYS.filter((k) => p.causes[k] > 0).map((k) => (
                        <li key={k}>{causeText(k, p.causes[k])}</li>
                      ))}
                    </ul>
                    {p.top_misconception_id && (
                      <p className="muted">Salah paling sering: {misLabel(p.pattern_id, p.top_misconception_id)}</p>
                    )}
                  </>
                )}
              </li>
            ))}
          </ul>

          <h2 className="sub">Sesi terakhir</h2>
          <ul className="list" role="list">
            {data.dash.sessions.map((s) => (
              <li key={s.attempt_id} className="panel">
                {s.mode === 'simulasi' ? 'Simulasi' : `Drill ${label(s.pattern_id)}`},{' '}
                {new Date(s.started_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}: benar{' '}
                {s.correct} dari {s.answered} ({fmtPct(s.correct, s.answered)})
              </li>
            ))}
          </ul>
        </>
      )}

      {data && data.dash.queue.length > 0 && (
        <>
          <h2 className="sub">Latihan berikutnya</h2>
          <ul className="list" role="list">
            {data.dash.queue.map((q) => (
              <li key={q.pattern_id} className="panel">
                <strong>{label(q.pattern_id)}</strong>
                <span className="muted">{queueReason(q)}</span>
                <Link className="btn" to={`/drill?pattern=${encodeURIComponent(q.pattern_id)}`}>
                  Latih sekarang
                </Link>
              </li>
            ))}
          </ul>
        </>
      )}
    </section>
  )
}
