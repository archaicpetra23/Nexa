import { describe, it, expect } from 'vitest'
import client from '../client.js'

// Ini bukan test perilaku HTTP — yang diuji adalah konfigurasi instance Axios,
// karena di situlah bug session biasanya bersembunyi (baseURL salah, kredensial
// tidak ikut terkirim).
describe('api client', () => {
  it('memakai baseURL /api/v1', () => {
    expect(client.defaults.baseURL).toBe('/api/v1')
  })

  it('mengirim cookie kredensial (withCredentials)', () => {
    expect(client.defaults.withCredentials).toBe(true)
  })

  it('adalah instance Axios dengan method HTTP', () => {
    expect(typeof client.get).toBe('function')
    expect(typeof client.post).toBe('function')
    expect(typeof client.put).toBe('function')
    expect(typeof client.delete).toBe('function')
  })

  it('mengirim Content-Type JSON', () => {
    expect(client.defaults.headers?.['Content-Type']).toBe('application/json')
  })
})
