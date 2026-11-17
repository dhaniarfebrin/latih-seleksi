# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Pelamar seleksi ASN (CPNS lebih dulu) dan, pada fase berikutnya, Rekrutmen Bersama BUMN. Profil kerja: lulusan S1, usia ≤35, banyak yang sedang bekerja (docs/prd.md §3). Mereka berlatih untuk SKD (TWK, TIU, TKP) dan ingin memperbaiki jenis soal yang sering salah, bukan sekadar menambah jumlah soal.

Perangkat: ponsel dan laptop setara (dikonfirmasi pemilik, 2026-11-17). Layout 360–430 px dan desktop harus sama kuat; simulasi 100 menit dan latihan singkat sama-sama ada. "Sesi pendek di ponsel" di PRD tetap **asumsi (belum divalidasi)**: belum ada wawancara pengguna.

## Product Purpose

Aplikasi latihan seleksi dengan loop tertutup: simulasi → diagnosis per pola → drill terarah → ulang. Kunci jawaban dapat diperiksa, penjelasan per soal ada, dan dasbor membandingkan hasil dengan ambang per subtes. Dua tujuan sekaligus: alat latihan yang berguna, dan proyek terbuka dengan backend production, auth, tes, dan rilis nyata.

Sasaran rilis: publik dan gratis (dikonfirmasi pemilik, 2026-11-17). Rilis menunggu FR-01..FR-13 lengkap (NOTES.md, keputusan 2026-11-17). Pembayaran ditunda sampai ada pengguna (docs/prd.md §4).

## Positioning

Bukan bank soal terbesar. Pembeda yang tidak bisa ditiru tanpa membangun hal yang sama:
- Kunci jawaban tidak pernah bergantung pada LLM semata: dari kode (TIU numerik, silogisme), kutipan sumber (TWK), atau manusia (bobot TKP).
- Penjelasan muncul otomatis setelah salah, dan pengecoh punya label miskonsepsi.
- Soal numerik dibangkitkan ulang dengan angka baru, jadi yang dilatih adalah pola, bukan hafalan.
- Penjadwalan ulang fakta TWK dan pola TIU dengan FSRS.

## Operating Context

- Spesifikasi ujian: KepmenPANRB 321/2024 (komposisi 110 soal: TWK 30, TIU 35, TKP 45; 100 menit; skor TWK/TIU 5/0, TKP 1–5). Ambang umum 65/80/166. Disimpan sebagai data berversi dan berlabel tahun; jangan dicampur dengan spesifikasi lain (TKP 156 vs 166).
- Login wajib untuk drill dan simulasi (Firebase Authentication, Google dan email); hanya `sub` yang disimpan.
- Server memegang timer dan skoring; kunci tidak dikirim ke klien sebelum submit; jawaban simulasi tahan reload.
- Navigasi mobile: 5 tab (Beranda, Drill, Simulasi, Riwayat, Akun), disembunyikan saat ujian berjalan.

## Capabilities and Constraints

- FR-01..FR-13 sudah ada di kode (NOTES.md, Open): simulasi, hasil vs ambang, review per soal dan laporan soal, drill TIU dan TWK, dasbor kelemahan, FSRS, akun dengan ekspor dan hapus data, disclaimer.
- Bobot TKP adalah penilaian penyusun yang ditinjau manusia, bukan kunci resmi; UI wajib tetap menyebutnya.
- Label penyebab salah (ceroboh, tekanan waktu, kurang paham) adalah perkiraan dari waktu, bukan diagnosis pasti.
- Tidak ada data pribadi pengguna dalam prompt LLM. Bank soal pra-generasi, bukan on-the-fly.
- Bukan sasaran MVP: leaderboard, wawancara suara, video, soal figural, modul BUMN, kalibrasi IRT, app native, pembayaran.
- Tidak diputuskan: domain produksi, hosting final (rencana: Cloudflare statis + Cloud Run + Neon, belum dideploy), peninjau kedua untuk bobot TKP.

## Brand Commitments

- "Latih Seleksi" adalah nama kerja dan boleh berubah (dikonfirmasi pemilik, 2026-11-17).
- Suara: bahasa Indonesia, sapaan informal "kamu", langsung dan tanpa basa-basi (seperti README dan landing). Mengikat.
- Semua soal orisinal. Disclaimer independen wajib tampil (onboarding dan footer) dengan tautan ke simulasi resmi cat.bkn.go.id. Dilarang memakai nama atau logo BKN, KemenPANRB, atau FHCI.
- Repo publik: tidak ada data pribadi di berkas yang dilacak; target generik semua pelamar ASN/BUMN.

## Evidence on Hand

- Riset pasar, keluhan pengguna, dan spesifikasi resmi: `docs/research/` (laporan-riset.md, notes/). Klaim PRD merujuk ke sana.
- Bank soal: TWK 75 item `reviewed_by: human`, TKP 45 item `reviewed`/`human`; kutipan verbatim dari naskah resmi di `api/internal/bank/data/`.
- Tidak ada: testimoni, jumlah pengguna, metrik retensi, atau hasil wawancara pengguna. Jangan dibuat-buat di landing atau di mana pun.
- Mockup lama: `docs/mockup/latih-seleksi.html` (punya perbaikan belum di-render; bukan otoritas visual).

## Product Principles

1. Dapat dipercaya sebelum menarik: kunci, sumber, dan batas klaim (bobot TKP bukan kunci resmi, label penyebab = perkiraan) selalu tampak jujur di UI.
2. Diagnosis menuju tindakan: setiap hasil mengarah ke drill pola berikutnya, bukan berhenti di skor.
3. Tahan gangguan: sesi ujian harus pulih setelah reload atau sinyal putus; ponsel dan laptop sama-sama kelas satu.
4. Orisinal dan independen: tidak meniru atau mengklaim afiliasi dengan penyelenggara seleksi.
5. Tidak ada klaim tanpa bukti: angka, testimoni, dan jaminan kelulusan tidak ditampilkan kecuali ada sumbernya.

## Accessibility & Inclusion

Dari docs/prd.md §11: navigasi keyboard, kontras memadai, ukuran teks yang bisa diatur. Status tidak pernah hanya warna (ikon + teks); ikon dekoratif `aria-hidden` dan setiap kontrol punya label teks (NOTES.md, konvensi 2026-11-15). Standar formal (mis. WCAG AA) belum ditetapkan pemilik.
