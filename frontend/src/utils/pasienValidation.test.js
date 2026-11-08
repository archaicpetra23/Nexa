import { describe, it, expect } from 'vitest'
import { validatePasienForm } from '@/utils/pasienValidation'

describe('validatePasienForm', () => {
  const validForm = {
    nik: '1234567890123456',
    nama: 'Budi Santoso',
    tanggal_lahir: '1990-05-15',
    jenis_kelamin: 'L',
    alamat: 'Jl. Raya No. 123',
    no_bpjs: '1234567890123',
  }

  it('returns valid=true untuk form lengkap & valid', () => {
    const { valid, errors } = validatePasienForm(validForm)
    expect(valid).toBe(true)
    expect(Object.keys(errors)).toHaveLength(0)
  })

  it('returns valid=false kalau NIK kosong', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, nik: '' })
    expect(valid).toBe(false)
    expect(errors.nik).toBe('NIK wajib diisi')
  })

  it('returns valid=false kalau NIK bukan 16 digit', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, nik: '123' })
    expect(valid).toBe(false)
    expect(errors.nik).toBe('NIK harus 16 digit angka')
  })

  it('returns valid=false kalau NIK mengandung non-digit', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, nik: '123456789012345a' })
    expect(valid).toBe(false)
    expect(errors.nik).toBe('NIK harus 16 digit angka')
  })

  it('returns valid=false kalau nama kosong', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, nama: '' })
    expect(valid).toBe(false)
    expect(errors.nama).toBe('Nama wajib diisi')
  })

  it('returns valid=false kalau tanggal_lahir kosong', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, tanggal_lahir: '' })
    expect(valid).toBe(false)
    expect(errors.tanggal_lahir).toBe('Tanggal lahir wajib diisi')
  })

  it('returns valid=false kalau jenis_kelamin kosong', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, jenis_kelamin: '' })
    expect(valid).toBe(false)
    expect(errors.jenis_kelamin).toBe('Jenis kelamin wajib dipilih')
  })

  it('returns valid=false kalau alamat kosong', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, alamat: '' })
    expect(valid).toBe(false)
    expect(errors.alamat).toBe('Alamat wajib diisi')
  })

  it('returns valid=true kalau no_bpjs kosong (opsional)', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, no_bpjs: '' })
    expect(valid).toBe(true)
    expect(errors.no_bpjs).toBeUndefined()
  })

  it('returns valid=false kalau no_bpjs diisi tapi bukan 13 digit', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, no_bpjs: '123' })
    expect(valid).toBe(false)
    expect(errors.no_bpjs).toBe('No. BPJS harus 13 digit angka')
  })

  it('returns valid=false kalau no_bpjs 13 digit tapi mengandung huruf', () => {
    const { valid, errors } = validatePasienForm({ ...validForm, no_bpjs: '123456789012a' })
    expect(valid).toBe(false)
    expect(errors.no_bpjs).toBe('No. BPJS harus 13 digit angka')
  })
})