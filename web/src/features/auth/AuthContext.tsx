import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { authConfigured, onAuth, type AuthUser } from './auth'

interface AuthState {
  user: AuthUser | null
  // False until Firebase has reported the first auth state.
  ready: boolean
  setUser: (u: AuthUser | null) => void
}

const Ctx = createContext<AuthState>({ user: null, ready: false, setUser: () => {} })

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [ready, setReady] = useState(!authConfigured)
  useEffect(
    () =>
      onAuth((u) => {
        setUser(u)
        setReady(true)
      }),
    [],
  )
  const value = useMemo(() => ({ user, ready, setUser }), [user, ready])
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>
}

export const useAuth = () => useContext(Ctx)
