import { LogIn } from 'lucide-react'
import { useState, type FormEvent } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from './AuthContext'
import {
  authConfigured,
  authErrorMessage,
  refreshUser,
  resendVerification,
  signInEmail,
  signInGoogle,
  signOut,
  signUpEmail,
} from './auth'

export function Login() {
  const { user, ready } = useAuth()
  const loc = useLocation()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  if (!authConfigured) {
    return (
      <section className="screen narrow">
        <h1 className="d">Login belum dikonfigurasi</h1>
        <p className="muted">Isi VITE_FIREBASE_* lalu bangun ulang aplikasi.</p>
      </section>
    )
  }
  if (!ready) return <p className="screen narrow muted">Memuat...</p>
  if (user?.emailVerified) {
    const from = (loc.state as { from?: string } | null)?.from
    return <Navigate to={from && from !== '/login' ? from : '/'} replace />
  }
  if (user) return <Verify />

  const run = async (fn: () => Promise<void>) => {
    setBusy(true)
    setError('')
    try {
      await fn()
    } catch (e) {
      setError(authErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  const submit = (e: FormEvent, signUp: boolean) => {
    e.preventDefault()
    void run(() => (signUp ? signUpEmail(email, password) : signInEmail(email, password)))
  }

  return (
    <section className="screen narrow">
      <p className="eyebrow">Masuk</p>
      <h1 className="d">
        Masuk ke <mark>Latih Seleksi</mark>
      </h1>
      <p className="muted">Hasil simulasi disimpan per akun.</p>
      <button className="btn hi" disabled={busy} onClick={() => void run(signInGoogle)}>
        Masuk dengan Google
      </button>
      <form className="field" onSubmit={(e) => submit(e, false)}>
        <label className="lbl">
          Email
          <input type="email" autoComplete="email" required value={email} onChange={(e) => setEmail(e.target.value)} />
        </label>
        <label className="lbl">
          Kata sandi
          <input
            type="password"
            autoComplete="current-password"
            required
            minLength={6}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
        <div className="row">
          <button className="btn" type="submit" disabled={busy}>
            <LogIn className="ic" aria-hidden="true" />&nbsp;Masuk
          </button>
          <button className="btn" type="button" disabled={busy || !email || password.length < 6} onClick={(e) => submit(e, true)}>
            Daftar
          </button>
        </div>
      </form>
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
    </section>
  )
}

export function Verify() {
  const { user, setUser } = useAuth()
  const [msg, setMsg] = useState('')
  const [busy, setBusy] = useState(false)
  const act = async (fn: () => Promise<string>) => {
    setBusy(true)
    try {
      setMsg(await fn())
    } catch (e) {
      setMsg(authErrorMessage(e))
    } finally {
      setBusy(false)
    }
  }
  return (
    <section className="screen narrow">
      <p className="eyebrow">Akun</p>
      <h1 className="d">Verifikasi email</h1>
      <p className="muted">Kami mengirim tautan verifikasi ke {user?.email}. Buka tautan itu, lalu muat ulang status.</p>
      <button
        className="btn hi"
        disabled={busy}
        onClick={() =>
          void act(async () => {
            const u = await refreshUser()
            setUser(u)
            return u?.emailVerified ? '' : 'Email belum terverifikasi'
          })
        }
      >
        Saya sudah verifikasi
      </button>
      <button
        className="btn"
        disabled={busy}
        onClick={() =>
          void act(async () => {
            await resendVerification()
            return 'Email verifikasi dikirim ulang'
          })
        }
      >
        Kirim ulang email
      </button>
      <button className="link" onClick={() => void signOut()}>
        Keluar
      </button>
      {msg && (
        <p role="status" className="muted">
          {msg}
        </p>
      )}
    </section>
  )
}
