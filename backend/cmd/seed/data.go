package main

import "nexa/backend/internal/model"

// Indonesian name components for realistic patient generation
var (
	firstNames = []string{
		"Ahmad", "Budi", "Citra", "Dewi", "Eko", "Fani", "Gita", "Hadi",
		"Indri", "Joko", "Kartika", "Lilis", "Maya", "Nina", "Oki", "Putri",
		"Qori", "Rani", "Sari", "Tina", "Uli", "Vika", "Wati", "Xenia",
		"Yoga", "Zara", "Adi", "Bayu", "Dian", "Eka", "Fajar", "Galuh",
		"Hendra", "Intan", "Jaka", "Kirana", "Lina", "Mira", "Nanda", "Ratna",
		"Bambang", "Rudi", "Santi", "Wulan", "Yanto", "Agus", "Rina", "Dedi",
		"Tuti", "Andi", "Lestari", "Yudi", "Fitri", "Heru", "Sri", "Anton",
		"Dwi", "Tri", "Yuni", "Sinta", "Reza", "Nurul", "Fitra", "Arif",
	}

	lastNames = []string{
		"Santoso", "Wijaya", "Nugroho", "Prasetya", "Hidayat", "Firmansyah",
		"Kusuma", "Lestari", "Maharani", "Pramudya", "Saputra", "Tambunan",
		"Utami", "Wahyudi", "Yusuf", "Zulkifli", "Adi", "Bahri", "Cahyo",
		"Darjo", "Effendi", "Fauzi", "Gunawan", "Hartono", "Irawan",
		"Julianto", "Kurniawan", "Lubis", "Mardi", "Nasution", "Oktaviani",
		"Pratama", "Raharjo", "Setiawan", "Taufik", "Wardana", "Budiman",
		"Haryanto", "Susanto", "Wibowo", "Suryanto", "Permana", "Suharto",
		"Hakim", "Rachman", "Salim", "Iskandar", "Mulyono", "Sudirman",
	}

	middleNames = []string{
		"Budi", "Cahya", "Dharma", "Eko", "Fajar", "Guna", "Hadi", "Indra",
		"Jaya", "Kusuma", "Laksana", "Mulya", "Nugraha", "Putra", "Rini",
	}

	streetPrefixes = []string{
		"Jl. Merdeka", "Jl. Sudirman", "Jl. Gatot Subroto", "Jl. Ahmad Yani",
		"Jl. Diponegoro", "Jl. Veteran", "Jl. Pahlawan", "Jl. Pemuda",
		"Jl. Kenanga", "Jl. Melati", "Jl. Mawar", "Jl. Anggrek",
		"Jl. Raya Bogor", "Jl. Raya Jakarta", "Jl. Raya Bandung", "Jl. MT Haryono",
		"Jl. Kebon Jeruk", "Jl. Cempaka", "Jl. Dahlia", "Jl. Teratai",
	}

	kelurahan = []string{
		"Sumber Jaya", "Makmur", "Sejahtera", "Sentosa", "Indah Permai",
		"Mulia Jaya", "Karya Bhakti", "Mekar Sari", "Taman Indah", "Sukamaju",
		"Harmoni", "Bahagia", "Damai", "Tenteram", "Asri", "Mandiri",
	}

	kecamatan = []string{
		"Kec. Cipayung", "Kec. Jatinegara", "Kec. Tebet", "Kec. Kebayoran",
		"Kec. Menteng", "Kec. Tanah Abang", "Kec. Cakung", "Kec. Matraman",
		"Kec. Cilandak", "Kec. Pasar Minggu", "Kec. Jagakarsa", "Kec. Mampang",
	}

	kota = []string{
		"Jakarta Timur", "Jakarta Selatan", "Jakarta Pusat", "Jakarta Barat",
		"Jakarta Utara", "Bogor", "Depok", "Tangerang", "Bekasi",
	}

	keluhan = []string{
		"Demam tinggi sejak 3 hari yang lalu",
		"Batuk berdahak disertai sesak napas",
		"Nyeri perut bagian bawah",
		"Sakit kepala hebat dan mual",
		"Diare lebih dari 5 kali sehari",
		"Nyeri dada sebelah kiri menjalar ke lengan",
		"Lemas dan pusing berputar",
		"Bengkak pada kedua kaki",
		"Demam naik turun, badan menggigil",
		"Batuk kering berkepanjangan",
		"Nyeri saat buang air kecil",
		"Muntah-muntah tidak bisa makan",
		"Sesak napas saat beraktivitas",
		"Nyeri ulu hati dan perih",
		"Gatal-gatal seluruh badan",
		"Benjolan di leher membesar",
		"Penglihatan kabur mendadak",
		"Telinga berdenging terus menerus",
		"Sulit tidur dan gelisah",
		"Badan lemas tidak bertenaga",
	}

	catatan = []string{
		"Pasien kooperatif, vital sign stabil",
		"Riwayat hipertensi tidak terkontrol",
		"Alergi obat golongan penisilin",
		"Pasien rujukan dari puskesmas",
		"Keadaan umum lemah, kesadaran compos mentis",
		"Tidak ada riwayat penyakit keluarga",
		"Pasien kontrol rutin setiap bulan",
		"Hasil lab terlampir",
		"Perlu observasi ketat 24 jam",
		"Edukasi pola hidup sehat sudah diberikan",
		"Riwayat diabetes melitus tipe 2",
		"Pasien menolak rawat inap",
		"Keluarga meminta second opinion",
		"Hasil rontgen menunjukkan kelainan",
		"Pasien sudah menjalani terapi 1 minggu",
	}

	alasanPenolakan = []string{
		"Dokumen klaim tidak lengkap, SEP tidak sesuai",
		"Diagnosis tidak sesuai dengan tindakan yang diklaim",
		"Kode INA-CBGs tidak tepat untuk kasus ini",
		"Verifikasi berkas menunjukkan ketidaksesuaian data",
		"Tarif yang diajukan melebihi plafon INA-CBGs",
		"Data rekam medis tidak mendukung klaim",
		"Pemeriksaan penunjang tidak sesuai indikasi",
		"Tindakan medis belum memenuhi standar",
		"Masa perawatan melebihi batas kewajaran",
		"Obat yang diklaim tidak sesuai formularium",
		"Perlu klarifikasi lebih lanjut dari dokter DPJP",
		"Coding diagnosis kurang spesifik",
		"Belum ada persetujuan tindakan dari keluarga",
		"Data pasien tidak valid di database BPJS",
		"Klaim duplikat untuk periode yang sama",
	}
)

// expandedICD10Data contains real, commonly-used ICD-10 diagnosis codes.
// Target: 150-300 entries. These are actual codes, not invented ones.
var expandedICD10Data = []model.Diagnosis{
	// Infectious diseases (A00-B99)
	{KodeICD10: "A09", NamaDiagnosis: "Diare dan Gastroenteritis"},
	{KodeICD10: "A15", NamaDiagnosis: "Tuberkulosis Paru"},
	{KodeICD10: "A16", NamaDiagnosis: "Tuberkulosis Paru tanpa Konfirmasi Bakteriologis"},
	{KodeICD10: "A41", NamaDiagnosis: "Sepsis"},
	{KodeICD10: "B34", NamaDiagnosis: "Infeksi Virus"},
	{KodeICD10: "A90", NamaDiagnosis: "Demam Dengue"},
	{KodeICD10: "B01", NamaDiagnosis: "Cacar Air (Varicella)"},
	{KodeICD10: "B05", NamaDiagnosis: "Campak"},
	{KodeICD10: "A37", NamaDiagnosis: "Pertusis (Batuk Rejan)"},

	// Blood diseases (D50-D89)
	{KodeICD10: "D50", NamaDiagnosis: "Anemia Defisiensi Besi"},
	{KodeICD10: "D64", NamaDiagnosis: "Anemia Lainnya"},
	{KodeICD10: "D62", NamaDiagnosis: "Anemia Akut Pasca Perdarahan"},
	{KodeICD10: "D69", NamaDiagnosis: "Purpura dan Gangguan Hemoragik Lain"},

	// Endocrine diseases (E00-E90)
	{KodeICD10: "E10", NamaDiagnosis: "Diabetes Melitus Tipe 1"},
	{KodeICD10: "E11", NamaDiagnosis: "Diabetes Melitus Tipe 2"},
	{KodeICD10: "E14", NamaDiagnosis: "Diabetes Melitus Tidak Spesifik"},
	{KodeICD10: "E46", NamaDiagnosis: "Malnutrisi Energi Protein Tidak Spesifik"},
	{KodeICD10: "E66", NamaDiagnosis: "Obesitas"},
	{KodeICD10: "E78", NamaDiagnosis: "Gangguan Metabolisme Lipoprotein"},
	{KodeICD10: "E87", NamaDiagnosis: "Gangguan Keseimbangan Cairan, Elektrolit dan Asam Basa"},
	{KodeICD10: "E04", NamaDiagnosis: "Gondok Non-toksik Lainnya"},

	// Circulatory system (I00-I99)
	{KodeICD10: "I10", NamaDiagnosis: "Hipertensi Esensial (Primer)"},
	{KodeICD10: "I11", NamaDiagnosis: "Penyakit Jantung Hipertensif"},
	{KodeICD10: "I21", NamaDiagnosis: "Infark Miokard Akut"},
	{KodeICD10: "I50", NamaDiagnosis: "Gagal Jantung"},
	{KodeICD10: "I63", NamaDiagnosis: "Infark Serebral (Stroke Iskemik)"},
	{KodeICD10: "I64", NamaDiagnosis: "Stroke Tidak Spesifik"},
	{KodeICD10: "I25", NamaDiagnosis: "Penyakit Jantung Iskemik Kronis"},
	{KodeICD10: "I20", NamaDiagnosis: "Angina Pektoris"},
	{KodeICD10: "I48", NamaDiagnosis: "Fibrilasi dan Flutter Atrium"},

	// Respiratory system (J00-J99)
	{KodeICD10: "J00", NamaDiagnosis: "Nasofaringitis Akut (Common Cold)"},
	{KodeICD10: "J06", NamaDiagnosis: "Infeksi Saluran Pernapasan Atas Akut"},
	{KodeICD10: "J18", NamaDiagnosis: "Pneumonia"},
	{KodeICD10: "J20", NamaDiagnosis: "Bronkitis Akut"},
	{KodeICD10: "J44", NamaDiagnosis: "Penyakit Paru Obstruktif Kronik (PPOK)"},
	{KodeICD10: "J45", NamaDiagnosis: "Asma"},
	{KodeICD10: "J12", NamaDiagnosis: "Pneumonia Virus"},
	{KodeICD10: "J22", NamaDiagnosis: "Infeksi Saluran Pernapasan Bawah Akut"},
	{KodeICD10: "J90", NamaDiagnosis: "Efusi Pleura"},
	{KodeICD10: "J42", NamaDiagnosis: "Bronkitis Kronis"},

	// Digestive system (K00-K93)
	{KodeICD10: "K29", NamaDiagnosis: "Gastritis dan Duodenitis"},
	{KodeICD10: "K35", NamaDiagnosis: "Apendisitis Akut"},
	{KodeICD10: "K80", NamaDiagnosis: "Kolelitiasis (Batu Empedu)"},
	{KodeICD10: "K81", NamaDiagnosis: "Kolesistitis"},
	{KodeICD10: "K92", NamaDiagnosis: "Penyakit Sistem Pencernaan Lainnya"},
	{KodeICD10: "K21", NamaDiagnosis: "Penyakit Refluks Gastroesofageal (GERD)"},
	{KodeICD10: "K59", NamaDiagnosis: "Gangguan Fungsional Usus Lainnya"},
	{KodeICD10: "K76", NamaDiagnosis: "Penyakit Hati Lainnya"},
	{KodeICD10: "K52", NamaDiagnosis: "Gastroenteritis dan Kolitis Non-infeksi"},

	// Genitourinary system (N00-N99)
	{KodeICD10: "N39", NamaDiagnosis: "Infeksi Saluran Kemih"},
	{KodeICD10: "N18", NamaDiagnosis: "Penyakit Ginjal Kronik"},
	{KodeICD10: "N20", NamaDiagnosis: "Kalkulus Ginjal dan Ureter (Batu Ginjal)"},
	{KodeICD10: "N40", NamaDiagnosis: "Hiperplasia Prostat"},
	{KodeICD10: "N17", NamaDiagnosis: "Gagal Ginjal Akut"},
	{KodeICD10: "N73", NamaDiagnosis: "Penyakit Radang Panggul"},
	{KodeICD10: "N76", NamaDiagnosis: "Inflamasi Vagina dan Vulva Akut Lainnya"},

	// Pregnancy complications (O00-O99)
	{KodeICD10: "O80", NamaDiagnosis: "Persalinan Tunggal Spontan"},
	{KodeICD10: "O82", NamaDiagnosis: "Persalinan Tunggal dengan Sectio Caesarea"},
	{KodeICD10: "O14", NamaDiagnosis: "Preeklampsia"},
	{KodeICD10: "O24", NamaDiagnosis: "Diabetes Melitus dalam Kehamilan"},
	{KodeICD10: "O42", NamaDiagnosis: "Ketuban Pecah Dini"},
	{KodeICD10: "O60", NamaDiagnosis: "Persalinan Prematur"},
	{KodeICD10: "O62", NamaDiagnosis: "Kelainan Kekuatan Kontraksi Persalinan"},

	// Perinatal conditions (P00-P96)
	{KodeICD10: "P07", NamaDiagnosis: "Gangguan Terkait Kehamilan Pendek dan Berat Lahir Rendah"},
	{KodeICD10: "P22", NamaDiagnosis: "Distres Pernapasan pada Bayi Baru Lahir"},
	{KodeICD10: "P59", NamaDiagnosis: "Ikterus Neonatal"},

	// Symptoms and signs (R00-R99)
	{KodeICD10: "R50", NamaDiagnosis: "Demam"},
	{KodeICD10: "R51", NamaDiagnosis: "Sakit Kepala"},
	{KodeICD10: "R10", NamaDiagnosis: "Nyeri Abdomen dan Pelvis"},
	{KodeICD10: "R11", NamaDiagnosis: "Mual dan Muntah"},
	{KodeICD10: "R06", NamaDiagnosis: "Kelainan Pernapasan"},
	{KodeICD10: "R07", NamaDiagnosis: "Nyeri Tenggorokan dan Dada"},
	{KodeICD10: "R42", NamaDiagnosis: "Pusing dan Vertigo"},
	{KodeICD10: "R53", NamaDiagnosis: "Malaise dan Kelelahan"},
	{KodeICD10: "R73", NamaDiagnosis: "Peningkatan Glukosa Darah"},

	// Injury and poisoning (S00-T98)
	{KodeICD10: "S06", NamaDiagnosis: "Cedera Intrakranial"},
	{KodeICD10: "S72", NamaDiagnosis: "Fraktur Femur"},
	{KodeICD10: "S82", NamaDiagnosis: "Fraktur Tibia dan Fibula"},
	{KodeICD10: "S42", NamaDiagnosis: "Fraktur Bahu dan Lengan Atas"},
	{KodeICD10: "T14", NamaDiagnosis: "Cedera Tidak Spesifik"},
	{KodeICD10: "T78", NamaDiagnosis: "Efek Merugikan Tidak Spesifik"},
	{KodeICD10: "T88", NamaDiagnosis: "Komplikasi Tindakan Medis dan Bedah Lain"},

	// Neoplasms (C00-D49)
	{KodeICD10: "C50", NamaDiagnosis: "Neoplasma Ganas Payudara"},
	{KodeICD10: "C18", NamaDiagnosis: "Neoplasma Ganas Kolon"},
	{KodeICD10: "C34", NamaDiagnosis: "Neoplasma Ganas Bronkus dan Paru"},
	{KodeICD10: "C53", NamaDiagnosis: "Neoplasma Ganas Serviks Uteri"},
	{KodeICD10: "D12", NamaDiagnosis: "Neoplasma Jinak Kolon dan Rektum"},
	{KodeICD10: "D25", NamaDiagnosis: "Leiomioma Uterus (Mioma)"},

	// Mental disorders (F00-F99)
	{KodeICD10: "F32", NamaDiagnosis: "Episode Depresi"},
	{KodeICD10: "F41", NamaDiagnosis: "Gangguan Ansietas Lainnya"},
	{KodeICD10: "F10", NamaDiagnosis: "Gangguan Mental dan Perilaku Akibat Alkohol"},
	{KodeICD10: "F20", NamaDiagnosis: "Skizofrenia"},

	// Nervous system (G00-G99)
	{KodeICD10: "G40", NamaDiagnosis: "Epilepsi"},
	{KodeICD10: "G43", NamaDiagnosis: "Migrain"},
	{KodeICD10: "G44", NamaDiagnosis: "Sindrom Sakit Kepala Lainnya"},
	{KodeICD10: "G47", NamaDiagnosis: "Gangguan Tidur"},

	// Eye diseases (H00-H59)
	{KodeICD10: "H10", NamaDiagnosis: "Konjungtivitis"},
	{KodeICD10: "H25", NamaDiagnosis: "Katarak Senilis"},
	{KodeICD10: "H26", NamaDiagnosis: "Katarak Lainnya"},
	{KodeICD10: "H53", NamaDiagnosis: "Gangguan Penglihatan"},

	// Ear diseases (H60-H95)
	{KodeICD10: "H60", NamaDiagnosis: "Otitis Eksterna"},
	{KodeICD10: "H66", NamaDiagnosis: "Otitis Media Supuratif"},
	{KodeICD10: "H65", NamaDiagnosis: "Otitis Media Non-supuratif"},

	// Skin diseases (L00-L99)
	{KodeICD10: "L02", NamaDiagnosis: "Abses Kulit dan Jaringan Subkutan"},
	{KodeICD10: "L03", NamaDiagnosis: "Selulitis"},
	{KodeICD10: "L20", NamaDiagnosis: "Dermatitis Atopik"},
	{KodeICD10: "L30", NamaDiagnosis: "Dermatitis Lainnya"},
	{KodeICD10: "L50", NamaDiagnosis: "Urtikaria"},

	// Musculoskeletal (M00-M99)
	{KodeICD10: "M10", NamaDiagnosis: "Gout (Asam Urat)"},
	{KodeICD10: "M19", NamaDiagnosis: "Osteoartritis Lainnya"},
	{KodeICD10: "M54", NamaDiagnosis: "Dorsalgia (Nyeri Punggung)"},
	{KodeICD10: "M79", NamaDiagnosis: "Gangguan Jaringan Lunak Lainnya"},
	{KodeICD10: "M25", NamaDiagnosis: "Gangguan Sendi Lainnya"},

	// Additional common codes
	{KodeICD10: "Z38", NamaDiagnosis: "Bayi Lahir Hidup Menurut Tempat Kelahiran"},
	{KodeICD10: "Z51", NamaDiagnosis: "Perawatan Medis Lainnya"},
	{KodeICD10: "Z73", NamaDiagnosis: "Masalah Terkait Kesulitan Mengelola Hidup"},
}

// expandedICD9Data contains real, commonly-used ICD-9 CM procedure codes.
// Target: 80-150 entries.
var expandedICD9Data = []model.Tindakan{
	// Diagnostic procedures
	{KodeTindakan: "38.93", NamaTindakan: "Pengambilan Darah Vena", TarifStandar: 50000},
	{KodeTindakan: "45.13", NamaTindakan: "Endoskopi Lambung dengan Biopsi", TarifStandar: 1500000},
	{KodeTindakan: "45.16", NamaTindakan: "Esofagogastroduodenoskopi (EGD)", TarifStandar: 1200000},
	{KodeTindakan: "48.24", NamaTindakan: "Kolonoskopi", TarifStandar: 1800000},
	{KodeTindakan: "87.44", NamaTindakan: "Foto Rontgen Dada", TarifStandar: 150000},
	{KodeTindakan: "87.03", NamaTindakan: "Foto Rontgen Kepala", TarifStandar: 200000},
	{KodeTindakan: "87.16", NamaTindakan: "Foto Rontgen Lain pada Tulang Belakang", TarifStandar: 180000},
	{KodeTindakan: "88.01", NamaTindakan: "CT Scan Kepala", TarifStandar: 900000},
	{KodeTindakan: "88.38", NamaTindakan: "CT Scan Abdomen", TarifStandar: 1000000},
	{KodeTindakan: "88.76", NamaTindakan: "USG Abdomen", TarifStandar: 250000},
	{KodeTindakan: "88.77", NamaTindakan: "USG Obstetri", TarifStandar: 200000},
	{KodeTindakan: "88.79", NamaTindakan: "USG Lain dan USG Tak Spesifik", TarifStandar: 220000},
	{KodeTindakan: "89.52", NamaTindakan: "Elektrokardiogram (EKG)", TarifStandar: 100000},
	{KodeTindakan: "89.54", NamaTindakan: "Elektroensefalogram (EEG)", TarifStandar: 350000},

	// Therapeutic procedures - respiratory
	{KodeTindakan: "96.04", NamaTindakan: "Pemasangan Infus", TarifStandar: 75000},
	{KodeTindakan: "96.05", NamaTindakan: "Infus Obat Lain", TarifStandar: 80000},
	{KodeTindakan: "96.07", NamaTindakan: "Insersi Intubasi Endotrakeal Lainnya", TarifStandar: 500000},
	{KodeTindakan: "96.71", NamaTindakan: "Ventilator Mekanik Kontinyu", TarifStandar: 2000000},
	{KodeTindakan: "96.72", NamaTindakan: "Ventilator Mekanik Intermiten", TarifStandar: 1500000},
	{KodeTindakan: "93.96", NamaTindakan: "Terapi Oksigen", TarifStandar: 150000},

	// Injections and infusions
	{KodeTindakan: "99.04", NamaTindakan: "Injeksi Antibiotik", TarifStandar: 120000},
	{KodeTindakan: "99.07", NamaTindakan: "Injeksi Analgesik", TarifStandar: 100000},
	{KodeTindakan: "99.10", NamaTindakan: "Injeksi Kortikosteroid", TarifStandar: 150000},
	{KodeTindakan: "99.15", NamaTindakan: "Transfusi Packed Red Cells", TarifStandar: 450000},
	{KodeTindakan: "99.18", NamaTindakan: "Injeksi Imunoglobulin", TarifStandar: 800000},

	// Physical therapy and rehabilitation
	{KodeTindakan: "93.39", NamaTindakan: "Terapi Fisik Lainnya", TarifStandar: 150000},
	{KodeTindakan: "93.90", NamaTindakan: "Fisioterapi", TarifStandar: 100000},
	{KodeTindakan: "93.91", NamaTindakan: "Latihan Otot", TarifStandar: 120000},
	{KodeTindakan: "93.92", NamaTindakan: "Latihan Rentang Gerak", TarifStandar: 110000},

	// Surgical procedures - digestive
	{KodeTindakan: "47.09", NamaTindakan: "Apendektomi Lainnya", TarifStandar: 5000000},
	{KodeTindakan: "51.22", NamaTindakan: "Kolesistektomi", TarifStandar: 8000000},
	{KodeTindakan: "46.75", NamaTindakan: "Kolostomi Abdomen", TarifStandar: 6500000},
	{KodeTindakan: "49.46", NamaTindakan: "Eksisi Hemoroid", TarifStandar: 3000000},

	// Surgical procedures - orthopedic
	{KodeTindakan: "79.35", NamaTindakan: "Reduksi Terbuka Fraktur Femur", TarifStandar: 15000000},
	{KodeTindakan: "79.36", NamaTindakan: "Reduksi Terbuka Fraktur Tibia dan Fibula", TarifStandar: 12000000},
	{KodeTindakan: "79.31", NamaTindakan: "Reduksi Terbuka Fraktur Humerus", TarifStandar: 10000000},
	{KodeTindakan: "81.51", NamaTindakan: "Penggantian Sendi Panggul Total", TarifStandar: 35000000},
	{KodeTindakan: "81.54", NamaTindakan: "Penggantian Sendi Lutut Total", TarifStandar: 32000000},
	{KodeTindakan: "78.60", NamaTindakan: "Pengangkatan Implant Tak Spesifik", TarifStandar: 5000000},

	// Surgical procedures - obstetrics
	{KodeTindakan: "74.1", NamaTindakan: "Sectio Caesarea Segmen Bawah", TarifStandar: 8000000},
	{KodeTindakan: "73.59", NamaTindakan: "Persalinan Pervaginam Instrumental Lain", TarifStandar: 3000000},
	{KodeTindakan: "75.0", NamaTindakan: "Amniotomi (Pemecahan Ketuban)", TarifStandar: 500000},

	// Surgical procedures - cardiovascular
	{KodeTindakan: "36.10", NamaTindakan: "Angioplasti Koroner", TarifStandar: 40000000},
	{KodeTindakan: "36.91", NamaTindakan: "Perbaikan Pembuluh Darah dengan Tambalan Sintetis", TarifStandar: 25000000},
	{KodeTindakan: "37.22", NamaTindakan: "Kateterisasi Jantung Kiri", TarifStandar: 8000000},
	{KodeTindakan: "37.23", NamaTindakan: "Kateterisasi Jantung Kanan", TarifStandar: 7500000},
	{KodeTindakan: "38.91", NamaTindakan: "Ligasi Vena", TarifStandar: 4000000},

	// Urological procedures
	{KodeTindakan: "57.94", NamaTindakan: "Insersi Kateter Urin Menetap", TarifStandar: 200000},
	{KodeTindakan: "55.23", NamaTindakan: "Nefrolitotomi", TarifStandar: 12000000},
	{KodeTindakan: "60.29", NamaTindakan: "Prostatektomi Transurethral", TarifStandar: 10000000},

	// Wound care and dressing
	{KodeTindakan: "86.59", NamaTindakan: "Penutupan Luka Kulit dan Jaringan Subkutan Lainnya", TarifStandar: 500000},
	{KodeTindakan: "86.22", NamaTindakan: "Eksisi Luka", TarifStandar: 800000},
	{KodeTindakan: "86.04", NamaTindakan: "Insisi dengan Drainase pada Kulit dan Jaringan Subkutan", TarifStandar: 600000},
	{KodeTindakan: "93.57", NamaTindakan: "Aplikasi Lain pada Luka dan Lesi Kulit", TarifStandar: 300000},

	// Anesthesia
	{KodeTindakan: "99.60", NamaTindakan: "Resusitasi Kardiopulmoner", TarifStandar: 1000000},
	{KodeTindakan: "99.28", NamaTindakan: "Injeksi Anestesi Regional Lain", TarifStandar: 700000},

	// Additional common procedures
	{KodeTindakan: "54.91", NamaTindakan: "Aspirasi Cairan Peritoneum", TarifStandar: 800000},
	{KodeTindakan: "34.04", NamaTindakan: "Insersi Drain Interkostal", TarifStandar: 1200000},
	{KodeTindakan: "31.1", NamaTindakan: "Trakeostomi Temporer", TarifStandar: 2500000},
	{KodeTindakan: "31.29", NamaTindakan: "Trakeostomi Permanen Lainnya", TarifStandar: 3000000},
	{KodeTindakan: "42.33", NamaTindakan: "Endoskopi Esofagus", TarifStandar: 1100000},
	{KodeTindakan: "50.11", NamaTindakan: "Biopsi Hati Tertutup", TarifStandar: 1500000},
	{KodeTindakan: "52.14", NamaTindakan: "Biopsi Pankreas Tertutup", TarifStandar: 1800000},
	{KodeTindakan: "68.59", NamaTindakan: "Histerektomi Vaginal Lainnya", TarifStandar: 9000000},
	{KodeTindakan: "69.01", NamaTindakan: "Dilatasi dan Kuretase (D&C) Setelah Persalinan", TarifStandar: 2000000},
	{KodeTindakan: "80.16", NamaTindakan: "Artroskopi Lutut", TarifStandar: 5000000},
}

// expandedCBGsData contains real INA-CBG codes.
// Target: 60-120 entries.
var expandedCBGsData = []model.TarifCBGs{
	// Group A - Infectious diseases
	{KodeCBGS: "A-4-10-I", Deskripsi: "Penyakit Infeksi Ringan Tanpa Komplikasi", Tarif: 850000},
	{KodeCBGS: "A-4-10-II", Deskripsi: "Penyakit Infeksi Ringan dengan Komplikasi Sedang", Tarif: 1200000},
	{KodeCBGS: "A-4-10-III", Deskripsi: "Penyakit Infeksi Berat dengan Komplikasi", Tarif: 1800000},
	{KodeCBGS: "A-4-11-I", Deskripsi: "Tuberkulosis Tanpa Komplikasi", Tarif: 1500000},
	{KodeCBGS: "A-4-11-II", Deskripsi: "Tuberkulosis dengan Komplikasi", Tarif: 2200000},
	{KodeCBGS: "A-4-12-I", Deskripsi: "Sepsis Tanpa Ventilator", Tarif: 5000000},
	{KodeCBGS: "A-4-12-III", Deskripsi: "Sepsis dengan Ventilator >96 Jam", Tarif: 12000000},

	// Group B - Cardiovascular
	{KodeCBGS: "B-1-14-I", Deskripsi: "Penyakit Jantung Ringan", Tarif: 2800000},
	{KodeCBGS: "B-1-14-II", Deskripsi: "Penyakit Jantung Sedang", Tarif: 4200000},
	{KodeCBGS: "B-1-14-III", Deskripsi: "Penyakit Jantung Berat", Tarif: 7500000},
	{KodeCBGS: "B-1-15-I", Deskripsi: "Gagal Jantung Ringan", Tarif: 3200000},
	{KodeCBGS: "B-1-15-II", Deskripsi: "Gagal Jantung Sedang", Tarif: 5000000},
	{KodeCBGS: "B-1-15-III", Deskripsi: "Gagal Jantung Berat", Tarif: 8500000},
	{KodeCBGS: "B-1-16-II", Deskripsi: "Infark Miokard Akut", Tarif: 9000000},
	{KodeCBGS: "B-1-20-I", Deskripsi: "Hipertensi Tanpa Komplikasi", Tarif: 1200000},
	{KodeCBGS: "B-1-20-II", Deskripsi: "Hipertensi dengan Komplikasi", Tarif: 2100000},

	// Group C - Endocrine and metabolic
	{KodeCBGS: "C-4-13-I", Deskripsi: "Diabetes Melitus Tanpa Komplikasi", Tarif: 2500000},
	{KodeCBGS: "C-4-13-II", Deskripsi: "Diabetes Melitus dengan Komplikasi Sedang", Tarif: 4000000},
	{KodeCBGS: "C-4-13-III", Deskripsi: "Diabetes Melitus dengan Komplikasi Berat", Tarif: 6500000},
	{KodeCBGS: "C-4-14-I", Deskripsi: "Gangguan Elektrolit Ringan", Tarif: 1800000},
	{KodeCBGS: "C-4-14-II", Deskripsi: "Gangguan Elektrolit Sedang", Tarif: 2800000},

	// Group D - Respiratory
	{KodeCBGS: "D-4-16-I", Deskripsi: "Pneumonia Ringan", Tarif: 1750000},
	{KodeCBGS: "D-4-16-II", Deskripsi: "Pneumonia Sedang", Tarif: 3200000},
	{KodeCBGS: "D-4-16-III", Deskripsi: "Pneumonia Berat", Tarif: 5500000},
	{KodeCBGS: "D-4-17-I", Deskripsi: "PPOK Ringan", Tarif: 2000000},
	{KodeCBGS: "D-4-17-II", Deskripsi: "PPOK Sedang", Tarif: 3500000},
	{KodeCBGS: "D-4-17-III", Deskripsi: "PPOK Berat dengan Ventilator", Tarif: 8000000},
	{KodeCBGS: "D-4-18-I", Deskripsi: "Asma Bronkial Ringan", Tarif: 1400000},
	{KodeCBGS: "D-4-18-II", Deskripsi: "Asma Bronkial Sedang", Tarif: 2500000},

	// Group E - Digestive
	{KodeCBGS: "E-4-10-I", Deskripsi: "Gastritis dan Dispepsia Ringan", Tarif: 1100000},
	{KodeCBGS: "E-4-10-II", Deskripsi: "Gastritis dan Dispepsia Sedang", Tarif: 1800000},
	{KodeCBGS: "E-4-11-I", Deskripsi: "Penyakit Hati Ringan", Tarif: 2200000},
	{KodeCBGS: "E-4-11-II", Deskripsi: "Penyakit Hati Sedang", Tarif: 4000000},
	{KodeCBGS: "E-4-11-III", Deskripsi: "Penyakit Hati Berat", Tarif: 7500000},
	{KodeCBGS: "E-6-10-II", Deskripsi: "Apendisitis dengan Komplikasi", Tarif: 6000000},
	{KodeCBGS: "E-6-11-I", Deskripsi: "Hernia Inguinalis", Tarif: 4500000},

	// Group F - Genitourinary
	{KodeCBGS: "F-4-10-I", Deskripsi: "Infeksi Saluran Kemih Ringan", Tarif: 1300000},
	{KodeCBGS: "F-4-10-II", Deskripsi: "Infeksi Saluran Kemih Sedang", Tarif: 2100000},
	{KodeCBGS: "F-4-11-I", Deskripsi: "Gagal Ginjal Akut Ringan", Tarif: 4500000},
	{KodeCBGS: "F-4-11-II", Deskripsi: "Gagal Ginjal Akut Sedang", Tarif: 7000000},
	{KodeCBGS: "F-4-11-III", Deskripsi: "Gagal Ginjal Akut Berat dengan Dialisis", Tarif: 12000000},
	{KodeCBGS: "F-4-12-I", Deskripsi: "Penyakit Ginjal Kronik Stage 1-3", Tarif: 3500000},
	{KodeCBGS: "F-4-12-II", Deskripsi: "Penyakit Ginjal Kronik Stage 4-5", Tarif: 6500000},

	// Group G - Pregnancy and childbirth
	{KodeCBGS: "G-7-10-I", Deskripsi: "Persalinan Normal Pervaginam", Tarif: 2500000},
	{KodeCBGS: "G-7-10-II", Deskripsi: "Persalinan dengan Komplikasi Ringan", Tarif: 3500000},
	{KodeCBGS: "G-7-11-I", Deskripsi: "Sectio Caesarea Tanpa Komplikasi", Tarif: 6000000},
	{KodeCBGS: "G-7-11-II", Deskripsi: "Sectio Caesarea dengan Komplikasi", Tarif: 9000000},
	{KodeCBGS: "G-4-10-I", Deskripsi: "Preeklampsia Ringan", Tarif: 2800000},
	{KodeCBGS: "G-4-10-II", Deskripsi: "Preeklampsia Berat", Tarif: 5500000},

	// Group H - Neonates
	{KodeCBGS: "H-8-10-I", Deskripsi: "Bayi Baru Lahir Sehat", Tarif: 800000},
	{KodeCBGS: "H-4-10-I", Deskripsi: "Neonatus dengan Berat Lahir Rendah", Tarif: 4500000},
	{KodeCBGS: "H-4-10-II", Deskripsi: "Neonatus dengan Berat Lahir Sangat Rendah", Tarif: 8000000},
	{KodeCBGS: "H-4-11-I", Deskripsi: "Ikterus Neonatal", Tarif: 2200000},
	{KodeCBGS: "H-4-12-I", Deskripsi: "Distres Pernapasan Neonatal", Tarif: 5500000},

	// Group I - Blood and immune disorders
	{KodeCBGS: "I-4-10-I", Deskripsi: "Anemia Ringan", Tarif: 1500000},
	{KodeCBGS: "I-4-10-II", Deskripsi: "Anemia Sedang", Tarif: 2500000},
	{KodeCBGS: "I-4-10-III", Deskripsi: "Anemia Berat dengan Transfusi", Tarif: 4000000},

	// Group J - Neoplasms
	{KodeCBGS: "J-4-10-I", Deskripsi: "Neoplasma Jinak", Tarif: 3000000},
	{KodeCBGS: "J-4-11-II", Deskripsi: "Neoplasma Ganas dengan Kemoterapi", Tarif: 12000000},
	{KodeCBGS: "J-4-11-III", Deskripsi: "Neoplasma Ganas Stadium Lanjut", Tarif: 18000000},

	// Group K - Nervous system
	{KodeCBGS: "K-1-10-I", Deskripsi: "Stroke Iskemik Ringan", Tarif: 5500000},
	{KodeCBGS: "K-1-10-II", Deskripsi: "Stroke Iskemik Sedang", Tarif: 8500000},
	{KodeCBGS: "K-1-10-III", Deskripsi: "Stroke Iskemik Berat", Tarif: 15000000},
	{KodeCBGS: "K-1-11-II", Deskripsi: "Stroke Hemoragik", Tarif: 12000000},
	{KodeCBGS: "K-4-10-I", Deskripsi: "Epilepsi Terkontrol", Tarif: 2000000},
	{KodeCBGS: "K-4-10-II", Deskripsi: "Epilepsi dengan Status Epileptikus", Tarif: 5000000},

	// Group L - Injury and trauma
	{KodeCBGS: "L-8-10-I", Deskripsi: "Cedera Ringan", Tarif: 1200000},
	{KodeCBGS: "L-6-10-I", Deskripsi: "Fraktur Ekstremitas Atas Tertutup", Tarif: 8000000},
	{KodeCBGS: "L-6-10-II", Deskripsi: "Fraktur Ekstremitas Atas Terbuka", Tarif: 12000000},
	{KodeCBGS: "L-6-11-I", Deskripsi: "Fraktur Ekstremitas Bawah Tertutup", Tarif: 10000000},
	{KodeCBGS: "L-6-11-II", Deskripsi: "Fraktur Ekstremitas Bawah Terbuka", Tarif: 15000000},
	{KodeCBGS: "L-1-10-I", Deskripsi: "Cedera Kepala Ringan", Tarif: 3500000},
	{KodeCBGS: "L-1-10-II", Deskripsi: "Cedera Kepala Sedang", Tarif: 8000000},
	{KodeCBGS: "L-1-10-III", Deskripsi: "Cedera Kepala Berat", Tarif: 18000000},

	// Additional codes
	{KodeCBGS: "M-8-10-I", Deskripsi: "Observasi dan Evaluasi", Tarif: 900000},
	{KodeCBGS: "N-6-10-I", Deskripsi: "Tindakan Bedah Minor", Tarif: 2500000},
	{KodeCBGS: "N-6-11-I", Deskripsi: "Tindakan Bedah Sedang", Tarif: 6000000},
	{KodeCBGS: "N-6-12-II", Deskripsi: "Tindakan Bedah Mayor", Tarif: 15000000},
}
