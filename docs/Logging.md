# LOGGING.md — Konvensi Dev Log untuk Proyek Nexa

Dokumen ini mengatur bagaimana **AI agent** (dan kontributor manusia) mencatat pekerjaan yang dilakukan di proyek Nexa, supaya pemilik proyek selalu bisa menjawab: *"terakhir lagi bikin apa? sudah fix apa? error apa yang masih ngambang?"* — walau pengerjaan berpindah-pindah sesi/model AI.

File log aktualnya ada di: **`docs/DEVLOG.md`** (dibuat terpisah, lihat berkas pendamping). Dokumen ini adalah **aturan penulisannya**.

---

## 1. Kapan Harus Menulis Entri Log

Tulis satu entri baru setiap kali sebuah sesi kerja AI menghasilkan salah satu dari:

- Fitur baru selesai dibuat (walau sebagian/partial).
- Bug diperbaiki.
- Refactor / perubahan struktur kode.
- Perubahan skema database / migrasi.
- Error/bug ditemukan tapi **belum** selesai diperbaiki (dicatat sebagai `Open`).
- Perubahan terkait keamanan (RBAC, auth, dsb).

Tidak perlu menandai kategori dengan ikon — cukup isi kolom **Tipe** pada template entri (§3).

Jangan tulis entri untuk hal remeh (typo komentar, format ulang tanpa perubahan logic) — supaya log tidak jadi noise.

---

## 2. Lokasi & Struktur File

```
docs/
└── DEVLOG.md      # log kronologis, entri terbaru di PALING ATAS
```

Kalau `DEVLOG.md` sudah sangat panjang (>±6 bulan histori), boleh diarsipkan per tahun: `docs/devlog/2026.md`, `docs/devlog/2027.md`, dst — tapi format entri tetap sama.

---

## 3. Format Entri (wajib diikuti persis)

```markdown
## [YYYY-MM-DD HH:MM] — <Judul Singkat Pekerjaan>

- **Agent:** <nama model AI, mis. Claude Sonnet 5 / GPT-5 Codex / Gemini CLI>
- **Tipe:** Fitur Baru | Fix Bug | Refactor | Migrasi DB | Investigasi Error | Keamanan
- **Status:** Selesai | Sebagian | Open/Belum Fix
- **Modul:** <mis. Rekam Medis / Klaim / Auth / Dashboard>
- **File terdampak:**
  - `backend/internal/usecase/klaim_usecase.go`
  - `frontend/src/stores/claimStore.js`
- **Deskripsi:**
  Penjelasan singkat 1–4 kalimat: apa yang dikerjakan, kenapa, dan bagaimana caranya.
- **Error/Kendala (jika ada):**
  Pesan error atau gejala bug, plus root cause kalau sudah diketahui.
- **Next Step (jika belum selesai):**
  Apa yang perlu dilanjutkan di sesi berikutnya.
```

### Aturan pengisian tanggal & jam
- Format **wajib**: `YYYY-MM-DD HH:MM` (24 jam), zona waktu mengikuti waktu lokal user (WIB, UTC+7) kecuali disepakati lain.
- Agent AI **wajib mengambil tanggal/jam aktual saat itu** (jangan mengarang atau memakai placeholder) — kalau agent tidak punya akses waktu real-time, tanyakan ke user atau minta user isi manual di kolom itu.
- Entri baru selalu ditambahkan di **paling atas** file (reverse-chronological), supaya update terbaru langsung terlihat tanpa scroll.

---

## 4. Contoh Entri

```markdown
## [2026-09-17 14:20] — Tambah endpoint transisi status klaim

- **Agent:** Claude Sonnet 5
- **Tipe:** Fitur Baru
- **Status:** Selesai
- **Modul:** Klaim
- **File terdampak:**
  - `backend/internal/delivery/http/klaim_handler.go`
  - `backend/internal/usecase/klaim_usecase.go`
- **Deskripsi:**
  Menambahkan `PATCH /api/v1/klaim/:id/status` untuk transisi status draft→pending/disetujui/ditolak, termasuk validasi wajib isi `alasan_pending_tolak` saat status pending/ditolak.
- **Error/Kendala:** –
- **Next Step:** Tambahkan unit test untuk skenario status ditolak tanpa alasan (harus reject 422).

---

## [2026-09-17 10:05] — Fix bug pencarian ICD-10 lambat

- **Agent:** Claude Sonnet 5
- **Tipe:** Fix Bug
- **Status:** Selesai
- **Modul:** Master Data / Autocomplete
- **File terdampak:**
  - `backend/internal/repository/master_repository.go`
- **Deskripsi:**
  Query autocomplete ICD-10 sebelumnya full table scan karena index trigram belum kepakai (`ILIKE` tanpa `%%` pattern yang sesuai). Diperbaiki agar memakai `gin_trgm_ops` index yang sudah ada di DDL.
- **Error/Kendala:**
  Respon awalnya >2 detik untuk >10rb baris master ICD-10, seharusnya <300ms sesuai NFR-02.
- **Next Step:** –
```

---

## 5. Bagian Tambahan: Tabel Error Terbuka (opsional tapi disarankan)

Di bagian **paling bawah** `docs/DEVLOG.md`, jaga satu tabel ringkasan berisi error yang masih `Open`, supaya tidak perlu scroll semua histori:

```markdown
## Error / Isu yang Masih Terbuka

| Tanggal Ditemukan | Modul | Deskripsi Singkat | Status |
| :--- | :--- | :--- | :--- |
| 2026-09-15 | Klaim | Rekonsiliasi export CSV salah format tanggal | Open |
```

Setiap kali sebuah error di tabel ini selesai diperbaiki, **pindahkan barisnya** menjadi entri log normal (§3) dengan status Selesai, lalu hapus dari tabel ini.

---

## 6. Instruksi Khusus untuk AI Agent

1. Sebelum mulai kerja, **baca beberapa entri terakhir** di `docs/DEVLOG.md` untuk tahu konteks pekerjaan sebelumnya (siapa yang bikin apa, apa yang masih pending) — jangan mengerjakan ulang sesuatu yang sudah selesai, dan jangan lupakan error yang masih `Open`.
2. Setelah selesai kerja, **selalu tambahkan entri baru** — ini wajib, lihat `AGENT.md` §4.
3. Isi kolom **Agent** dengan jujur sesuai model yang sedang dipakai — jangan mengosongkan atau menyamarkan.
4. Jangan mengedit/menghapus entri log lama milik sesi lain, kecuali memindahkan status error dari tabel §5 ke entri baru sesuai §5.
5. Kalau task yang dikerjakan gagal total / dibatalkan, tetap catat sebagai entri dengan status Open/Belum Fix dan alasan kegagalannya — supaya user tahu itu sudah pernah dicoba dan tidak berhasil.
6. Jangan gunakan emoji/emoticon di mana pun dalam log — gunakan label teks polos sesuai template (`Selesai`, `Sebagian`, `Open/Belum Fix`), sesuai kebijakan ikon di `CODING_STYLE.md` §0.2.
