package model

import "nexa/backend/internal/models"

import "time"

type Pasien struct {
	IDPasien     uint       `gorm:"primaryKey;column:id_pasien" json:"id_pasien"`
	NIK          string     `gorm:"column:nik;uniqueIndex;not null" json:"nik"`
	NoBPJS       *string    `gorm:"column:no_bpjs;uniqueIndex" json:"no_bpjs"`
	Nama         string     `gorm:"not null" json:"nama"`
	TanggalLahir time.Time  `gorm:"column:tanggal_lahir;not null" json:"tanggal_lahir"`
	JenisKelamin string     `gorm:"column:jenis_kelamin;not null;check:jenis_kelamin IN ('L','P')" json:"jenis_kelamin"`
	Alamat       string     `gorm:"not null" json:"alamat"`
	NoHP         *string    `gorm:"column:no_hp" json:"no_hp"`
	CreatedAt    time.Time  `gorm:"default:now();column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:now();column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

func (Pasien) TableName() string {
	return "pasien"
}

type RekamMedis struct {
	IDRekam          uint       `gorm:"primaryKey;column:id_rekam" json:"id_rekam"`
	IDPasien         uint       `gorm:"column:id_pasien;not null" json:"id_pasien"`
	IDDokter         uint       `gorm:"column:id_dokter;not null" json:"id_dokter"`
	TanggalKunjungan time.Time  `gorm:"column:tanggal_kunjungan;not null;default:now()" json:"tanggal_kunjungan"`
	TanggalPulang    *time.Time `gorm:"column:tanggal_pulang" json:"tanggal_pulang"`
	JenisPerawatan   string     `gorm:"column:jenis_perawatan;not null;check:jenis_perawatan IN ('rawat_jalan','rawat_inap')" json:"jenis_perawatan"`
	Keluhan          string     `gorm:"not null" json:"keluhan"`
	Catatan          *string    `json:"catatan"`
	CreatedAt        time.Time  `gorm:"default:now();column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"default:now();column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	Pasien *Pasien      `gorm:"foreignKey:IDPasien" json:"pasien,omitempty"`
	Dokter *models.User `gorm:"foreignKey:IDDokter" json:"dokter,omitempty"`
}

func (RekamMedis) TableName() string {
	return "rekam_medis"
}

type Diagnosis struct {
	KodeICD10     string  `gorm:"primaryKey;column:kode_icd10;type:text" json:"kode_icd10"`
	NamaDiagnosis string  `gorm:"column:nama_diagnosis;not null" json:"nama_diagnosis"`
	Kategori      *string `json:"kategori"`
}

func (Diagnosis) TableName() string {
	return "diagnosis"
}

type RekamDiagnosis struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	IDRekam   uint   `gorm:"column:id_rekam;not null" json:"id_rekam"`
	KodeICD10 string `gorm:"column:kode_icd10;type:text;not null" json:"kode_icd10"`
	Jenis     string `gorm:"not null;check:jenis IN ('primer','sekunder')" json:"jenis"`

	RekamMedis *RekamMedis `gorm:"foreignKey:IDRekam" json:"rekam_medis,omitempty"`
	// references: MUST be explicit. With only foreignKey:, GORM infers the
	// referenced primary key and rewrites diagnosis.kode_icd10 to bigint, which
	// breaks the whole ICD seed with "invalid input syntax for type bigint".
	Diagnosis *Diagnosis `gorm:"foreignKey:KodeICD10;references:KodeICD10" json:"diagnosis,omitempty"`
}

func (RekamDiagnosis) TableName() string {
	return "rekam_diagnosis"
}

type Tindakan struct {
	KodeTindakan string  `gorm:"primaryKey;column:kode_tindakan;type:text" json:"kode_tindakan"`
	NamaTindakan string  `gorm:"column:nama_tindakan;not null" json:"nama_tindakan"`
	TarifStandar float64 `gorm:"column:tarif_standar;not null;default:0;check:tarif_standar >= 0" json:"tarif_standar"`
}

func (Tindakan) TableName() string {
	return "tindakan"
}

type DetailTindakan struct {
	IDDetail     uint   `gorm:"primaryKey;column:id_detail" json:"id_detail"`
	IDRekam      uint   `gorm:"column:id_rekam;not null" json:"id_rekam"`
	KodeTindakan string `gorm:"column:kode_tindakan;type:text;not null" json:"kode_tindakan"`
	Jumlah       int    `gorm:"not null;default:1;check:jumlah > 0" json:"jumlah"`

	RekamMedis *RekamMedis `gorm:"foreignKey:IDRekam" json:"rekam_medis,omitempty"`
	// references: MUST be explicit — same reason as RekamDiagnosis.Diagnosis.
	Tindakan *Tindakan `gorm:"foreignKey:KodeTindakan;references:KodeTindakan" json:"tindakan,omitempty"`
}

func (DetailTindakan) TableName() string {
	return "detail_tindakan"
}

type TarifCBGs struct {
	KodeCBGS  string  `gorm:"primaryKey;column:kode_cbgs;type:text" json:"kode_cbgs"`
	Deskripsi string  `gorm:"not null" json:"deskripsi"`
	Tarif     float64 `gorm:"not null;check:tarif >= 0" json:"tarif"`
}

func (TarifCBGs) TableName() string {
	return "tarif_cbgs"
}

type Klaim struct {
	IDKlaim            uint       `gorm:"primaryKey;column:id_klaim" json:"id_klaim"`
	IDRekam            uint       `gorm:"column:id_rekam;unique;not null" json:"id_rekam"`
	KodeCBGS           string     `gorm:"column:kode_cbgs;not null" json:"kode_cbgs"`
	IDPetugasCasemix   uint       `gorm:"column:id_petugas_casemix;not null" json:"id_petugas_casemix"`
	StatusKlaim        string     `gorm:"column:status_klaim;not null;check:status_klaim IN ('draft','pending','disetujui','ditolak')" json:"status_klaim"`
	TanggalKlaim       time.Time  `gorm:"column:tanggal_klaim;not null;default:CURRENT_DATE" json:"tanggal_klaim"`
	NominalKlaim       float64    `gorm:"column:nominal_klaim;not null;check:nominal_klaim >= 0" json:"nominal_klaim"`
	AlasanPendingTolak *string    `gorm:"column:alasan_pending_tolak" json:"alasan_pending_tolak"`
	CreatedAt          time.Time  `gorm:"default:now();column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"default:now();column:updated_at" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	RekamMedis         *RekamMedis `gorm:"foreignKey:IDRekam" json:"rekam_medis,omitempty"`
	// references: MUST be explicit — same reason as RekamDiagnosis.Diagnosis.
	KodeCBGSNavigation *TarifCBGs   `gorm:"foreignKey:KodeCBGS;references:KodeCBGS" json:"tarif_cbgs,omitempty"`
	PetugasCasemix     *models.User `gorm:"foreignKey:IDPetugasCasemix" json:"petugas_casemix,omitempty"`
}

func (Klaim) TableName() string {
	return "klaim"
}

type LogAktivitas struct {
	IDLog          uint      `gorm:"primaryKey;column:id_log" json:"id_log"`
	IDUser         *uint     `gorm:"column:id_user" json:"id_user"`
	TabelTerdampak string    `gorm:"column:tabel_terdampak;not null" json:"tabel_terdampak"`
	Aktivitas      string    `gorm:"not null;check:aktivitas IN ('INSERT','UPDATE','DELETE','LOGIN','LOGOUT')" json:"aktivitas"`
	DataID         *int      `gorm:"column:data_id" json:"data_id"`
	Waktu          time.Time `gorm:"default:now()" json:"waktu"`

	User *models.User `gorm:"foreignKey:IDUser" json:"user,omitempty"`
}

func (LogAktivitas) TableName() string {
	return "log_aktivitas"
}

// Session already exists in session.go
// Role, User, Unit already exist in role.go
