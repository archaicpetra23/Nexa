# Coding Style Guide — Nexa (Casemix Management System)

Dokumen ini menjadi rujukan gaya penulisan kode untuk seluruh kontributor proyek **Nexa**, mencakup Backend (Go/Gin/GORM) dan Frontend (Vue 3/PrimeVue/Pinia). Tujuannya: konsistensi, keamanan (OWASP Top 10), dan kemudahan maintenance jangka panjang oleh tim RSUD.

---

## 0. Prinsip Umum

1. **Konsistensi > Preferensi Pribadi.** Ikuti pola yang sudah ada di codebase, bukan gaya pribadi.
2. **Explicit > Implicit.** Hindari "magic value"; gunakan konstanta bernama untuk status (`draft`, `pending`, `disetujui`, `ditolak`, role, dll).
3. **Fail Safe, Fail Loud.** Error tidak boleh ditelan diam-diam (`_ = err` dilarang kecuali dikomentari alasannya).
4. **Security by Default.** Setiap endpoint baru wajib melewati checklist RBAC + IDOR (lihat §3.4) sebelum merge.
5. Semua kode, komentar teknis, nama variabel/fungsi ditulis **Bahasa Inggris**. Istilah domain medis/BPJS yang sudah baku (`icd10`, `cbgs`, `dpjp`, `bpjs`) boleh dipertahankan apa adanya, mengikuti penamaan tabel di §7.

---

## 1. Struktur Direktori (mengikuti §13 PRD)

```
nexa-casemix/
├── backend/
│   └── internal/
│       ├── config/        # env & db config
│       ├── middleware/    # jwt, rbac, rate limiter, cors
│       ├── model/         # GORM entities & DTOs
│       ├── repository/    # akses DB murni (tanpa business logic)
│       ├── usecase/       # business logic + validasi IDOR
│       └── delivery/http/ # Gin handlers & routes
└── frontend/
    └── src/
        ├── components/    # UI reusable, tanpa API call langsung
        ├── views/         # halaman, boleh panggil store/service
        ├── stores/        # Pinia (authStore, claimStore, dst)
        ├── services/      # axios client + interceptors
        └── router/        # route + navigation guards
```

**Aturan lapisan (backend):** `delivery` → `usecase` → `repository` → `model`. Handler **tidak boleh** memanggil `repository` langsung, dan `repository` **tidak boleh** berisi business rule (mis. validasi `alasan_pending_tolak` wajib diisi ada di usecase, bukan repository).

---

## 2. Backend — Go (Gin + GORM)

### 2.1 Format & Linting
- Wajib `gofmt` / `goimports` sebelum commit.
- Gunakan `golangci-lint` (`golint`, `govet`, `errcheck`, `staticcheck` aktif).
- Panjang baris disarankan < 120 karakter.

### 2.2 Penamaan
| Elemen | Konvensi | Contoh |
| :--- | :--- | :--- |
| Package | huruf kecil, singkat, tanpa underscore | `usecase`, `middleware` |
| Struct (exported) | PascalCase | `RekamMedis`, `KlaimUsecase` |
| Interface | PascalCase, akhiran deskriptif bukan `-er` paksa | `PasienRepository` |
| Fungsi/Method exported | PascalCase, verb di depan | `CreatePasien`, `ValidateNIK` |
| Fungsi privat | camelCase | `hashPassword`, `buildFilter` |
| Konstanta status | PascalCase group / UPPER_SNAKE jika enum string DB | `StatusDraft = "draft"` |
| DTO Request/Response | akhiran jelas | `CreatePasienRequest`, `KlaimResponse` |
| File | snake_case, sufiks layer | `pasien_usecase.go`, `pasien_repository.go`, `pasien_handler.go` |

### 2.3 Struct Model vs DTO
- **Model** (`internal/model`) = representasi tabel (GORM tags), **tidak** diekspos langsung ke JSON response bila mengandung field sensitif (mis. `password_hash`).
- **DTO** dipakai di layer `delivery` untuk request/response. Jangan `json:"-"` sebagai satu-satunya pengaman — gunakan DTO terpisah.

```go
type Pasien struct {
    IDPasien   uint   `gorm:"primaryKey;column:id_pasien"`
    NIK        string `gorm:"column:nik;uniqueIndex"`
    // ...
}

type PasienResponse struct {
    IDPasien uint   `json:"id_pasien"`
    NIK      string `json:"nik"`
    Nama     string `json:"nama"`
}
```

### 2.4 Error Handling
- Selalu `if err != nil { return err }` — jangan diabaikan.
- Gunakan **sentinel error** / custom error type per domain agar handler bisa memetakan ke HTTP status yang tepat, bukan `errors.New` string bebas di banyak tempat.

```go
var (
    ErrNotFound      = errors.New("data tidak ditemukan")
    ErrForbidden     = errors.New("akses ditolak")
    ErrDuplicateData = errors.New("data duplikat")
)
```

- Response error **selalu** mengikuti format standar di §8 (`success: false`), tanpa membocorkan stack trace / detail SQL ke klien (lihat NFR & OWASP A05 di PRD).

### 2.5 Transaksi & Integritas Data
- Operasi yang menyentuh >1 tabel (mis. buat klaim + update status + tulis `log_aktivitas`) **wajib** dibungkus transaksi GORM (`Begin`/`Commit`/`Rollback`), sesuai NFR-03.
- Soft delete **wajib** via `deleted_at` (GORM default `gorm.DeletedAt`). **Dilarang** memanggil `Unscoped().Delete()` pada tabel transaksional (rekam medis, klaim, pasien).

### 2.6 Audit Trail
- Setiap usecase yang melakukan `INSERT`/`UPDATE`/`DELETE` pada data sensitif wajib memanggil helper pencatat log (mis. `auditlog.Record(ctx, userID, "klaim", "UPDATE", idKlaim)`) — jangan duplikasi logic pencatatan log di banyak tempat, sentralisasikan di satu helper/interceptor (§FR-08).

### 2.7 Keamanan (wajib per endpoint baru)
1. Middleware `RBAC(allowedRoles...)` terpasang di route.
2. Di usecase, validasi kepemilikan data sebelum mutasi (contoh persis dari PRD §11.1):
   ```go
   if rekam.IDDokter != currentUser.ID && !currentUser.IsAdmin {
       return ErrForbidden
   }
   ```
3. Password **hanya** lewat `bcrypt`, cost ≥ 10. Jangan pernah log/return password_hash.
4. Query dinamis **wajib** parameterized (GORM `Where("nik = ?", nik)`), jangan concat string SQL.
5. Rate limiter aktif untuk endpoint auth.

### 2.8 Komentar & Dokumentasi
- Setiap fungsi exported diberi komentar Go-doc standar (`// CreateKlaim creates ...`).
- Business rule non-trivial (mis. rumus LoS, validasi `chk_alasan_pending_tolak`) diberi komentar penjelas rumus/alasan, bukan menerjemahkan kode baris demi baris.

---

## 3. Frontend — Vue 3 (Composition API)

### 3.1 Format & Linting
- ESLint (`eslint-plugin-vue`) + Prettier, jalankan `lint-staged` pada pre-commit.
- **Wajib** `<script setup>` + Composition API. Hindari Options API pada komponen baru.

### 3.2 Penamaan
| Elemen | Konvensi | Contoh |
| :--- | :--- | :--- |
| Nama file komponen | PascalCase | `PasienForm.vue`, `KlaimStatusTag.vue` |
| Komponen di template | PascalCase | `<PasienForm />` |
| Composable | camelCase, prefiks `use` | `useDebounceSearch.js` |
| Pinia store | camelCase, akhiran `Store` | `authStore.js`, `claimStore.js` |
| Props/emits | camelCase | `patientId`, `emit('update:status')` |
| Variabel reaktif | camelCase, deskriptif | `isLoading`, `selectedDiagnosis` |

### 3.3 Struktur Komponen
Urutan blok dalam `<script setup>`: imports → props/emits → state (`ref`/`reactive`) → computed → watchers → lifecycle hooks → fungsi.

```vue
<script setup>
import { ref, computed } from 'vue'
import { useClaimStore } from '@/stores/claimStore'

const props = defineProps({ claimId: { type: Number, required: true } })
const emit = defineEmits(['statusChanged'])

const claimStore = useClaimStore()
const isSubmitting = ref(false)

const statusLabel = computed(() => claimStore.getStatusLabel(props.claimId))

async function submitStatusChange(payload) {
  isSubmitting.value = true
  try {
    await claimStore.updateStatus(props.claimId, payload)
    emit('statusChanged', payload.status_klaim)
  } finally {
    isSubmitting.value = false
  }
}
</script>
```

### 3.4 State & API Call
- **Komponen tidak boleh memanggil axios langsung.** Semua request lewat `services/` (axios instance dengan interceptor `withCredentials: true` untuk cookie JWT) dan/atau Pinia store.
- Store bertanggung jawab atas caching/state global (mis. `authStore.user`, `claimStore.list`); komponen hanya konsumsi state via `computed`/`storeToRefs`.

### 3.5 Autocomplete & Debounce
- Semua input pencarian ICD-10/ICD-9 **wajib** pakai composable `useDebounceSearch` dengan delay **300ms** (konsisten dengan FR/NFR di PRD), jangan implementasi debounce ad-hoc di tiap komponen.

### 3.6 Role-Based Rendering
- Jangan sembunyikan menu hanya via CSS `v-if` di sisi client sebagai satu-satunya proteksi — ini **UX only**, otorisasi sebenarnya tetap divalidasi backend. Gunakan `authStore.role` + `router` navigation guard untuk mencegah akses route.

### 3.7 Styling
- Utility-first dengan **Tailwind**; komponen kompleks (`DataTable`, `AutoComplete`, `Dialog`, `Tag`) pakai **PrimeVue**.
- Warna status klaim mengikuti mapping baku (§12 PRD) — definisikan sekali di satu util/constant (`statusColors.js`), jangan hardcode warna berulang di tiap komponen.

---

## 4. Format Response API (kontrak wajib)

Semua endpoint backend mengembalikan bentuk seragam ini (§8 PRD):

```json
{
  "success": true,
  "message": "Deskripsi status operasi",
  "data": {}
}
```

- `success: false` untuk error, `data` boleh `null`, `message` **harus** ramah pengguna (bukan pesan error Go/SQL mentah).
- Untuk list dengan pagination, `data` berisi `{ "items": [...], "total": 0, "page": 1, "limit": 10 }`.

---

## 5. Konvensi Endpoint REST

- Penamaan path: `kebab-case`, plural noun, versi eksplisit (`/api/v1/rekam-medis`).
- Method HTTP sesuai semantik CRUD: `GET` (baca), `POST` (buat), `PUT` (update penuh), `PATCH` (update sebagian/transisi status), `DELETE` (soft delete).
- Filter & pagination via query string (`?search=&page=&limit=`), bukan body.

---

## 6. Konvensi Database (PostgreSQL)

Mengikuti DDL §7 PRD — **jangan menyimpang** dari pola berikut untuk tabel baru:

- Nama tabel & kolom: `snake_case`, Bahasa Indonesia sesuai domain existing (`id_pasien`, `tanggal_kunjungan`).
- Primary key: `id_<nama_tabel_singular>` (`id_pasien`, `id_klaim`), `SERIAL`.
- Timestamp wajib: `created_at`, `updated_at` (`TIMESTAMPTZ DEFAULT NOW()`), `deleted_at` (nullable, soft delete).
- Constraint check diberi nama eksplisit: `chk_<deskripsi>` (`chk_alasan_pending_tolak`).
- Unique constraint gabungan: `uq_<tabel1>_<tabel2>` (`uq_rekam_icd10`).
- Index performa: `idx_<tabel>_<kolom>`; index trigram: `idx_trgm_<tabel>_<kolom>`.
- Foreign key: default `ON DELETE RESTRICT` untuk data transaksional (jangan `CASCADE` sembarangan pada `pasien`, `klaim`, `rekam_medis`).

---

## 7. Git Commit Convention

Gunakan **Conventional Commits**:

```
<type>(<scope>): <deskripsi singkat>

feat(klaim): tambah endpoint transisi status klaim
fix(auth): perbaiki validasi cookie httponly
refactor(pasien): pisahkan validasi NIK ke usecase
docs(readme): update instruksi docker compose
```

Tipe yang dipakai: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `security`.

---

## 8. Testing

- Backend: unit test per `usecase` (business rule, mis. validasi `chk_alasan_pending_tolak`, cek IDOR) menggunakan `testing` + `testify`; gunakan mock untuk `repository`.
- Frontend: minimal test komponen kritikal (form pasien, transisi status klaim) dengan Vitest.
- Target sebelum PR merge: alur "happy path" + minimal 1 skenario penolakan/otorisasi (403/422) per endpoint baru.

---

## 9. Checklist Sebelum Pull Request

- [ ] `gofmt`/ESLint bersih, tanpa warning lint.
- [ ] Endpoint baru sudah dicek RBAC middleware + validasi kepemilikan (IDOR).
- [ ] Operasi multi-tabel dibungkus transaksi.
- [ ] Mutasi data sensitif tercatat ke `log_aktivitas`.
- [ ] Tidak ada `console.log` / `fmt.Println` debug tersisa.
- [ ] Response API mengikuti format standar §4.
- [ ] Tidak ada credential/secret hardcoded (gunakan `.env`).
