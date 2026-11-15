import { useNavigate } from 'react-router-dom'
import { useAuth } from './AuthContext'
import { authErrorMessage, signOut } from './auth'
import { useState } from 'react'

export function Account() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const out = async () => {
    try {
      await signOut()
      navigate('/login', { replace: true })
    } catch (e) {
      setError(authErrorMessage(e))
    }
  }
  return (
    <section className="screen narrow">
      <p className="eyebrow">Akun</p>
      <h1 className="d">Akun</h1>
      <div className="panel">
        <p>
          Email: <strong>{user?.email ?? '-'}</strong>
        </p>
        <p>Status email: {user?.emailVerified ? 'terverifikasi' : 'belum terverifikasi'}</p>
      </div>
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      <button className="btn" onClick={() => void out()}>
        Keluar
      </button>
    </section>
  )
}
