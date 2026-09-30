-- Nexa Casemix Initial Schema Migration
-- PRD §7 Database Architecture & DDL (PostgreSQL 16)
-- Run this on a fresh database or after dropping existing tables

-- Enable pg_trgm extension for trigram search on medical master data
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 1. Tabel Roles
CREATE TABLE IF NOT EXISTS roles (
    id_role SERIAL PRIMARY KEY,
    nama_role VARCHAR(50) UNIQUE NOT NULL
);

-- 2. Tabel Units
CREATE TABLE IF NOT EXISTS units (
    id_unit SERIAL PRIMARY KEY,
    nama_unit VARCHAR(50) UNIQUE NOT NULL
);

-- 3. Tabel Users
CREATE TABLE IF NOT EXISTS users (
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
CREATE TABLE IF NOT EXISTS pasien (
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
CREATE TABLE IF NOT EXISTS rekam_medis (
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
CREATE TABLE IF NOT EXISTS diagnosis (
    kode_icd10 VARCHAR(10) PRIMARY KEY,
    nama_diagnosis VARCHAR(255) NOT NULL,
    kategori VARCHAR(100)
);

-- 7. Tabel Relasi Rekam Diagnosis (Junction Table)
CREATE TABLE IF NOT EXISTS rekam_diagnosis (
    id SERIAL PRIMARY KEY,
    id_rekam INT NOT NULL REFERENCES rekam_medis(id_rekam) ON DELETE CASCADE,
    kode_icd10 VARCHAR(10) NOT NULL REFERENCES diagnosis(kode_icd10) ON DELETE RESTRICT,
    jenis VARCHAR(10) NOT NULL CHECK (jenis IN ('primer', 'sekunder')),
    CONSTRAINT uq_rekam_icd10 UNIQUE (id_rekam, kode_icd10)
);

-- 8. Tabel Master Tindakan (ICD-9 CM)
CREATE TABLE IF NOT EXISTS tindakan (
    kode_tindakan VARCHAR(10) PRIMARY KEY,
    nama_tindakan VARCHAR(255) NOT NULL,
    tarif_standar NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (tarif_standar >= 0)
);

-- 9. Tabel Relasi Detail Tindakan (Junction Table)
CREATE TABLE IF NOT EXISTS detail_tindakan (
    id_detail SERIAL PRIMARY KEY,
    id_rekam INT NOT NULL REFERENCES rekam_medis(id_rekam) ON DELETE CASCADE,
    kode_tindakan VARCHAR(10) NOT NULL REFERENCES tindakan(kode_tindakan) ON DELETE RESTRICT,
    jumlah INT NOT NULL DEFAULT 1 CHECK (jumlah > 0),
    CONSTRAINT uq_rekam_tindakan UNIQUE (id_rekam, kode_tindakan)
);

-- 10. Tabel Master Tarif INA-CBGs
CREATE TABLE IF NOT EXISTS tarif_cbgs (
    kode_cbgs VARCHAR(15) PRIMARY KEY,
    deskripsi VARCHAR(255) NOT NULL,
    tarif NUMERIC(14, 2) NOT NULL CHECK (tarif >= 0)
);

-- 11. Tabel Klaim Casemix
CREATE TABLE IF NOT EXISTS klaim (
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
CREATE TABLE IF NOT EXISTS log_aktivitas (
    id_log SERIAL PRIMARY KEY,
    id_user INT REFERENCES users(id_user) ON DELETE SET NULL,
    tabel_terdampak VARCHAR(50) NOT NULL,
    aktivitas VARCHAR(20) NOT NULL CHECK (aktivitas IN ('INSERT', 'UPDATE', 'DELETE', 'LOGIN', 'LOGOUT')),
    data_id INT,
    waktu TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for Performance
CREATE INDEX IF NOT EXISTS idx_pasien_nik ON pasien(nik);
CREATE INDEX IF NOT EXISTS idx_pasien_no_bpjs ON pasien(no_bpjs);
CREATE INDEX IF NOT EXISTS idx_rekam_medis_pasien ON rekam_medis(id_pasien);
CREATE INDEX IF NOT EXISTS idx_rekam_medis_dokter ON rekam_medis(id_dokter);
CREATE INDEX IF NOT EXISTS idx_klaim_status ON klaim(status_klaim);
CREATE INDEX IF NOT EXISTS idx_log_waktu ON log_aktivitas(waktu);

-- GIN Trigram Indexes for Autocomplete (< 300ms)
CREATE INDEX IF NOT EXISTS idx_trgm_diagnosis_nama ON diagnosis USING gin (nama_diagnosis gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_trgm_tindakan_nama ON tindakan USING gin (nama_tindakan gin_trgm_ops);