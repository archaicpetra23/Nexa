package database

import (
	"log"
	"nexa/backend/internal/model"
	"nexa/backend/internal/models"
	"nexa/backend/pkg"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Role{}, &models.Unit{}, &models.User{}, &models.Session{},
		&model.Pasien{}, &model.RekamMedis{}, &model.Diagnosis{}, &model.RekamDiagnosis{},
		&model.Tindakan{}, &model.DetailTindakan{}, &model.TarifCBGs{},
		&model.Klaim{}, &model.LogAktivitas{},
	); err != nil {
		return nil, err
	}

	seedData(db)
	return db, nil
}

func seedData(db *gorm.DB) {
	hash, _ := pkg.HashPassword("password123")

	roles := []string{"admin_ti", "petugas_rm", "dokter_dpjp", "perawat", "petugas_casemix", "keuangan", "manajemen"}
	for _, r := range roles {
		db.Where(models.Role{NamaRole: r}).FirstOrCreate(&models.Role{NamaRole: r})
	}

	db.Where(models.Unit{NamaUnit: "Unit Umum"}).FirstOrCreate(&models.Unit{NamaUnit: "Unit Umum"})

	userMap := map[string]uint{
		"admin":     1,
		"petugas":   2,
		"dokter":    3,
		"perawat":   4,
		"casemix":   5,
		"keuangan":  6,
		"manajemen": 7,
	}

	userNames := []string{"admin", "petugas", "dokter", "perawat", "casemix", "keuangan", "manajemen"}
	displayNames := []string{"Admin TI", "Petugas RM", "Dr. Dokter", "Perawat", "Petugas Casemix", "Staf Keuangan", "Direksi"}
	professions := []string{"IT", "Rekam Medis", "Dokter", "Perawat", "Koder", "Keuangan", "Manajemen"}

	for i, un := range userNames {
		db.Where(models.User{Username: un}).FirstOrCreate(&models.User{
			Username:     un,
			Nama:         displayNames[i],
			Profesi:      professions[i],
			PasswordHash: hash,
			IDRole:       userMap[un],
			IDUnit:       1,
		})
	}

	// Seed ICD-10 Diagnosis
	icd10Data := []model.Diagnosis{
		{KodeICD10: "A09", NamaDiagnosis: "Diare dan Gastroenteritis"},
		{KodeICD10: "B34", NamaDiagnosis: "Infeksi Virus"},
		{KodeICD10: "D64", NamaDiagnosis: "Anemia"},
		{KodeICD10: "E11", NamaDiagnosis: "Diabetes Melitus Tipe 2"},
		{KodeICD10: "I10", NamaDiagnosis: "Hipertensi"},
		{KodeICD10: "J06", NamaDiagnosis: "Infeksi Saluran Pernapasan Atas"},
		{KodeICD10: "J18", NamaDiagnosis: "Pneumonia"},
		{KodeICD10: "K35", NamaDiagnosis: "Apendisitis"},
		{KodeICD10: "N39", NamaDiagnosis: "Infeksi Saluran Kemih"},
		{KodeICD10: "R50", NamaDiagnosis: "Demam"},
	}
	// Seed reference data. Errors are logged, never swallowed — a silent seed failure
	// leaves the ICD autocomplete and the dashboard empty without any visible symptom.
	for _, d := range icd10Data {
		if err := db.Where("kode_icd10 = ?", d.KodeICD10).
			FirstOrCreate(&d, model.Diagnosis{
				KodeICD10:     d.KodeICD10,
				NamaDiagnosis: d.NamaDiagnosis,
				Kategori:      d.Kategori,
			}).Error; err != nil {
			log.Printf("[SEED] gagal seed ICD-10 %s: %v", d.KodeICD10, err)
		}
	}

	// Seed ICD-9 CM Tindakan
	icd9Data := []model.Tindakan{
		{KodeTindakan: "38.93", NamaTindakan: "Pengambilan Darah Vena", TarifStandar: 50000},
		{KodeTindakan: "45.13", NamaTindakan: "Endoskopi Lambung", TarifStandar: 1500000},
		{KodeTindakan: "87.44", NamaTindakan: "Foto Rontgen Dada", TarifStandar: 150000},
		{KodeTindakan: "88.01", NamaTindakan: "CT Scan Kepala", TarifStandar: 900000},
		{KodeTindakan: "88.76", NamaTindakan: "USG Abdomen", TarifStandar: 250000},
		{KodeTindakan: "93.90", NamaTindakan: "Fisioterapi", TarifStandar: 100000},
		{KodeTindakan: "96.04", NamaTindakan: "Pemasangan Infus", TarifStandar: 75000},
		{KodeTindakan: "96.71", NamaTindakan: "Ventilator Mekanik", TarifStandar: 2000000},
		{KodeTindakan: "99.04", NamaTindakan: "Suntikan Antibiotik", TarifStandar: 120000},
		{KodeTindakan: "99.15", NamaTindakan: "Transfusi Darah", TarifStandar: 450000},
	}
	for _, t := range icd9Data {
		if err := db.Where("kode_tindakan = ?", t.KodeTindakan).
			FirstOrCreate(&t, model.Tindakan{
				KodeTindakan: t.KodeTindakan,
				NamaTindakan: t.NamaTindakan,
				TarifStandar: t.TarifStandar,
			}).Error; err != nil {
			log.Printf("[SEED] gagal seed ICD-9 CM %s: %v", t.KodeTindakan, err)
		}
	}

	// Seed INA-CBGs Tarif
	cbgsData := []model.TarifCBGs{
		{KodeCBGS: "A-4-10-I", Deskripsi: "Kasus Infeksi Ringan", Tarif: 850000},
		{KodeCBGS: "B-1-14-II", Deskripsi: "Penyakit Jantung Tingkat Sedang", Tarif: 4200000},
		{KodeCBGS: "C-4-13-III", Deskripsi: "Diabetes Komplikasi Berat", Tarif: 6500000},
		{KodeCBGS: "D-4-16-I", Deskripsi: "Pneumonia Ringan", Tarif: 1750000},
		{KodeCBGS: "E-4-10-II", Deskripsi: "Hipertensi dengan Komplikasi", Tarif: 2100000},
	}
	for _, c := range cbgsData {
		if err := db.Where("kode_cbgs = ?", c.KodeCBGS).
			FirstOrCreate(&c, model.TarifCBGs{
				KodeCBGS:  c.KodeCBGS,
				Deskripsi: c.Deskripsi,
				Tarif:     c.Tarif,
			}).Error; err != nil {
			log.Printf("[SEED] gagal seed INA-CBGs %s: %v", c.KodeCBGS, err)
		}
	}

	log.Println("Database seeded successfully")
}
