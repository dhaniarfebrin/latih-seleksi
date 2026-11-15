import { History as HistoryIcon, House, Target, Timer, UserRound } from 'lucide-react'
import { Link, Navigate, NavLink, Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { AboutPage } from '@/features/about/AboutPage'
import { AccountPage } from '@/features/account/AccountPage'
import { AuthProvider, useAuth } from '@/features/auth/AuthContext'
import { authConfigured } from '@/features/auth/auth'
import { DrillPage } from '@/features/drill/DrillPage'
import { ExamRoute } from '@/features/exam/ExamPage'
import { HistoryPage } from '@/features/history/HistoryPage'
import { HomePage } from '@/features/home/HomePage'
import { LoginPage, Verify } from '@/features/auth/LoginPage'
import { StartPage } from '@/features/exam/StartPage'

export function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route element={<Shell />}>
          <Route path="login" element={<LoginPage />} />
          <Route path="about" element={<AboutPage />} />
          <Route element={<RequireAuth />}>
            <Route index element={<HomePage />} />
            <Route path="drill" element={<DrillPage />} />
            <Route path="simulation" element={<StartPage />} />
            <Route path="simulation/:id" element={<ExamRoute />} />
            <Route path="history" element={<HistoryPage />} />
            <Route path="account" element={<AccountPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}

// The login gate is a UI convenience; the API enforces the token on drill and exams.
function RequireAuth() {
  const { user, ready } = useAuth()
  const loc = useLocation()
  if (!authConfigured) {
    return (
      <section className="screen narrow">
        <h1 className="d">Login belum dikonfigurasi</h1>
        <p className="muted">Isi VITE_FIREBASE_* lalu bangun ulang aplikasi.</p>
      </section>
    )
  }
  if (!ready) return <p className="screen narrow muted">Memuat...</p>
  if (!user) return <Navigate to="/login" replace state={{ from: loc.pathname + loc.search }} />
  if (!user.emailVerified) return <Verify />
  return <Outlet />
}

function Shell() {
  const { user } = useAuth()
  const signedIn = !!user?.emailVerified
  return (
    <div className="app">
      <header className="top">
        <div className="in">
          <Link to={signedIn ? '/' : '/login'} className="brand">
            Latih <mark>Seleksi</mark>
          </Link>
          {signedIn && (
            <nav className="nav" aria-label="Menu utama">
              <NavLink to="/" end>
                <House className="ic" aria-hidden="true" />
                Beranda
              </NavLink>
              <NavLink to="/drill">
                <Target className="ic" aria-hidden="true" />
                Drill
              </NavLink>
              <NavLink to="/simulation">
                <Timer className="ic" aria-hidden="true" />
                Simulasi
              </NavLink>
              <NavLink to="/history">
                <HistoryIcon className="ic" aria-hidden="true" />
                Riwayat
              </NavLink>
              <NavLink to="/account">
                <UserRound className="ic" aria-hidden="true" />
                Akun
              </NavLink>
            </nav>
          )}
        </div>
      </header>
      <main className="main">
        <Outlet />
      </main>
      <footer className="foot">
        <div className="in">
          Aplikasi latihan independen, bukan produk resmi penyelenggara seleksi. Soal buatan sendiri dan tidak menjamin kelulusan.{' '}
          <Link to="/about">Tentang</Link>{' '}
          <a href="https://cat.bkn.go.id" target="_blank" rel="noopener noreferrer">
            Simulasi resmi (cat.bkn.go.id)
          </a>
        </div>
      </footer>
    </div>
  )
}
