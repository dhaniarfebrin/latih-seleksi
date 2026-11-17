import { Download, LogOut, MailCheck, MailWarning, Trash2 } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/features/auth/AuthContext'
import { authErrorMessage, deleteAccount, signOut } from '@/features/auth/auth'
import { errorMessage } from '@/shared/lib/error'
import { useState } from 'react'
import { deleteMyData, exportMyData } from './api'

export function AccountPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [downloading, setDownloading] = useState(false)
  const out = async () => {
    try {
      await signOut()
      navigate('/login', { replace: true })
    } catch (e) {
      setError(authErrorMessage(e))
    }
  }
  const download = async () => {
    setError('')
    setDownloading(true)
    try {
      const data = await exportMyData()
      const url = URL.createObjectURL(new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' }))
      const a = document.createElement('a')
      a.href = url
      a.download = 'latih-seleksi-data.json'
      a.click()
      URL.revokeObjectURL(url)
    } catch (e) {
      setError(errorMessage(e) || 'Gagal mengunduh data')
    } finally {
      setDownloading(false)
    }
  }
  const remove = async () => {
    if (!window.confirm('Hapus akun dan semua data latihan? Tindakan ini tidak bisa dibatalkan.')) return
    setError('')
    try {
      await deleteMyData()
    } catch (e) {
      setError(errorMessage(e) || 'Gagal menghapus data')
      return
    }
    try {
      await deleteAccount()
      navigate('/login', { replace: true })
    } catch (e) {
      const code = (e as { code?: string }).code
      setError(
        code === 'auth/requires-recent-login'
          ? authErrorMessage(e)
          : 'Data latihan sudah dihapus, tetapi akun login belum. Coba tekan Hapus akun dan data sekali lagi.',
      )
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
      <button className="btn" disabled={downloading} onClick={() => void download()}>
        <Download className="ic" aria-hidden="true" />&nbsp;Unduh data saya
      </button>
      <button className="btn" onClick={() => void remove()}>
        <Trash2 className="ic" aria-hidden="true" />&nbsp;Hapus akun dan data
      </button>
    </section>
  )
}
