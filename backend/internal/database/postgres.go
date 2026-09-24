package database

import (
	"log"
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

	if err := db.AutoMigrate(&models.Role{}, &models.Unit{}, &models.User{}, &models.Session{}); err != nil {
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

	log.Println("Database seeded successfully")
}
