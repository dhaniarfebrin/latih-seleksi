import { Link } from 'react-router-dom'
import { TrendBars } from '@/features/dashboard/TrendBars'
import { ExamRow, listExams, listSpecs } from '@/features/exam'
import type { ExamSummary } from '@/features/exam/types'
import { useLoad } from '@/shared/hooks/useLoad'

// Finished exams, oldest to newest, as percent of the maximum score of what was taken.
const examPoints = (exams: ExamSummary[]) =>
  exams
    .filter((e) => e.result && e.result.max_total > 0)
    .sort((a, b) => a.started_at.localeCompare(b.started_at))
    .slice(-10)
    .map((e) => ({
      value: Math.round((e.result!.total * 100) / e.result!.max_total),
      date: e.started_at,
      detail: `Simulasi ${new Date(e.started_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}: skor ${e.result!.total} dari ${e.result!.max_total}`,
    }))

export function HistoryPage() {
  const { data, error, loading } = useLoad(async () => {
    const [exams, specs] = await Promise.all([listExams(), listSpecs().catch(() => [])])
    return { exams, specs }
  })
  return (
    <section className="screen wide">
      <h1 className="d">Riwayat simulasi</h1>
      {loading && <p className="muted">Memuat...</p>}
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      {data && data.exams.length === 0 && (
        <div className="panel">
          <p className="muted">Belum ada simulasi.</p>
          <Link className="btn hi" to="/simulation">
            Mulai simulasi
          </Link>
        </div>
      )}
      {data && data.exams.length > 0 && (
        <>
          <section className="panel" aria-labelledby="skor-h">
            <h2 id="skor-h">Skor tiap simulasi (% dari nilai maksimum)</h2>
            <TrendBars points={examPoints(data.exams)} empty="Grafik muncul setelah 3 simulasi selesai." />
          </section>
          <ul className="list" role="list">
            {data.exams.map((e) => (
              <ExamRow key={e.id} e={e} specs={data.specs} />
            ))}
          </ul>
        </>
      )}
    </section>
  )
}
