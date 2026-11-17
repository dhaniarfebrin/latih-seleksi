import { Link } from 'react-router-dom'
import { DiagnosisSheet } from './DiagnosisSheet'
import { SampleQuestion } from './SampleQuestion'

const SOURCES = [
  { name: 'Kode', body: 'Soal numerik dan silogisme dihitung program, angkanya diacak ulang tiap latihan.' },
  { name: 'Kutipan', body: 'Kunci TWK diambil dari naskah resmi, dengan kutipan dan pasalnya.' },
  { name: 'Penilaian manusia', body: 'Bobot TKP adalah penilaian penyusun yang ditinjau manusia, bukan kunci resmi.' },
]

const STEPS = [
  'Simulasi berwaktu, dinilai server.',
  'Diagnosis per pola: angka tiap jenis soal dan jenis salahmu.',
  'Drill pada pola yang paling lemah, lalu ulang lewat jadwal FSRS.',
]

// Auth-free so it can be rendered to static HTML at build time.
export function LandingContent() {
  return (
    <div className="lp">
      <section className="lp-hero">
        <div className="lp-hero-copy">
          <h1 className="d">
            Tahu persis <mark>jenis soal</mark> yang bikin kamu salah
          </h1>
          <p className="lead">
            Ikut simulasi SKD, lihat angka tiap pola soal, lalu latih yang paling lemah. Semua soal dibuat sendiri dan
            kunci jawabannya dihitung atau dikutip dari sumber, bukan ditebak.
          </p>
          <Link to="/login" className="btn hi">
            Mulai berlatih
          </Link>
        </div>
        <DiagnosisSheet />
      </section>
      <section className="lp-proof" aria-labelledby="lp-salah">
        <h2 id="lp-salah">Setiap salah dijelaskan, bukan cuma dinilai</h2>
        <div className="lp-proof-grid">
          <SampleQuestion />
          <div className="lp-sources">
            <h3>Kunci jawaban tidak pernah cuma tebakan model bahasa</h3>
            <figure className="lp-quote">
              <blockquote>
                &ldquo;bersumber dari bahasa yang diikrarkan dalam Sumpah Pemuda tanggal 28 Oktober 1928&rdquo;
              </blockquote>
              <figcaption>
                Kunci sebuah soal TWK, dari UU 24/2009 Pasal 25 ayat (1). Soal ini sudah ditinjau manusia.
              </figcaption>
            </figure>
            <dl>
              {SOURCES.map(({ name, body }) => (
                <div key={name}>
                  <dt>{name}</dt>
                  <dd>{body}</dd>
                </div>
              ))}
            </dl>
          </div>
        </div>
      </section>
      <section aria-labelledby="lp-format">
        <h2 id="lp-format">Format simulasi mengikuti Kepmen 321/2024</h2>
        <dl className="lp-spec">
          <div>
            <dt>Soal</dt>
            <dd className="num">110</dd>
          </div>
          <div>
            <dt>Waktu</dt>
            <dd className="num">100 menit</dd>
          </div>
          <div>
            <dt>TWK / TIU / TKP</dt>
            <dd className="num">30 / 35 / 45</dd>
          </div>
          <div>
            <dt>Ambang umum</dt>
            <dd className="num">65 / 80 / 166</dd>
          </div>
        </dl>
        <p className="muted">Waktu dan skor dipegang server, dan jawabanmu tetap aman kalau halaman dimuat ulang.</p>
      </section>
      <section aria-labelledby="lp-langkah">
        <h2 id="lp-langkah">Cara kerjanya</h2>
        <ol className="lp-steps">
          {STEPS.map((s) => (
            <li key={s}>{s}</li>
          ))}
        </ol>
        <p className="muted">Datamu milik kamu: unduh seluruh riwayat atau hapus akun kapan saja.</p>
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
