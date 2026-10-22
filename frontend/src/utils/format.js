// Formatter angka dan tanggal yang dipakai lintas view.
// Sebelumnya formatRupiah disalin inline di 4 view (KlaimView, DashboardView,
// RekamMedisView, MasterICDView) dan formatTanggalWaktu di AdminView —
// sekarang satu sumber supaya bisa diuji dan tidak ada yang tertinggal.

// formatRupiah memakai Intl.NumberFormat id-ID dengan IDR, tanpa desimal.
// Nilai kosong (null/undefined/""/NaN) diformat jadi "Rp 0" supaya tabel tidak
// menampilkan "RpNaN" saat data belum termuat.
//
// Catatan: pembatas "Rp " adalah NON-BREAKING SPACE (U+00A0), bukan spasi biasa,
// supaya hasilnya identik dengan output Intl di semua jalur kode.
export function formatRupiah(value) {
  const numeric = Number(value)
  if (!Number.isFinite(numeric)) return 'Rp\u00a00'

  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(numeric)
}

// formatTanggal mengembalikan "22 Okt 2026" (tahun penuh).
// timeZone: 'UTC' dipakai supaya hasil tidak bergantung pada timezone test
// runner — kalau tidak, assertion akan pecah di mesin dengan zona berbeda.
// Tahun dipaksa 'numeric' supaya konsisten di semua view; sebelumnya ada view
// yang memakai '2-digit' sehingga tahun tampil "26".
export function formatTanggal(isoString) {
  if (!isoString) return '-'

  const date = new Date(isoString)
  if (Number.isNaN(date.getTime())) return '-'

  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    timeZone: 'UTC'
  }).format(date)
}

// formatTanggalSingkat mengembalikan "22/10/2026" — format yang dipakai tabel
// daftar (KlaimView, RekamMedisView, PasienView) sebelum helper ini ada.
export function formatTanggalSingkat(isoString) {
  if (!isoString) return '-'

  const date = new Date(isoString)
  if (Number.isNaN(date.getTime())) return '-'

  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).format(date)
}

// formatTanggalWaktu mengembalikan "22/10/2026 08:58" — sama dengan helper
// inline yang dulu ada di AdminView.vue.
export function formatTanggalWaktu(date) {
  if (!date) return '-'

  const d = date instanceof Date ? date : new Date(date)
  if (Number.isNaN(d.getTime())) return '-'

  return d.toLocaleString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
