import { Link } from 'react-router-dom'
import { listExams, listSpecs } from './api'
import { getToken } from './auth'
import { ExamRow } from './ExamRow'
import { useLoad } from './useLoad'

export function History() {
  const { data, error, loading } = useLoad(async () => {
    const [exams, specs] = await Promise.all([getToken().then((t) => listExams(t)), listSpecs().catch(() => [])])
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
