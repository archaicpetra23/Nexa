<template>
  <div :class="['app-shell', { 'is-sidebar-hidden': !sidebarVisible, 'is-admin-page': route.name === 'admin' }]">
    <Sidebar
      :mobile-open="mobileSidebarOpen"
      :desktop-hidden="!sidebarVisible"
      @close="mobileSidebarOpen = false"
    />

    <button
      v-if="mobileSidebarOpen"
      class="mobile-backdrop"
      type="button"
      aria-label="Tutup menu navigasi"
      @click="mobileSidebarOpen = false"
    ></button>

    <div class="app-main-column">
      <Topbar :sidebar-expanded="sidebarExpanded" :is-mobile="isMobile" @toggle-sidebar="toggleSidebar" />

      <main class="app-content">
        <div class="app-content-inner">
          <slot />
        </div>
      </main>

      <footer class="app-footer">
        <span>Nexa Casemix Management System</span>
        <span class="app-footer-secure">
          <i class="pi pi-lock" aria-hidden="true"></i> Akses aman &amp; berbasis peran
        </span>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import Topbar from '../components/Topbar.vue'

const sidebarVisible = ref(true)
const mobileSidebarOpen = ref(false)
const isMobile = ref(false)
const route = useRoute()
let mobileQuery

const sidebarExpanded = computed(() => isMobile.value ? mobileSidebarOpen.value : sidebarVisible.value)

function updateViewport() {
  isMobile.value = mobileQuery.matches
  mobileSidebarOpen.value = false
}

function toggleSidebar() {
  if (isMobile.value) {
    mobileSidebarOpen.value = !mobileSidebarOpen.value
    return
  }

  sidebarVisible.value = !sidebarVisible.value
}

onMounted(() => {
  mobileQuery = window.matchMedia('(max-width: 767px)')
  updateViewport()
  mobileQuery.addEventListener('change', updateViewport)
})

onUnmounted(() => {
  mobileQuery?.removeEventListener('change', updateViewport)
})

watch(() => route.path, () => {
  mobileSidebarOpen.value = false
})
</script>

<style scoped>
.app-shell {
  display: flex;
  min-height: 100vh;
  align-items: stretch;
  background: var(--color-bg-base);
}

.app-shell.is-admin-page {
  height: 100vh;
  height: 100dvh;
  min-height: 0;
  overflow: hidden;
}

.app-main-column {
  display: flex;
  min-width: 0;
  min-height: 100vh;
  flex: 1 1 auto;
  flex-direction: column;
}

.is-admin-page .app-main-column {
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.app-content {
  flex: 1 1 auto;
  width: 100%;
}

.is-admin-page .app-content {
  min-height: 0;
  overflow: hidden;
}

.app-content-inner {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 30px 40px;
}

.is-admin-page .app-content-inner {
  height: 100%;
  overflow: hidden;
  padding-top: 16px;
  padding-bottom: 16px;
}

.app-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin: 0 30px;
  padding: 15px 0 19px;
  border-top: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: 10px;
}

.is-admin-page .app-footer {
  flex: 0 0 auto;
  margin-top: 0;
  padding-top: 10px;
  padding-bottom: 12px;
}

.app-footer-secure {
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.app-footer i {
  color: var(--color-primary);
}

.mobile-backdrop {
  display: none;
}

@media (max-width: 900px) {
  .app-content-inner {
    padding: 24px 20px 32px;
  }

  .app-footer {
    margin: 0 20px;
  }

  .is-admin-page .app-content-inner {
    padding-top: 14px;
    padding-bottom: 14px;
  }
}

@media (max-width: 767px) {
  .mobile-backdrop {
    position: fixed;
    z-index: 20;
    inset: 0;
    display: block;
    border: 0;
    background: rgb(15 39 64 / 34%);
  }

  .app-content-inner {
    padding: 20px 16px 28px;
  }

  .app-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    margin: 0 16px;
    padding-bottom: 16px;
  }

  .is-admin-page .app-footer {
    gap: 4px;
    padding-top: 8px;
    padding-bottom: 8px;
  }
}
</style>
