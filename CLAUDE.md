# Latih Seleksi

Aplikasi latihan seleksi CPNS/BUMN. Baca `docs/prd.md` sebelum merancang atau mengubah fitur. Bukti dan sumbernya ada di `docs/research/`.

## Invarian (lihat PRD §7)
- LLM tidak pernah menjadi satu-satunya sumber kebenaran kunci jawaban. Kunci berasal dari kode (TIU numerik, silogisme), kutipan sumber (TWK), atau manusia (bobot TKP).
- Penjelasan ditulis setelah kunci tetap, lalu diperiksa otomatis sebelum terbit.
- Bank soal dipra-generasi secara batch. LLM on-the-fly hanya untuk chat berkuota.
- Semua soal orisinal. Tidak ada nama atau logo BKN/KemenPANRB/FHCI, dan disclaimer wajib tampil.
- Tidak ada data pribadi pengguna dalam prompt LLM.
- Spesifikasi ujian (komposisi, durasi, skor, ambang) disimpan sebagai data berversi, tidak di-hardcode.
- Setiap generator atau solver wajib punya unit test yang mencocokkan kunci dengan hitungan independen.
