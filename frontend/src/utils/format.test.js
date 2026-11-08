import { describe, it, expect } from 'vitest'
import {
  formatRupiah,
  formatTanggal,
  formatTanggalSingkat,
  formatTanggalWaktu,
} from '@/utils/format'

describe('formatRupiah', () => {
  it('formats number to IDR currency', () => {
    expect(formatRupiah(1000000)).toBe('Rp\u00a01.000.000')
  })

  it('formats string number to IDR', () => {
    expect(formatRupiah('500000')).toBe('Rp\u00a0500.000')
  })

  it('returns Rp 0 for null', () => {
    expect(formatRupiah(null)).toBe('Rp\u00a00')
  })

  it('returns Rp 0 for undefined', () => {
    expect(formatRupiah(undefined)).toBe('Rp\u00a00')
  })

  it('returns Rp 0 for empty string', () => {
    expect(formatRupiah('')).toBe('Rp\u00a00')
  })

  it('returns Rp 0 for NaN', () => {
    expect(formatRupiah(NaN)).toBe('Rp\u00a00')
  })

  it('handles large numbers', () => {
    expect(formatRupiah(1234567890)).toBe('Rp\u00a01.234.567.890')
  })

  it('handles zero', () => {
    expect(formatRupiah(0)).toBe('Rp\u00a00')
  })
})

describe('formatTanggal', () => {
  it('formats ISO date to DD MMM YYYY', () => {
    expect(formatTanggal('2026-10-22')).toBe('22 Okt 2026')
  })

  it('formats ISO datetime to DD MMM YYYY', () => {
    expect(formatTanggal('2026-10-22T08:30:00Z')).toBe('22 Okt 2026')
  })

  it('returns dash for empty string', () => {
    expect(formatTanggal('')).toBe('-')
  })

  it('returns dash for null', () => {
    expect(formatTanggal(null)).toBe('-')
  })

  it('returns dash for undefined', () => {
    expect(formatTanggal(undefined)).toBe('-')
  })

  it('returns dash for invalid date string', () => {
    expect(formatTanggal('not-a-date')).toBe('-')
  })
})

describe('formatTanggalSingkat', () => {
  it('formats ISO date to DD/MM/YYYY', () => {
    expect(formatTanggalSingkat('2026-10-22')).toBe('22/10/2026')
  })

  it('formats ISO datetime to DD/MM/YYYY', () => {
    expect(formatTanggalSingkat('2026-10-22T08:30:00Z')).toBe('22/10/2026')
  })

  it('returns dash for empty string', () => {
    expect(formatTanggalSingkat('')).toBe('-')
  })

  it('returns dash for null', () => {
    expect(formatTanggalSingkat(null)).toBe('-')
  })

  it('returns dash for undefined', () => {
    expect(formatTanggalSingkat(undefined)).toBe('-')
  })

  it('returns dash for invalid date string', () => {
    expect(formatTanggalSingkat('not-a-date')).toBe('-')
  })
})

describe('formatTanggalWaktu', () => {
  it('formats Date to DD/MM/YYYY, HH.MM', () => {
    const date = new Date('2026-10-22T08:30:00')
    expect(formatTanggalWaktu(date)).toBe('22/10/2026, 08.30')
  })

  it('formats ISO string to DD/MM/YYYY, HH.MM', () => {
    expect(formatTanggalWaktu('2026-10-22T08:30:00')).toBe('22/10/2026, 08.30')
  })

  it('returns dash for empty string', () => {
    expect(formatTanggalWaktu('')).toBe('-')
  })

  it('returns dash for null', () => {
    expect(formatTanggalWaktu(null)).toBe('-')
  })

  it('returns dash for undefined', () => {
    expect(formatTanggalWaktu(undefined)).toBe('-')
  })

  it('returns dash for invalid date', () => {
    expect(formatTanggalWaktu('not-a-date')).toBe('-')
  })
})