// Exam simulation shapes. Mirror api/internal/exam/exam.go.
import type { Option } from '@/shared/api/types'

export interface SpecSection {
  id: string
  label: string
  items: number
  available: boolean
}

export interface Track {
  id: string
  label: string
  thresholds?: Record<string, number>
  total_min?: number
}

export interface Spec {
  id: string
  label: string
  scope: string
  duration_min: number
  sections: SpecSection[]
  tracks: Track[]
}

export interface ExamResponse {
  answer: string
  time_ms: number
  flagged: boolean
  seq: number
  correct?: boolean
  // Weight of the chosen option; only for a closed exam and a weighted (TKP) item.
  weight?: number
  misconception_id?: string
}

export interface ExamItem {
  pos: number
  section: string
  stem: string
  options: Option[]
  response?: ExamResponse
  answer_key?: string
  steps?: string[]
  // Median time of the user's correct answers in the same pattern; closed exam only.
  median_ms?: number
}

export interface Exam {
  id: string
  spec_id: string
  track: string
  sections: string[]
  status: 'in_progress' | 'submitted' | 'expired'
  started_at: string
  deadline_at: string
  duration_seconds: number
  duration_basis: string
  submitted_at?: string
  items: ExamItem[]
}

export interface SectionResult {
  id: string
  items: number
  answered: number
  score: number
  max: number
  min?: number
  met: boolean
}

export interface ExamResult {
  exam_id: string
  spec_id: string
  track: string
  status: string
  partial: boolean
  duration_basis: string
  sections: SectionResult[]
  total: number
  max_total: number
  submitted_at: string
}

// An exam without items; result is set once it is closed.
export interface ExamSummary {
  id: string
  spec_id: string
  track: string
  sections: string[]
  status: Exam['status']
  started_at: string
  deadline_at: string
  duration_seconds: number
  duration_basis: string
  submitted_at?: string
  result?: ExamResult
}

export type ReportReason = 'kunci_salah' | 'soal_ambigu' | 'typo' | 'lainnya'
