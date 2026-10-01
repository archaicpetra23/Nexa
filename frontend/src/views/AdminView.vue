<template>
  <DefaultLayout>
    <div class="admin-page">
      <header class="admin-hero">
        <div class="admin-hero-copy">
          <span class="admin-eyebrow"><i class="pi pi-shield" aria-hidden="true"></i> PUSAT KONTROL</span>
          <h1>Panel Admin</h1>
          <p>Kelola pengguna, klaim, unit, dan aktivitas sistem dari satu tempat.</p>
        </div>
        <div class="admin-hero-actions">
          <button v-if="activeTab === 'users'" @click="openUserModal" class="admin-primary-action" type="button">
            <i class="pi pi-plus" aria-hidden="true"></i>
            <span>Tambah Pengguna</span>
          </button>
          <button v-if="activeTab === 'units'" @click="openUnitModal" class="admin-primary-action" type="button">
            <i class="pi pi-plus" aria-hidden="true"></i>
            <span>Tambah Unit</span>
          </button>
        </div>
      </header>

      <div v-if="loading" class="admin-loading">
        <i class="pi pi-spin pi-spinner" aria-hidden="true"></i>
        <p>Memuat data panel...</p>
      </div>

      <template v-else>
        <div class="admin-tabs-wrap">
          <nav class="admin-tabs" role="tablist" aria-label="Bagian panel admin">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            @click="activeTab = tab.key"
            :class="{ 'is-active': activeTab === tab.key }"
            role="tab"
            :aria-selected="activeTab === tab.key"
            type="button"
          >
            {{ tab.label }}
          </button>
          </nav>
        </div>

      <!-- Users Tab -->
      <section v-if="activeTab === 'users'" class="admin-section">
        <div class="admin-section-heading">
          <div>
            <h2>Manajemen Pengguna</h2>
            <p>Atur akun, profesi, akses, dan status pengguna.</p>
          </div>
          <span class="admin-count"><i class="pi pi-users" aria-hidden="true"></i>{{ usersTotal }} pengguna</span>
        </div>
        <div class="admin-filter-bar">
          <label class="admin-search">
            <i class="pi pi-search" aria-hidden="true"></i>
            <InputText v-model="userSearch" placeholder="Cari nama, username, atau role..." class="text-sm" @input="debouncedFetchUsers" />
          </label>
          <div class="admin-select-filter">
            <i class="pi pi-filter" aria-hidden="true"></i>
            <Select
              v-model="userRoleFilter"
              :options="roleFilterOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Semua Role"
              :pt="{ overlay: { class: 'admin-role-dropdown' } }"
              class="admin-filter-select text-sm"
              aria-label="Filter berdasarkan role"
              @change="resetUserPage"
            />
          </div>
        </div>

        <ModernTable
          :value="users"
          :loading="usersLoading"
          paginator
          :rows="userRows"
          :first="(userPage - 1) * userRows"
          :totalRecords="usersTotal"
          :rowsPerPageOptions="[10, 20, 50]"
          lazy
          dataKey="id_user"
          ariaLabel="Daftar pengguna"
          @page="onUserPageChange"
        >
          <template #header>
            <th scope="col">ID</th>
            <th scope="col">Username</th>
            <th scope="col">Nama</th>
            <th scope="col">Profesi</th>
            <th scope="col">Spesialisasi</th>
            <th scope="col">No. STR</th>
            <th scope="col">Role</th>
            <th scope="col">Unit</th>
            <th scope="col">Status</th>
            <th scope="col">Aksi</th>
          </template>
          <template #row="{ item }">
            <td class="font-mono">{{ item.id_user }}</td>
            <td class="font-mono">{{ item.username }}</td>
            <td class="user-name-cell">{{ item.nama }}</td>
            <td>{{ item.profesi }}</td>
            <td>{{ item.spesialisasi || '-' }}</td>
            <td class="font-mono">{{ item.no_str }}</td>
            <td class="role-cell">
              <Select
                v-model="item.role"
                :options="roleNameOptions"
                optionLabel="label"
                optionValue="value"
                :pt="{ overlay: { class: 'admin-role-dropdown' } }"
                class="admin-role-select w-full text-sm"
                :disabled="usersLoading"
                @change="onUserRoleChange(item)"
              />
            </td>
            <td>{{ item.unit_nama || '-' }}</td>
            <td><StatusBadge :status="item.deleted_at ? 'nonaktif' : 'aktif'" type="user" /></td>
            <td>
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" severity="secondary" text rounded @click="editUser(item)" :disabled="usersLoading" aria-label="Edit pengguna" />
                <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDeleteUser(item)" :disabled="usersLoading" aria-label="Hapus pengguna" />
              </div>
            </td>
          </template>
        </ModernTable>

      </section>

      <!-- Claims Tab -->
      <section v-if="activeTab === 'claims'" class="admin-section">
        <div class="admin-section-heading">
          <div>
            <h2>Manajemen Klaim</h2>
            <p>Tinjau pengajuan dan perbarui status klaim.</p>
          </div>
          <span class="admin-count"><i class="pi pi-file-check" aria-hidden="true"></i>{{ claimsTotal }} klaim</span>
        </div>
        <div class="admin-filter-bar">
          <label class="admin-search">
            <i class="pi pi-search" aria-hidden="true"></i>
            <InputText v-model="claimSearch" placeholder="Cari nama pasien..." class="text-sm" @input="debouncedFetchClaims" />
          </label>
          <div class="admin-select-filter">
            <i class="pi pi-filter" aria-hidden="true"></i>
            <Select
              v-model="claimStatusFilter"
              :options="claimStatusOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Semua Status"
              class="admin-filter-select text-sm"
              aria-label="Filter berdasarkan status klaim"
              @change="resetClaimPage"
            />
          </div>
          <DatePicker
            v-model="claimDateFrom"
            placeholder="Dari tanggal"
            dateFormat="dd/mm/yy"
            showIcon
            class="admin-date-filter"
            aria-label="Tanggal klaim mulai"
          />
          <DatePicker
            v-model="claimDateTo"
            placeholder="Sampai tanggal"
            dateFormat="dd/mm/yy"
            showIcon
            class="admin-date-filter"
            aria-label="Tanggal klaim sampai"
          />
          <button
            v-if="authStore.user?.role === 'keuangan'"
            @click="exportClaimsCSV"
            class="admin-secondary-action"
            type="button"
          >
            <i class="pi pi-download" aria-hidden="true"></i>
            Ekspor CSV
          </button>
        </div>

        <ModernTable
          :value="claims"
          :loading="claimsLoading"
          paginator
          :rows="claimRows"
          :first="(claimPage - 1) * claimRows"
          :totalRecords="claimsTotal"
          :rowsPerPageOptions="[10, 20, 50]"
          lazy
          dataKey="id_klaim"
          ariaLabel="Daftar klaim"
          @page="onClaimPageChange"
        >
          <template #header>
            <th scope="col">ID Klaim</th>
            <th scope="col">No. RM</th>
            <th scope="col">Pasien</th>
            <th scope="col">Kode CBGs</th>
            <th scope="col">Status</th>
            <th scope="col">Nominal</th>
            <th scope="col">Tanggal Klaim</th>
            <th scope="col">Tanggal Update</th>
            <th scope="col">Petugas</th>
            <th scope="col">Aksi</th>
          </template>
          <template #row="{ item }">
            <td class="font-mono">{{ item.id_klaim }}</td>
            <td class="font-mono">{{ item.id_rekam }}</td>
            <td class="user-name-cell">{{ item.pasien_nama }}</td>
            <td class="font-mono">{{ item.kode_cbgs }}</td>
            <td><StatusBadge :status="item.status_klaim" type="claim" /></td>
            <td class="font-mono amount-cell">{{ formatNumber(item.nominal_klaim) }}</td>
            <td class="font-mono">{{ formatDate(item.tanggal_klaim) }}</td>
            <td class="font-mono">{{ formatDate(item.updated_at) }}</td>
            <td>{{ item.petugas_nama }}</td>
            <td>
                <div class="flex gap-1">
                  <Button
                    v-if="item.status_klaim === 'draft' || item.status_klaim === 'pending'"
                    icon="pi pi-calculator"
                    severity="primary"
                    text
                    rounded
                    @click="openScoreModal(item)"
                    :disabled="claimsLoading"
                    aria-label="Skoring klaim"
                  />
                  <span v-else class="text-secondary text-sm self-center">-</span>
                </div>
            </td>
          </template>
        </ModernTable>

      </section>

        <!-- Audit Trail Tab -->
        <section v-if="activeTab === 'audit'" class="admin-section">
          <div class="admin-section-heading">
            <div>
              <h2>Audit Trail</h2>
              <p>Lacak aktivitas penting dan perubahan data di sistem.</p>
            </div>
            <span class="admin-count"><i class="pi pi-history" aria-hidden="true"></i>{{ auditTotal }} aktivitas</span>
          </div>
          <div class="admin-filter-bar">
            <div class="admin-select-filter">
              <i class="pi pi-user" aria-hidden="true"></i>
              <Select
                v-model="auditUserFilter"
                :options="auditUserOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Semua User"
                class="admin-filter-select"
                aria-label="Filter berdasarkan pengguna"
                @change="resetAuditPage"
              />
            </div>
            <div class="admin-select-filter">
              <i class="pi pi-filter" aria-hidden="true"></i>
              <Select
                v-model="auditActivityFilter"
                :options="auditActivityOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Semua Aktivitas"
                class="admin-filter-select"
                aria-label="Filter berdasarkan aktivitas"
                @change="resetAuditPage"
              />
            </div>
            <DatePicker
              v-model="auditDateFrom"
              placeholder="Dari tanggal"
              dateFormat="dd/mm/yy"
              showIcon
              class="admin-date-filter"
              aria-label="Aktivitas sejak tanggal"
            />
            <DatePicker
              v-model="auditDateTo"
              placeholder="Sampai tanggal"
              dateFormat="dd/mm/yy"
              showIcon
              class="admin-date-filter"
              aria-label="Aktivitas sampai tanggal"
            />
          </div>

          <ModernTable
            :value="auditTrail"
            :loading="auditLoading"
            paginator
            :rows="auditRows"
            :first="(auditPage - 1) * auditRows"
            :totalRecords="auditTotal"
            :rowsPerPageOptions="[20, 50, 100]"
            lazy
            ariaLabel="Riwayat aktivitas sistem"
            @page="onAuditPageChange"
          >
            <template #header>
              <th scope="col">Waktu</th>
              <th scope="col">User</th>
              <th scope="col">Aktivitas</th>
              <th scope="col">Tabel</th>
              <th scope="col">Data ID</th>
            </template>
            <template #row="{ item }">
              <td class="font-mono">{{ formatDateTime(item.waktu) }}</td>
              <td class="font-mono user-name-cell">{{ item.username || '-' }}</td>
              <td><StatusBadge :status="item.aktivitas" type="activity" /></td>
              <td>{{ item.tabel_terdampak }}</td>
              <td class="font-mono">{{ item.data_id || '-' }}</td>
            </template>
          </ModernTable>

        </section>

        <!-- Units Tab -->
        <section v-if="activeTab === 'units'" class="admin-section">
          <div class="admin-section-heading">
            <div>
              <h2>Manajemen Unit</h2>
              <p>Kelola unit kerja yang tersedia untuk pengguna.</p>
            </div>
            <span class="admin-count"><i class="pi pi-building" aria-hidden="true"></i>{{ unitsTotal }} unit</span>
          </div>

          <ModernTable
            :value="units"
            :loading="unitsLoading"
            paginator
            :rows="unitRows"
            :first="(unitPage - 1) * unitRows"
            :totalRecords="unitsTotal"
            :rowsPerPageOptions="[10, 20, 50]"
            dataKey="id_unit"
            ariaLabel="Daftar unit"
            @page="onUnitPageChange"
          >
            <template #header>
              <th scope="col">ID</th>
              <th scope="col">Nama Unit</th>
              <th scope="col">Aksi</th>
            </template>
            <template #row="{ item }">
              <td class="font-mono">{{ item.id_unit }}</td>
              <td class="user-name-cell">{{ item.nama_unit }}</td>
              <td>
                <div class="flex gap-1">
                  <Button icon="pi pi-pencil" severity="secondary" text rounded @click="editUnit(item)" :disabled="unitsLoading" aria-label="Edit unit" />
                  <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDeleteUnit(item)" :disabled="unitsLoading" aria-label="Hapus unit" />
                </div>
              </td>
            </template>
          </ModernTable>

        </section>

        <!-- Roles Tab (Read-only) -->
        <section v-if="activeTab === 'roles'" class="admin-section">
          <div class="admin-section-heading">
            <div>
              <h2>Manajemen Role</h2>
              <p>Daftar role yang mengatur akses di dalam sistem.</p>
            </div>
            <span class="admin-count"><i class="pi pi-key" aria-hidden="true"></i>{{ roles.length }} role</span>
          </div>
          <p class="admin-info-note"><i class="pi pi-info-circle" aria-hidden="true"></i> Penugasan role dilakukan melalui halaman Manajemen Pengguna.</p>

          <ModernTable
            :value="roles"
            :loading="rolesLoading"
            dataKey="id_role"
            ariaLabel="Daftar role"
          >
            <template #header>
              <th scope="col">ID</th>
              <th scope="col">Nama Role</th>
              <th scope="col">Deskripsi</th>
            </template>
            <template #row="{ item }">
              <td class="font-mono">{{ item.id_role }}</td>
              <td class="user-name-cell">{{ item.nama_role }}</td>
              <td class="description-cell">{{ item.deskripsi || '-' }}</td>
            </template>
          </ModernTable>
        </section>
      </template>

      <UserFormModal
        :show="showUserModal"
        :user="editingUser"
        :roles="roles"
        :units="units"
        :saving="userSaving"
        @update:show="showUserModal = $event"
        @saved="saveUser"
      />

      <!-- Unit Form Modal -->
      <UnitFormModal
        :show="showUnitModal"
        @update:show="showUnitModal = $event"
        :unit="editingUnit"
        @saved="saveUnit"
      />

      <!-- Delete Confirm Modal -->
      <DeleteConfirm
        :show="showDeleteConfirm"
        @update:show="showDeleteConfirm = $event"
        :title="deleteConfirmTitle"
        :message="deleteConfirmMessage"
        @confirm="executeDelete"
      />

      <!-- Score Modal -->
      <div v-if="selectedKlaim" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-surface rounded-lg shadow-xl max-w-lg w-full p-6">
          <h3 class="text-lg font-semibold text-primary mb-4">Skoring Klaim #{{ selectedKlaim.id_klaim }}</h3>
          <p class="text-sm text-secondary mb-4">Pasien: {{ selectedKlaim.pasien_nama }}</p>

          <div class="mb-4">
            <label class="block text-sm font-medium text-secondary mb-1">Status Baru</label>
            <Select v-model="newStatus" :options="claimStatusOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>

          <div class="mb-4" v-if="newStatus === 'pending' || newStatus === 'ditolak'">
            <label class="block text-sm font-medium text-secondary mb-1">Alasan (Wajib)</label>
            <textarea v-model="alasan" rows="3" class="w-full px-3 py-2 border border-border rounded focus:outline-none focus:ring-2 focus:ring-primary" placeholder="Masukkan alasan pending/ditolak"></textarea>
          </div>

          <div class="flex justify-end gap-3">
            <Button label="Batal" severity="secondary" @click="closeScoreModal" />
            <Button label="Simpan" @click="submitScore" :loading="scoreLoading" />
          </div>
        </div>
      </div>

      <Toast />
    </div>
  </DefaultLayout>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '../stores/auth'
import { useAdminStore } from '../stores/adminStore'
import client from '../api/client'

import DefaultLayout from '../layouts/DefaultLayout.vue'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Button from 'primevue/button'
import Toast from 'primevue/toast'

import StatusBadge from '../components/StatusBadge.vue'
import ModernTable from '../components/ModernTable.vue'
import UserFormModal from '../components/UserFormModal.vue'
import UnitFormModal from '../components/UnitFormModal.vue'
import DeleteConfirm from '../components/DeleteConfirm.vue'

const authStore = useAuthStore()
const adminStore = useAdminStore()
const toast = useToast()

const tabs = [
  { key: 'users', label: 'Pengguna' },
  { key: 'claims', label: 'Klaim' },
  { key: 'audit', label: 'Audit Trail' },
  { key: 'units', label: 'Unit' },
  { key: 'roles', label: 'Role' }
]

const activeTab = ref('users')
const loading = ref(true)

// Users
const users = ref([])
const usersTotal = ref(0)
const usersLoading = ref(false)
const userSearch = ref('')
const userRoleFilter = ref(null)
const roleFilterOptions = computed(() => [
  { label: 'Semua Role', value: null },
  ...roles.value.map(r => ({ label: formatRoleName(r.nama_role), value: r.nama_role }))
])
const roleNameOptions = computed(() => roles.value.map(r => ({
  label: formatRoleName(r.nama_role),
  value: r.nama_role
})))
const showUserModal = ref(false)
const editingUser = ref(null)
const userSaving = ref(false)
const userPage = ref(1)
const userRows = ref(10)
let userSearchTimeout = null

// Claims
const claims = ref([])
const claimsTotal = ref(0)
const claimsLoading = ref(false)
const claimStatusFilter = ref('')
const claimSearch = ref('')
const claimDateFrom = ref(null)
const claimDateTo = ref(null)
const claimStatusOptions = computed(() => [
  { label: 'Semua Status', value: '' },
  { label: 'Draft', value: 'draft' },
  { label: 'Pending', value: 'pending' },
  { label: 'Disetujui', value: 'disetujui' },
  { label: 'Ditolak', value: 'ditolak' }
])
const selectedKlaim = ref(null)
const newStatus = ref('pending')
const alasan = ref('')
const scoreLoading = ref(false)
const claimPage = ref(1)
const claimRows = ref(10)
let claimSearchTimeout = null

// Audit Trail
const auditTrail = ref([])
const auditTotal = ref(0)
const auditLoading = ref(false)
const auditUserFilter = ref(null)
const auditActivityFilter = ref('')
const auditDateFrom = ref(null)
const auditDateTo = ref(null)
const auditUserOptions = computed(() => [
  { label: 'Semua User', value: null },
  ...adminStore.users.map(u => ({ label: u.nama, value: u.id_user }))
  ])
const auditActivityOptions = [
  { label: 'Semua Aktivitas', value: '' },
  { label: 'INSERT', value: 'INSERT' },
  { label: 'UPDATE', value: 'UPDATE' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'LOGIN', value: 'LOGIN' },
  { label: 'LOGOUT', value: 'LOGOUT' }
]
const auditPage = ref(1)
const auditRows = ref(20)

// Units
const units = ref([])
const unitsTotal = ref(0)
const unitsLoading = ref(false)
const showUnitModal = ref(false)
const editingUnit = ref(null)
const unitPage = ref(1)
const unitRows = ref(10)

// Roles
const roles = ref([])
const rolesLoading = ref(false)

// Delete Confirm
const showDeleteConfirm = ref(false)
const deleteConfirmTitle = ref('Konfirmasi Hapus')
const deleteConfirmMessage = ref('')
const deleteTarget = ref(null)
const deleteType = ref('')

// Watch activeTab to fetch data
watch(activeTab, async (newTab) => {
  switch (newTab) {
    case 'users':
      if (users.value.length === 0) await fetchUsers()
      break
    case 'claims':
      if (claims.value.length === 0) await fetchClaims()
      break
    case 'audit':
      if (auditTrail.value.length === 0) await fetchAuditTrail()
      break
    case 'units':
      if (units.value.length === 0) await fetchUnits()
      break
    case 'roles':
      if (roles.value.length === 0) await fetchRoles()
      break
  }
})

watch([claimDateFrom, claimDateTo], resetClaimPage)
watch([auditDateFrom, auditDateTo], resetAuditPage)

// Initial load
onMounted(async () => {
  await Promise.all([fetchRoles(), fetchUnits()])
  await fetchUsers()
  loading.value = false
})

// ============ Users ============
async function fetchUsers() {
  usersLoading.value = true
  try {
    const params = { page: userPage.value, limit: userRows.value }
    if (userSearch.value) params.search = userSearch.value
    if (userRoleFilter.value) params.role = userRoleFilter.value
    await adminStore.fetchUsers(params)
    users.value = adminStore.users
    usersTotal.value = adminStore.usersTotal
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data pengguna', life: 3000 })
  } finally {
    usersLoading.value = false
  }
}

function onUserRoleChange(user) {
  // Inline role reassignment persists through the same user update endpoint
  const role = roles.value.find(r => r.nama_role === user.role)
  if (!role || role.id_role === user.id_role) return

  const payload = {
    nama: user.nama,
    profesi: user.profesi,
    spesialisasi: user.spesialisasi || '',
    no_str: user.no_str,
    id_role: role.id_role,
    id_unit: user.id_unit
  }

  adminStore.updateUser(user.id_user, payload)
    .then(result => {
      if (!result.success) {
        toast.add({ severity: 'error', summary: 'Error', detail: result.message || 'Gagal mengubah role', life: 3000 })
      } else {
        toast.add({ severity: 'success', summary: 'Sukses', detail: `Role diubah menjadi ${user.role}`, life: 3000 })
      }
      fetchUsers()
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengubah role', life: 3000 })
      fetchUsers()
    })
}

function debouncedFetchUsers() {
  clearTimeout(userSearchTimeout)
  userSearchTimeout = setTimeout(resetUserPage, 300)
}

function onUserPageChange(event) {
  userPage.value = event.page + 1
  userRows.value = event.rows
  fetchUsers()
}

function resetUserPage() {
  userPage.value = 1
  fetchUsers()
}

function openUserModal() {
  revealUserForm(null)
}

function editUser(user) {
  revealUserForm(user)
}

function revealUserForm(user) {
  editingUser.value = user
  showUserModal.value = true
}

async function saveUser(payload) {
  userSaving.value = true
  try {
    let result
    if (editingUser.value) {
      result = await adminStore.updateUser(editingUser.value.id_user, payload)
    } else {
      result = await adminStore.createUser(payload)
    }
    if (result.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editingUser.value ? 'Pengguna diperbarui' : 'Pengguna ditambahkan', life: 3000 })
      showUserModal.value = false
      await fetchUsers()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: result.message || 'Gagal menyimpan pengguna', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menyimpan pengguna', life: 3000 })
  } finally {
    userSaving.value = false
  }
}

function confirmDeleteUser(user) {
  deleteTarget.value = user
  deleteType.value = 'user'
  deleteConfirmTitle.value = 'Hapus Pengguna'
  deleteConfirmMessage.value = `Apakah Anda yakin ingin menghapus pengguna "${user.nama}" (${user.username})? Data akan di-soft delete.`
  showDeleteConfirm.value = true
}

// ============ Claims ============
async function fetchClaims() {
  claimsLoading.value = true
  try {
    const params = { page: claimPage.value, limit: claimRows.value }
    if (claimStatusFilter.value) params.status = claimStatusFilter.value
    if (claimSearch.value) params.search = claimSearch.value
    if (claimDateFrom.value) params.dari = formatDateForApi(claimDateFrom.value)
    if (claimDateTo.value) params.sampai = formatDateForApi(claimDateTo.value)
    await adminStore.fetchClaims(params)
    claims.value = adminStore.claims
    claimsTotal.value = adminStore.claimsTotal
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data klaim', life: 3000 })
  } finally {
    claimsLoading.value = false
  }
}

function debouncedFetchClaims() {
  clearTimeout(claimSearchTimeout)
  claimSearchTimeout = setTimeout(resetClaimPage, 300)
}

function onClaimPageChange(event) {
  claimPage.value = event.page + 1
  claimRows.value = event.rows
  fetchClaims()
}

function resetClaimPage() {
  claimPage.value = 1
  fetchClaims()
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
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Alasan wajib diisi untuk status pending atau ditolak', life: 3000 })
    return
  }

  scoreLoading.value = true
  try {
    const result = await adminStore.updateClaimStatus(selectedKlaim.value.id_klaim, {
      status_klaim: newStatus.value,
      alasan_pending_tolak: alasan.value.trim() || null
    })
    if (result.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Status klaim diperbarui', life: 3000 })
      closeScoreModal()
      await fetchClaims()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: result.message || 'Gagal memperbarui status', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui status', life: 3000 })
  } finally {
    scoreLoading.value = false
  }
}

async function exportClaimsCSV() {
  try {
    const res = await client.get('/admin/submissions/export', {
      params: {
        status: claimStatusFilter.value,
        dari: claimDateFrom.value ? formatDateForApi(claimDateFrom.value) : null,
        sampai: claimDateTo.value ? formatDateForApi(claimDateTo.value) : null,
        search: claimSearch.value
      },
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `klaim_export_${new Date().toISOString().split('T')[0]}.csv`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengekspor data', life: 3000 })
  }
}

// ============ Audit Trail ============
async function fetchAuditTrail() {
  auditLoading.value = true
  try {
    const params = { page: auditPage.value, limit: auditRows.value }
    if (auditUserFilter.value) params.user_id = auditUserFilter.value
    if (auditActivityFilter.value) params.aktivitas = auditActivityFilter.value
    if (auditDateFrom.value) params.dari = formatDateForApi(auditDateFrom.value)
    if (auditDateTo.value) params.sampai = formatDateForApi(auditDateTo.value)
    await adminStore.fetchAuditTrail(params)
    auditTrail.value = adminStore.auditTrail
    auditTotal.value = adminStore.auditTotal
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat audit trail', life: 3000 })
  } finally {
    auditLoading.value = false
  }
}

function onAuditPageChange(event) {
  auditPage.value = event.page + 1
  auditRows.value = event.rows
  fetchAuditTrail()
}

function resetAuditPage() {
  auditPage.value = 1
  fetchAuditTrail()
}

// ============ Units ============
async function fetchUnits() {
  unitsLoading.value = true
  try {
    await adminStore.fetchUnits()
    units.value = adminStore.units
    unitsTotal.value = adminStore.units.length
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data unit', life: 3000 })
  } finally {
    unitsLoading.value = false
  }
}

function onUnitPageChange(event) {
  unitPage.value = event.page + 1
  unitRows.value = event.rows
}

function openUnitModal() {
  editingUnit.value = null
  showUnitModal.value = true
}

function editUnit(unit) {
  editingUnit.value = unit
  showUnitModal.value = true
}

async function saveUnit(payload) {
  try {
    let result
    if (editingUnit.value) {
      result = await adminStore.updateUnit(editingUnit.value.id_unit, payload)
    } else {
      result = await adminStore.createUnit(payload)
    }
    if (result.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editingUnit.value ? 'Unit diperbarui' : 'Unit ditambahkan', life: 3000 })
      await fetchUnits()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: result.message || 'Gagal menyimpan unit', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menyimpan unit', life: 3000 })
  }
}

function confirmDeleteUnit(unit) {
  deleteTarget.value = unit
  deleteType.value = 'unit'
  deleteConfirmTitle.value = 'Hapus Unit'
  deleteConfirmMessage.value = `Apakah Anda yakin ingin menghapus unit "${unit.nama_unit}"?`
  showDeleteConfirm.value = true
}

// ============ Roles ============
async function fetchRoles() {
  rolesLoading.value = true
  try {
    await adminStore.fetchRoles()
    roles.value = adminStore.roles
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data role', life: 3000 })
  } finally {
    rolesLoading.value = false
  }
}

// ============ Delete Execute ============
async function executeDelete() {
  try {
    let result
    if (deleteType.value === 'user') {
      result = await adminStore.deleteUser(deleteTarget.value.id_user)
    } else if (deleteType.value === 'unit') {
      result = await adminStore.deleteUnit(deleteTarget.value.id_unit)
    }
    if (result.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Data dihapus', life: 3000 })
      if (deleteType.value === 'user') await fetchUsers()
      if (deleteType.value === 'unit') await fetchUnits()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: result.message || 'Gagal menghapus', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus', life: 3000 })
  }
  showDeleteConfirm.value = false
  deleteTarget.value = null
  deleteType.value = ''
}

// ============ Helpers ============
function formatNumber(num) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num)
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function formatDateTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatDateForApi(date) {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatRoleName(role) {
  const acronyms = { ti: 'TI', rm: 'RM', cbgs: 'CBGs' }
  return role
    .replace(/[_-]+/g, ' ')
    .split(' ')
    .filter(Boolean)
    .map(word => acronyms[word.toLowerCase()] || `${word[0].toUpperCase()}${word.slice(1).toLowerCase()}`)
    .join(' ')
}
</script>

<style scoped>
.admin-page {
  display: flex;
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
  color: var(--color-text-primary);
}

.admin-page > header,
.admin-page > .admin-tabs-wrap {
  flex: 0 0 auto;
}

.admin-hero {
  position: relative;
  display: flex;
  min-height: 104px;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  overflow: hidden;
  padding: 18px 24px;
  border: 1px solid rgb(255 255 255 / 8%);
  border-radius: 16px;
  background: linear-gradient(115deg, #102b43 0%, #1c3d5a 62%, #315875 100%);
  box-shadow: 0 12px 28px rgb(19 48 72 / 12%);
  color: white;
}

.admin-hero::after {
  position: absolute;
  top: -90px;
  right: 10%;
  width: 250px;
  height: 250px;
  border: 1px solid rgb(255 255 255 / 10%);
  border-radius: 50%;
  box-shadow: 0 0 0 35px rgb(255 255 255 / 3%), 0 0 0 70px rgb(255 255 255 / 2%);
  content: '';
  pointer-events: none;
}

.admin-hero-copy,
.admin-hero-actions {
  position: relative;
  z-index: 1;
}

.admin-eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  color: #c6d8e5;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.14em;
}

.admin-eyebrow i {
  color: #a8c9dd;
}

.admin-hero h1 {
  margin: 0;
  color: white;
  font-size: clamp(24px, 3vw, 30px);
  font-weight: 600;
  letter-spacing: -0.035em;
  line-height: 1.2;
}

.admin-hero p {
  margin: 8px 0 0;
  color: #d2dee8;
  font-size: 12px;
}

.admin-primary-action,
.admin-secondary-action {
  display: inline-flex;
  min-height: 42px;
  align-items: center;
  justify-content: center;
  gap: 9px;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  transition: background-color 140ms ease, border-color 140ms ease, transform 140ms ease;
  white-space: nowrap;
}

.admin-primary-action {
  padding: 0 16px;
  background: white;
  color: var(--color-primary);
  box-shadow: 0 3px 10px rgb(8 25 40 / 14%);
}

.admin-primary-action:hover {
  transform: translateY(-1px);
  background: #f1f6f9;
}

.admin-secondary-action {
  padding: 0 13px;
  border-color: var(--color-border);
  background: white;
  color: var(--color-text-primary);
}

.admin-secondary-action:hover {
  border-color: #b9c8d3;
  background: #f7fafc;
}

.admin-tabs-wrap {
  overflow-x: auto;
  padding: 5px;
  border: 1px solid var(--color-border);
  border-radius: 11px;
  background: #eeede9;
}

.admin-tabs {
  display: flex;
  min-width: max-content;
  gap: 4px;
}

.admin-tabs button {
  min-height: 39px;
  padding: 0 16px;
  border: 0;
  border-radius: 7px;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  font-weight: 500;
  transition: background-color 140ms ease, color 140ms ease, box-shadow 140ms ease;
}

.admin-tabs button:hover:not(.is-active) {
  background: rgb(255 255 255 / 55%);
  color: var(--color-text-primary);
}

.admin-tabs button.is-active {
  background: white;
  box-shadow: 0 1px 3px rgb(31 41 55 / 10%);
  color: var(--color-primary);
  font-weight: 600;
}

.admin-section {
  display: flex;
  width: 100%;
  flex: 1 1 auto;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
  animation: admin-fade-in 180ms ease-out both;
}

.admin-section-heading {
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.admin-section-heading h2 {
  margin: 0;
  color: var(--color-text-primary);
  font-size: 19px;
  font-weight: 600;
  letter-spacing: -0.02em;
}

.admin-section-heading p {
  margin: 4px 0 0;
  color: var(--color-text-secondary);
  font-size: 12px;
}

.admin-count {
  display: inline-flex;
  min-height: 30px;
  align-items: center;
  gap: 7px;
  padding: 0 10px;
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: white;
  color: #536474;
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
}

.admin-count i {
  color: #617f96;
}

.admin-filter-bar {
  display: flex;
  flex: 0 0 auto;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  background: white;
}

.admin-search {
  display: flex;
  min-width: min(100%, 270px);
  min-height: 38px;
  flex: 1 1 270px;
  align-items: center;
  gap: 9px;
  padding: 0 11px;
  border: 1px solid #dfe3e6;
  border-radius: 7px;
  background: #fbfcfc;
  color: #7b8790;
}

.admin-search:focus-within {
  border-color: #8ba6ba;
  box-shadow: 0 0 0 3px rgb(28 61 90 / 8%);
}

.admin-search :deep(.p-inputtext) {
  width: 100%;
  min-width: 0;
  padding: 8px 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

.admin-select-filter {
  display: flex;
  min-height: 38px;
  flex: 0 0 210px;
  align-items: center;
  gap: 9px;
  padding: 0 10px;
  border: 1px solid #dfe3e6;
  border-radius: 7px;
  background: #fbfcfc;
  color: #71808a;
  transition: border-color 140ms ease, box-shadow 140ms ease;
}

.admin-select-filter:focus-within {
  border-color: #8ba6ba;
  box-shadow: 0 0 0 3px rgb(28 61 90 / 8%);
}

.admin-select-filter > i {
  flex: 0 0 auto;
  font-size: 12px;
}

.admin-select-filter .admin-filter-select {
  width: 100%;
  min-width: 0;
}

.admin-select-filter :deep(.p-select) {
  width: 100%;
  min-height: 34px;
  border: 0;
  background: transparent;
  box-shadow: none;
  color: #344653;
  font-size: 12px;
}

.admin-select-filter :deep(.p-select-label) {
  padding: 7px 0;
}

.admin-select-filter :deep(.p-select-dropdown) {
  width: 28px;
}

.admin-filter-select {
  min-width: 150px;
}

.admin-date-filter {
  display: flex;
  width: 165px;
  min-width: 150px;
  min-height: 38px;
  flex: 0 0 165px;
}

.admin-filter-bar :deep(.admin-date-filter.p-datepicker) {
  overflow: hidden;
  border: 1px solid #dfe3e6;
  border-radius: 7px;
  background: #fbfcfc;
  transition: border-color 140ms ease, box-shadow 140ms ease;
}

.admin-filter-bar :deep(.admin-date-filter.p-datepicker:focus-within) {
  border-color: #8ba6ba;
  box-shadow: 0 0 0 3px rgb(28 61 90 / 8%);
}

.admin-filter-bar :deep(.admin-date-filter .p-datepicker-input) {
  width: 100%;
  min-width: 0;
  padding: 7px 10px;
  border: 0;
  background: #fbfcfc;
  box-shadow: none;
  color: #344653;
  font-size: 12px;
}

.admin-filter-bar :deep(.admin-date-filter .p-datepicker-dropdown) {
  width: 36px;
  border: 0;
  border-left: 1px solid #e5e9eb;
  border-radius: 0;
  background: #f4f7f8;
  color: #526777;
}

.admin-filter-bar :deep(.admin-date-filter .p-datepicker-dropdown:hover) {
  background: #eaf0f3;
  color: var(--color-primary);
}

.admin-info-note {
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 9px;
  margin: 0;
  padding: 11px 13px;
  border: 1px solid #dce6ec;
  border-radius: 8px;
  background: #f5f9fb;
  color: #526777;
  font-size: 12px;
}

.admin-info-note i {
  color: #527d99;
}

.admin-role-select :deep(.p-select) {
  min-height: 34px;
  border-color: #dce4e8;
  border-radius: 7px;
  background: #fff;
  color: #344653;
  font-size: 12px;
  transition: border-color 140ms ease, box-shadow 140ms ease;
}

.admin-role-select :deep(.p-select:not(.p-disabled):hover) {
  border-color: #aabdc9;
}

.admin-role-select :deep(.p-select:not(.p-disabled).p-focus) {
  border-color: #64849a;
  box-shadow: 0 0 0 3px rgb(28 61 90 / 9%);
}

.admin-loading {
  display: flex;
  flex: 1 1 auto;
  min-height: 0;
  min-height: 220px;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 10px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: white;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.admin-loading i {
  color: var(--color-primary);
  font-size: 22px;
}

.admin-loading p {
  margin: 0;
}

@keyframes admin-fade-in {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 640px) {
  .admin-page {
    gap: 10px;
  }

  .admin-hero {
    min-height: 92px;
    align-items: center;
    flex-direction: row;
    gap: 8px;
    padding: 12px 13px;
  }

  .admin-hero-actions {
    flex: 0 0 auto;
  }

  .admin-hero h1 {
    font-size: 21px;
  }

  .admin-hero p {
    display: none;
  }

  .admin-primary-action {
    min-height: 36px;
    gap: 6px;
    padding: 0 9px;
    font-size: 10px;
  }

  .admin-section-heading {
    align-items: flex-start;
    gap: 8px;
  }

  .admin-section-heading h2 {
    font-size: 16px;
  }

  .admin-section-heading p {
    max-width: 30ch;
    font-size: 10px;
  }

  .admin-count {
    margin-top: 2px;
    font-size: 10px;
  }

  .admin-filter-bar > * {
    width: 100%;
  }

  .admin-select-filter {
    flex-basis: auto;
  }

  .admin-date-filter {
    width: 100%;
    flex-basis: auto;
  }

  .admin-filter-bar {
    max-height: 122px;
    gap: 6px;
    overflow-y: auto;
    padding: 8px;
  }

  .admin-search {
    min-width: 100%;
    min-height: 34px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .admin-section {
    animation: none;
  }
}
</style>