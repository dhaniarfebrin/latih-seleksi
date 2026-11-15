export function About() {
  return (
    <section className="screen narrow">
      <p className="eyebrow">Tentang</p>
      <h1 className="d">
        Tentang <mark>Latih Seleksi</mark>
      </h1>
      <div className="panel">
        <h2 className="sub">Catatan penting</h2>
        <p>
          Aplikasi latihan independen, bukan produk resmi penyelenggara seleksi atau instansi pemerintah. Semua soal dibuat sendiri
          dan tidak menjamin kelulusan.
        </p>
      </div>
      <div className="panel">
        <h2 className="sub">Sumber spesifikasi ujian</h2>
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
          Simulasi saat ini parsial (TIU saja) dengan durasi prorata, bukan durasi resmi. Ambang batas tiap jalur mengikuti sumber di atas;
          ambang CPNS 2026 belum diumumkan.
        </p>
      </div>
      <div className="panel">
        <h2 className="sub">Pembuat</h2>
        <p>Dibuat oleh Dhaniar. Saran dan masukan boleh lewat:</p>
        <ul className="steps">
          <li>
            <a href="https://dhaniar.pages.dev" target="_blank" rel="noopener noreferrer">
              Web: dhaniar.pages.dev
            </a>
          </li>
          <li>
            <a href="https://www.linkedin.com/in/dhaniar" target="_blank" rel="noopener noreferrer">
              LinkedIn: dhaniar
            </a>
          </li>
          <li>
            <a href="https://www.instagram.com/holyfebrn" target="_blank" rel="noopener noreferrer">
              Instagram: @holyfebrn
            </a>
          </li>
        </ul>
      </div>
      <p>
        <a href="https://cat.bkn.go.id" target="_blank" rel="noopener noreferrer">
          Simulasi resmi (cat.bkn.go.id)
        </a>
      </p>
    </section>
  )
}
