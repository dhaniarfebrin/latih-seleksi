import { Link } from 'react-router-dom'
import { ExamRow, listExams, listSpecs } from '@/features/exam'
import { useLoad } from '@/shared/hooks/useLoad'

export function HistoryPage() {
  const { data, error, loading } = useLoad(async () => {
    const [exams, specs] = await Promise.all([listExams(), listSpecs().catch(() => [])])
    return { exams, specs }
  })
  return (
    <section className="screen wide">
      <p className="eyebrow">Riwayat</p>
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
        <ul className="list" role="list">
          {data.exams.map((e) => (
            <ExamRow key={e.id} e={e} specs={data.specs} />
          ))}
        </ul>
      )}
    </section>
  )
}
