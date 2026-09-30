<template>
  <DefaultLayout>
    <div class="dashboard">
      <header class="dashboard-heading">
        <div class="heading-copy">
          <span class="kicker">RINGKASAN RUANG KERJA</span>
          <h1>{{ greeting }}, {{ firstName }}</h1>
          <p>Berikut ringkasan akun dan akses Anda pada sistem Nexa.</p>
        </div>
        <time class="date-chip" :datetime="dateTime">
          <i class="pi pi-calendar" aria-hidden="true"></i>
          {{ dateLabel }}
        </time>
      </header>

      <section class="welcome-panel" aria-label="Pengguna saat ini">
        <span class="welcome-avatar" aria-hidden="true">{{ initials }}</span>
        <div class="welcome-copy">
          <span class="welcome-overline">SELAMAT DATANG</span>
          <h2>{{ authStore.user?.nama || authStore.user?.username || 'Pengguna' }}</h2>
          <p>
            {{ roleLabel }}
            <template v-if="authStore.user?.unit_nama"> &middot; {{ authStore.user.unit_nama }}</template>
          </p>
        </div>
        <span class="session-chip"><span class="session-dot"></span> Sesi aktif</span>
      </section>

      <div class="dashboard-grid">
        <section class="panel" aria-labelledby="account-title">
          <header class="panel-heading">
            <div>
              <span class="kicker">AKUN SAYA</span>
              <h2 id="account-title">Informasi akun</h2>
            </div>
            <span class="icon-tile"><i class="pi pi-id-card" aria-hidden="true"></i></span>
          </header>

          <dl class="detail-list">
            <div class="detail-row">
              <dt>Username</dt>
              <dd class="font-mono">{{ authStore.user?.username || '-' }}</dd>
            </div>
            <div class="detail-row">
              <dt>Peran sistem</dt>
              <dd>{{ roleLabel }}</dd>
            </div>
            <div class="detail-row">
              <dt>Profesi</dt>
              <dd>{{ authStore.user?.profesi || '-' }}</dd>
            </div>
            <div class="detail-row">
              <dt>Unit kerja</dt>
              <dd>{{ authStore.user?.unit_nama || '-' }}</dd>
            </div>
          </dl>
        </section>

        <section class="panel" aria-labelledby="workspace-title">
          <header class="panel-heading">
            <div>
              <span class="kicker">NAVIGASI</span>
              <h2 id="workspace-title">Ruang kerja</h2>
            </div>
            <span class="icon-tile"><i class="pi pi-sitemap" aria-hidden="true"></i></span>
          </header>

          <router-link v-if="isAdmin" to="/admin" class="workspace-link">
            <span class="workspace-icon"><i class="pi pi-sliders-h" aria-hidden="true"></i></span>
            <span class="workspace-text">
              <strong>Panel admin</strong>
              <small>Pengguna, unit, klaim, dan audit trail</small>
            </span>
            <i class="pi pi-arrow-right workspace-arrow" aria-hidden="true"></i>
          </router-link>

          <div v-else class="workspace-empty">
            <i class="pi pi-info-circle" aria-hidden="true"></i>
            <p>Belum ada modul operasional untuk peran <strong>{{ roleLabel }}</strong>.</p>
          </div>

          <p class="workspace-note">
            Menu dan tindakan mengikuti hak akses akun Anda. Setiap perubahan data tercatat pada audit trail.
          </p>
        </section>
      </div>

      <section class="dashboard-note">
        <span class="icon-tile"><i class="pi pi-lock" aria-hidden="true"></i></span>
        <div>
          <strong>Data terlindungi sesuai peran</strong>
          <p>Akses berbasis peran aktif untuk menjaga kerahasiaan rekam medis dan data klaim.</p>
        </div>
      </section>
    </div>
  </DefaultLayout>
</template>

<script setup>
import { computed } from 'vue'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import { useAuthStore } from '../stores/auth'

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

const roleLabel = computed(() => roleLabels[authStore.user?.role] || authStore.user?.role || 'Staf rumah sakit')
const isAdmin = computed(() => authStore.user?.role === 'admin_ti')
const firstName = computed(() => {
  const name = authStore.user?.nama || authStore.user?.username || 'Pengguna'
  return name.trim().split(/\s+/)[0]
})
const initials = computed(() => {
  const name = authStore.user?.nama || authStore.user?.username || 'N'
  return name.trim().split(/\s+/).slice(0, 2).map((part) => part.charAt(0).toUpperCase()).join('')
})

const now = new Date()
const dateTime = [
  now.getFullYear(),
  String(now.getMonth() + 1).padStart(2, '0'),
  String(now.getDate()).padStart(2, '0')
].join('-')
const dateLabel = new Intl.DateTimeFormat('id-ID', {
  weekday: 'long',
  day: 'numeric',
  month: 'long',
  year: 'numeric'
}).format(now)
const greeting = (() => {
  const hour = now.getHours()
  if (hour < 11) return 'Selamat pagi'
  if (hour < 15) return 'Selamat siang'
  if (hour < 18) return 'Selamat sore'
  return 'Selamat malam'
})()
</script>

<style scoped>
/* ---------- struktur halaman ---------- */
.dashboard {
  display: grid;
  width: 100%;
  gap: 18px;
  animation: dashboard-in 240ms ease-out both;
}

.dashboard-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
}

.kicker {
  display: block;
  color: var(--color-text-secondary);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1.25px;
}

.dashboard-heading h1 {
  margin: 8px 0 5px;
  color: var(--color-primary-dark);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.7px;
  line-height: 1.25;
}

.heading-copy > p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.date-chip {
  display: inline-flex;
  min-height: 40px;
  align-items: center;
  gap: 9px;
  padding: 0 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  font-size: 12px;
  white-space: nowrap;
}

.date-chip i {
  color: var(--color-primary);
}

/* ---------- panel selamat datang ---------- */
.welcome-panel {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 16px;
  padding: 20px 22px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: var(--color-surface);
  box-shadow: var(--shadow-panel);
}

.welcome-avatar {
  display: grid;
  width: 48px;
  height: 48px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 50%;
  background: #EEF2F5;
  color: var(--color-primary);
  font-size: 14px;
  font-weight: 600;
}

.welcome-copy {
  min-width: 0;
  flex: 1;
}

.welcome-overline {
  display: block;
  color: var(--color-text-secondary);
  font-size: 9px;
  font-weight: 600;
  letter-spacing: 1px;
}

.welcome-copy h2 {
  overflow: hidden;
  margin: 4px 0 3px;
  color: var(--color-primary-dark);
  font-size: 17px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.welcome-copy p {
  overflow: hidden;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-chip {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 8px;
  padding: 7px 12px;
  border: 1px solid #D5E3DA;
  border-radius: 999px;
  background: #F4F9F6;
  color: var(--color-status-approved);
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}

.session-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--color-status-approved);
}

/* ---------- grid dua panel ---------- */
.dashboard-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 18px;
}

.panel {
  display: flex;
  min-width: 0;
  flex-direction: column;
  padding: 20px 22px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: var(--color-surface);
  box-shadow: var(--shadow-panel);
}

.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--color-border);
}

.panel-heading h2 {
  margin: 5px 0 0;
  color: var(--color-primary-dark);
  font-size: 16px;
  font-weight: 600;
}

.icon-tile {
  display: grid;
  width: 36px;
  height: 36px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 6px;
  background: #EEF2F5;
  color: var(--color-primary);
  font-size: 14px;
}

/* ---------- daftar detail akun ---------- */
.detail-list {
  display: grid;
  gap: 13px;
  margin: 16px 0 0;
}

.detail-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 14px;
}

.detail-row dt {
  flex: 0 0 auto;
  color: var(--color-text-secondary);
  font-size: 12px;
}

.detail-row dd {
  overflow: hidden;
  margin: 0;
  color: var(--color-text-primary);
  font-size: 13px;
  font-weight: 500;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-row dd.font-mono {
  font-size: 12px;
}

/* ---------- tautan ruang kerja ---------- */
.workspace-link {
  display: flex;
  min-height: 62px;
  align-items: center;
  gap: 13px;
  margin-top: 16px;
  padding: 10px 13px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  text-decoration: none;
  transition: background-color var(--motion-fast), border-color var(--motion-fast);
}

.workspace-link:hover {
  border-color: #C8D1D9;
  background: var(--color-bg-base);
}

.workspace-icon {
  display: grid;
  width: 38px;
  height: 38px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 6px;
  background: #EEF2F5;
  color: var(--color-primary);
  font-size: 15px;
}

.workspace-text {
  display: grid;
  min-width: 0;
  flex: 1;
  gap: 3px;
}

.workspace-text strong {
  color: var(--color-text-primary);
  font-size: 13px;
  font-weight: 600;
}

.workspace-text small {
  overflow: hidden;
  color: var(--color-text-secondary);
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.workspace-arrow {
  flex: 0 0 auto;
  color: var(--color-text-secondary);
  font-size: 12px;
}

.workspace-empty {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-top: 16px;
  padding: 13px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: var(--color-bg-base);
  color: var(--color-text-secondary);
}

.workspace-empty > i {
  margin-top: 2px;
  color: var(--color-primary);
}

.workspace-empty p {
  margin: 0;
  font-size: 12px;
  line-height: 1.6;
}

.workspace-empty strong {
  color: var(--color-text-primary);
  font-weight: 600;
}

.workspace-note {
  margin: 14px 0 0;
  color: var(--color-text-secondary);
  font-size: 11px;
  line-height: 1.65;
}

/* ---------- catatan bawah ---------- */
.dashboard-note {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 16px 20px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: var(--color-surface);
  box-shadow: var(--shadow-panel);
}

.dashboard-note > div {
  min-width: 0;
  flex: 1;
}

.dashboard-note strong {
  color: var(--color-text-primary);
  font-size: 13px;
  font-weight: 600;
}

.dashboard-note p {
  margin: 3px 0 0;
  color: var(--color-text-secondary);
  font-size: 11px;
  line-height: 1.6;
}

/* ---------- animasi ---------- */
@keyframes dashboard-in {
  from {
    opacity: 0;
    transform: translateY(6px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ---------- responsif ---------- */
@media (max-width: 900px) {
  .dashboard-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 640px) {
  .dashboard {
    gap: 14px;
  }

  .dashboard-heading {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
  }

  .dashboard-heading h1 {
    font-size: 21px;
  }

  .date-chip {
    min-height: 36px;
    font-size: 11px;
  }

  .welcome-panel {
    gap: 12px;
    padding: 16px;
  }

  .welcome-avatar {
    width: 42px;
    height: 42px;
  }

  .session-chip {
    display: none;
  }

  .panel {
    padding: 16px;
  }

  .dashboard-note {
    align-items: flex-start;
    padding: 14px 16px;
  }
}
</style>
