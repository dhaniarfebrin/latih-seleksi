import { CircleCheck, CircleX } from 'lucide-react'

// Static sample of a corrected drill question (own example data). Real text stays in the markup for crawlers.
const OPTIONS = [
  { k: 'A', text: 'Rp208.000' },
  { k: 'B', text: 'Rp216.000', state: 'ok' },
  { k: 'C', text: 'Rp220.000', state: 'bad' },
  { k: 'D', text: 'Rp224.000' },
  { k: 'E', text: 'Rp240.000' },
] as const

export function SampleQuestion() {
  return (
    <div className="panel" role="group" aria-label="Contoh soal drill yang dikoreksi">
      <p className="muted">Contoh drill, pola Persentase berantai</p>
      <p className="stem">Harga barang Rp200.000 naik 20%, lalu turun 10%. Berapa harga akhirnya?</p>
      <ul className="opts">
        {OPTIONS.map((o) => {
          const state = 'state' in o ? o.state : undefined
          return (
            <li key={o.k} className={state ? `opt ${state}` : 'opt'}>
              <span className="k">{o.k}</span>
              <span>{o.text}</span>
              {state === 'ok' && (
                <span className="tag">
                  <CircleCheck className="ic" aria-hidden="true" /> Benar
                </span>
              )}
              {state === 'bad' && (
                <span className="tag">
                  <CircleX className="ic" aria-hidden="true" /> Pilihanmu
                </span>
              )}
            </li>
          )
        })}
      </ul>
      <ol className="steps">
        <li>Naik 20%: 200.000 x 1,2 = 240.000.</li>
        <li>Turun 10%: 240.000 x 0,9 = 216.000.</li>
      </ol>
      <p className="chip warn">Jenis salah: menjumlahkan persen langsung (+20% - 10% = +10%)</p>
    </div>
  )
}
