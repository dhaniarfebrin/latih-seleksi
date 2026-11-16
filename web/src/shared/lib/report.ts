export type ReportReason = 'kunci_salah' | 'soal_ambigu' | 'typo' | 'lainnya'

export const REASONS: { id: ReportReason; label: string }[] = [
  { id: 'kunci_salah', label: 'Kunci jawaban salah' },
  { id: 'soal_ambigu', label: 'Soal ambigu' },
  { id: 'typo', label: 'Salah ketik' },
  { id: 'lainnya', label: 'Lainnya' },
]

// Longest note the API accepts for a question report.
export const MAX_REPORT_NOTE = 500
