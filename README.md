# Nexa — Sistem Informasi Casemix Terintegrasi

**Nexa** adalah aplikasi web rekam medis dan penagihan terintegrasi yang dirancakan untuk mempercepat siklus operasional klaim BPJS Kesehatan di Rumah Sakit Umum Daerah (RSUD). Sistem ini menjembatani _data silos_ antara loket pendaftaran, dokter DPJP, perawat bangsal, koder Casemix, dan bagian keuangan.

## Tech Stack

- **Backend:** Go 1.25 (Gin Gonic, GORM), PostgreSQL 16
- **Frontend:** Vue 3 (Composition API, Vite, PrimeVue, Tailwind CSS, Pinia)
- **Arsitektur:** Decoupled Monorepo, RESTful API, JWT + RBAC

---

## Prerequisites

Sebelum mulai, pastikan hal berikut sudah terpenuhi:

| Kebutuhan             | Keterangan                                                                   |
| :-------------------- | :--------------------------------------------------------------------------- |
| **Docker Desktop**    | Windows / macOS / Linux. Untuk Windows wajib mengaktifkan backend **WSL 2**. |
| **Docker Compose v2** | Sudah termasuk di Docker Desktop. Cek dengan `docker compose version`.       |
| **Port bebas**        | Port `5432`, `8080`, dan `5173` tidak dipakai aplikasi lain.                 |
| **Koneksi internet**  | Diperlukan saat build pertama (download image + dependency).                 |

> **Catatan untuk pengguna Windows:** seluruh perintah di panduan ini dijalankan di **PowerShell** atau **Git Bash**. Kalau memakai CMD, perintah `cp` diganti `copy`.

---

## Quick Start (Docker) — Cara Termudah & Disarankan

Cara ini **tidak perlu** install Go, Node.js, atau PostgreSQL secara manual. Semua sudah dibungkus di dalam Docker.

### 1. Masuk ke Folder Project

```bash
cd "Nexa"
```

### 2. Buat File `.env`

File `.env` menyimpan konfigurasi database dan keamanan. Template-nya sudah tersedia di `.env.example`.

```bash
cp .env.example .env
```

Lalu buka `.env` dan isi nilainya:

```dotenv
# PostgreSQL
POSTGRES_USER=nexa
POSTGRES_PASSWORD=nexa123
POSTGRES_DB=nexa

# Backend
JWT_SECRET=GANTI_DENGAN_HASIL_OPENSSL
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false
SERVER_PORT=:8080
```

**Generate `JWT_SECRET`** (wajib, minimal 32 karakter):

```bash
# Linux / macOS / Git Bash
openssl rand -base64 32
```

```powershell
# Windows PowerShell
$b = New-Object byte[] 32
[System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($b)
[Convert]::ToBase64String($b)
```

Copy hasilnya ke baris `JWT_SECRET=` pada file `.env`.

> **Penting:** jangan pakai `echo "JWT_SECRET=..." > .env`. Tanda `>` akan **menimpa seluruh isi** `.env` dan menghapus `POSTGRES_*`. Kalau mau menambah lewat terminal, pakai `>>` (append) atau edit manual dengan text editor.

### 3. Build & Run

```bash
# Build semua image + start containers di background
docker compose up --build -d
```

> **Build pertama memakan waktu 3–10 menit** (download image Go & Node, `go mod download`, `npm install`). Ini normal. Build berikutnya jauh lebih cepat karena cache.

### 4. Verifikasi Semua Service Jalan

Tunggu ±30 detik setelah build selesai, lalu cek:

```bash
docker compose ps
```

Semua service harus berstatus **`Up`** (dan `postgres` harus **`Up (healthy)`**):

```
NAME              SERVICE    STATUS
nexa-postgres-1   postgres   Up (healthy)   0.0.0.0:5432->5432/tcp
nexa-backend-1    backend    Up             0.0.0.0:8080->8080/tcp
nexa-frontend-1   frontend   Up             0.0.0.0:5173->5173/tcp
```

Cek health endpoint backend:

```bash
curl http://localhost:8080/health
```

Harus membalas:

```json
{ "message": "Nexa API running", "success": true }
```

> Kalau ada service yang **`Exit`** atau **`Restarting`**, cek log-nya: `docker compose logs backend` (atau ganti `backend` dengan nama service lain).

### 5. Akses Aplikasi & Login

Buka browser:

```
http://localhost:5173
```

Saat pertama kali dijalankan, backend **otomatis membuat tabel database (auto-migrate)** dan **otomatis mengisi akun demo (auto-seed)**. Jadi kamu bisa langsung login tanpa setup database manual.

**Akun demo yang tersedia** — password semuanya `password123`:

| Username    | Role              | Password      |
| :---------- | :---------------- | :------------ |
| `admin`     | `admin_ti`        | `password123` |
| `petugas`   | `petugas_rm`      | `password123` |
| `dokter`    | `dokter_dpjp`     | `password123` |
| `perawat`   | `perawat`         | `password123` |
| `casemix`   | `petugas_casemix` | `password123` |
| `keuangan`  | `keuangan`        | `password123` |
| `manajemen` | `manajemen`       | `password123` |

> Gunakan akun `admin` untuk akses penuh, dan akun lain untuk mencoba pembatasan hak akses (RBAC) per role.

### 6. Stop & Cleanup

```bash
# Stop semua service (data database TETAP tersimpan)
docker compose down

# Stop + hapus volume (⚠️ data DB dan akun akan hilang, akan di-seed ulang saat start)
docker compose down -v
```

---

## Troubleshooting

| Masalah                             | Penyebab                                         | Solusi                                                                                        |
| :---------------------------------- | :----------------------------------------------- | :-------------------------------------------------------------------------------------------- |
| `port is already allocated`         | Port 5432/8080/5173 dipakai proses lain          | Matikan aplikasi tersebut, atau ubah port kiri di `docker-compose.yml` (contoh `"5433:5432"`) |
| `FATAL: DB_DSN not set in .env`     | File `.env` belum dibuat / kosong                | Ulangi langkah 2                                                                              |
| `FATAL: JWT_SECRET not set in .env` | `JWT_SECRET` masih kosong atau masih placeholder | Isi dengan hasil `openssl rand -base64 32`                                                    |
| Backend `Restarting` terus          | Kredensial Postgres tidak cocok                  | Cek `POSTGRES_*` di `.env`, lalu `docker compose down -v && docker compose up --build -d`     |
| Halaman blank / `Failed to fetch`   | Frontend belum selesai `npm install`             | Tunggu 1–2 menit, cek `docker compose logs -f frontend`                                       |
| Perubahan kode tidak muncul         | Vite cache                                       | `docker compose restart frontend`                                                             |
| Tidak bisa login                    | Database belum ter-seed                          | `docker compose down -v && docker compose up --build -d`                                      |

---

## Struktur Project

```
Nexa/
├── backend/              # API Go (Gin + GORM)
│   ├── cmd/server/       # Entry point server
│   ├── cmd/seed/         # Seeder opsional
│   ├── internal/         # config, database, model, routes, handler
│   └── scripts/          # SQL schema referensi
├── frontend/             # SPA Vue 3 + Vite
│   └── src/              # views, components, stores, layouts
├── docs/                 # PRD, Coding Style, DevLog, Agent.md
├── docker-compose.yml    # Orkestrasi postgres + backend + frontend
├── .env.example          # Template konfigurasi
└── README.md
```

---

## Quick Start (Local) — Tanpa Docker

Untuk development yang butuh hot-reload lebih cepat atau debugging langsung. **Lewati bagian ini kalau sudah berhasil pakai Docker.**

### 1. Install Dependencies

**Windows** — install manual:

- **Go 1.25+**: [https://go.dev/dl/](https://go.dev/dl/)
- **PostgreSQL 16**: [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)
- **Node.js 18+**: [https://nodejs.org/](https://nodejs.org/)

Atau pakai Chocolatey:

```powershell
choco install go postgresql nodejs
```

**Ubuntu / Debian:**

```bash
sudo apt update
sudo apt install -y golang postgresql nodejs npm
```

Verify:

```bash
go version          # minimal 1.25
psql --version      # minimal 16
node --version      # minimal 18
```

### 2. Start PostgreSQL Service

**Windows:**

- PostgreSQL otomatis jalan sebagai service setelah install
- Atau buka **pgAdmin** → connect ke postgres

**Ubuntu / Debian:**

```bash
sudo systemctl enable --now postgresql
sudo -u postgres createuser --superuser $USER
sudo -u postgres createdb nexa
```

Atau via DBeaver: connect ke postgres → Create New Database → nama: `nexa`

### 3. Setup Backend

```bash
cd backend
cp .env.example .env
```

Generate JWT secret:

```bash
openssl rand -base64 32
```

Edit `backend/.env`:

```dotenv
DB_DSN="host=localhost user=postgres password=YOUR_PASSWORD dbname=nexa port=5432 sslmode=disable"
JWT_SECRET="OUTPUT_DARI_OPENSSL"
SERVER_PORT=":8080"
COOKIE_DOMAIN="localhost"
COOKIE_SECURE="false"
```

Jalankan:

```bash
go mod tidy
go run cmd/server/main.go
```

→ Backend jalan di `http://localhost:8080`. Tabel dan akun demo otomatis dibuat saat start pertama.

### 4. Setup Frontend

```bash
cd frontend
npm install
npm run dev
```

→ Frontend jalan di `http://localhost:5173`

> Saat mode local, Vite mem-proxy `/api` ke `http://backend:8080` (nama host Docker). Untuk menjalankan tanpa Docker, ubah target proxy di `frontend/vite.config.js` menjadi `http://localhost:8080`.

### 5. Stop Services

```bash
pkill -f "go run cmd/server"
pkill -f "vite"
```

---

## Security Notes

- `.env` TIDAK di-commit (ada di `.gitignore`)
- `.env.example` hanya template dengan placeholder
- JWT secret harus random min 32 byte: `openssl rand -base64 32`
- Cookie: `HttpOnly`, `SameSite=Lax`, `Secure=false` (local only)
- bcrypt cost: 12
- **Akun demo di atas hanya untuk keperluan development/demo.** Untuk production, ganti seluruh password dan hapus seeder akun demo.

---

## Dokumen Lengkap

### Dokumen Produk & Ide

| File                                               | Deskripsi                                                                   |
| :------------------------------------------------- | :-------------------------------------------------------------------------- |
| [`docs/ide_proyek.md`](docs/ide_proyek.md)         | Ide proyek, masalah, target pengguna, fitur inti, dan kriteria keberhasilan |
| [`docs/Nexa (PRD).md`](<docs/Nexa%20(PRD).md>)     | Product Requirement Document (PRD) lengkap dalam format Markdown            |
| [`docs/Nexa (PRD).pdf`](<docs/Nexa%20(PRD).pdf>)   | PRD dalam format PDF                                                        |
| [`docs/Nexa (PRD).docx`](<docs/Nexa%20(PRD).docx>) | PRD dalam format Word                                                       |

### Dokumen Panduan Pengembangan

| File                                           | Deskripsi                                                                           |
| :--------------------------------------------- | :---------------------------------------------------------------------------------- |
| [`docs/Agent.md`](docs/Agent.md)               | Panduan perilaku AI coding agent (main rules, alur kerja per task, larangan khusus) |
| [`docs/Coding_Style.md`](docs/Coding_Style.md) | Standar gaya penulisan kode Backend (Go/Gin/GORM) & Frontend (Vue 3/PrimeVue/Pinia) |
| [`docs/Design_UI.md`](docs/Design_UI.md)       | Panduan desain antarmuka dan komponen UI                                            |
| [`docs/Logging.md`](docs/Logging.md)           | Konvensi dev log — aturan pencatatan pekerjaan ke `DEVLOG.md`                       |
| [`docs/DevLog.md`](docs/DevLog.md)             | Dev log kronologis — riwayat fitur, fix, refactor, dan error yang masih terbuka     |

## Masalah yang Diselesaikan

1. **Pencatatan Terfragmentasi** — alur data manual antar-unit memicu inkonsistensi berkas rekam medis.
2. **Duplikasi & Human Error** — kesalahan input kode diagnosis (ICD-10) dan tindakan (ICD-9 CM) tanpa validasi otomatis.
3. **Pencarian Berkas Lambat** — temu-balik rekam medis fisik memakan waktu 15–60 menit.
4. **Klaim Pending/Dispute** — berkas tertahan tanpa pencatatan alasan penolakan terpusat.

## Target Pengguna (7 Role RBAC)

| Peran             | Tanggung Jawab Utama                         |
| :---------------- | :------------------------------------------- |
| `admin_ti`        | Manajemen pengguna, konfigurasi, audit trail |
| `petugas_rm`      | Registrasi pasien, verifikasi NIK/No BPJS    |
| `dokter_dpjp`     | Input resume klinis, diagnosis ICD-10        |
| `perawat`         | Input tindakan medis ICD-9 CM                |
| `petugas_casemix` | Verifikasi koding, grouping INA-CBGs         |
| `keuangan`        | Monitoring klaim disetujui, rekonsiliasi     |
| `manajemen`       | Akses read-only dasbor analitik              |

## Fitur Inti

1. **Autentikasi & RBAC Multi-Role** — sesi aman JWT via _HttpOnly Cookie_
2. **Manajemen Pasien** — CRUD + validasi NIK (16 digit), No BPJS (13 digit)
3. **Rekam Medis Kunjungan** — rawat jalan/inap + kalkulasi _Length of Stay_
4. **Pencarian Cepat ICD** — autocomplete _fuzzy search_ < 300ms (indeks trigram)
5. **Engine Rekomendasi INA-CBGs** — pemetaan paket tarif otomatis
6. **Siklus Klaim (State Machine)** — `draft` → `pending` → `disetujui`/`ditolak`
7. **Dasbor Analitik** — Top 10 ICD-10 + metrik status klaim _real-time_
8. **Audit Trail** — pencatatan otomatis seluruh mutasi data klinis

## Keamanan

- OWASP Top 10 compliant
- Autentikasi JWT via _HttpOnly_, _Secure_, _SameSite=Lax_ cookie
- RBAC middleware + mitigasi IDOR/BOLA
- _Rate limiter_ login (5 percobaan/menit/IP)
- _Soft delete_ + _ACID transaction_ PostgreSQL

## Target Keberhasilan

- Pencarian berkas: **< 5 detik** (optimal < 1 detik)
- Duplikasi kode medis: **0%**
- Efisiensi koding Casemix: **terpangkas >= 50%**
- Transparansi dispute: **100%** berkas pending/tolak tercatat alasan
- Audit trail: **100%** aksi mutasi tercatat di `log_aktivitas`
