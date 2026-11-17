import { Link } from 'react-router-dom'
import { listPatterns } from '@/features/drill'
import { useLoad } from '@/shared/hooks/useLoad'
import { getDashboard } from './api'
import { causeText, fmtDuration, fmtPct } from './format'
import { NextPractice } from './NextPractice'
import { TrendBars } from './TrendBars'
import type { CauseCounts, PatternStat, SessionStat } from './types'

const CAUSE_KEYS: (keyof CauseCounts)[] = ['tekanan_waktu', 'ceroboh', 'kurang_paham', 'belum_jelas']

const pct = (correct: number, answered: number) => Math.round((correct * 100) / answered)

// Weakest first; patterns with more answers win a tie because their number is firmer.
const weakestFirst = (a: PatternStat, b: PatternStat) =>
  a.correct / a.answered - b.correct / b.answered || b.answered - a.answered

// The API returns sessions newest first; the trend reads oldest to newest.
const sessionPoints = (sessions: SessionStat[], label: (id: string) => string) =>
  sessions
    .filter((s) => s.answered > 0)
    .reverse()
    .map((s) => ({
      value: pct(s.correct, s.answered),
      date: s.started_at,
      detail: `${s.mode === 'simulasi' ? 'Simulasi' : `Drill ${label(s.pattern_id)}`}, ${new Date(s.started_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}: benar ${s.correct} dari ${s.answered} (${fmtPct(s.correct, s.answered)})`,
    }))

export function DashboardPage() {
  const { data, error, loading } = useLoad(async () => {
    const [dash, patterns] = await Promise.all([getDashboard(), listPatterns().catch(() => [])])
    return { dash, patterns }
  })
  const label = (id: string) => data?.patterns.find((p) => p.id === id)?.label ?? id
  const misLabel = (patternId: string, misId: string) =>
    data?.patterns.find((p) => p.id === patternId)?.misconceptions[misId] ?? misId

  const points = data ? sessionPoints([...data.dash.sessions], label) : []
  const latest = points.at(-1)
  const earlier = points.slice(0, -1)
  const earlierAvg = earlier.length > 0 ? Math.round(earlier.reduce((n, p) => n + p.value, 0) / earlier.length) : 0
  const ranked = data ? data.dash.patterns.filter((p) => p.answered > 0).sort(weakestFirst) : []

  return (
    <section className="screen wide">
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
      {data && data.dash.patterns.length === 0 && <NextPractice queue={data.dash.queue} label={label} plain />}
      {data && data.dash.patterns.length > 0 && (
        <>
          <p className="muted">
            Berdasarkan 2000 jawaban terakhir dari drill dan simulasi (sesi tertua bisa terpotong). Penyebab salah
            hanyalah perkiraan dari waktu menjawab.
          </p>
          <div className="dash-top">
            <section className="panel" aria-labelledby="trend-h">
              <h2 id="trend-h">Sesi terakhir, benar dari dijawab (%)</h2>
              <TrendBars points={points} empty="Tren muncul setelah 3 sesi. Kerjakan satu drill lagi." />
              {latest && earlier.length >= 2 && (
                <p className="muted">
                  Sesi terbaru {latest.value}%, rata-rata {earlier.length} sesi sebelumnya {earlierAvg}%.
                </p>
              )}
              <details className="trend-detail">
                <summary>Rincian sesi</summary>
                <ul className="list" role="list">
                  {data.dash.sessions.map((s) => (
                    <li key={s.attempt_id}>
                      {s.mode === 'simulasi' ? 'Simulasi' : `Drill ${label(s.pattern_id)}`},{' '}
                      {new Date(s.started_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}: benar{' '}
                      {s.correct} dari {s.answered} ({fmtPct(s.correct, s.answered)})
                    </li>
                  ))}
                </ul>
              </details>
            </section>
            <NextPractice queue={data.dash.queue} label={label} />
          </div>

          <section className="panel" aria-labelledby="pola-h">
            <h2 id="pola-h">Per pola soal, terlemah dulu</h2>
            <ul className="sheet-rows" role="list">
              {ranked.map((p, i) => {
                const weak = i === 0 && p.answered > p.correct
                return (
                  <li key={p.pattern_id} className={weak ? 'sheet-row weak' : 'sheet-row'}>
                    <span className="sheet-name">{label(p.pattern_id)}</span>
                    <span className="sheet-bar" aria-hidden="true">
                      <i style={{ width: `${pct(p.correct, p.answered)}%`, animationDelay: `${Math.min(i, 8) * 70}ms` }} />
                    </span>
                    <span className="num">{fmtPct(p.correct, p.answered)}</span>
                    <div className="sheet-why">
                      <span>
                        {p.correct}/{p.answered} benar, rata-rata {fmtDuration(p.avg_time_ms)}
                      </span>
                      {p.answered > p.correct && (
                        <details open={weak}>
                          <summary>Penyebab salah</summary>
                          <ul>
                            {CAUSE_KEYS.filter((k) => p.causes[k] > 0).map((k) => (
                              <li key={k}>{causeText(k, p.causes[k])}</li>
                            ))}
                          </ul>
                          {p.top_misconception_id && (
                            <p>Salah paling sering: {misLabel(p.pattern_id, p.top_misconception_id)}</p>
                          )}
                        </details>
                      )}
                    </div>
                  </li>
                )
              })}
            </ul>
          </section>

          <section className="panel" aria-labelledby="subtes-h">
            <h2 id="subtes-h">Per subtes</h2>
            <div className="tbl-wrap">
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
            </div>
          </section>
        </>
      )}
    </section>
  )
}
