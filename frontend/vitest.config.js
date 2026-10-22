import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// Konfigurasi Vitest. Plugin vue dipakai supaya file .vue bisa diimpor langsung
// oleh test kalau nanti ada test komponen.
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) }
  },
  test: {
    environment: 'jsdom',
    globals: false,
    include: ['src/**/*.test.js'],
    css: false
  }
})
