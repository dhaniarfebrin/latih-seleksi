# Ulasan Google Play aplikasi latihan CPNS SKD dan BUMN, dikumpulkan 14 November 2026

Catatan metodologi (penting bagi penulis laporan). Ulasan publik Google Play diambil dengan pustaka `google-play-scraper` 1.2.7 (lang `id`, country `id`, jeda sekitar 1 sampai 1,5 detik per permintaan) pada 14 November 2026, dari 16 aplikasi yang dipilih dari hasil pencarian Play Store (kueri seputar CPNS, SKD, CAT, BUMN, AKHLAK). Per aplikasi diambil sampai 200 ulasan urutan terbaru dan sampai 200 urutan paling relevan; setelah dedup per ID ulasan terkumpul 3.384 ulasan, 2.331 di antaranya bertanggal 2024 sampai 2026 (769 dari 2024, 597 dari 2025, 965 dari 2026). Hanya skor, teks, jumlah suka, versi aplikasi, tanggal, dan penanda balasan pengembang yang disimpan; nama dan foto pengulas tidak disimpan dan tidak dikutip. Data mentah hanya di folder sementara, bukan di repo. Klasifikasi memakai aturan kata kunci (regex) lalu dicek dengan membaca sampel manual; angka di bawah adalah **hasil kata kunci**, bukan hitungan manual semua ulasan. Batas metode ada di bagian Gaps.

Cara menghitung. Populasi utama keluhan: ulasan bintang 1 sampai 3 tahun 2024 sampai 2026 dari 13 aplikasi SKD/BUMN murni ("inti", n=204 dari 1.718 ulasan). Tiga aplikasi campuran (Analitica, Al Faiz, Tryout.id; isinya juga UTBK dan ujian lain) dilaporkan terpisah. Satu ulasan boleh masuk beberapa kategori, jadi persentase tidak dijumlahkan. Ulasan Tryout.id bulan Oktober 2025 (278 ulasan, hampir semua bintang 5 dan bernada seragam) dikeluarkan karena diduga ulasan promosi.

## 1. Skala dan ulasan yang terkumpul

### Takeaway
Tidak ada aplikasi yang dominan di sisi pengguna, tetapi skalanya jauh lebih besar daripada yang tercatat di laporan: JadiASN, JadiBUMN, dan MisiPNS masing-masing punya ribuan rating Play dan rata-rata di atas 4,5. Ulasan sangat condong positif (88% ulasan inti bintang 4 sampai 5), sehingga keluhan harus dibaca sebagai proporsi dari minoritas bintang rendah, bukan dari semua ulasan.

### Cited Findings
Data per aplikasi (installs, skor, jumlah rating, dan tanggal pembaruan dari halaman Play yang dibaca pada 14 Nov 2026; n sampel = ulasan 2024 sampai 2026 yang diambil, bukan seluruh ulasan aplikasi):

| Aplikasi | Install Play | Skor / rating | Update terakhir | n sampel | n bintang 1-3 | Sumber |
|---|---|---|---|---|---|---|
| JadiASN | 500 rb+ | 4,53 / 9.177 | 5 Nov 2026 | 200 | 52 | [Play](https://play.google.com/store/apps/details?id=id.jadiasn.app&hl=id&gl=id) |
| MisiPNS | 50 rb+ | 4,89 / 8.043 | 13 Okt 2026 | 311 | 18 | [Play](https://play.google.com/store/apps/details?id=com.misipns.app&hl=id&gl=id) |
| BisaCPNS | 100 rb+ | 4,64 / 5.450 | 7 Okt 2024 | 97 | 14 | [Play](https://play.google.com/store/apps/details?id=com.senius.bisacpns&hl=id&gl=id) |
| Simulasi CAT CPNS | 100 rb+ | 4,73 / 4.880 | 29 Jul 2026 | 226 | 17 | [Play](https://play.google.com/store/apps/details?id=id.my.frameless.simulasicat&hl=id&gl=id) |
| JadiBUMN | 100 rb+ | 4,68 / 3.655 | 5 Nov 2026 | 214 | 44 | [Play](https://play.google.com/store/apps/details?id=id.jadibumn.app&hl=id&gl=id) |
| Bimbel & Simulasi CAT CPNS | 10 rb+ | 4,86 / 868 | 4 Okt 2026 | 74 | 7 | [Play](https://play.google.com/store/apps/details?id=tescpnsasn.simulasicatcpns.bimbelcpnsjadiasn&hl=id&gl=id) |
| Materi & Simulasi CAT CPNS | 10 rb+ | 4,58 / 489 | 16 Okt 2025 | 60 | 4 | [Play](https://play.google.com/store/apps/details?id=com.daus.simulasicatcpns&hl=id&gl=id) |
| Soal CPNS 2026 - Latihan CAT | 100 rb+ | 4,56 / 368 | 7 Agu 2026 | 109 | 26 | [Play](https://play.google.com/store/apps/details?id=com.ellstudiosapp.soalcpns&hl=id&gl=id) |
| IDCPNS | 10 rb+ | 4,89 / 363 | 11 Nov 2026 | 239 | 11 | [Play](https://play.google.com/store/apps/details?id=com.idcpns.idcpns_app&hl=id&gl=id) |
| SMARTCPNS | 10 rb+ | 4,58 / 342 | 20 Jul 2026 | 9 | 1 | [Play](https://play.google.com/store/apps/details?id=com.smartcpns.edu&hl=id&gl=id) |
| EKSAM | 5 rb+ | 4,95 / 153 | 28 Agu 2026 | 115 | 3 | [Play](https://play.google.com/store/apps/details?id=com.id.eksam&hl=id&gl=id) |
| Apta BUMN | 10 rb+ | 4,8 / 132 | 31 Okt 2025 | 59 | 4 | [Play](https://play.google.com/store/apps/details?id=cib.apta.bumn&hl=id&gl=id) |
| Tes BUMN TKD AKHLAK | 5 rb+ | 3,33 / 21 | tidak tercatat | 5 | 3 | [Play](https://play.google.com/store/apps/details?id=com.tes.tkdbumn&hl=id&gl=id) |
| Analitica (campuran UTBK/CPNS) | 100 rb+ | 4,86 / 3.793 | 10 Nov 2026 | 256 | 12 | [Play](https://play.google.com/store/apps/details?id=co.xomodo.analitica&hl=id&gl=id) |
| Al Faiz (campuran UTBK/BUMN/SKD) | 100 rb+ | 4,83 / 230 | 4 Nov 2026 | 70 | 34 | [Play](https://play.google.com/store/apps/details?id=com.idekode.alfaiz&hl=id&gl=id) |
| Tryout.id (platform ujian umum) | 10 rb+ | 4,76 / 376 | tidak tercatat | 9 (setelah membuang 278 ulasan Okt 2025) | 6 | [Play](https://play.google.com/store/apps/details?id=id.tryout.app&hl=id&gl=id) |

- Rentang install Play ditampilkan sebagai ambang (misal "100 rb+" berarti 100 rb sampai di bawah 500 rb), jadi tidak bertentangan dengan perkiraan AppGoblin untuk MisiPNS (64 rb+) atau EllStudios (235 rb+). Tetapi JadiASN tercatat 500 rb+ install di Play, di atas "235 ribu+" yang dicatat laporan sebagai terbesar. Ads/IAP dari metadata Play: Soal CPNS 2026 - Latihan CAT, Simulasi CAT CPNS, Bimbel & Simulasi CAT, Materi & Simulasi CAT, dan Tes BUMN TKD AKHLAK mencantumkan iklan; MisiPNS, Apta BUMN, Analitica, dan Simulasi CAT CPNS mencantumkan pembelian dalam aplikasi. (Halaman Play di tabel di atas, dibaca 14 Nov 2026 lewat pustaka scraper, bukan lewat browser.)
- Dua aplikasi yang sudah ada sebagai app BUMN khusus: JadiBUMN (3.655 rating) dan Apta BUMN (132 rating), ditambah Tes BUMN TKD AKHLAK (21 rating). Ruang BUMN tidak kosong di Play.
- Distribusi bintang 1.718 ulasan inti: bintang 5 = 1.369, bintang 4 = 145, bintang 3 = 55, bintang 2 = 38, bintang 1 = 111. Bintang 1 sampai 3 = 204 (11,9%).
- Bintang rendah tidak merata: proporsi bintang 1 sampai 3 dalam sampel adalah 26% (JadiASN), 21% (JadiBUMN), 24% (Soal CPNS 2026 - Latihan CAT), 14% (BisaCPNS), 7% (Simulasi CAT CPNS), 6% (MisiPNS), 5% (IDCPNS), 3% (EKSAM); di Al Faiz 49%. Catatan: sampel adalah gabungan "terbaru" dan "paling relevan", jadi bukan persentase populasi.

### Inferences
- Pasar tidak sesepi dugaan laporan: ada setidaknya lima app dengan rating Play ribuan (JadiASN, MisiPNS, BisaCPNS, Simulasi CAT CPNS, JadiBUMN) selain Analitica. Skala ini mengubah klaim "tidak ada pemain dominan" menjadi "tidak ada pemain dominan yang terverifikasi di iOS/AppGoblin, tetapi ada beberapa pemain menengah di Play" (inferensi dari angka di atas).
- Karena mayoritas ulasan bintang 5 dan pendek ("membantu", "bagus"), proporsi keluhan di bawah harus dibaca sebagai proporsi di antara minoritas yang menulis keluhan, dan jumlah mutlaknya kecil (204 ulasan).

### Gaps
- Install dan rating hanya dari halaman Play pada satu hari; tren tidak diketahui.
- Aplikasi dengan rating di bawah sekitar 100 (Viracun 16, Impia 9, INFO CPNS 97 dan lainnya) tidak diambil ulasannya.
- `com.ayiip.cpns` mengembalikan "App not found" pada 14 Nov 2026 sehingga tidak dianalisis.

## 2. Keluhan: apa yang dikeluhkan pengguna bintang rendah

### Takeaway
Keluhan terbanyak di ulasan bintang 1 sampai 3 adalah **masalah teknis dan transaksi** (aplikasi tidak bisa dibuka, error, login/perangkat, pembayaran pending atau premium belum terbuka, admin lambat; gabungan 26,5% dari bintang 1 sampai 3 inti), diikuti **kunci jawaban/pembahasan salah atau tidak sinkron** (12,7%; 15,7% bila kualitas soal seperti typo dan jawaban rancu ditambah). Harga/paywall (9,3%) dan janji tidak sesuai/"penipuan" (8,3%) menyusul. Keluhan yang sering diduga (iklan, soal berulang, timer tidak sesuai CAT, TKP tidak jelas, tidak ada analitik) muncul jarang (2 sampai 4% dari bintang 1 sampai 3, di bawah 1% dari semua ulasan).

### Cited Findings
Frekuensi kata kunci pada 204 ulasan bintang 1 sampai 3 inti (2024 sampai 2026); kolom terakhir adalah jumlah ulasan bintang 4 sampai 5 dengan kategori yang sama (dari 1.514 ulasan bintang 4 sampai 5), yang sebagian adalah keluhan ringan dari pengguna yang tetap memberi bintang tinggi:

| Kategori | n | % dari bintang 1-3 | n di bintang 4-5 |
|---|---|---|---|
| Bug/aplikasi tidak bisa dibuka, error, data hilang | 25 | 12,3% | 7 |
| Harga/paywall ("harus premium", mahal, fitur terkunci) | 19 | 9,3% | 7 |
| Kunci jawaban salah / jawaban tidak sinkron | 18 | 8,8% | 5 |
| Janji tidak sesuai / "penipuan" / menyesal | 17 | 8,3% | 0 |
| Pembayaran pending / premium belum terbuka / status kembali gratis | 16 | 7,8% | 2 |
| Kualitas soal (typo, rancu, gambar tidak muncul) | 14 | 6,9% | 1 |
| Soal tidak mirip ujian asli / terlalu mudah | 13 | 6,4% | 2 |
| Login/akun/batas perangkat | 13 | 6,4% | 2 |
| Pembahasan kurang/salah/tidak jelas | 10 | 4,9% | 5 |
| Konten kurang lengkap | 9 | 4,4% | 2 |
| Iklan | 8 | 3,9% | 3 |
| Admin/CS tidak merespons | 7 | 3,4% | 0 |
| Soal/pembahasan tidak diperbarui | 6 | 2,9% | 2 |
| Soal berulang | 4 | 2,0% | 3 |
| Timer/CAT tidak sesuai | 4 | 2,0% | 1 |
| TKP membingungkan | 3 | 1,5% | 1 |
| BUMN/AKHLAK kurang | 3 | 1,5% | 0 |
| Tidak ada analitik/riwayat | 1 | 0,5% | 1 |

- Sebanyak 127 dari 204 ulasan bintang 1 sampai 3 (62%) cocok minimal satu kategori; 77 sisanya adalah ulasan sangat pendek atau umum ("jelek", "bagus", "ok") atau di luar kategori (mis. pendaftaran ribet, UX). Gabungan: kunci/pembahasan = 26 ulasan (12,7%); bug/pembayaran/login/CS = 54 ulasan (26,5%); pembayaran/harga/janji = 44 ulasan (21,6%); tidak mirip/outdated/konten kurang = 24 ulasan (11,8%). Sumber: ulasan bintang 1 sampai 3 pada 13 halaman Play di tabel Bagian 1, dibaca 14 Nov 2026.
- Pola kunci/pembahasan (paling relevan untuk gap "kepercayaan pada kunci"): kata kunci kunci/pembahasan muncul di 8 dari 13 aplikasi inti, terbanyak di Soal CPNS 2026 - Latihan CAT (9 dari 26 ulasan bintang 1 sampai 3), Simulasi CAT CPNS (4 dari 17), dan JadiBUMN (4 dari 44), lalu JadiASN (3), BisaCPNS, MisiPNS, Apta, dan Bimbel & Simulasi CAT; pengulas menulis bahwa jawaban benar ditandai salah, pembahasan tidak cocok dengan soal, atau jawaban yang sama berubah antar minggu. Contoh parafrase: "jawaban saya benar saat menjawab tetapi di hasil jadi salah", "pembahasan benar tetapi kunci jawaban salah", "soal dan jawaban berubah-ubah", dan seorang pengulas Apta menulis soal dan jawaban banyak tidak sinkron dan bahwa soalnya sendiri yang ngaco. Satu pengulas IDCPNS menulis sudah melaporkan soal tetapi tidak ada perbaikan; satu pengulas menulis tidak ada cara melaporkan soal di aplikasi. Bukti kuat ada, tetapi hanya sekitar 2% dari semua ulasan inti.
- Pola pembayaran/premium: pengulas menulis sudah bayar tetapi status pending dan admin tidak membalas (BisaCPNS, JadiBUMN, EKSAM, Al Faiz), "sudah premium tetapi tiba-tiba free user", sudah membeli tetapi materi atau fitur tetap terkunci, salah beli paket tanpa bisa ganti atau refund, dan satu pengulas yang gagal berlangganan lewat dompet digital. Ini mengonfirmasi bahwa bug entitlement di catatan versi MisiPNS (App Store) bukan kasus tunggal.
- Pola login dan perangkat: tidak bisa login, sering logout sendiri, batas perangkat (khususnya Al Faiz: 15 dari 34 ulasan bintang 1 sampai 3 menyangkut bug, 9 menyangkut login/perangkat), aplikasi tidak terbuka atau layar biru setelah pembaruan, "network" saat membuat akun.
- Pola harga/paywall dan janji: "harus premium semua", TO gratis syaratnya mengikuti akun, menandai teman, dan mengirim tangkapan layar (IDCPNS), fitur baru terkunci walau sudah membeli 1 tahun, paket mahal yang isinya tidak sesuai (kasus KAI, SPPI, kementerian tertentu; JadiBUMN dan JadiASN), kelas online yang dijanjikan tidak ada, dan beberapa ulasan bintang 1 menyebut "penipuan" atau "bimbel penipuan". Beberapa ulasan membeli pada bulan-bulan menjelang tes.
- Pola kesesuaian ujian: soal TWK terlalu mudah atau hafalan padahal ujian sekarang penalaran, tidak ada soal figural 9 kotak, TIU masih pecahan, "tidak ada satu pun soal di sini yang masuk" saat tes KAI, dan satu pengulas MisiPNS bahwa soal TO terasa relatif mudah dibanding simulasi resmi BKN. Ini menyentuh gap "kebaruan kisi-kisi" dan "kalibrasi skor tryout ke skor asli", walau tidak ada ulasan yang menulis skor tryout X vs skor asli Y.
- Iklan: 8 ulasan bintang 1 sampai 3 dan 3 ulasan bintang 4 sampai 5. Pada lima aplikasi yang mencantumkan iklan di metadata Play (Soal CPNS 2026 - Latihan CAT, Simulasi CAT CPNS, Bimbel & Simulasi CAT, Materi & Simulasi CAT, Tes BUMN TKD AKHLAK; 474 ulasan inti), 11 ulasan (2,3%) menyebut iklan sebagai keluhan. Contoh parafrase: "iklan terlalu panjang", "dikit-dikit iklan", "iklan saat mengerjakan soal mengganggu konsentrasi".
- Soal berulang: 4 ulasan bintang 1 sampai 3 (satu paket SKB yang tiga try out-nya mengulang soal yang sama, satu ulasan "soalnya banyak yang keulang") dan satu ulasan bintang 4 soal sama muncul setelah beberapa kali mencoba. Timer: satu pengulas mengeluh waktu TO hanya 1 menit 40 detik per soal setelah membeli, beberapa mengeluh timer berjalan atau data hilang saat error (Al Faiz, Soal CPNS 2026 - Latihan CAT, MisiPNS). TKP: tiga ulasan (TKP/TWK terasa kekanak-kanakan, tidak masuk akal, pilihan terlalu mencolok); tidak satu pun menulis soal skor 1 sampai 5 per opsi.
- Aplikasi campuran (bintang 1 sampai 3; terpisah dari angka di atas): Al Faiz 34 ulasan (bug 15, login/perangkat 9, janji 5, pembayaran 3, harga 2); Analitica 12 ulasan (harga 2, pembayaran 2, bug 2; satu pengulas menulis Analitica "makin pelit" pada pengguna gratis); Tryout.id 6 ulasan setelah burst dikeluarkan.
- Estimasi presisi (membaca sampel): dari 127 ulasan bintang 1 sampai 3 inti yang cocok, seluruhnya dibaca satu per satu; kategori kunci salah, iklan, login, dan bug hampir semuanya benar. Kesalahan yang saya temukan: kategori pembayaran/premium memuat 3 sampai 4 ulasan yang sebenarnya keluhan kualitas soal setelah membeli (sekitar 12 dari 16 benar), dan "janji/penipuan" mencampur ulasan yang hanya kecewa. Pada bintang 4 sampai 5 saya membaca semua kecocokan kategori kunci, pembahasan, iklan, berulang, dan paywall (sekitar 30 ulasan): sekitar tiga perempat keluhan sungguhan, sisanya pujian (mis. "kunci jawaban lengkap" dan "TO gratis tiap minggu").

### Inferences
- **Terbukti dengan ulasan 2024 sampai 2026**: kepercayaan pada kunci/pembahasan (12,7% dari bintang 1 sampai 3, di 8 dari 13 aplikasi inti) dan kerapuhan akses premium/pembayaran/login (26,5% gabungan). Perlu dicatat bahwa kunci/pembahasan hanya sekitar 2% dari seluruh ulasan inti, jadi bobotnya sebagai pembeda produk adalah hipotesis (inferensi, bukan temuan).
- **Lemah di ulasan**: TKP tidak jelas (3 ulasan), tidak ada analitik (1 keluhan, 2 permintaan riwayat), iklan (2,3% di aplikasi beriklan), soal berulang (4 ulasan), timer tidak sesuai CAT (4 ulasan). Absennya keluhan bukan bukti absennya kebutuhan, karena pengulas Play jarang menulis tentang fitur yang tidak mereka ketahui ada.
- Harga/paywall dan "janji tidak sesuai" adalah keluhan komersial yang sering di paket bimbel (JadiASN, JadiBUMN, BisaCPNS), bukan di app gratis berbasis iklan. Ini menunjukkan sensitivitas harga setelah pembelian, bukan sebelum pembelian (inferensi).
- Admin/CS yang tidak membalas muncul di 7 ulasan bintang 1 sampai 3 dan berpasangan dengan keluhan pembayaran; sebaliknya admin yang responsif muncul sebagai pujian (Bagian 3). Ini menyiratkan bahwa dukungan pengguna menjadi pembeda pada pelanggan berbayar (inferensi).

### Gaps
- Aturan kata kunci tidak mengenali ironi, negasi kompleks, atau bahasa gaul; pengecekan dengan membaca hanya mencakup bintang 1 sampai 3 inti dan sekitar 30 ulasan bintang 4 sampai 5 yang dipilih, bukan sampel acak atas semua ulasan.
- 77 ulasan bintang 1 sampai 3 tidak punya kategori; sebagian bisa jadi keluhan UX atau pendaftaran yang tidak ada kategorinya.
- Ulasan hanya berasal dari pengguna yang menulis; tidak ada data demografi, tidak tahu apakah pengulas pelamar CPNS atau BUMN, dan tidak ada ukuran keluhan dalam populasi penuh pengguna.
- Balasan pengembang tidak dianalisis.

## 3. Pujian: fitur yang dihargai pengguna

### Takeaway
Dari 1.514 ulasan inti bintang 4 sampai 5, 764 (50%) cocok minimal satu kategori pujian; sisanya umum ("bagus", "semoga lolos"). Pujian terbanyak bersifat umum (membantu persiapan/lolos, 38,5%), diikuti tampilan/kemudahan (7,3%), soal mirip ujian asli/HOTS (6,2%), video materi dan mentor (4,9%), banyaknya soal (4,2%), harga gratis/murah (3,7%), dan pembahasan yang jelas (2,7%). Fitur statistik/ranking dan simulasi CAT disebut jarang (0,9% dan 1,3%).

### Cited Findings
Frekuensi kata kunci pada 1.514 ulasan inti bintang 4 sampai 5 (13 aplikasi, 2024 sampai 2026):

| Kategori pujian | n | % |
|---|---|---|
| Membantu persiapan / membantu lolos (umum) | 583 | 38,5% |
| Tampilan/kemudahan penggunaan | 110 | 7,3% |
| Soal mirip ujian asli / HOTS | 94 | 6,2% |
| Video materi / mentor / kelas | 74 | 4,9% |
| Banyak soal / variasi | 64 | 4,2% |
| Gratis / murah | 56 | 3,7% |
| Pembahasan jelas | 41 | 2,7% |
| Simulasi mirip CAT | 19 | 1,3% |
| Statistik / ranking | 13 | 0,9% |
| Admin / dukungan responsif | 12 | 0,8% |

- Contoh parafrase pujian: "pembahasan mudah dipahami", "TO gratis tiap minggu dan paling tidak pelit" (IDCPNS), "soal HOTS dan update", "admin cepat menangani kendala" (JadiASN), "bisa dipakai offline" (Bimbel & Simulasi CAT dan Simulasi CAT CPNS), "peringkat nilai peserta dan pembahasan memuaskan".
- Pada Analitica (campuran UTBK/CPNS; 244 ulasan bintang 4 sampai 5) sekitar 91 ulasan (37%, hitungan kata kunci "AI/Ica/Minca" kasar) menyebut asisten AI penjelas sebagai fitur favorit; kata kunci pembahasan jelas, banyak soal, dan video materi juga muncul. Sebagian besar konteksnya UTBK/ujian mandiri, jadi bukti ini tidak langsung untuk SKD. Sumber: [Analitica di Play](https://play.google.com/store/apps/details?id=co.xomodo.analitica&hl=id&gl=id), dibaca 14 Nov 2026.

### Inferences
- Pembahasan jelas dan soal mirip asli disebut sebagai alasan puas, tetapi dengan proporsi kecil; pujian terbesar bersifat umum. Sinyal "AI membantu menjelaskan" muncul kuat di Analitica, yang mendukung bahwa penjelasan AI diterima pengguna (di konteks UTBK), bukan bahwa ia meningkatkan nilai (inferensi).
- Pujian atas statistik dan ranking sedikit (0,9%); tidak ada bukti ulasan bahwa leaderboard dicari.

### Gaps
- Kategori pujian lebar dan tumpang tindih ("hots" dan "bagus" dihitung sebagai mirip asli); gunakan sebagai urutan kasar, bukan ukuran.
- Burst ulasan bintang 5 pada beberapa app (EKSAM 91 ulasan pada Januari 2026, Apta 53 pada Maret 2025, IDCPNS sekitar 50 pada Januari/Februari 2024 dan Februari 2026) bisa organik atau didorong promosi; tidak dikeluarkan dari angka di atas, tetapi menggembungkan pujian umum untuk app tersebut.

## 4. Permintaan fitur yang disebutkan eksplisit

### Takeaway
Sekitar 97 ulasan terjaring pola permintaan/saran (di luar burst Tryout.id); setelah dibaca, permintaan eksplisit yang nyata hanya sekitar 20 sampai 25 dan tersebar. Permintaan terbanyak: **soal lebih mutakhir/HOTS/penalaran dan lebih sulit** (4), **riwayat hasil ujian** (2), **timer untuk latihan soal** (2, bertentangan dengan 1 permintaan menghapus timer latihan), dan **iklan dikurangi** (2).

### Cited Findings
Permintaan eksplisit (jumlah ulasan dalam kurung; semua 2024 sampai 2026 dari halaman Play di Bagian 1, dibaca 14 Nov 2026):
- Soal lebih baru/HOTS/penalaran, TIU dan TWK diperbarui, TO lebih sulit mendekati simulasi BKN (4, di antaranya Simulasi CAT CPNS dan Bimbel & Simulasi CAT).
- Riwayat ujian yang sudah dikerjakan agar bisa dibandingkan (2, Simulasi CAT CPNS).
- Timer pada latihan soal agar mirip tes sungguhan (2, Bimbel & Simulasi CAT); sebaliknya 1 pengulas meminta latihan soal tanpa waktu dan hanya tryout yang berwaktu.
- Drilling soal (1; pengulas menyebut sudah mengisi form saran di grup tetapi belum ada); pengguna Analitica memuji fitur drill yang sudah ada.
- Video pembahasan pada TO (1, Al Faiz); materi literasi (1); tryout bisa diulang selama masa premium (1); gabungkan TIU, TWK, dan TKP dalam satu sesi (1, Bimbel & Simulasi CAT; satu ulasan lain meminta simulasi SKB per kategori digabung dengan satu timer); mode offline lebih baik (1; dua ulasan lain memuji fitur offline yang sudah ada).
- Tambah bank soal (1); SKB untuk instansi tertentu (2: Kementan dan Kemenkumham); filter jurusan untuk lowongan (1); mode gelap (1); pembayaran lewat dompet digital (1 pengulas yang gagal membayar); diskon (1); iklan dikurangi (2, Soal CPNS 2026 - Latihan CAT dan Simulasi CAT CPNS); ruang untuk melaporkan soal (1 yang tidak menemukannya, 1 yang sudah melapor tanpa tanggapan).

### Inferences
- Permintaan fitur belajar terarah (drill, riwayat, kelemahan per pola) ada tetapi kecil (sekitar 4 ulasan). Itu bukan penolakan terhadap loop adaptif; ia menunjukkan bahwa sebagian besar pengguna tidak meminta apa yang tidak mereka lihat tersedia (inferensi, belum diuji).
- Pembaruan konten ke format ujian terkini (penalaran, figural) adalah permintaan terbanyak di kelompok fitur belajar.

### Gaps
- Pola pencarian saran mungkin melewatkan permintaan yang ditulis tanpa kata saran. Tidak ada klasifikasi manual penuh atas 2.331 ulasan.

## 5. Batas dan temuan tambahan

### Takeaway
Sampel adalah campuran "terbaru" dan "paling relevan" sehingga tidak acak; hasilnya cocok untuk memberi peringkat tema, tidak untuk memperkirakan tingkat keluhan populasi.

### Cited Findings
- Urutan "paling relevan" Play menarik ulasan lama; hanya ulasan 2024 sampai 2026 yang dipakai (2.331 dari 3.384).
- Tryout.id: 278 dari 287 ulasan bertanggal Oktober 2025 hampir seragam, bintang 5, dan umumnya menyebut UKOM/UTBK/CPNS sekaligus; dikeluarkan (diduga promosi, **Asumsi (belum pasti)**: diverifikasi dengan membaca balasan pengembang atau memeriksa akun pengulas, yang tidak dilakukan).
- Aplikasi campuran Analitica (80 dari 256 ulasan menyebut istilah UTBK/UM/TKA, hanya 12 menyebut istilah CPNS/SKD/BUMN), Al Faiz, dan Tryout.id dipisahkan dari angka inti.

### Inferences
- Karena burst bintang 5 dan pengulas yang cenderung menulis saat puas atau marah, proporsi keluhan sebenarnya dalam seluruh pengguna tidak diketahui.

### Gaps
- Tidak ada pengambilan ulang untuk menguji kestabilan sampel; pembaruan Play dapat mengubah ulasan yang tampil. Skrip dan data mentah ada di folder sementara sesi, bukan di repo.
