// Firebase Authentication wrapper. The SDK is initialised lazily from the
// VITE_FIREBASE_* values (public web config); without them auth is off and
// the drill keeps working. Tokens are never stored or logged here.
import { initializeApp, type FirebaseApp } from 'firebase/app'
import {
  createUserWithEmailAndPassword,
  getAuth,
  GoogleAuthProvider,
  onAuthStateChanged,
  sendEmailVerification,
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut as fbSignOut,
  type Auth,
  type User,
} from 'firebase/auth'

export interface AuthUser {
  uid: string
  email: string | null
  emailVerified: boolean
}

const cfg = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY as string | undefined,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN as string | undefined,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID as string | undefined,
}

export const authConfigured = !!(cfg.apiKey && cfg.authDomain && cfg.projectId)

let app: FirebaseApp | undefined
function auth(): Auth {
  if (!authConfigured) throw new Error('Login belum dikonfigurasi')
  app ??= initializeApp({ apiKey: cfg.apiKey, authDomain: cfg.authDomain, projectId: cfg.projectId })
  return getAuth(app)
}

const view = (u: User | null): AuthUser | null =>
  u ? { uid: u.uid, email: u.email, emailVerified: u.emailVerified } : null

export const signInGoogle = async () => {
  await signInWithPopup(auth(), new GoogleAuthProvider())
}

export const signInEmail = async (email: string, password: string) => {
  await signInWithEmailAndPassword(auth(), email, password)
}

export const signUpEmail = async (email: string, password: string) => {
  const cred = await createUserWithEmailAndPassword(auth(), email, password)
  await sendEmailVerification(cred.user)
}

export const resendVerification = async () => {
  const u = auth().currentUser
  if (!u) throw new Error('Belum login')
  await sendEmailVerification(u)
}

// Reloads the profile so a just-verified email shows up without re-login.
export const refreshUser = async (): Promise<AuthUser | null> => {
  const u = auth().currentUser
  if (u) {
    await u.reload()
    await u.getIdToken(true)
  }
  return view(auth().currentUser)
}

export const signOut = () => fbSignOut(auth())

// Calls cb with null when signed out. Without config it reports null once.
export function onAuth(cb: (u: AuthUser | null) => void): () => void {
  if (!authConfigured) {
    cb(null)
    return () => {}
  }
  return onAuthStateChanged(auth(), (u) => cb(view(u)))
}

// getIdToken refreshes an expired token by itself.
export const getToken = async (): Promise<string | null> => {
  if (!authConfigured) return null
  const u = auth().currentUser
  return u ? u.getIdToken() : null
}

// Maps Firebase error codes to short Indonesian messages.
export function authErrorMessage(e: unknown): string {
  const code = (e as { code?: string }).code ?? ''
  switch (code) {
    case 'auth/invalid-credential':
    case 'auth/wrong-password':
    case 'auth/user-not-found':
      return 'Email atau kata sandi salah'
    case 'auth/email-already-in-use':
      return 'Email sudah terdaftar'
    case 'auth/weak-password':
      return 'Kata sandi terlalu lemah (minimal 6 karakter)'
    case 'auth/invalid-email':
      return 'Format email tidak valid'
    case 'auth/popup-closed-by-user':
    case 'auth/cancelled-popup-request':
      return 'Login dibatalkan'
    case 'auth/too-many-requests':
      return 'Terlalu banyak percobaan, coba lagi nanti'
    case 'auth/network-request-failed':
      return 'Koneksi gagal'
    default:
      return (e as Error).message || 'Login gagal'
  }
}
