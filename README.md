# Latih Seleksi

Latihan seleksi ASN dan BUMN (SKD CPNS lebih dulu) yang fokus pada satu hal: **pola soal yang sering kamu salah**.

Banyak aplikasi latihan memberi ribuan soal, tapi kamu tetap salah di jenis soal yang sama. Di sini alurnya dibalik: kerjakan soal, lihat kenapa jawabanmu salah, lalu ulangi pola itu dengan angka baru sampai terbiasa.

## Yang membuatnya berbeda

- **Kunci jawaban bisa dipercaya.** Untuk soal hitungan, kunci dan langkah penyelesaian dihitung oleh kode, bukan ditebak model bahasa. Kunci jawaban tidak pernah bergantung pada LLM semata.
- **Penjelasan muncul sendiri saat kamu salah.** Tidak perlu minta. Kamu langsung melihat langkahnya.
- **Pengecoh punya nama.** Tiap opsi salah mewakili kesalahan yang umum, misalnya "menjumlahkan persen langsung". Jadi kamu tahu salahnya di konsep mana, bukan sekadar "kurang teliti".
- **Soal tak habis.** Angkanya dibangkitkan ulang, jadi yang dilatih adalah polanya, bukan hafalan jawaban.
- **Semua soal buatan sendiri.** Bukan soal ujian sebenarnya.

## Status

Proyek ini masih dikembangkan dan dipakai sendiri dulu.

Yang sudah ada:
- Generator soal **persentase berantai** (TIU numerik), dengan kunci dan langkah dari kode.
- API Go + PostgreSQL untuk mode drill: 8 soal per sesi, penjelasan setelah menjawab, waktu dan jawaban tercatat.
- Aplikasi web (React) yang nyaman di layar ponsel.

Yang direncanakan: pola TIU lain (kerja-sama, kecepatan relatif, rata-rata gabungan, deret angka, perbandingan, silogisme), simulasi SKD penuh dengan timer, TWK dan TKP, dasbor kelemahan, lalu akun pengguna. Rinciannya ada di [`docs/prd.md`](docs/prd.md).

## Coba di komputermu

Butuh Docker.

```bash
cp .env.example .env
# isi POSTGRES_PASSWORD di .env dengan kata sandi apa saja
docker compose up --build
```

Buka http://localhost:8081. Dari ponsel di WiFi yang sama, buka `http://<alamat-ip-komputermu>:8081`.

Hanya server web (nginx, berisi file statis) yang dibuka ke jaringan. API dan database hanya bisa diakses dari komputermu sendiri.

## Menjalankan tes

```bash
# Backend (Go 1.22+)
cd api
gofmt -l . && go vet ./... && go test ./...
# Tes yang butuh database aktif jika DATABASE_URL diisi, dan dilewati jika kosong.

# Frontend (Node 20+)
cd web
npm ci && npm run build && npx vitest run
```

Catatan: pada beberapa Mac dengan Go 1.22.x, `go test` bisa berhenti dengan galat `dyld: missing LC_UUID load command`. Jalankan dengan `go test -ldflags=-linkmode=external ./...`.

## Isi repo

| Folder | Isi |
|---|---|
| `api/` | API Go: generator soal (`internal/gen`), layanan drill, penyimpanan Postgres, Dockerfile |
| `web/` | Aplikasi web Vite + React + TypeScript |
| `docs/prd.md` | Kebutuhan produk dan aturan konten yang tidak boleh dilanggar |
| `docs/architecture.md` | Rancangan arsitektur (stack: SPA Vite, API Go, Postgres) |
| `docs/research/` | Riset kompetitor, keluhan pengguna, dan spesifikasi resmi beserta sumbernya |
| `docs/mockup/` | Mockup tampilan |

Catatan: `docs/prd.md` §10 masih menyebut Next.js; untuk stack, `docs/architecture.md` yang berlaku (SPA Vite statis di Cloudflare, API Go di Cloud Run, Postgres di Neon).

## Disclaimer

Aplikasi latihan independen, bukan produk resmi penyelenggara seleksi. Soal buatan sendiri dan tidak menjamin kelulusan. Untuk simulasi resmi, gunakan situs penyelenggara.
