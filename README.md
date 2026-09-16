## Nexa: Sistem Informasi Pengelolaan Data Casemix Terintegrasi

| Metadata | Keterangan |
| :--- | :--- |
| **Nama Proyek** | Nexa (Casemix Management System) |
| **Versi Dokumen** | 1.0.0 |
| **Status Dokumen** | Ready for Development / Approved |
| **Target Lingkungan** | On-Premise / Local Docker Virtualization (RSUD) |
| **Tech Stack Utama** | Go 1.22+ (Gin Gonic, GORM), Vue 3 (Vite, PrimeVue, Pinia), PostgreSQL 16 |

---

## 1. Executive Summary & Ringkasan Sistem

**Nexa** adalah aplikasi web rekam medis dan penagihan terintegrasi yang dirancang untuk mempercepat siklus operasional klaim BPJS Kesehatan di Rumah Sakit Umum Daerah (RSUD). Sistem ini menjembatani jurang data (*data silos*) antara loket pendaftaran, dokter penanggung jawab pelayanan (DPJP), perawat bangsal, koder Casemix, dan bagian keuangan.

Nexa mengadopsi arsitektur **Decoupled Monorepo** berbasis **RESTful API**:
* **Backend:** REST API berperforma tinggi menggunakan Go (Gin Gonic) dengan transaksi ACID PostgreSQL 16.
* **Frontend:** Single Page Application (SPA) reaktif berbasis Vue 3 (Composition API) + PrimeVue + Tailwind CSS.
* **Keamanan:** Kepatuhan OWASP Top 10, autentikasi sesi JWT via `HttpOnly` Cookie, otorisasi RBAC berlapis, mitigasi IDOR, dan pencatatan riwayat audit medis otomatis (*Audit Trail*).

---

## 2. Problem Statement & Latar Belakang Operasional

### 2.1 Latar Belakang
Unit Casemix merupakan pusat penerimaan pendapatan rumah sakit dari pasien peserta Jaminan Kesehatan Nasional (JKN/BPJS). Besaran penggantian klaim bergantung pada kombinasi ketepatan penegakan diagnosis (ICD-10), tindakan medis (ICD-9 CM), serta penentuan kode tarif INA-CBGs (*Indonesia Case Based Groups*).

### 2.2 Masalah Utama (*Pain Points*)
1. **Pencatatan Terfragmentasi & Silo Data:** Alur data manual antar loket pendaftaran, poliklinik/bangsal rawat, dan unit Casemix memicu inkonsistensi berkas rekam medis.
2. **Duplikasi & Human Error:** Kesalahan input kode diagnosis dan tindakan akibat input bebas tanpa validasi kamus ICD resmi dan ketiadaan *anti-duplication rule*.
3. **Pencarian Berkas Lambat:** Temu-balik rekam medis fisik memakan waktu 15–60 menit, sedangkan sistem digital lama membutuhkan waktu 5–15 menit per berkas.
4. **Tingginya Tingkat Klaim *Pending* / *Dispute*:** Berkas klaim tertahan atau ditolak oleh verifikator BPJS karena diskrepansi data atau ketiadaan catatan alasan penolakan yang tersentralisasi untuk perbaikan cepat.

---

## 3. Goals & Non-Goals

### 3.1 Project Goals
* **Integrasi End-to-End:** Mengonsolidasikan pendaftaran pasien, pencatatan rekam medis klinis, koding tarif INA-CBGs, dan pengarsipan klaim ke dalam satu basis data relasional.
* **Performa Pencarian Berkas:** Menurunkan latensi pencarian data pasien, berkas kunjungan, dan kamus master ICD menjadi $< 1\text{ detik}$ (maksimal $< 5\text{ detik}$ di bawah beban kerja tinggi).
* **Zero Duplicate Coding:** Mencegah insiden duplikasi kode diagnosis atau tindakan pada nomor kunjungan yang sama melalui *unique constraint database*.
* **Otomasi Rekomendasi Tarif INA-CBGs:** Memberikan pemetaan paket tarif klaim otomatis begitu diagnosis primer dan tindakan selesai diinput.
* **Pelacakan Siklus Klaim 100% Transparan:** Memastikan setiap perubahan status berkas klaim (*draft*, *pending*, *disetujui*, *ditolak*) terdokumentasi lengkap bersama justifikasi penolakan.
* **Audit Trail & Keamanan Ketat:** Mengamankan catatan rekam medis dari modifikasi tidak sah (mitigasi BOLA/IDOR) dan mencatat seluruh mutasi data secara otomatis.

### 3.2 Non-Goals (Out of Scope)
* **Farmasi & Logistik:** Tidak mencakup inventaris fisik gudang obat, resep racikan, atau bahan medis habis pakai (BMHP).
* **Direct Bridging Server BPJS VClaim:** Tidak melakukan panggilan API *direct production bridging* ke TrustMark/VClaim BPJS; sistem berfokus pada *grouping internal*, kalkulasi plafon, dan kesiapan berkas (*readiness*).
* **Akuntansi Umum (General Ledger):** Tidak memproses jurnal debit/kredit umum rumah sakit, hanya terbatas pada pencatatan nilai klaim Casemix dan status pencairannya.

---

## 4. Target User & Role-Based Access Control (RBAC)

Sistem memberlakukan kontrol akses berbasis peran (*Role-Based Access Control* / RBAC):

| Peran (Role) | Target Pengguna | Tanggung Jawab & Hak Akses Utama |
| :--- | :--- | :--- |
| `admin_ti` | Administrator TI | Manajemen pengguna (*user accounts*), konfigurasi sistem, audit trail log monitor. |
| `petugas_rm` | Petugas Registrasi & RM | Registrasi identitas pasien (CRUD), verifikasi NIK/No BPJS, input pendaftaran kunjungan. |
| `dokter_dpjp` | Dokter Penanggung Jawab | Input resume klinis kunjungan, anamnesis keluhan, penegakan diagnosis ICD-10 (primer/sekunder). |
| `perawat` | Perawat Bangsal / Poli | Input rincian tindakan medis (ICD-9 CM), kuantitas tindakan, dan catatan observasi. |
| `petugas_casemix` | Koder & Verifikator Klaim | Validasi koding medis, penentuan tarif paket INA-CBGs, pembaruan status klaim, input alasan dispute. |
| `keuangan` | Staf Kasir / Keuangan RS | Monitoring daftar klaim berstatus `disetujui`, rekonsiliasi nilai klaim, laporan keuangan Casemix. |
| `manajemen` | Direksi & Komite Medis | Akses *read-only* dasbor analitik, laporan 10 besar penyakit, tren rasio penerimaan klaim. |

### Matriks Otorisasi Modul

| Modul / Fitur | `admin_ti` | `petugas_rm` | `dokter_dpjp` | `perawat` | `petugas_casemix` | `keuangan` | `manajemen` |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Manajemen Pengguna & Unit | **CRUD** | - | - | - | - | - | - |
| Registrasi Pasien | R | **CRUD** | R | R | R | R | R |
| Rekam Medis (Kunjungan) | R | C/R | **CRUD** | R | R | R | R |
| Input Diagnosis (ICD-10) | R | - | **CRUD** | - | R | - | R |
| Input Prosedur (ICD-9) | R | - | R | **CRUD** | R | - | R |
| Verifikasi & Grouping CBGs | R | - | R | R | **CRUD** | R | R |
| Manajemen Status Klaim | R | - | - | - | **CRUD** | U (*read/verify*) | R |
| Audit Trail System Log | **Read** | - | - | - | - | - | - |
| Dasbor & Laporan Analitik | Full | Dasar | Dasar | Dasar | Full | Laporan Klaim | Full |

---

## 5. User Stories & Acceptance Criteria

* **US-01 (Pencegahan Duplikasi Pasien):**
  * *Sebagai* Petugas RM, *saya ingin* mencari data pasien via NIK (16 digit) atau No BPJS (13 digit) sebelum membuat data baru, *agar* tidak terjadi pencatatan rekam ganda.
  * **Acceptance Criteria:** Sistem memvalidasi keunikan NIK/No BPJS; jika sudah terdaftar, tampilkan peringatan dan buka profil yang sudah ada.
* **US-02 (Autocomplete Diagnosis DPJP):**
  * *Sebagai* Dokter DPJP, *saya ingin* mencari kode penyakit ICD-10 secara instan via teks diagnosis, *agar* penegakan diagnosis selesai dalam beberapa detik.
  * **Acceptance Criteria:** Input field memiliki *debounce* $300\text{ ms}$; pencarian mendukung *fuzzy matching* / trigram dengan respon $< 300\text{ ms}$; wajib menandai tepat 1 diagnosis primer per kunjungan.
* **US-03 (Pencatatan Prosedur Medis):**
  * *Sebagai* Perawat, *saya ingin* memasukkan seluruh tindakan medis (ICD-9 CM) beserta frekuensi kuantitasnya pada kunjungan rawat, *agar* data billing tindakan lengkap.
  * **Acceptance Criteria:** Sistem menolak nilai jumlah tindakan $\le 0$; mencegah duplikasi kode tindakan yang sama pada satu `id_rekam`.
* **US-04 (Grouping Otomatis INA-CBGs):**
  * *Sebagai* Petugas Casemix, *saya ingin* sistem memberikan rekomendasi kode paket tarif INA-CBGs berdasarkan kombinasi diagnosis primer dan tindakan yang dipilih, *agar* durasi koding berkurang minimal 50%.
  * **Acceptance Criteria:** Sistem menampilkan kalkulasi tarif acuan INA-CBGs secara otomatis saat membuka berkas rekam medis terkait.
* **US-05 (Logging Alasan Klaim Pending / Tolak):**
  * *Sebagai* Petugas Casemix, *saya ingin* mengisi uraian kendala saat mengubah status klaim menjadi `pending` atau `ditolak`, *agar* dokter/perawat mengetahui letak perbaikan berkas.
  * **Acceptance Criteria:** Kolom `alasan_pending_tolak` bersifat wajib diisi (*mandatory*) apabila status klaim diubah ke `pending` atau `ditolak`.
* **US-06 (Rekonsiliasi Keuangan Klaim Disetujui):**
  * *Sebagai* Staf Keuangan, *saya ingin* memfilter daftar klaim `disetujui` berdasarkan rentang tanggal tertentu dan mengekspornya ke format tabular, *agar* proses rekonsiliasi berjalan presisi.
  * **Acceptance Criteria:** Tersedia filter rentang tanggal, agregasi total nominal klaim disetujui, dan tombol ekspor data.
* **US-07 (Executive Dashboard):**
  * *Sebagai* Direksi / Manajemen, *saya ingin* melihat grafik visualisasi 10 penyakit terbanyak dan rasio klaim secara real-time, *agar* dapat mengambil keputusan strategis rumah sakit.
  * **Acceptance Criteria:** Dasbor menampilkan metrik jumlah pasien, pie-chart status klaim, dan bar-chart Top 10 ICD-10 teratas.
* **US-08 (Integritas & Audit Trail):**
  * *Sebagai* Admin TI, *saya ingin* seluruh operasi mutasi (INSERT, UPDATE, DELETE) pada rekam medis tercatat secara otomatis, *agar* kepatuhan hukum rekam medis terjamin.
  * **Acceptance Criteria:** Tabel `log_aktivitas` mencatat user ID, aksi, nama tabel, ID record terkait, dan stempel waktu presisi.

---

## 6. Functional Requirements (FR)

* **FR-01 (Autentikasi & Sesi Aman):**
  * Hashing sandi wajib menggunakan `bcrypt` dengan *cost factor* $\ge 10$.
  * Penerbitan token autentikasi JWT (`RS256` atau `HS256`) dengan masa aktif terkonfigurasi (misal: 8–24 jam).
  * JWT disimpan secara eksklusif dalam cookie dengan atribut `HttpOnly = true`, `SameSite = Lax`, dan `Secure = true` (pada lingkungan HTTPS).
* **FR-02 (Manajemen Pasien Terpadu):**
  * CRUD identitas pasien: NIK (wajib 16 digit angka), No BPJS (opsional, 13 digit angka), Nama Lengkap, Tanggal Lahir, Jenis Kelamin (`L`/`P`), Alamat, No HP.
  * Fitur pencarian instan berbasis NIK, Nomor BPJS, atau Nama Pasien.
* **FR-03 (Manajemen Rekam Medis & Rawat Inap/Jalan):**
  * Pencatatan kunjungan medis: Relasi ke `id_pasien` dan `id_dokter` (DPJP).
  * Pemilihan jenis perawatan: `rawat_jalan` atau `rawat_inap`.
  * Otomasi kalkulasi *Length of Stay* (LoS) untuk rawat inap:
    $$\text{LoS} = \max\left(1, \lceil \text{tanggal\_pulang} - \text{tanggal\_kunjungan} \rceil\right)$$
  * Validasi integritas waktu: $\text{tanggal\_pulang} \ge \text{tanggal\_kunjungan}$.
* **FR-04 (Master Data & Junction Medis ICD-10 / ICD-9 CM):**
  * Pencarian cepat master ICD-10 (diagnosis) dan ICD-9 CM (tindakan) menggunakan indeks trigram PostgreSQL.
  * Klasifikasi diagnosis: Tepat 1 diagnosis berstatus `primer`, diagnosis lainnya berstatus `sekunder`.
  * Pencegahan entri duplikat pada tingkat kunjungan yang sama.
* **FR-05 (Engine Grouping INA-CBGs):**
  * Mesin rekomendasi tarif internal: Memetakan kombinasi diagnosis primer (`kode_icd10`) dan tindakan (`kode_tindakan`) ke master `tarif_cbgs`.
  * Menampilkan deskripsi kasus dan plafon nominal tarif klaim BPJS secara real-time.
* **FR-06 (Siklus Berkas Klaim Casemix):**
  * Pengelolaan alur status klaim (*State Machine*): `draft` $\rightarrow$ `pending` $\rightarrow$ `disetujui` / `ditolak`.
  * Validasi wajib mengisi teks `alasan_pending_tolak` ketika status adalah `pending` atau `ditolak`.
  * Relasi unik satu-ke-satu (*1-to-1*) antara `rekam_medis` dan berkas `klaim`.
* **FR-07 (Dasbor Analisis & Pelaporan Real-time):**
  * Agregasi statistik: Total klaim bulan berjalan, rasio klaim diterima vs pending/ditolak, total pendapatan klaim disetujui.
  * Grafik visualisasi: 10 Diagnosis terbanyak (Top 10 ICD-10) berbasis tanggal kunjungan.
* **FR-08 (Audit Trail Otomatis):**
  * Pencatatan otomatis ke tabel `log_aktivitas` untuk setiap operasi mutasi data sensitif (`INSERT`, `UPDATE`, `DELETE`) dan sesi pengguna (`LOGIN`, `LOGOUT`).

---

## 7. Database Architecture & DDL (PostgreSQL 16)

Basis data dinormalisasi hingga tahap 3NF, dilengkapi kunci asing, aturan integritas (*Check Constraints*), serta mekanisme *Soft Delete* (`deleted_at`).

```
+---------------+       +------------------+       +---------------+
|     roles     |       |      users       |       |     units     |
+---------------+       +------------------+       +---------------+
| id_role (PK)  |<------| id_user (PK)     |------>| id_unit (PK)  |
| nama_role     |       | id_role (FK)     |       | nama_unit     |
+---------------+       | id_unit (FK)     |       +---------------+
                        +------------------+
                                 |
                                 | (id_dokter)
                                 v
+---------------+       +------------------+       +---------------------+
|    pasien     |<------|   rekam_medis    |------>|        klaim        |
+---------------+       +------------------+       +---------------------+
| id_pasien(PK) |       | id_rekam (PK)    |       | id_klaim (PK)       |
| nik (UQ)      |       | id_pasien (FK)   |       | id_rekam (FK, UQ)   |
| no_bpjs (UQ)  |       | id_dokter (FK)   |       | kode_cbgs (FK)      |
+---------------+       +------------------+       | id_petugas (FK)     |
                                |      |           | status_klaim        |
        +-----------------------+      +----+      +---------------------+
        |                                   |                 |
        v                                   v                 v
+------------------+               +-----------------+ +---------------+
| rekam_diagnosis  |               | detail_tindakan | |  tarif_cbgs   |
+------------------+               +-----------------+ +---------------+
| id (PK)          |               | id_detail (PK)  | | kode_cbgs(PK) |
| id_rekam (FK)    |               | id_rekam (FK)   | | deskripsi     |
| kode_icd10 (FK)  |               | kode_tindak(FK) | | tarif         |
| jenis            |               | jumlah          | +---------------+
+------------------+               +-----------------+
        |                                   |
        v                                   v
+------------------+               +-----------------+
|    diagnosis     |               |    tindakan     |
+------------------+               +-----------------+
| kode_icd10 (PK)  |               | kode_tindak(PK) |
| nama_diagnosis   |               | nama_tindakan   |
+------------------+               +-----------------+
```

### 7.1 Skrip DDL Lengkap (PostgreSQL 16)

```sql
-- Inisialisasi Ekstensi Trigram untuk Pencarian Cepat Teks Medis
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 1. Tabel Roles
CREATE TABLE roles (
    id_role SERIAL PRIMARY KEY,
    nama_role VARCHAR(50) UNIQUE NOT NULL
);

-- 2. Tabel Units
CREATE TABLE units (
    id_unit SERIAL PRIMARY KEY,
    nama_unit VARCHAR(50) UNIQUE NOT NULL
);

-- 3. Tabel Users
CREATE TABLE users (
    id_user SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    profesi VARCHAR(50) NOT NULL,
    spesialisasi VARCHAR(50),
    no_str VARCHAR(30) UNIQUE,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    id_role INT NOT NULL REFERENCES roles(id_role) ON DELETE RESTRICT,
    id_unit INT NOT NULL REFERENCES units(id_unit) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 4. Tabel Pasien
CREATE TABLE pasien (
    id_pasien SERIAL PRIMARY KEY,
    nik CHAR(16) UNIQUE NOT NULL,
    no_bpjs VARCHAR(13) UNIQUE,
    nama VARCHAR(100) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin CHAR(1) NOT NULL CHECK (jenis_kelamin IN ('L', 'P')),
    alamat TEXT NOT NULL,
    no_hp VARCHAR(15),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 5. Tabel Rekam Medis
CREATE TABLE rekam_medis (
    id_rekam SERIAL PRIMARY KEY,
    id_pasien INT NOT NULL REFERENCES pasien(id_pasien) ON DELETE RESTRICT,
    id_dokter INT NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,
    tanggal_kunjungan TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tanggal_pulang TIMESTAMPTZ,
    jenis_perawatan VARCHAR(20) NOT NULL CHECK (jenis_perawatan IN ('rawat_jalan', 'rawat_inap')),
    keluhan TEXT NOT NULL,
    catatan TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_tanggal_pulang CHECK (tanggal_pulang IS NULL OR tanggal_pulang >= tanggal_kunjungan)
);

-- 6. Tabel Master Diagnosis (ICD-10)
CREATE TABLE diagnosis (
    kode_icd10 VARCHAR(10) PRIMARY KEY,
    nama_diagnosis VARCHAR(255) NOT NULL,
    kategori VARCHAR(100)
);

-- 7. Tabel Relasi Rekam Diagnosis (Junction Table)
CREATE TABLE rekam_diagnosis (
    id SERIAL PRIMARY KEY,
    id_rekam INT NOT NULL REFERENCES rekam_medis(id_rekam) ON DELETE CASCADE,
    kode_icd10 VARCHAR(10) NOT NULL REFERENCES diagnosis(kode_icd10) ON DELETE RESTRICT,
    jenis VARCHAR(10) NOT NULL CHECK (jenis IN ('primer', 'sekunder')),
    CONSTRAINT uq_rekam_icd10 UNIQUE (id_rekam, kode_icd10)
);

-- 8. Tabel Master Tindakan (ICD-9 CM)
CREATE TABLE tindakan (
    kode_tindakan VARCHAR(10) PRIMARY KEY,
    nama_tindakan VARCHAR(255) NOT NULL,
    tarif_standar NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (tarif_standar >= 0)
);

-- 9. Tabel Relasi Detail Tindakan (Junction Table)
CREATE TABLE detail_tindakan (
    id_detail SERIAL PRIMARY KEY,
    id_rekam INT NOT NULL REFERENCES rekam_medis(id_rekam) ON DELETE CASCADE,
    kode_tindakan VARCHAR(10) NOT NULL REFERENCES tindakan(kode_tindakan) ON DELETE RESTRICT,
    jumlah INT NOT NULL DEFAULT 1 CHECK (jumlah > 0),
    CONSTRAINT uq_rekam_tindakan UNIQUE (id_rekam, kode_tindakan)
);

-- 10. Tabel Master Tarif INA-CBGs
CREATE TABLE tarif_cbgs (
    kode_cbgs VARCHAR(15) PRIMARY KEY,
    deskripsi VARCHAR(255) NOT NULL,
    tarif NUMERIC(14, 2) NOT NULL CHECK (tarif >= 0)
);

-- 11. Tabel Klaim Casemix
CREATE TABLE klaim (
    id_klaim SERIAL PRIMARY KEY,
    id_rekam INT UNIQUE NOT NULL REFERENCES rekam_medis(id_rekam) ON DELETE RESTRICT,
    kode_cbgs VARCHAR(15) NOT NULL REFERENCES tarif_cbgs(kode_cbgs) ON DELETE RESTRICT,
    id_petugas_casemix INT NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,
    status_klaim VARCHAR(20) NOT NULL CHECK (status_klaim IN ('draft', 'pending', 'disetujui', 'ditolak')),
    tanggal_klaim DATE NOT NULL DEFAULT CURRENT_DATE,
    nominal_klaim NUMERIC(14, 2) NOT NULL CHECK (nominal_klaim >= 0),
    alasan_pending_tolak TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_alasan_pending_tolak CHECK (
        (status_klaim IN ('pending', 'ditolak') AND alasan_pending_tolak IS NOT NULL AND LENGTH(TRIM(alasan_pending_tolak)) > 0)
        OR (status_klaim IN ('draft', 'disetujui'))
    )
);

-- 12. Tabel Log Aktivitas (Audit Trail)
CREATE TABLE log_aktivitas (
    id_log SERIAL PRIMARY KEY,
    id_user INT REFERENCES users(id_user) ON DELETE SET NULL,
    tabel_terdampak VARCHAR(50) NOT NULL,
    aktivitas VARCHAR(20) NOT NULL CHECK (aktivitas IN ('INSERT', 'UPDATE', 'DELETE', 'LOGIN', 'LOGOUT')),
    data_id INT,
    waktu TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Pembuatan Indeks Performa & Trigram
CREATE INDEX idx_pasien_nik ON pasien(nik);
CREATE INDEX idx_pasien_no_bpjs ON pasien(no_bpjs);
CREATE INDEX idx_rekam_medis_pasien ON rekam_medis(id_pasien);
CREATE INDEX idx_rekam_medis_dokter ON rekam_medis(id_dokter);
CREATE INDEX idx_klaim_status ON klaim(status_klaim);
CREATE INDEX idx_log_waktu ON log_aktivitas(waktu);

-- GIN Trigram Indexes untuk Autocomplete Instan (< 300 ms)
CREATE INDEX idx_trgm_diagnosis_nama ON diagnosis USING gin (nama_diagnosis gin_trgm_ops);
CREATE INDEX idx_trgm_tindakan_nama ON tindakan USING gin (nama_tindakan gin_trgm_ops);
```

---

## 8. API Specifications & Endpoints Contract

Seluruh komunikasi API mengikuti kaidah standar RESTful dengan format muatan JSON:

```json
{
  "success": true,
  "message": "Deskripsi status operasi",
  "data": {}
}
```

### 8.1 Modul Autentikasi (`/api/v1/auth`)
* `POST /api/v1/auth/login`
  * **Rate Limit:** 5 request / menit / IP.
  * **Body:** `{"username": "...", "password": "..."}`
  * **Response:** Mengeset Cookie `jwt` (`HttpOnly`, `Secure`, `SameSite=Lax`), payload user info.
* `POST /api/v1/auth/logout`
  * **Response:** Mengosongkan cookie `jwt`, mencatat aktivitas `LOGOUT` ke `log_aktivitas`.
* `GET /api/v1/auth/me`
  * **Auth:** Cookie JWT.
  * **Response:** Identitas staf login, data role, dan unit.

### 8.2 Modul Pasien (`/api/v1/pasien`)
* `GET /api/v1/pasien?search=&page=1&limit=10` — Mengambil daftar pasien terpaginasi.
* `GET /api/v1/pasien/:id` — Detail pasien berdasarkan `id_pasien`.
* `POST /api/v1/pasien` — Mendaftarkan pasien baru (Validasi: NIK 16 digit, format BPJS).
* `PUT /api/v1/pasien/:id` — Memperbarui data profil pasien.
* `DELETE /api/v1/pasien/:id` — *Soft delete* data pasien (Role: `petugas_rm`, `admin_ti`).

### 8.3 Modul Rekam Medis (`/api/v1/rekam-medis`)
* `GET /api/v1/rekam-medis?pasien_id=&status=&page=1` — Daftar riwayat rekam medis.
* `GET /api/v1/rekam-medis/:id` — Data detail rekam medis, daftar ICD-10, dan rincian ICD-9.
* `POST /api/v1/rekam-medis` — Membuka kunjungan baru rawat inap/jalan.
* `PUT /api/v1/rekam-medis/:id` — Memperbarui data kunjungan / menetapkan tanggal pulang.
* `POST /api/v1/rekam-medis/:id/diagnosis` — Menambahkan diagnosis ICD-10 (primer/sekunder).
* `DELETE /api/v1/rekam-medis/:id/diagnosis/:id_rekam_diag` — Menghapus rincian diagnosis.
* `POST /api/v1/rekam-medis/:id/tindakan` — Menambahkan prosedur tindakan ICD-9 CM.
* `DELETE /api/v1/rekam-medis/:id/tindakan/:id_detail` — Menghapus tindakan dari rekam medis.

### 8.4 Modul Master Medis & Autocomplete (`/api/v1/master`)
* `GET /api/v1/master/icd10?q={query}` — Pencarian diagnosis (Trigram fuzzy search, batas 20 hasil).
* `GET /api/v1/master/icd9?q={query}` — Pencarian prosedur tindakan (Trigram fuzzy search, batas 20 hasil).
* `GET /api/v1/master/cbg-recommendation?icd10_primer={kode}&tindakan={kode1,kode2}` — Mengembalikan rekomendasi kode paket `tarif_cbgs`.

### 8.5 Modul Casemix & Klaim (`/api/v1/klaim`)
* `GET /api/v1/klaim?status=&tanggal_mulai=&tanggal_selesai=&page=1` — Filter berkas klaim.
* `GET /api/v1/klaim/:id` — Mengambil detail kalkulasi klaim beserta data rekam medis.
* `POST /api/v1/klaim` — Membuat draft klaim baru dari rekam medis yang sudah lengkap.
* `PATCH /api/v1/klaim/:id/status` — Transisi status klaim (`draft`, `pending`, `disetujui`, `ditolak`).
  * **Payload:** `{"status_klaim": "pending", "alasan_pending_tolak": "Resume medis kurang tanda tangan DPJP"}`

### 8.6 Modul Dasbor & Laporan (`/api/v1/analytics`)
* `GET /api/v1/analytics/summary` — Menampilkan KPI: Pasien aktif, Total klaim, Klaim dispute.
* `GET /api/v1/analytics/top-diseases?limit=10&year=2026` — Data grafik 10 besar diagnosis ICD-10.
* `GET /api/v1/analytics/export-klaim?format=csv&status=disetujui` — Ekspor data tabular klaim.

---

## 9. Casemix Workflow & State Machine

```
   +-------------------------------------------------------------+
   | Loket RM: Pendaftaran Pasien & Buka Kunjungan (rekam_medis) |
   +-------------------------------------------------------------+
                                  |
                                  v
   +-------------------------------------------------------------+
   | Dokter DPJP: Anamnesis & Penegakan Diagnosis (ICD-10)       |
   | (Primer = 1, Sekunder = 0..N)                               |
   +-------------------------------------------------------------+
                                  |
                                  v
   +-------------------------------------------------------------+
   | Perawat Bangsal: Pencatatan Rincian Tindakan (ICD-9 CM)     |
   +-------------------------------------------------------------+
                                  |
                                  v
   +-------------------------------------------------------------+
   | Koder Casemix: Otomasi Grouping INA-CBGs & Pembuatan Klaim  |
   | State: DRAFT                                                |
   +-------------------------------------------------------------+
                                  |
            +---------------------+---------------------+
            |                                           |
            v (Berkas Lengkap & Lolos)                  v (Ada ketidaksesuaian berkas)
   +-------------------+                     +------------------------------------+
   | State: DISETUJUI  |                     | State: PENDING / DITOLAK           |
   +-------------------+                     | (Wajib isi alasan_pending_tolak)   |
            |                                +------------------------------------+
            v                                           |
   +-------------------+                                v
   | Rekonsiliasi      |                     +------------------------------------+
   | Staf Keuangan     |                     | Revisi berkas oleh DPJP / Casemix  |
   +-------------------+                     +------------------------------------+
```

### Aturan Transisi Status Klaim
1. **DRAFT:** Berkas baru dibentuk oleh koder Casemix. Dapat diedit sepenuhnya.
2. **PENDING:** Berkas ditangguhkan karena syarat verifikasi klinis atau administrasi belum lengkap. Nilai `alasan_pending_tolak` tidak boleh kosong.
3. **DITOLAK:** Berkas tidak memenuhi kriteria jaminan klaim BPJS. Nilai `alasan_pending_tolak` wajib mencantumkan klausul penolakan.
4. **DISETUJUI:** Berkas terverifikasi final. Nilai `nominal_klaim` dikunci dan siap ditagihkan/direkonsiliasi oleh tim keuangan.

---

## 10. Non-Functional Requirements (NFR)

* **NFR-01 (Keamanan & Proteksi Data):**
  * Standar OWASP Top 10 terpenuhi secara sistematis.
  * Autentikasi sesi berbasis HttpOnly cookie untuk mencegah pencurian token lewat skrip Cross-Site Scripting (XSS).
  * Proteksi Broken Object Level Authorization (BOLA/IDOR): Pemeriksaan validitas kepemilikan data dokter pada lapisan logika usecase.
* **NFR-02 (Performa & Kecepatan Akses):**
  * Response time pencarian autocomplete master ICD $< 300\text{ ms}$ memanfaatkan indeks GIN Trigram.
  * Query data detail pasien dan riwayat kunjungan dieksekusi $< 1\text{ detik}$.
  * Skalabilitas beban konkurensi: Mampu melayani sedikitnya 100 *concurrent requests* tanpa degradasi performa (didukung Go Goroutines & Gin engine).
* **NFR-03 (Integritas Basis Data):**
  * Transaksi simpan klaim dan junction medis menerapkan *ACID Transaction* (`Begin()`, `Commit()`, `Rollback()` di GORM).
  * Menjamin tidak ada entitas data yatim (*orphaned data*) dengan penerapan foreign key constraints ketat (`ON DELETE RESTRICT` untuk data transaksional induk).
* **NFR-04 (Ketersediaan & Keandalan):**
  * Log error aplikasi terpusat dengan level severity (`INFO`, `WARN`, `ERROR`).
  * Middleware Gin `Recovery()` aktif untuk mencegah aplikasi *crash* akibat runtime panic.

---

## 11. Web Security Architecture & OWASP Compliance

### 11.1 Matriks Implementasi OWASP Top 10

| OWASP Category | Kerentanan Potensial | Solusi Arsitektural pada Nexa |
| :--- | :--- | :--- |
| **A01:2021 - Broken Access Control** | Pengguna mengakses atau mengubah data milik dokter/pasien lain (BOLA/IDOR). | RBAC Middleware di Gin + validasi usecase `if rekam.IDDokter != currentUser.ID && !currentUser.IsAdmin { return ErrForbidden }`. |
| **A02:2021 - Cryptographic Failures** | Kebocoran kata sandi dan manipulasi sesi JWT. | Enkripsi kata sandi via `bcrypt` ($cost \ge 10$); JWT ditandatangani secara kriptografis; transmisi wajib HTTPS. |
| **A03:2021 - Injection** | SQL Injection via form filter dan parameter pencarian. | Pemakaian GORM Parameterized Queries / Prepared Statements secara konsisten; sanitasi karakter input. |
| **A04:2021 - Insecure Design** | Manipulasi status klaim atau data ganda yang merugikan keuangan RS. | *Check Constraints* pada skema basis data, *anti-duplicate unique keys*, serta *audit trail* rekaman medis. |
| **A05:2021 - Security Misconfiguration** | Port terbuka, *debug error stack trace* terekspos ke klien. | Mode produksi Gin `GIN_MODE=release`; respons galat seragam tanpa membocorkan struktur internal skema database. |
| **A07:2021 - Identification & Auth Failures** | Serangan *Brute Force Login* pada portal tenaga kesehatan. | Rate limiter middleware pada rute `/api/v1/auth/login` (maksimal 5 percobaan gagal per menit per alamat IP). |

### 11.2 Spesifikasi Proteksi Sesi Cookie & CORS
```go
// Konfigurasi Cookie pada Handler Login Go Gin
c.SetSameSite(http.SameSiteLaxMode)
c.SetCookie(
    "jwt",          // Nama cookie
    tokenString,    // Nilai JWT
    3600 * 8,       // Masa berlaku (8 jam)
    "/",            // Path
    "",             // Domain (kosongkan untuk host saat ini)
    true,           // Secure: true (hanya lewat HTTPS)
    true,           // HttpOnly: true (mencegah akses document.cookie di JS)
)
```

Konfigurasi CORS pada backend Go wajib membatasi origin spesifik (misal: `http://localhost:5173` saat tahap dev) dengan opsi `AllowCredentials: true`.

---

## 12. UI/UX & Design Guidelines

1. **Prinsip Tampilan Dashboard Medis:**
   * Desain minimalis, bersih (*clinical look*), kontras warna memenuhi standar aksesibilitas WCAG 2.1 AA.
   * Framework UI: **PrimeVue** (Komponen DataTable, AutoComplete, Dialog, Tag) dipadukan dengan **Tailwind CSS**.
2. **Navigasi Dinamis (Role-Based Rendering):**
   * Bilah menu samping (*sidebar*) hanya me-render menu navigasi sesuai izin peran token pengguna yang disimpan pada state Pinia.
3. **Pencarian Cepat dengan Debounce:**
   * Seluruh komponen input teks master diagnosis dan tindakan ICD menggunakan fungsi *debounce* ($300\text{ ms}$) untuk meminimalkan beban komputasi server.
4. **Status Badges & Notifikasi:**
   * Penanda status klaim berbasis warna lencana (*PrimeVue Tag*):
     * `draft`: Biru / Slate
     * `pending`: Oranye / Kuning (*Warning*)
     * `disetujui`: Hijau (*Success*)
     * `ditolak`: Merah (*Danger*)
   * Umpan balik aksi menggunakan *PrimeVue Toast* (sukses/gagal) dan konfirmasi modal sebelum aksi penting (*Soft Delete* / Penolakan Berkas).

---

## 13. System Architecture & Repository Structure

### 13.1 Struktur Decoupled Monorepo
```
nexa-casemix/
├── docker-compose.yml
├── README.md
│
├── backend/                       # Go 1.22+ RESTful API
│   ├── cmd/
│   │   └── api/
│   │       └── main.go            # Entry point backend
│   ├── internal/
│   │   ├── config/                # Konfigurasi env & database
│   │   ├── middleware/            # JWT, RBAC, RateLimiter, CORS
│   │   ├── model/                 # Struct entitas GORM & DTO
│   │   ├── repository/            # Lapisan akses database query
│   │   ├── usecase/               # Logika bisnis & validasi IDOR
│   │   └── delivery/http/         # Gin HTTP Handlers & Routes
│   ├── pkg/
│   │   ├── token/                 # JWT helper
│   │   └── response/              # Format JSON standar
│   ├── go.mod
│   └── go.sum
│
└── frontend/                      # Vue 3 Single Page Application
    ├── src/
    │   ├── assets/                # Styling Tailwind & icons
    │   ├── components/            # Reusable UI (Navbar, Sidebar, Modals)
    │   ├── layouts/               # Layout per peran pengguna
    │   ├── router/                # Vue Router & Navigation Guards
    │   ├── stores/                # Pinia (authStore, claimStore)
    │   ├── services/              # Axios HTTP client & Interceptors
    │   ├── views/                 # Halaman utama (Dashboard, Pasien, Klaim)
    │   ├── App.vue
    │   └── main.js
    ├── package.json
    └── vite.config.js
```

### 13.2 Konfigurasi Docker Compose (`docker-compose.yml`)
```yaml
version: '3.8'

services:
  postgres_db:
    image: postgres:16-alpine
    container_name: nexa_postgres
    restart: always
    environment:
      POSTGRES_DB: nexa_casemix
      POSTGRES_USER: nexa_admin
      POSTGRES_PASSWORD: a_secure_password
    ports:
      - "5432:5432"
    volumes:
      - nexa_pgdata:/var/lib/postgresql/data
      - ./backend/scripts/init.sql:/docker-entrypoint-initdb.d/init.sql

volumes:
  nexa_pgdata:
    driver: local
```

---

## 14. Key Performance Indicators (Success Metrics)

| Metrik Evaluasi | Kondisi Eksisting (Sebelum Nexa) | Target Nexa | Metode Pengukuran |
| :--- | :--- | :--- | :--- |
| **Kecepatan Pencarian Berkas** | $15 - 60\text{ menit}$ (Fisik) / $5 - 15\text{ menit}$ (Manual PC) | $< 5\text{ detik}$ (Optimal $< 1\text{ detik}$) | Pengujian waktu respon antarmuka dan *query execution plan* (EXPLAIN ANALYZE). |
| **Integritas Kode Medis** | Sering terjadi duplikasi koding pada lembar rekam | $0\%$ Duplikasi Kode Medis | Audit basis data pada tabel `rekam_diagnosis` dan `detail_tindakan`. |
| **Efisiensi Koding Casemix** | $10 - 20\text{ menit}$ per berkas klaim | Terpangkas minimal $50\%$ ($\le 5\text{ menit}$) | Observasi *time-to-complete* pada alur penentuan kode INA-CBGs. |
| **Transparansi Berkas Dispute** | Tidak tercatat sistematis (sering tercecer) | $100\%$ Berkas Pending/Tolak memiliki catatan alasan | Pengecekan constraint `chk_alasan_pending_tolak` pada tabel `klaim`. |
| **Audit Trail Compliance** | Tidak ada pencatatan mutasi data terpusat | $100\%$ Aksi CUD tercatat di `log_aktivitas` | Verifikasi integritas log pasca pengujian fungsional modul rekam medis. |

---

## 15. Risks & Mitigation Matrix

| Risiko | Tingkat Risiko | Dampak Potensial | Strategi Mitigasi Teknis |
| :--- | :---: | :--- | :--- |
| **Latensi Pencarian Kamus ICD** | Sedang | UI freeze atau lag ketika koder mengetik di input diagnosis. | Pasang indeks GIN Trigram pada PostgreSQL; implementasikan teknik *debounce* ($300\text{ ms}$) pada komponen PrimeVue AutoComplete. |
| **Inkonsistensi Status & Nilai Klaim** | Tinggi | Selisih nilai tagihan RS dengan tarif acuan INA-CBGs BPJS. | Kunci transaksi dengan *ACID transactions*; buat relasi 1-to-1 yang unik antara `id_rekam` dan `klaim`. |
| **Penghapusan Berkas Medis Sensitif** | Sangat Tinggi | Hilangnya riwayat rekam medis pasien yang melanggar hukum. | Terapkan *Soft Delete* (`deleted_at`) di GORM; nonaktifkan perintah SQL `HARD DELETE` secara permanen pada tabel transaksi. |
| **Eksploitasi BOLA / IDOR** | Sangat Tinggi | Dokter/pengguna tidak bertanggung jawab mengubah data pasien lain. | Validasi kepemilikan objek pada usecase Go sebelum mutasi basis data diproses. |
| **Serangan Brute Force Login** | Tinggi | Akun staf administrasi atau koder Casemix berhasil diretas. | Terapkan Gin Rate Limiter middleware (maks. 5 kali percobaan / menit / IP); enkripsi bcrypt $cost \ge 10$. |

---

## 16. Development Roadmap & Sprints

* **Sprint 1 (Pekan 1–2): Environment Setup & Database Modeling**
  * Setup repositori monorepo dan Docker Compose PostgreSQL 16.
  * Eksekusi migrasi DDL lengkap, pembuatan indeks trigram, dan penyusunan data awal (*seeders*: roles, units, master ICD-10 & ICD-9 dasar).
* **Sprint 2 (Pekan 3–4): Core Backend & Identity Management**
  * Implementasi arsitektur berlapis Go Gin (Model, Repository, Usecase, Handler).
  * Pembuatan middleware autentikasi JWT via HttpOnly cookie, middleware RBAC, dan rate limiter.
  * Endpoint CRUD data pasien beserta validasi NIK/No BPJS.
* **Sprint 3 (Pekan 5–6): Rekam Medis & Pencarian Kode Klinis**
  * Endpoint CRUD kunjungan rekam medis terintegrasi.
  * Fitur penambahan diagnosis ICD-10 (primer/sekunder) dan tindakan ICD-9 CM anti-duplikasi.
  * Endpoint master autocomplete dengan latensi respon $< 300\text{ ms}$.
  * Implementasi pencegahan BOLA/IDOR pada modul rekam medis.
* **Sprint 4 (Pekan 7–8): Logika Casemix & Engine Klaim**
  * Logika pencocokan tarif INA-CBGs otomatis.
  * Pengelolaan siklus hidup klaim (*State Machine*: draft, pending, disetujui, ditolak).
  * Validasi mutlak kolom alasan pending/tolak.
  * Implementasi trigger/interceptor pencatatan mutasi ke `log_aktivitas`.
* **Sprint 5 (Pekan 9–10): Frontend Development & API Integration**
  * Inisialisasi antarmuka Vue 3 + Tailwind CSS + PrimeVue.
  * Konfigurasi Pinia auth store, Axios interceptor credentials, dan Vue Router navigation guards.
  * Halaman manajemen pasien, formulir rekam medis interaktif, dan dasbor verifikasi Casemix.
  * Integrasi Chart.js untuk visualisasi 10 besar diagnosis penyakit.
* **Sprint 6 (Pekan 11–12): Testing, Security Hardening & Handover**
  * Pengujian fungsional *End-to-End* (E2E) dan *Black-Box Testing*.
  * Uji penetrasi mandiri terhadap ancaman OWASP Top 10 (SQLi, IDOR, XSS, Brute Force).
  * Optimasi query dan benchmarking waktu tanggap sistem.
  * Finalisasi dokumentasi teknis dan panduan operasional pengguna.
