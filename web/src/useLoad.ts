import { useEffect, useState } from 'react'

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
      (e: Error) => live && setState({ data: undefined, error: e.message, loading: false }),
    )
    return () => {
      live = false
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])
  return state
}
