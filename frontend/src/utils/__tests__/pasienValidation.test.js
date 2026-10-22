import { describe, it, expect } from 'vitest'
import { validatePasienForm } from '../pasienValidation.js'

// Kasus valid penuh dipakai ulang di beberapa test.
const validForm = () => ({
  nik: '2290134500000000',
  nama: 'Tuti Saputra',
  tanggal_lahir: new Date('1995-05-15'),
  jenis_kelamin: 'L',
  alamat: 'Jl. Melati No. 12, Tangerang',
  no_bpjs: '3117600000000'
})

describe('validatePasienForm', () => {
  it('menerima form yang lengkap dan valid', () => {
    const result = validatePasienForm(validForm())

    expect(result.valid).toBe(true)
    expect(result.errors).toEqual({})
  })

  it('menolak NIK kosong', () => {
    const form = { ...validForm(), nik: '' }
    const result = validatePasienForm(form)

    expect(result.valid).toBe(false)
    expect(result.errors.nik).toBe('NIK wajib diisi')
  })

  it('menolak NIK yang bukan 16 digit angka', () => {
    const result = validatePasienForm({ ...validForm(), nik: '12345' })

    expect(result.valid).toBe(false)
    expect(result.errors.nik).toBe('NIK harus 16 digit angka')
  })

  it('menolak NIK yang mengandung huruf', () => {
    const result = validatePasienForm({ ...validForm(), nik: '1234abcd56789012' })

    expect(result.valid).toBe(false)
    expect(result.errors.nik).toBe('NIK harus 16 digit angka')
  })

  it('menerima 16 digit angka sebagai NIK yang sah', () => {
    const result = validatePasienForm({ ...validForm(), nik: '1122334455667788' })

    expect(result.errors.nik).toBeUndefined()
  })

  it('menerima nama kosong sebagai error terpisah', () => {
    const result = validatePasienForm({ ...validForm(), nama: '' })

    expect(result.valid).toBe(false)
    expect(result.errors.nama).toBe('Nama wajib diisi')
  })

  it('menerima No BPJS kosong karena bersifat opsional', () => {
    const result = validatePasienForm({ ...validForm(), no_bpjs: '' })

    expect(result.valid).toBe(true)
    expect(result.errors.no_bpjs).toBeUndefined()
  })

  it('menerima No BPJS null karena bersifat opsional', () => {
    const result = validatePasienForm({ ...validForm(), no_bpjs: null })

    expect(result.valid).toBe(true)
    expect(result.errors.no_bpjs).toBeUndefined()
  })

  it('menolak No BPJS yang bukan 13 digit', () => {
    const result = validatePasienForm({ ...validForm(), no_bpjs: '123' })

    expect(result.valid).toBe(false)
    expect(result.errors.no_bpjs).toBe('No. BPJS harus 13 digit angka')
  })

  it('menerima No BPJS 13 digit', () => {
    const result = validatePasienForm({ ...validForm(), no_bpjs: '0001234567890' })

    expect(result.errors.no_bpjs).toBeUndefined()
  })

  it('menolak jenis kelamin kosong', () => {
    const result = validatePasienForm({ ...validForm(), jenis_kelamin: '' })

    expect(result.valid).toBe(false)
    expect(result.errors.jenis_kelamin).toBe('Jenis kelamin wajib dipilih')
  })

  it('menerima jenis kelamin L dan P', () => {
    for (const jenisKelamin of ['L', 'P']) {
      const result = validatePasienForm({ ...validForm(), jenis_kelamin: jenisKelamin })
      expect(result.errors.jenis_kelamin).toBeUndefined()
    }
  })

  it('menolak tanggal lahir kosong', () => {
    const result = validatePasienForm({ ...validForm(), tanggal_lahir: null })

    expect(result.valid).toBe(false)
    expect(result.errors.tanggal_lahir).toBe('Tanggal lahir wajib diisi')
  })

  it('menolak alamat kosong', () => {
    const result = validatePasienForm({ ...validForm(), alamat: '' })

    expect(result.valid).toBe(false)
    expect(result.errors.alamat).toBe('Alamat wajib diisi')
  })

  it('melaporkan beberapa error sekaligus', () => {
    const result = validatePasienForm({
      nik: 'salah',
      nama: '',
      tanggal_lahir: null,
      jenis_kelamin: '',
      alamat: '',
      no_bpjs: '12'
    })

    expect(result.valid).toBe(false)
    expect(Object.keys(result.errors)).toEqual([
      'nik',
      'nama',
      'tanggal_lahir',
      'jenis_kelamin',
      'alamat',
      'no_bpjs'
    ])
  })
})
