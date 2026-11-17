import { CalendarClock, ChartNoAxesColumn, Download, ListChecks, Timer } from 'lucide-react'
import { Link } from 'react-router-dom'
import { AppPreview } from './AppPreview'

const FEATURES = [
  {
    Icon: ListChecks,
    title: 'Drill dengan langkah',
    body: 'Soal TIU dan TWK dengan pembahasan langkah demi langkah dan penanda jenis kesalahan Anda.',
  },
  {
    Icon: Timer,
    title: 'Simulasi SKD berwaktu',
    body: 'Ujian berwaktu yang dinilai server, dengan ambang sesuai Kepmen yang berlaku.',
  },
  {
    Icon: ChartNoAxesColumn,
    title: 'Dasbor titik lemah',
    body: 'Lihat pola soal yang paling sering salah dan apa yang sebaiknya dilatih lebih dulu.',
  },
  {
    Icon: CalendarClock,
    title: 'Ulangan terjadwal',
    body: 'Fakta TWK dan pola TIU yang salah dijadwalkan ulang dengan algoritma FSRS.',
  },
  {
    Icon: Download,
    title: 'Data milik Anda',
    body: 'Unduh seluruh riwayat latihan atau hapus akun dan data kapan saja.',
  },
]

// Auth-free so it can be rendered to static HTML at build time.
export function LandingContent() {
  return (
    <div className="lp">
      <section className="lp-hero">
        <p className="eyebrow">Latihan seleksi CPNS</p>
        <h1 className="d">
          Latihan SKD yang menunjukkan <mark>di mana Anda salah</mark>
        </h1>
        <p className="lead">
          Latih TIU dan TWK, coba simulasi SKD berwaktu, dan lihat titik lemah Anda. Semua soal dibuat sendiri dan kunci
          jawabannya dihitung atau dikutip dari sumber, bukan ditebak.
        </p>
        <Link to="/login" className="btn hi">
          Mulai berlatih
        </Link>
      </section>
      <section aria-labelledby="lp-fitur">
        <h2 id="lp-fitur">Yang Anda dapatkan</h2>
        <ul className="lp-feats">
          {FEATURES.map(({ Icon, title, body }) => (
            <li key={title}>
              <h3>
                <Icon className="ic" aria-hidden="true" /> {title}
              </h3>
              <p>{body}</p>
            </li>
          ))}
        </ul>
      </section>
      <section aria-labelledby="lp-intip">
        <h2 id="lp-intip">Intip aplikasinya</h2>
        <AppPreview />
      </section>
      <section aria-labelledby="lp-langkah">
        <h2 id="lp-langkah">Cara memakainya</h2>
        <ol className="lp-steps">
          <li>Buat akun dan verifikasi email Anda.</li>
          <li>Kerjakan drill pada pola yang ingin dilatih, lalu baca langkah dan jenis salah Anda.</li>
          <li>Ikuti simulasi berwaktu, lalu periksa dasbor untuk menentukan latihan berikutnya.</li>
        </ol>
      </section>
      <section className="lp-cta" aria-labelledby="lp-mulai">
        <h2 id="lp-mulai">Siap berlatih?</h2>
        <Link to="/login" className="btn hi">
          Masuk atau daftar
        </Link>
        <p className="lp-note">
          Aplikasi latihan independen, bukan produk resmi penyelenggara seleksi atau instansi pemerintah. Soal buatan
          sendiri dan tidak menjamin kelulusan.
        </p>
      </section>
    </div>
  )
}
