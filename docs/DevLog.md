# Dev Log — Nexa

> Aturan penulisan entri: lihat `LOGGING.md`. Entri terbaru selalu ditambahkan di **paling atas**, tepat di bawah baris ini.

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
