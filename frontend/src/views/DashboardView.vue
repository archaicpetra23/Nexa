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

      <!-- Stats Cards -->
      <div class="stat-cards">
        <div class="stat-card">
          <div class="stat-icon" style="background: #EEF2F5; color: var(--color-primary)">
            <i class="pi pi-users" aria-hidden="true"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">Total Pasien</span>
            <span class="stat-value font-mono">{{ formatNumber(stats.total_pasien) }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="background: #EEF2F5; color: var(--color-primary)">
            <i class="pi pi-file-check" aria-hidden="true"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">Total Klaim</span>
            <span class="stat-value font-mono">{{ formatNumber(stats.total_klaim) }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="background: #FEF3C7; color: #B7791F">
            <i class="pi pi-clock" aria-hidden="true"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">Klaim Pending</span>
            <span class="stat-value font-mono">{{ formatNumber(stats.klaim_pending) }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="background: #D1FAE5; color: #2F7A4F">
            <i class="pi pi-check" aria-hidden="true"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">Klaim Disetujui</span>
            <span class="stat-value font-mono">{{ formatNumber(stats.klaim_disetujui) }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="background: #FEE2E2; color: #B54245">
            <i class="pi pi-times" aria-hidden="true"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">Klaim Ditolak</span>
            <span class="stat-value font-mono">{{ formatNumber(stats.klaim_ditolak) }}</span>
          </div>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="charts-grid">
        <!-- Top Diagnosis Chart -->
        <div class="chart-panel">
          <div class="chart-header">
            <h3>Penyakit Terbanyak</h3>
            <span class="chart-caption">Berdasarkan kode ICD-10</span>
          </div>
          <div v-if="stats.top_diagnosis && stats.top_diagnosis.length" class="chart-container">
            <canvas ref="diagnosisChart"></canvas>
          </div>
          <div v-else class="chart-empty">
            <i class="pi pi-inbox" aria-hidden="true"></i>
            <p>Belum ada data diagnosis</p>
          </div>
        </div>

        <!-- Claim Status Chart -->
        <div class="chart-panel">
          <div class="chart-header">
            <h3>Status Klaim</h3>
            <span class="chart-caption">Distribusi status klaim BPJS</span>
          </div>
          <div v-if="hasClaimData" class="chart-container">
            <canvas ref="statusChart"></canvas>
          </div>
          <div v-else class="chart-empty">
            <i class="pi pi-inbox" aria-hidden="true"></i>
            <p>Belum ada data klaim</p>
          </div>
        </div>
      </div>

      <!-- Nominal Summary -->
      <div class="nominal-summary">
        <i class="pi pi-money-bill" aria-hidden="true"></i>
        <span class="summary-label">Total Nominal Disetujui:</span>
        <span class="summary-value font-mono">{{ formatRupiah(stats.total_nominal_disetujui) }}</span>
      </div>

      <!-- Recent Claims -->
      <section class="recent-claims-panel" aria-labelledby="recent-title">
        <header class="panel-heading">
          <div>
            <h2 id="recent-title">Klaim Terbaru</h2>
            <p>5 klaim terbaru dalam sistem</p>
          </div>
          <router-link to="/klaim" class="view-all-link">Lihat semua</router-link>
        </header>

        <ModernTable
          :value="recentClaims"
          :loading="claimsLoading"
          :paginator="false"
          dataKey="id_klaim"
          ariaLabel="Daftar klaim terbaru"
        >
          <template #header>
            <th scope="col">ID Klaim</th>
            <th scope="col">Pasien</th>
            <th scope="col">Kode ICD</th>
            <th scope="col">Nominal</th>
            <th scope="col">Status</th>
          </template>
          <template #row="{ item }">
            <td class="font-mono">{{ item.id_klaim }}</td>
            <td class="user-name-cell">{{ item.pasien_nama || '-' }}</td>
            <td class="font-mono">{{ item.kode_cbgs || '-' }}</td>
            <td class="font-mono amount-cell">{{ formatRupiah(item.nominal_klaim) }}</td>
            <td><StatusBadge :status="item.status_klaim" type="claim" /></td>
          </template>
        </ModernTable>
      </section>

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
import { computed, nextTick, onMounted, ref } from 'vue'
import { Chart, registerables } from 'chart.js'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import ModernTable from '../components/ModernTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useAuthStore } from '../stores/auth'
import client from '../api/client'

Chart.register(...registerables)

const authStore = useAuthStore()

const stats = ref({
  total_pasien: 0,
  total_klaim: 0,
  klaim_pending: 0,
  klaim_disetujui: 0,
  klaim_ditolak: 0,
  klaim_draft: 0,
  total_nominal_disetujui: 0,
  top_diagnosis: []
})

const recentClaims = ref([])
const claimsLoading = ref(false)

const diagnosisChart = ref(null)
const statusChart = ref(null)
let diagnosisChartInstance = null
let statusChartInstance = null

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

const hasClaimData = computed(() => {
  return stats.value.klaim_disetujui > 0 || stats.value.klaim_pending > 0 || stats.value.klaim_ditolak > 0
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

async function fetchStats() {
  try {
    const res = await client.get('/dashboard/stats')
    if (res.data.success) {
      stats.value = res.data.data
      await nextTick()
      renderCharts()
      fetchRecentClaims()
    }
  } catch (e) {
    console.error('Failed to fetch dashboard stats:', e)
  }
}

function renderCharts() {
  renderDiagnosisChart()
  renderStatusChart()
}

function renderDiagnosisChart() {
  if (!diagnosisChart.value || !stats.value.top_diagnosis || !stats.value.top_diagnosis.length) return

  if (diagnosisChartInstance) {
    diagnosisChartInstance.destroy()
  }

  const ctx = diagnosisChart.value.getContext('2d')
  const labels = stats.value.top_diagnosis.map(d => d.kode_icd10)
  const primerData = stats.value.top_diagnosis.map(d => d.jumlah_primer)
  const sekunderData = stats.value.top_diagnosis.map(d => d.jumlah_sekunder)
  const totalData = stats.value.top_diagnosis.map(d => d.jumlah_primer + d.jumlah_sekunder)

  const gradientPrimer = ctx.createLinearGradient(0, 0, 0, 240)
  gradientPrimer.addColorStop(0, 'rgba(28, 61, 90, 0.20)')
  gradientPrimer.addColorStop(1, 'rgba(28, 61, 90, 0.01)')

  const gradientSekunder = ctx.createLinearGradient(0, 0, 0, 240)
  gradientSekunder.addColorStop(0, 'rgba(46, 139, 154, 0.15)')
  gradientSekunder.addColorStop(1, 'rgba(46, 139, 154, 0.01)')

  const gradientTotal = ctx.createLinearGradient(0, 0, 0, 240)
  gradientTotal.addColorStop(0, 'rgba(217, 119, 6, 0.12)')
  gradientTotal.addColorStop(1, 'rgba(217, 119, 6, 0.01)')

  diagnosisChartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels: labels,
      datasets: [
        {
          label: 'Primer',
          data: primerData,
          borderColor: '#1C3D5A',
          backgroundColor: gradientPrimer,
          borderWidth: 2.5,
          pointBackgroundColor: '#1C3D5A',
          pointBorderColor: '#fff',
          pointBorderWidth: 2,
          pointRadius: 4,
          pointHoverRadius: 6,
          fill: true,
          tension: 0.3
        },
        {
          label: 'Sekunder',
          data: sekunderData,
          borderColor: '#2E8B9A',
          backgroundColor: gradientSekunder,
          borderWidth: 2,
          pointBackgroundColor: '#2E8B9A',
          pointBorderColor: '#fff',
          pointBorderWidth: 2,
          pointRadius: 3,
          pointHoverRadius: 5,
          fill: true,
          tension: 0.3
        },
        {
          label: 'Total',
          data: totalData,
          borderColor: '#D97706',
          backgroundColor: gradientTotal,
          borderWidth: 2,
          pointBackgroundColor: '#D97706',
          pointBorderColor: '#fff',
          pointBorderWidth: 2,
          pointRadius: 3,
          pointHoverRadius: 5,
          fill: true,
          tension: 0.3
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      plugins: {
        legend: {
          position: 'top',
          labels: { font: { size: 11 }, padding: 12 }
        }
      },
      scales: {
        x: {
          beginAtZero: true,
          grid: { drawBorder: false, color: 'rgba(0, 0, 0, 0.05)' },
          ticks: { font: { size: 11, family: 'IBM Plex Mono' } }
        },
        y: {
          beginAtZero: true,
          grid: { color: 'rgba(0, 0, 0, 0.05)' },
          ticks: { font: { size: 11 } }
        }
      }
    }
  })
}

function renderStatusChart() {
  if (!statusChart.value || !hasClaimData.value) return

  if (statusChartInstance) {
    statusChartInstance.destroy()
  }

  const ctx = statusChart.value.getContext('2d')
  statusChartInstance = new Chart(ctx, {
    type: 'pie',
    data: {
      labels: ['Disetujui', 'Pending', 'Ditolak'],
      datasets: [
        {
          data: [
            stats.value.klaim_disetujui,
            stats.value.klaim_pending,
            stats.value.klaim_ditolak
          ],
          backgroundColor: [
            '#2F7A4F',
            '#B7791F',
            '#B54245'
          ],
          borderColor: '#fff',
          borderWidth: 2,
          hoverOffset: 8
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      plugins: {
        legend: {
          position: 'bottom',
          labels: { font: { size: 11 }, padding: 16 }
        }
      }
    }
  })
}

async function fetchRecentClaims() {
  claimsLoading.value = true
  try {
    const res = await client.get('/klaim', { params: { page: 1, limit: 5 } })
    if (res.data.success) {
      recentClaims.value = res.data.data.items || []
    }
  } catch (e) {
    console.error('Failed to fetch recent claims:', e)
  } finally {
    claimsLoading.value = false
  }
}

function formatNumber(value) {
  return new Intl.NumberFormat('id-ID').format(value)
}

function formatRupiah(value) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value)
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
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
  font-size: var(--text-2xs);
  font-weight: 600;
  letter-spacing: 1.25px;
  text-transform: uppercase;
}

.dashboard-heading h1 {
  margin: 8px 0 5px;
  color: var(--color-primary-dark);
  font-size: var(--text-2xl);
  font-weight: 600;
  letter-spacing: -0.7px;
  line-height: 1.25;
}

.heading-copy > p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
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
  font-size: var(--text-xs);
  white-space: nowrap;
}

.date-chip i {
  color: var(--color-primary);
}

.stat-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
  box-shadow: var(--shadow-panel);
}

.stat-icon {
  display: grid;
  width: 44px;
  height: 44px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: var(--radius-md);
  font-size: var(--text-xl);
}

.stat-content {
  display: grid;
  gap: 6px;
  min-width: 0;
}

.stat-label {
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  font-weight: 500;
}

.stat-value {
  color: var(--color-text-primary);
  font-size: var(--text-xl);
  font-weight: 700;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 14px;
}

.chart-panel {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
  box-shadow: var(--shadow-panel);
  overflow: hidden;
}

.chart-header {
  padding: 16px;
  border-bottom: 1px solid var(--color-border);
}

.chart-header h3 {
  margin: 0 0 4px;
  color: var(--color-text-primary);
  font-size: var(--text-md);
  font-weight: 600;
}

.chart-caption {
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
}

.chart-container {
  position: relative;
  height: 240px;
  padding: 12px 16px;
}

.chart-empty {
  display: flex;
  min-height: 240px;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  color: var(--color-text-secondary);
  text-align: center;
}

.chart-empty i {
  font-size: var(--text-2xl);
  color: var(--color-text-disabled);
}

.nominal-summary {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
  box-shadow: var(--shadow-panel);
}

.nominal-summary i {
  color: var(--color-primary);
  font-size: var(--text-lg);
}

.summary-label {
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-weight: 500;
}

.summary-value {
  color: var(--color-text-primary);
  font-size: var(--text-lg);
  font-weight: 700;
}

.recent-claims-panel {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
  box-shadow: var(--shadow-panel);
  overflow: hidden;
}

.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  border-bottom: 1px solid var(--color-border);
}

.panel-heading h2 {
  margin: 0 0 4px;
  color: var(--color-text-primary);
  font-size: var(--text-md);
  font-weight: 600;
}

.panel-heading p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
}

.view-all-link {
  color: var(--color-primary);
  text-decoration: none;
  font-size: var(--text-xs);
  font-weight: 600;
  transition: color 140ms ease;
}

.view-all-link:hover {
  color: var(--color-primary-hover);
}

.amount-cell {
  text-align: right;
}

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
  font-size: var(--text-base);
  font-weight: 600;
}

.welcome-copy {
  min-width: 0;
  flex: 1;
}

.welcome-overline {
  display: block;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  font-weight: 600;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.welcome-copy h2 {
  overflow: hidden;
  margin: 4px 0 3px;
  color: var(--color-primary-dark);
  font-size: var(--text-lg);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.welcome-copy p {
  overflow: hidden;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
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
  font-size: var(--text-2xs);
  font-weight: 500;
  white-space: nowrap;
}

.session-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--color-status-approved);
}

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

.icon-tile {
  display: grid;
  width: 36px;
  height: 36px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: var(--radius-sm);
  background: #EEF2F5;
  color: var(--color-primary);
  font-size: var(--text-base);
}

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
  font-size: var(--text-xs);
}

.detail-row dd {
  overflow: hidden;
  margin: 0;
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-weight: 500;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-row dd.font-mono {
  font-size: var(--text-xs);
}

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
  font-size: var(--text-md);
}

.workspace-text {
  display: grid;
  min-width: 0;
  flex: 1;
  gap: 3px;
}

.workspace-text strong {
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-weight: 600;
}

.workspace-text small {
  overflow: hidden;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.workspace-arrow {
  flex: 0 0 auto;
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
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
  font-size: var(--text-xs);
  line-height: 1.6;
}

.workspace-empty strong {
  color: var(--color-text-primary);
  font-weight: 600;
}

.workspace-note {
  margin: 14px 0 0;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  line-height: 1.65;
}

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
  font-size: var(--text-sm);
  font-weight: 600;
}

.dashboard-note p {
  margin: 3px 0 0;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  line-height: 1.6;
}

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

@media (max-width: 900px) {
  .dashboard-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .charts-grid {
    grid-template-columns: 1fr;
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
    font-size: var(--text-2xs);
  }

  .stat-cards {
    grid-template-columns: 1fr;
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

@media (prefers-reduced-motion: reduce) {
  .dashboard {
    animation: none;
  }
}
</style>

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
  font-size: var(--text-2xs);
  font-weight: 600;
  letter-spacing: 1.25px;
  text-transform: uppercase;
}

.dashboard-heading h1 {
  margin: 8px 0 5px;
  color: var(--color-primary-dark);
  font-size: var(--text-2xl);
  font-weight: 600;
  letter-spacing: -0.7px;
  line-height: 1.25;
}

.heading-copy > p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
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
  font-size: var(--text-xs);
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
  font-size: var(--text-base);
  font-weight: 600;
}

.welcome-copy {
  min-width: 0;
  flex: 1;
}

.welcome-overline {
  display: block;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  font-weight: 600;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.welcome-copy h2 {
  overflow: hidden;
  margin: 4px 0 3px;
  color: var(--color-primary-dark);
  font-size: var(--text-lg);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.welcome-copy p {
  overflow: hidden;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
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
  font-size: var(--text-2xs);
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
  border-radius: var(--radius-sm);
  background: #EEF2F5;
  color: var(--color-primary);
  font-size: var(--text-base);
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
  font-size: var(--text-xs);
}

.detail-row dd {
  overflow: hidden;
  margin: 0;
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-weight: 500;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-row dd.font-mono {
  font-size: var(--text-xs);
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
  font-size: var(--text-md);
}

.workspace-text {
  display: grid;
  min-width: 0;
  flex: 1;
  gap: 3px;
}

.workspace-text strong {
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-weight: 600;
}

.workspace-text small {
  overflow: hidden;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.workspace-arrow {
  flex: 0 0 auto;
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
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
  font-size: var(--text-xs);
  line-height: 1.6;
}

.workspace-empty strong {
  color: var(--color-text-primary);
  font-weight: 600;
}

.workspace-note {
  margin: 14px 0 0;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
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
  font-size: var(--text-sm);
  font-weight: 600;
}

.dashboard-note p {
  margin: 3px 0 0;
  color: var(--color-text-secondary);
  font-size: var(--text-2xs);
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
    font-size: var(--text-2xs);
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
