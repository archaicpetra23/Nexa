# Nexa — Sistem Informasi Casemix Terintegrasi

**Nexa** adalah aplikasi web rekam medis dan penagihan terintegrasi yang dirancakan untuk mempercepat siklus operasional klaim BPJS Kesehatan di Rumah Sakit Umum Daerah (RSUD). Sistem ini menjembatani _data silos_ antara loket pendaftaran, dokter DPJP, perawat bangsal, koder Casemix, dan bagian keuangan.

## Tech Stack

- **Backend:** Go 1.22+ (Gin Gonic, GORM), PostgreSQL 16
- **Frontend:** Vue 3 (Composition API, Vite, PrimeVue, Tailwind CSS, Pinia)
- **Arsitektur:** Decoupled Monorepo, RESTful API, JWT + RBAC

---

## Quick Start (Docker)

### 1. Install Docker

Download & install Docker Desktop sesuai OS masing-masing:

- **[https://www.docker.com/get-started/](https://www.docker.com/get-started/)** (Windows, macOS, Linux)

Pastikan Docker running sebelum lanjut ke langkah berikutnya.

### 2. Build & Run

```bash
# Build semua image + start containers
docker compose up --build -d
```

Service yang jalan:
| Service | Port | Fungsi |
|---------|------|--------|
| PostgreSQL | 5432 | Database |
| Backend | 8080 | API |
| Frontend | 5173 | UI |

### 3. Generate JWT Secret (Opsional)

Generate secret random dan set di environment:
```bash
# Buat file .env di root project
echo "JWT_SECRET=$(openssl rand -base64 32)" > .env

# Restart backend pakai secret baru
docker compose up -d backend
```

### 4. Akses Aplikasi

```bash
# Buka browser
http://localhost:5173
```

### 5. Stop & Cleanup

```bash
# Stop semua service
docker compose down

# Stop + hapus volume (⚠️ data DB hilang)
docker compose down -v
```

---

## Quick Start (Local)

### 1. Install Dependencies

**Windows:**
- **Go 1.22+**: [https://go.dev/dl/](https://go.dev/dl/)
- **PostgreSQL 16**: [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)
- **Node.js 18+**: [https://nodejs.org/](https://nodejs.org/)

Atau pakai Chocolatey:
```powershell
choco install go postgresql nodejs
```

**Ubuntu / Debian:**
```bash
sudo apt update
sudo apt install -y golang-1.22 postgresql-16 nodejs npm
```

Atau install manual:
- **Go 1.22+**: [https://go.dev/dl/](https://go.dev/dl/)
- **PostgreSQL 16**: [https://www.postgresql.org/download/linux/ubuntu/](https://www.postgresql.org/download/linux/ubuntu/)
- **Node.js 18+**: [https://nodejs.org/](https://nodejs.org/)

Verify:
```bash
go version          # minimal 1.22
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

Edit `.env`:
```
DB_DSN="host=localhost user=postgres password=YOUR_PASSWORD dbname=nexa port=5432 sslmode=disable"
JWT_SECRET="OUTPUT_DARI_OPENSSL"
SERVER_PORT=":8080"
COOKIE_DOMAIN="localhost"
COOKIE_SECURE="false"
```

Jalankan:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```
→ Backend jalan di `http://localhost:8080`

### 4. Setup Frontend

```bash
cd frontend
npm install
npm run dev
```
→ Frontend jalan di `http://localhost:5173`

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

## Development Roadmap

| Sprint          | Fokus                                  |
| :-------------- | :------------------------------------- |
| 1 (Pekan 1-2)   | Environment Setup & Database Modeling  |
| 2 (Pekan 3-4)   | Core Backend & Identity Management     |
| 3 (Pekan 5-6)   | Rekam Medis & Pencarian Kode Klinis    |
| 4 (Pekan 7-8)   | Logika Casemix & Engine Klaim          |
| 5 (Pekan 9-10)  | Frontend Development & API Integration |
| 6 (Pekan 11-12) | Testing, Security Hardening & Handover |
