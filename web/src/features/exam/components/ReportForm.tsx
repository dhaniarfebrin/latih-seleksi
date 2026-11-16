import { CircleAlert, CircleCheck } from 'lucide-react'
import { useState } from 'react'
import { reportItem } from '../api'
import { MAX_REPORT_NOTE } from '../exam'
import type { ReportReason } from '../types'

const REASONS: { id: ReportReason; label: string }[] = [
  { id: 'kunci_salah', label: 'Kunci jawaban salah' },
  { id: 'soal_ambigu', label: 'Soal ambigu' },
  { id: 'typo', label: 'Salah ketik' },
  { id: 'lainnya', label: 'Lainnya' },
]

type Status = { kind: 'idle' } | { kind: 'sending' } | { kind: 'sent' } | { kind: 'failed'; message: string }

export function ReportForm({ examId, pos }: { examId: string; pos: number }) {
  const [reason, setReason] = useState<ReportReason>('kunci_salah')
  const [note, setNote] = useState('')
  const [status, setStatus] = useState<Status>({ kind: 'idle' })

  const send = async () => {
    setStatus({ kind: 'sending' })
    try {
      await reportItem(examId, pos, reason, note.trim())
      setStatus({ kind: 'sent' })
    } catch (e) {
      setStatus({ kind: 'failed', message: e instanceof Error ? e.message : 'Gagal mengirim laporan' })
    }
  }

  return (
    <details>
      <summary>Laporkan soal</summary>
      <div className="field">
        <label className="lbl">
          Alasan
          <select value={reason} onChange={(e) => setReason(e.target.value as ReportReason)}>
            {REASONS.map((r) => (
              <option key={r.id} value={r.id}>
                {r.label}
              </option>
            ))}
          </select>
        </label>
        <label className="lbl">
          Catatan (opsional)
          <textarea rows={3} maxLength={MAX_REPORT_NOTE} value={note} onChange={(e) => setNote(e.target.value)} />
        </label>
        <button type="button" className="btn" disabled={status.kind === 'sending'} onClick={send}>
          {status.kind === 'sending' ? 'Mengirim...' : 'Kirim laporan'}
        </button>
        <p role="status" className={status.kind === 'failed' ? 'badt' : 'okt'}>
          {status.kind === 'sent' && (
            <>
              <CircleCheck className="ic" aria-hidden="true" /> Laporan terkirim
            </>
          )}
          {status.kind === 'failed' && (
            <>
              <CircleAlert className="ic" aria-hidden="true" /> {status.message}
            </>
          )}
        </p>
      </div>
    </details>
  )
}
