# Format resmi Rekrutmen Bersama BUMN (RBB)

Dibaca langsung pada 18 Nov 2026. Tujuan: memutuskan apakah modul BUMN (PRD fase 2) sudah bisa dispesifikasikan sebagai data ujian berversi, sesuai invarian "spesifikasi ujian disimpan sebagai data berversi" (CLAUDE.md, PRD §7).

**Ringkasan:** tahapan dan daftar komponen tes RBB terbukti dari sumber resmi (Kementerian BUMN). Jumlah soal, durasi, subtes, dan ambang batas TIDAK ada di satu pun sumber resmi yang bisa dibuka; angkanya hanya dari media dan situs bimbel, dan angka itu saling bertentangan. Definisi nilai AKHLAK resmi terbit dan terbukti. Contoh soal resmi tidak ditemukan.

## Lingkup pencarian dan batasnya
- Pencarian web (mode standar, sekali mode extended) dengan kata kunci RBB, TKD, AKHLAK, learning agility, bahasa Inggris, ambang batas, FHCI, BP BUMN, Danantara; dibuka: situs bumn.go.id (rilis dan dua LKIP PDF), jdih.bumn.go.id (SE-7/2020), ANTARA, detik, Kompas, Suara, IDN Times, Investortrust, Liputan6, Dealls.
- **Portal resmi peserta tidak bisa dibaca.** `rekrutmenbersama2025.fhcibumn.id`: WebFetch gagal ("certificate has expired"); percobaan saya mengambilnya dengan verifikasi TLS dimatikan ditolak oleh sistem izin dan tidak saya lanjutkan. `rekrutmenbersama.bumn.go.id`, `fhcibumn.id`, `fhcibumn.com`, `rekrutmenbersama2024.fhcibumn.id` tidak resolve dari lingkungan ini (DNS gagal, bisa jadi pembatasan lingkungan, belum tentu situs mati). Jadi kisi-kisi, FAQ, dan pengumuman di portal itu belum dibaca langsung.
- Instagram `@fhci.bumn` dan `@kementerianbumn` tidak saya buka (akun media sosial, tidak ada alat untuk itu di sesi ini); semua angka yang katanya dari sana hanya lewat kutipan media.
- Halaman bimbel "panduan lengkap tes dari FHCI" (fahum.umsu.ac.id) dialihkan ke domain lain yang tidak dikenal; tidak saya buka.
- Tidak ada kata kunci Inggris dan tidak ada pencarian di arsip PDF FHCI/Danantara di luar yang tersebut.

## 1. Tahapan dan komponen tes, dan batch terbaru

**Batch terbaru yang bisa dibuktikan: RBB 2025.** RBB 2026 tidak ditemukan sama sekali (lihat di bawah).

| Hal | Isi | Sumber (label) |
|---|---|---|
| Penyelenggara | Kementerian BUMN bekerja sama dengan Forum Human Capital Indonesia (FHCI); berjalan sejak 2019, dibuka kembali 2023 dan 2024 | LKIP Deputi SDMTI Kementerian BUMN 2024, hlm. 30 (nomor halaman cetak) / hlm. PDF 33 — [bumn.go.id PDF](https://www.bumn.go.id/storage/kontenlaporan/files/files_1755070917.pdf). Resmi |
| Tahapan RBB 2024 | Pendaftaran; Seleksi Administrasi; Tes Online Tahap 1 (Tes Kemampuan Dasar, AKHLAK, Wawasan Kebangsaan); Tes Online Tahap 2 (Tes Bahasa Inggris dan Learning Agility); Tes Kemampuan Bidang di BUMN (Psikotest, Wawancara, Medical Check Up, Social Media Analytic); Pengumuman Final. Pelaksanaan 23 Maret s.d. 11 Juli 2024 | LKIP 2024 (sama), hlm. PDF 33. Resmi |
| Tahapan RBB 2023 | Pendaftaran; Tes tahap 1 (Tes kemampuan dasar dan AKHLAK); Tes tahap 2 (Tes Bahasa Inggris); seleksi di BUMN (tes kemampuan bidang, psikotest, MCU, social media analytic); Pengumuman final | LKIP 2023 hlm. PDF 46 — [bumn.go.id PDF](https://www.bumn.go.id/storage/kontenlaporan/files/files_1728276962.pdf). Resmi. Dokumen menulis tanggal "11 Mei s.d. 1 September 2024" untuk RBB 2023; kemungkinan salah ketik, tidak saya pastikan |
| Penambahan 2024 | "sebagai penyempurnaan tahun ini dilakukan pula Tes Wawasan Kebangsaan dan Learning Agility"; "sistem gugur di setiap tahapannya"; jenjang SMA/Sederajat hingga S2, "Fresh Graduate dan Experience Hired", disabilitas lewat mekanisme terpisah | Rilis Kementerian BUMN, tanggal terbit 2024-04-23 (dateline 22 Maret 2024) — [bumn.go.id](https://bumn.go.id/publikasi/berita/rilis/detail/6627750b32729). Resmi |
| Tahapan RBB 2025 | "registrasi, seleksi administrasi, tes daring" yang terdiri dari "tes kemampuan dasar, tes AKHLAK, wawasan kebangsaan, tes bahasa Inggris, dan learning agility", lalu wawancara dan tes kesehatan dalam tes kemampuan bidang tiap BUMN. Mulai Jumat 7 Maret 2025; jenjang SMA/SMK, D3, D4, S1, S2 | Rilis Kementerian BUMN (tanggal terbit situs 2025-04-09, dateline 5 Maret 2025) — [bumn.go.id](https://bumn.go.id/publikasi/berita/rilis/detail/67f5f71ed5c0d). Resmi |
| Pembagian tahap 1 dan 2 (2025) | Online Test 1 = TKD, AKHLAK, Wawasan Kebangsaan; Online Test 2 = Bahasa Inggris dan Learning Agility | ANTARA 9 Mar 2025, mengutip FHCI — [antaranews](https://www.antaranews.com/berita/4698141/mengenal-jenis-materi-tes-seleksi-rekrutmen-bersama-bumn-2025). Media (kutipan FHCI). Konsisten dengan LKIP 2024 |
| Jadwal 2025 (jalur reguler) | Tes tahap 2 20-22 Mei 2025; hasil tahap 2 diumumkan 13 Juni 2025; tahap 3 (TKB, wawancara, MCU, analisis media sosial, digital mindset) Juni-Juli 2025; pengumuman final minggu kedua Juli 2025 | detikSulsel 13 Jun 2025 — [detik](https://www.detik.com/sulsel/berita/d-7962377/link-dan-cara-cek-pengumuman-hasil-tes-tahap-2-rekrutmen-bersama-bumn-2025), "dikutip dari laman FHCI BUMN". Media |
| Jalur disabilitas 2025 | Tes tahap 1 (TKD, Akhlak, TWK) dan tahap 2 (English, Learning Agility) di minggu kedua-ketiga April 2025, lalu TKB April-Mei | Kompas.com 9 Apr 2025 — [kompas](https://www.kompas.com/edu/read/2025/04/09/124842071/rekrutmen-bersama-bumn-2025-link-pengumuman-tahapan-dan-jenis-tes). Media |
| Data hasil | 2024: 108 perusahaan, kebutuhan 1.777, pelamar 842.584, lulus 1.378. 2023: 123 perusahaan, kebutuhan 2.223, pelamar 883.236, lulus 1.420 | LKIP 2024 hlm. PDF 33-34; LKIP 2023 hlm. PDF 46. Resmi |

**Fresh graduate vs experienced hire.** Sumber resmi hanya menyebut kedua jalur ada (rilis 2024). Saya tidak menemukan sumber yang menyebut komponen tes berbeda antar kedua jalur; perbedaan yang tampak di hasil pencarian hanya syarat (usia, IPK).
- **Asumsi (belum pasti):** komponen tes sama untuk fresh graduate dan experienced hire — tidak ada sumber yang membedakannya, tetapi tidak ada sumber resmi yang menyatakan sama. Cara cek: panduan peserta di portal RBB.
- Pengecualian per jenjang: ANTARA 9 Mar 2025 mengutip "bagi lulusan jenjang pendidikan SMA/sederajat, tidak perlu mengikuti tes bahasa Inggris" (media, kutipan FHCI). Sumber lain tidak sepakat soal pembagian AKHLAK/Learning Agility untuk SMA (hasil pencarian, belum dibuka); anggap belum jelas.

**RBB 2026: tidak ada format resmi.**
- Liputan6 Cek Fakta, terbit 9 Jun 2026, 19:01 WIB — [liputan6](https://www.liputan6.com/cek-fakta/read/7839160/cek-fakta-tidak-benar-link-pendaftaran-rekrutmen-bersama-bumn-2026): "Hingga Selasa 9 Juni 2026, belum ada informasi Rekrutmen Bersama BUMN 2026 di media sosial FHCI"; tautan pendaftaran RBB 2026 di Facebook (28 Mei 2026) dinyatakan tidak benar dan meminta data pribadi.
- Dealls (terbit 29 Apr 2026) — [dealls](https://dealls.com/pengembangan-karir/rekrutmen-bersama-bumn-2026): "dibuka sebentar lagi", jadwal "estimasi", tanpa perubahan format 2026.
- Pencarian di portal berita besar (detik, Kompas, ANTARA, CNBC, Kontan, Bisnis, CNN, Liputan6, kumparan) tidak mengembalikan berita RBB 2026 yang resmi sampai hasil terbaru yang muncul (Jun 2026). Satu situs kecil (Bungko) mengaku RBB 2026 dibuka 20 Mei-10 Juni 2026; halamannya 404 dan bertentangan dengan Liputan6, jadi tidak dipakai.
- **Asumsi (belum pasti):** per 18 Nov 2026 belum ada batch RBB 2026 — pencarian tidak mencakup berita Jul-Nov 2026 dengan baik dan Instagram FHCI tidak saya buka. Cara cek: `@fhci.bumn` dan portal FHCI langsung.
- **Asumsi (belum pasti):** peran Kementerian BUMN sudah diganti BP BUMN (UU 16/2025 menurut judul hasil pencarian NU Online, Jakarta Post, dan lain-lain; tidak saya buka). Artikel Dealls April 2026 masih menulis "Kementerian BUMN" sebagai penyelenggara. Tidak ada sumber yang menyebut Danantara menyelenggarakan RBB.

## 2. Detail per komponen (RBB 2025, kecuali dicatat lain)

**Tidak satu pun angka di tabel ini terbukti dari sumber resmi yang saya buka.** Rilis dan LKIP Kementerian BUMN tidak memuat jumlah soal, durasi, atau ambang batas (dikonfirmasi dari halaman rilis 2024 dan 2025; LKIP dibaca teksnya hlm. PDF 33-34 dan 45-46). Angka di bawah semuanya "media/bimbel", sebagian mengaku mengutip FHCI.

| Komponen | Jumlah soal / durasi | Ambang batas / bobot | Sumber (label) |
|---|---|---|---|
| TKD | 100 soal / 73 menit; 4 subtes @25 soal: Verbal Logical Reasoning (18 mnt), Number Sequence (27 mnt), Word Classification (9 mnt), Diagram Reasoning (20 mnt) (jumlah menit subtes = 74, bukan 73) | min 58; bobot 40% | Dealls 18 Apr 2025 — [dealls](https://dealls.com/pengembangan-karir/contoh-soal-tkd-bumn). Bimbel; mengutip materi RBB 2023 dan 2024. Total 100/73 juga di detikSulsel 27 Apr 2025 "dilansir dari laman resmi Rekrutmen Bersama BUMN 2025" — [detik](https://www.detik.com/sulsel/berita/d-7888082/nilai-ambang-batas-bumn-2025-berdasarkan-aturan-rbb-yuk-cek-di-sini). Media (mengaku dari laman resmi, tidak saya verifikasi) |
| AKHLAK | 90 soal / 30 menit | min 65; bobot 50% | Dealls (sama); detikSulsel 27 Apr 2025 (90/30, min 65). Media/bimbel |
| Wawasan Kebangsaan (TWK) | 10 soal; durasi 5 menit (Dealls) atau 10 menit (detikSulsel) | min 50; bobot 10% | Dealls; detikSulsel. **Bertentangan** pada durasi |
| English | 85 soal / 80 menit; min 450 poin; fokus reading, grammar, vocabulary | min 450 | ANTARA 24 Mei 2025, "tes dirancang FHCI" — [antaranews](https://www.antaranews.com/berita/4855385/ingat-ini-nilai-ambang-batas-untuk-lulus-tes-tahap-2-rbb-bumn-2025); IDN Times 21 Mei 2025 (85/80, "450 dari 677", setara TOEFL ITP — klaim IDN) — [idntimes](https://www.idntimes.com/life/career/nilai-maksimal-tbi-bumn-00-98756-h8rmv9). Media. Dealls dan Tirto (hasil pencarian, tidak dibuka) menyebut 140 soal/115 menit (2023) dan 100/90 (2024): format berubah antar tahun |
| Learning Agility | ANTARA dan IDN Times: 100 soal / 35 menit, min 16 poin (IDN: "16 dari 40"). detikSulsel: 100 soal / 40 menit, min 36. Dealls: 50 soal / 30 menit, min 197 dari 300; dimensi mental, change, people, results agility, self-awareness | min 16 / 36 / 197 | **Tiga set angka saling bertentangan.** ANTARA 24 Mei 2025; IDN 21 Mei 2025; detikSulsel 27 Apr 2025; Dealls 2 Mei 2025 — [dealls](https://dealls.com/pengembangan-karir/learning-agility-bumn) (halaman itu sendiri menyebut "2024" di pengantar). Media/bimbel |

**Aturan peringkat (2025):**
- Ambang tahap 1 per komponen (TKD 58, AKHLAK 65, TWK 50), bobot TKD 40 / AKHLAK 50 / TWK 10 untuk SMA/sederajat (bobot untuk non-SMA tidak tercantum); rasio 1:200 per posisi, nilai sama dengan peringkat terakhir tetap lolos — detikSumut 18 Mei 2025, mengutip Instagram `@fhci.bumn` — [detik](https://www.detik.com/sumut/berita/d-7919391/rbb-2025-umumkan-hasil-tes-online-tahap-1-simak-cara-cek-nilai-ambang-batasnya). Media (kutipan akun FHCI, postingan asli tidak dibuka).
- Tahap 2: "peserta yang memenuhi atau melampaui ambang batas diperingkat berdasarkan skor"; rumus skor dan skala tidak dijelaskan — ANTARA 24 Mei 2025. Bobot 60/40 antar tes online 1 dan 2 dan rasio 1:10 hanya dari ringkasan hasil pencarian, tidak ada di halaman ANTARA yang saya buka: abaikan.
- Aturan salah-benar (penalti jawaban salah) untuk semua komponen: **tidak ditemukan.**
- Ambang batas English: 450 (ANTARA/IDN); detikSulsel menulis "tidak disebut" untuk English. Ambang dan skor tiap komponen berubah antar sumber, jadi tidak boleh di-hardcode.

**Platform dan aturan tes (2025):** tes daring lewat Safe Exam Browser, webcam menyala, wajah terlihat, internet minimal 20 Mbps, Windows/macOS, Chrome disarankan, dilarang kalkulator, sempoa, ponsel, headset, dual monitor, penutup wajah — ANTARA 14 Apr 2025 (tidak mengutip FHCI secara langsung) — [antaranews](https://www.antaranews.com/berita/4768749/aturan-mengikuti-tes-online-rekrutmen-bersama-bumn-2025). Media. Konsekuensi untuk aplikasi: kalkulator dilarang, jadi soal numerik latihan harus bisa dihitung tanpa kalkulator.

**Konflik isi TKD.** Dua deskripsi isi TKD bertentangan:
- Kompas.com 9 Apr 2025 dan Suara.com 12 Apr 2025: TKD terdiri dari TWK, TIU (verbal, numerik, logika), dan TKP; Suara menulis sendiri bahwa isinya "berdasarkan informasi umum dari pelaksanaan sebelumnya" dan memperkirakan "100-120 soal" (bukan resmi) — [suara](https://www.suara.com/news/2025/04/12/112928/ini-kisi-kisi-rekrutmen-bersama-bumn-2025-tahap-1?page=all). Media.
- Dealls: empat subtes (verbal logical reasoning, number sequence, word classification, diagram reasoning). Bimbel.
- ANTARA 9 Mar 2025 (kutipan FHCI): TKD "akan menilai pengetahuan, keterampilan berpikir kritis dan analisis, dan pemecahan masalah" — hanya kemampuan, tanpa subtes.
- Tidak ada sumber resmi yang menyebut nama subtes. Jangan modelkan nama subtes Dealls sebagai fakta.

## 3. AKHLAK: teks resmi dan contoh soal

**Teks nilai dan definisi resmi terbit dan terbukti.** SE Menteri BUMN No. SE-7/MBU/07/2020, "Nilai-Nilai Utama (Core Values) Sumber Daya Manusia Badan Usaha Milik Negara", ditetapkan 1 Juli 2020 (Erick Thohir), 3 hlm. — [jdih.bumn.go.id PDF](https://jdih.bumn.go.id/storage/peraturan/SE-7-MBU-07-2020.pdf). Dibaca dari render gambar tiap halaman, Angka E.1 (hlm. 2-3). Enam nilai, tiap nilai punya definisi dan tiga panduan perilaku:

| Nilai | Definisi (verbatim) | Panduan perilaku (verbatim) |
|---|---|---|
| Amanah | Memegang teguh kepercayaan yang diberikan. | (1) Memenuhi janji dan komitmen. (2) Bertanggung jawab atas tugas, keputusan, dan tindakan yang dilakukan. (3) Berpegang teguh kepada nilai moral dan etika. |
| Kompeten | Terus belajar dan mengembangkan kapabilitas. | (1) Meningkatkan kompetensi diri untuk menjawab tantangan yang selalu berubah. (2) Membantu orang lain belajar. (3) Menyelesaikan tugas dengan kualitas terbaik. |
| Harmonis | Saling peduli dan menghargai perbedaan. | (1) Menghargai setiap orang apapun latar belakangnya. (2) Suka menolong orang lain. (3) Membangun lingkungan kerja yang kondusif. |
| Loyal | Berdedikasi dan mengutamakan kepentingan Bangsa dan Negara. | (1) Menjaga nama baik sesama karyawan, pimpinan, BUMN, dan Negara. (2) Rela berkorban untuk mencapai tujuan yang lebih besar. (3) Patuh kepada pimpinan sepanjang tidak bertentangan dengan hukum dan etika. |
| Adaptif | Terus berinovasi dan antusias dalam menggerakkan ataupun menghadapi perubahan. | (1) Cepat menyesuaikan diri untuk menjadi lebih baik. (2) Terus-menerus melakukan perbaikan mengikuti perkembangan teknologi. (3) Bertindak proaktif. |
| Kolaboratif | Membangun kerja sama yang sinergis. | (1) Memberi kesempatan kepada berbagai pihak untuk berkontribusi. (2) Terbuka dalam bekerja sama untuk menghasilkan nilai tambah. (3) Menggerakkan pemanfaatan berbagai sumber daya untuk tujuan bersama. |

- Angka E.2-E.3 (hlm. 3): wajib dipakai "secara utuh tanpa pengurangan atau penambahan" pada seluruh group BUMN.
- **Beda dari BerAKHLAK ASN.** SE ini 6 nilai (A-K-H-L-A-K); SE MenPANRB 20/2021 punya 7 nilai BerAKHLAK (ada Berorientasi Pelayanan) dengan perilaku berbeda (lihat `sumber_twk_tkp.md`). Bank TKP yang ada memakai SE 20/2021; modul AKHLAK BUMN wajib memakai SE-7/2020, jangan dicampur.
- Status berlaku SE-7/2020: **belum diverifikasi.** Hasil pencarian tidak menemukan pencabutan, tetapi halaman JDIH BUMN (daftar status) tidak saya buka. SK-115/MBU/05/2022 (AKHLAK Culture Journey) merujuk ke SE ini menurut ringkasan pencarian (tidak dibuka). Cara cek: halaman peraturan di jdih.bumn.go.id atau peraturan.bpk.go.id.
- **Asumsi (belum pasti):** tes AKHLAK RBB mengukur nilai-nilai ini; sumber media menyebut "karakter dan penerapan nilai" (ANTARA 9 Mar 2025) dan memakai urutan enam nilai yang sama, tetapi tidak ada sumber resmi yang menautkan item tes ke SE-7/2020.
- Kunci/bobot jawaban AKHLAK: **tidak ada pedoman resmi.** Statusnya sama dengan bobot TKP (sumber_twk_tkp.md): kunci buatan aplikasi, bukan kunci resmi, harus berlabel begitu.

**Contoh soal resmi: tidak ditemukan.** Pencarian "contoh soal resmi / kisi-kisi resmi FHCI" hanya mengembalikan kumpulan soal pihak ketiga (Dealls, Tirto, Mamikos, Sonora, Kitalulus, detik). Beberapa artikel menyebut "laman resmi" memuat materi, tetapi portalnya tidak bisa dibaca (lihat Lingkup). **Asumsi (belum pasti):** portal RBB 2025 mungkin memuat kisi-kisi atau tes uji coba resmi — kami hanya melihat media menyebut "uji coba" tanpa membukanya. Cara cek: buka portal dari peramban dengan sertifikat valid atau minta sesama peserta tangkapan layar.

## 4. Tumpang tindih dengan mesin SKD aplikasi

Generator bertipe di `api/internal/gen/`: analitis, aritmetika, deretangka, derethuruf, himpunan, kecepatan, kerjasama, perbandingan, persentase, ratarata, silogisme, umur (hasil `ls` direktori; tidak saya buka kodenya). Bank TWK/TKP di `api/internal/bank/data`. Pemetaan di bawah adalah penilaian saya atas deskripsi media (tidak resmi), bukan spesifikasi.

| Komponen RBB (deskripsi media) | Kemungkinan dipakai ulang |
|---|---|
| TKD "Number Sequence" / numerik | `deretangka`, `aritmetika`, `perbandingan`, `persentase`, `ratarata`, dll. Tumpang tindih tinggi. Kalkulator dilarang, jadi angka harus kecil |
| TKD "Verbal Logical Reasoning" / logika | `silogisme`, `analitis`, `himpunan`. Tumpang tindih tinggi |
| TKD "Word Classification" (klasifikasi kata) | Belum ada generator; butuh daftar kategori kata atau bank kurasi. Tidak ada sumber resmi bahwa subtes ini ada |
| TKD "Diagram Reasoning" / gambar | Butuh soal figural (SVG generator, PRD fase 3). Belum tersedia |
| Wawasan Kebangsaan | Sebagian dengan bank TWK (Pancasila, UUD 1945, Bhinneka Tunggal Ika); materi detail tidak resmi |
| AKHLAK (situasi kerja, pilihan ganda) | Struktur item TKP (situasi + opsi berbobot) bisa dipakai; isi dan bobot perlu bank baru berdasar SE-7/2020 |
| English (struktur, reading, listening?) | Tidak ada di aplikasi; butuh bank teks. Listening disebut Tribunnews (tidak dibuka); jangan dimodelkan |
| Learning Agility | Tidak ada di aplikasi; tidak ada format yang bisa dipastikan |

## 5. Hak cipta dan branding

- PRD/CLAUDE.md melarang nama atau logo BKN/KemenPANRB/FHCI dan mewajibkan soal orisinal. **FHCI adalah penyelenggara RBB**, jadi modul BUMN tidak boleh memakai nama atau logo FHCI, dan sebaiknya tidak menamai fitur "RBB resmi"/"simulasi RBB" yang menyiratkan afiliasi. Sebutan faktual "gaya tes seleksi BUMN" dengan disclaimer lebih aman.
- **Asumsi (belum pasti):** Kementerian BUMN/BP BUMN dan logo Garuda perlu diperlakukan sama seperti BKN/KemenPANRB (tidak dipakai) — larangan di PRD hanya menyebut tiga nama, perlu keputusan pemilik produk. SE-7/2020 sendiri berkepala Garuda dan nama Menteri BUMN; kutip isinya, jangan pakai kepala suratnya.
- Kutipan definisi AKHLAK: naskah SE-7/2020 (3 hlm., dibaca sebagai gambar) tidak memuat pernyataan hak cipta atau larangan penggandaan. Dasar yang dipakai untuk SE MenPANRB 20/2021 (UU 28/2014 Pasal 43 huruf b dan 44 ayat (1) huruf a, sumber disebut) kemungkinan berlaku sama, tetapi analisis itu belum diulang untuk SE-7/2020 dan bukan nasihat hukum. Tetap sebut sumber lengkap tiap kutipan (mis. "SE MBU-7/2020 angka E.1").
- Soal dari bimbel/Tirto/Dealls/media memiliki hak cipta pihak masing-masing: jangan disalin, dan jangan dipakai sebagai kunci atau dasar subtes. Sumber di atas hanya dipakai untuk fakta format, itupun bertanda tidak resmi.
- Soal asli peserta RBB yang bocor ("bocoran") tidak boleh dipakai. Suara.com sendiri menyatakan kisi-kisinya bukan dokumen resmi.

## Siap dispesifikasikan?

**Rekomendasi: belum untuk spesifikasi ujian lengkap. Hanya sebagian yang bisa dimodelkan, dengan status sumber di data.**

| Komponen | Bisa dimodelkan sekarang? | Alasan |
|---|---|---|
| Urutan tahapan dan daftar komponen (TKD, AKHLAK, Wawasan Kebangsaan, English, Learning Agility) | Ya, sebagai metadata berlabel "resmi, RBB 2024/2025" | Terbukti di rilis dan LKIP Kementerian BUMN |
| Nilai AKHLAK, definisi, panduan perilaku | Ya, sebagai data korpus verbatim | Terbukti dari SE-7/2020 (status berlaku belum dicek) |
| Jumlah soal, durasi, subtes TKD | Tidak sebagai angka baku | Hanya media/bimbel dan saling bertentangan |
| Ambang batas, bobot, aturan peringkat | Tidak sebagai angka baku | Hanya media; Learning Agility punya tiga set angka; berubah antar tahun |
| English dan Learning Agility | Tidak | Format berubah tiap tahun (140/115, 100/90, 85/80) dan tidak ada soal contoh resmi |
| Wawasan Kebangsaan | Sebagian | Topik umum jelas (Pancasila, UUD 1945), jumlah/durasi tidak resmi |
| Batch 2026 | Tidak | Tidak ada pengumuman resmi yang ditemukan |

Langkah yang saya sarankan sebelum menulis spesifikasi:
1. Buka portal peserta RBB 2025 dari peramban biasa dan simpan tangkapan layar kisi-kisi/FAQ resmi (inilah satu-satunya jalur ke angka resmi). Bila ada panduan peserta PDF, itu menggantikan semua angka media di atas.
2. Cek `@fhci.bumn` untuk RBB 2026 sebelum memutuskan target tahun.
3. Bila tetap ingin lanjut tanpa itu: modelkan spesifikasi sebagai "profil latihan" berversi (`rbb-latihan-v1`) yang jelas bukan spesifikasi resmi, dengan jumlah soal dan durasi yang bisa diubah, label "non-resmi" di UI, dan jangan menampilkan "lulus/tidak lulus" berdasarkan ambang media.
4. Prioritas konten yang aman: TKD numerik dan logika (pakai ulang generator), bank AKHLAK situasional dari SE-7/2020 dengan kunci berlabel buatan aplikasi. Tunda English, Learning Agility, subtes gambar dan klasifikasi kata.
