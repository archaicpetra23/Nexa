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
