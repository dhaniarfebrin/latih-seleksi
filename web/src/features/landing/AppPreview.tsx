import { CircleCheck, CircleX } from 'lucide-react'

// Static sample of the app (own example data). Real text stays in the markup for crawlers.
export function AppPreview() {
  return (
    <div className="lp-preview">
      <div className="panel" aria-label="Contoh soal drill">
        <p className="eyebrow">Contoh drill - Persentase berantai</p>
        <p className="stem">Harga barang Rp200.000 naik 20%, lalu turun 10%. Berapa harga akhirnya?</p>
        <ul className="opts">
          <li className="opt ok">
            <span className="k">B</span>
            <span>Rp216.000</span>
            <span className="tag">
              <CircleCheck className="ic" aria-hidden="true" />
              Benar
            </span>
          </li>
          <li className="opt bad">
            <span className="k">C</span>
            <span>Rp220.000</span>
            <span className="tag">
              <CircleX className="ic" aria-hidden="true" />
              Pilihan Anda
            </span>
          </li>
        </ul>
        <ol className="steps">
          <li>Naik 20%: 200.000 x 1,2 = 240.000.</li>
          <li>Turun 10%: 240.000 x 0,9 = 216.000.</li>
        </ol>
        <p className="chip warn">Jenis salah: menjumlahkan persen langsung (+20% - 10% = +10%)</p>
      </div>
      <div className="panel" aria-label="Contoh panel titik lemah">
        <p className="eyebrow">Contoh titik lemah</p>
        <ul className="tally">
          <li>
            <span>Persentase berantai</span>
            <span className="num">55%</span>
          </li>
          <li>
            <span>Deret angka</span>
            <span className="num">78%</span>
          </li>
          <li>
            <span>Silogisme</span>
            <span className="num">91%</span>
          </li>
        </ul>
        <p className="chip warn">Latih dulu: Persentase berantai</p>
      </div>
    </div>
  )
}
