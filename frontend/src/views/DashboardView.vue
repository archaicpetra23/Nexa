<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-2xl font-bold text-blue-600">Nexa</h1>
          </div>
          <div class="flex items-center">
            <span class="text-gray-700 mr-4">{{ user?.nama }}</span>
            <button 
              @click="handleLogout"
              class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>
    
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-2xl font-bold text-gray-800 mb-4">Dashboard</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="bg-blue-50 p-6 rounded-lg">
              <h3 class="text-lg font-semibold text-blue-800 mb-2">Informasi Pengguna</h3>
              <div class="space-y-2">
                <p><span class="font-medium">Nama:</span> {{ user?.nama }}</p>
                <p><span class="font-medium">Username:</span> {{ user?.username }}</p>
                <p><span class="font-medium">Profesi:</span> {{ user?.profesi }}</p>
                <p><span class="font-medium">Role:</span> <span class="bg-blue-200 px-2 py-1 rounded">{{ user?.role }}</span></p>
              </div>
            </div>
            
            <div class="bg-green-50 p-6 rounded-lg">
              <h3 class="text-lg font-semibold text-green-800 mb-2">Hak Akses</h3>
              <div class="space-y-2">
                <p v-if="user?.role === 'admin_ti'" class="text-sm">✓ Manajemen Pengguna & Unit</p>
                <p v-if="user?.role === 'petugas_rm'" class="text-sm">✓ Registrasi Pasien</p>
                <p v-if="user?.role === 'dokter_dpjp'" class="text-sm">✓ Input Diagnosis (ICD-10)</p>
                <p v-if="user?.role === 'perawat'" class="text-sm">✓ Input Prosedur (ICD-9)</p>
                <p v-if="user?.role === 'petugas_casemix'" class="text-sm">✓ Verifikasi & Grouping CBGs</p>
                <p v-if="user?.role === 'keuangan'" class="text-sm">✓ Manajemen Status Klaim</p>
                <p v-if="user?.role === 'manajemen'" class="text-sm">✓ Dashboard & Laporan Analitik</p>
                <p v-if="user?.role" class="text-sm mt-4 text-gray-600">Akses sesuai role: {{ getRoleDescription(user?.role) }}</p>
              </div>
            </div>
          </div>
          
          <div class="mt-6 p-4 bg-yellow-50 border-l-4 border-yellow-400">
            <p class="text-sm text-yellow-800">
              <strong>Status:</strong> Aplikasi Nexa berhasil berjalan! Backend (Go) dan Frontend (Vue) sudah terintegrasi.
            </p>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

function getRoleDescription(role) {
  const descriptions = {
    'admin_ti': 'Administrator TI - Full system access',
    'petugas_rm': 'Petugas Rekam Medis - Patient registration',
    'dokter_dpjp': 'Dokter DPJP - Clinical records & diagnosis',
    'perawat': 'Perawat - Medical procedures',
    'petugas_casemix': 'Petugas Casemix - Claims verification',
    'keuangan': 'Keuangan - Financial reconciliation',
    'manajemen': 'Manajemen - Analytics & reports'
  }
  return descriptions[role] || role
}
</script>