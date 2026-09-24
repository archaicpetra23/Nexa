# Nexa — UI Design Specification
## Sistem Informasi Pengelolaan Data Casemix Terintegrasi

**Versi:** 1.0.0
**Stack:** Vue 3 (Composition API) + PrimeVue + Tailwind CSS

---

## 1. Prinsip Desain

Nexa dipakai oleh staf rumah sakit (petugas RM, dokter, perawat, koder Casemix, keuangan, manajemen) dalam sesi kerja panjang, untuk data yang berdampak medis dan finansial. Desain harus mengutamakan **kejelasan, kecepatan input, dan kepercayaan** — bukan estetika dekoratif.

1. **Kejelasan di atas dekorasi** — tanpa gradient, tanpa bayangan tebal, border-radius kecil dan konsisten.
2. **Status selalu dua penanda** — warna + label/ikon, karena status klaim terkait keputusan finansial & medis yang tidak boleh ambigu.
3. **Data kritis selalu terlihat langsung** — nomor RM, status klaim, nama DPJP tidak boleh disembunyikan di balik hover/tooltip.
4. **Minim animasi** — hanya pada transisi status penting (misal klaim berubah jadi "Disetujui"), bukan di setiap hover card.
5. **Aksesibel untuk sesi kerja lama** — kontras tinggi tapi tidak stark white, mengurangi kelelahan mata.
6. **Menu & tampilan menyesuaikan role** — bukan menampilkan semua lalu men-disable yang tidak boleh diakses.

---

## 2. Color Palette

### Primary
| Token | Hex | Fungsi |
|---|---|---|
| `primary` | `#1C3D5A` | Sidebar, header, tombol aksi utama |
| `primary-hover` | `#2D5578` | Hover/active state |
| `primary-dark` | `#0F2740` | Item sidebar aktif, teks di atas primary |

### Background & Surface
| Token | Hex | Fungsi |
|---|---|---|
| `bg-base` | `#F7F6F3` | Background utama aplikasi |
| `surface` | `#FFFFFF` | Card, table, form container |
| `border` | `#E5E2DB` | Divider, border tipis antar elemen |

### Teks
| Token | Hex | Fungsi |
|---|---|---|
| `text-primary` | `#1F2937` | Teks utama/konten |
| `text-secondary` | `#6B7280` | Label, teks sekunder |
| `text-disabled` | `#9CA3AF` | Placeholder, elemen nonaktif |

### Status Klaim (semantic — khusus badge status, tidak dipakai di elemen lain)
| Status | Hex | Catatan |
|---|---|---|
| Draft | `#9CA3AF` | Abu netral |
| Pending | `#B7791F` | Amber gelap |
| Disetujui | `#2F7A4F` | Hijau tua |
| Ditolak | `#B54245` | Merah bata |

**Aturan pemakaian:**
- Semantic color status hanya untuk badge/label status klaim — jangan dipakai untuk tombol atau elemen dekoratif lain.
- `primary` adalah satu-satunya warna aksen kuat di UI — hindari menambah warna brand kedua agar hierarki visual tetap jelas.
- Kontras `text-primary` di atas `bg-base` sudah AA-compliant.

---

## 3. Tipografi

| Peran | Typeface | Alasan |
|---|---|---|
| UI & body text | **IBM Plex Sans** | Karakter institusional-profesional, angka tabular rapi untuk nominal klaim |
| Kode & data teknis | **IBM Plex Mono** | Membedakan data terstruktur (kode ICD-10, ICD-9, No. RM, NIK, kode CBGs) dari teks naratif secara sekilas pandang |

### Type Scale
| Elemen | Ukuran | Weight |
|---|---|---|
| Judul halaman (H1) | 24px | 600 |
| Judul section (H2) | 20px | 600 |
| Judul card (H3) | 16px | 600 |
| Body / isi table | 14px | 400 |
| Label / caption | 13px | 500 |
| Kode (ICD, NIK, No. RM) | 13px, mono | 400 |

Line-height body: 1.5. Hindari judul besar dramatis ala landing page — ini software kerja, hierarki cukup halus.

---

## 4. Layout

```
┌──────────┬─────────────────────────────────────┐
│          │  Topbar: breadcrumb, nama user, unit │
│ Sidebar  ├─────────────────────────────────────┤
│ (menu    │                                       │
│  sesuai  │   Content area (max-width ~1200px,   │
│  role)   │   left-aligned)                       │
│          │                                       │
│          │   Table / Form / Card grid            │
└──────────┴─────────────────────────────────────┘
```

- **Sidebar**: collapsible, isi menu berbeda per role (bukan seragam lalu di-disable).
- **Konten left-aligned**, bukan center — ini software data-dense, bukan halaman marketing.
- **Table**: zebra-striping halus, sticky header, bukan card-list — untuk data pasien/klaim yang jumlahnya banyak.
- **Form rekam medis**: satu section per tab (Anamnesis / Diagnosis / Tindakan), bukan satu form panjang, agar tidak melelahkan dokter/perawat saat input.

### Spacing & Radius
| Token | Nilai |
|---|---|
| `radius` | 6px (konsisten di semua elemen: button, input, card) |
| `spacing-unit` | 4px (kelipatan 4/8/12/16/24) |
| `content-max-width` | 1200px |

---

## 5. Komponen Kunci

### Sidebar Navigasi
- Item aktif: background `primary-dark`, teks putih.
- Item non-aktif: teks `text-secondary`, hover ke `bg-base` sedikit lebih gelap.
- Ikon di kiri tiap label, konsisten set ikon (misal PrimeIcons).

### Tombol
| Jenis | Style |
|---|---|
| Primary | Background `primary`, teks putih, hover `primary-hover` |
| Secondary | Border `border`, teks `text-primary`, background transparan |
| Danger (misal tolak klaim) | Border/teks `#B54245`, background transparan sampai hover |

### Badge Status Klaim
- Bentuk pill kecil, radius penuh, dengan **ikon + teks** (bukan warna saja):
  - Draft: ikon dokumen, abu
  - Pending: ikon jam, amber
  - Disetujui: ikon centang, hijau
  - Ditolak: ikon silang, merah bata

### Input Autocomplete (ICD-10 / ICD-9)
- Debounce 300ms sesuai FR.
- Dropdown hasil menampilkan kode (mono) + nama diagnosis/tindakan berdampingan, kode di-highlight bold.
- State loading: skeleton kecil, bukan spinner besar yang mengganggu ritme ketik.

### Table
- Header sticky, background `surface`, teks `text-secondary` uppercase kecil (13px, bukan tracking berlebihan).
- Baris zebra: selang-seling `surface` dan `bg-base` sangat tipis.
- Kolom kode (NIK, No. RM, ICD) pakai font mono agar mudah dipindai mata.

### Form Rekam Medis
- Tab per section, indikator progres kecil di atas (bukan wizard step besar).
- Validasi inline real-time (misal cek NIK duplikat saat mengetik), bukan alert setelah submit.

### Empty & Error States
- Kosong: pesan singkat + ajakan aksi jelas (misal "Belum ada kunjungan. Tambah kunjungan baru").
- Error: jelaskan apa yang salah dan cara memperbaikinya, nada netral-informatif, tidak meminta maaf berlebihan.

---

## 6. Motion

- Transisi hover: 120–150ms, ease-out, hanya pada elemen interaktif (button, row hover).
- Transisi status klaim berubah: highlight singkat pada badge baru (misal fade+scale 200ms) sebagai konfirmasi visual — satu-satunya tempat motion "diberi panggung".
- Hormati `prefers-reduced-motion`.

---

## 7. Aksesibilitas

- Kontras teks minimal AA (4.5:1 untuk body text).
- Focus state keyboard terlihat jelas (outline 2px `primary` pada elemen fokus).
- Status tidak pernah disampaikan lewat warna saja — selalu disertai ikon/teks.
- Ukuran target klik/tap minimal 40x40px untuk elemen interaktif di table actions.

---

## 8. Referensi Cepat (Design Tokens)

```css
:root {
  --color-primary: #1C3D5A;
  --color-primary-hover: #2D5578;
  --color-primary-dark: #0F2740;

  --color-bg-base: #F7F6F3;
  --color-surface: #FFFFFF;
  --color-border: #E5E2DB;

  --color-text-primary: #1F2937;
  --color-text-secondary: #6B7280;
  --color-text-disabled: #9CA3AF;

  --color-status-draft: #9CA3AF;
  --color-status-pending: #B7791F;
  --color-status-approved: #2F7A4F;
  --color-status-rejected: #B54245;

  --radius-base: 6px;
  --spacing-unit: 4px;
  --content-max-width: 1200px;

  --font-sans: 'IBM Plex Sans', sans-serif;
  --font-mono: 'IBM Plex Mono', monospace;
}
```
