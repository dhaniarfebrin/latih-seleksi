import { Link, Navigate, NavLink, Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { About } from './About'
import { Account } from './Account'
import { AuthProvider, useAuth } from './AuthContext'
import { authConfigured } from './auth'
import { Drill } from './Drill'
import { ExamRoute } from './ExamPage'
import { History } from './History'
import { Home } from './Home'
import { Login, Verify } from './Login'
import { Start } from './Start'

export function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route element={<Shell />}>
          <Route path="login" element={<Login />} />
          <Route path="about" element={<About />} />
          <Route element={<RequireAuth />}>
            <Route index element={<Home />} />
            <Route path="drill" element={<Drill />} />
            <Route path="simulation" element={<Start />} />
            <Route path="simulation/:id" element={<ExamRoute />} />
            <Route path="history" element={<History />} />
            <Route path="account" element={<Account />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}

// The login gate is a UI convenience only: the drill API stays public.
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
                Home
              </NavLink>
              <NavLink to="/drill">Drill</NavLink>
              <NavLink to="/simulation">Simulation</NavLink>
              <NavLink to="/history">History</NavLink>
            </nav>
          )}
          {signedIn && (
            <NavLink to="/account" className="acct">
              Akun
            </NavLink>
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
