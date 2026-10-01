package main

import (
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

	// Seed roles and units (if not exist)
	seedRolesAndUnits(db)

	// Seed regular users only (no admin - will be created manually)
	seedRegularUsers(db)

	// Seed master reference data
	seedMasterData(db)

	log.Println("Seeding completed successfully")
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

// seedRegularUsers creates users for non-admin roles only with random names
func seedRegularUsers(db *gorm.DB) {
	// Define non-admin roles (all except admin_ti)
	nonAdminRoles := []string{
		"petugas_rm",      // Petugas Registrasi & RM
		"dokter_dpjp",     // Dokter Penanggung Jawab
		"perawat",         // Perawat Bangsal / Poli
		"petugas_casemix", // Koder & Verifikator Klaim
		"keuangan",        // Staf Kasir / Keuangan RS
		"manajemen",       // Direksi & Komite Medis
	}

	// Get role IDs for mapping
	var roleMap = make(map[string]uint)
	for _, roleName := range nonAdminRoles {
		var role models.Role
		if err := db.Where("nama_role = ?", roleName).First(&role); err != nil {
			log.Printf("Warning: Could not find role %s: %v", roleName, err)
			continue
		}
		roleMap[roleName] = role.IDRole
	}

	// Get unit ID (assuming "Unit Umum" exists)
	var unit models.Unit
	if err := db.Where("nama_unit = ?", "Unit Umum").First(&unit).Error; err != nil {
		log.Fatalf("Failed to get unit 'Unit Umum': %v", err)
	}

	// First name and last name lists for random generation
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

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create one user per non-admin role with random name
	for _, roleName := range nonAdminRoles {
		roleID, ok := roleMap[roleName]
		if !ok {
			log.Printf("Skipping user creation for role %s (role not found)", roleName)
			continue
		}

		// Generate random name
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		fullName := firstName + " " + lastName

		// Generate username (lowercase, no spaces)
		username := firstName + lastName
		// Ensure uniqueness by checking if exists and adding number if needed
		var count int64
		db.Model(&models.User{}).Where("username = ?", username).Count(&count)
		if count > 0 {
			username = username + "01" // Add suffix if duplicate
		}

		// Hash password (default password as per existing system)
		passwordHash, err := pkg.HashPassword("password123")
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", username, err)
			continue
		}

		// Check if user with this username already exists
		var existingUser models.User
		if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
			log.Printf("User %s already exists, skipping", username)
			continue
		}

		// Create the user
		user := models.User{
			Nama:         fullName,
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

		log.Printf("Created regular user: %s (%s) - Role: %s", username, fullName, roleName)
	}
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
