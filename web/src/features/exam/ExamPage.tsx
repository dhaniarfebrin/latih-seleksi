import { useEffect, useReducer, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ApiError } from '@/shared/api/client'
import { getExam, getExamResult } from './api'
import { examReduce, initialExamSession } from './exam'
import { Result } from './components/Result'
import { Running } from './components/Running'
import { errorMessage } from '@/shared/lib/error'

// The exam id comes from the URL, so a refresh or a shared link resumes the
// same exam. Keyed by id so switching exams starts from a clean state.
export function ExamRoute() {
  const { id = '' } = useParams()
  return <ExamPage key={id} id={id} />
}

function ExamPage({ id }: { id: string }) {
  const [s, dispatch] = useReducer(examReduce, { ...initialExamSession, signedIn: true, screen: 'start' })
  const [error, setError] = useState('')
  const [missing, setMissing] = useState(false)

  useEffect(() => {
    let live = true
    void (async () => {
      try {
        const exam = await getExam(id)
        const result = exam.status === 'in_progress' ? undefined : await getExamResult(id)
        if (live) dispatch({ type: 'resumed', exam, result })
      } catch (e) {
        if (!live) return
        if (e instanceof ApiError && e.status === 404) setMissing(true)
        else setError(errorMessage(e))
      }
    })()
    return () => {
      live = false
    }
  }, [id])

  if (missing) {
    return (
      <section className="screen narrow">
        <h1 className="d">Simulasi tidak ditemukan</h1>
        <p className="muted">Simulasi ini tidak ada atau bukan milik akun ini.</p>
        <Link className="btn hi" to="/simulation">
          Mulai simulasi baru
        </Link>
      </section>
    )
  }
  if (error) {
    return (
      <section className="screen narrow">
        <p role="alert" className="err">
          {error}
        </p>
        <Link className="btn" to="/history">
          Ke riwayat
        </Link>
      </section>
    )
  }
  if (s.screen === 'running' && s.exam) {
    return (
      <Running
        exam={s.exam}
        onAnswered={(pos, response) => dispatch({ type: 'answered', pos, response })}
        onFinished={(result, exam) => dispatch({ type: 'finished', result, exam })}
      />
    )
  }
  if (s.screen === 'result' && s.result) return <Result result={s.result} exam={s.exam} />
  return <p className="screen narrow muted">Memuat...</p>
}
