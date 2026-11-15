import type { CauseCounts, QueueItem } from './types'

// Whole-number percent, or a dash when nothing was answered.
export const fmtPct = (correct: number, answered: number) =>
  answered === 0 ? '-' : `${Math.round((correct * 100) / answered)}%`

// Milliseconds as "12 dtk" under a minute, otherwise "m:ss".
export const fmtDuration = (ms: number) => {
  const s = Math.round(ms / 1000)
  if (s < 60) return `${s} dtk`
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}

// The server only estimates causes from timing, so the wording says so.
export const CAUSE_LABELS: Record<keyof CauseCounts, string> = {
  tekanan_waktu: 'Lambat menjawab (mungkin tekanan waktu)',
  ceroboh: 'Terlalu cepat (mungkin ceroboh)',
  kurang_paham: 'Waktu wajar (mungkin kurang paham)',
  belum_jelas: 'Belum bisa diperkirakan (jawaban benar di pola ini kurang dari 3)',
}

export const causeText = (key: keyof CauseCounts, count: number) => `${CAUSE_LABELS[key]}: ${count} soal`

export const queueReason = (q: QueueItem) =>
  q.reason === 'akurasi_rendah' ? `Akurasi ${q.accuracy_pct}%, perlu diulang` : 'Belum dicoba'
