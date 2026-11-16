// Pure helpers for the exam simulation screen: timer, seq, screen reducer.
import type { Exam, ExamItem, ExamResponse, ExamResult, Spec } from './types'

// Whole seconds left until the deadline, never negative. offsetMs is server
// clock minus client clock, so the client's clock error does not move the timer.
export function remainingSeconds(deadlineISO: string, nowMs: number, offsetMs: number): number {
  const left = Date.parse(deadlineISO) - (nowMs + offsetMs)
  return Number.isNaN(left) ? 0 : Math.max(0, Math.ceil(left / 1000))
}

// Highest seq already stored for the exam; the next answer sends one more.
export function nextSeq(items: ExamItem[]): number {
  return items.reduce((m, it) => Math.max(m, it.response?.seq ?? 0), 0) + 1
}

export type Screen = 'login' | 'start' | 'running' | 'result'

export interface ExamSession {
  screen: Screen
  signedIn: boolean
  exam: Exam | null
  result: ExamResult | null
}

export type ExamAction =
  | { type: 'signedIn' }
  | { type: 'signedOut' }
  | { type: 'started'; exam: Exam }
  | { type: 'resumed'; exam: Exam; result?: ExamResult }
  | { type: 'answered'; pos: number; response: ExamResponse }
  | { type: 'finished'; result: ExamResult; exam?: Exam }
  | { type: 'reset' }

export const initialExamSession: ExamSession = { screen: 'login', signedIn: false, exam: null, result: null }

export function examReduce(s: ExamSession, a: ExamAction): ExamSession {
  switch (a.type) {
    case 'signedIn':
      // A sign-in event (token refresh, tab refocus) never pulls the user out of a running exam or its result.
      return s.screen === 'login' ? { ...s, signedIn: true, screen: 'start' } : { ...s, signedIn: true }
    case 'signedOut':
      return { ...initialExamSession }
    case 'started':
      return { ...s, screen: 'running', exam: a.exam, result: null }
    case 'resumed':
      return a.exam.status === 'in_progress'
        ? { ...s, screen: 'running', exam: a.exam, result: null }
        : { ...s, screen: 'result', exam: a.exam, result: a.result ?? null }
    case 'answered': {
      if (s.screen !== 'running' || !s.exam) return s
      const items = s.exam.items.map((it) =>
        it.pos === a.pos && (it.response?.seq ?? 0) <= a.response.seq ? { ...it, response: a.response } : it,
      )
      return { ...s, exam: { ...s.exam, items } }
    }
    case 'finished':
      if (s.screen === 'result') return s
      return { ...s, screen: 'result', result: a.result, exam: a.exam ?? s.exam }
    case 'reset':
      return s.signedIn ? { ...s, screen: 'start', exam: null, result: null } : { ...initialExamSession }
  }
}

// The weighted subtest: scored by option weight (the authors' judgement), not right or wrong.
export const WEIGHTED_SECTION = 'TKP'

// Ids of the sections of spec that can be taken now, in spec order.
export const availableSections = (sp: Spec): string[] => sp.sections.filter((s) => s.available).map((s) => s.id)

// Toggles id in selected, never leaving it empty (an exam needs a section).
export function toggleSection(selected: string[], id: string): string[] {
  if (!selected.includes(id)) return [...selected, id]
  return selected.length > 1 ? selected.filter((x) => x !== id) : selected
}

export type ItemStatus = 'empty' | 'answered' | 'flagged'

export const itemStatus = (it: ExamItem): ItemStatus =>
  !it.response?.answer ? 'empty' : it.response.flagged ? 'flagged' : 'answered'

export interface SectionSummary {
  id: string
  total: number
  answered: number
  flagged: number
}

// Per-section progress for the running screen, in first-seen section order.
export function summarizeSections(items: ExamItem[]): SectionSummary[] {
  const by = new Map<string, SectionSummary>()
  for (const it of items) {
    const s = by.get(it.section) ?? { id: it.section, total: 0, answered: 0, flagged: 0 }
    s.total++
    const st = itemStatus(it)
    if (st !== 'empty') s.answered++
    if (st === 'flagged') s.flagged++
    by.set(it.section, s)
  }
  return [...by.values()]
}

// Longest note the API accepts for a question report.
export const MAX_REPORT_NOTE = 500

// Whole seconds from milliseconds, at least 1 so a quick answer never reads "0 dtk".
export const formatSeconds = (ms: number): string => `${Math.max(1, Math.round(ms / 1000))} dtk`

// "Waktu: 42 dtk (median kamu 30 dtk)"; the median part only when there is one.
export function timeLine(timeMS: number, medianMS?: number): string {
  const base = `Waktu: ${formatSeconds(timeMS)}`
  return medianMS === undefined ? base : `${base} (median kamu ${formatSeconds(medianMS)})`
}
