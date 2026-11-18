# PRD: Latih Seleksi

Status: draf v0.1. Sumber utama: `docs/research/laporan-riset.md` (riset pasar, keluhan pengguna, pendekatan AI, spesifikasi resmi) dan catatan di `docs/research/notes/`. Klaim di sini merujuk ke dokumen itu; angka tanpa rujukan ditandai **Asumsi (belum pasti)**.

## 1. Ringkasan

Aplikasi web mobile-first untuk latihan seleksi CPNS (SKD: TWK, TIU, TKP) dan, pada fase berikutnya, Rekrutmen Bersama BUMN (TKD, AKHLAK, English, learning agility). Fokusnya bukan bank soal terbesar, melainkan loop tertutup **simulasi → diagnosis per pola → drill terarah → ulang**, dengan kunci jawaban yang dapat diperiksa dan penjelasan per soal.

Dua tujuan sekaligus:
1. Alat latihan untuk calon peserta seleksi ASN dan BUMN yang ingin memperbaiki pola soal yang sering salah.
2. Proyek terbuka: backend production, auth, testing, dan app yang benar-benar dirilis.

## 2. Masalah dan bukti

| Masalah | Bukti (lihat laporan, bagian "Celah") | Kekuatan |
|---|---|---|
| Kunci jawaban atau pembahasan salah | 12,7% dari ulasan Play bintang 1–3 di 8 dari 13 app inti; catatan versi MisiPNS memperbaiki kunci dua kali | Sedang |
| Pembayaran, login, dan akses premium rapuh | 26,5% dari ulasan bintang 1–3 (bug, pembayaran, login) | Sedang |
| Soal tidak mirip ujian asli atau usang | 6,4% dari ulasan bintang 1–3 | Lemah–sedang |
| Tidak ada loop dari kelemahan ke drill terarah | Tidak ada klaim "adaptif" di halaman produk yang dibaca; permintaan eksplisit di ulasan hanya sekitar 4 | Lemah, ini taruhan diferensiasi |
| Umpan balik setelah salah harus otomatis | RCT Khanmigo: siswa hanya bertanya ke tutor di 17% sesi yang berisi kesalahan (NBER w35620) | Sedang, prinsip desain umum |
| Skor TKP sulit dipahami | Sumber lama (2018) dan komersial; jarang di ulasan Play (1,5%) | Lemah |

Hipotesis kerja (belum diuji ke banyak pengguna): kesalahan SKD terkonsentrasi pada pola rutin TIU (persentase berantai, kerja-sama, kecepatan relatif, rata-rata gabungan, silogisme), bukan pada topik secara umum.

## 3. Pengguna

- **Primer:** pelamar CPNS lulusan S1, usia ≤35, banyak yang sedang bekerja, latihan di ponsel dalam sesi pendek. Bagian "sesi pendek" berstatus **Asumsi (belum pasti)**, karena tidak ada di ulasan dan perlu divalidasi lewat wawancara.
- **Sekunder (fase 2):** pelamar BUMN, baik jalur fresh graduate maupun berpengalaman.

## 4. Sasaran dan bukan sasaran

**Sasaran MVP**
- Simulasi SKD yang setia pada KepmenPANRB 321/2024.
- Setiap soal punya kunci yang dapat diperiksa dan penjelasan.
- Drill otomatis dari pola terlemah.
- Dasbor terhadap ambang per subtes.

**Bukan sasaran MVP**
- Leaderboard atau battle.
- Wawancara suara.
- Video.
- Soal figural.
- Modul BUMN English dan learning agility, karena formatnya belum jelas.
- Kalibrasi IRT.
- App native.
- Pembayaran. Ini tertunda sampai ada pengguna; kalau dibangun nanti, entitlement wajib dikelola di sisi server.

## 5. Metrik keberhasilan

- **Efektivitas latihan:** pada pengguna uji, skor TIU di drill pola naik ke ≥7/10 dalam 9 menit, dan simulasi penuh lolos ketiga ambang.
- **Kualitas konten:** laporan soal salah kurang dari 1% per 1.000 jawaban. Target ini **Asumsi (belum pasti)** dan perlu dikalibrasi setelah data pertama masuk.
- **Produk (jika dibuka publik):** retensi D7 dan jumlah drill per pengguna aktif. Target ditetapkan setelah beta.

## 6. Kebutuhan fungsional MVP

| ID | Kebutuhan | Kriteria penerimaan |
|---|---|---|
| FR-01 | Simulasi SKD | Susunan 110 soal (TWK 30, TIU 35, TKP 45). Timer 100 menit. Navigasi nomor soal dengan status terjawab/belum. Submit otomatis saat waktu habis. Skor TWK/TIU 5/0, TKP 1–5, kosong 0. |
| FR-02 | Hasil terhadap ambang | Skor per subtes dan total dibandingkan ambang dari spesifikasi aktif: umum 65/80/166, Cumlaude/Diaspora 311 dengan TIU ≥85, afirmasi 286 dengan TIU ≥60. Status gugur per subtes tampil jelas. Spesifikasi diberi label tahun. |
| FR-03 | Mode review | Per soal: jawaban user, kunci, waktu yang dipakai dibanding median, penjelasan, dan untuk TKP alasan bobot tiap opsi. Tombol "Laporkan soal". |
| FR-04 | Penjelasan otomatis setelah salah | Di mode drill, penjelasan langsung muncul setelah jawaban salah tanpa diminta. |
| FR-05 | Drill per pola | User memilih pola atau sistem memilih pola terlemah. Sesi 5–10 soal dengan timer opsional dan target detik per soal. |
| FR-06 | Generator TIU numerik | Minimal 6 keluarga templat: persentase berantai, kerja-sama, kecepatan relatif, rata-rata gabungan, deret angka, perbandingan. Kunci dan langkah dihitung kode. Pengecoh berasal dari miskonsepsi berlabel. Generator deterministik untuk seed yang sama. |
| FR-07 | Generator silogisme | Dibangun dari aturan logika formal. Keabsahan tiap opsi diperiksa kode. |
| FR-08 | Bank TWK dengan grounding | Setiap soal menyimpan kutipan sumber (pasal UUD atau dokumen sejarah) beserta rujukannya. Soal tanpa kutipan tidak bisa terbit. |
| FR-09 | Bank TKP | Situasi dan opsi boleh didraf LLM, tapi bobot 1–5 ditetapkan dan ditinjau manusia. Setiap soal terpetakan ke satu dari 6 aspek resmi. |
| FR-10 | Dasbor kelemahan | Akurasi dan waktu per subtes dan per pola, tren antar sesi, label penyebab kesalahan (ceroboh, tekanan waktu, kurang paham) dari sinyal murah, dan antrean drill berikutnya. |
| FR-11 | Penjadwalan ulang | Soal fakta TWK dan pola yang salah dijadwalkan ulang dengan FSRS parameter bawaan. Semua respons dicatat sejak hari pertama untuk kalibrasi nanti. |
| FR-12 | Akun minimal | Login (email magic link atau OAuth). Progres tersinkron antar perangkat. Ekspor dan hapus data sendiri. |
| FR-13 | Disclaimer dan tautan resmi | Disclaimer independen tampil di onboarding dan footer, ada tautan ke simulasi resmi cat.bkn.go.id, tanpa nama atau logo BKN/KemenPANRB/FHCI. |

## 7. Invarian konten (wajib, tidak boleh dilanggar)

1. **LLM tidak pernah menjadi satu-satunya sumber kebenaran kunci jawaban.** Kunci harus berasal dari kode (TIU numerik, silogisme), kutipan sumber (TWK), atau manusia (bobot TKP).
2. **Urutan pipeline: kunci tetap → LLM menulis penjelasan yang dibatasi pada kunci → pemeriksaan otomatis → review → terbit.** Pemeriksaan otomatis mencakup angka akhir yang cocok dan pasal yang benar-benar dikutip.
3. **Soal verbal dan TWK buatan LLM diuji solver tiga kali dengan urutan opsi diacak.** Soal yang tidak disepakati masuk review manusia, begitu juga sampel acak dari soal yang disepakati.
4. **Bank soal dibuat lebih dulu (pra-generasi) secara batch,** bukan on-the-fly. Pra-generasi lebih murah dan lebih aman (laporan, bagian biaya). LLM on-the-fly hanya untuk chat penjelasan berkuota.
5. **Semua soal orisinal.** Dilarang menyalin soal bocoran, rekaman ujian, atau bank soal pihak lain.
6. **Tidak ada data pribadi pengguna dalam prompt LLM.**
7. **Status soal:** `draft → verified → published → retired`. Laporan soal yang valid menurunkan status dan tercatat di riwayat koreksi yang bisa dilihat.

## 8. Spesifikasi ujian sebagai konfigurasi berversi

Spesifikasi disimpan sebagai data, tidak di-hardcode. Contoh: `skd-2024` dari KepmenPANRB 321/2024 berisi komposisi, durasi, skoring, ambang per jalur, dan kisi-kisi per aspek. Saat Kepmen CPNS 2026 terbit, cukup tambah versi baru. Saat ini (Nov 2026) Kepmen itu belum ditemukan, lihat `docs/research/notes/spesifikasi_resmi.md`.

## 9. Sketsa model data

- `exam_spec`: id, versi, sumber hukum, komposisi, durasi, aturan skor, ambang.
- `topic` / `pattern`: hierarki subtes → aspek → pola, misalnya `TIU.numerik.persentase-berantai`.
- `item`: id, pattern_id, jenis (`template|bank`), stem, options, kunci (atau bobot per opsi), penjelasan, kutipan sumber, status, versi, provenance (generator/seed/model/reviewer).
- `item_template`: id, pattern_id, parameter, solver, pembuat pengecoh beserta label miskonsepsi.
- `attempt` / `response`: user, item (atau template + seed), jawaban, benar/skor, waktu_ms, mode (`simulasi|drill`).
- `mastery`: user × pattern, akurasi, waktu, status FSRS.
- `item_report`: item, user, alasan, status tindak lanjut.

## 10. Arsitektur dan stack

Rincian dan alasannya ada di `docs/architecture.md` §3.

- **Backend:** Go + PostgreSQL. Dipilih karena bisa dikemas sebagai satu binary kecil dalam container. Isinya generator templat, skoring, penjadwal FSRS, dan API.
- **Frontend:** SPA statis Vite + React + TypeScript, mobile-first. PWA (service worker, antrean jawaban offline) direncanakan tetapi belum terpasang.
- **Auth:** Firebase Authentication (Google dan email/sandi). Backend memverifikasi ID token sendiri lewat JWKS dan hanya menyimpan `sub`.
- **Hosting (rencana, belum pernah di-deploy):** Cloudflare (SPA statis), Cloud Run (API Go), Neon (PostgreSQL). Lokal memakai Docker Compose.
- **Pipeline konten:** CLI/batch terpisah dari runtime (generate → verify → review → publish), dengan LLM murah via Batch API.
- **Mobile native (fase 3):** React Native, setelah web stabil.

## 11. Kebutuhan non-fungsional

- **Mobile:** layar 360–430px, tanpa scroll horizontal. Soal termuat kurang dari 1 detik di 4G (**Asumsi (belum pasti)**, target internal).
- **Simulasi tahan gangguan:** jawaban tersimpan lokal, sesi bisa dilanjutkan setelah reload atau sinyal putus.
- **Privasi (UU 27/2022):** minimisasi data, persetujuan eksplisit, kebijakan privasi, hapus akun.
- **Aksesibilitas:** navigasi keyboard, kontras, ukuran teks yang bisa diatur.
- **Testing:** unit test untuk setiap solver/generator (kunci = hasil hitung independen) dan skoring, plus E2E alur simulasi.

## 12. Tahapan rilis

| Fase | Isi | Keluaran |
|---|---|---|
| 0. Validasi | 3–5 wawancara pelamar, prototipe 1 keluarga templat (persentase berantai) | Keputusan lanjut atau ubah |
| 1. MVP SKD | FR-01 s.d. FR-13, bank awal sesuai kapasitas review | Dipakai sendiri, lalu beta tertutup |
| 2. BUMN | TKD + AKHLAK (sumber resmi nilai AKHLAK dicek ulang), English struktur dan bacaan | Modul BUMN |
| 3. Perluasan | Chat penjelasan berkuota, soal figural (SVG generator), React Native, monetisasi | Rilis publik |

## 13. Pertanyaan terbuka

1. Stack final: Go atau Laravel? Keycloak self-host atau auth terkelola?
2. Hosting dan anggaran bulanan. Batas budget AI personal sudah ada; untuk hosting belum ditetapkan.
3. Siapa peninjau kedua untuk bobot TKP dan soal TWK?
4. Apakah sejak awal dibuka publik, atau hanya untuk pribadi sampai fase 1 selesai?
5. Nama produk final dan domain. "Latih Seleksi" masih nama kerja.

## 14. Referensi

- Laporan riset: `docs/research/laporan-riset.md`
- Catatan riset: `docs/research/notes/` (kompetitor, keluhan pengguna, ulasan Play Store, AI dan adaptif, spesifikasi resmi)
- Dasar hukum SKD: KepmenPANRB 321/2024, https://data-jdih.menpan.go.id/dokumen/2024kepmenpanrb321.pdf
