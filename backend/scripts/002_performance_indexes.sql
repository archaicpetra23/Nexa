-- 002_performance_indexes.sql
-- Performance indexes for bulk-seeded data
-- Apply manually after bulk seed:
-- docker compose exec -T postgres psql -U postgres -d nexa < backend/scripts/002_performance_indexes.sql

-- Enable pg_trgm extension for trigram matching on text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Trigram index for patient name search (ILIKE '%query%' optimization)
-- pasienRepository.Search() runs ILIKE on nama for every list request
CREATE INDEX IF NOT EXISTS idx_trgm_pasien_nama ON pasien USING gin (nama gin_trgm_ops);

-- Foreign key indexes for join performance
CREATE INDEX IF NOT EXISTS idx_rekam_medis_pasien ON rekam_medis(id_pasien);
CREATE INDEX IF NOT EXISTS idx_rekam_medis_dokter ON rekam_medis(id_dokter);
CREATE INDEX IF NOT EXISTS idx_rekam_diagnosis_rekam ON rekam_diagnosis(id_rekam);
CREATE INDEX IF NOT EXISTS idx_detail_tindakan_rekam ON detail_tindakan(id_rekam);
CREATE INDEX IF NOT EXISTS idx_klaim_rekam ON klaim(id_rekam);
CREATE INDEX IF NOT EXISTS idx_klaim_petugas ON klaim(id_petugas_casemix);

-- Filter indexes for common queries
CREATE INDEX IF NOT EXISTS idx_klaim_status ON klaim(status_klaim) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_rekam_medis_jenis ON rekam_medis(jenis_perawatan) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_rekam_medis_tanggal ON rekam_medis(tanggal_kunjungan) WHERE deleted_at IS NULL;

-- Composite index for dashboard stats query (top_diagnosis)
CREATE INDEX IF NOT EXISTS idx_rekam_diagnosis_kode ON rekam_diagnosis(kode_icd10);
