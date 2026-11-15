// Thin client for the Go API. Field names mirror api/internal/drill/drill.go.
const BASE = import.meta.env.VITE_API_URL ?? ''

export interface Pattern {
  id: string
  label: string
  misconceptions: Record<string, string>
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

export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message)
  }
}

// Server clock minus client clock in ms, from the Date header of the latest
// response (1 s resolution). The exam timer uses it instead of the raw client clock.
let serverOffsetMs = 0
export const getServerOffset = () => serverOffsetMs

async function call<T>(method: string, path: string, body?: unknown, token?: string | null): Promise<T> {
  const headers: Record<string, string> = {}
  if (body !== undefined) headers['Content-Type'] = 'application/json'
  if (token) headers.Authorization = `Bearer ${token}`
  const res = await fetch(BASE + path, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  const date = Date.parse(res.headers.get('Date') ?? '')
  if (!Number.isNaN(date)) serverOffsetMs = date - Date.now()
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const j = (await res.json()) as { error?: string }
      if (j.error) msg = j.error
    } catch {
      // keep the status message
    }
    if (res.status === 401) msg = 'Sesi login berakhir, silakan login lagi'
    else if (res.status === 403 && msg === 'email not verified') msg = 'Email belum diverifikasi. Cek kotak masuk lalu muat ulang status.'
    throw new ApiError(res.status, msg)
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

// Exam simulation (auth required). Shapes mirror api/internal/exam/exam.go.
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

export const listSpecs = () => call<{ specs: Spec[] }>('GET', '/v1/specs').then((r) => r.specs)

export const createExam = (token: string | null, spec_id: string, track: string, sections: string[]) =>
  call<Exam>('POST', '/v1/exams', { spec_id, track, sections }, token)

export const getExam = (token: string | null, id: string) => call<Exam>('GET', `/v1/exams/${id}`, undefined, token)

// seq must rise per exam: the server ignores a request whose seq is not newer.
export const putExamResponse = (
  token: string | null,
  id: string,
  pos: number,
  r: { answer: string; time_ms: number; flagged: boolean; seq: number },
) => call<ExamResponse>('PUT', `/v1/exams/${id}/responses/${pos}`, r, token)

export const submitExam = (token: string | null, id: string) =>
  call<ExamResult>('POST', `/v1/exams/${id}/submit`, undefined, token)

export const getExamResult = (token: string | null, id: string) =>
  call<ExamResult>('GET', `/v1/exams/${id}/result`, undefined, token)

// Newest first; the server caps limit at 50.
export const listExams = (token: string | null, limit?: number) =>
  call<{ exams: ExamSummary[] }>('GET', `/v1/exams${limit ? `?limit=${limit}` : ''}`, undefined, token).then((r) => r.exams)
