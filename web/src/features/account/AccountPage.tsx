import { Download, LogOut, MailCheck, MailWarning, Trash2 } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/features/auth/AuthContext'
import { authErrorMessage, deleteAccount, signOut } from '@/features/auth/auth'
import { errorMessage } from '@/shared/lib/error'
import { useState } from 'react'
import { ConfirmDialog } from '@/shared/components/ConfirmDialog'
import { deleteMyData, exportMyData } from './api'

export function AccountPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [downloading, setDownloading] = useState(false)
  const [notice, setNotice] = useState('')
  const [confirming, setConfirming] = useState(false)
  const [deleting, setDeleting] = useState(false)
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
    setNotice('')
    setDownloading(true)
    try {
      const data = await exportMyData()
      const url = URL.createObjectURL(new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' }))
      const a = document.createElement('a')
      a.href = url
      a.download = 'latih-seleksi-data.json'
      a.click()
      URL.revokeObjectURL(url)
      setNotice('Data berhasil diunduh.')
    } catch (e) {
      setError(errorMessage(e) || 'Gagal mengunduh data')
    } finally {
      setDownloading(false)
    }
  }
  const remove = async () => {
    setError('')
    setNotice('')
    setDeleting(true)
    try {
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
    } finally {
      setDeleting(false)
      setConfirming(false)
    }
  }
  const verified = !!user?.emailVerified
  const initial = (user?.email ?? '?').charAt(0).toUpperCase()
  return (
    <section className="screen narrow">
      <h1 className="d">Akun</h1>
      <div className="panel">
        <div className="acct-head">
          <span className="avatar" aria-hidden="true">
            {initial}
          </span>
          <div>
            <div className="em">{user?.email ?? '-'}</div>
            <span className={verified ? 'badge ok' : 'badge warn'}>
              {verified ? <MailCheck className="ic" aria-hidden="true" /> : <MailWarning className="ic" aria-hidden="true" />}
              {verified ? 'Email terverifikasi' : 'Email belum terverifikasi'}
            </span>
          </div>
        </div>
      </div>
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      {notice && (
        <p role="status" className="okt">
          {notice}
        </p>
      )}
      <div className="panel">
        <h2 className="sub">Data saya</h2>
        <p className="muted">Unduh riwayat latihan dan simulasi Anda dalam berkas JSON.</p>
        <button className="btn ghost" disabled={downloading} onClick={() => void download()}>
          <Download className="ic" aria-hidden="true" />&nbsp;Unduh data saya
        </button>
      </div>
      <div className="panel">
        <h2 className="sub">Sesi</h2>
        <button className="btn ghost" onClick={() => void out()}>
          <LogOut className="ic" aria-hidden="true" />&nbsp;Keluar
        </button>
      </div>
      <div className="panel danger">
        <h2 className="sub">Zona berbahaya</h2>
        <p className="muted">Menghapus akun juga menghapus semua data latihan dan simulasi.</p>
        <button
          className="btn danger"
          onClick={() => {
            setError('')
            setNotice('')
            setConfirming(true)
          }}
        >
          <Trash2 className="ic" aria-hidden="true" />&nbsp;Hapus akun dan data
        </button>
      </div>
      <ConfirmDialog
        open={confirming}
        title="Hapus akun dan data?"
        confirmLabel={deleting ? 'Menghapus...' : 'Ya, hapus akun'}
        busy={deleting}
        onConfirm={() => void remove()}
        onCancel={() => setConfirming(false)}
      >
        <p>Akun login dan semua riwayat drill, simulasi, dan jadwal ulangan Anda akan dihapus. Tindakan ini tidak bisa dibatalkan.</p>
        <p className="muted">Ingin menyimpan salinan? Tutup dialog ini dan tekan Unduh data saya dulu.</p>
      </ConfirmDialog>
    </section>
  )
}
