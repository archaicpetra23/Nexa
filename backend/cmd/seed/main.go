package main

import (
	"log"
	"math/rand"
	"time"

	"nexa/backend/internal/config"
	"nexa/backend/internal/database"
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

	log.Println("Seeding completed successfully")
}

// seedRolesAndUnits creates default roles and units if they don't exist
func seedRolesAndUnits(db *gorm.DB) {
	// Roles from PRD
	roles := []string{
		"admin_ti",      // Admin TI
		"petugas_rm",    // Petugas Registrasi & RM
		"dokter_dpjp",   // Dokter Penanggung Jawab
		"perawat",       // Perawat Bangsal / Poli
		"petugas_casemix", // Koder & Verifikator Klaim
		"keuangan",      // Staf Kasir / Keuangan RS
		"manajemen",     // Direksi & Komite Medis
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
		"petugas_rm",    // Petugas Registrasi & RM
		"dokter_dpjp",   // Dokter Penanggung Jawab
		"perawat",       // Perawat Bangsal / Poli
		"petugas_casemix", // Koder & Verifikator Klaim
		"keuangan",      // Staf Kasir / Keuangan RS
		"manajemen",     // Direksi & Komite Medis
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