import { describe, it, expect } from 'vitest'
import { validateKlaimForm } from '@/utils/klaimValidation'

describe('validateKlaimForm', () => {
  const validForm = {
    id_rekam: 1,
    kode_cbg: 'A1',
    id_petugas_casemix: 1,
    status_klaim: 'draft',
    nominal_klaim: 1000000,
  }

  it('returns valid=true untuk form lengkap & valid', () => {
    const { valid, errors } = validateKlaimForm(validForm)
    expect(valid).toBe(true)
    expect(Object.keys(errors)).toHaveLength(0)
  })

  it('returns valid=false kalau id_rekam kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, id_rekam: '' })
    expect(valid).toBe(false)
    expect(errors.id_rekam).toBe('ID Rekam Medis wajib diisi')
  })

  it('returns valid=false kalau kode_cbg kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, kode_cbg: '' })
    expect(valid).toBe(false)
    expect(errors.kode_cbg).toBe('Kode CBGS wajib diisi')
  })

  it('returns valid=false kalau id_petugas_casemix kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, id_petugas_casemix: '' })
    expect(valid).toBe(false)
    expect(errors.id_petugas_casemix).toBe('Petugas Casemix wajib dipilih')
  })

  it('returns valid=false kalau status_klaim kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, status_klaim: '' })
    expect(valid).toBe(false)
    expect(errors.status_klaim).toBe('Status Klaim wajib dipilih')
  })

  it('returns valid=false kalau nominal_klaim kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, nominal_klaim: '' })
    expect(valid).toBe(false)
    expect(errors.nominal_klaim).toBe('Nominal Klaim wajib diisi')
  })

  it('returns valid=false kalau nominal_klaim negatif', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, nominal_klaim: -100 })
    expect(valid).toBe(false)
    expect(errors.nominal_klaim).toBe('Nominal Klaim tidak boleh negatif')
  })

  it('returns valid=true kalau status_klaim pending dan alasan diisi', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, status_klaim: 'pending', alasan_pending_tolak: 'berkas kurang' })
    expect(valid).toBe(true)
    expect(errors.alasan_pending_tolak).toBeUndefined()
  })

  it('returns valid=false kalau status_klaim pending tapi alasan kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, status_klaim: 'pending', alasan_pending_tolak: '' })
    expect(valid).toBe(false)
    expect(errors.alasan_pending_tolak).toBe('Alasan wajib diisi saat status pending/ditolak')
  })

  it('returns valid=true kalau status_klaim ditolak dan alasan diisi', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, status_klaim: 'ditolak', alasan_pending_tolak: 'diagnosa tidak sesuai' })
    expect(valid).toBe(true)
  })

  it('returns valid=false kalau status_klaim ditolak tapi alasan kosong', () => {
    const { valid, errors } = validateKlaimForm({ ...validForm, status_klaim: 'ditolak', alasan_pending_tolak: '' })
    expect(valid).toBe(false)
    expect(errors.alasan_pending_tolak).toBe('Alasan wajib diisi saat status pending/ditolak')
  })
})