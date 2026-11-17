import { call } from '@/shared/api/client'

// First pattern of the dashboard's drill queue (api/internal/insight), or null
// when the queue is empty. Only the field drill needs is typed here.
export const getWeakestPatternId = () =>
  call<{ queue: { pattern_id: string }[] }>('GET', '/v1/me/dashboard').then((d) => d.queue[0]?.pattern_id ?? null)
