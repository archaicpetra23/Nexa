<template>
  <header class="app-topbar">
    <div class="topbar-start">
      <button
        class="menu-toggle"
        type="button"
        :aria-label="sidebarExpanded ? 'Sembunyikan menu navigasi' : 'Tampilkan menu navigasi'"
        :aria-expanded="sidebarExpanded"
        aria-controls="app-sidebar"
        @click="emit('toggle-sidebar')"
      >
        <i :class="isMobile ? 'pi pi-bars' : sidebarExpanded ? 'pi pi-angle-left' : 'pi pi-angle-right'" aria-hidden="true"></i>
      </button>
      <nav class="topbar-breadcrumb" aria-label="Breadcrumb">
        <span class="breadcrumb-root">Nexa</span>
        <i class="pi pi-angle-right" aria-hidden="true"></i>
        <span class="breadcrumb-current">{{ currentPageLabel }}</span>
      </nav>
    </div>

    <div class="topbar-user">
      <span class="topbar-user-copy">
        <strong>{{ authStore.user?.nama || authStore.user?.username || 'Pengguna' }}</strong>
        <small>{{ authStore.user?.unit_nama || roleLabel }}</small>
      </span>
      <span class="topbar-avatar">{{ userInitials }}</span>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

defineProps({
  sidebarExpanded: { type: Boolean, default: true },
  isMobile: { type: Boolean, default: false }
})

const emit = defineEmits(['toggle-sidebar'])
const route = useRoute()
const authStore = useAuthStore()

const pageLabels = {
  '/dashboard': 'Dashboard',
  '/admin': 'Panel admin'
}

const roleLabels = {
  admin_ti: 'Administrator TI',
  petugas_rm: 'Petugas Rekam Medis',
  dokter_dpjp: 'Dokter DPJP',
  petugas_casemix: 'Petugas Casemix',
  keuangan: 'Keuangan',
  manajemen: 'Manajemen'
}

const currentPageLabel = computed(() => pageLabels[route.path] || 'Dashboard')
const roleLabel = computed(() => roleLabels[authStore.user?.role] || 'Staf rumah sakit')
const userInitials = computed(() => {
  const name = authStore.user?.nama || authStore.user?.username || 'N'
  return name.trim().split(/\s+/).slice(0, 2).map((part) => part.charAt(0).toUpperCase()).join('')
})
</script>

<style scoped>
.app-topbar {
  position: sticky;
  z-index: 10;
  top: 0;
  display: flex;
  min-height: 66px;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 0 30px;
  border-bottom: 1px solid var(--color-border);
  background: rgb(255 255 255 / 96%);
}

.topbar-start,
.topbar-breadcrumb,
.topbar-user {
  display: flex;
  align-items: center;
}

.topbar-start {
  gap: 15px;
}

.topbar-breadcrumb {
  gap: 10px;
  font-size: 11px;
}

.breadcrumb-root {
  color: var(--color-text-secondary);
}

.topbar-breadcrumb > i {
  color: var(--color-text-disabled);
  font-size: 10px;
}

.breadcrumb-current {
  color: var(--color-primary-dark);
  font-weight: 600;
}

.topbar-user {
  gap: 11px;
}

.topbar-user-copy {
  display: grid;
  gap: 2px;
  text-align: right;
}

.topbar-user-copy strong {
  color: var(--color-text-primary);
  font-size: 11px;
  font-weight: 600;
}

.topbar-user-copy small {
  color: var(--color-text-secondary);
  font-size: 9px;
}

.topbar-avatar {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border: 1px solid #D8E0E6;
  border-radius: 50%;
  background: #EEF2F5;
  color: var(--color-primary);
  font-size: 10px;
  font-weight: 600;
}

.menu-toggle {
  display: grid;
  width: 38px;
  height: 38px;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
  color: var(--color-primary);
  cursor: pointer;
}

@media (max-width: 767px) {
  .app-topbar {
    min-height: 60px;
    padding: 0 16px;
  }

}

@media (max-width: 420px) {
  .topbar-user-copy {
    display: none;
  }
}
</style>
