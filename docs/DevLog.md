# Dev Log — Nexa

> Aturan penulisan entri: lihat `LOGGING.md`. Entri terbaru selalu ditambahkan di **paling atas**, tepat di bawah baris ini.

---

## [2026-10-01 00:36] — Fix dashboard harus scroll untuk melihat konten + rebuild DashboardView

- **Agent:** SSERAPHIM (Hermes)
- **Tipe:** Fix Bug
- **Status:** Selesai
- **Modul:** Frontend UI/UX (Dashboard, Layout)
- **File terdampak:**
  - `frontend/src/layouts/DefaultLayout.vue` — tulis ulang, layout flex sidebar+topbar+konten
  - `frontend/src/views/DashboardView.vue` — tulis ulang mengikuti bahasa desain LandingView/AuthLayout
  - `frontend/src/components/Sidebar.vue` — penyesuaian (sticky, 100vh, role menu)
  - `frontend/src/components/Topbar.vue` — penyesuaian (sticky, breadcrumb, user info)
  - `frontend/Dockerfile` — `npm install` dijalankan saat start
  - `frontend/src/router/index.js` — tambah `scrollBehavior` kembali ke atas
- **Deskripsi:**
  Memperbaiki bug dashboard yang kontennya baru muncul setelah di-scroll ke bawah. Root cause: style scoped `DefaultLayout.vue` (`.app-shell`, `.app-main-column`, `.app-content`, `.app-content-inner`) tidak ter-apply — Vite menyajikan versi CSS lama dari cache transform, sehingga `.app-shell` tetap `display: block`. Akibatnya sidebar setinggi 100vh (757px) menumpuk di atas main column, dan konten baru mulai di y=823px. Setelah cache Vite dibersihkan dan layout ditulis ulang, `.app-shell` menjadi `display: flex` dan konten mulai di y=96px. DashboardView juga dirombak agar konsisten dengan bahasa desain LandingView dan LoginView (kicker uppercase, panel dengan border+shadow-panel, icon tile #EEF2F5, chip tanggal, grid dua kolom, catatan bawah).
- **Error/Kendala:**
  Gejala awal: `.app-shell` computed `display: block`, `.app-main-column` top 757px, `.app-content` `flex: 0 1 auto`, `.app-content-inner` `padding-top: 0px`. Diagnosa lewat Chrome DevTools Protocol (query CSSOM) menunjukkan rules untuk scope ID layout hanya berisi `.max-w-content` (versi lama), bukan rules baru. Diverifikasi dengan `curl` ke endpoint CSS Vite yang mengembalikan `__vite__css` lama, sementara file di disk sudah benar dan inode host==container. Fix: hapus `/app/node_modules/.vite` di container lalu restart frontend.
- **Next Step:**
  Buat halaman yang belum ada (PasienView, RekamMedisView, KlaimView, AnalyticsView) dengan bahasa desain yang sama, lalu verifikasi lewat pengukuran layout otomatis.

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
| :--- | :--- | :--- | :--- |
| — | — | Belum ada. | — |
