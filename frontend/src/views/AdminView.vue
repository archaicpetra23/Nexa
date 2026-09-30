<template>
  <div class="min-h-screen bg-base">
    <header class="bg-primary text-white py-4 px-8">
      <div class="max-w-screen-xl mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-semibold">Nexa</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm">{{ authStore.user?.nama }}</span>
          <span class="text-sm opacity-75">{{ authStore.user?.role }}</span>
          <router-link to="/dashboard" class="px-4 py-2 bg-white text-primary rounded text-sm hover:bg-gray-100 transition-colors">
            Kembali ke Dashboard
          </router-link>
          <button @click="handleLogout" class="px-4 py-2 bg-white text-primary rounded text-sm hover:bg-gray-100 transition-colors">
            Logout
          </button>
        </div>
      </div>
    </header>

    <main class="max-w-screen-xl mx-auto px-8 py-10">
      <h1 class="text-2xl font-semibold text-primary mb-8">Panel Admin</h1>

      <div v-if="loading" class="text-center py-12">
        <p class="text-secondary">Memuat data...</p>
      </div>

      <template v-else>
        <section class="mb-12">
          <h2 class="text-lg font-semibold text-primary mb-4">Manajemen Pengguna</h2>
          <div class="bg-surface rounded border border-border overflow-hidden">
            <table class="w-full">
              <thead class="bg-base">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">ID</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Username</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Nama</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Profesi</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Role</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr v-for="user in users" :key="user.id_user" class="hover:bg-base">
                  <td class="px-4 py-3 text-sm font-mono">{{ user.id_user }}</td>
                  <td class="px-4 py-3 text-sm font-mono">{{ user.username }}</td>
                  <td class="px-4 py-3 text-sm">{{ user.nama }}</td>
                  <td class="px-4 py-3 text-sm">{{ user.profesi }}</td>
                  <td class="px-4 py-3 text-sm">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary/10 text-primary">
                      {{ user.role }}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-sm">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                      :class="user.deleted_at ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'">
                      {{ user.deleted_at ? 'Nonaktif' : 'Aktif' }}
                    </span>
                  </td>
                </tr>
                <tr v-if="users.length === 0">
                  <td colspan="6" class="px-4 py-8 text-center text-secondary">Belum ada data pengguna</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section>
          <h2 class="text-lg font-semibold text-primary mb-4">Manajemen Klaim</h2>
          <div class="bg-surface rounded border border-border overflow-hidden">
            <table class="w-full">
              <thead class="bg-base">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">ID Klaim</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">No. RM</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Pasien</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Kode CBGs</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Status</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Nominal</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Petugas</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-secondary uppercase">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr v-for="klaim in klaims" :key="klaim.id_klaim" class="hover:bg-base">
                  <td class="px-4 py-3 text-sm font-mono">{{ klaim.id_klaim }}</td>
                  <td class="px-4 py-3 text-sm font-mono">{{ klaim.id_rekam }}</td>
                  <td class="px-4 py-3 text-sm">{{ klaim.pasien_nama }}</td>
                  <td class="px-4 py-3 text-sm font-mono">{{ klaim.kode_cbgs }}</td>
                  <td class="px-4 py-3 text-sm">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                      :class="statusClass(klaim.status_klaim)">
                      {{ klaim.status_klaim }}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-sm font-mono text-right">{{ formatNumber(klaim.nominal_klaim) }}</td>
                  <td class="px-4 py-3 text-sm">{{ klaim.petugas_nama }}</td>
                  <td class="px-4 py-3">
                    <button
                      v-if="klaim.status_klaim === 'draft' || klaim.status_klaim === 'pending'"
                      @click="openScoreModal(klaim)"
                      class="text-sm text-primary hover:underline"
                    >
                      Skoring
                    </button>
                    <span v-else class="text-sm text-secondary">-</span>
                  </td>
                </tr>
                <tr v-if="klaims.length === 0">
                  <td colspan="8" class="px-4 py-8 text-center text-secondary">Belum ada data klaim</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </template>

      <!-- Score Modal -->
      <div v-if="selectedKlaim" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-surface rounded-lg shadow-xl max-w-lg w-full p-6">
          <h3 class="text-lg font-semibold text-primary mb-4">Skoring Klaim #{{ selectedKlaim.id_klaim }}</h3>
          <p class="text-sm text-secondary mb-4">Pasien: {{ selectedKlaim.pasien_nama }}</p>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-secondary mb-1">Status Baru</label>
            <select v-model="newStatus" class="w-full px-3 py-2 border border-border rounded focus:outline-none focus:ring-2 focus:ring-primary">
              <option value="pending">Pending</option>
              <option value="disetujui">Disetujui</option>
              <option value="ditolak">Ditolak</option>
            </select>
          </div>

          <div class="mb-4" v-if="newStatus === 'pending' || newStatus === 'ditolak'">
            <label class="block text-sm font-medium text-secondary mb-1">Alasan (Wajib)</label>
            <textarea v-model="alasan" rows="3" class="w-full px-3 py-2 border border-border rounded focus:outline-none focus:ring-2 focus:ring-primary" placeholder="Masukkan alasan pending/ditolak"></textarea>
          </div>

          <div class="flex justify-end gap-3">
            <button @click="closeScoreModal" class="px-4 py-2 border border-border rounded text-secondary hover:bg-base transition-colors">
              Batal
            </button>
            <button @click="submitScore" class="px-4 py-2 bg-primary text-white rounded hover:bg-primary-hover transition-colors">
              Simpan
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import client from '../api/client'

const router = useRouter()
const authStore = useAuthStore()

const users = ref([])
const klaims = ref([])
const loading = ref(true)

const selectedKlaim = ref(null)
const newStatus = ref('pending')
const alasan = ref('')

async function fetchUsers() {
  try {
    const res = await client.get('/admin/users')
    if (res.data.success) users.value = res.data.data.items || res.data.data
  } catch (e) {
    console.error('Fetch users error:', e)
  }
}

async function fetchKlaims() {
  try {
    const res = await client.get('/admin/submissions')
    if (res.data.success) klaims.value = res.data.data.items || res.data.data
  } catch (e) {
    console.error('Fetch klaims error:', e)
  }
}

async function loadData() {
  loading.value = true
  await Promise.all([fetchUsers(), fetchKlaims()])
  loading.value = false
}

onMounted(loadData)

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

function statusClass(status) {
  const classes = {
    draft: 'bg-gray-100 text-gray-700',
    pending: 'bg-amber-100 text-amber-800',
    disetujui: 'bg-green-100 text-green-700',
    ditolak: 'bg-red-100 text-red-700'
  }
  return classes[status] || 'bg-gray-100 text-gray-700'
}

function formatNumber(num) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num)
}

function openScoreModal(klaim) {
  selectedKlaim.value = klaim
  newStatus.value = 'pending'
  alasan.value = ''
}

function closeScoreModal() {
  selectedKlaim.value = null
  newStatus.value = 'pending'
  alasan.value = ''
}

async function submitScore() {
  if (!selectedKlaim.value) return
  if ((newStatus.value === 'pending' || newStatus.value === 'ditolak') && !alasan.value.trim()) {
    alert('Alasan wajib diisi untuk status pending atau ditolak')
    return
  }

  try {
    await client.patch(`/klaim/${selectedKlaim.value.id_klaim}/status`, {
      status_klaim: newStatus.value,
      alasan_pending_tolak: alasan.value.trim() || null
    })
    closeScoreModal()
    await fetchKlaims()
  } catch (e) {
    console.error('Update status error:', e)
    alert('Gagal memperbarui status: ' + (e.response?.data?.message || 'Unknown error'))
  }
}
</script>