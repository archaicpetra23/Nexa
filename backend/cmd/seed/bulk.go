package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/models"

	"gorm.io/gorm"
)

// BulkConfig holds all tunable parameters for the bulk seeder.
// All counts are approximate and will be randomized within specified bands.
type BulkConfig struct {
	TargetPatients int // Approximate patient count (randomized ±10%)
	RandSeed       int64

	// Derived count bands (relative to actual randomized patient count)
	RekamMedisBandMin float64 // 0.60
	RekamMedisBandMax float64 // 0.90
	DiagnosisBandMin  float64 // 1.00 per rekam
	DiagnosisBandMax  float64 // 1.40 per rekam
	TindakanBandMin   float64 // 1.00 per rekam
	TindakanBandMax   float64 // 1.40 per rekam
	KlaimBandMin      float64 // 0.40 per rekam
	KlaimBandMax      float64 // 0.70 per rekam

	// Batch size for inserts
	BatchSize int
}

// newBulkConfig returns a config with default values
func newBulkConfig(targetPatients int, randSeed int64) *BulkConfig {
	return &BulkConfig{
		TargetPatients:    targetPatients,
		RandSeed:          randSeed,
		RekamMedisBandMin: 0.60,
		RekamMedisBandMax: 0.90,
		DiagnosisBandMin:  1.00,
		DiagnosisBandMax:  1.40,
		TindakanBandMin:   1.00,
		TindakanBandMax:   1.40,
		KlaimBandMin:      0.40,
		KlaimBandMax:      0.70,
		BatchSize:         500,
	}
}

// randomizeCounts computes actual counts with randomization.
// The patient count must NEVER be exactly the target.
func (c *BulkConfig) randomizeCounts(rng *rand.Rand) (int, int, int, int, int) {
	// Randomize patient count ±10%, ensure it's never exactly the target
	patientCount := c.TargetPatients + rng.Intn(c.TargetPatients/5+1) - c.TargetPatients/10
	if patientCount == c.TargetPatients {
		if rng.Float64() < 0.5 {
			patientCount++
		} else {
			patientCount--
		}
	}

	// Derive each count independently with randomization
	rekamCount := int(float64(patientCount) * (c.RekamMedisBandMin + rng.Float64()*(c.RekamMedisBandMax-c.RekamMedisBandMin)))
	diagnosisCount := int(float64(rekamCount) * (c.DiagnosisBandMin + rng.Float64()*(c.DiagnosisBandMax-c.DiagnosisBandMin)))
	tindakanCount := int(float64(rekamCount) * (c.TindakanBandMin + rng.Float64()*(c.TindakanBandMax-c.TindakanBandMin)))
	klaimCount := int(float64(rekamCount) * (c.KlaimBandMin + rng.Float64()*(c.KlaimBandMax-c.KlaimBandMin)))

	return patientCount, rekamCount, diagnosisCount, tindakanCount, klaimCount
}

// RunBulkSeed executes the bulk data generation with idempotency check
func RunBulkSeed(db *gorm.DB, cfg *BulkConfig) error {
	// Idempotency: check if patients already exist
	var existingCount int64
	if err := db.Model(&model.Pasien{}).Count(&existingCount).Error; err != nil {
		return fmt.Errorf("failed to check existing patients: %w", err)
	}

	if existingCount > 0 {
		log.Printf("Bulk seed skipped: %d patients already exist", existingCount)
		log.Println("To reseed, wipe the volume: docker compose down -v && docker compose up -d --build")
		return nil
	}

	rng := rand.New(rand.NewSource(cfg.RandSeed))

	// Compute randomized counts
	patientCount, rekamCount, diagnosisCount, tindakanCount, klaimCount := cfg.randomizeCounts(rng)

	log.Println("=== BULK SEED PLAN ===")
	log.Printf("Target patients (--pasien): %d", cfg.TargetPatients)
	log.Printf("Randomized patient count: %d", patientCount)
	log.Printf("Planned rekam_medis: %d", rekamCount)
	log.Printf("Planned rekam_diagnosis: %d", diagnosisCount)
	log.Printf("Planned detail_tindakan: %d", tindakanCount)
	log.Printf("Planned klaim: %d", klaimCount)
	log.Println("======================")

	// Assert hard invariants before proceeding
	if klaimCount > rekamCount {
		return fmt.Errorf("invariant violated: klaimCount (%d) > rekamCount (%d)", klaimCount, rekamCount)
	}

	startTime := time.Now()

	// Generate and insert patients
	log.Println("Generating patients...")
	patients, err := generatePatients(patientCount, rng)
	if err != nil {
		return fmt.Errorf("failed to generate patients: %w", err)
	}
	if err := insertBatch(db, patients, cfg.BatchSize, "pasien"); err != nil {
		return err
	}
	firstPatientID := patients[0].IDPasien
	log.Printf("✓ pasien: %d rows inserted (IDs: %d-%d)", len(patients), firstPatientID, patients[len(patients)-1].IDPasien)

	// Query dokter_dpjp and petugas_casemix users once
	dokterIDs, err := getUserIDsByRole(db, "dokter_dpjp")
	if err != nil || len(dokterIDs) == 0 {
		return fmt.Errorf("no dokter_dpjp users found: %w", err)
	}
	casemixIDs, err := getUserIDsByRole(db, "petugas_casemix")
	if err != nil || len(casemixIDs) == 0 {
		return fmt.Errorf("no petugas_casemix users found: %w", err)
	}

	// Generate and insert rekam_medis
	log.Println("Generating rekam_medis...")
	rekams, err := generateRekamMedis(rekamCount, firstPatientID, patientCount, dokterIDs, rng)
	if err != nil {
		return fmt.Errorf("failed to generate rekam_medis: %w", err)
	}
	if err := insertBatch(db, rekams, cfg.BatchSize, "rekam_medis"); err != nil {
		return err
	}
	firstRekamID := rekams[0].IDRekam
	log.Printf("✓ rekam_medis: %d rows inserted (IDs: %d-%d)", len(rekams), firstRekamID, rekams[len(rekams)-1].IDRekam)

	// Generate and insert rekam_diagnosis
	log.Println("Generating rekam_diagnosis...")
	diagnoses, err := generateRekamDiagnosis(diagnosisCount, rekams, rng)
	if err != nil {
		return fmt.Errorf("failed to generate rekam_diagnosis: %w", err)
	}
	if err := insertBatch(db, diagnoses, cfg.BatchSize, "rekam_diagnosis"); err != nil {
		return err
	}
	log.Printf("✓ rekam_diagnosis: %d rows inserted", len(diagnoses))

	// Generate and insert detail_tindakan
	log.Println("Generating detail_tindakan...")
	tindakans, err := generateDetailTindakan(tindakanCount, rekams, rng)
	if err != nil {
		return fmt.Errorf("failed to generate detail_tindakan: %w", err)
	}
	if err := insertBatch(db, tindakans, cfg.BatchSize, "detail_tindakan"); err != nil {
		return err
	}
	log.Printf("✓ detail_tindakan: %d rows inserted", len(tindakans))

	// Generate and insert klaim (DISTINCT rekam ids via shuffle)
	log.Println("Generating klaim...")
	klaims, err := generateKlaim(klaimCount, rekams, casemixIDs, rng)
	if err != nil {
		return fmt.Errorf("failed to generate klaim: %w", err)
	}
	if err := insertBatch(db, klaims, cfg.BatchSize, "klaim"); err != nil {
		return err
	}
	log.Printf("✓ klaim: %d rows inserted", len(klaims))

	// Write exactly ONE summary log entry for the bulk seed
	// Rationale: log_aktivitas is an audit trail of user actions. A bulk seed is not a
	// per-row user action, so we record one summary row rather than thousands of INSERT rows.
	summaryLog := model.LogAktivitas{
		IDUser:         nil, // No specific user; this is a seeder operation
		TabelTerdampak: "seed",
		Aktivitas:      "INSERT",
		DataID:         nil,
		Waktu:          time.Now(),
	}
	if err := db.Create(&summaryLog).Error; err != nil {
		log.Printf("Warning: failed to write log_aktivitas summary: %v", err)
	}

	duration := time.Since(startTime)
	log.Println("=== BULK SEED COMPLETE ===")
	log.Printf("Duration: %v", duration)
	log.Printf("Total rows: %d", patientCount+rekamCount+diagnosisCount+tindakanCount+klaimCount)
	log.Println("==========================")

	return nil
}

// generatePatients creates patient records with unique NIK by construction
func generatePatients(count int, rng *rand.Rand) ([]model.Pasien, error) {
	patients := make([]model.Pasien, count)

	// NIK uniqueness by construction: random 8-digit base + monotonic counter
	nikBase := rng.Intn(90000000) + 10000000 // 8-digit base

	for i := 0; i < count; i++ {
		// NIK: exactly 16 digits, unique by arithmetic
		nik := fmt.Sprintf("%08d%08d", nikBase, i)

		// NoBPJS: ~85% filled, 13 digits, unique when present
		var noBPJS *string
		if rng.Float64() < 0.85 {
			bpjsBase := rng.Intn(90000) + 10000 // 5-digit base
			bpjs := fmt.Sprintf("%05d%08d", bpjsBase, i)
			noBPJS = &bpjs
		}

		// Name: first + occasional middle + last
		firstName := firstNames[rng.Intn(len(firstNames))]
		lastName := lastNames[rng.Intn(len(lastNames))]
		nama := firstName + " " + lastName
		if rng.Float64() < 0.3 {
			middleName := middleNames[rng.Intn(len(middleNames))]
			nama = firstName + " " + middleName + " " + lastName
		}

		// TanggalLahir: 1940-01-01 to 2025-12-31, weighted toward ages 20-70
		// That's roughly 1954-2006 for current year 2026
		birthYear := 1940 + rng.Intn(86)
		if rng.Float64() < 0.6 { // 60% fall in the 20-70 age range
			birthYear = 1956 + rng.Intn(51) // 1956-2006
		}
		birthMonth := rng.Intn(12) + 1
		birthDay := rng.Intn(28) + 1 // Avoid month-end edge cases
		tanggalLahir := time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.UTC)

		// JenisKelamin: ~50/50
		jenisKelamin := "L"
		if rng.Float64() < 0.5 {
			jenisKelamin = "P"
		}

		// Alamat: street + number + kelurahan + kecamatan + kota
		streetPrefix := streetPrefixes[rng.Intn(len(streetPrefixes))]
		streetNumber := rng.Intn(200) + 1
		kel := kelurahan[rng.Intn(len(kelurahan))]
		kec := kecamatan[rng.Intn(len(kecamatan))]
		kotaName := kota[rng.Intn(len(kota))]
		alamat := fmt.Sprintf("%s No. %d, Kel. %s, %s, %s", streetPrefix, streetNumber, kel, kec, kotaName)

		// NoHP: ~90% filled, 08 + 9-11 digits, total ≤15 chars
		var noHP *string
		if rng.Float64() < 0.90 {
			hpDigits := 9 + rng.Intn(3) // 9, 10, or 11 digits after "08"
			hpNumber := ""
			for d := 0; d < hpDigits; d++ {
				hpNumber += fmt.Sprintf("%d", rng.Intn(10))
			}
			hp := "08" + hpNumber
			noHP = &hp
		}

		patients[i] = model.Pasien{
			NIK:          nik,
			NoBPJS:       noBPJS,
			Nama:         nama,
			TanggalLahir: tanggalLahir,
			JenisKelamin: jenisKelamin,
			Alamat:       alamat,
			NoHP:         noHP,
		}
	}

	return patients, nil
}

// generateRekamMedis creates medical records referencing generated patients
func generateRekamMedis(count int, firstPatientID uint, patientCount int, dokterIDs []uint, rng *rand.Rand) ([]model.RekamMedis, error) {
	rekams := make([]model.RekamMedis, count)

	// Date range: last 24 months
	now := time.Now()
	twoYearsAgo := now.AddDate(-2, 0, 0)

	for i := 0; i < count; i++ {
		// IDPasien: pick from generated patient ID range
		idPasien := firstPatientID + uint(rng.Intn(patientCount))

		// IDDokter: pick from dokter_dpjp pool
		idDokter := dokterIDs[rng.Intn(len(dokterIDs))]

		// TanggalKunjungan: random within last 24 months
		daysSince := rng.Intn(730) // ~2 years
		tanggalKunjungan := twoYearsAgo.AddDate(0, 0, daysSince)

		// JenisPerawatan: ~70% rawat_jalan, ~30% rawat_inap
		jenisPerawatan := "rawat_jalan"
		if rng.Float64() < 0.30 {
			jenisPerawatan = "rawat_inap"
		}

		// TanggalPulang: nil for rawat_jalan; for rawat_inap = kunjungan + 1-14 days
		var tanggalPulang *time.Time
		if jenisPerawatan == "rawat_inap" {
			lossDays := rng.Intn(14) + 1
			pulang := tanggalKunjungan.AddDate(0, 0, lossDays)
			tanggalPulang = &pulang
		}

		// Keluhan: random phrase
		keluhanText := keluhan[rng.Intn(len(keluhan))]

		// Catatan: ~60% filled
		var catatanText *string
		if rng.Float64() < 0.60 {
			cat := catatan[rng.Intn(len(catatan))]
			catatanText = &cat
		}

		rekams[i] = model.RekamMedis{
			IDPasien:         idPasien,
			IDDokter:         idDokter,
			TanggalKunjungan: tanggalKunjungan,
			TanggalPulang:    tanggalPulang,
			JenisPerawatan:   jenisPerawatan,
			Keluhan:          keluhanText,
			Catatan:          catatanText,
		}
	}

	return rekams, nil
}

// generateRekamDiagnosis creates diagnosis records: 1 primer + 0-2 sekunder per rekam.
// Enforces no duplicate (id_rekam, kode_icd10) pairs.
func generateRekamDiagnosis(targetCount int, rekams []model.RekamMedis, rng *rand.Rand) ([]model.RekamDiagnosis, error) {
	diagnoses := make([]model.RekamDiagnosis, 0, targetCount)

	// Build ICD-10 code list from expanded master data
	icd10Codes := make([]string, len(expandedICD10Data))
	for i, d := range expandedICD10Data {
		icd10Codes[i] = d.KodeICD10
	}

	if len(icd10Codes) == 0 {
		return nil, fmt.Errorf("no ICD-10 codes available")
	}

	for _, rekam := range rekams {
		// Each rekam gets exactly 1 'primer' + 0-2 'sekunder'
		numSekunder := rng.Intn(3) // 0, 1, or 2
		totalForRekam := 1 + numSekunder

		usedCodes := make(map[string]bool)

		for j := 0; j < totalForRekam; j++ {
			jenis := "primer"
			if j > 0 {
				jenis = "sekunder"
			}

			// Pick a code not yet used for this rekam
			var kodeICD10 string
			attempts := 0
			for {
				kodeICD10 = icd10Codes[rng.Intn(len(icd10Codes))]
				if !usedCodes[kodeICD10] {
					break
				}
				attempts++
				if attempts > 100 {
					// Fallback: just pick any, collision is unlikely with 150+ codes
					break
				}
			}
			usedCodes[kodeICD10] = true

			diagnoses = append(diagnoses, model.RekamDiagnosis{
				IDRekam:   rekam.IDRekam,
				KodeICD10: kodeICD10,
				Jenis:     jenis,
			})

			if len(diagnoses) >= targetCount {
				return diagnoses, nil
			}
		}
	}

	return diagnoses, nil
}

// generateDetailTindakan creates procedure records: 1-2 distinct per rekam.
// Enforces no duplicate (id_rekam, kode_tindakan) pairs.
func generateDetailTindakan(targetCount int, rekams []model.RekamMedis, rng *rand.Rand) ([]model.DetailTindakan, error) {
	tindakans := make([]model.DetailTindakan, 0, targetCount)

	// Build ICD-9 code list from expanded master data
	icd9Codes := make([]string, len(expandedICD9Data))
	for i, t := range expandedICD9Data {
		icd9Codes[i] = t.KodeTindakan
	}

	if len(icd9Codes) == 0 {
		return nil, fmt.Errorf("no ICD-9 codes available")
	}

	for _, rekam := range rekams {
		// Each rekam gets 1-2 distinct tindakan
		numTindakan := 1 + rng.Intn(2) // 1 or 2

		usedCodes := make(map[string]bool)

		for j := 0; j < numTindakan; j++ {
			// Pick a code not yet used for this rekam
			var kodeTindakan string
			attempts := 0
			for {
				kodeTindakan = icd9Codes[rng.Intn(len(icd9Codes))]
				if !usedCodes[kodeTindakan] {
					break
				}
				attempts++
				if attempts > 100 {
					break
				}
			}
			usedCodes[kodeTindakan] = true

			// Jumlah: 1-5
			jumlah := rng.Intn(5) + 1

			tindakans = append(tindakans, model.DetailTindakan{
				IDRekam:      rekam.IDRekam,
				KodeTindakan: kodeTindakan,
				Jumlah:       jumlah,
			})

			if len(tindakans) >= targetCount {
				return tindakans, nil
			}
		}
	}

	return tindakans, nil
}

// generateKlaim creates claim records with DISTINCT id_rekam via shuffle.
// Each klaim references a unique rekam (UNIQUE constraint on id_rekam).
func generateKlaim(count int, rekams []model.RekamMedis, casemixIDs []uint, rng *rand.Rand) ([]model.Klaim, error) {
	if count > len(rekams) {
		return nil, fmt.Errorf("invariant violated: cannot create %d klaims from %d rekams (UNIQUE constraint)", count, len(rekams))
	}

	klaims := make([]model.Klaim, count)

	// Shuffle rekam indices to ensure distinctness
	rekamIndices := make([]int, len(rekams))
	for i := range rekamIndices {
		rekamIndices[i] = i
	}
	rng.Shuffle(len(rekamIndices), func(i, j int) {
		rekamIndices[i], rekamIndices[j] = rekamIndices[j], rekamIndices[i]
	})

	// Build CBGs code list
	cbgsCodes := make([]string, len(expandedCBGsData))
	cbgsTarifs := make(map[string]float64)
	for i, c := range expandedCBGsData {
		cbgsCodes[i] = c.KodeCBGS
		cbgsTarifs[c.KodeCBGS] = c.Tarif
	}

	if len(cbgsCodes) == 0 {
		return nil, fmt.Errorf("no CBGs codes available")
	}

	for i := 0; i < count; i++ {
		rekam := rekams[rekamIndices[i]]

		// KodeCBGS: random from master set
		kodeCBGS := cbgsCodes[rng.Intn(len(cbgsCodes))]
		baseTarif := cbgsTarifs[kodeCBGS]

		// IDPetugasCasemix: random from pool
		idPetugasCasemix := casemixIDs[rng.Intn(len(casemixIDs))]

		// StatusKlaim: draft ~10%, pending ~25%, disetujui ~50%, ditolak ~15%
		roll := rng.Float64()
		var statusKlaim string
		if roll < 0.10 {
			statusKlaim = "draft"
		} else if roll < 0.35 {
			statusKlaim = "pending"
		} else if roll < 0.85 {
			statusKlaim = "disetujui"
		} else {
			statusKlaim = "ditolak"
		}

		// AlasanPendingTolak: required for pending/ditolak, nil otherwise
		// Build FROM status, never the reverse
		var alasanPendingTolak *string
		if statusKlaim == "pending" || statusKlaim == "ditolak" {
			alasan := alasanPenolakan[rng.Intn(len(alasanPenolakan))]
			alasanPendingTolak = &alasan
		}

		// TanggalKlaim: on or after tanggal_kunjungan
		daysAfterKunjungan := rng.Intn(30) // 0-29 days after
		tanggalKlaim := rekam.TanggalKunjungan.AddDate(0, 0, daysAfterKunjungan)

		// NominalKlaim: baseTarif ±10%, never negative
		variance := (rng.Float64()*0.2 - 0.1) // -0.1 to +0.1
		nominalKlaim := baseTarif * (1 + variance)
		if nominalKlaim < 0 {
			nominalKlaim = 0
		}

		klaims[i] = model.Klaim{
			IDRekam:            rekam.IDRekam,
			KodeCBGS:           kodeCBGS,
			IDPetugasCasemix:   idPetugasCasemix,
			StatusKlaim:        statusKlaim,
			TanggalKlaim:       tanggalKlaim,
			NominalKlaim:       nominalKlaim,
			AlasanPendingTolak: alasanPendingTolak,
		}
	}

	return klaims, nil
}

// getUserIDsByRole queries users by role and returns their IDs
func getUserIDsByRole(db *gorm.DB, roleName string) ([]uint, error) {
	var role models.Role
	if err := db.Where("nama_role = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}

	var users []models.User
	if err := db.Where("id_role = ? AND deleted_at IS NULL", role.IDRole).Find(&users).Error; err != nil {
		return nil, err
	}

	ids := make([]uint, len(users))
	for i, u := range users {
		ids[i] = u.IDUser
	}

	return ids, nil
}

// insertBatch inserts rows in batches with transaction wrapping
func insertBatch(db *gorm.DB, data interface{}, batchSize int, tableName string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(data, batchSize).Error; err != nil {
			return fmt.Errorf("failed to insert %s: %w", tableName, err)
		}
		return nil
	})
}
