import { useState } from 'react'
import type { Pattern } from '../types'
import { errorMessage } from '@/shared/lib/error'
import { groupPatterns, sectionOf } from '../patterns'

const N = 8
// Banked sections (TWK) have no difficulty bands; the API still takes one.
const FIXED_BAND = 2
const BANDED_SECTION = 'TIU'
// Reserved pattern id of the review drill (api/internal/drill ReviewPatternID).
const REVIEW_PATTERN_ID = 'ULANG'
const TARGETS: [number | null, string][] = [
  [null, 'Mati'],
  [30, '30 d'],
  [60, '60 d'],
  [90, '90 d'],
]

export function DrillStart(props: {
  // Preselected pattern from the dashboard queue; an unknown id falls back to the first pattern.
  initialPatternId: string
  patterns: Pattern[]
  loadError: string
  // Cards due for review; the review button shows only when above zero.
  dueCount: number
  onStart: (patternId: string, n: number, band: number, targetSec: number | null) => Promise<void>
  // Resolves to the weakest pattern id from the dashboard queue, or null when there is none yet.
  onPickWeakest: () => Promise<string | null>
}) {
  const { initialPatternId, patterns, loadError, dueCount, onStart, onPickWeakest } = props
  const [patternId, setPatternId] = useState(initialPatternId)
  const [band, setBand] = useState(2)
  const [targetSec, setTargetSec] = useState<number | null>(null)
  const [noWeakest, setNoWeakest] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  const pattern = patterns.find((p) => p.id === patternId) ?? patterns[0]
  const banded = pattern ? sectionOf(pattern.id) === BANDED_SECTION : true
  const begin = async (id: string, useBand: boolean) => {
    setBusy(true)
    setError('')
    try {
      await onStart(id, N, useBand ? band : FIXED_BAND, targetSec)
    } catch (e) {
      setError(errorMessage(e))
      setBusy(false)
    }
  }
  const start = () => pattern && begin(pattern.id, banded)
  const startReview = () => begin(REVIEW_PATTERN_ID, false)

  const pickWeakest = async () => {
    setBusy(true)
    setError('')
    try {
      const id = await onPickWeakest()
      if (id === null) setNoWeakest(true)
      else setPatternId(id)
    } catch (e) {
      setError(errorMessage(e))
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <h1 className="d">
        <mark>Drill</mark>
      </h1>
      <p className="muted">
        {pattern ? `${N} soal ${pattern.label}` : 'Memuat pola soal...'}. Jawab, lalu lihat kunci, penjelasan, dan jenis salahnya.
      </p>
      <fieldset className="field">
        <legend>Pola soal</legend>
        {groupPatterns(patterns).map((g) => (
          <div key={g.section}>
            <p className="eyebrow">{g.section}</p>
            <div className="seg">
              {g.patterns.map((p) => (
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
          </div>
        ))}
        <button className="btn" disabled={busy || noWeakest} onClick={pickWeakest}>
          Pola terlemah
        </button>
        {noWeakest && <p className="muted">Belum ada data dasbor untuk menentukan pola terlemah.</p>}
      </fieldset>
      {banded && (
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
      )}
      <fieldset className="field">
        <legend>Target waktu per soal</legend>
        <div className="seg">
          {TARGETS.map(([v, label]) => (
            <label key={label} className={targetSec === v ? 'on' : ''}>
              <input type="radio" name="target" checked={targetSec === v} onChange={() => setTargetSec(v)} />
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
      {dueCount > 0 && (
        <button className="btn" disabled={busy} onClick={startReview}>
          Ulang yang jatuh tempo ({dueCount})
        </button>
      )}
    </>
  )
}
