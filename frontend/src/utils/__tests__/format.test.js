import { describe, it, expect } from 'vitest'
import { formatRupiah, formatTanggal, formatTanggalSingkat, formatTanggalWaktu } from '../format.js'

describe('formatRupiah', () => {
  it('memformat angka besar dengan pemisah ribuan titik', () => {
    expect(formatRupiah(147091183232.21338)).toBe('Rp 147.091.183.232')
  })

  it('memformat nol', () => {
    expect(formatRupiah(0)).toBe('Rp 0')
  })

  it('membulatkan, bukan memotong', () => {
    expect(formatRupiah(1500.75)).toBe('Rp 1.501')
    expect(formatRupiah(1500.25)).toBe('Rp 1.500')
  })

  it('meletakkan tanda minus di depan mata uang', () => {
    expect(formatRupiah(-5000)).toBe('-Rp 5.000')
  })

  it('memformat nilai kosong menjadi Rp0, bukan RpNaN', () => {
    expect(formatRupiah(null)).toBe('Rp 0')
    expect(formatRupiah(undefined)).toBe('Rp 0')
    expect(formatRupiah('')).toBe('Rp 0')
    expect(formatRupiah(NaN)).toBe('Rp 0')
  })

  it('menerima string angka', () => {
    expect(formatRupiah('2500')).toBe('Rp 2.500')
  })

  it('selalu memakai rupiah dan pemisah ribuan titik untuk nilai besar', () => {
    const formatted = formatRupiah(147091183232)

    expect(formatted).toContain('Rp')
    expect(formatted).toContain('147.091.183.232')
    expect(formatted).not.toContain(',')
  })
})

describe('formatTanggal', () => {
  it('memformat ISO string ke tanggal Indonesia', () => {
    expect(formatTanggal('2026-10-22T08:58:16.000Z')).toBe('22 Okt 2026')
  })

  it('hasilnya tidak bergantung pada timezone runner (dinormalisasi ke UTC)', () => {
    // Tanpa timeZone: 'UTC', jam 08:58 UTC bisa jadi 21 Okt di zona negatif.
    expect(formatTanggal('2026-10-22T08:58:16.000Z')).toBe('22 Okt 2026')
  })

  it('mengembalikan strip untuk nilai kosong', () => {
    expect(formatTanggal(null)).toBe('-')
    expect(formatTanggal('')).toBe('-')
    expect(formatTanggal(undefined)).toBe('-')
  })

  it('mengembalikan strip untuk tanggal tidak valid', () => {
    expect(formatTanggal('bukan-tanggal')).toBe('-')
  })
})

describe('formatTanggalSingkat', () => {
  it('memformat ke tanggal tabel Indonesia dengan tahun penuh', () => {
    expect(formatTanggalSingkat('2026-10-22T08:58:16.000Z')).toBe('22/10/2026')
  })

  it('mengembalikan strip untuk nilai kosong', () => {
    expect(formatTanggalSingkat(null)).toBe('-')
    expect(formatTanggalSingkat('')).toBe('-')
  })

  it('mengembalikan strip untuk tanggal tidak valid', () => {
    expect(formatTanggalSingkat('bukan-tanggal')).toBe('-')
  })
})

describe('formatTanggalWaktu', () => {
  it('memformat ke tanggal dan jam Indonesia', () => {
    expect(formatTanggalWaktu(new Date('2026-10-22T08:58:16Z'))).toContain('2026')
  })

  it('menerima string tanggal', () => {
    expect(formatTanggalWaktu('2026-10-22T08:58:16Z')).toContain('2026')
  })

  it('mengembalikan strip untuk nilai kosong', () => {
    expect(formatTanggalWaktu(null)).toBe('-')
    expect(formatTanggalWaktu('')).toBe('-')
    expect(formatTanggalWaktu(undefined)).toBe('-')
  })

  it('mengembalikan strip untuk tanggal tidak valid', () => {
    expect(formatTanggalWaktu('bukan-tanggal')).toBe('-')
  })
})
