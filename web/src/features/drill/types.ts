// Field names mirror api/internal/drill/drill.go.
import type { Option } from '@/shared/api/types'

export interface Pattern {
  id: string
  label: string
  misconceptions: Record<string, string>
}

export interface ResponseView {
  answer: string
  correct: boolean
  time_ms: number
  answer_key: string
  steps: string[]
  misconception_id?: string
  answered_at: string
}

export interface Item {
  pos: number
  stem: string
  options: Option[]
  response?: ResponseView
}

export interface Attempt {
  id: string
  mode: string
  pattern_id: string
  band: number
  n: number
  started_at: string
  submitted_at?: string
  items: Item[]
}

export interface Result {
  attempt_id: string
  n: number
  answered: number
  score: number
  avg_time_ms: number
  misconceptions: Record<string, number>
  submitted_at: string
}
