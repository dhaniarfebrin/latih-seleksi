---
name: Latih Seleksi
description: Latihan seleksi ASN yang tenang dan padat, bergaya buku catatan belajar dengan stabilo kuning.
colors:
  paper: "#f4f0e8"
  card: "#fbf9f4"
  ink: "#17191c"
  muted: "#5a5e5b"
  line: "#d8d2c4"
  highlighter: "#f7d547"
  correct: "#1b6446"
  correct-wash: "#dceadf"
  wrong: "#ae2f1d"
  wrong-wash: "#f5d9d0"
  caution: "#7a5000"
  caution-wash: "#f4e4b4"
  paper-night: "#121413"
  card-night: "#1b1e1c"
  ink-night: "#eeebe1"
  muted-night: "#a5a99f"
  line-night: "#333730"
  edge-night: "#686d63"
  highlighter-night: "#e6c73a"
  correct-night: "#78d3a2"
  correct-wash-night: "#173226"
  wrong-night: "#ff917c"
  wrong-wash-night: "#3b1e18"
  caution-night: "#f0ce6a"
  caution-wash-night: "#36300f"
typography:
  display:
    fontFamily: "Bricolage Grotesque, Trebuchet MS, system-ui, sans-serif"
    fontSize: "34px"
    fontWeight: 800
    lineHeight: 1.05
    letterSpacing: "-0.025em"
  headline:
    fontFamily: "Bricolage Grotesque, Trebuchet MS, system-ui, sans-serif"
    fontSize: "26px"
    fontWeight: 700
    lineHeight: 1.1
    letterSpacing: "-0.02em"
  title:
    fontFamily: "Bricolage Grotesque, Trebuchet MS, system-ui, sans-serif"
    fontSize: "17px"
    fontWeight: 700
    lineHeight: 1.45
  body:
    fontFamily: "IBM Plex Sans, system-ui, -apple-system, sans-serif"
    fontSize: "16px"
    fontWeight: 400
    lineHeight: 1.45
  question:
    fontFamily: "IBM Plex Sans, system-ui, -apple-system, sans-serif"
    fontSize: "18px"
    fontWeight: 400
    lineHeight: 1.45
  label:
    fontFamily: "IBM Plex Mono, ui-monospace, Menlo, monospace"
    fontSize: "12px"
    fontWeight: 400
    lineHeight: 1.45
    letterSpacing: "0.06em"
  numeric:
    fontFamily: "IBM Plex Mono, ui-monospace, Menlo, monospace"
    fontSize: "16px"
    fontWeight: 600
    lineHeight: 1.45
    fontFeature: "tnum"
rounded:
  cell: "10px"
  field: "12px"
  control: "14px"
  panel: "18px"
  nav: "28px"
  round: "50%"
spacing:
  xs: "4px"
  sm: "8px"
  md: "12px"
  lg: "16px"
  xl: "24px"
  gutter: "18px"
components:
  button-primary:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.paper}"
    rounded: "{rounded.control}"
    height: "54px"
    padding: "0 20px"
  button-highlight:
    backgroundColor: "{colors.highlighter}"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    height: "54px"
    padding: "0 20px"
  button-ghost:
    backgroundColor: "transparent"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    height: "54px"
    padding: "0 20px"
  button-danger:
    backgroundColor: "{colors.wrong}"
    textColor: "{colors.paper}"
    rounded: "{rounded.control}"
    height: "54px"
    padding: "0 20px"
  panel:
    backgroundColor: "{colors.card}"
    textColor: "{colors.ink}"
    rounded: "{rounded.panel}"
    padding: "16px"
  option:
    backgroundColor: "{colors.card}"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    height: "54px"
    padding: "8px 14px"
  option-correct:
    backgroundColor: "{colors.correct-wash}"
    textColor: "{colors.correct}"
  option-wrong:
    backgroundColor: "{colors.wrong-wash}"
    textColor: "{colors.wrong}"
  option-picked:
    backgroundColor: "{colors.highlighter}"
    textColor: "{colors.ink}"
  chip-caution:
    backgroundColor: "{colors.caution-wash}"
    textColor: "{colors.caution}"
    rounded: "{rounded.field}"
    padding: "6px 12px"
  segment:
    backgroundColor: "{colors.card}"
    textColor: "{colors.ink}"
    rounded: "{rounded.field}"
    height: "48px"
    padding: "8px 12px"
  segment-selected:
    backgroundColor: "{colors.highlighter}"
    textColor: "{colors.ink}"
  input:
    backgroundColor: "{colors.card}"
    textColor: "{colors.ink}"
    rounded: "{rounded.field}"
    height: "48px"
    padding: "8px 12px"
  nav-pill:
    backgroundColor: "{colors.card}"
    textColor: "{colors.muted}"
    rounded: "{rounded.nav}"
    padding: "6px"
  nav-item-active:
    backgroundColor: "{colors.highlighter}"
    textColor: "{colors.ink}"
    rounded: "{rounded.panel}"
    height: "52px"
---

# Design System: Latih Seleksi

## Overview

**Creative North Star: "The Study Notebook"**

Sistem ini terasa seperti buku catatan belajar yang rapi: kertas krem hangat, tinta arang pekat, dan satu stabilo kuning untuk hal yang sedang penting. Angka dan label memakai huruf mono seperti tulisan tangan teknis di margin. Kartu hanya dibatasi garis tipis. Tidak ada bayangan yang menghias, tidak ada gradien, tidak ada maskot. Suasananya tenang, padat, dan hangat, sejalan dengan suara informal "kamu" di PRODUCT.md.

Kepadatan sedang: satu kolom sempit (480 px) untuk alur fokus seperti login dan drill, kolom lebar (1040 px) untuk beranda dan dasbor, kolom baca (640–720 px) untuk Tentang, kartu soal, dan ringkasan, dengan satu pita waktu di atas kartu soal saat ujian dan drill. Target sentuh besar (48–54 px) berlaku di ponsel dan laptop, karena keduanya pengguna utama.

Penolakan visual yang dikonfirmasi pemilik (2026-11-17): bukan aplikasi edukasi yang ramai. Tidak ada maskot, gamifikasi mencolok, gradien ungu-biru ala startup edtech, maupun tampilan portal pemerintah.

**Catatan implementasi:** `web/src/styles.css` mendeklarasikan Bricolage Grotesque, IBM Plex Sans, dan IBM Plex Mono, dan ketiganya dimuat self-host lewat paket fontsource di `web/src/main.tsx` (Bricolage Grotesque variabel `wght.css`; IBM Plex Sans latin 400/600/700; IBM Plex Mono latin 400/600; `font-display: swap`). Fallback (Trebuchet MS, system-ui) hanya tampil sesaat sebelum font siap. Belum ada preload atau penyesuaian metrik fallback.

**Key Characteristics:**
- Kertas krem dan tinta arang; satu aksen kuning stabilo.
- Datar: garis 1–1,5 px menggantikan bayangan.
- Status selalu ikon + teks, tidak pernah hanya warna.
- Angka dalam mono tabular; judul dalam display berat 800.
- Mode gelap otomatis lewat `prefers-color-scheme`, bukan toggle.
- Komponen besar dan tegas: tombol 54 px, sudut 14–18 px.

## Colors

Palet buku catatan: dua krem hangat, satu tinta, satu stabilo, dan tiga warna umpan balik (benar, salah, hati-hati) yang masing-masing punya warna isi dan warna teks. Semua didefinisikan sebagai custom property di `:root`, dengan set gelap menimpa lewat `prefers-color-scheme: dark`.

### Primary
- **Tinta Arang** (`#17191c`; `#eeebe1` di mode gelap): teks, tombol utama, garis fokus 3 px, progress, dan batas komponen terpilih. Tinta adalah warna aksi; tombol utama berwarna tinta, bukan kuning.
- **Stabilo Kuning** (`#f7d547`; `#e6c73a` di mode gelap): penanda hal yang aktif atau penting: tab navigasi aktif, segmen terpilih, opsi yang dipilih, tombol sorot (`.btn.hi`), avatar, dan `mark` pada nama atau frasa kunci. Selalu dengan teks tinta di atasnya.

### Secondary
- **Hijau Tinta Benar** (`#1b6446`, isi `#dceadf`; mode gelap `#78d3a2` / `#173226`): jawaban benar, soal terjawab, status sukses.
- **Merah Koreksi** (`#ae2f1d`, isi `#f5d9d0`; mode gelap `#ff917c` / `#3b1e18`): jawaban salah, galat, tindakan merusak (hapus akun), timer hampir habis.

### Tertiary
- **Cokelat Peringatan** (`#7a5000`, isi `#f4e4b4`; mode gelap `#f0ce6a` / `#36300f`): chip jenis kesalahan, soal ditandai, label perkiraan.

### Neutral
- **Kertas Krem Hangat** (`#f4f0e8`; `#121413` di mode gelap): latar halaman dan bilah sticky.
- **Kartu Putih Gading** (`#fbf9f4`; `#1b1e1c`): permukaan kartu, panel, opsi, input, dialog.
- **Abu Tinta Pudar** (`#5a5e5b`; `#a5a99f`): teks sekunder, eyebrow, ikon nav tidak aktif.
- **Garis Kertas** (`#d8d2c4`; `#333730`): pembatas dan garis kartu. Di mode gelap garis tepi kontrol (`--edge`) lebih terang, `#686d63`, agar batas kontrol tetap terlihat.

### Named Rules
**The One Highlighter Rule.** Stabilo kuning menandai satu hal yang sedang aktif atau penting per area; ia bukan dekorasi. Di mode gelap, `mark` berubah dari garis bawah setengah tinggi menjadi chip solid dengan teks gelap, karena teks terang di atas kuning tidak terbaca.

**The Never Color Alone Rule.** Benar, salah, terjawab, ditandai, dan peringatan selalu disertai ikon dan teks (konvensi NOTES.md 2026-11-15). Warna hanya memperkuat.

## Typography

**Display Font:** Bricolage Grotesque (dengan Trebuchet MS, system-ui)
**Body Font:** IBM Plex Sans (dengan system-ui, -apple-system)
**Label/Mono Font:** IBM Plex Mono (dengan ui-monospace, Menlo)

**Character:** Judul grotesque yang tebal dan agak playful bertemu isi Plex yang netral dan label mono seperti anotasi teknis. Pasangan ini menjaga nuansa belajar yang ramah tanpa jadi kekanak-kanakan. Lihat catatan pemuatan font di Overview.

### Hierarchy
- **Display** (800, 34 px; 52 px untuk hero landing di ≥900 px, line-height 1,05, letter-spacing −0,025em): judul halaman (`h1.d`), sering dengan nama pengguna dalam `mark`.
- **Headline** (700, 26 px, line-height 1,1, −0,02em): judul bagian di landing (`.lp h2`).
- **Title** (700, 17 px; 22 px untuk judul dialog): judul panel dan subbagian (`.sub`, `.lp-feats h3`).
- **Body** (400, 16 px, line-height 1,45): teks isi. Teks lead landing 18 px dengan lebar maksimum 56ch.
- **Question** (400, 18 px): batang soal (`.stem`).
- **Label** (mono 12 px, +0,06em, huruf kapital): label kelompok kecil (bukan di atas judul halaman). Teks bantu kecil 13–14 px bobot 600 untuk label field, badge, dan tag.
- **Numeric** (mono 600, tabular-nums): skor, timer, persentase, nomor soal (`.num`).

### Named Rules
**The Mono For Numbers Rule.** Setiap angka yang dibandingkan atau dihitung (skor, waktu, persentase, nomor) memakai mono tabular agar rata dan mudah dipindai.

## Layout

Satu stylesheet global (`web/src/styles.css`), tanpa framework CSS. Halaman adalah `.screen` kolom tunggal dengan gap 16 px dan padding 22 px 18 px; `.narrow` (480 px) untuk alur fokus, `.wide` (1040 px) untuk beranda, dasbor, dan riwayat. Satu breakpoint utama: **900 px**.

- **Di bawah 900 px (ponsel):** navigasi berupa pil melayang di bawah (maks 420 px, 5 tab: Beranda, Drill, Simulasi, Riwayat, Akun), menyisakan padding bawah 96 px. Pil disembunyikan saat ujian berjalan (`.app:has(.exam)`) agar tidak tersentuh tak sengaja dan posisi soal tidak hilang.
- **Mulai 900 px (laptop):** navigasi pindah ke bilah atas sticky (maks 1040 px), tab sejajar dengan Akun di kanan. Panel pilihan jadi dua kolom bergaris pemisah vertikal. Saat ujian, pita waktu menempel di atas dengan tab subtes di baris jam. Baris riwayat jadi grid empat kolom.
- Kelas `.mo` / `.dk` menampilkan elemen hanya di ponsel atau hanya di desktop.
- Ritme jarak: 4, 8, 10, 12, 16, 24 px. Landing: gap bagian 40 px (56 px di desktop).
- Target sentuh: tombol 54 px, segmen dan input 48 px, sel penggaris nomor 34–36 px (di atas minimum WCAG 24 px), item nav 52 px.

## Elevation & Depth

Datar secara default. Kedalaman datang dari lapisan tonal (kertas → kartu) dan garis 1–1,5 px, bukan dari bayangan. Satu-satunya elemen yang melayang adalah pil navigasi ponsel, dan dialog yang diberi backdrop gelap. Bayangan tidak muncul sebagai respons status.

### Shadow Vocabulary
- **Pil melayang** (`box-shadow: 0 6px 24px rgb(0 0 0 / .18)`): hanya pil navigasi ponsel. Bila `backdrop-filter` didukung, latarnya jadi kartu 88% transparan dengan blur 12 px.
- **Backdrop dialog** (`rgb(0 0 0 / .55)`): lapisan di belakang dialog.

### Named Rules
**The Flat-By-Default Rule.** Kartu, panel, opsi, dan input tidak pernah diberi bayangan. Batas dan isi warna yang membedakan.

## Shapes

Sudut membulat besar dan konsisten: sel 10 px, field 12 px, kontrol (tombol, opsi, baris) 14 px, panel dan dialog 18 px, pil navigasi 28 px. Lingkaran penuh (50%) untuk avatar dan huruf pilihan A–E. Garis kontrol 1,5 px; garis kartu dan pembatas 1 px. Fokus: kontur 3 px berwarna tinta dengan offset 2 px (`:focus-visible`).

## Components

Karakter: besar, tegas, mudah disentuh, dengan garis tinta yang jelas.

### Buttons
- **Shape:** sudut 14 px, tinggi minimum 54 px, lebar penuh, garis 1,5 px.
- **Primary:** latar tinta, teks kertas (`.btn`).
- **Highlight:** latar stabilo, teks tinta, garis tinta (`.btn.hi`), untuk aksi lanjut seperti "Lanjutkan simulasi".
- **Ghost:** transparan, garis tepi; hover menggelapkan garis ke tinta.
- **Hierarki:** satu aksi utama per layar (terisi: tinta atau stabilo); aksi lain memakai ghost (`.btn.ghost`) atau tautan (`.more`). Varian `.btn.sm` = lebar menyesuaikan, tinggi 44 px, untuk aksi bantu di dalam formulir seperti "Pola terlemah".
- **Danger:** latar merah koreksi, dipakai untuk hapus akun.
- **Disabled:** opacity 0,55, kursor `not-allowed`.
- Tautan bisa memakai kelas `.btn` (`a.btn` tanpa garis bawah).

### Chips
- **Caution chip:** blok isi cokelat peringatan, sudut 12 px, teks 14 px bobot 600, misalnya "Jenis salah: menjumlahkan persen langsung". Dipakai untuk label penyebab kesalahan dan saran latihan.

### Cards / Containers
- **Corner Style:** 18 px.
- **Background:** kartu putih gading di atas kertas krem.
- **Shadow Strategy:** tanpa bayangan (lihat Elevation).
- **Border:** 1 px garis kertas. Varian `.panel.danger` memakai batas dan judul merah koreksi.
- **Internal Padding:** 16 px, gap 4 px antar isi.

### Answer Options
- Baris tombol 54 px dengan huruf pilihan bulat 30 px di kiri dan tag status di kanan. Setelah dijawab: benar = isi hijau + teks hijau + ikon centang + teks "Benar"; salah = isi merah + teks merah + ikon silang + "Pilihan Anda". Terpilih (sebelum submit) = stabilo kuning.

### Segmented Choice
- Grid radio berbentuk kartu 48 px (`.seg`), terpilih = stabilo kuning dengan garis tinta. Radio asli disembunyikan tetapi tetap fokus lewat `:focus-within`.
- **Daftar yang bisa dipilih (`.seg.pick`):** kotak 20 px di kiri bila boleh pilih lebih dari satu (centang), titik bila hanya satu (`.seg.pick.radio`). Kuning saja bukan penanda; legend juga menyebut "pilih satu" atau "pilih satu atau lebih".
- **Chip mengalir (`.seg.tags`):** lebar mengikuti teks dan membungkus, agar baris terakhir tidak meninggalkan sel yatim. `.seg.c4` = empat kolom sama lebar untuk empat opsi pendek.

### Inputs / Fields
- **Style:** tinggi 48 px, garis 1,5 px, sudut 12 px, latar kartu, label di atas bobot 600 14 px.
- **Focus:** kontur tinta 3 px dengan offset 2 px.
- **Error:** blok `.err` dengan isi merah koreksi, sudut 12 px, `role="alert"`.

### Navigation
- **Ponsel:** pil melayang berisi 5 tab vertikal (ikon 20 px di atas label 11 px bobot 600), item 52 px, tab aktif berisi stabilo kuning.
- **Desktop:** bilah sticky dengan tab horizontal (ikon + teks), sudut 10 px, Akun didorong ke kanan; merek "Latih Seleksi" dalam display 800 20 px.
- Hanya satu tab aktif; warna item tidak aktif abu tinta pudar.

### Time Band and Ruler
- **Pita (`.band`):** kartu bergaris di atas kartu soal, sticky saat ujian (`.pin`). Berisi jam mono besar (28 px ponsel, 44 px desktop; merah bila ≤60 detik), tab subtes (`.tab`, aktif = isi tinta), progress, dan penggaris nomor. Di ponsel tab dan penggaris masuk ke `details` "Nomor soal (x/y terjawab)". Drill memakai pita yang sama dengan jam berjalan dan target.
- **Penggaris (`Ruler`, `shared/components/Ruler.tsx`):** satu sel mono 36 px per soal. Ujian: terjawab = isi tinta + centang, ditandai = garis putus-putus + ikon bendera, kosong = kartu. Drill dan ringkasan: benar = isi tinta + centang, salah = merah koreksi + ikon silang. Soal saat ini = stabilo kuning. Dengan `onPick` sel menjadi tombol; tanpa itu hanya penanda (`role="img"`).

### Question Card
- `.qcard`: kartu 720 px di tengah, sudut 18 px, berisi batang soal 19 px bobot 600, opsi, dan baris Sebelumnya / Tandai ragu / Berikutnya. Opsi terpilih = stabilo; setelah dijawab di drill, benar hijau dan salah merah, dengan panel langkah di bawah kartu.

### Result Rows
- `.sheet-rows` / `.sheet-row`: nama, batang isi tinta (lebar = persentase), angka mono, dan satu baris keterangan. Baris terlemah memakai `.weak`. Dipakai di Dasbor (per pola), hasil Ujian (per subtes), dan landing. Tren sesi memakai `TrendBars` (kolom per sesi, minimal 3 titik).

### Reading Page
- `.doc` (640 px): catatan penting sebagai blok stabilo, bagian dipisah garis 1 px, tanpa kartu bertumpuk. Dipakai Tentang.

### Dialog
- `dialog.dlg`: lebar maksimum 440 px, sudut 18 px, kartu bergaris, backdrop gelap 55%. Judul display 22 px.

## Do's and Don'ts

### Do:
- **Do** pakai custom property di `:root` untuk warna, jangan menulis hex baru di komponen; set gelap menimpa lewat `prefers-color-scheme`.
- **Do** tampilkan status dengan ikon + teks, bukan hanya warna; ikon dekoratif `aria-hidden` dan setiap kontrol punya label teks.
- **Do** pakai stabilo kuning untuk satu hal aktif per area, dengan teks tinta di atasnya.
- **Do** set angka dalam mono tabular (`.num`).
- **Do** jaga target sentuh ≥48 px, tombol utama 54 px, di ponsel maupun laptop.
- **Do** hormati `prefers-reduced-motion`; gerak hanya transisi latar 0,15 detik pada tab nav.
- **Do** sembunyikan navigasi melayang saat ujian berjalan di ponsel.
- **Do** tampilkan disclaimer independen di footer dan jaga bahasa informal "kamu".

### Don't:
- **Don't** pakai bayangan pada kartu, panel, opsi, atau input; kedalaman dari garis dan lapisan warna.
- **Don't** tambahkan maskot, gamifikasi mencolok, atau gradien ungu-biru ala startup edtech.
- **Don't** tiru tampilan portal pemerintah atau pakai nama/logo BKN, KemenPANRB, atau FHCI.
- **Don't** pakai teks terang di atas stabilo kuning; di mode gelap `mark` harus chip solid bergelap.
- **Don't** menandai benar/salah hanya dengan warna.
- **Don't** menampilkan testimoni, angka pengguna, atau jaminan kelulusan tanpa sumber (PRODUCT.md, Product Principles).
- **Don't** menaruh eyebrow di atas judul halaman yang hanya mengulang judulnya; eyebrow hanya boleh jadi label kelompok (misalnya subtes di mulai drill).
