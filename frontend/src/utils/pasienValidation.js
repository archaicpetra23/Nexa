// Validasi form pasien.
// Aturan di sini adalah ekstraksi 1:1 dari `validate()` inline di PasienView.vue,
// termasuk pesan error dan regex yang sama persis — jangan diperketat tanpa
// menyalahi perilaku form yang sudah ada.

// validatePasienForm memeriksa objek form pasien dan mengembalikan
// { valid, errors } dengan errors berisi pesan per field yang gagal.
export function validatePasienForm(form) {
  const errors = {}

  if (!form.nik) {
    errors.nik = 'NIK wajib diisi'
  } else if (!/^\d{16}$/.test(form.nik)) {
    errors.nik = 'NIK harus 16 digit angka'
  }

  if (!form.nama) errors.nama = 'Nama wajib diisi'
  if (!form.tanggal_lahir) errors.tanggal_lahir = 'Tanggal lahir wajib diisi'
  if (!form.jenis_kelamin) errors.jenis_kelamin = 'Jenis kelamin wajib dipilih'
  if (!form.alamat) errors.alamat = 'Alamat wajib diisi'

  // No BPJS opsional — regex hanya jalan kalau nilainya terisi.
  if (form.no_bpjs && !/^\d{13}$/.test(form.no_bpjs)) {
    errors.no_bpjs = 'No. BPJS harus 13 digit angka'
  }

  return { valid: Object.keys(errors).length === 0, errors }
}
