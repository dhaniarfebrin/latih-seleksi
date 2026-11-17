import { useEffect, useReducer, useState } from 'react'
import { createAttempt, getDueCount, listPatterns, putResponse, submitAttempt } from './api'
import { initialSession, reduce } from './session'
import type { Pattern } from './types'
import { errorMessage } from '@/shared/lib/error'

// Server state of a drill: the pattern list, the session reducer and the calls
// that move it. Components get plain callbacks and never touch the API.
export function useDrill() {
  const [s, dispatch] = useReducer(reduce, initialSession)
  const [patterns, setPatterns] = useState<Pattern[]>([])
  const [patternsError, setPatternsError] = useState('')
  const [dueCount, setDueCount] = useState(0)

  useEffect(() => {
    let live = true
    listPatterns().then(
      (p) => live && setPatterns(p),
      (e) => live && setPatternsError(errorMessage(e)),
    )
    // The review button is optional: a failed count just hides it.
    getDueCount().then(
      (n) => live && setDueCount(n),
      () => undefined,
    )
    return () => {
      live = false
    }
  }, [])

  const start = async (patternId: string, n: number, band: number, targetSec: number | null) => {
    dispatch({ type: 'started', attempt: await createAttempt(patternId, n, band), targetSec })
  }

  const answer = async (attemptId: string, pos: number, key: string, timeMs: number) => {
    dispatch({ type: 'answered', pos, response: await putResponse(attemptId, pos, key, timeMs) })
  }

  return {
    s,
    patterns,
    patternsError,
    dueCount,
    start,
    answer,
    submit: submitAttempt,
    next: () => dispatch({ type: 'next' }),
    reset: () => dispatch({ type: 'reset' }),
  }
}
