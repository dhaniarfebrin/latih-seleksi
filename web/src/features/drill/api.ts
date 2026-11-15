import { call } from '@/shared/api/client'
import type { Attempt, Pattern, ResponseView, Result } from './types'

export const listPatterns = () => call<{ patterns: Pattern[] }>('GET', '/v1/patterns').then((r) => r.patterns)

export const createAttempt = (pattern_id: string, n: number, band: number) =>
  call<Attempt>('POST', '/v1/attempts', { pattern_id, n, band })

export const getAttempt = (id: string) => call<Attempt>('GET', `/v1/attempts/${id}`)

// pos is 1-based; the first answer for a position wins and the response
// reveals the key, steps and misconception_id.
export const putResponse = (id: string, pos: number, answer: string, time_ms: number) =>
  call<ResponseView>('PUT', `/v1/attempts/${id}/responses/${pos}`, { answer, time_ms })

export const submitAttempt = (id: string) => call<Result>('POST', `/v1/attempts/${id}/submit`)
