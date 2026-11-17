import { useAuth } from '@/features/auth/AuthContext'
import { HomePage } from '@/features/home/HomePage'
import { Verify } from '@/features/auth/LoginPage'
import { LandingContent } from './LandingContent'

// "/" shows the landing to guests (and while auth is loading), Beranda to verified users.
export function LandingPage() {
  const { user, ready } = useAuth()
  if (ready && user) return user.emailVerified ? <HomePage /> : <Verify />
  return <LandingContent />
}
