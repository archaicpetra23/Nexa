import path from 'path'

import { defineConfig, splitVendorChunkPlugin } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue(), splitVendorChunkPlugin()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  server: {
    host: '0.0.0.0',
    port: 5174,
    strictPort: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['vue', 'primevue', 'chart.js', 'axios', 'pinia', 'vue-router'],
          'admin': ['@/views/AdminView.vue'],
          'dashboard': ['@/views/DashboardView.vue'],
          'klaim': ['@/views/KlaimView.vue'],
          'pasien': ['@/views/PasienView.vue'],
          'rekam-medis': ['@/views/RekamMedisView.vue'],
          'master-icd': ['@/views/MasterICDView.vue'],
        },
      },
    },
  }
})