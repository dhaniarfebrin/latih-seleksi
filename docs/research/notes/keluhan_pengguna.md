# Keluhan dan keinginan pengguna aplikasi latihan CPNS SKD dan BUMN (Rekrutmen Bersama BUMN), per November 2026

Catatan metodologi (penting bagi penulis laporan): tool yang tersedia hanya WebSearch (ringkasan hasil, tanpa Reddit/domain reddit.com diblokir) dan WebFetch. Halaman ulasan Google Play tidak bisa dimuat (WebFetch hanya mengembalikan header halaman), Reddit tidak dapat diakses, X/TikTok/YouTube tidak terindeks, dan Bisnis.com mengembalikan HTTP 403. Akibatnya, bukti ulasan pengguna primer sangat tipis. Hampir semua klaim di bawah adalah bukti lemah (artikel media, blog komersial, atau satu-dua ulasan App Store). Tidak ada frekuensi kuantitatif keluhan yang berhasil ditemukan. Jangan menyajikan sinyal frekuensi di bawah ini sebagai statistik.

## 1. Keluhan umum di ulasan aplikasi tryout CPNS/BUMN (kunci jawaban salah, pembahasan buruk, soal tidak sesuai, iklan, paywall, bug, timer, TKP)

### Takeaway
Ulasan Google Play tidak berhasil dibaca. Satu-satunya bukti ulasan primer adalah halaman App Store aplikasi MisiPNS (730 rating, rata-rata 4.9), yang menunjukkan masalah nyata di sisi langganan/paywall dan koreksi kunci jawaban lewat catatan versi. Kategori keluhan lain (iklan, timer, soal terlalu mudah) hanya didukung sinyal tidak langsung. Bukti lemah.

### Cited Findings
- Seorang pengguna MisiPNS di App Store (ulasan bertanggal 11 Sep, tahun tidak ditampilkan di teks halaman; kemungkinan 2026) melaporkan bahwa setelah membayar langganan, soal latihan tidak termuat, dan logout-login mengembalikan akun ke status gratis. Itu satu ulasan dari empat yang tampil; tiga lainnya positif (menyebut membantu persiapan CPNS 2027, dan seorang pengguna menyebut kerabatnya menilai isi soal sangat mirip ujian 2024). Bukti lemah (n=4). — [MisiPNS di App Store](https://apps.apple.com/id/app/misipns-tryout-cpns-pppk/id6759079870?l=id)
- Catatan versi MisiPNS mengonfirmasi bahwa kunci jawaban dan pembahasan memang pernah salah dan diperbaiki: v2.7.0 (23 Jun) "merapikan dan memperbaiki kunci jawaban TWK, TKP, TIU"; v2.18.0 (13 Okt) "memperbaiki kunci jawaban dan pembahasan beberapa soal". Ini bukti tidak langsung dari pengembang, bukan ulasan pengguna. — [MisiPNS di App Store](https://apps.apple.com/id/app/misipns-tryout-cpns-pppk/id6759079870?l=id)
- Catatan versi yang sama menunjukkan bug paywall berulang: penawaran premium muncul berulang saat aplikasi dibuka (v2.7.0), akses premium kadang macet setelah pembelian (v2.9.0, v2.17.0), halaman langganan membeku di iOS (v2.6.0), pembelian paket lain menurunkan paket yang ada (v2.17.0), tanggal kedaluwarsa salah (v2.18.0). Tahun versi tidak tampil di teks; urutan versi dan "Sep/Okt" konsisten dengan 2026. — [MisiPNS di App Store](https://apps.apple.com/id/app/misipns-tryout-cpns-pppk/id6759079870?l=id)
- Sebuah hasil pencarian merangkum ulasan App Store lain tentang bug tampilan di mana jawaban benar kadang berwarna merah (ringkasan hasil pencarian atas halaman yang sama; tidak terlihat langsung di fetch saya). Bukti sangat lemah. — [MisiPNS di App Store](https://apps.apple.com/id/app/misipns-tryout-cpns-pppk/id6759079870?l=id)
- Model bisnis umum: simulasi terbatas gratis, kunci jawaban lengkap, pembahasan, dan riwayat tryout tak terbatas dikunci di balik pembelian dalam aplikasi (contoh "CPNS 2026 - Simulasi CAT SKD", update terakhir 2026-09-06, 296+ instal, 0 rating). Ini deskripsi struktur paywall, bukan keluhan pengguna. — [AppGoblin](https://appgoblin.info/apps/com.zonjaxebdevs.cpns_simulasicatskd)
- Rekrutmencpns.id membatasi pembahasan lengkap dan kunci jawaban detail ke paket Premium (menurut ringkasan hasil pencarian; halaman tidak saya fetch). Bukti lemah. — [rekrutmencpns.id/tryout](https://rekrutmencpns.id/tryout)
- Satu artikel review umum (blog, tanpa tanggal terverifikasi) menyebut kualitas soal sebuah platform cukup akurat, tetapi sebagian soal dianggap terlalu mudah (dari ringkasan hasil pencarian; saya tidak membuka halamannya). Bukti sangat lemah dan kemungkinan konten promosi. — [omgindonesia.com](https://omgindonesia.com/pendidikan/review-tryout-online-cpns-terbaru-gratis-mana-yang-paling-akurat/)
- Opini lama (16 Des 2019) peserta: tryout membuat orang terpaku pada soal yang pernah dilatih, malas bernalar karena jawaban sudah tersedia, dan pikiran menjadi "tumpul" saat ujian asli. Ini opini satu penulis dan sudah lama. — [Kompasiana, 2019](https://www.kompasiana.com/mulyo26281/5df7b55fd541df676b1587d5/manfaat-positif-dan-negatif-ikut-try-out-cpns)
- Klaim marketing sisi tandingan: Tryout.id menyatakan aplikasinya "tanpa iklan yang mengganggu". Ini mengimplikasikan iklan dipandang sebagai pain point oleh kompetitor, tetapi bukan keluhan pengguna. — [tryout.id](https://tryout.id/)
- Tidak ditemukan: ulasan Google Play yang bisa dikutip untuk Ayiip, Eksam, Viracun, "Simulasi Tes CPNS (CAT) Pro", dan lain-lain. Satu agregator melaporkan 804 ulasan/4,82 untuk "Simulasi Tes CPNS (CAT) Pro" tetapi menandai confidence rendah dengan 0 ulasan dianalisis. — [marlvel.ai](https://marlvel.ai/api/llm/apps/events/tescpnsasn-simulasicatcpns-bimbelcpnsjadiasn)

### Inferences
- Bug paywall/langganan (akses premium tidak aktif setelah bayar, penawaran premium berulang) adalah kategori kegagalan teknis yang nyata untuk aplikasi tryout berlangganan, terbukti dari ulasan dan 7 catatan perbaikan dalam satu tahun di satu aplikasi. Ini sinyal kuat untuk satu aplikasi, tapi belum terbukti umum.
- Kesalahan kunci jawaban tampaknya biasa pada bank soal buatan pihak ketiga (MisiPNS memperbaikinya dua kali dalam beberapa bulan), sehingga mekanisme pelaporan soal/koreksi cepat bisa jadi diferensiasi. Ini inferensi dari satu aplikasi.
- Pengguna yang bayar kemungkinan menuntut kepastian akses; kehilangan status premium merusak kepercayaan lebih dari soal yang kurang bagus (inferensi, tanpa data).

### Gaps
- Ulasan Google Play sebenarnya (kata kunci "iklan", "kunci jawaban salah", "pembahasan", "timer", "paksa langganan") tidak dapat diakses; perlu diambil manual atau lewat scraper (misalnya google-play-scraper, tanggal dan app id: com.id.eksam, com.misipns.app, com.ayiip.cpns, com.cekotomatis.eformbri, com.zonjaxebdevs.cpns_simulasicatskd) untuk menghitung frekuensi.
- Tidak ada data frekuensi untuk keluhan apa pun.
- Keluhan "soal terlalu mudah / tidak HOTS", "soal berulang", "timer tidak sesuai CAT" tidak punya bukti langsung dari pengguna dalam riset ini.

## 2. Diskusi Reddit, X, Kaskus, Quora, YouTube, TikTok: apa yang membantu lolos vs membuang waktu

### Takeaway
Reddit tidak dapat dijangkau; Kaskus hanya menghasilkan thread lama (2018-2020). Bukti terkuat adalah esai peserta dan berita yang merekam pola lama: TKP dikeluhkan paling banyak, passing grade per subtes menjebak orang berskor total tinggi. Bukti 2024-2026 dari forum sosial: tidak ditemukan.

### Cited Findings
- Esai Mojok (12 Nov 2018, tentang SKD 2018): peserta menyebut TKP sebagai penyebab kegagalan terbanyak karena semua opsi tampak positif, tanpa dasar jelas untuk memilih, dan penilaiannya terasa subjektif; passing grade TKP 143 disebut tinggi; BKN menyebut hanya sekitar 9% peserta nasional memenuhi passing grade. Beberapa peserta juga bilang kisi-kisi yang mereka pelajari berbeda-beda dan soal TKP "di luar dugaan". Esai opini/satir, bukan data. — [Mojok, 2018](https://mojok.co/esai/hikmah-protes-soal-seleksi-kompetensi-dasar-cpns/)
- Thread Kaskus tentang petisi online soal tes CPNS (2018) mempertanyakan soal TKP tanpa acuan jawaban pasti dan apakah soal sudah melalui uji validitas/reliabilitas; satu peserta mengaku mengira skornya cukup padahal tidak mencapai passing grade TKP 143 meski total tinggi. Dikutip dari ringkasan hasil pencarian, tidak saya fetch langsung. Bukti lemah dan lama. — [Kaskus, 2018](https://www.kaskus.co.id/thread/5be8ba0312e2577b178b4567/banyak-peserta-tak-lulus-muncul-petisi-online-soal-tes-cpns)
- Esai Mojok (26 Mar 2020): penulis honorer yang gagal SKD 2017 dan 2018 lalu lolos passing grade 2020, menyebut yang membantu adalah buku bank soal berpuluh-puluh halaman dan keanggotaan Instagram berbayar untuk latihan soal; menyebut keanggotaan itu mahal bagi guru honorer yang digaji per tiga bulan; menyimpulkan keberuntungan (ranking 27, hanya 3 teratas lanjut SKB) sama pentingnya dengan kemampuan. Satu kisah, 2020. — [Mojok, 2020](https://mojok.co/terminal/pintar-aja-nggak-cukup-lolos-seleksi-cpns-tuh-yang-paling-penting-hoki/)
- Sebuah thread Kaskus berjudul "Kiat sukses lolos CPNS lewat bimbingan belajar" dan berita yang dikutip menyebut tidak ada jaminan lolos meski ikut bimbel (dari ringkasan hasil pencarian; judul thread terverifikasi di hasil, isi tidak). Bukti sangat lemah. — [Kaskus](https://www.kaskus.co.id/thread/59ba1096c0cb177f0a8b456a/kiat-sukses-lolos-cpns-lewat-bimbingan-belajar/2)

### Inferences
- TKP adalah titik sakit berulang lintas tahun (2018 sampai artikel 2026), dan ketidakjelasan alasan skor TKP adalah kebutuhan yang belum terpenuhi. Kesimpulan berdasarkan pola lintas sumber, tapi sumber 2024-2026 berasal dari blog komersial, bukan pengguna.
- Biaya akses (keanggotaan berbayar) adalah hambatan bagi pelamar berpenghasilan rendah (honorer), mendukung tier gratis yang bermakna. Satu anekdot.

### Gaps
- Tidak ada thread Reddit, tweet, komentar YouTube, atau ringkasan TikTok yang bisa dibaca. Perlu pencarian manual di Reddit (r/indonesia, r/finansial, r/JobsIndonesia jika ada) dan X dengan kata kunci "tryout SKD", "nilai tryout vs asli", "bimbel CPNS nipu".
- Tidak ada testimoni 2024-2026 tentang "apa yang membantu lolos" dari pengguna non-komersial.

## 3. Cara belajar kelemahan saat ini dan apa yang diinginkan

### Takeaway
Hanya sumber komersial (blog SEO) dan satu paper sistem tryout yang membahas analisis kelemahan; tidak ada bukti langsung tentang pengguna memakai spreadsheet atau catatan. Kebutuhan analisis per subtes, pola kesalahan, dan manajemen waktu tampak kuat di materi pemasaran, tapi ini bukan suara pengguna.

### Cited Findings
- Blog komersial (7 Agu 2026, jelas promosi layanan tryout, tanpa sumber): menyebut penyebab "salah diagnosis kelemahan"; tiga pola gagal: total tinggi tapi satu subtes di bawah ambang, kehabisan waktu karena terlalu lama di TWK, dan salah membaca skala TKP (memilih jawaban yang terdengar paling bermoral, bukan yang skornya tertinggi); menyarankan tryout yang melaporkan skor per subtes, status lulus per skema, jumlah benar/salah/kosong, waktu pengerjaan, dan log pola kesalahan, lalu ulang sekitar dua kali seminggu. Bukti lemah (promosi). — [SZ Media, 2026](https://www.semarangzoo.com/2026/08/strategi-tryout-cpns-setelah-gagal-skd-dan-materai-elektronik.html)
- Jurnal sistem tryout CPNS (UNESA, tanpa tanggal terverifikasi) merancang sistem yang mengidentifikasi soal salah dan merekomendasikan materi terkait, serta menyatakan peserta perlu fokus ke kelemahan. Ini desain sistem akademik, bukan survei pengguna. — [JINACS UNESA](https://ejournal.unesa.ac.id/index.php/jinacs/article/download/75138/53488/192615)
- Satu sumber menyebut banyak peserta gagal karena waktu habis pada TIU, bukan karena tidak tahu jawabannya; menyarankan eliminasi opsi dan latihan pola. Blog tanpa data. — [pdiperjuanganbali.id](https://www.pdiperjuanganbali.id/trik-cepat-menjawab-soal-tiu-skd-cpns)
- Simulasi resmi BKN (cat.bkn.go.id/simulasi) hanya untuk mengenal sistem CAT, bukan kisi-kisi, dan memiliki kuota harian terbatas. — [Kompas TV, 2024](https://www.kompas.tv/nasional/540866/jadwal-tryout-skd-cpns-2024-online-resmi-dari-bkn-ini-link-cara-daftar-dan-kisi-kisinya?page=all)

### Inferences
- Fitur yang konsisten disebut di konten penyedia tryout: skor per subtes, status lulus per skema ambang, analisis waktu, pola kesalahan per tipe soal, dan penjelasan skor TKP. Kemungkinan ini adalah fitur yang belum dianggap standar (karena dijual sebagai pembeda), tapi belum diverifikasi lewat pengguna.
- Praktik spreadsheet/catatan manual tidak punya bukti dalam riset ini (lihat Gaps).

### Gaps
- Tidak ada bukti langsung bahwa orang memakai spreadsheet, buku catatan salah, atau mengulang tryout; tidak ada kutipan "saya berharap ada ...".
- Perlu survei/wawancara pengguna atau pencarian Reddit/X untuk kata kunci "catatan soal salah", "spreadsheet tryout", "buku kesalahan".

## 4. Pain khusus pekerja profesional dan tes BUMN (AKHLAK, Inggris, learning agility, psikotes)

### Takeaway
Satu-satunya keluhan peserta BUMN yang terdokumentasi adalah learning agility (RBB 2024): dinilai sulit dan mengecoh meski materi dikuasai. Banyak sumber latihan gratis (Dealls, Tirto, Bisnis, Kitalulus), tetapi kecocokannya dengan soal resmi tidak dapat diverifikasi. Tidak ada bukti tentang kebutuhan khusus pekerja (waktu terbatas, mobile).

### Cited Findings
- Bisnis.com (28 Mei 2024, hasil pencarian merangkum karena halaman 403): tes learning agility "viral" sulit pada RBB 2024; peserta mengeluh sulit mendapat nilai tinggi meski menguasai soal, sebagian bilang soal mengecoh sehingga banyak jawaban tidak sesuai; artikel menjelaskan cara hitung penilaian. Tidak ada data resmi tingkat kesulitan. Bukti lemah-sedang (media nasional, kutipan tidak terverifikasi langsung). — [Bisnis.com, 28 Mei 2024](https://kabar24.bisnis.com/read/20240528/79/1768946/viral-tes-learning-agility-sulit-ini-cara-hitung-penilaiannya-agar-lulus)
- Learning agility adalah subtes baru di RBB 2024; format 2025 disebut berbeda antarsumber (Dealls: 50 soal/30 menit; Tirto: 100 soal/maks 35 menit); Bahasa Inggris juga bervariasi (Dealls: 140 soal/115 menit; Tirto: 85 soal/80 menit). Perbedaan ini menandakan format tidak transparan atau berubah; materi latihan pihak ketiga bisa tidak cocok. — [Dealls](https://dealls.com/pengembangan-karir/learning-agility-bumn); [Tirto](https://tirto.id/tips-jitu-fhci-tes-online-tahap-2-rekrutmen-bersama-bumn-2025-hb2k)
- Passing grade RBB 2024 yang dilaporkan: AKHLAK 65, Wawasan Kebangsaan 50, learning agility 197 (skala skor, bukan jumlah soal); sumber lain menyebut TKD minimal 58 dari 100 dan AKHLAK minimal 65 dari 90 soal; angka beda antarsumber. Pembobotan 2024: Bahasa Inggris 80% dan learning agility 20%. Tidak ada angka resmi terbaru yang dikonfirmasi. — [Bisnis.com, 2024](https://kabar24.bisnis.com/read/20240528/79/1768946/viral-tes-learning-agility-sulit-ini-cara-hitung-penilaiannya-agar-lulus); [Dealls AKHLAK](https://dealls.com/pengembangan-karir/nilai-akhlak-bumn)
- Tersedia banyak konten latihan gratis untuk AKHLAK, learning agility, Bahasa Inggris, bahkan digital mindset (Dealls, Tirto, CNN Indonesia, Kitalulus, Suara), yang berarti pasar latihan konten statis padat; sumber keluhan tentang kualitasnya tidak ditemukan. — [Dealls digital mindset](https://dealls.com/pengembangan-karir/digital-mindset-bumn); [Kitalulus](https://www.kitalulus.com/blog/info-cpns/contoh-tes-learning-agility-bumn/)

### Inferences
- Learning agility adalah subtes yang peserta anggap tidak jelas penilaiannya; aplikasi yang menjelaskan metode skor dan menyediakan simulasi format asli bisa mengisi celah. Berdasarkan satu laporan media 2024, belum ada data 2025-2026.
- Hubungan antara tes BUMN dan pekerja profesional (kebanyakan pelamar RBB sudah bekerja atau fresh graduate) bersifat asumsi; tidak ada sumber di sini.

### Gaps
- Tidak ada bukti spesifik tentang pelamar yang bekerja: keterbatasan waktu, penggunaan mobile, sesi singkat. **Asumsi (belum pasti):** pekerja butuh sesi latihan singkat di ponsel; belum diverifikasi, perlu survei atau wawancara.
- Tidak ada keluhan peserta tentang AKHLAK, Bahasa Inggris, atau psikotes BUMN dari forum nyata; sumber RBB 2025 dan 2026 tidak ditemukan.

## 5. Selisih skor tryout vs skor SKD asli

### Takeaway
Tidak ada data pengguna yang membandingkan skor tryout dengan skor asli. Hanya ada penjelasan struktural dari penyedia tryout (passing grade per subtes, pola TKP berbeda, simulator hanya meniru format). Bukti lemah.

### Cited Findings
- Ringkasan hasil pencarian atas situs penyedia tryout (2024-2026): skenario "total tinggi tapi gagal" umum karena ambang per subtes; ambang SKD 2024 yang sering dikutip TWK 65, TIU 80, TKP 166 dari maksimal 550; sumber memuat inkonsistensi (100 vs 110 soal). Klaim penyedia, bukan pengguna. — [Kompas TV, 2024](https://www.kompas.tv/lifestyle/540103/10-link-tryout-cpns-2024-gratis-untuk-latihan-soal-dan-nilai-ambang-batas-yang-harus-diketahui?page=all); [Kawan Kelas](https://kawankelas.id/tryout-cpns-gratis-2026-10-link-dan-cara-memaksimalkannya)
- cat.bkn.go.id adalah satu-satunya simulasi CAT resmi BKN; simulator lain hanya meniru format sehingga skor belum tentu setara ujian asli (ringkasan hasil pencarian; tanpa data pembanding). — [Kompas.com, 10 Sep 2024](https://www.kompas.com/tren/read/2024/09/10/103000865/cara-ikut-simulasi-cat-bkn-gratis-untuk-latihan-skd-cpns-2024?page=all)
- Skor TKP berkisar 45-225 dengan tiap opsi bernilai 1-5; pola pemberian skor berbeda dari yang dibiasakan di tryout dapat menurunkan skor asli (klaim penyedia tryout, tanpa data). — [Kawan Kelas](https://kawankelas.id/tryout-cpns-gratis-2026-10-link-dan-cara-memaksimalkannya)
- Soal SKD 2018-2019: BKN menyatakan standar soal sama tetapi bentuknya lebih mengecoh/HOTS; detik 2019: passing grade turun namun soal lebih sulit. Berita lama, bukan 2024-2026. — [detik 2018](https://finance.detik.com/lowongan-kerja/d-4296167/bkn-soal-tes-cpns-tahun-ini-mengecoh); [detik 2019](https://finance.detik.com/cpns/d-4874451/jangan-terlena-passing-grade-turun-soal-tes-cpns-lebih-susah-lho)
- Berita 2024: tiga peserta peraih nilai tertinggi tes CPNS 2024 hanya salah satu soal (artikel dikutip di hasil pencarian; detail tidak diverifikasi). Menunjukkan skor tinggi mungkin tetapi bukan indikator rata-rata. — [Kompas TV](https://www.kompas.tv/ekonomi/550469/3-orang-peraih-nilai-tertinggi-tes-cpns-2024-salah-cuma-1-soal?page=all)

### Inferences
- Pain point yang masuk akal: pengguna salah menafsirkan skor tryout total tanpa melihat ambang per subtes, dan skor TKP tryout tidak mengkalibrasi ke skala asli. Penjelas skor per subtes dan estimasi probabilitas lolos bisa bernilai, tetapi keinginan pengguna terhadap fitur ini tidak terverifikasi.
- Kalibrasi skor (membandingkan skor tryout dengan hasil asli yang dilaporkan pengguna) tidak dilakukan siapa pun dalam sumber yang ditemukan; peluang diferensiasi, tapi belum terbukti ada permintaan.

### Gaps
- Tidak ada laporan pengguna "skor tryout X, skor asli Y" yang ditemukan. Perlu pencarian manual di Reddit, X, dan grup Telegram/Facebook.
- Tidak ada data pembanding resmi tentang validitas prediktif tryout pihak ketiga.
