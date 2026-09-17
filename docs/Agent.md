# AGENT.md — Panduan AI Agent untuk Proyek Nexa

Dokumen ini dibaca oleh **AI coding agent apa pun** (Claude Code, Cursor, GitHub Copilot Chat/Agent, Codex, Windsurf, Gemini CLI, dll.) sebelum mengerjakan task di repo ini. Tujuannya: hasil kerja tetap konsisten meskipun dikerjakan lintas model/sesi yang berbeda-beda.

> **Urutan baca wajib bagi agent:** `AGENT.md` (dokumen ini) → `CODING_STYLE.md` → `LOGGING.md` → baru mulai kerja.

---

## 1. Konteks Proyek (ringkas)

- **Nama:** Nexa — Sistem Informasi Pengelolaan Data Casemix Terintegrasi (RSUD, klaim BPJS/JKN).
- **Stack:** Go 1.22+ (Gin, GORM) · Vue 3 (Composition API, Vite, PrimeVue, Pinia, Tailwind) · PostgreSQL 16.
- **Arsitektur:** Decoupled monorepo (`backend/` REST API, `frontend/` SPA), layer backend `delivery → usecase → repository → model`.
- **Domain sensitif:** rekam medis pasien, diagnosis ICD-10, tindakan ICD-9, klaim BPJS. Ini adalah **data kesehatan** — perlakukan setiap perubahan dengan asumsi ada implikasi hukum & privasi.
- Detail lengkap requirement ada di PRD proyek; detail gaya kode ada di `CODING_STYLE.md`.

---

## 2. Aturan Main (berlaku untuk semua agent, model apa pun)

1. **Baca sebelum menulis.** Jangan generate ulang file yang sudah ada tanpa membaca isinya dulu. Jangan asumsikan struktur folder — cek langsung.
2. **Ikuti `CODING_STYLE.md` secara ketat** — penamaan, layering, format response API, konvensi DB. Jangan perkenalkan pola baru (state management lain, ORM lain, format response lain) tanpa izin eksplisit dari user.
3. **Jangan pernah membuat perubahan skema database secara langsung ke DDL produksi.** Perubahan struktur tabel harus dalam bentuk file migrasi baru, bukan mengedit `init.sql` yang sudah dipakai.
4. **Jangan menonaktifkan middleware keamanan** (`RBAC`, JWT auth, rate limiter, CORS) untuk "mempermudah testing" kecuali diminta eksplisit dan disertai peringatan bahwa itu tidak boleh terbawa ke commit.
5. **Jangan hardcode secret/credential** (password DB, JWT secret, API key) di kode. Gunakan `.env` / config yang sudah ada.
6. **Jangan melakukan hard delete** pada tabel transaksional (`pasien`, `rekam_medis`, `klaim`, dst). Selalu soft delete via `deleted_at`.
7. **Setiap endpoint/fitur baru yang menyentuh data milik user lain wajib ada pengecekan kepemilikan (anti-IDOR)** sesuai contoh di `CODING_STYLE.md` §2.7.
8. **Jangan mengubah alur bisnis kritikal** (state machine status klaim, rumus LoS, validasi `alasan_pending_tolak` wajib) tanpa konfirmasi dari user — ini berdampak ke keuangan RS.
9. **Kalau ambigu, tanya — jangan menebak** untuk hal yang berisiko tinggi (skema DB, alur klaim, RBAC). Untuk hal berisiko rendah (styling, penamaan variabel lokal), boleh ambil keputusan wajar dan sebutkan asumsinya.
10. **Jangan menghapus atau menimpa kode/fitur yang tidak diminta untuk diubah.** Kalau perlu refactor besar di luar scope task, usulkan dulu, jangan langsung eksekusi.

---

## 3. Alur Kerja Standar per Task

1. **Pahami task** — baca kode terkait, cek apakah ada fitur serupa yang sudah ada sebagai referensi pola.
2. **Rencanakan singkat** — file apa saja yang akan dibuat/diubah, apakah menyentuh DB schema, apakah menyentuh auth/RBAC.
3. **Implementasi** mengikuti `CODING_STYLE.md`.
4. **Uji** — minimal jalankan/pastikan tidak ada error kompilasi (`go build`), lint bersih, dan untuk fitur penting sertakan/perbarui test bila memungkinkan.
5. **Catat ke log** — tambahkan entri baru di `docs/DEVLOG.md` sesuai format di `LOGGING.md`, **setiap kali** menyelesaikan pekerjaan (fitur baru, fix bug, refactor) — lihat §4.
6. **Ringkas ke user** — jelaskan apa yang berubah, file apa saja yang tersentuh, dan apakah ada follow-up/risiko yang perlu diperhatikan.

---

## 4. Kewajiban Logging (wajib, lihat `LOGGING.md`)

Setiap agent AI yang menyelesaikan pekerjaan (fitur/fix/refactor/investigasi error) **wajib** menambahkan satu entri baru ke `docs/DEVLOG.md`, mengikuti template dan aturan di `LOGGING.md`. Tujuannya supaya user bisa lihat riwayat "sedang bikin apa, sudah fix apa, dan error apa yang masih terbuka" — walau pengerjaan berpindah-pindah model AI.

**Ini bukan opsional.** Jika sebuah sesi kerja tidak menghasilkan entri log baru, anggap task belum selesai.

---

## 5. Identitas Agent

Karena project ini akan dikerjakan oleh beberapa model AI berbeda, setiap agent **wajib mengisi nama modelnya sendiri** di kolom "Agent" pada entri log (lihat `LOGGING.md`), misalnya `Claude (Sonnet 5)`, `GPT-5 Codex`, `Gemini CLI`, `Copilot Agent`, dst — supaya user tahu siapa mengerjakan apa.

---

## 6. Larangan Khusus (Hard Rules)

- [DILARANG] Menghapus data audit trail (`log_aktivitas`) atau menonaktifkan pencatatannya.
- [DILARANG] Mengembalikan `password_hash` atau field sensitif lain di response API.
- [DILARANG] Menyimpan JWT di `localStorage`/`sessionStorage` di frontend (harus tetap via `HttpOnly` cookie).
- [DILARANG] Query SQL manual yang tidak parameterized (rawan SQL Injection).
- [DILARANG] Commit langsung ke branch utama tanpa melalui PR (kecuali diarahkan lain oleh user).
- [DILARANG] Mengklaim "sudah ditest" tanpa benar-benar menjalankan build/test.

> Catatan gaya penulisan: dokumen ini dan seluruh dokumentasi proyek **tidak menggunakan emoji/emoticon** — lihat `CODING_STYLE.md` §0.2. Gunakan label teks (`[DILARANG]`, `[WAJIB]`, dst.) atau icon pack (`lucide-vue-next`/`PrimeIcons`) di UI.

---

## 7. Kalau Agent Tidak Yakin

Jika sebuah instruksi dari user tampak bertentangan dengan dokumen ini (misalnya diminta menonaktifkan RBAC untuk demo), agent boleh mengikuti instruksi user setelah **mengingatkan risikonya secara eksplisit** — jangan diam-diam melanggar aturan tanpa memberi tahu.
