import { useEffect, useState } from 'react'
import { errorMessage } from '@/shared/lib/error'

export interface Loaded<T> {
  data: T | undefined
  error: string
  loading: boolean
}

// Runs load once on mount and reports its outcome; a result that arrives after
// unmount is dropped.
export function useLoad<T>(load: () => Promise<T>): Loaded<T> {
  const [state, setState] = useState<Loaded<T>>({ data: undefined, error: '', loading: true })
  useEffect(() => {
    let live = true
    load().then(
      (data) => live && setState({ data, error: '', loading: false }),
      (e) => live && setState({ data: undefined, error: errorMessage(e), loading: false }),
    )
    return () => {
      live = false
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps -- load runs once on mount by contract; callers pass inline closures
  return state
}
