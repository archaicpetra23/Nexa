// Aturan transisi status klaim.
// State machine ini mengikuti validateKlaimStatus di backend
// (backend/internal/usecase/klaim_usecase.go) supaya frontend dan backend sepakat.
// Status disetujui dan ditolak adalah terminal — tidak ada jalan keluar.

export const ALLOWED_TRANSITIONS = Object.freeze({
  draft: ['pending', 'disetujui', 'ditolak'],
  pending: ['disetujui', 'ditolak', 'draft'],
  disetujui: [],
  ditolak: []
})

// canTransition menolak kalau status tidak dikenal, transisi tidak ada di
// daftar, status tujuan sama dengan status awal, atau tujuan butuh alasan
// (pending/ditolak) tapi alasannya kosong.
export function canTransition(current, next, alasan = null) {
  const allowed = ALLOWED_TRANSITIONS[current]
  if (!allowed || !allowed.includes(next)) return false
  if (current === next) return false

  const needsAlasan = next === 'pending' || next === 'ditolak'
  if (needsAlasan && (alasan === null || alasan === undefined || String(alasan).trim() === '')) {
    return false
  }

  return true
}

// transitionError mengembalikan pesan kalau transisi ditolak, atau null kalau
// transisi diizinkan — dipakai UI untuk menonaktifkan opsi yang tidak sah.
export function transitionError(current, next, alasan = null) {
  const allowed = ALLOWED_TRANSITIONS[current]
  if (!allowed) return `Status ${current} tidak dikenal`
  if (current === next) return 'Status baru tidak boleh sama dengan status saat ini'
  if (!allowed.includes(next)) {
    return `Tidak dapat mengubah status ${current} menjadi ${next}`
  }

  const needsAlasan = next === 'pending' || next === 'ditolak'
  if (needsAlasan && (alasan === null || alasan === undefined || String(alasan).trim() === '')) {
    return 'Alasan wajib diisi untuk status pending atau ditolak'
  }

  return null
}

// validateKlaimForm validasi form klaim sebelum submit
export function validateKlaimForm(form) {
  const errors = {}

  if (!form.id_rekam || String(form.id_rekam).trim() === '') {
    errors.id_rekam = 'ID Rekam Medis wajib diisi'
  }

  if (!form.kode_cbg || String(form.kode_cbg).trim() === '') {
    errors.kode_cbg = 'Kode CBGS wajib diisi'
  }

  if (!form.id_petugas_casemix || String(form.id_petugas_casemix).trim() === '') {
    errors.id_petugas_casemix = 'Petugas Casemix wajib dipilih'
  }

  if (!form.status_klaim || String(form.status_klaim).trim() === '') {
    errors.status_klaim = 'Status Klaim wajib dipilih'
  }

  if (form.nominal_klaim === undefined || form.nominal_klaim === null || String(form.nominal_klaim).trim() === '') {
    errors.nominal_klaim = 'Nominal Klaim wajib diisi'
  } else if (Number(form.nominal_klaim) < 0) {
    errors.nominal_klaim = 'Nominal Klaim tidak boleh negatif'
  }

  const needsAlasan = form.status_klaim === 'pending' || form.status_klaim === 'ditolak'
  if (needsAlasan && (!form.alasan_pending_tolak || String(form.alasan_pending_tolak).trim() === '')) {
    errors.alasan_pending_tolak = 'Alasan wajib diisi saat status pending/ditolak'
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors
  }
}
