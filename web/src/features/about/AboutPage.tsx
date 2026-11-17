import { ExternalLink, Globe } from 'lucide-react'

// A reading page: one column, sections divided by rules instead of boxed cards.
export function AboutPage() {
  return (
    <article className="screen mid doc">
      <h1 className="d">
        Tentang <mark>Latih Seleksi</mark>
      </h1>
      <p className="note" role="note">
        Aplikasi latihan independen, bukan produk resmi penyelenggara seleksi atau instansi pemerintah. Semua soal dibuat sendiri
        dan tidak menjamin kelulusan.
      </p>
      <section aria-labelledby="src-h">
        <h2 id="src-h" className="sub">
          Sumber spesifikasi ujian
        </h2>
        <ul className="steps">
          <li>
            <a href="https://data-jdih.menpan.go.id/dokumen/2024kepmenpanrb321.pdf" target="_blank" rel="noopener noreferrer">
              KepmenPANRB 321/2024
            </a>{' '}
            (SKD CPNS 2024)
          </li>
          <li>
            <a href="https://data-jdih.menpan.go.id/dokumen/2026kepmenpanrb406.pdf" target="_blank" rel="noopener noreferrer">
              KepmenPANRB 406/2026
            </a>{' '}
            (SKD sekolah kedinasan 2026)
          </li>
        </ul>
        <p className="muted">
          Simulasi lengkap memakai durasi resmi; memilih sebagian subtes memakai durasi prorata, bukan durasi resmi. Bobot TKP adalah penilaian penyusun, bukan kunci resmi. Ambang batas tiap jalur mengikuti sumber di atas;
          ambang CPNS 2026 belum diumumkan.
        </p>
        <p>
          <a href="https://cat.bkn.go.id" target="_blank" rel="noopener noreferrer">
            Simulasi resmi (cat.bkn.go.id)
          </a>
        </p>
      </section>
      <section aria-labelledby="mk-h">
        <h2 id="mk-h" className="sub">
          Pembuat
        </h2>
        <p>Dibuat oleh Dhaniar. Saran dan masukan boleh lewat:</p>
        <ul className="steps">
          <li>
            <a href="https://dhaniar.pages.dev" target="_blank" rel="noopener noreferrer">
              <Globe className="ic" aria-hidden="true" /> Web: dhaniar.pages.dev
            </a>
          </li>
          <li>
            <a href="https://www.linkedin.com/in/dhaniar" target="_blank" rel="noopener noreferrer">
              <ExternalLink className="ic" aria-hidden="true" /> LinkedIn: dhaniar
            </a>
          </li>
          <li>
            <a href="https://www.instagram.com/holyfebrn" target="_blank" rel="noopener noreferrer">
              <ExternalLink className="ic" aria-hidden="true" /> Instagram: @holyfebrn
            </a>
          </li>
        </ul>
      </section>
    </article>
  )
}
