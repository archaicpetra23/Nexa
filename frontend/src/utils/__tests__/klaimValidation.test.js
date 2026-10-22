import { describe, it, expect } from 'vitest'
import {
  ALLOWED_TRANSITIONS,
  canTransition,
  transitionError
} from '../klaimValidation.js'

describe('ALLOWED_TRANSITIONS', () => {
  it('memuat tepat empat status', () => {
    // sort() mengurutkan berdasarkan code point, jadi 'disetujui' < 'ditolak'.
    expect(Object.keys(ALLOWED_TRANSITIONS).sort()).toEqual([
      'disetujui',
      'ditolak',
      'draft',
      'pending'
    ])
  })

  it('bersifat beku supaya tidak bisa diubah saat runtime', () => {
    expect(Object.isFrozen(ALLOWED_TRANSITIONS)).toBe(true)
  })

  it('menandai disetujui dan ditolak sebagai status terminal', () => {
    expect(ALLOWED_TRANSITIONS.disetujui).toEqual([])
    expect(ALLOWED_TRANSITIONS.ditolak).toEqual([])
  })
})

describe('canTransition', () => {
  it('mengizinkan draft menuju pending, disetujui, atau ditolak', () => {
    expect(canTransition('draft', 'pending', 'Menunggu berkas')).toBe(true)
    expect(canTransition('draft', 'disetujui')).toBe(true)
    expect(canTransition('draft', 'ditolak', 'Berkas tidak lengkap')).toBe(true)
  })

  it('mengizinkan pending menuju disetujui, ditolak, atau kembali ke draft', () => {
    expect(canTransition('pending', 'disetujui')).toBe(true)
    expect(canTransition('pending', 'ditolak', 'Koding salah')).toBe(true)
    expect(canTransition('pending', 'draft')).toBe(true)
  })

  it('menolak semua transisi dari status disetujui', () => {
    expect(canTransition('disetujui', 'draft')).toBe(false)
    expect(canTransition('disetujui', 'pending')).toBe(false)
    expect(canTransition('disetujui', 'ditolak', 'alasan')).toBe(false)
  })

  it('menolak semua transisi dari status ditolak', () => {
    expect(canTransition('ditolak', 'pending')).toBe(false)
    expect(canTransition('ditolak', 'disetujui')).toBe(false)
    expect(canTransition('ditolak', 'draft')).toBe(false)
  })

  it('menolak transisi dari status yang tidak dikenal', () => {
    for (const next of ['draft', 'pending', 'disetujui', 'ditolak']) {
      expect(canTransition('arsip', next)).toBe(false)
    }
    expect(canTransition(undefined, 'draft')).toBe(false)
  })

  it('menolak transisi ke status yang sama', () => {
    expect(canTransition('draft', 'draft')).toBe(false)
    expect(canTransition('pending', 'pending', 'alasan')).toBe(false)
  })

  it('menolak pending atau ditolak tanpa alasan', () => {
    expect(canTransition('pending', 'ditolak', '')).toBe(false)
    expect(canTransition('pending', 'ditolak', null)).toBe(false)
    expect(canTransition('pending', 'ditolak')).toBe(false)
    expect(canTransition('pending', 'ditolak', '   ')).toBe(false)
    expect(canTransition('draft', 'pending', '')).toBe(false)
  })

  it('menerima pending atau ditolak dengan alasan terisi', () => {
    expect(canTransition('pending', 'ditolak', 'Berkas tidak lengkap')).toBe(true)
    expect(canTransition('draft', 'pending', 'Menunggu hasil lab')).toBe(true)
  })

  it('tidak mewajibkan alasan untuk status disetujui', () => {
    expect(canTransition('pending', 'disetujui', null)).toBe(true)
    expect(canTransition('pending', 'disetujui', undefined)).toBe(true)
    expect(canTransition('pending', 'disetujui', '')).toBe(true)
  })
})

describe('transitionError', () => {
  it('mengembalikan null saat transisi diizinkan', () => {
    expect(transitionError('draft', 'disetujui')).toBeNull()
    expect(transitionError('pending', 'ditolak', 'Koding salah')).toBeNull()
  })

  it('mengembalikan pesan saat transisi ditolak', () => {
    expect(typeof transitionError('disetujui', 'draft')).toBe('string')
    expect(transitionError('disetujui', 'draft').length).toBeGreaterThan(0)
    expect(transitionError('draft', 'draft').length).toBeGreaterThan(0)
  })

  it('menjelaskan alasan yang hilang', () => {
    const message = transitionError('pending', 'ditolak', '')

    expect(typeof message).toBe('string')
    expect(message).toMatch(/alasan/i)
  })
})
