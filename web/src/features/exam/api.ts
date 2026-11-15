import { call } from '@/shared/api/client'
import type { Exam, ExamResponse, ExamResult, ExamSummary, Spec } from './types'

export const listSpecs = () => call<{ specs: Spec[] }>('GET', '/v1/specs').then((r) => r.specs)

export const createExam = (spec_id: string, track: string, sections: string[]) =>
  call<Exam>('POST', '/v1/exams', { spec_id, track, sections })

export const getExam = (id: string) => call<Exam>('GET', `/v1/exams/${id}`)

// seq must rise per exam: the server ignores a request whose seq is not newer.
export const putExamResponse = (
  id: string,
  pos: number,
  r: { answer: string; time_ms: number; flagged: boolean; seq: number },
) => call<ExamResponse>('PUT', `/v1/exams/${id}/responses/${pos}`, r)

export const submitExam = (id: string) => call<ExamResult>('POST', `/v1/exams/${id}/submit`)

export const getExamResult = (id: string) => call<ExamResult>('GET', `/v1/exams/${id}/result`)

// Newest first; the server caps limit at 50.
export const listExams = (limit?: number) =>
  call<{ exams: ExamSummary[] }>('GET', `/v1/exams${limit ? `?limit=${limit}` : ''}`).then((r) => r.exams)
