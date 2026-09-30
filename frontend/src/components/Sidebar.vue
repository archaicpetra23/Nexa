<template>
  <aside id="app-sidebar" :class="['app-sidebar', { 'is-open': mobileOpen, 'is-hidden': desktopHidden }]" aria-label="Navigasi aplikasi">
    <div class="sidebar-header">
      <router-link to="/dashboard" class="brand" @click="emit('close')">
        <span class="brand-icon"><i class="pi pi-plus" aria-hidden="true"></i></span>
        <span class="brand-name">Nexa</span>
      </router-link>
      <button class="mobile-close" type="button" aria-label="Tutup navigasi" @click="emit('close')">
        <i class="pi pi-times" aria-hidden="true"></i>
      </button>
    </div>

    <div class="sidebar-section-label">MENU UTAMA</div>
    <nav class="sidebar-navigation" aria-label="Menu utama">
      <router-link
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        :class="['navigation-link', { 'is-active': isActiveRoute(item.to) }]"
        :aria-current="isActiveRoute(item.to) ? 'page' : undefined"
        @click="emit('close')"
      >
        <i :class="item.icon" aria-hidden="true"></i>
        <span>{{ item.label }}</span>
      </router-link>
    </nav>

    <div class="sidebar-bottom">
      <div class="security-card">
        <i class="pi pi-shield" aria-hidden="true"></i>
        <div>
          <strong>Akses terlindungi</strong>
          <span>Menu mengikuti peran akun</span>
        </div>
      </div>
      <div class="user-card">
        <span class="user-avatar">{{ initials }}</span>
        <span class="user-information">
          <strong>{{ authStore.user?.nama || authStore.user?.username || 'Pengguna' }}</strong>
          <small>{{ roleLabel }}</small>
        </span>
        <button class="logout-button" type="button" aria-label="Keluar" title="Keluar" @click="handleLogout">
          <i class="pi pi-sign-out" aria-hidden="true"></i>
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

defineProps({
  mobileOpen: { type: Boolean, default: false },
  desktopHidden: { type: Boolean, default: false }
})

const emit = defineEmits(['close'])
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const roleLabels = {
  admin_ti: 'Administrator TI',
  petugas_rm: 'Petugas Rekam Medis',
  dokter_dpjp: 'Dokter DPJP',
  dokter: 'Dokter',
  perawat: 'Perawat',
  petugas_casemix: 'Petugas Casemix',
  koder_casemix: 'Koder Casemix',
  keuangan: 'Keuangan',
  manajemen: 'Manajemen'
}

const menuItems = computed(() => {
  const items = [{ label: 'Dashboard', to: '/dashboard', icon: 'pi pi-home' }]
  if (authStore.user?.role === 'admin_ti') {
    items.push({ label: 'Panel admin', to: '/admin', icon: 'pi pi-sliders-h' })
  }
  return items
})

const roleLabel = computed(() => roleLabels[authStore.user?.role] || authStore.user?.role || 'Staf rumah sakit')
const initials = computed(() => {
  const name = authStore.user?.nama || authStore.user?.username || 'N'
  return name.trim().split(/\s+/).slice(0, 2).map((part) => part.charAt(0).toUpperCase()).join('')
})

function isActiveRoute(to) {
  return route.path === to || route.path.startsWith(`${to}/`)
}

async function handleLogout() {
  await authStore.logout()
  emit('close')
  router.push('/login')
}
</script>

<style scoped>
.app-sidebar {
  position: sticky;
  top: 0;
  display: flex;
  width: 224px;
  height: 100vh;
  flex: 0 0 224px;
  flex-direction: column;
  padding: 20px 14px 14px;
  overflow-y: auto;
  background: var(--color-primary);
  color: white;
  transition: width 180ms ease-out, flex-basis 180ms ease-out, padding 180ms ease-out;
}

.app-sidebar.is-hidden {
  width: 0;
  flex-basis: 0;
  padding-right: 0;
  padding-left: 0;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  min-height: 48px;
  align-items: center;
  justify-content: space-between;
  padding: 0 7px 17px;
  border-bottom: 1px solid #36536D;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: white;
  text-decoration: none;
}

.brand-icon {
  display: grid;
  width: 31px;
  height: 31px;
  place-items: center;
  border-radius: var(--radius-base);
  background: white;
  color: var(--color-primary);
  font-size: 13px;
}

.brand-name {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.4px;
}

.mobile-close {
  display: none;
}

.sidebar-section-label {
  margin: 27px 10px 10px;
  color: #C0CCD5;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.7px;
}

.sidebar-navigation {
  display: grid;
  gap: 5px;
}

.navigation-link {
  display: flex;
  min-height: 42px;
  align-items: center;
  gap: 11px;
  padding: 0 11px;
  border-radius: var(--radius-base);
  color: #E1E7EC;
  font-size: 13px;
  text-decoration: none;
  transition: background-color var(--motion-fast), color var(--motion-fast);
}

.navigation-link > i {
  width: 17px;
  color: #C0CCD5;
  font-size: 14px;
  text-align: center;
}

.navigation-link:hover {
  background: var(--color-primary-hover);
  color: white;
}

.navigation-link.is-active {
  background: var(--color-primary-dark);
  color: white;
  font-weight: 500;
}

.navigation-link.is-active > i {
  color: white;
}

.sidebar-bottom {
  display: grid;
  gap: 14px;
  margin-top: auto;
  padding-top: 20px;
}

.security-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 11px 10px;
  border: 1px solid #36536D;
  border-radius: var(--radius-base);
  background: #234563;
}

.security-card > i {
  color: white;
  font-size: 15px;
}

.security-card > div,
.user-information {
  display: grid;
  min-width: 0;
  gap: 3px;
}

.security-card strong,
.user-information strong {
  overflow: hidden;
  color: white;
  font-size: 11px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.security-card span,
.user-information small {
  overflow: hidden;
  color: #D3DDE4;
  font-size: 10px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-card {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 9px;
  padding: 12px 2px 0;
  border-top: 1px solid #36536D;
}

.user-avatar {
  display: grid;
  width: 34px;
  height: 34px;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid #557087;
  border-radius: 50%;
  background: #355672;
  color: white;
  font-size: 10px;
  font-weight: 600;
}

.user-information {
  flex: 1;
}

.logout-button {
  display: grid;
  width: 38px;
  height: 38px;
  flex: 0 0 auto;
  place-items: center;
  border: 0;
  border-radius: var(--radius-base);
  background: transparent;
  color: #E1E7EC;
  cursor: pointer;
}

.logout-button:hover {
  background: var(--color-primary-dark);
  color: white;
}

@media (max-width: 767px) {
  .app-sidebar {
    position: fixed;
    z-index: 40;
    inset: 0 auto 0 0;
    width: min(288px, calc(100vw - 48px));
    height: 100dvh;
    transform: translateX(-100%);
    transition: transform 180ms ease-out, width 180ms ease-out, flex-basis 180ms ease-out;
  }

  .app-sidebar.is-open {
    transform: translateX(0);
  }

  .app-sidebar.is-hidden {
    width: min(288px, calc(100vw - 48px));
    flex-basis: auto;
    padding: 20px 14px 14px;
    overflow-y: auto;
  }

  .mobile-close {
    display: grid;
    width: 38px;
    height: 38px;
    place-items: center;
    border: 0;
    border-radius: var(--radius-base);
    background: transparent;
    color: white;
    cursor: pointer;
  }
}
</style>
