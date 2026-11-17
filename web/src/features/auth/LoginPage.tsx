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

export function LoginPage() {
  const { user, ready } = useAuth()
  const loc = useLocation()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [mode, setMode] = useState<'in' | 'up'>('in')

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
  const submit = (e: FormEvent) => {
    e.preventDefault()
    void run(() => (mode === 'up' ? signUpEmail(email, password) : signInEmail(email, password)))
  }

  const label = mode === 'up' ? 'Daftar' : 'Masuk'
  return (
    <section className="screen narrow login-wrap">
      <div className="login-note">
        <h2 className="sub">Kenapa perlu akun?</h2>
        <ul>
          <li>Riwayat drill dan simulasi tersimpan per akun.</li>
          <li>Titik lemah dan jadwal ulangan dihitung dari jawaban Anda.</li>
          <li>Data bisa diunduh atau dihapus kapan saja.</li>
        </ul>
      </div>
      <div className="col">
        <h1 className="d">
          {label} ke <mark>Latih Seleksi</mark>
        </h1>
        <p className="muted">Hasil simulasi disimpan per akun.</p>
        <fieldset className="field">
          <legend>Pilihan</legend>
          <div className="seg" style={{ gridTemplateColumns: '1fr 1fr' }}>
            {(['in', 'up'] as const).map((m) => (
              <label key={m} className={mode === m ? 'on' : ''}>
                <input type="radio" name="mode" checked={mode === m} onChange={() => setMode(m)} />
                {m === 'in' ? 'Masuk' : 'Daftar'}
              </label>
            ))}
          </div>
        </fieldset>
        <form className="field" onSubmit={submit}>
          <label className="lbl">
            Email
            <input type="email" autoComplete="email" required value={email} onChange={(e) => setEmail(e.target.value)} />
          </label>
          <label className="lbl">
            Kata sandi
            <input
              type="password"
              autoComplete={mode === 'up' ? 'new-password' : 'current-password'}
              required
              minLength={6}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>
          <button className="btn hi" type="submit" disabled={busy}>
            <LogIn className="ic" aria-hidden="true" />&nbsp;{label}
          </button>
        </form>
        <div className="or">atau</div>
        <button className="btn ghost" disabled={busy} onClick={() => void run(signInGoogle)}>
          Masuk dengan Google
        </button>
        {error && (
          <p role="alert" className="err">
            {error}
          </p>
        )}
      </div>
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
        className="btn ghost"
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
