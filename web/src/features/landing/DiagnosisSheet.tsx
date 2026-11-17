import { CircleX } from 'lucide-react'
import { Link } from 'react-router-dom'

// Static sample result sheet (own example data, labelled as such). Real text stays in the markup for crawlers.
const ROWS = [
  { name: 'Persentase berantai', pct: 55, weakest: true },
  { name: 'Deret angka', pct: 78, weakest: false },
  { name: 'Silogisme', pct: 91, weakest: false },
]

export function DiagnosisSheet() {
  return (
    <div className="panel lp-sheet" role="group" aria-labelledby="lp-sheet-h">
      <div className="lp-sheet-head">
        <h2 id="lp-sheet-h">Titik lemahmu</h2>
        <span className="chip">Contoh</span>
      </div>
      <ul className="sheet-rows">
        {ROWS.map(({ name, pct, weakest }, i) => (
          <li key={name} className={weakest ? 'sheet-row weak' : 'sheet-row'}>
            <span className="sheet-name">{name}</span>
            <span className="sheet-bar" aria-hidden="true">
              <i style={{ width: `${pct}%`, animationDelay: `${i * 90}ms` }} />
            </span>
            <span className="num">{pct}%</span>
            {weakest && (
              <span className="sheet-why">
                <CircleX className="ic" aria-hidden="true" />
                Sering menjumlahkan persen langsung: +20% lalu -10% dianggap +10%.
              </span>
            )}
          </li>
        ))}
      </ul>
      <Link to="/login" className="btn lp-sheet-cta">
        Latih pola ini
      </Link>
    </div>
  )
}
