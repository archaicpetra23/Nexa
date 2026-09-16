## TUGAS MATKUL RPL - Semester 3 - Razan Rafi Akmaluzzuhair || 2510101014

## TEMA 1: Nexa: Sistem Informasi Pengelolaan Data Casemix Terintegrasi

### 1. Masalah Nyata yang Akan Diselesaikan

Unit Casemix merupakan sentral pengelolaan pembiayaan klaim jaminan kesehatan (BPJS) di Rumah Sakit Umum Daerah (RSUD). Namun, alur pencatatan data saat ini masih manual dan terfragmentasi antar-unit (loket pendaftaran, poliklinik/bangsal, rekam medis, Casemix, dan bagian keuangan). Kondisi ini menyebabkan:

* **Pencarian berkas rekam medis yang sangat lambat**, memakan waktu 15–60 menit untuk berkas fisik atau 5–15 menit pada sistem digital yang tidak terintegrasi.

* **Tingginya kesalahan input (*human error*) dan duplikasi kode medis**, seperti ketidaksesuaian kode diagnosis (ICD-10) dan tindakan medis (ICD-9 CM) akibat ketiadaan validasi otomatis saat penginputan oleh dokter dan perawat.

* **Tingginya angka klaim tertahan (*pending*) atau ditolak (*dispute*) oleh verifikator BPJS**, tanpa adanya sistem pencatatan terpusat mengenai alasan penolakan berkas sehingga unit terkait lambat memperbaiki data yang keliru.

### 2. Profil Target Pengguna

* **Petugas Registrasi / RM**: Mendaftarkan identitas pasien baru secara akurat, memverifikasi NIK/Nomor BPJS, serta memastikan tidak terjadi pencatatan rekam medis ganda.

* **Dokter DPJP (Penanggung Jawab Pelayanan)**: Menginput anamnesis, instruksi klinis, serta memilih kode diagnosis ICD-10 (primer dan sekunder) secara cepat dan tepat melalui sistem.

* **Perawat Bangsal / Poli**: Mencatat rincian prosedur dan tindakan medis (ICD-9 CM) yang telah dilakukan beserta kuantitasnya sebelum berkas diverifikasi.

* **Petugas / Koder Casemix**: Memvalidasi koding medis, menetapkan paket tarif INA-CBGs secara otomatis, memverifikasi kelayakan klaim, dan mendokumentasikan alasan klaim jika berstatus *pending* atau ditolak.

* **Staf Keuangan**: Memantau daftar berkas klaim yang berstatus *disetujui* untuk kebutuhan rekonsiliasi penerimaan dana klaim rumah sakit.

* **Manajemen / Direksi RS**: Mengakses laporan eksekutif dan grafik visualisasi penyakit terbanyak untuk pengambilan keputusan strategis rumah sakit.

* **Admin TI**: Mengelola hak akses akun pengguna dan memantau riwayat aktivitas mutasi data rekam medis (*audit trail*).

### 3. Manfaat Aplikasi

* Mengintegrasikan alur kerja seluruh unit layanan rumah sakit ke dalam satu basis data terpusat, menghilangkan *silo* data antar-departemen.

* Memangkas waktu temu-balik (*retrieval*) berkas pasien dan master kode ICD dari puluhan menit menjadi di bawah hitungan detik ($< 1 - 5\text{ detik}$).

* Menghilangkan insiden duplikasi input kode diagnosis dan tindakan pada lembar kunjungan yang sama melalui validasi *database constraint*.

* Mengurangi beban kerja koder Casemix minimal 50% melalui fitur rekomendasi tarif paket INA-CBGs otomatis.

* Meningkatkan transparansi dan akuntabilitas siklus klaim dengan pencatatan alasan kendala berkas yang jelas saat klaim ditolak atau tertunda.

* Menjamin integritas dan keamanan data rekam medis pasien sesuai standar OWASP Top 10 dan regulasi rekam medis (dilengkapi *Role-Based Access Control*, proteksi IDOR, dan pencatatan riwayat audit otomatis).

### 4. Daftar Fitur Inti

* **Sistem Autentikasi & RBAC Multi-Role**: Pembatasan hak akses menu berbasis 7 peran pengguna menggunakan sesi aman JWT via *HttpOnly Cookie*.

* **Manajemen Identitas Pasien (CRUD)**: Pendaftaran data pasien dengan validasi wajib format NIK (16 digit), No BPJS (13 digit), serta fitur pencarian instan.

* **Pencatatan Rekam Medis Kunjungan**: Pengelolaan data kunjungan rawat jalan/rawat inap, anamnesis keluhan, instruksi klinis, serta kalkulasi lama rawat (*Length of Stay*) otomatis.

* **Pencarian Cepat Master ICD (Autocomplete)**: Fitur pencarian teks diagnosis (ICD-10) dan tindakan (ICD-9 CM) dengan respon cepat ($< 300\text{ ms}$) menggunakan indeks trigram, lengkap dengan pemisahan status diagnosis primer/sekunder.

* **Engine Rekomendasi Tarif INA-CBGs**: Otomasi pencocokan paket tarif klaim BPJS berdasarkan kombinasi diagnosis primer dan tindakan yang dimasukkan.

* **Manajemen Siklus Berkas Klaim (State Machine)**: Pemantauan status berkas (*draft*, *pending*, *disetujui*, *ditolak*) dengan kewajiban input uraian catatan kendala jika status berkas *pending* atau ditolak.

* **Dasbor Analitik Eksekutif**: Visualisasi grafik 10 besar penyakit terbanyak (*Top 10 ICD-10*) dan ringkasan metrik status penerimaan klaim secara *real-time*.

* **Audit Trail System Log**: Pencatatan otomatis setiap aksi mutasi data klinis (*INSERT*, *UPDATE*, *DELETE*) dan riwayat *login/logout* ke dalam tabel log audit.

### 5. Fitur yang Tidak Dikerjakan

* Sistem inventaris fisik gudang farmasi, stok obat-obatan, dan bahan medis habis pakai (BMHP).

* *Direct-bridging* otomatis ke server produksi BPJS Kesehatan (API VClaim/TrustMark).

* Pembukuan akuntansi umum rumah sakit (*General Ledger*, neraca keuangan, dan kas keluar/masuk umum).

* Modul kasir pembayaran pasien umum non-BPJS atau integrasi *payment gateway*.

* Fitur transkripsi suara ke teks (*voice-to-text*) untuk penulisan resume rekam medis dokter.

### 6. Kriteria Aplikasi Dinyatakan Berhasil

* Alur kerja operasional dapat berjalan secara menyeluruh tanpa galat: Registrasi Pasien $\rightarrow$ Input Resume DPJP $\rightarrow$ Input Tindakan Perawat $\rightarrow$ Grouping INA-CBGs & Verifikasi Klaim Casemix.

* Pencarian data pasien dan kode master ICD-10/ICD-9 CM beroperasi sangat cepat dengan waktu tanggap $< 1 - 5\text{ detik}$.

* Tidak terjadi duplikasi kode diagnosis atau tindakan pada kunjungan rekam medis yang sama (0% duplikasi pada basis data).

* Sistem menolak perubahan status klaim menjadi *pending* atau *ditolak* apabila kolom catatan alasan kendala tidak diisi.

* Antarmuka web (Vue 3 + PrimeVue) me-render menu navigasi secara dinamis sesuai hak akses akun yang aktif.

* Seluruh data transaksi klinis tersimpan secara aman dan persisten di basis data PostgreSQL 16 dalam kontainer Docker, dengan pencatatan riwayat perubahan data pada tabel log aktivitas.

#### TEMA 2: FinanceTracker / FinForecast

#### TEMA 3: Docker Container Monitor & Manager