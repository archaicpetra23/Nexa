# Nexa — Sistem Informasi Casemix Terintegrasi

**Nexa** adalah aplikasi web rekam medis dan penagihan terintegrasi yang dirancang untuk mempercepat siklus operasional klaim BPJS Kesehatan di Rumah Sakit Umum Daerah (RSUD). Sistem ini menjembatani *data silos* antara loket pendaftaran, dokter DPJP, perawat bangsal, koder Casemix, dan bagian keuangan.

## Tech Stack

- **Backend:** Go 1.22+ (Gin Gonic, GORM), PostgreSQL 16
- **Frontend:** Vue 3 (Composition API, Vite, PrimeVue, Tailwind CSS, Pinia)
- **Arsitektur:** Decoupled Monorepo, RESTful API, JWT + RBAC

## Struktur Dokumen

| File | Deskripsi |
| :--- | :--- |
| [`ide_proyek.md`](ide_proyek.md) | Ide proyek, masalah, target pengguna, fitur inti, dan kriteria keberhasilan |
| [`docs/product_requirement_document_prd_nexa.md`](docs/product_requirement_document_prd_nexa.md) | Product Requirement Document (PRD) lengkap dalam format Markdown |
| [`docs/Nexa (PRD).pdf`](docs/Nexa%20(PRD).pdf) | PRD dalam format PDF |

## Masalah yang Diselesaikan

1. **Pencatatan Terfragmentasi** — alur data manual antar-unit memicu inkonsistensi berkas rekam medis.
2. **Duplikasi & Human Error** — kesalahan input kode diagnosis (ICD-10) dan tindakan (ICD-9 CM) tanpa validasi otomatis.
3. **Pencarian Berkas Lambat** — temu-balik rekam medis fisik memakan waktu 15–60 menit.
4. **Klaim Pending/Dispute** — berkas tertahan tanpa pencatatan alasan penolakan terpusat.

## Target Pengguna (7 Role RBAC)

| Peran | Tanggung Jawab Utama |
| :--- | :--- |
| `admin_ti` | Manajemen pengguna, konfigurasi, audit trail |
| `petugas_rm` | Registrasi pasien, verifikasi NIK/No BPJS |
| `dokter_dpjp` | Input resume klinis, diagnosis ICD-10 |
| `perawat` | Input tindakan medis ICD-9 CM |
| `petugas_casemix` | Verifikasi koding, grouping INA-CBGs |
| `keuangan` | Monitoring klaim disetujui, rekonsiliasi |
| `manajemen` | Akses read-only dasbor analitik |

## Fitur Inti

1. **Autentikasi & RBAC Multi-Role** — sesi aman JWT via *HttpOnly Cookie*
2. **Manajemen Pasien** — CRUD + validasi NIK (16 digit), No BPJS (13 digit)
3. **Rekam Medis Kunjungan** — rawat jalan/inap + kalkulasi *Length of Stay*
4. **Pencarian Cepat ICD** — autocomplete *fuzzy search* < 300ms (indeks trigram)
5. **Engine Rekomendasi INA-CBGs** — pemetaan paket tarif otomatis
6. **Siklus Klaim (State Machine)** — `draft` → `pending` → `disetujui`/`ditolak`
7. **Dasbor Analitik** — Top 10 ICD-10 + metrik status klaim *real-time*
8. **Audit Trail** — pencatatan otomatis seluruh mutasi data klinis

## Keamanan

- OWASP Top 10 compliant
- Autentikasi JWT via *HttpOnly*, *Secure*, *SameSite=Lax* cookie
- RBAC middleware + mitigasi IDOR/BOLA
- *Rate limiter* login (5 percobaan/menit/IP)
- *Soft delete* + *ACID transaction* PostgreSQL

## Target Keberhasilan

- Pencarian berkas: **< 5 detik** (optimal < 1 detik)
- Duplikasi kode medis: **0%**
- Efisiensi koding Casemix: **terpangkas ≥ 50%**
- Transparansi dispute: **100%** berkas pending/tolak tercatat alasan
- Audit trail: **100%** aksi mutasi tercatat di `log_aktivitas`

## 🗺️ Development Roadmap

| Sprint | Fokus |
| :--- | :--- |
| 1 (Pekan 1–2) | Environment Setup & Database Modeling |
| 2 (Pekan 3–4) | Core Backend & Identity Management |
| 3 (Pekan 5–6) | Rekam Medis & Pencarian Kode Klinis |
| 4 (Pekan 7–8) | Logika Casemix & Engine Klaim |
| 5 (Pekan 9–10) | Frontend Development & API Integration |
| 6 (Pekan 11–12) | Testing, Security Hardening & Handover |

## Tim

Praktikum Rekayasa Perangkat Lunak — Semester 3, Pradita University
