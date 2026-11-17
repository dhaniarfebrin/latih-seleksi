import { describe, expect, it } from 'vitest'
import type { Attempt, ResponseView } from './types'
import { initialSession, misconceptionLabel, reduce, summarize } from './session'

const resp = (correct: boolean, time_ms: number, mis?: string): ResponseView => ({
  answer: correct ? 'A' : 'B',
  correct,
  time_ms,
  answer_key: 'A',
  steps: ['x'],
  misconception_id: mis,
  answered_at: '2026-11-14T00:00:00Z',
})

const attempt = (n: number): Attempt => ({
  id: 'id',
  mode: 'drill',
  pattern_id: 'p',
  band: 2,
  n,
  started_at: '2026-11-14T00:00:00Z',
  items: Array.from({ length: n }, (_, i) => ({
    pos: i + 1,
    stem: 's',
    options: [{ key: 'A', text: '1' }],
  })),
})

describe('session screens', () => {
  it('goes start -> question -> summary -> start', () => {
    let s = reduce(initialSession, { type: 'started', attempt: attempt(2), targetSec: null })
    expect(s.screen).toBe('question')
    expect(s.index).toBe(0)
    s = reduce(s, { type: 'next' })
    expect(s).toMatchObject({ screen: 'question', index: 1 })
    s = reduce(s, { type: 'next' })
    expect(s.screen).toBe('summary')
    s = reduce(s, { type: 'reset' })
    expect(s).toEqual(initialSession)
  })

  it('keeps the first response for a position', () => {
    let s = reduce(initialSession, { type: 'started', attempt: attempt(1), targetSec: null })
    s = reduce(s, { type: 'answered', pos: 1, response: resp(true, 1000) })
    s = reduce(s, { type: 'answered', pos: 1, response: resp(false, 2000) })
    expect(s.attempt?.items[0].response?.correct).toBe(true)
  })
})

describe('summarize', () => {
  it('counts score, average time and misconceptions', () => {
    let s = reduce(initialSession, { type: 'started', attempt: attempt(4), targetSec: null })
    s = reduce(s, { type: 'answered', pos: 1, response: resp(true, 10000) })
    s = reduce(s, { type: 'answered', pos: 2, response: resp(false, 20000, 'tanda-terbalik') })
    s = reduce(s, { type: 'answered', pos: 3, response: resp(false, 30000, 'tanda-terbalik') })
    const sum = summarize(s.attempt!.items)
    expect(sum).toEqual({
      n: 4,
      answered: 3,
      score: 1,
      avgMs: 20000,
      misconceptions: { 'tanda-terbalik': 2 },
      withinTarget: null,
    })
  })

  it('counts answers within the target', () => {
    let s = reduce(initialSession, { type: 'started', attempt: attempt(3), targetSec: 30 })
    expect(s.targetSec).toBe(30)
    s = reduce(s, { type: 'answered', pos: 1, response: resp(true, 30000) })
    s = reduce(s, { type: 'answered', pos: 2, response: resp(true, 30001) })
    expect(summarize(s.attempt!.items, 30).withinTarget).toBe(1)
  })

  it('has zero average with no answers', () => {
    expect(summarize(attempt(2).items).avgMs).toBe(0)
  })
})

describe('misconceptionLabel', () => {
  const labels = { 'tanda-terbalik': 'Tanda tertukar' }

  it('returns the label of the pattern', () => {
    expect(misconceptionLabel(labels, 'tanda-terbalik')).toBe('Tanda tertukar')
  })

  it('falls back to the id when the pattern has no label for it', () => {
    expect(misconceptionLabel(labels, 'tak-dikenal')).toBe('tak-dikenal')
    expect(misconceptionLabel({}, 'tanda-terbalik')).toBe('tanda-terbalik')
  })
})
