import { useEffect, useReducer, useState } from 'react'
import { createAttempt, listPatterns, putResponse, submitAttempt } from './api'
import { initialSession, reduce } from './session'
import type { Pattern } from './types'
import { errorMessage } from '@/shared/lib/error'

// Server state of a drill: the pattern list, the session reducer and the calls
// that move it. Components get plain callbacks and never touch the API.
export function useDrill() {
  const [s, dispatch] = useReducer(reduce, initialSession)
  const [patterns, setPatterns] = useState<Pattern[]>([])
  const [patternsError, setPatternsError] = useState('')

  useEffect(() => {
    let live = true
    listPatterns().then(
      (p) => live && setPatterns(p),
      (e) => live && setPatternsError(errorMessage(e)),
    )
    return () => {
      live = false
    }
  }, [])

  const start = async (patternId: string, n: number, band: number) => {
    dispatch({ type: 'started', attempt: await createAttempt(patternId, n, band) })
  }

  const answer = async (attemptId: string, pos: number, key: string, timeMs: number) => {
    dispatch({ type: 'answered', pos, response: await putResponse(attemptId, pos, key, timeMs) })
  }

  return {
    s,
    patterns,
    patternsError,
    start,
    answer,
    submit: submitAttempt,
    next: () => dispatch({ type: 'next' }),
    reset: () => dispatch({ type: 'reset' }),
  }
}
