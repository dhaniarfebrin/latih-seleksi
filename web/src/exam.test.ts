import { describe, expect, it } from 'vitest'
import type { Exam, ExamItem, ExamResult } from './api'
import {
  examReduce,
  initialExamSession,
  itemStatus,
  nextSeq,
  remainingSeconds,
  summarizeSections,
} from './exam'

const item = (pos: number, extra: Partial<ExamItem> = {}): ExamItem => ({
  pos,
  section: 'TIU',
  stem: 's',
  options: [],
  ...extra,
})
const exam = (status: Exam['status'] = 'in_progress', items: ExamItem[] = [item(1), item(2)]): Exam => ({
  id: 'e1',
  spec_id: 'sp',
  track: 't',
  sections: ['TIU'],
  status,
  started_at: '2026-01-01T00:00:00Z',
  deadline_at: '2026-01-01T00:10:00Z',
  duration_seconds: 600,
  duration_basis: 'prorated',
  items,
})
const result: ExamResult = {
  exam_id: 'e1',
  spec_id: 'sp',
  track: 't',
  status: 'submitted',
  partial: true,
  duration_basis: 'prorated',
  sections: [],
  total: 0,
  max_total: 0,
  submitted_at: '2026-01-01T00:05:00Z',
}
const resp = (seq: number, extra = {}) => ({ answer: 'A', time_ms: 1, flagged: false, seq, ...extra })

describe('remainingSeconds', () => {
  const dl = '2026-01-01T00:10:00Z'
  const t0 = Date.parse('2026-01-01T00:00:00Z')
  it('counts down from the deadline', () => {
    expect(remainingSeconds(dl, t0, 0)).toBe(600)
  })
  it('applies the server clock offset', () => {
    // client clock is 30s behind the server: server now = t0 + 30s
    expect(remainingSeconds(dl, t0, 30_000)).toBe(570)
    expect(remainingSeconds(dl, t0, -30_000)).toBe(630)
  })
  it('never goes negative and rounds partial seconds up', () => {
    expect(remainingSeconds(dl, t0 + 700_000, 0)).toBe(0)
    expect(remainingSeconds(dl, t0 + 599_500, 0)).toBe(1)
    expect(remainingSeconds('bad', t0, 0)).toBe(0)
  })
})

describe('nextSeq', () => {
  it('starts at 1 and rises above the stored maximum', () => {
    expect(nextSeq([item(1), item(2)])).toBe(1)
    expect(nextSeq([item(1, { response: resp(3) }), item(2, { response: resp(7) })])).toBe(8)
  })
})

describe('examReduce', () => {
  const signed = examReduce(initialExamSession, { type: 'signedIn' })
  it('moves login -> start on sign-in and back on sign-out', () => {
    expect(signed.screen).toBe('start')
    expect(examReduce(signed, { type: 'signedOut' })).toEqual(initialExamSession)
  })
  it('goes running -> result and does not return to start by itself', () => {
    const running = examReduce(signed, { type: 'started', exam: exam() })
    expect(running.screen).toBe('running')
    const done = examReduce(running, { type: 'finished', result })
    expect(done.screen).toBe('result')
    expect(examReduce(done, { type: 'signedIn' }).screen).toBe('result')
    expect(examReduce(done, { type: 'finished', result }).screen).toBe('result')
    expect(examReduce(running, { type: 'signedIn' }).screen).toBe('running')
    expect(examReduce(done, { type: 'reset' }).screen).toBe('start')
  })
  it('resumes an in-progress exam and shows a closed one as result', () => {
    expect(examReduce(signed, { type: 'resumed', exam: exam() }).screen).toBe('running')
    expect(examReduce(signed, { type: 'resumed', exam: exam('expired'), result }).screen).toBe('result')
  })
  it('keeps the newest answer per position', () => {
    let s = examReduce(signed, { type: 'started', exam: exam() })
    s = examReduce(s, { type: 'answered', pos: 1, response: resp(2, { answer: 'B' }) })
    s = examReduce(s, { type: 'answered', pos: 1, response: resp(1, { answer: 'C' }) })
    expect(s.exam?.items[0].response?.answer).toBe('B')
  })
})

describe('summarizeSections', () => {
  it('counts answered and flagged per section', () => {
    const items = [
      item(1, { response: resp(1) }),
      item(2, { response: resp(2, { flagged: true }) }),
      item(3),
      item(4, { section: 'TWK' }),
    ]
    expect(summarizeSections(items)).toEqual([
      { id: 'TIU', total: 3, answered: 2, flagged: 1 },
      { id: 'TWK', total: 1, answered: 0, flagged: 0 },
    ])
    expect(items.map(itemStatus)).toEqual(['answered', 'flagged', 'empty', 'empty'])
  })
})
