import type { Attempt, Item, ResponseView } from './types'

export type Screen = 'start' | 'question' | 'summary'

export interface Session {
  screen: Screen
  attempt: Attempt | null
  index: number // 0-based position of the current question
}

export type Action =
  | { type: 'started'; attempt: Attempt }
  | { type: 'answered'; pos: number; response: ResponseView }
  | { type: 'next' }
  | { type: 'reset' }

export const initialSession: Session = { screen: 'start', attempt: null, index: 0 }

export function reduce(s: Session, a: Action): Session {
  switch (a.type) {
    case 'started':
      return { screen: 'question', attempt: a.attempt, index: 0 }
    case 'answered': {
      if (!s.attempt) return s
      const items = s.attempt.items.map((it) =>
        it.pos === a.pos && !it.response ? { ...it, response: a.response } : it,
      )
      return { ...s, attempt: { ...s.attempt, items } }
    }
    case 'next': {
      if (s.screen !== 'question' || !s.attempt) return s
      return s.index + 1 < s.attempt.items.length
        ? { ...s, index: s.index + 1 }
        : { ...s, screen: 'summary' }
    }
    case 'reset':
      return initialSession
  }
}

export interface Summary {
  n: number
  answered: number
  score: number
  avgMs: number
  misconceptions: Record<string, number>
}

export function summarize(items: Item[]): Summary {
  let answered = 0
  let score = 0
  let totalMs = 0
  const misconceptions: Record<string, number> = {}
  for (const it of items) {
    const r = it.response
    if (!r) continue
    answered++
    totalMs += r.time_ms
    if (r.correct) score++
    if (r.misconception_id) misconceptions[r.misconception_id] = (misconceptions[r.misconception_id] ?? 0) + 1
  }
  return { n: items.length, answered, score, avgMs: answered ? Math.round(totalMs / answered) : 0, misconceptions }
}

// labels is Pattern.misconceptions of the attempt's pattern; ids are only
// unique per pattern. Unknown ids fall back to the id itself.
export const misconceptionLabel = (labels: Record<string, string>, id: string) => labels[id] ?? id
