import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createExam, listSpecs } from './api'
import { availableSections, toggleSection } from './exam'
import type { Spec } from './types'
import { errorMessage } from '@/shared/lib/error'

export function StartPage() {
  const navigate = useNavigate()
  const [specs, setSpecs] = useState<Spec[]>([])
  const [specId, setSpecId] = useState('')
  const [trackId, setTrackId] = useState('')
  const [picked, setPicked] = useState<string[] | null>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    listSpecs()
      .then((all) => setSpecs(all.filter((sp) => availableSections(sp).length > 0)))
      .catch((e) => setError(errorMessage(e)))
  }, [])

  const spec = specs.find((sp) => sp.id === specId) ?? specs[0]
  const track = spec?.tracks.find((t) => t.id === trackId) ?? spec?.tracks[0]
  const open = spec ? availableSections(spec) : []
  const chosen = picked ?? open
  const partial = !!spec && chosen.length < spec.sections.length

  const start = async () => {
    if (!spec || !track) return
    setBusy(true)
    setError('')
    try {
      const exam = await createExam(spec.id, track.id, chosen)
      navigate(`/simulation/${exam.id}`)
    } catch (e) {
      setError(errorMessage(e))
      setBusy(false)
    }
  }

  return (
    <section className="screen narrow">
      <h1 className="d">
        Simulasi <mark>SKD</mark>
      </h1>
      <p className="muted">
        {partial ? 'Simulasi parsial, durasi prorata (bukan durasi resmi). ' : 'Simulasi lengkap. '}Waktu dihitung server dan jalan terus walau halaman ditutup.
      </p>
      {specs.length > 1 && (
        <fieldset className="field">
          <legend>Spesifikasi</legend>
          <div className="seg one">
            {specs.map((sp) => (
              <label key={sp.id} className={spec?.id === sp.id ? 'on' : ''}>
                <input
                  type="radio"
                  name="spec"
                  checked={spec?.id === sp.id}
                  onChange={() => {
                    setSpecId(sp.id)
                    setTrackId('')
                    setPicked(null)
                  }}
                />
                {sp.label}
              </label>
            ))}
          </div>
        </fieldset>
      )}
      {spec && (
        <fieldset className="field">
          <legend>Subtes</legend>
          <div className="seg one">
            {spec.sections.map((s) => (
              <label key={s.id} className={chosen.includes(s.id) ? 'on' : ''}>
                <input
                  type="checkbox"
                  checked={chosen.includes(s.id)}
                  disabled={!s.available}
                  onChange={() => setPicked(toggleSection(chosen, s.id))}
                />
                {s.label} ({s.items} soal)
              </label>
            ))}
          </div>
        </fieldset>
      )}
      {spec && spec.tracks.length > 0 && (
        <fieldset className="field">
          <legend>Jalur</legend>
          <div className="seg one">
            {spec.tracks.map((t) => (
              <label key={t.id} className={track?.id === t.id ? 'on' : ''}>
                <input type="radio" name="track" checked={track?.id === t.id} onChange={() => setTrackId(t.id)} />
                {t.label}
              </label>
            ))}
          </div>
        </fieldset>
      )}
      {error && (
        <p role="alert" className="err">
          {error}
        </p>
      )}
      <button className="btn hi" disabled={!spec || !track || busy} onClick={() => void start()}>
        {busy ? 'Menyiapkan...' : 'Mulai simulasi'}
      </button>
    </section>
  )
}
