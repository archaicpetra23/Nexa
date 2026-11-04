import { defineConfig, splitVendorChunkPlugin } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue(), splitVendorChunkPlugin()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    proxy: {
      '/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
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
  },
})