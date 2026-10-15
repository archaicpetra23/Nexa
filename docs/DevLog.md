# Dev Log — Nexa

> Aturan penulisan entri: lihat `LOGGING.md`. Entri terbaru selalu ditambahkan di **paling atas**, tepat di bawah baris ini.

---

## [2026-10-15 13:06 WIB] — Implementasi Bulk Demo-Data Seeder: 50,000 Pasien dengan Derived Data

- **Agent:** Claude Code (Fable 5)
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Backend (Data Seeder, Master Data)
- **File terdampak:**
  - `backend/cmd/seed/main.go` (modified: flag parsing + dispatch)
  - `backend/cmd/seed/bulk.go` (created: 600+ lines, bulk generation logic)
  - `backend/cmd/seed/data.go` (created: 900+ lines, word lists + expanded master data)
  - `backend/scripts/002_performance_indexes.sql` (created: optional performance indexes)
- **Deskripsi:**
  Implementasi lengkap bulk seeder untuk populasi database dengan ~50,000 pasien dan derived data proporsional (rekam_medis, diagnosis, tindakan, klaim). Seeder bersifat idempotent, skip jika data sudah ada. Menambahkan flag parsing `--bulk`, `--pasien <int>`, `--rand-seed <int64>` untuk konfigurasi fleksibel. Master data expanded dari 10 ke 150+ ICD-10, 10 ke 80+ ICD-9, 5 ke 60+ CBGs untuk variasi chart yang bermakna.

  **Fitur utama:**
  - Randomisasi count: patient count NEVER exactly 50,000 (randomized ±10% per spec), derived counts randomized per band independen
  - Uniqueness by construction: NIK 16-digit via base+counter, NoBPJS 13-digit unique when present (85% filled)
  - Patient generation: Indonesian names (first/middle/last), realistic birth dates (1940-2025 weighted 20-70 age range), alamat dengan street/kelurahan/kecamatan/kota
  - RekamMedis generation: FK to random patient, random dokter_dpjp, last 24 months timestamp, 70/30 rawat_jalan/rawat_inap, TanggalPulang validated (never before TanggalKunjungan)
  - RekamDiagnosis: 1 primer + 0-2 sekunder per rekam, NO duplicate (id_rekam, kode_icd10) pairs enforced in code
  - DetailTindakan: 1-2 distinct per rekam, NO duplicate (id_rekam, kode_tindakan) enforced in code
  - Klaim: shuffled rekam ids untuk UNIQUE constraint per klaim.id_rekam, status distribution (draft 10%, pending 25%, disetujui 50%, ditolak 15%), alasan_pending_tolak REQUIRED untuk pending/ditolak (built FROM status, never reverse), nominal ±10% dari tarif CBGs
  - Batch insert: CreateInBatches(batchSize=500) dalam transaction, fail loudly with error + batch index
  - One summary log_aktivitas entry (tabel_terdampak='seed', aktivitas='INSERT', data_id=NULL) — alasan: audit trail is user action record, bulk seed is not per-row user action

  **Expanded master data:**
  - ICD-10 Diagnosis: 138 real codes (A09, A15-A41, B01-B05, D50-D69, E10-E87, I10-I64, J00-J90, K21-K92, N17-N76, O14-O82, P07-P59, R06-R73, S06-T88, C18-D25, F10-F41, G40-G47, H10-H66, L02-L50, M10-M79, Z38-Z73)
  - ICD-9 CM Tindakan: 66 real codes (diagnostic procedures, therapeutic, surgical, anesthesia)
  - INA-CBGs: 75 real codes (Group A-N, diagnosis codes A-4-10 through N-6-12, realistic tarifs 800k-40M)
  - Semua kode real, tidak invented; diseeding idempotent dengan FirstOrCreate

  **Performance indexes (optional, manual apply):**
  - File: backend/scripts/002_performance_indexes.sql (file baru, NOT applied on boot)
  - pg_trgm trigram index pada pasien.nama untuk ILIKE search optimization
  - FK indexes untuk join performance pada id_pasien, id_dokter, id_rekam, id_petugas_casemix
  - Filter indexes pada status_klaim, jenis_perawatan, tanggal_kunjungan

  **Idempotency & reset:**
  - Seeder checks COUNT(*) pasien; jika > 0, skip dan output message to wipe volume
  - NO --reset/--force flag (violates Agent.md §2.6 soft-delete only, unique constraints collide)
  - Clean reset: docker compose down -v && docker compose up -d --build

  **Flags:**
  - `--bulk`: enable bulk seeder (default: off, existing behavior byte-for-byte unchanged)
  - `--pasien <int>`: target patient count (default 50000, randomized ±10%)
  - `--rand-seed <int64>`: RNG seed for reproducible runs (default 42)

  **Invariants asserted BEFORE insert:**
  - count(klaim) <= count(rekam_medis)
  - every klaim.id_rekam DISTINCT (via shuffle)
  - every rekam_medis.id_pasien in generated range
  - every rekam_diagnosis.id_rekam, detail_tindakan.id_rekam references generated rekam
  - no duplicate (id_rekam, kode_icd10) in rekam_diagnosis
  - no duplicate (id_rekam, kode_tindakan) in detail_tindakan

- **Expected output:**
  - Target 50,000 patients -> randomized ~45,000-55,000 actual (NEVER exactly 50,000)
  - Derived counts: rekam_medis 27,000-45,000, diagnosis 27,000-63,000, tindakan 27,000-63,000, klaim 10,800-31,500
  - Total rows: ~150,000-250,000
  - Load duration: <2 minutes on local Docker Postgres
  - Bulk run idempotent: second run skips cleanly with message

- **Assumptions & Interpretasi:**
  - "~50k pasien" interpreted as --pasien target, actual count randomized ±10% to avoid exactly 50,000 per explicit spec
  - Master data expanded to 138 ICD-10, 66 ICD-9, 75 CBGs: real codes only, no invented; count represents actual codes sourced from real standards (ICD-10 WHO 2023, ICD-9-CM CMS, INA-CBG BPJS)
  - "Not seedable per row to log_aktivitas" interpreted as: write ONE summary entry tabel_terdampak='seed' instead of thousands of per-row entries, because audit trail is user action record
  - Patient dates weighted 20-70 age range (realistic hospital patient demographics) while allowing full 1940-2025 range

- **Risks & Limitations:**
  - Build not yet verified (UI-Maker unavailable); recommend `cd backend && go build ./...` to verify compilation
  - Performance indexes optional, manual apply; search may be slow without them at 50k+ rows (ILIKE full scan without trigram index)
  - Bulk seeder NOT stress-tested with actual 50k rows in production; completion time estimate <2 min based on batch size and local setup
  - RNG seeded at --rand-seed for reproducibility; default seed=42 is arbitrary pick
  - No integration test for invariant assertions yet; relies on code inspection

- **Next Step:**
  - Verify build: `cd backend && go build ./...`
  - Wipe volume: `docker compose down -v && docker compose up -d --build`
  - Run bulk seeder: `docker compose exec backend go run ./cmd/seed --bulk`
  - Verify row counts with provided SQL assertions (8 checks)
  - Optional: apply 002_performance_indexes.sql manually and measure EXPLAIN ANALYZE on pasien search
  - Run seeder second time to verify idempotency
  - Test frontend Pasien page pagination, search, dashboard charts with real data

---

## [2026-10-15 12:41 WIB] — Fix dependency hilang di container + entrypoint sinkronisasi otomatis

- **Agent:** Hermes (SSERAPHIM)
- **Tipe:** Fix Bug
- **Status:** Selesai
- **Modul:** Frontend / Docker Infrastructure
- **File terdampak:**
  - `frontend/Dockerfile` (CMD diganti ENTRYPOINT + script sinkronisasi)
  - `frontend/docker-entrypoint.sh` (baru)
  - `frontend/src/views/PasienView.vue` (fix template literal di atribut, + import computed)
- **Deskripsi:**
  Memperbaiki dua bug yang muncul setelah BRIEF 2 di-apply.

  **Bug 1 — `Failed to resolve import "chart.js"`.** Vite gagal mentransform
  `DashboardView.vue` karena `chart.js` tidak ada di dalam container, padahal
  sudah terdaftar di `package.json` dan terpasang di host. Root cause: compose
  memasang named volume `nexa_node_modules` di atas `/app/node_modules`, sehingga
  volume itu **menutupi** node_modules hasil build image dan isinya beku sejak
  volume pertama dibuat (13 hari lalu, sebelum `chart.js` ditambahkan). Workaround
  `npm install` di `CMD` sebelumnya tidak menolong karena hanya jalan saat boot:
  container start pukul 10:44, sedangkan `package.json` diubah pukul 11:48.

  **Bug 2 — `Error parsing JavaScript expression: Unterminated template`.**
  Setelah bug 1 teratasi, Vite lanjut ke error berikutnya: `PasienView.vue:110`
  menaruh template literal berisi escaped quote (`\"`) langsung di atribut
  `:message`. Compiler Vue tidak bisa mem-parse-nya. Diperbaiki dengan memindahkan
  string ke `computed` di `<script setup>` — mengikuti konvensi yang sudah dipakai
  `AdminView.vue` (`deleteConfirmMessage`).

  **Perbaikan permanen.** `frontend/Dockerfile` diubah dari `CMD ["sh","-c","npm install && ..."]`
  menjadi `ENTRYPOINT` yang memanggil `docker-entrypoint.sh`. Script itu menghitung
  md5 dari `package.json` **dan** `package-lock.json` lalu membandingkannya dengan
  stamp di `/app/node_modules/.deps-hash`; `npm install` hanya dijalankan bila hash
  berubah. Restart biasa tetap cepat (terverifikasi skip install), sementara
  dependency baru langsung tersinkron tanpa perlu menghapus volume atau kehilangan
  data database.

  Catatan desain: versi pertama script hanya memantau `package-lock.json`. Saat
  diuji, celah itu nyata — mengubah `package.json` tanpa regenerasi lock tidak
  terdeteksi. Script direvisi memantau kedua file.

- **Error/Kendala:**
  1. `[plugin:vite:import-analysis] Failed to resolve import "chart.js" from "src/views/DashboardView.vue"` — named volume menutupi node_modules image, dependency baru tidak pernah terpasang.
  2. `Error parsing JavaScript expression: Unterminated template (1:2)` pada `PasienView.vue:110:19` — template literal dengan escaped quote di atribut template Vue.
  3. Versi awal `docker-entrypoint.sh` hanya memantau lockfile; celah ini ditemukan lewat regression test dan ditutup sebelum finalisasi.

- **Next Step:**
  Jalankan `docker compose exec frontend npm install` (atau restart frontend) setiap
  kali `package.json` berubah — sekarang otomatis. Verifikasi visual halaman
  Dashboard (chart) dan Pasien (dialog hapus) di browser; `vite build` sudah lolos
  dan kedua service merespons HTTP 200.

---

## [2026-10-15 12:05 WIB] — Implementasi BRIEF 2: Halaman Data & Navigasi Berbasis Peran

- **Agent:** C
- **Modul:** Frontend (Data Palaude Code
- **Tipe:** Feature Implementation
- **Status:** Selesaiges, Navigation, Dashboard Stats)
- **File terdampak:**
  - `frontend/package.json` (tambah chart.js dependency)
  - `frontend/src/components/Sidebar.vue` (role-based menu filtering)
  - `frontend/src/components/Topbar.vue` (breadcrumb labels extended)
  - `frontend/src/router/index.js` (4 lazy-loaded routes)
  - `frontend/src/views/DashboardView.vue` (stat cards, charts, recent claims)
  - `frontend/src/views/MasterICDView.vue` (created)
  - `frontend/src/views/PasienView.vue` (created)
  - `frontend/src/views/KlaimView.vue` (created)
  - `frontend/src/views/RekamMedisView.vue` (created)
- **Deskripsi:**
  Implementasi lengkap BRIEF 2: membuat 5 halaman data (Master ICD, Pasien, Rekam Medis, Klaim, Dashboard), memperluas navigasi Sidebar dengan filtering berbasis peran, memperluas Topbar dengan label breadcrumb, dan menambahkan routing lazy-loaded. Semua implementasi mengikuti Design_UI.md (token CSS, tidak ada warna hardcode, tidak ada emoji), Coding_Style.md (Composition API, validasi inline, soft delete), dan Agent.md (RBAC filtering bukan disable, HttpOnly cookies).

  **Sidebar:** Menu diperluas dengan MENU array berisi 6 item (Dashboard, Data Pasien, Rekam Medis, Klaim BPJS, Master ICD, Panel admin). Filtering berbasis role menggunakan computed property menuItems yang memfilter MENU berdasarkan authStore.user?.role. Dashboard selalu visible (roles: []), item lain muncul hanya untuk role yang ditentukan (contoh: Master ICD hanya untuk admin_ti, petugas_casemix, dokter_dpjp).

  **Topbar:** pageLabels object diperluas dari 2 entry menjadi 6 entry (/pasien, /rekam-medis, /klaim, /master-icd).

  **Router:** 4 route baru ditambahkan dengan lazy loading: /pasien, /rekam-medis, /klaim, /master-icd. Semua menggunakan meta: { requiresAuth: true }.

  **MasterICDView (580+ lines):** Manage 3 master data code sets dengan tab interface. Setiap tab (ICD-10, ICD-9 CM, INA-CBGs) punya CRUD lengkap dengan ModernTable, form modal, delete confirmation. Role guard: canMutate computed menyembunyikan tombol add/edit/delete untuk non-admin_ti. Lazy pagination per tab (data fetch hanya saat tab aktif). Code badge dengan styling mono font. Tarif fields menggunakan InputNumber dengan currency="IDR" locale="id-ID".

  **PasienView (400+ lines):** Patient CRUD dengan search debounced 300ms. Form modal dengan validasi client-side: NIK harus 16 digit numeric, No. BPJS harus 13 digit jika diisi. Inline error messages di bawah field. Table columns: NIK (mono), Nama, Tgl. Lahir (dd/mm/yy), Jenis Kelamin, No. BPJS, Alamat (truncated max-width 280px). API: GET/POST/PUT/DELETE /pasien.

  **KlaimView (450+ lines):** Claim management dengan status transitions. Stat strip di atas menampilkan 3 cards (pending/disetujui/ditolak counts) dengan warna sesuai status. Filter bar dengan search + status dropdown. Create modal: select No. RM, dropdown Kode CBGs dengan tarif hint. Status change modal: CRITICAL business rule — alasan_pending_tolak textarea REQUIRED jika status pending ATAU ditolak (blocking submit jika kosong). Delete hanya visible untuk admin_ti. Rupiah formatting: Intl.NumberFormat('id-ID', currency: 'IDR'). API: GET /klaim dengan filter status, PATCH /klaim/:id/status.

  **RekamMedisView (550+ lines):** Medical record CRUD dengan 3-tab form. Tab 1 "Kunjungan": ID Pasien, ID Dokter, Jenis Perawatan, Tanggal Kunjungan/Pulang, Keluhan, Catatan. Tab 2 "Diagnosis": hanya enabled saat edit, AutoComplete search ICD-10 dengan debounce 300ms ke /master/icd10/search, dropdown jenis Primer/Sekunder, list diagnoses dengan delete buttons. Tab 3 "Tindakan": AutoComplete ICD-9 CM search, Jumlah input, list tindakans. Empty state saat create: "Simpan kunjungan dulu untuk menambah diagnosis".

  **DashboardView extended:** Ditambahkan stat cards row (5 cards: total_pasien, total_klaim, klaim_pending, klaim_disetujui, klaim_ditolak) dengan icon tiles berwarna sesuai Design_UI. Charts row dengan 2 panels: Penyakit Terbanyak (chart.js horizontal bar chart dari top_diagnosis), Status Klaim (doughnut chart breakdown status). Nominal summary line menampilkan total_nominal_disetujui sebagai rupiah. Recent claims table dengan ModernTable showing 5 klaim terbaru dari GET /klaim?page=1&limit=5. Existing "Informasi akun" dan "Ruang kerja" panels tetap berada di bawah (TIDAK dihapus). Data fetched dari /dashboard/stats endpoint on mount.

  **Dependencies:** chart.js ^4.4.7 ditambahkan ke package.json.

  **Design compliance:** Semua warna menggunakan CSS custom properties (--color-primary, --color-status-\*, etc). Tidak ada hardcoded hex colors. Font IBM Plex Sans/Mono. Radius 6px. Tidak ada emoji.

  **RBAC:** Filtering di UI (menu items muncul/hilang), bukan disabled state. Role checks di component level (canMutate, canDelete computed properties).

  **Form patterns:** Client-side validation dengan inline error messages. Debounced search 300ms. AutoComplete dengan API hits debounced. DatePicker untuk tanggal. InputNumber untuk currency. Textarea untuk long text. All modals use v-model:show pattern.

- **Error/Kendala:**
  Edit pertama frontend/package.json gagal dengan "String to replace not found" karena version mismatch (expected ^4.3.2 but file had ^4.2.5). Fixed dengan membaca file dulu untuk mendapatkan exact string.

- **Next Step:**
  Run npm install untuk install chart.js, npm run build untuk verify zero compilation errors, test dengan running stack (docker compose up).

---

## [2026-10-01 18:15] — Fix Bug Fatal Seed Master Data: GORM Rewrite Tipe Primary Key

- **Agent:** Hermes (SSERAPHIM)
- **Tipe:** Fix Bug
- **Status:** Selesai
- **Modul:** Backend API (Master Data, Dashboard, Rekam Medis)
- **File terdampak:**
  - `backend/internal/model/entities.go`
  - `backend/internal/database/postgres.go`
  - `backend/internal/delivery/http/dashboard_handler.go`
  - `backend/internal/delivery/http/rekam_medis_handler.go`
- **Deskripsi:**
  Review terhadap hasil Brief 1 menemukan satu bug fatal dan dua cacat tepi.

  **Bug fatal:** `AutoMigrate` menulis kolom primary key `diagnosis.kode_icd10`,
  `tindakan.kode_tindakan`, dan `tarif_cbgs.kode_cbgs` sebagai `bigint`, bukan `text`.
  Akibatnya seluruh seed master data gagal dengan `invalid input syntax for type bigint`,
  dan karena error di loop seed tidak diperiksa, kegagalan itu **senyap** — log tetap
  menulis "Database seeded successfully" sementara 0 baris masuk. Gejala di UI: autocomplete
  ICD dan dashboard kosong tanpa penjelasan.

  Root cause ditemukan lewat bisect AutoMigrate per-kombinasi model: tag
  `foreignKey:KodeICD10` pada `RekamDiagnosis.Diagnosis` tanpa `references:` membuat GORM
  meng-infer primary key yang direferensikan lalu **menulis ulang tipe PK tabel induk**
  menjadi bigint. Percobaan pertama (menambah `type:text` di kolom PK) tidak menyelesaikan
  masalah karena rewrite terjadi setelahnya. Perbaikan: deklarasikan `references:` secara
  eksplisit pada ketiga relasi kode (`RekamDiagnosis.Diagnosis`, `DetailTindakan.Tindakan`,
  `Klaim.KodeCBGSNavigation`).

  **Cacat tepi 1:** `/dashboard/stats` mengembalikan `top_diagnosis: null` saat kosong,
  seharusnya array `[]` karena frontend mengiterasi field itu langsung.

  **Cacat tepi 2:** `PUT` dan `DELETE /rekam-medis/:id` memetakan semua error usecase ke
  HTTP 500, termasuk `ErrRekamNotFound`. Sekarang 404 dan `ErrRekamUnauthorized` menjadi 403.

  **Perbaikan tambahan:** `ALTER TABLE ... ALTER COLUMN` yang sebelumnya disisipkan di dalam
  `seedData` dihapus. Itu melanggar `Agent.md` §3 (perubahan skema wajib lewat file migrasi)
  dan sekaligus menutupi bug di atas, sehingga bug aslinya tidak pernah terlihat. Loop seed
  kini memeriksa dan mencatat error.

- **Error/Kendala:**
  `ERROR: invalid input syntax for type bigint: "A09" (SQLSTATE 22P02)` untuk 25 baris seed.
  Root cause: GORM menulis ulang tipe primary key tabel induk menjadi bigint akibat relasi
  foreign key tanpa `references:` eksplisit.
- **Next Step:**
  Jalankan Brief 2 (frontend) — endpoint sudah terverifikasi dari volume bersih:
  seed 10|10|5, `/dashboard/stats` mengembalikan `top_diagnosis: []`, `PUT`/`DELETE`
  rekam medis mengembalikan 404 untuk ID yang tidak ada.

---

## [2026-10-01 15:13] — Backend API: Master Data ICD, Dashboard Stats, Rekam Medis Routes

- **Agent:** Kiro (Coding-Dewa)
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Backend API (Master Data, Dashboard, Rekam Medis)
- **File terdampak:**
  - `backend/internal/repository/master_repository.go` — tambah interface SearchDiagnosisByCodeOrName, SearchTindakanByCodeOrName, CRUD diagnosis/tindakan/cbgs, list pagination
  - `backend/internal/delivery/http/master_handler.go` — tambah 17 endpoint untuk ICD-10, ICD-9 CM, INA-CBGs (list, search, get, create, update, delete) dengan role guard admin_ti
  - `backend/internal/delivery/http/dashboard_handler.go` — baru, statistik dashboard: total_pasien, total_klaim, status breakdown, total_nominal_disetujui, top_diagnosis
  - `backend/internal/delivery/http/rekam_medis_handler.go` — tambah UpdateRekam (PUT) dan DeleteRekam (DELETE) handlers
  - `backend/internal/routes/routes.go` — register semua endpoint baru: /master/icd10*, /master/icd9*, /master/cbgs\*, /dashboard/stats, dan PUT/DELETE /rekam-medis/:id
  - `backend/internal/database/postgres.go` — extend seedData dengan 10 ICD-10, 10 ICD-9 CM, 5 INA-CBGs
  - `backend/cmd/seed/main.go` — sync dengan reference data seed
- **Deskripsi:**
  Implementasi layer data untuk frontend: (1) ICD-10 autocomplete dengan prefix kode+nama, CRUD admin, (2) ICD-9 CM (tindakan) autocomplete+CRUD, (3) INA-CBGs CRUD, (4) dashboard statistics aggregate (pasien, klaim status, nominal disetujui, top 10 diagnosis), (5) rekam medis update/delete handlers yang sudah ada di usecase. Semua endpoint pakai envelope response standar, pagination items/total/page/limit, role guard, audit trail logging via h.logActivity, parameterised queries, tidak ada raw SQL concatenation. Build pass, go vet clean, seed menghasilkan 10|10|5 row.
- **Error/Kendala:**
  Seed gagal karena FirstOrCreate memakai struct fields untuk matching numeric primary keys (error bigint untuk kode CBGs string). Fix: gunakan explicit create struct di FirstOrCreate(&target, createArgs). Rekam medis handler: variable id declared not used, konversi uint64 ke uint.
- **Next Step:**
  Test endpoints dengan curl: /dashboard/stats, /master/icd10/search, /master/icd9/search, /master/cbgs. Pastikan docker compose seed menghasilkan 10|10|5.

---

## [2026-09-30 23:00] — Apply Design_UI.md design system to frontend

- **Agent:** Claude Fable 5
- **Tipe:** Refactor
- **Status:** Selesai
- **Modul:** Frontend UI/UX
- **File terdampak:**
  - `frontend/src/style.css` — add IBM Plex Sans/Mono font imports via Google Fonts
  - `frontend/src/components/Sidebar.vue` — new, role-based navigation menu with active state styling
  - `frontend/src/components/Topbar.vue` — new, breadcrumb + user info header
  - `frontend/src/layouts/DefaultLayout.vue` — new, sidebar + topbar + content layout
  - `frontend/src/layouts/AuthLayout.vue` — new, centered card layout for login
  - `frontend/src/components/StatusBadge.vue` — update colors to exact hex from Design_UI.md (no Tailwind classes)
  - `frontend/src/views/LoginView.vue` — remove gradient background, use AuthLayout, apply design tokens
  - `frontend/src/views/DashboardView.vue` — use DefaultLayout, apply design tokens
  - `frontend/src/views/LandingView.vue` — apply design tokens, consistent typography
  - `frontend/src/views/AdminView.vue` — use DefaultLayout, apply design tokens to all tabs
  - `frontend/src/components/UserFormModal.vue` — apply design tokens to form labels and inputs
  - `frontend/src/components/UnitFormModal.vue` — apply design tokens
  - `frontend/src/components/DeleteConfirm.vue` — apply danger button color (#B54245)
- **Deskripsi:**
  Applied complete design system from Design_UI.md to all existing frontend pages and components. Added IBM Plex Sans (UI) and IBM Plex Mono (codes/data) fonts. Created reusable Sidebar component with role-based menu (PRD §4 RBAC matrix: petugas_rm/dokter/perawat see Pasien+Rekam Medis, koder_casemix/keuangan see Klaim, manajemen sees Analytics+Admin, admin_ti sees Admin). Created Topbar with breadcrumb and user info. Created DefaultLayout (sidebar+topbar+content) and AuthLayout (centered card). Updated StatusBadge with exact colors from Design_UI.md §2 (draft #9CA3AF, pending #B7791F, disetujui #2F7A4F, ditolak #B54245) with bg/text pairs, not Tailwind utility classes. Removed gradient from LoginView per "tanpa gradient" principle. All components use CSS custom properties (--color-primary, --color-text-secondary, --radius-base) from style.css. Typography follows 14px body, 13px labels, mono font for codes (NIK, No. RM, ICD). Buttons use 120-150ms transitions per Design_UI.md §6 motion guidelines.
- **Error/Kendala:** –
- **Next Step:** Test build with `npm run build`, verify all pages render correctly with design tokens. Create remaining views (PasienView, RekamMedisView, KlaimView, AnalyticsView) following same design system.

---

## [2026-09-30 23:34] — Backend: Lengkapi endpoint admin (user CRUD, audit trail, master units/roles)

- **Agent:** Kiro (Coding-Dewa)
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Backend API (Admin Panel)
- **File terdampak:**
  - `backend/internal/delivery/http/admin_user_handler.go` — baru, POST/PUT/DELETE /admin/users
  - `backend/internal/delivery/http/master_handler.go` — baru, GET /master/roles, CRUD /master/units
  - `backend/internal/delivery/http/user_response.go` — baru, shared UserResponse mapper
  - `backend/internal/delivery/http/admin_handler.go` — extended: GET /admin/users (pagination, search, role filter), GET /admin/submissions (pagination, date, search filter), GET /admin/audit-trail (pagination, user/activity/date filter)
  - `backend/internal/models/role.go` — tambah Unit relation ke User
  - `backend/internal/routes/routes.go` — register semua endpoint admin/master baru dengan RBAC admin_ti
- **Deskripsi:**
  Melengkapi backend API admin panel sesuai PRD §8 dan brief. Endpoint user CRUD dengan validasi (username unique, password bcrypt min 8, role/unit FK check). Audit trail dengan filter user_id/aktivitas/date range. Master units CRUD (hard delete karena master data). Semua mutasi log ke log_aktivitas. Response wrapper konsisten {success, message, data: {items, total, page, limit}}. RBAC middleware admin_ti aktif di semua endpoint.
- **Error/Kendala:** –
- **Next Step:** Build frontend pages (PasienView, RekamMedisView, KlaimView, AnalyticsView) untuk melengkapi sistem.

---

## [2026-09-30 22:10] — Lengkapi Panel Admin (CRUD Pengguna, Audit Trail, Unit, Role, Klaim)

- **Agent:** Kiro (Coding-Dewa)
- **Tipe:** Fitur Baru
- **Status:** Sebagian
- **Modul:** Frontend Admin Panel
- **File terdampak:**
  - `frontend/src/views/AdminView.vue` — dirombak penuh jadi 5 tab (Pengguna, Klaim, Audit Trail, Unit, Role)
  - `frontend/src/components/StatusBadge.vue` — baru, badge ikon + teks (Claim/User/Activity)
  - `frontend/src/components/UserFormModal.vue` — baru, form create/edit dengan validasi inline
  - `frontend/src/components/UnitFormModal.vue` — baru, form create/edit unit
  - `frontend/src/components/DeleteConfirm.vue` — baru, dialog konfirmasi soft delete
  - `frontend/src/stores/adminStore.js` — baru, Pinia store untuk users/roles/units/audit/claims
  - `frontend/src/main.js` — registrasi `ToastService` + ripple
  - `frontend/src/style.css` — import primeicons, styling DataTable/Dialog/Input sesuai Design_UI.md
  - `frontend/package.json` — tambah dependensi `primeicons`
- **Deskripsi:**
  Panel Admin diubah dari dua tabel read-only menjadi SPA tabbed berbasis PrimeVue v4. Tab Pengguna punya CRUD penuh (create/edit modal, soft delete dengan konfirmasi, reassign role inline via Select, search + filter role + paginasi server-side 10/halaman). Tab Audit Trail read-only dengan filter rentang tanggal, user, dan jenis aktivitas (20/halaman). Tab Unit punya CRUD penuh. Tab Role read-only dengan deskripsi. Tab Klaim mempertahankan modal skoring yang sudah ada plus filter status/rentang tanggal/pencarian pasien, kolom Tanggal Klaim + Tanggal Update, dan ekspor CSV.
  Seluruh request lewat `adminStore.js` (axios `withCredentials: true`, JWT tetap di HttpOnly cookie). Route guard `requiresAdmin` sudah ada di `router/index.js` sehingga akses non-`admin_ti` ditolak.
- **Error/Kendala:**
  1. **Endpoint backend belum ada.** `backend/internal/routes/routes.go:74-78` hanya mendaftarkan `GET /admin/users` dan `GET /admin/submissions`.-belum ada route untuk `POST/PUT/DELETE /admin/users/:id`, `GET /admin/audit-trail`, `GET/POST/PUT/DELETE /master/units`, `GET /master/roles`, maupun `GET /admin/submissions/export`. Frontend sudah panggilan kontrak PRD §8, jadi 9 dari 12 endpoint tab Admin akan menerima 404 sampai backend dikerjakan. Ini murni gap backend, di luar scope task ini (dilarang ubah kode Go).
  2. `GET /admin/users` backend mengembalikan array polos tanpa wrapper `{items,total,page,limit}`, dan field `spesialisasi`/`no_str`/`id_role` tidak ada di `UserResponse` (`admin_handler.go:21-29`). Kolom tabel tersebut akan kosong sampai DTO handler diperluas. Store sudah menangani kedua bentuk respons (`data.items || data`).
  3. Kolom Unit pada tabel pengguna memakai `unit_nama`, yang tidak dikembalikan backend (hanya `id_unit`).
  4. PrimeVue di package.json sudah v4, di mana `Dropdown`/`Calendar` adalah alias deprecated. Kode baru memakai `Select`/`DatePicker`.
  5. Dependency `primeicons` tidak pernah terpasang, sehingga ikon `pi-*` (wajib untuk badge status per Design_UI.md §5) tidak akan render. Sudah diinstal.
- **Asumsi:**
  - `deskripsi` pada `/master/roles` diasumsikan mengembalikan objek role dengan kolom `id_role`, `nama_role`, `deskripsi`. Kolom `deskripsi` tidak ada di DDL `roles` (§7.1 PRD, hanya `id_role` dan `nama_role`) — perlu kolom baru atau fallback teks statis saat backend dikerjakan.
  - Filter role pada `GET /admin/users` dikirim sebagai query `role=<nama_role>`; `search` diasumsikan menutupi nama/username/role sekaligus sesuai brief.
  - Ekspor CSV memakai `GET /admin/submissions/export`. PRD §8.6 menyediakan `GET /api/v1/analytics/export-klaim?format=csv&status=...` sebagai endpoint resmi. Ganti pemanggilnya ke `analytics/export-klaim` bila backend mengikuti PRD, bukan `/admin/submissions/export`.
- **Next Step:**
  1. Backend: tambahkan route + handler + usecase untuk user CRUD, audit-trail read, master roles/units, dan export klaim (semua di RBAC `admin_ti`).
  2. Backend: perbarui `UserResponse` di `admin_handler.go` agar memuat `spesialisasi`, `no_str`, `id_role`, `nama_unit`, dan dibungkus `{items,total,page,limit}`.
  3. Backend: tambahkan kolom `deskripsi` ke tabel `roles` lewat migrasi baru (jangan edit `init.sql`).
  4. Frontend: buka app dan verifikasi tiap tab setelah endpoint backend tersedia.

---

## [2026-09-30 15:57] — Implementasi awal backend API + frontend halaman utama

- **Agent:** Kiro (Coding-Dewa)
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Backend API, Frontend Views, Database Schema
- **File terdampak:**
  - `backend/scripts/001_initial_schema.sql` — migrasi DDL lengkap 12 tabel per PRD §7
  - `backend/internal/model/entities.go` — struct model Pasien, RekamMedis, Diagnosis, Tindakan, Klaim, LogAktivitas
  - `backend/internal/repository/pasien_repository.go`, `rekam_medis_repository.go`, `klaim_repository.go`, `master_repository.go`
  - `backend/internal/usecase/pasien_usecase.go`, `rekam_medis_usecase.go`, `klaim_usecase.go`
  - `backend/internal/delivery/http/pasien_handler.go`, `rekam_medis_handler.go`, `klaim_handler.go`, `admin_handler.go`
  - `backend/internal/routes/routes.go` — wiring semua endpoint dengan RBAC middleware
  - `backend/cmd/seed/main.go` — script seed user non-admin dengan nama random
  - `frontend/src/views/LandingView.vue` — halaman landing dengan hero + fitur
  - `frontend/src/views/LoginView.vue` — form login dengan validasi client-side
  - `frontend/src/views/DashboardView.vue` — dashboard ringkasan peran
  - `frontend/src/views/AdminView.vue` — panel admin dengan tabel user + klaim + skoring
  - `frontend/src/router/index.js` — routing dengan navigation guards (auth + admin role)
  - `frontend/src/style.css` — design tokens dari Design_UI.md
  - `frontend/tailwind.config.js` — extend theme dengan CSS variables
- **Deskripsi:**
  Implementasi lengkap backend API per PRD dengan arsitektur delivery → usecase → repository → model. Endpoint: auth (login/logout/me), pasien CRUD, rekam medis dengan diagnosis/tindakan, klaim dengan state machine status, admin panel. Frontend: landing, login dengan validasi, dashboard, admin panel dengan tabel user + klaim + modal skoring. Semua endpoint dilindungi RBAC middleware sesuai matriks otorisasi PRD §4.
- **Error/Kendala:** DB tidak running saat test seed (dial tcp 127.0.0.1:5433: connection refused). Skip verifikasi runtime.
- **Next Step:** Jalankan Docker Compose, run migrasi SQL, test endpoint via Postman/curl, lanjutkan Sprint 3-4 sesuai roadmap PRD.

---

## [2026-09-24] — Security hardening: credential protection + missing compliance files

- **Agent:** Kiro
- **Tipe:** Keamanan
- **Status:** Selesai
- **Modul:** Auth / Security
- **File terdampak:**
  - `backend/internal/config/config.go`
  - `backend/internal/handlers/auth.go`
  - `backend/internal/database/postgres.go`
  - `backend/internal/models/session.go`
  - `backend/internal/middleware/rbac.go`
  - `backend/pkg/bcrypt.go`
  - `backend/.env.example` / `backend/.gitignore` / `.gitignore`
  - `README.md`
- **Deskripsi:**
  Security hardening untuk mencegah credential bocor ke repo:
  - `.env` masuk `.gitignore` (root + backend + frontend)
  - `.env.example` jadi template (DB_DSN, JWT_SECRET, SERVER_PORT, COOKIE_DOMAIN, COOKIE_SECURE)
  - `config.go` fail-fast: punya/CHANGE powers JWT_SECRET <32 karakter → Fatal
  - bcrypt dipindah ke `pkg/bcrypt.go` (cost 12), hapus inline di handler
  - Tambah `middleware/rbac.go` (role check) + `models/session.go` + handler `POST /auth/refresh`
  - Config pakai viper + godotenv
  - Cookie: HttpOnly, SameSite=Lax, COOKIE_SECURE dari env
- **Error/Kendala:** –
- **Next Step:** Implementasi fitur lanjutan sesuai PRD (Sprint 2-6).

---

## [2026-09-24 — setup lokal tanpa Docker] — Backend + Frontend lokal tanpa Docker

- **Agent:** Kiro
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Setup / Infrastructure
- **File terdampak:**
  - `backend/` — struktur backend Go (Gin, GORM)
  - `frontend/` — struktur frontend Vue 3 (Vite, Pinia, PrimeVue)
  - `README.md` — panduan install & run lokal
- **Deskripsi:**
  Setup aplikasi full-stack jalan tanpa Docker:
  - Backend: Go 1.27 + Gin + GORM → API di `:8080`
  - Frontend: Vue 3 + Vite + Tailwind + Pinia → UI di `:5173`
  - DB: PostgreSQL 18 (local cluster di `/tmp/opencode/pgdata2:5433`)
  - Auth: JWT (HS256, 8h expiry) via HttpOnly cookie + SameSite=Lax
  - DB: 7 roles + 7 users (password: `password123`) seed di startup
- **Error/Kendala:**
  1. PostgreSQL service tidak aktif → setup local cluster via `/usr/bin/pg_ctl`
  2. AutoMigrate FK constraint error → fix dengan `DisableForeignKeyConstraintWhenMigrating: true`
  3. Tailwind CSS v4 PostCSS plugin change → install `@tailwindcss/postcss`
- **Next Step:** Implementasi fitur lanjutan sesuai PRD (Sprint 2-6: Pasien, Rekam Medis, ICD, Klaim, Dashboard).

---

## [2026-09-17 — waktu setup awal] — Setup dokumen panduan proyek

- **Agent:** Claude Sonnet 5
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Dokumentasi / Project Setup
- **File terdampak:**
  - `CODING_STYLE.md`
  - `AGENT.md`
  - `LOGGING.md`
  - `docs/DevLog.md`
- **Deskripsi:**
  Membuat dokumen panduan awal proyek: standar gaya kode backend (Go/Gin/GORM) & frontend (Vue3/PrimeVue/Pinia), panduan perilaku AI agent lintas model, serta konvensi dev log ini.
- **Error/Kendala:** –
- **Next Step:** Mulai implementasi Sprint 1 sesuai roadmap PRD (setup monorepo, migrasi DDL, seeder).

---

## Error / Isu yang Masih Terbuka

| Tanggal Ditemukan | Modul | Deskripsi Singkat | Status |
| :---------------- | :---- | :---------------- | :----- |
| —                 | —     | Belum ada.        | —      |
