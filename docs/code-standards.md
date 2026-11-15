# Code standards

Aturan ini berlaku untuk kode baru dan refactor. Kode lama yang belum sesuai dirapikan bertahap; jangan merapikan file yang tidak sedang disentuh tugasnya.

Prinsipnya SOLID, tetapi diterapkan dengan idiom tiap bahasa: Go bukan bahasa OOP, jadi tidak ada lapisan "Clean Architecture" penuh, repository per entitas, atau interface yang hanya punya satu implementasi tanpa alasan (tes atau sumber kedua). Lebih sedikit kode yang jelas lebih baik dari banyak lapisan.

## 1. Prinsip umum

| Prinsip | Artinya di repo ini |
|---|---|
| SRP | Satu package/folder = satu fitur atau satu tanggung jawab. Satu file = satu topik (handler drill, handler exam, middleware). File > ~400 baris (di luar tes dan generator soal) adalah tanda perlu dipecah. |
| OCP | Variasi ditambah lewat data atau implementasi baru, bukan `if/switch` baru di inti. Contoh yang sudah benar: `gen.Family` + `catalog`, `exam.ItemSource`, spesifikasi ujian di `spec/data/*.json`. |
| LSP | Setiap implementasi interface memenuhi kontrak yang tertulis di komentar interface (mis. `Generate` deterministik per seed). Kontrak diuji, bukan diasumsikan. |
| ISP | Interface kecil, didefinisikan di sisi **pemakai**, hanya berisi method yang dipakai pemakai itu. |
| DIP | Logika domain (`drill`, `exam`, `scoring`) tidak bergantung pada tipe konkret infrastruktur (`*store.Store`, `pgx`, `net/http`, Firebase). Dependensi disuntik lewat konstruktor dan dirakit hanya di `cmd/*/main.go` (atau `App.tsx`/provider di web). |
| DRY | Helper yang sama di dua package dipindah ke satu tempat. Tiga baris yang mirip belum tentu perlu diabstraksi; logika domain yang sama (validasi jawaban, snapshot soal) wajib satu sumber. |

Invarian produk di `CLAUDE.md` (kunci jawaban, spesifikasi berversi, tes generator) tetap berlaku di atas semua aturan ini.

## 2. Backend Go (`api/`)

### 2.1 Struktur

```
api/
  cmd/<binary>/main.go     composition root: baca env, rakit dependensi, jalankan. Tanpa logika domain.
  db/                      migrasi SQL (embed)
  internal/
    <fitur>/               drill, exam, bank, spec, scoring, auth
      <fitur>.go           tipe domain, Service, konstruktor
      ports.go             interface yang dibutuhkan Service (store, source, clock)
      view.go              tipe JSON yang dikirim ke client (jika banyak)
      errors.go            error sentinel fitur (jika lebih dari 2-3)
    gen/                   kontrak generator + satu package per keluarga soal
    question/              tipe bersama soal beku: snapshot, opsi, validasi jawaban A-E, seed acak
    httpx/                 transport HTTP saja
      router.go            NewHandler + daftar route
      middleware.go        request id, CORS, auth
      respond.go           writeJSON, decode, pemetaan error -> status
      drill.go, exam.go    handler per fitur
    store/                 PostgreSQL saja
      store.go             Open, Close, helper
      migrate.go
      attempt.go, exam.go, user.go   query per agregat
```

- Package-by-feature, bukan package-by-layer (tidak ada `models/`, `services/`, `controllers/`).
- Nama package satu kata, huruf kecil, tanpa underscore (`deretangka`, bukan `deret_angka`).
- Package `internal/` tidak saling impor secara melingkar; arah dependensi: `httpx` -> fitur -> `question`/`gen`/`spec`; `store` mengimplementasikan interface fitur, fitur tidak mengimpor `store`.

### 2.2 Dependensi dan interface

```go
// internal/drill/ports.go
type Store interface {
    CreateAttempt(ctx context.Context, a Attempt, items []Item) (Attempt, error)
    GetAttempt(ctx context.Context, id string) (Attempt, []Item, []Response, error)
    SaveResponse(ctx context.Context, id string, r Response) (Response, error)
    MarkSubmitted(ctx context.Context, id string) (time.Time, error)
}

func New(st Store, families []gen.Family) *Service
```

- Service menerima interface, mengembalikan struct konkret ("accept interfaces, return structs").
- Tidak ada pemanggilan dependensi global di dalam konstruktor (mis. `catalog.All()` di `drill.New`); dependensi dikirim dari `main.go`.
- Waktu dan keacakan disuntik (`now func() time.Time`, sumber seed) agar bisa diuji tanpa DB dan tanpa `sleep`.
- Tipe data yang lewat batas store-domain didefinisikan di package domain atau `question`, bukan di `store`. Konstanta domain (status exam `in_progress/submitted/expired`) milik package `exam`.
- Tes service memakai fake in-memory yang memenuhi interface `Store`. Tes `store` sendiri tetap memakai Postgres sungguhan (`DATABASE_URL`).

### 2.3 Fungsi dan parameter

- Lebih dari 4 parameter selain `ctx` -> pakai struct input (`exam.AnswerInput{Pos, Answer, TimeMS, Flagged, Seq}`).
- Validasi input di service (bukan di handler); handler hanya decode, panggil, encode.
- `ctx context.Context` selalu parameter pertama untuk operasi I/O.

### 2.4 Error

- Error sentinel per fitur (`ErrNotFound`, `ErrInvalid`, ...) di-wrap dengan `%w`; pembuat pesan memakai `fmt.Errorf("%w: ...", ErrInvalid, ...)`.
- Error yang sama maknanya di beberapa fitur boleh berbagi satu sentinel di `question` (atau package `apperr` kecil) supaya `httpx` tidak harus memetakan `drill.ErrX` dan `exam.ErrX` satu per satu.
- `store` menerjemahkan error driver (`pgx.ErrNoRows`, SQLSTATE `22P02`) menjadi `ErrNotFound`; `pgx` tidak bocor ke luar `store`.
- Galat internal (generator gagal, DB) -> 500 dan di-log dengan request id; galat input -> 400. Jangan bungkus galat internal sebagai `ErrInvalid`.
- Tidak ada `panic` di jalur request.

### 2.5 HTTP

- `httpx` tidak berisi logika domain. Bentuk JSON request didefinisikan di handler, bentuk response (`...View`) di package fitur.
- Route terdaftar di satu tempat (`router.go`), middleware dibungkus berurutan dan jelas.
- Body dibatasi (`maxBodyBytes`), field tak dikenal ditolak.
- Token dan data pribadi tidak pernah di-log.

### 2.6 Gaya dan tes

- `gofmt`, `go vet` wajib bersih (perintah lengkap di `NOTES.md` Setup).
- Setiap identifier yang diekspor punya komentar doc yang diawali namanya.
- Komentar menjelaskan *kenapa* dan kontrak, bukan mengulang kode.
- Tes table-driven; nama `TestXxx_kondisi`. Tes generator mengikuti konvensi di `NOTES.md` (kunci dihitung ulang dari teks soal).

## 3. Frontend React + TypeScript (`web/`)

### 3.1 Struktur

```
web/src/
  main.tsx
  app/                    App.tsx (routes), Shell, RequireAuth, provider global
  features/
    drill/                DrillPage.tsx, components/, useDrill.ts, session.ts (+ .test.ts), api.ts, types.ts
    exam/                 ExamPage.tsx, StartPage.tsx, components/ (Running, Result, ExamRow), exam.ts (+ .test.ts), api.ts, types.ts
    auth/                 AuthContext.tsx, auth.ts, LoginPage.tsx
    account/, history/, home/, about/
  shared/
    api/client.ts         fetch wrapper, ApiError, offset jam server, header Authorization
    hooks/                useLoad.ts
    lib/                  format waktu, helper murni
    ui/                   komponen presentasional umum (Button, Alert, Segmented)
  styles.css
```

- Satu file = satu komponen yang diekspor. Komponen kecil privat boleh satu file dengan induknya hanya jika < ~40 baris.
- Halaman (route) diberi akhiran `Page`. Komponen presentasional tidak memanggil API.
- Feature tidak mengimpor isi feature lain; yang dipakai bersama pindah ke `shared/`.
- Impor antarfolder lewat path relatif yang pendek; tambahkan alias `@/` di `tsconfig`/`vite.config` saat struktur ini dibuat.

### 3.2 Tanggung jawab

- **Komponen**: render dan event. Tidak berisi `fetch`, tidak menyusun URL.
- **Hook** (`useXxx`): state server dan efek (memuat, mengirim jawaban, submit). Menangani unmount (flag `live` atau `AbortController`).
- **Reducer / fungsi murni** (`session.ts`, `exam.ts`): logika transisi state, diuji dengan Vitest tanpa DOM.
- **api.ts per feature**: fungsi endpoint bertipe; memakai `shared/api/client.ts`. Token diambil client sendiri (DIP), bukan dioper sebagai argumen pertama di setiap panggilan.

### 3.3 TypeScript

- `strict` tetap aktif. Tidak ada `any`; pakai `unknown` lalu persempit.
- Tipe response API ditulis sekali per feature (`types.ts`) dan mencerminkan `...View` di Go; komentar menyebut file Go sumbernya.
- Union literal untuk status (`'in_progress' | 'submitted' | 'expired'`), bukan `string`.
- Pesan galat diambil lewat satu helper (`errorMessage(e: unknown): string`), bukan `(e as Error).message` berulang.

### 3.4 Gaya

- Format: Prettier (tanpa titik koma, kutip tunggal, `printWidth` 120, trailing comma) sesuai gaya kode sekarang.
- Lint: ESLint flat config dengan `typescript-eslint` dan `eslint-plugin-react-hooks`; `exhaustive-deps` aktif. Setiap `eslint-disable` wajib ada alasan di baris yang sama.
- Dependensi baru dipin versi tepat (lihat `NOTES.md`), dan versinya minimal berumur dua minggu.
- Teks UI berbahasa Indonesia; identifier, komentar, dan pesan commit berbahasa Inggris.
- Aksesibilitas: tombol punya label, ikon dekoratif `aria-hidden`, galat memakai `role="alert"`.

## 4. Proses

- Refactor tidak mengubah perilaku: setiap langkah refactor diakhiri tes yang sama tetap hijau (Go: termasuk tes DB terhadap Postgres sementara; web: `npm run build` dan Vitest).
- Satu commit = satu langkah refactor yang bisa di-review sendiri; jangan campur refactor dengan fitur.
- Pindah file memakai `git mv` agar riwayat tetap terbaca.
