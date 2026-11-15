import { LogOut, MailCheck, MailWarning } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/features/auth/AuthContext'
import { authErrorMessage, signOut } from '@/features/auth/auth'
import { useState } from 'react'

export function AccountPage() {
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
        <p>
          {user?.emailVerified ? <MailCheck className="ic" aria-hidden="true" /> : <MailWarning className="ic" aria-hidden="true" />} Status
          email: {user?.emailVerified ? 'terverifikasi' : 'belum terverifikasi'}
        </p>
      </div>
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      <button className="btn" onClick={() => void out()}>
        <LogOut className="ic" aria-hidden="true" />&nbsp;Keluar
      </button>
    </section>
  )
}
