package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"nexa/backend/internal/config"
	"nexa/backend/internal/database"
	"nexa/backend/internal/model"
	"nexa/backend/internal/models"
	"nexa/backend/pkg"

	"gorm.io/gorm"
)

// ptrString returns a pointer to the string passed in
func ptrString(s string) *string {
	return &s
}

func main() {
	// Parse flags for bulk seeder
	bulkFlag := flag.Bool("bulk", false, "Run bulk seeder to generate ~50k patients with derived data")
	pasienCount := flag.Int("pasien", 49000, "Approximate patient count for bulk seed (randomized ±4%, range 47k-51k)")
	randSeed := flag.Int64("rand-seed", 42, "RNG seed for reproducible bulk seed runs")
	flag.Parse()

	// Load config
	cfg := config.Load()
	if cfg == nil {
		log.Fatal("Failed to load config")
	}

	// Connect to database
	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connected successfully")

	// Dispatch: bulk seeder or default seeder
	if *bulkFlag {
		log.Println("Starting BULK SEEDER...")
		bulkCfg := newBulkConfig(*pasienCount, *randSeed)
		if err := RunBulkSeed(db, bulkCfg); err != nil {
			log.Fatalf("Bulk seed failed: %v", err)
		}
		// Seed expanded master data after bulk seed
		seedExpandedMasterData(db)

		// Seed audit trail (100 aktivitas acak)
		seedAuditTrail(db)

		log.Println("Bulk seeding completed successfully")
	} else {
		// Default path: existing behavior unchanged
		// Seed roles and units (if not exist)
		seedRolesAndUnits(db)

		// Seed admin user (default admin account)
		seedAdminUser(db)

		// Seed regular users
		seedRegularUsers(db)

		// Seed master reference data
		seedMasterData(db)

		log.Println("Seeding completed successfully")
	}
}

// seedRolesAndUnits creates default roles and units if they don't exist
func seedRolesAndUnits(db *gorm.DB) {
	// Roles from PRD
	roles := []string{
		"admin_ti",        // Admin TI
		"petugas_rm",      // Petugas Registrasi & RM
		"dokter_dpjp",     // Dokter Penanggung Jawab
		"perawat",         // Perawat Bangsal / Poli
		"petugas_casemix", // Koder & Verifikator Klaim
		"keuangan",        // Staf Kasir / Keuangan RS
		"manajemen",       // Direksi & Komite Medis
	}

	for _, roleName := range roles {
		var role models.Role
		result := db.Where(models.Role{NamaRole: roleName}).FirstOrCreate(&role, models.Role{NamaRole: roleName})
		if result.Error != nil {
			log.Printf("Failed to create role %s: %v", roleName, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Created role: %s", roleName)
		}
	}

	// Default unit
	var unit models.Unit
	result := db.Where(models.Unit{NamaUnit: "Unit Umum"}).FirstOrCreate(&unit, models.Unit{NamaUnit: "Unit Umum"})
	if result.Error != nil {
		log.Printf("Failed to create unit: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Created unit: %s", unit.NamaUnit)
	}
}

// seedAdminUser creates the default admin account if it doesn't exist
func seedAdminUser(db *gorm.DB) {
	var existingUser models.User
	if err := db.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		log.Println("Admin user already exists, skipping")
		return
	}

	// Get admin_ti role
	var adminRole models.Role
	if err := db.Where("nama_role = ?", "admin_ti").First(&adminRole).Error; err != nil {
		log.Printf("Warning: Could not find admin_ti role: %v", err)
		return
	}

	// Get default unit
	var unit models.Unit
	if err := db.Where("nama_unit = ?", "Unit Umum").First(&unit).Error; err != nil {
		log.Printf("Warning: Could not find Unit Umum: %v", err)
		return
	}

	passwordHash, err := pkg.HashPassword("password123")
	if err != nil {
		log.Printf("Failed to hash password for admin: %v", err)
		return
	}

	admin := models.User{
		Nama:         "Admin TI",
		Profesi:      "IT",
		Spesialisasi: ptrString("Sistem Informasi"),
		Username:     "admin",
		PasswordHash: passwordHash,
		IDRole:       adminRole.IDRole,
		IDUnit:       unit.IDUnit,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return
	}
	log.Println("Created admin user: admin (Role: admin_ti)")
}

// seedRegularUsers membuat 25 user non-admin dengan role acak
func seedRegularUsers(db *gorm.DB) {
	nonAdminRoles := []string{
		"petugas_rm", "dokter_dpjp", "perawat",
		"petugas_casemix", "keuangan", "manajemen",
	}

	// Map role name -> ID
	roleMap := make(map[string]uint)
	for _, roleName := range nonAdminRoles {
		var role models.Role
		if err := db.Where("nama_role = ?", roleName).First(&role).Error; err != nil {
			log.Printf("Warning: role %s not found: %v", roleName, err)
			continue
		}
		roleMap[roleName] = role.IDRole
	}

	var unit models.Unit
	if err := db.Where("nama_unit = ?", "Unit Umum").First(&unit).Error; err != nil {
		log.Fatalf("Failed to get unit: %v", err)
	}

	firstNames := []string{
		"Ahmad", "Budi", "Citra", "Dewi", "Eko", "Fani", "Gita", "Hadi",
		"Indri", "Joko", "Kartika", "Lilis", "Maya", "Nina", "Oki", "Putri",
		"Qori", "Rani", "Sari", "Tina", "Uli", "Vika", "Wati", "Xenia",
		"Yoga", "Zara",
	}
	lastNames := []string{
		"Santoso", "Wijaya", "Nugroho", "Prasetya", "Hidayat", "Firmansyah",
		"Kusuma", "Lestari", "Maharani", "Pramudya", "Saputra", "Tambunan",
		"Utami", "Wahyudi", "Yusuf", "Zulkifli", "Adi", "Bahri", "Cahyo",
		"Darjo", "Effendi", "Fauzi", "Gunawan", "Hartono", "Irawan",
		"Julianto", "Kurniawan", "Lubis", "Mardi", "Nasution", "Oktaviani",
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 25; i++ {
		// Role acak untuk tiap user
		roleName := nonAdminRoles[rng.Intn(len(nonAdminRoles))]
		roleID, ok := roleMap[roleName]
		if !ok {
			continue
		}

		firstName := firstNames[rng.Intn(len(firstNames))]
		lastName := lastNames[rng.Intn(len(lastNames))]
		username := fmt.Sprintf("%s%s%02d", firstName, lastName, i)

		var existingUser models.User
		if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
			continue
		}

		passwordHash, err := pkg.HashPassword("password123")
		if err != nil {
			continue
		}

		user := models.User{
			Nama:         firstName + " " + lastName,
			Profesi:      getProfesiForRole(roleName),
			Spesialisasi: ptrString(getSpesialisasiForRole(roleName)),
			Username:     username,
			PasswordHash: passwordHash,
			IDRole:       roleID,
			IDUnit:       unit.IDUnit,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", username, err)
			continue
		}
		log.Printf("Created user: %s (%s) - Role: %s", username, user.Nama, roleName)
	}
}

// seedAuditTrail membuat 100 log aktivitas acak
func seedAuditTrail(db *gorm.DB) {
	log.Println("Seeding audit trail (100 entries)...")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Ambil semua user ID untuk referensi
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Warning: failed to load users for audit trail: %v", err)
		return
	}

	if len(users) == 0 {
		log.Println("Warning: no users found, skipping audit trail")
		return
	}

	tables := []string{"pasien", "rekam_medis", "rekam_diagnosis", "klaim", "users", "diagnosis", "tindakan"}
	activities := []string{"INSERT", "UPDATE", "DELETE"}

	for i := 0; i < 100; i++ {
		u := users[rng.Intn(len(users))]
		tableName := tables[rng.Intn(len(tables))]
		activity := activities[rng.Intn(len(activities))]

		hoursAgo := rng.Intn(720)
		waktu := time.Now().Add(-time.Duration(hoursAgo) * time.Hour)

		logEntry := model.LogAktivitas{
			IDUser:         &u.IDUser,
			TabelTerdampak: tableName,
			Aktivitas:      activity,
			DataID:         nil,
			Waktu:          waktu,
		}

		if err := db.Create(&logEntry).Error; err != nil {
			log.Printf("Warning: failed to create audit log %d: %v", i, err)
		}
	}

	log.Println("Audit trail seeding completed (100 entries)")
}

// getProfesiForRole returns appropriate profession for a role
func getProfesiForRole(roleName string) string {
	switch roleName {
	case "petugas_rm":
		return "Rekam Medis"
	case "dokter_dpjp":
		return "Dokter"
	case "perawat":
		return "Perawat"
	case "petugas_casemix":
		return "Koder Casemix"
	case "keuangan":
		return "Keuangan"
	case "manajemen":
		return "Manajemen"
	default:
		return "Staff"
	}
}

// getSpesialisasiForRole returns appropriate specialization for a role
func getSpesialisasiForRole(roleName string) string {
	switch roleName {
	case "petugas_rm":
		return "Registrasi"
	case "dokter_dpjp":
		return "Umum"
	case "perawat":
		return "Umum"
	case "petugas_casemix":
		return "Casemix"
	case "keuangan":
		return "Keuangan Rumah Sakit"
	case "manajemen":
		return "Manajemen Rumah Sakit"
	default:
		return ""
	}
}

// seedMasterData seeds ICD-10, ICD-9 CM, and INA-CBGs reference data
func seedMasterData(db *gorm.DB) {
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
	for _, d := range icd10Data {
		result := db.Where("kode_icd10 = ?", d.KodeICD10).FirstOrCreate(&d)
		if result.Error != nil {
			log.Printf("Failed to seed ICD-10 %s: %v", d.KodeICD10, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Seeded ICD-10: %s - %s", d.KodeICD10, d.NamaDiagnosis)
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
		result := db.Where("kode_tindakan = ?", t.KodeTindakan).FirstOrCreate(&t)
		if result.Error != nil {
			log.Printf("Failed to seed ICD-9 %s: %v", t.KodeTindakan, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Seeded ICD-9: %s - %s", t.KodeTindakan, t.NamaTindakan)
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
		result := db.Where("kode_cbgs = ?", c.KodeCBGS).FirstOrCreate(&c)
		if result.Error != nil {
			log.Printf("Failed to seed CBGs %s: %v", c.KodeCBGS, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Seeded CBGs: %s - %s", c.KodeCBGS, c.Deskripsi)
		}
	}

	log.Println("Master data seeding completed")
}

// seedExpandedMasterData seeds the expanded master data sets for bulk seeding.
// Called after bulk seed to populate diagnosis, tindakan, and tarif_cbgs with
// 150-300 ICD-10, 80-150 ICD-9, and 60-120 CBGs codes respectively.
func seedExpandedMasterData(db *gorm.DB) {
	log.Println("Seeding expanded master data...")

	// Seed expanded ICD-10 data
	for _, d := range expandedICD10Data {
		result := db.Where("kode_icd10 = ?", d.KodeICD10).FirstOrCreate(&d)
		if result.Error != nil {
			log.Printf("[MASTER] gagal seed ICD-10 %s: %v", d.KodeICD10, result.Error)
		}
	}

	// Seed expanded ICD-9 data
	for _, t := range expandedICD9Data {
		result := db.Where("kode_tindakan = ?", t.KodeTindakan).FirstOrCreate(&t)
		if result.Error != nil {
			log.Printf("[MASTER] gagal seed ICD-9 %s: %v", t.KodeTindakan, result.Error)
		}
	}

	// Seed expanded CBGs data
	for _, c := range expandedCBGsData {
		result := db.Where("kode_cbgs = ?", c.KodeCBGS).FirstOrCreate(&c)
		if result.Error != nil {
			log.Printf("[MASTER] gagal seed CBGs %s: %v", c.KodeCBGS, result.Error)
		}
	}

	var icd10Count, icd9Count, cbgsCount int64
	db.Model(&model.Diagnosis{}).Count(&icd10Count)
	db.Model(&model.Tindakan{}).Count(&icd9Count)
	db.Model(&model.TarifCBGs{}).Count(&cbgsCount)

	log.Printf("Expanded master data seeding completed: %d ICD-10, %d ICD-9, %d CBGs", icd10Count, icd9Count, cbgsCount)
}
