// Thin client for the Go API. Field names mirror api/internal/drill/drill.go.
const BASE = import.meta.env.VITE_API_URL ?? ''

export interface Pattern {
  id: string
  label: string
}

export interface Option {
  key: string
  text: string
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

async function call<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(BASE + path, {
    method,
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const j = (await res.json()) as { error?: string }
      if (j.error) msg = j.error
    } catch {
      // keep the status message
    }
    throw new Error(msg)
  }
  return (await res.json()) as T
}

export const listPatterns = () => call<{ patterns: Pattern[] }>('GET', '/v1/patterns').then((r) => r.patterns)

export const createAttempt = (pattern_id: string, n: number, band: number) =>
  call<Attempt>('POST', '/v1/attempts', { pattern_id, n, band })

export const getAttempt = (id: string) => call<Attempt>('GET', `/v1/attempts/${id}`)

// pos is 1-based; the first answer for a position wins and the response
// reveals the key, steps and misconception_id.
export const putResponse = (id: string, pos: number, answer: string, time_ms: number) =>
  call<ResponseView>('PUT', `/v1/attempts/${id}/responses/${pos}`, { answer, time_ms })

export const submitAttempt = (id: string) => call<Result>('POST', `/v1/attempts/${id}/submit`)
