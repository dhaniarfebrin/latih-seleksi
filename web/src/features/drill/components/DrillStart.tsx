import { useState } from 'react'
import type { Pattern } from '../types'
import { errorMessage } from '@/shared/lib/error'

const N = 8

export function DrillStart(props: {
  // Preselected pattern from the dashboard queue; an unknown id falls back to the first pattern.
  initialPatternId: string
  patterns: Pattern[]
  loadError: string
  onStart: (patternId: string, n: number, band: number) => Promise<void>
}) {
  const { initialPatternId, patterns, loadError, onStart } = props
  const [patternId, setPatternId] = useState(initialPatternId)
  const [band, setBand] = useState(2)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  const pattern = patterns.find((p) => p.id === patternId) ?? patterns[0]
  const start = async () => {
    if (!pattern) return
    setBusy(true)
    setError('')
    try {
      await onStart(pattern.id, N, band)
    } catch (e) {
      setError(errorMessage(e))
      setBusy(false)
    }
  }

  return (
    <>
      <p className="eyebrow">{N} soal</p>
      <h1 className="d">
        Drill <mark>TIU</mark>
      </h1>
      <p className="muted">{pattern ? pattern.label : 'Memuat pola soal...'} (TIU). Jawab, lalu lihat langkah dan jenis salah hitungnya.</p>
      <fieldset className="field">
        <legend>Pola soal</legend>
        <div className="seg">
          {patterns.map((p) => (
            <label key={p.id} className={pattern?.id === p.id ? 'on' : ''}>
              <input
                type="radio"
                name="pattern"
                checked={pattern?.id === p.id}
                onChange={() => setPatternId(p.id)}
              />
              {p.label}
            </label>
          ))}
        </div>
      </fieldset>
      <fieldset className="field">
        <legend>Tingkat kesulitan</legend>
        <div className="seg">
          {[
            [1, 'Mudah'],
            [2, 'Sedang'],
            [3, 'Sulit'],
          ].map(([v, label]) => (
            <label key={v} className={band === v ? 'on' : ''}>
              <input type="radio" name="band" checked={band === v} onChange={() => setBand(v as number)} />
              {label}
            </label>
          ))}
        </div>
      </fieldset>
      {(loadError || error) && (
        <p role="alert" className="err">
          {error || loadError}
        </p>
      )}
      <button className="btn hi" disabled={!pattern || busy} onClick={start}>
        {busy ? 'Menyiapkan...' : 'Mulai latihan'}
      </button>
    </>
  )
}
