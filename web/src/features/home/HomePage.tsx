import { ChartNoAxesColumn, Play, Target, Timer } from 'lucide-react'
import { Link } from 'react-router-dom'
import { ExamRow, listExams, listSpecs } from '@/features/exam'
import { useAuth } from '@/features/auth/AuthContext'
import { useLoad } from '@/shared/hooks/useLoad'

export function HomePage() {
  const { user } = useAuth()
  const { data, error, loading } = useLoad(async () => {
    const [exams, specs] = await Promise.all([listExams(5), listSpecs().catch(() => [])])
    return { exams, specs }
  })
  const running = data?.exams[0]?.status === 'in_progress' ? data.exams[0] : undefined
  const lastDone = data?.exams.find((e) => e.status !== 'in_progress')
  return (
    <section className="screen wide">
      <p className="eyebrow">Beranda</p>
      <h1 className="d">
        Halo, <mark>{user?.email?.split('@')[0] ?? 'peserta'}</mark>
      </h1>
      {running && (
        <div className="panel">
          <h2 className="sub">Simulasi sedang berjalan</h2>
          <p className="muted">Waktu terus berjalan di server. Lanjutkan sebelum habis.</p>
          <Link className="btn hi" to={`/simulation/${running.id}`}>
            <Play className="ic" aria-hidden="true" />&nbsp;Lanjutkan simulasi
          </Link>
        </div>
      )}
      <div className="panel choices">
        <div>
          <h2 className="sub">
            <Target className="ic" aria-hidden="true" /> Drill
          </h2>
          <p className="muted">Latihan 8 soal TIU dengan langkah dan jenis salah hitung.</p>
          <Link className="btn" to="/drill">
            Mulai drill
          </Link>
        </div>
        <div>
          <h2 className="sub">
            <Timer className="ic" aria-hidden="true" /> Simulasi
          </h2>
          <p className="muted">Simulasi TIU berwaktu, dinilai server terhadap ambang.</p>
          <Link className="btn" to="/simulation">
            Mulai simulasi
          </Link>
        </div>
      </div>
      <div className="panel">
        <h2 className="sub">
          <ChartNoAxesColumn className="ic" aria-hidden="true" /> Titik lemah
        </h2>
        <p className="muted">Akurasi dan waktu per pola soal, serta saran latihan berikutnya.</p>
        <Link className="btn" to="/dashboard">
          Lihat dasbor
        </Link>
      </div>
      <section>
        <h2 className="sub">Hasil terakhir</h2>
        {loading && <p className="muted">Memuat...</p>}
        {error && (
          <p role="alert" className="err">
            {error}
          </p>
        )}
        {data && !lastDone && <p className="muted">Belum ada simulasi yang selesai.</p>}
        {data && lastDone && (
          <ul className="list" role="list">
            <ExamRow e={lastDone} specs={data.specs} />
          </ul>
        )}
        {data && data.exams.length > 0 && (
          <p>
            <Link to="/history">Lihat semua riwayat</Link>
          </p>
        )}
      </section>
    </section>
  )
}
