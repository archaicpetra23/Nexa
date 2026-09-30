<template>
  <div class="min-h-screen bg-base">
    <header class="bg-primary text-white py-4 px-8">
      <div class="max-w-screen-xl mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-semibold">Nexa</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm">{{ authStore.user?.nama }}</span>
          <span class="text-sm opacity-75">{{ authStore.user?.role }}</span>
          <button @click="handleLogout" class="px-4 py-2 bg-white text-primary rounded text-sm hover:bg-gray-100 transition-colors">
            Logout
          </button>
        </div>
      </div>
    </header>

    <main class="max-w-screen-xl mx-auto px-8 py-10">
      <h1 class="text-2xl font-semibold text-primary mb-8">Dashboard</h1>

      <div class="grid md:grid-cols-3 gap-6 mb-10">
        <div class="bg-surface p-6 rounded border border-border">
          <p class="text-secondary text-sm mb-2">Role</p>
          <p class="text-xl font-semibold text-primary">{{ authStore.user?.role }}</p>
        </div>
        <div class="bg-surface p-6 rounded border border-border">
          <p class="text-secondary text-sm mb-2">Profesi</p>
          <p class="text-xl font-semibold text-primary">{{ authStore.user?.profesi }}</p>
        </div>
        <div class="bg-surface p-6 rounded border border-border">
          <p class="text-secondary text-sm mb-2">Username</p>
          <p class="text-xl font-semibold text-primary font-mono">{{ authStore.user?.username }}</p>
        </div>
      </div>

      <div class="bg-surface p-8 rounded border border-border">
        <h2 class="text-lg font-semibold text-primary mb-4">Modul Tersedia</h2>
        <p class="text-secondary leading-relaxed">
          Dashboard menampilkan ringkasan KPI operasional dan grafik visualisasi klaim.
          Fitur halaman modul (Pasien, Rekam Medis, Klaim) sedang dalam tahap implementasi Sprint 3-4.
        </p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
