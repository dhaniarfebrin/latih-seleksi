# Format resmi Rekrutmen Bersama BUMN (RBB)

Dibaca langsung pada 18 Nov 2026. Tujuan: memutuskan apakah modul BUMN (PRD fase 2) sudah bisa dispesifikasikan sebagai data ujian berversi, sesuai invarian "spesifikasi ujian disimpan sebagai data berversi" (CLAUDE.md, PRD §7).

**Ringkasan:** tahapan dan daftar komponen tes RBB terbukti dari sumber resmi (Kementerian BUMN). Jumlah soal, durasi, dan subtes TIDAK ada di satu pun sumber resmi yang bisa dibuka (koreksi 19 Nov 2026: ambang, bobot, skala skor, dan aturan platform 2023 dan 2024 ternyata resmi, dibaca lewat arsip portal FHCI; lihat bagian 7); angka butir dan menit hanya dari media dan situs bimbel, dan saling bertentangan. Definisi nilai AKHLAK resmi terbit dan terbukti. Contoh soal resmi tidak ditemukan. Tambahan 19 Nov 2026: bagian 6 memverifikasi ulang RBB 2025 (benar berlangsung) dan merinci RBB 2024; kesimpulan tidak berubah, hanya ambang/bobot 2024 lebih konsisten.

## Lingkup pencarian dan batasnya
- Pencarian web (mode standar, sekali mode extended) dengan kata kunci RBB, TKD, AKHLAK, learning agility, bahasa Inggris, ambang batas, FHCI, BP BUMN, Danantara; dibuka: situs bumn.go.id (rilis dan dua LKIP PDF), jdih.bumn.go.id (SE-7/2020), ANTARA, detik, Kompas, Suara, IDN Times, Investortrust, Liputan6, Dealls.
- **[Dikoreksi di bagian 7: portal 2023 dan 2024 terbaca lewat arsip Wayback; 2025 hanya SPA. Paragraf ini menggambarkan kondisi sebelum itu.]** **Portal resmi peserta tidak bisa dibaca.** `rekrutmenbersama2025.fhcibumn.id`: WebFetch gagal ("certificate has expired"); percobaan saya mengambilnya dengan verifikasi TLS dimatikan ditolak oleh sistem izin dan tidak saya lanjutkan. `rekrutmenbersama.bumn.go.id`, `fhcibumn.id`, `fhcibumn.com`, `rekrutmenbersama2024.fhcibumn.id` tidak resolve dari lingkungan ini (DNS gagal, bisa jadi pembatasan lingkungan, belum tentu situs mati). Jadi kisi-kisi, FAQ, dan pengumuman di portal itu belum dibaca langsung.
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

**Tidak satu pun angka di tabel ini terbukti dari sumber resmi yang saya buka (RBB 2025). Koreksi 19 Nov 2026: untuk 2023 dan 2024 ambang, bobot, dan skala skor resmi (bagian 7.2); skala nilai berbeda dari bacaan media "dari jumlah soal".** Rilis dan LKIP Kementerian BUMN tidak memuat jumlah soal, durasi, atau ambang batas (dikonfirmasi dari halaman rilis 2024 dan 2025; LKIP dibaca teksnya hlm. PDF 33-34 dan 45-46). Angka di bawah semuanya "media/bimbel", sebagian mengaku mengutip FHCI.

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
- Status berlaku SE-7/2020: peraturan.bpk.go.id (dicek 2026-11-20, hasil pencarian "SE-7/MBU/07/2020") menampilkan entri "Berlaku mulai 6 tahun yang lalu", berlaku sejak 1 Juli 2020; halaman itu tidak memuat relasi mencabut, diubah, atau dicabut (tidak ada bagian "Status Peraturan"). Itu ketiadaan di halaman, bukan bukti tidak ada pencabut; halaman jdih.bumn.go.id tidak bisa dibuka (404). SK-115/MBU/05/2022 merujuk ke SE ini menurut ringkasan pencarian (tidak dibuka).
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

## 6. RBB 2024 (tambahan 19 Nov 2026)

Dibaca langsung pada 19 Nov 2026. Semua halaman di bawah saya buka sendiri kecuali bila ditandai. **Hasil utama: tahapan resmi 2024 sudah jelas; ambang batas, bobot, dan aturan peringkat 2024 konsisten di banyak media bertanggal 2024 yang mengutip FHCI; jumlah soal, durasi, dan subtes 2024 tetap tidak ada di sumber yang bisa dibuka.**

### 6.1 Verifikasi ulang: apakah RBB 2025 benar terjadi?

Tidak ada salah baca. Kutipan literal dari rilis Kementerian BUMN [67f5f71ed5c0d](https://bumn.go.id/publikasi/berita/rilis/detail/67f5f71ed5c0d) (dibuka ulang 19 Nov 2026):
- Judul: "Rekrutmen Bersama BUMN 2025 Resmi Dibuka, Lebih Inklusif Merangkul Disabilitas Hingga Semua Suku di Nusantara"
- "Tanggal Terbit : 2025-04-09 11:27:10" dan dateline "Jakarta, 05 Maret 2025 - Kementerian BUMN bersama dengan Forum Human Capital Indonesia (FHCI),": tanggal terbit situs (9 Apr) berbeda dari dateline (5 Mar); rilis ini kemungkinan diunggah belakangan ke situs, saya tidak menemukan penjelasannya.
- "RBB 2025 dimulai pada Jumat, 7 Maret 2025 dan melibatkan beberapa tahapan seleksi"
- "Tahapan seleksi mencakup registrasi, seleksi administrasi, tes daring (termasuk tes kemampuan dasar," (kalimat dilanjutkan dengan tes lain, wawancara, tes kesehatan; WebFetch memotong kutipan panjang).
- "Peserta juga diharapkan lebih teliti dalam memperhatikan periodisasi setiap tahapan" (jadwal berbeda tiap jalur).

Rilis itu sendiri hanya membuktikan pengumuman pembukaan. Bukti bahwa tesnya benar berlangsung berasal dari berita independen bertanggal:
- ANTARA, "Sabtu, 24 Mei 2025 09:50 WIB": "tes tahap 2 yang berlangsung pada 20-22 Mei 2025" — [antaranews](https://www.antaranews.com/berita/4855385/ingat-ini-nilai-ambang-batas-untuk-lulus-tes-tahap-2-rbb-bumn-2025).
- detikJatim, "Senin, 21 Apr 2025 15:45 WIB": proses "sudah memasuki tes tahap 1", jendela tes tahap 1 14-20 April 2025 — [detik](https://www.detik.com/jatim/berita/d-7878653/persiapan-untuk-tahap-selanjutnya-rekrutmen-bersama-bumn-2025).
- Sebelumnya sudah dibuka: detikSulsel 13 Jun 2025 (hasil tahap 2) dan detikSumut 18 Mei 2025 (hasil tahap 1), lihat bagian 1 dan 2.

Kesimpulan: RBB 2025 terbukti berlangsung (pengumuman resmi + tiga outlet bertanggal 2025 yang melaporkan tes dan hasil). Yang belum terbukti dari sumber resmi hanyalah angka format, sama seperti 2024. Batas: portal resmi tidak bisa dibuka, jadi tidak ada konfirmasi dari FHCI langsung.

### 6.2 Tahapan dan jadwal RBB 2024

| Hal | Isi | Sumber (label) |
|---|---|---|
| Tahapan dan komponen | Lihat bagian 1: Tahap 1 (TKD, AKHLAK, Wawasan Kebangsaan), Tahap 2 (Bahasa Inggris, Learning Agility), TKB di BUMN, pengumuman final; 23 Mar-11 Jul 2024 | LKIP 2024 hlm. PDF 33. **Resmi** |
| Pendaftaran dan seleksi administrasi | 23 Maret-1 April 2024; hasil seleksi administrasi 23 April 2024 | detikEdu 23 Mar 2024 (angka itu); Kontan 23 Apr 2024 15:11 WIB — [kontan](https://industri.kontan.co.id/news/klik-rekrutmenbersama2024fhcibumnid-hasil-seleksi-rekrutmen-bumn-2024-diumumkan). Media |
| Jadwal setelah itu | Uji coba 24-26 April; Tes Online 1 27 April-4 Mei; Tes Online 2 **25-27 Mei**; tes seleksi BUMN 11-30 Juni; pengumuman final Juli 2024 | Kontan 23 Apr 2024 (sama). Media |
| Tes Online 2, tanggal lain | **28-29 Mei 2024** (hasil tes online 1 diumumkan Rabu 22 Mei 2024) | Liputan6 (diperbarui 23 Mei 2024), "situs resminya" — [liputan6](https://www.liputan6.com/regional/read/5602874/simak-jadwal-tes-online-2-rekrutmen-bersama-bumn-2024); CNN Indonesia 27 Mei 2024: "digelar selama dua hari pada 28 hingga 29 Mei 2024" — [cnn](https://www.cnnindonesia.com/edukasi/20240527172425-561-1102539/contoh-soal-tes-bahasa-inggris-bumn-2024-lengkap-kunci-jawaban). Media. **Bertentangan dengan Kontan.** **Asumsi (belum pasti):** jadwal diundur dari 25-27 ke 28-29 Mei; tidak ada sumber yang menyatakan pengunduran. Cara cek: pengumuman di portal 2024 |
| Tes ulang tahap 2 | Ada tes ulang English dan Learning Agility (Juni 2024); "Nilai tes diambil yang tertinggi", tetapi bila skor ulang lebih rendah dari percobaan pertama, dipakai skor pertama; lolos ditentukan "ranking dari nilai tahap 1 digabung dengan nilai tahap 2" | Kompas.tv 8 Jun 2024 mengutip Lieke Roosdianti Pangkey, Direktur Eksekutif FHCI (Kompas.com 7 Jun 2024) — [kompas.tv](https://www.kompas.tv/lifestyle/513721/tes-ulang-rekrutmen-bersama-bumn-2024-dapat-nilai-rendah-sistem-pakai-skor-yang-mana). Media (kutipan FHCI) |
| Jumlah perusahaan/kebutuhan | LKIP: 108 perusahaan, 1.777 kebutuhan. Kontan: "110 BUMN" dan "lebih dari 1.800 lowongan"; detikEdu: 1.830 posisi dari 688 lowongan | LKIP resmi; Kontan, detikEdu media. Angka media berbeda dari LKIP; pakai LKIP |

**Komponen berbeda per jenjang (penting).** Pada 2024 pembagian tahap berbeda untuk SMA dan D3-S2, sehingga ringkasan LKIP ("Tahap 1 = TKD, AKHLAK, Wawasan Kebangsaan") hanya cocok untuk D3-S2:
- D3-S2: Tahap 1 = TKD (bobot 40%, min 58), AKHLAK (50%, min 65), Wawasan Kebangsaan (10%, min 50); Tahap 2 = English (80%, min 450) + Learning Agility (20%, min 197). Gabungan: Tahap 1 60% + Tahap 2 40%.
- SMA/sederajat: Tahap 1 = TKD (min 58); Tahap 2 = AKHLAK (60%, min 65), Wawasan Kebangsaan (10%, min 50), Learning Agility (30%, min 197). Tanpa English. Gabungan: Tahap 1 40% + Tahap 2 60%.
- Sumber: Kompas.com 27 Mei 2024 mengutip Direktur Eksekutif FHCI (bobot dan ambang tahap 2 per jenjang) — [kompas](https://www.kompas.com/tren/read/2024/05/27/190000865/ramai-di-media-sosial-bagaimana-penilaian-tes-learning-agility-rekrutmen?page=all); IDN Times 30 Mar 2024 23:03 WIB (semua bobot dan bobot gabungan 60/40 dan 40/60) — [idntimes](https://www.idntimes.com/life/career/passing-grade-tes-rekrutmen-bersama-bumn-2024-00-98756-m9sw6f); detikJogja 24 Apr 2024 (bobot tahap 1, menyebut dokumen "Rincian Detail Tahapan Seleksi RBB 2024" sebagai dasar) — [detik](https://www.detik.com/jogja/bisnis/d-7308262/materi-tes-online-tahap-1-rekrutmen-bersama-bumn-2024-dan-passing-gradenya); Liputan6 23 Mei 2024 (komponen tahap 2 per jenjang, "situs resminya"); Kompas.tv 8 Jun 2024 (bobot gabungan 40/60 dan 60/40 dari FHCI). Media; label "mengutip FHCI/situs resmi" tidak saya verifikasi di sumbernya.
- Konflik kecil: detikJogja menulis dua hal berbeda untuk SMA di tahap 1 (bagian ambang: hanya TKD; bagian jadwal: "tes kemampuan dasar dan wawasan kebangsaan"). Kompas 2024 dan Liputan6 menaruh Wawasan Kebangsaan SMA di tahap 2, jadi versi "hanya TKD" lebih konsisten. Bobot AKHLAK/TWK SMA tahap 1 tidak diumumkan.
- Konteks: 2025 menaruh TKD, AKHLAK, TWK di satu tahap untuk semua jenjang (bagian 1-2), jadi struktur per jenjang bisa berubah antar tahun.

### 6.3 Angka per komponen, RBB 2024

| Komponen | Soal / durasi | Ambang / bobot | Sumber (label) |
|---|---|---|---|
| TKD | "100 soal" (IDN Times: "Minimal 58 dari 100 jumlah soal yang tersedia"; ambigu: bisa jumlah soal atau skala nilai). **Durasi dan subtes: tidak ditemukan di sumber 2024 yang saya buka** | min 58; bobot 40% (D3-S2) | IDN Times 30 Mar 2024; CNN Indonesia 30 Apr 2024 06:00 WIB menyebut "TKD minimal 58" dan mengaitkannya dengan FHCI — [cnn](https://www.cnnindonesia.com/edukasi/20240429134845-561-1091724/batas-nilai-lulus-minimal). Media |
| AKHLAK | "90" ("65 dari 90"); durasi tidak ditemukan | min 65; bobot 50% (D3-S2), 60% (SMA, tahap 2) | IDN Times 30 Mar 2024; CNN 30 Apr 2024 (hanya min 65); Kompas 27 Mei 2024 (SMA). Media |
| Wawasan Kebangsaan | Jumlah/durasi tidak ditemukan | min 50; bobot 10% (kedua jalur) | IDN Times; CNN 30 Apr 2024; Kompas 27 Mei 2024. Media |
| English | Jumlah/durasi 2024 **tidak ditemukan di sumber yang bisa dibuka** (CNN 27 Mei 2024 dan Kompas 27 Mei 2024 tidak menyebutnya) | min 450; bobot 80% (D3-S2) | Kompas 27 Mei 2024 (kutipan FHCI); IDN Times. Media. Hasil pencarian (tidak saya buka, Tirto 403) menyebut 100 soal/90 menit untuk 2024 dan 140/115 untuk 2023; satu hasil lain menyebut 85/80 tanpa tahun jelas. **Asumsi (belum pasti):** 85/80 (ANTARA/IDN 2025) tidak berlaku untuk 2024 |
| Learning Agility | "50 soal berisi beberapa opsi jawaban", 30 menit, skor maksimal 300, min 197; bobot 20% (D3-S2), 30% (SMA) | min 197 dari 300 | Kompas.com 27 Mei 2024 (artikel yang menampilkan wawancara Direktur Eksekutif FHCI; 50 soal/30 menit/300 ada di artikelnya, tidak jelas apakah semuanya ucapan FHCI; "tidak dijelaskan apakah setiap opsi punya bobot berbeda"). Media |

**Temuan penting:** angka Learning Agility "50 soal / 30 menit / 197 dari 300" yang di bagian 2 saya tandai sebagai milik Dealls "2025" ternyata sama dengan angka 2024 di Kompas 27 Mei 2024. **Asumsi (belum pasti):** Dealls menyalin format 2024; format 2025 (100 soal/35 menit/min 16 menurut ANTARA dan IDN 2025) kemungkinan memang berbeda skala. Cara cek: panduan peserta 2025 atau pengumuman FHCI.

**Aturan peringkat 2024:** rasio "1:10", yaitu untuk tiap posisi diambil sepuluh peringkat teratas yang di atas standar minimum (Kompas.tv 8 Jun 2024, kutipan FHCI). Bandingkan 2025: rasio 1:200 untuk tahap 1 (detikSumut 18 Mei 2025): bukan angka yang sama dan bukan kontradiksi yang terbukti, karena yang satu berlaku sesudah gabungan tahap 1+2 dan yang lain untuk lolos tahap 1; **Asumsi (belum pasti):** ini dua aturan berbeda.

**Platform 2024:** tidak ditemukan sumber 2024 yang menyebut Safe Exam Browser. Satu unggahan Scribd "Petunjuk Online Test Rekrutmen Bersama Bumn 2024" tidak dapat dibaca penuh (hanya ringkasan otomatis): "peserta wajib menggunakan PC/laptop dengan kamera" dan "trial test disarankan"; pengunggah adalah pengguna Scribd, bukan penerbit resmi. Tidak resmi, kontennya tidak terverifikasi.

**Fresh graduate vs experienced hire 2024:** detikEdu 23 Mar 2024 menulis lowongan "fresh graduate sampai experienced hire" dan tidak menyebut tes berbeda; semua sumber 2024 yang saya buka membedakan komponen berdasarkan jenjang pendidikan (SMA vs D3-S2), bukan berdasarkan pengalaman. **Asumsi (belum pasti):** tes sama untuk fresh graduate dan experienced hire pada jenjang yang sama; tidak ada sumber yang membedakannya, tetapi juga tidak ada yang menyatakan sama.

**Dokumen resmi yang disebut media tetapi tidak berhasil saya buka:** "Rincian Detail Tahapan Seleksi RBB 2024" dan "Petunjuk Tes Online Untuk Peserta Rekrutmen Bersama BUMN 2024" (detikJogja 24 Apr 2024; dikatakan terlampir di rekrutmenbersama2024.fhcibumn.id). Pencarian tidak menemukan salinan PDF resmi. [Diperbarui di 7.1: dokumen "Rincian Detail" 2024 berhasil dibuka lewat arsip memakai `curl` dengan TLS normal; WebFetch memang menolak `web.archive.org`. Isinya tidak memuat jumlah soal/durasi.]

**Kenapa 2024 lebih berguna dari 2025 untuk spesifikasi:** ambang dan bobot 2024 konsisten antara IDN Times (30 Mar), detikJogja (24 Apr), CNN (30 Apr), dan Kompas (27 Mei) tanpa angka bertentangan, berbeda dari Learning Agility 2025 yang punya tiga set angka. Tetapi semuanya tetap media, dan belum ada angka durasi atau subtes 2024.

## 7. Pass mendalam lintas RBB 2023, 2024, 2025 (tambahan 19 Nov 2026)

Dibaca langsung pada 19 Nov 2026. Bagian ini **mengoreksi** pernyataan di bagian Lingkup, 1, 2, dan 6 yang menyebut "tidak ada sumber resmi" untuk ambang, bobot, dan aturan tes: ada, lewat arsip.

### 7.1 Jalur ke dokumen resmi (dan batasnya)

- WebFetch menolak `web.archive.org` ("unable to fetch from web.archive.org"; disebut sekali). Jalur yang berhasil, dengan verifikasi TLS normal (tanpa `-k`, tanpa melewati sertifikat): `curl` biasa ke `https://web.archive.org/web/<timestamp>if_/<url>` dari Bash, ditambah `https://archive.org/wayback/available?url=...&timestamp=...` dan indeks CDX `https://web.archive.org/cdx/search/cdx`.
- Salinan resmi FHCI yang dibaca (semua **resmi**, diambil dari situs penyelenggara lewat arsip):
  - Portal 2024, snapshot 2 Jun 2024 (`rekrutmenbersama2024.fhcibumn.id`): beranda (jadwal, syarat usia dan IPK) dan FAQ.
  - Enam PDF "RINCIAN DETAIL TAHAPAN SELEKSI RBB 2024" di `/assets/frontend/alur_pendaftaran/` (snapshot 24 Mar 2024; lima dibaca: Tes Online 1, Tes Online 2, TKB, Tahap Regis, Pengumuman; Inaugurasi tidak diunduh).
  - Portal 2023, beranda snapshot 5 Jul 2023 dan FAQ snapshot 19 Okt 2023 (`rekrutmenbersama.fhcibumn.id`); PDF "Petunjuk_Online_Test_Peserta.pdf" (snapshot 19 Okt 2023, 7 hlm.; **judul di dalamnya "REKRUT BERSAMA BUMN 2022"**, jadi sisa dokumen 2022 yang dipakai ulang di portal 2023).
  - 2025: portal `rekrutmenbersama2025.fhcibumn.id` adalah SPA (HTML 4,6 KB); isi diambil dari `api-rbb.fhcibumn.id` dan indeks CDX hanya menyimpan 4 endpoint API (`get-config` = `{"is_active": 1}`, opsi register, cek agreement/privacy), tanpa isi tahapan atau FAQ. Satu PDF resmi 2025 tersimpan: "Tutorial Instalasi SEB.pdf" (17 hlm., snapshot 15 Apr 2025, dari `storage.googleapis.com/rbb_wr_files/assets/default/`).
- **Dokumen resmi 2023, 2024, dan 2025 yang saya baca tidak memuat jumlah soal, durasi, nama subtes, atau aturan penalti jawaban salah.** FAQ 2023 dan 2024 menyatakan sendiri: "Peserta akan diberikan instruksi pengerjaan soal dan contoh soal pada masing-masing subtes di bagian awal sebelum mengerjakan soal yang sebenarnya", yaitu angka itu hanya tampil di dalam tes, tidak diterbitkan. Frasa "pada masing-masing subtes" dan "setiap bagian soal harus selalu diiringi timer" (Petunjuk 2023, no. 20) menunjukkan tes dibagi subtes ber-timer, tanpa angkanya.
- Tidak berhasil/ditolak: Tirto dan Bisnis (HTTP 403), Tribunnews (halaman kosong), mirror `id.obnews.co` (sertifikat kedaluwarsa; tidak dilewati), Scribd (tidak terbaca), Instagram `@fhci.bumn` (tidak ada alat). Cermin arsip tidak tersedia untuk isi API 2025.

### 7.2 Fakta resmi per batch

| Fakta | 2023 | 2024 | 2025 | Sumber |
|---|---|---|---|---|
| Tahap 1, D3-S2 | TKD + AKHLAK; TKD min 58 (40%), AKHLAK min 65 (60%) | TKD min 58 (40%), AKHLAK min 65 (50%), Wawasan Kebangsaan min 50 (10%) | (tidak ada dokumen resmi terbaca) | 2023: FAQ arsip 19 Okt 2023; 2024: PDF "Tes Online 1" + FAQ arsip 2 Jun 2024. Resmi |
| Tahap 1, SMA | TKD min 58 | TKD min 58 (isi "pengetahuan dasar dan umum") | - | idem. Resmi |
| Tahap 2, D3-S2 | English min 450; bobot FAQ menulis "TKD & AKHLAK (60%) & Bahasa Inggris (40%)" | English min 450 (80%), Learning Agility min 197 (20%) | - | idem. Resmi |
| Tahap 2, SMA | AKHLAK min 65; FAQ menulis "TKD (40%) & AKHLAK (60%)" | AKHLAK min 65 (60%), Wawasan Kebangsaan min 50 (10%), Learning Agility min 197 (30%) | - | idem. Resmi |
| Bobot gabungan tahap 1 / tahap 2 | tidak dinyatakan | SMA 40/60; D3-S2 60/40 | - | PDF "Tes Online 2". Resmi |
| Skala nilai | FAQ: "65 dari jumlah soal yang tersedia", "58 dari jumlah soal yang tersedia" (ambigu); English "setara 450 (Toefl)"; konversi "Nilai bahasa inggris / 675 x 100" | FAQ: "AKHLAK 65 dari 100", "TKD 58 dari 100", "Wawasan Kebangsaan 50 dari 100", English "minimal setara 450 (Toefl ITP) dari 677", "Learning Agility 197 dari 300"; konversi "/ 675 x 100" | - | FAQ arsip. Resmi. Catatan: FAQ 2024 menulis 677 dan 675 dalam halaman yang sama |
| Peringkat | "Ranking sesuai kebutuhan BUMN terdaftar" | sama (rasio 1:10 hanya dari Kompas.tv, kutipan FHCI) | rasio 1:200 hanya media | FAQ resmi; media |
| Jadwal tes online | Tahap 1 12-20 Jun; tahap 2 16-20 Jul; TKB 5-26 Agu 2023 | Tahap 1 27 Apr-4 Mei; **tahap 2 28-29 Mei 2024**; TKB di BUMN 11-30 Jun; final Juli. Pendaftaran 23 Mar-1 Apr | tahap 2 20-22 Mei (ANTARA, media) | Beranda arsip. Resmi untuk 2023 dan 2024 |
| Peserta per jenjang | SMA/sederajat, D3, D4/S1, S2 | idem; usia maks Fresh Graduate: SMA/SMK 25, D3 27, S1/D4 30, S2 35; Experience Hire 45 (semua jenjang); IPK min SMA 75, D3/D4/S1 3,00, S2 3,25 | - | Beranda arsip 2024. Resmi |
| Platform | PC/laptop + webcam + mouse; Chrome/Firefox terbaru; Core i3, RAM 4 GB; ponsel/tablet dilarang; registrasi wajah; kamera menyala sepanjang tes; jawaban tersimpan; tombol "Selesai Test" hanya di soal terakhir. **SEB tidak disebut** | PC/laptop + webcam + mouse; **Safe Exam Browser wajib**; Core i3, RAM 4 GB (6 GB disarankan), HDD/SSD 500 GB; Windows/Linux/Mac; registrasi wajah; tanpa earphone; kamera mati 3 menit = tes berhenti dan peserta gugur; dilarang tab lain, browsing, dual monitor, perangkat eksternal, kalkulator/sempoa; login ulang maks 3x; jadwal tidak bisa diubah; "kecurangan auto gugur dan blacklist"; Prev/Next boleh selama waktu ada; timer jalan terus kalau login terlambat; waktu habis = submit otomatis; soal kosong tetap kosong | Tutorial SEB resmi: unduh SEB dari "Lamaran Saya > Jadwal Tes Online" (Windows/macOS); SEB tidak jalan bila Chrome Remote Desktop terpasang atau Zoom aktif; layar terkunci bila tersambung ke layar kedua. Aturan lain 2025 hanya dari media (Kompas.tv 20 Apr 2025: kamera menyala, tanpa topi/masker/headset/kalkulator, "AI proctoring", curang = gugur) | 2023: Petunjuk Online Test arsip + FAQ; 2024: FAQ + PDF Tes Online 1; 2025: Tutorial SEB. Resmi (2025 aturan lain: media) |
| Beda fresh graduate vs experienced hire | Tidak ada beda di dokumen terbaca | **Hanya batas usia** (lihat di atas) dan syarat dokumen; PDF Tes Online 1 dan 2 membedakan komponen hanya menurut jenjang pendidikan (SMA vs D3-S2), bukan jalur pengalaman | tidak ada dokumen resmi | Beranda dan PDF arsip 2024. Resmi |

Koreksi atas pernyataan sebelumnya di catatan ini:
- Ambang, bobot, dan aturan platform 2023 dan 2024 **resmi** (bukan hanya media). Pernyataan "tidak satu pun angka resmi" di bagian 2 hanya benar untuk jumlah soal, durasi, subtes, dan penalti.
- Pembacaan media "65 dari 90 soal" (IDN Times 30 Mar 2024) bertentangan dengan FAQ resmi 2024 yang menulis "AKHLAK 65 dari 100". Nilai ambang adalah **skor berskala**, bukan jumlah jawaban benar mentah: TKD dan AKHLAK dari 100, English dari 677 (dinormalkan ke 100 lewat /675), Learning Agility dari 300. **Asumsi (belum pasti):** 90 dan 100 adalah jumlah butir (bagian 2025/Dealls), skala 100 adalah hasil konversi; hubungan butir ke skor tidak dinyatakan resmi. Cara cek: layar instruksi/tes uji coba peserta.
- Tahap 2 2024: **28-29 Mei 2024** resmi; 25-27 Mei (Kontan) adalah jadwal awal. Bagian 6.2 yang menulis "Asumsi: jadwal diundur" sekarang terjawab oleh beranda arsip (jadwal akhir 28-29 Mei).
- Portal 2024 tidak "tidak bisa dibaca": bisa lewat arsip (7.1). Dokumen "Rincian Detail Tahapan Seleksi RBB 2024" yang di bagian 6 disebut tidak berhasil dibuka sudah dibuka; isinya tidak memuat jumlah soal/durasi. "Petunjuk Tes Online ... 2024" tidak ditemukan sebagai PDF; yang ada hanyalah FAQ 2024 dan Petunjuk (judul 2022) di portal 2023.
- Pembagian tahap 2024 yang di bagian 6.2 dari Kompas/IDN/Liputan6 dikonfirmasi PDF resmi (Tes Online 1 dan 2).

### 7.3 Angka yang hanya dari media, per batch (dengan hitungan sumber independen)

- **2023.**
  - RBG.id 12 Jun 2023 16:55 WIB — [rbg.id](https://www.rbg.id/ekonomi/9449114428/tes-rekrutmen-bersama-bumn-2023-jadwal-bobot-nilai-soal-waktu-ujian-jumlah-soal-dan-pengertian-tkd): TKD "100 soal dengan durasi 120 menit" (bobot 40%); AKHLAK "54 soal dengan durasi waktu 120 sudah digabung dengan TKD" (60%); English "lebih dari 450". Tidak mengutip sumber resmi; **satu media, bertentangan** dengan 90 soal AKHLAK di 2025 (format bisa berbeda antar tahun; tidak terbukti).
  - detikEdu 13 Jun 2023 14:30 WIB — [detik](https://www.detik.com/edu/detikpedia/d-6770130/nilai-minimal-kelulusan-tes-online-rekrutmen-bumn-2023-catat-ya) dan CNN Indonesia 14 Jun 2023 08:00 WIB — [cnn](https://www.cnnindonesia.com/edukasi/20230613172108-561-961353/passing-grade-rekrutmen-bersama-bumn-2023-cek-dulu-biar-lolos): hanya ambang 58/65/450 dan jadwal; **sesuai** FAQ resmi 2023 (kata "dari jumlah soal yang tersedia" disalin dari FAQ). Tanpa jumlah soal atau durasi.
  - English 140 soal / 115 menit (Listening 50/35 mnt, Structure 40/25, Reading 50/55) — Cakap, terbit 9 Mei 2023 (sebelum tes), diperbarui 21 Jun 2024, tanpa menyebut tahun RBB — [cakap](https://blog.cakap.com/tes-bahasa-inggris-bumn/). Bimbel, tidak resmi. kumparan 19 Mei 2025 menulis angka yang sama untuk RBB 2025 — [kumparan](https://kumparan.com/ragam-info/tes-bahasa-inggris-bumn-berapa-soal-peserta-tes-wajib-tahu-255ynsQ6TXG): **salinan, bukan sumber independen**, dan bertentangan dengan ANTARA 2025 (85/80). Tirto 2023 140/115 hanya dari ringkasan pencarian (halaman 403, tidak dibuka).
- **2024.**
  - IDN Times 30 Mar 2024 — [idntimes](https://www.idntimes.com/life/career/passing-grade-tes-rekrutmen-bersama-bumn-2024-00-98756-m9sw6f): "AKHLAK minimal 65 dari 90", "TKD minimal 58 dari 100", English 450, TWK 50, LA 197; **tidak ada durasi**. Satu media; "90" bertentangan dengan FAQ resmi "dari 100".
  - Learning Agility 50 soal / 30 menit / 300 / min 197 — Kompas.com 27 Mei 2024 (bagian 6.3). Min 197 dari 300 **dikonfirmasi resmi**; 50/30 satu media.
  - English 2024: "100 soal / 90 menit" hanya dari ringkasan pencarian yang menyebut Tirto dan KitaLulus (tidak dibuka). detik 22 Mar 2024 — [detik](https://finance.detik.com/bumn/d-7255512/sempat-dikeluhkan-tes-bahasa-inggris-rekrutmen-bumn-bakal-dievaluasi): Ketua FHCI Alexandra Askandar: tes bahasa Inggris "akan kita lihat kembali, tapi bukan berarti kita akan menurunkan standar"; Sekjen FHCI Dharma Syahputra: "banyak sekali yang menyatakan terlalu sulit sekali" — format 2023 dievaluasi, tanpa angka.
  - Subtes TKD 2024: tidak ada sumber yang saya buka. Dealls (18 Apr 2025) menulis TWK "subtes yang baru berlaku sejak RBB 2024" (konsisten dengan rilis resmi 2024).
- **2025.**
  - Dealls 18 Apr 2025 (dibuka) — [dealls](https://dealls.com/pengembangan-karir/contoh-soal-tkd-bumn): TKD 100 soal / 73 menit; empat subtes @25 soal: Verbal Logical Reasoning 18 mnt, Number Sequence 27 mnt, Word Classification 9 mnt, Diagram Reasoning 20 mnt (jumlah 74 mnt, bukan 73); AKHLAK 90 soal / 30 menit, min 65, bobot 50%; TWK "sekitar 10 soal", 5 menit, min 50, bobot 10%; min TKD 58 bobot 40%. Tanpa penalti salah. Sumber yang ia sebut: Katadata (materi 2023), Okezone (soal 2024), situs RBB 2024: **tidak ada yang mendukung angka 2025**, jadi asal angka tidak jelas.
  - detikSulsel 27 Apr 2025 (angka sama, mengaku dari laman resmi RBB 2025; TWK 10 menit): **sesuai Dealls pada total, bertentangan pada durasi TWK**. Dealls dan detikSulsel mungkin berbagi sumber yang sama; saya hitung sebagai dua laporan, bukan dua sumber independen.
  - Nama subtes di luar Dealls: ringkasan pencarian menunjukkan judul Tirto ("20 contoh soal word classification bumn 2025", "20 contoh soal number sequence bumn") dan Mamikos ("soal number sequence rekrutmen BUMN" bertahun 2024) yang memakai nama subtes sama; hanya judul, tidak dibuka, dan tidak memberi jumlah soal. **Hasil: nama subtes didukung >1 penerbit (judul), jumlah per subtes hanya Dealls.**
  - Kompas.tv 20 Apr 2025 12:06 WIB — [kompas.tv](https://www.kompas.tv/ekonomi/588053/simak-berikut-aturan-tes-online-rekrutmen-bersama-bumn-2025-dan-passing-grade-twk-tkd) (kutipan Direktur FHCI Lieke Roosdianti ke Kompas.com 17 Apr 2025): "peserta wajib memenuhi passing grade dari seluruh tes", periode tes online 19-28 April 2025 (**bertentangan** dengan detikJatim 21 Apr 2025: tahap 1 14-20 April; bagian 6.1), proctoring AI, kamera menyala, dilarang topi/masker/headset/kalkulator. Suara.com 12 Apr 2025 menyebut sumber resmi hanya `rekrutmenbersama2025.fhcibumn.id`, `fhcibumn.com`, Instagram `@kementerianbumn` dan `@fhci.bumn`, dan tidak memuat jumlah/durasi.
  - Mamikos (ringkasan pencarian, tidak dibuka): "AKHLAK 72% (65 dari 90), TKD 58% (58 dari 100), TWK 50 dari 100, LA 66% (197 dari 300)": tidak diverifikasi.
- **Aturan penilaian.** Penalti jawaban salah: **tidak ada sumber untuk tahun mana pun** (resmi maupun media; Dealls, IDN, CNN, Kompas, ANTARA tidak menyebutnya). Satu-satunya aturan resmi: soal yang tidak dijawab tetap kosong dan waktu habis = submit otomatis (FAQ 2024). Skor AKHLAK per opsi: hanya klaim bimbel (opsi terbaik bernilai tertinggi) dari situs bimbel yang tidak saya anggap andal; **tidak resmi**.

### 7.4 Tabel perbandingan tunggal

Kolom keyakinan: **resmi** (dokumen FHCI arsip), **media konsisten** (>=2 media bertanggal pada masa itu, tanpa angka bertentangan; independensi belum terbukti), **satu media**, **bertentangan**.

| Komponen / hal | 2023 | 2024 | 2025 | Keyakinan |
|---|---|---|---|---|
| Struktur tahap per jenjang | T1: TKD+AKHLAK (SMA: TKD); T2: English (SMA: AKHLAK) | T1: TKD, AKHLAK, TWK (SMA: TKD); T2: English, LA (SMA: AKHLAK, TWK, LA) | 5 komponen (TKD, AKHLAK, TWK, English, LA); T1 = TKD+AKHLAK+TWK, T2 = English+LA (ANTARA 9 Mar 2025) | 2023, 2024: **resmi**. 2025: rilis resmi menyebut komponen; pembagian tahap per jenjang media |
| Ambang TKD / AKHLAK / TWK | 58 / 65 / - | 58 / 65 / 50 (skala 100) | 58 / 65 / 50 | 2023, 2024: **resmi**. 2025: media konsisten (Dealls, detikSulsel, IDN, Kompas.tv menyebut ada) |
| Ambang English / LA | 450 / - | 450 (dari 677) / 197 dari 300 | 450 / LA bertentangan (16, 36, 197) | 2023, 2024: **resmi**. 2025 English media konsisten; 2025 LA **bertentangan** |
| Bobot | T1 D3-S2 TKD 40 / AKHLAK 60; T2 D3-S2 English (gabungan 60 : 40, lihat catatan di bawah tabel) | T1 D3-S2 40/50/10; T2 D3-S2 English 80, LA 20; SMA T2 60/10/30; gabungan 60/40 dan 40/60 | TKD 40 / AKHLAK 50 / TWK 10 (Dealls, detikSumut); lain tidak ada | 2023, 2024: **resmi** (teks FAQ 2023 untuk T2 ambigu); 2025: media konsisten untuk T1 |
| TKD: jumlah soal / durasi | 100 / 120 menit (RBG) | "100" (IDN, ambigu); durasi tidak ada | 100 / 73 menit (Dealls, detikSulsel) | 2023 satu media; 2024 satu media (jumlah) tanpa durasi; 2025 media konsisten pada total |
| TKD: subtes | tidak ada | tidak ada | 4 subtes @25 soal (nama didukung judul beberapa penerbit; hitungan hanya Dealls); menit 18/27/9/20 (jumlah 74) | 2025: satu media (hitungan), nama media konsisten; **bertentangan** dengan deskripsi TWK/TIU/TKP (Kompas, Suara) |
| AKHLAK: jumlah soal / durasi | 54 soal, digabung TKD, 120 menit (RBG) | "dari 90" (IDN) vs **resmi "dari 100"** | 90 / 30 menit | 2023 satu media; 2024 **bertentangan**; 2025 media konsisten (Dealls, detikSulsel; independensi tidak terbukti) |
| TWK: jumlah / durasi | tidak ada komponen | tidak ada | 10 soal; 5 menit (Dealls) atau 10 menit (detikSulsel) | 2025 **bertentangan** pada durasi |
| English: jumlah / durasi | 140 / 115 (Cakap; Listening 50/35, Structure 40/25, Reading 50/55) | 100 / 90 (Tirto, KitaLulus; belum dibuka) | 85 / 80 (ANTARA, IDN); kumparan 140/115 = salinan | 2023 satu media (bimbel); 2024 satu media (ringkasan); 2025 media konsisten, **bertentangan** dengan kumparan; format berubah antar tahun |
| Learning Agility | tidak ada | 50 soal / 30 menit / skor maks 300 (Kompas); min 197 resmi | 100 soal / 35 menit (ANTARA, IDN) vs 100/40 (detikSulsel) vs 50/30 (Dealls, kemungkinan salinan 2024) | 2024 min resmi, jumlah satu media; 2025 **bertentangan** |
| Skala skor | ambigu ("dari jumlah soal yang tersedia") | skor berskala: dari 100 / 677 (konversi /675) / 300 | tidak ada dokumen | 2024: **resmi**; 2023 ambigu; 2025 asumsi sama |
| Penalti jawaban salah | tidak ditemukan | tidak ditemukan | tidak ditemukan | **tidak ada sumber** |
| Platform dan proctoring | Chrome/Firefox, webcam wajib, registrasi wajah, tanpa SEB | SEB wajib, webcam, registrasi wajah, kamera mati 3 menit = gugur, dilarang kalkulator, login ulang maks 3x | SEB (tutorial resmi); aturan lain media (AI proctoring, tanpa kalkulator/headset/masker) | 2023, 2024: **resmi**; 2025: SEB resmi, aturan lain media konsisten |
| Fresh graduate vs experienced hire | tidak dibedakan | beda usia maks (25/27/30/35 vs 45) dan IPK; komponen tes sama menurut jenjang | tidak ada | 2024: **resmi** (usia), komponen tes tidak ada beda di dokumen |
| Aturan peringkat | "ranking sesuai kebutuhan BUMN terdaftar" | idem; rasio 1:10 (FHCI via Kompas.tv) | rasio 1:200 (detikSumut, kutipan @fhci.bumn) | 2023, 2024 FAQ **resmi**; rasio angka media (satu per tahun) |

Catatan bobot 2023: FAQ resmi 2023 menulis T1 D3-S2 "TKD (40%) & AKHLAK (60%)" dan T2 D3-S2 "TKD & AKHLAK (60%) & Bahasa Inggris (40%)", jadi gabungan T1 60 : T2 40, konsisten dengan 2024 (T1 60 : T2 40). Untuk SMA T2 2023 teksnya mengulang "TKD (40%) & AKHLAK (60%)" tanpa bobot terpisah; tidak saya tafsirkan lebih jauh.

## Siap dispesifikasikan?

**Rekomendasi (diperbarui 19 Nov 2026, setelah pass mendalam bagian 7): struktur, ambang, bobot, skala skor, dan aturan platform RBB 2023 dan 2024 sekarang resmi dan cukup untuk data ujian berversi. Jumlah soal, durasi, dan subtes belum resmi untuk batch mana pun, sehingga harus berupa parameter berlabel "tidak resmi" yang bisa diubah.** Modul BUMN dapat dispesifikasikan dengan spesifikasi berlapis: lapis resmi (struktur, ambang, bobot) dan lapis tidak resmi (butir, menit).

### Per komponen: cukup untuk spesifikasi berversi?

| Komponen / hal | Cukup? | Label keyakinan di data | Satu dokumen yang menutup celah tersisa |
|---|---|---|---|
| Struktur tahap per jenjang (SMA vs D3-S2), daftar komponen | Ya, per batch (`rbb-2023`, `rbb-2024`) | resmi (beranda dan PDF "Rincian Detail" arsip). 2025: komponen resmi, pembagian tahap media | Beranda/PDF "Rincian Detail Tahapan Seleksi RBB 2025" (tidak tersimpan di arsip karena SPA); tangkapan layar portal 2025 dari peserta |
| Ambang (TKD 58, AKHLAK 65, TWK 50, English 450, LA 197/300) | Ya untuk 2023, 2024 | resmi; 2025: media konsisten kecuali LA (bertentangan) | FAQ portal 2025 atau pengumuman hasil 2025 yang memuat ambang LA |
| Bobot per tahap dan gabungan (60/40, 40/60) | Ya untuk 2024; 2023 sebagian (T2 SMA ambigu) | resmi | PDF "Tes Online 2" edisi 2025 |
| Skala skor (TKD/AKHLAK/TWK dari 100; English dari 677 dinormalkan /675; LA dari 300) | Ya, sebagai skor berskala, bukan jumlah benar mentah | resmi (FAQ 2024) | Penjelasan butir-ke-skor (rumus konversi) tidak ada di dokumen mana pun; hanya layar instruksi/hasil peserta |
| TKD: jumlah soal / durasi | Sebagai parameter 100 soal, menit diatur (73 min dari media 2025) | tidak resmi (media konsisten 2025; satu media 2023/2024) | Layar instruksi per subtes pada tes uji coba/tes sebenarnya (FAQ menyatakan di sana), difoto peserta |
| TKD: subtes | Sebagai profil latihan dengan 4 subtes ber-label "tidak resmi" (hitungan satu media, nama banyak penerbit); jangan sebagai fakta | tidak resmi, satu media | idem (layar instruksi subtes) |
| AKHLAK: jumlah/durasi | Parameter 90 soal / 30 menit, tidak resmi; 2024 bertentangan dengan skala resmi "dari 100" | tidak resmi | idem |
| Aturan skor AKHLAK per opsi, penalti salah | Tidak. Gunakan skor tanpa penalti sebagai pilihan desain aplikasi, bukan klaim resmi | tidak ada sumber | idem; atau pernyataan FHCI tertulis |
| TWK: jumlah/durasi | Parameter 10 soal; menit 5 atau 10 (bertentangan) | tidak resmi, bertentangan | idem |
| English: jumlah/durasi/subtes | Tidak (140/115, 100/90, 85/80 berubah tiap tahun; listening tidak bisa dibuat); skip untuk fase 2 awal | tidak resmi, bertentangan | Layar instruksi tes English per tahun |
| Learning Agility | Hanya skala skor (dari 300, min 197) untuk 2024; jumlah/durasi/format bertentangan | min resmi 2024; sisanya satu media | Layar instruksi LA |
| Platform dan proctoring (aturan latihan: SEB, kamera, tanpa kalkulator) | Ya untuk 2023, 2024; 2025 SEB resmi | resmi (2023, 2024, SEB 2025) | Panduan peserta 2025 |
| Fresh graduate vs experienced hire | Tes sama per jenjang; hanya syarat usia/IPK/dokumen yang beda | resmi (usia/IPK 2024); tes: tidak ada beda di dokumen | Tidak perlu |
| Nilai AKHLAK (6 nilai, definisi, panduan perilaku) | Ya, korpus verbatim | resmi (SE-7/2020; status berlaku belum dicek) | - |
| Batch 2026 | Tidak | - | Pengumuman `@fhci.bumn` |

### Langkah yang disarankan

1. Pakai `rbb-2024` sebagai batch acuan: satu-satunya batch dengan struktur per jenjang, bobot, ambang, skala, dan aturan platform yang semuanya resmi. Isi jumlah soal/durasi dengan angka 2025 berlabel "tidak resmi, batch 2025, media" atau biarkan kosong dan wajib diisi pengguna.
2. Simpan sumber per kolom (URL arsip + timestamp snapshot + resmi/tidak) di data spesifikasi, sehingga angka tidak resmi bisa ditukar saat ada dokumen yang menutup celah.
3. Satu dokumen terbaik untuk menutup celah jumlah/durasi/subtes: tangkapan layar halaman instruksi tiap subtes (atau bagian "uji coba") dari peserta RBB 2024 atau 2025; FAQ resmi menyatakan angka itu tampil di sana. Alternatif: unggahan Instagram `@fhci.bumn` berisi kisi-kisi (tidak saya buka).
4. Jangan tampilkan "lulus/tidak lulus" sebagai klaim resmi untuk 2025 karena ambang Learning Agility bertentangan; untuk TKD, AKHLAK, TWK, ambang 2024 resmi dan stabil 2023-2025.
5. Jangan modelkan English dan Learning Agility di fase awal; format dan penilaian tidak stabil.
6. Cek `@fhci.bumn` untuk RBB 2026 sebelum menentukan target tahun (belum ada informasi resmi per sumber di bagian 1).

### Penilaian sebelumnya (19 Nov 2026, sebelum arsip resmi dibaca; dipertahankan sebagai riwayat)


**Rekomendasi lama: masih belum untuk spesifikasi resmi lengkap. Data RBB 2024 cukup untuk profil latihan berlabel "non-resmi" (struktur tahap, ambang, bobot), tetapi tidak cukup untuk jumlah soal, durasi, dan subtes.** Hanya sebagian yang bisa dimodelkan, dengan status sumber di data.

| Data RBB 2024 | Cukup untuk data ujian berversi? | Alasan |
|---|---|---|
| Struktur tahap per jenjang (D3-S2 dan SMA), komponen, bobot, ambang (TKD 58, AKHLAK 65, TWK 50, English 450, LA 197/300), bobot gabungan 60/40 dan 40/60, rasio 1:10 | Ya, sebagai `rbb-2024-tidak-resmi` dengan sumber per kolom | Konsisten di empat media bertanggal 2024 yang mengutip FHCI; bukan resmi |
| Jumlah soal dan durasi TKD, AKHLAK, TWK, English | Tidak | Tidak ada sumber 2024 yang dibuka; angka 2025 (100/73, 90/30, 10/10) berasal dari tahun lain |
| Learning Agility 50 soal / 30 menit / 300 poin | Sebagian | Hanya satu artikel Kompas 2024; berbeda dari angka 2025 |
| Subtes TKD | Tidak | Tidak ada sumber 2024 yang dibuka |

Penilaian awal (sebelum tambahan 2024), tetap berlaku untuk elemen lain:

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
0. (Tambahan 2024) Jika modul ditargetkan, pakai struktur 2024 di bagian 6 sebagai dasar profil latihan, karena lebih konsisten daripada 2025. Cari "Rincian Detail Tahapan Seleksi RBB 2024" dan "Petunjuk Tes Online" resmi (peserta 2024 mungkin menyimpan PDF-nya) untuk menutup celah jumlah soal dan durasi.
1. Buka portal peserta RBB 2025 (atau 2024) dari peramban biasa dan simpan tangkapan layar kisi-kisi/FAQ resmi (inilah satu-satunya jalur ke angka resmi). Bila ada panduan peserta PDF, itu menggantikan semua angka media di atas.
2. Cek `@fhci.bumn` untuk RBB 2026 sebelum memutuskan target tahun.
3. Bila tetap ingin lanjut tanpa itu: modelkan spesifikasi sebagai "profil latihan" berversi (`rbb-latihan-v1`) yang jelas bukan spesifikasi resmi, dengan jumlah soal dan durasi yang bisa diubah, label "non-resmi" di UI, dan jangan menampilkan "lulus/tidak lulus" berdasarkan ambang media.
4. Prioritas konten yang aman: TKD numerik dan logika (pakai ulang generator), bank AKHLAK situasional dari SE-7/2020 dengan kunci berlabel buatan aplikasi. Tunda English, Learning Agility, subtes gambar dan klasifikasi kata.
