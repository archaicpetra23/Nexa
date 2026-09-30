<template>
  <DefaultLayout>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-2xl font-semibold" style="color: var(--color-primary)">Panel Admin</h1>
      <div class="flex gap-2">
        <button
          v-if="activeTab === 'users'"
          @click="openUserModal"
          class="px-4 py-2 text-white rounded-base text-sm font-medium transition-colors duration-120"
          style="background-color: var(--color-primary)"
          @mouseenter="$event.target.style.backgroundColor = 'var(--color-primary-hover)'"
          @mouseleave="$event.target.style.backgroundColor = 'var(--color-primary)'"
        >
          Tambah Pengguna
        </button>
        <button
          v-if="activeTab === 'units'"
          @click="openUnitModal"
          class="px-4 py-2 text-white rounded-base text-sm font-medium transition-colors duration-120"
          style="background-color: var(--color-primary)"
          @mouseenter="$event.target.style.backgroundColor = 'var(--color-primary-hover)'"
          @mouseleave="$event.target.style.backgroundColor = 'var(--color-primary)'"
        >
          Tambah Unit
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p class="text-secondary text-sm">Memuat data...</p>
    </div>

    <template v-else>
      <!-- Tab Navigation -->
      <div class="border-b mb-6" style="border-color: var(--color-border)">
        <nav class="flex gap-1" role="tablist">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            @click="activeTab = tab.key"
            :class="[
              'px-4 py-3 text-sm font-medium rounded-base transition-colors duration-120',
              activeTab === tab.key
                ? 'text-white'
                : 'text-secondary'
            ]"
            :style="activeTab === tab.key
              ? { backgroundColor: 'var(--color-primary)' }
              : {}"
            role="tab"
            :aria-selected="activeTab === tab.key"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>

      <!-- Users Tab -->
      <section v-if="activeTab === 'users'" class="space-y-4">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-lg font-semibold" style="color: var(--color-primary)">Manajemen Pengguna</h2>
          <div class="flex gap-2">
            <InputText v-model="userSearch" placeholder="Cari nama, username, atau role..." class="w-64 text-sm" @input="debouncedFetchUsers" />
            <Select v-model="userRoleFilter" :options="roleFilterOptions" placeholder="Semua Role" class="w-40 text-sm" @change="fetchUsers" />
          </div>
        </div>

        <DataTable
          :value="users"
          :loading="usersLoading"
          :paginator="true"
          :rows="10"
          :totalRecords="usersTotal"
          :rowsPerPageOptions="[10, 20, 50]"
          :lazy="true"
          @page="onUserPageChange"
          responsiveLayout="scroll"
          selectionMode="single"
          :selection="selectedUser"
        >
          <Column field="id_user" header="ID" style="width: 6rem">
            <template #body="slotProps">
              <span class="font-mono text-sm">{{ slotProps.data.id_user }}</span>
            </template>
          </Column>
          <Column field="username" header="Username" style="width: 10rem">
            <template #body="slotProps">
              <span class="font-mono text-sm">{{ slotProps.data.username }}</span>
            </template>
          </Column>
          <Column field="nama" header="Nama" style="width: 14rem">
            <template #body="slotProps">
              <span class="text-sm">{{ slotProps.data.nama }}</span>
            </template>
          </Column>
          <Column field="profesi" header="Profesi" style="width: 10rem">
            <template #body="slotProps">
              <span class="text-sm">{{ slotProps.data.profesi }}</span>
            </template>
          </Column>
          <Column field="spesialisasi" header="Spesialisasi" style="width: 10rem">
            <template #body="slotProps">
              <span class="text-sm">{{ slotProps.data.spesialisasi || '-' }}</span>
            </template>
          </Column>
          <Column field="no_str" header="No. STR" style="width: 10rem">
            <template #body="slotProps">
              <span class="font-mono text-sm">{{ slotProps.data.no_str }}</span>
            </template>
          </Column>
          <Column field="role" header="Role" style="width: 12rem">
            <template #body="slotProps">
              <Select
                v-model="slotProps.data.role"
                :options="roleNameOptions"
                class="w-full text-sm"
                :disabled="usersLoading"
                @change="onUserRoleChange(slotProps.data)"
              />
            </template>
          </Column>
          <Column field="unit" header="Unit" style="width: 12rem">
            <template #body="slotProps">
              <span class="text-sm">{{ slotProps.data.unit_nama || '-' }}</span>
            </template>
          </Column>
          <Column field="status" header="Status" style="width: 8rem">
            <template #body="slotProps">
              <StatusBadge :status="slotProps.data.deleted_at ? 'nonaktif' : 'aktif'" type="user" />
            </template>
          </Column>
          <Column header="Aksi" style="width: 7rem">
            <template #body="slotProps">
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" severity="secondary" text @click="editUser(slotProps.data)" :disabled="usersLoading" aria-label="Edit pengguna" />
                <Button icon="pi pi-trash" severity="danger" text @click="confirmDeleteUser(slotProps.data)" :disabled="usersLoading" aria-label="Hapus pengguna" />
              </div>
            </template>
          </Column>
        </DataTable>

        <div v-if="users.length === 0 && !usersLoading" class="text-center py-8 text-secondary text-sm">
          Belum ada data pengguna. Klik "Tambah Pengguna" untuk menambah data baru.
        </div>
      </section>

      <!-- Claims Tab -->
      <section v-if="activeTab === 'claims'" class="space-y-4">
        <div class="flex flex-wrap gap-2 mb-4">
          <h2 class="text-lg font-semibold self-center" style="color: var(--color-primary)">Manajemen Klaim</h2>
          <Select v-model="claimStatusFilter" :options="claimStatusOptions" placeholder="Semua Status" class="w-40 text-sm" @change="fetchClaims" />
          <InputText v-model="claimSearch" placeholder="Cari nama pasien..." class="w-64 text-sm" @input="debouncedFetchClaims" />
          <div class="flex gap-2">
            <DatePicker v-model="claimDateFrom" placeholder="Dari" class="w-40 text-sm" @change="fetchClaims" />
            <DatePicker v-model="claimDateTo" placeholder="Sampai" class="w-40 text-sm" @change="fetchClaims" />
          </div>
          <button
            v-if="authStore.user?.role === 'keuangan'"
            @click="exportClaimsCSV"
            class="px-4 py-2 border rounded-base text-secondary hover:bg-base transition-colors duration-120 text-sm"
            style="border-color: var(--color-border)"
          >
            Ekspor CSV
          </button>
        </div>

        <DataTable
          :value="claims"
          :loading="claimsLoading"
          :paginator="true"
          :rows="10"
          :totalRecords="claimsTotal"
          :rowsPerPageOptions="[10, 20, 50]"
          :lazy="true"
          @page="onClaimPageChange"
          responsiveLayout="scroll"
          selectionMode="single"
          :selection="selectedClaim"
        >
          <Column field="id_klaim" header="ID Klaim" style="width: 8rem">
            <template #body="slotProps">
              <span class="font-mono text-sm">{{ slotProps.data.id_klaim }}</span>
            </template>
          </Column>
          <Column field="id_rekam" header="No. RM" style="width: 8rem">
            <template #body="slotProps">
              <span class="font-mono text-sm">{{ slotProps.data.id_rekam }}</span>
            </template>
          </Column>
          <Column field="pasien_nama" header="Pasien" style="width: 14rem">
            <template #body="slotProps">
              <span class="text-sm">{{ slotProps.data.pasien_nama }}</span>
            </template>
          </Column>
          <Column field="kode_cbgs" header="Kode CBGs" style="width: 8rem">
            <template #body="slotProps">
              <span class="font-mono text-sm">{{ slotProps.data.kode_cbgs }}</span>
            </template>
          </Column>
          <Column field="status_klaim" header="Status" style="width: 9rem">
            <template #body="slotProps">
              <StatusBadge :status="slotProps.data.status_klaim" type="claim" />
            </template>
          </Column>
          <Column field="nominal_klaim" header="Nominal" style="width: 10rem">
              <template #body="slotProps">
                <span class="font-mono text-sm text-right">{{ formatNumber(slotProps.data.nominal_klaim) }}</span>
              </template>
            </Column>
            <Column field="tanggal_klaim" header="Tanggal Klaim" style="width: 10rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ formatDate(slotProps.data.tanggal_klaim) }}</span>
              </template>
            </Column>
            <Column field="updated_at" header="Tanggal Update" style="width: 10rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ formatDate(slotProps.data.updated_at) }}</span>
              </template>
            </Column>
            <Column field="petugas_nama" header="Petugas" style="width: 12rem">
              <template #body="slotProps">
                <span class="text-sm">{{ slotProps.data.petugas_nama }}</span>
              </template>
            </Column>
            <Column header="Aksi" style="width: 7rem">
              <template #body="slotProps">
                <div class="flex gap-1">
                  <Button
                    v-if="slotProps.data.status_klaim === 'draft' || slotProps.data.status_klaim === 'pending'"
                    icon="pi pi-calculator"
                    severity="primary"
                    text
                    @click="openScoreModal(slotProps.data)"
                    :disabled="claimsLoading"
                    aria-label="Skoring klaim"
                  />
                  <span v-else class="text-secondary text-sm self-center">-</span>
                </div>
              </template>
            </Column>
          </DataTable>

          <div v-if="claims.length === 0 && !claimsLoading" class="text-center py-8 text-secondary">
            Belum ada data klaim.
          </div>
        </section>

        <!-- Audit Trail Tab -->
        <section v-if="activeTab === 'audit'" class="space-y-4">
          <h2 class="text-lg font-semibold text-primary">Audit Trail</h2>
          <div class="flex flex-wrap gap-2 mb-4">
            <Select v-model="auditUserFilter" :options="auditUserOptions" placeholder="Semua User" class="w-40" @change="fetchAuditTrail" />
            <Select v-model="auditActivityFilter" :options="auditActivityOptions" placeholder="Semua Aktivitas" class="w-40" @change="fetchAuditTrail" />
            <div class="flex gap-2">
              <DatePicker v-model="auditDateFrom" placeholder="Dari" class="w-40" @change="fetchAuditTrail" />
              <DatePicker v-model="auditDateTo" placeholder="Sampai" class="w-40" @change="fetchAuditTrail" />
            </div>
          </div>

          <DataTable
            :value="auditTrail"
            :loading="auditLoading"
            :paginator="true"
            :rows="20"
            :totalRecords="auditTotal"
            :rowsPerPageOptions="[20, 50, 100]"
            :lazy="true"
            @page="onAuditPageChange"
            responsiveLayout="scroll"
          >
            <Column field="waktu" header="Waktu" style="width: 12rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ formatDateTime(slotProps.data.waktu) }}</span>
              </template>
            </Column>
            <Column field="user" header="User" style="width: 10rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ slotProps.data.username || '-' }}</span>
              </template>
            </Column>
            <Column field="aktivitas" header="Aktivitas" style="width: 8rem">
              <template #body="slotProps">
                <StatusBadge :status="slotProps.data.aktivitas" type="activity" />
              </template>
            </Column>
            <Column field="tabel_terdampak" header="Tabel" style="width: 10rem">
              <template #body="slotProps">
                <span class="text-sm">{{ slotProps.data.tabel_terdampak }}</span>
              </template>
            </Column>
            <Column field="data_id" header="Data ID" style="width: 8rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ slotProps.data.data_id || '-' }}</span>
              </template>
            </Column>
          </DataTable>

          <div v-if="auditTrail.length === 0 && !auditLoading" class="text-center py-8 text-secondary">
            Belum ada data audit trail.
          </div>
        </section>

        <!-- Units Tab -->
        <section v-if="activeTab === 'units'" class="space-y-4">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-lg font-semibold text-primary">Manajemen Unit</h2>
          </div>

          <DataTable
            :value="units"
            :loading="unitsLoading"
            :paginator="true"
            :rows="10"
            :totalRecords="unitsTotal"
            :rowsPerPageOptions="[10, 20, 50]"
            :lazy="true"
            @page="onUnitPageChange"
            responsiveLayout="scroll"
          >
            <Column field="id_unit" header="ID" style="width: 5rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ slotProps.data.id_unit }}</span>
              </template>
            </Column>
            <Column field="nama_unit" header="Nama Unit" style="width: 20rem">
              <template #body="slotProps">
                <span class="text-sm">{{ slotProps.data.nama_unit }}</span>
              </template>
            </Column>
            <Column header="Aksi" style="width: 7rem">
              <template #body="slotProps">
                <div class="flex gap-1">
                  <Button icon="pi pi-pencil" severity="secondary" text @click="editUnit(slotProps.data)" :disabled="unitsLoading" aria-label="Edit unit" />
                  <Button icon="pi pi-trash" severity="danger" text @click="confirmDeleteUnit(slotProps.data)" :disabled="unitsLoading" aria-label="Hapus unit" />
                </div>
              </template>
            </Column>
          </DataTable>

          <div v-if="units.length === 0 && !unitsLoading" class="text-center py-8 text-secondary">
            Belum ada data unit. Klik "Tambah Unit" untuk menambah data baru.
          </div>
        </section>

        <!-- Roles Tab (Read-only) -->
        <section v-if="activeTab === 'roles'" class="space-y-4">
          <h2 class="text-lg font-semibold text-primary">Manajemen Role</h2>
          <p class="text-secondary text-sm mb-4">Penugasan role dilakukan melalui halaman Manajemen Pengguna (tab Pengguna).</p>

          <DataTable
            :value="roles"
            responsiveLayout="scroll"
            :paginator="false"
          >
            <Column field="id_role" header="ID" style="width: 5rem">
              <template #body="slotProps">
                <span class="font-mono text-sm">{{ slotProps.data.id_role }}</span>
              </template>
            </Column>
            <Column field="nama_role" header="Nama Role" style="width: 15rem">
              <template #body="slotProps">
                <span class="text-sm font-medium">{{ slotProps.data.nama_role }}</span>
              </template>
            </Column>
            <Column field="deskripsi" header="Deskripsi" style="width: 30rem">
              <template #body="slotProps">
                <span class="text-sm text-secondary">{{ slotProps.data.deskripsi || '-' }}</span>
              </template>
            </Column>
          </DataTable>
        </section>
      </template>

      <!-- User Form Modal -->
      <UserFormModal
        :show="showUserModal"
        @update:show="showUserModal = $event"
        :user="editingUser"
        :roles="roles"
        :units="units"
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
  </DefaultLayout>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '../stores/auth'
import { useAdminStore } from '../stores/adminStore'
import client from '../api/client'

import DefaultLayout from '../layouts/DefaultLayout.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Button from 'primevue/button'
import Toast from 'primevue/toast'

import StatusBadge from '../components/StatusBadge.vue'
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
  ...roles.value.map(r => ({ label: r.nama_role, value: r.nama_role }))
])
const roleNameOptions = computed(() => roles.value.map(r => ({ label: r.nama_role, value: r.nama_role })))
const selectedUser = ref(null)
const showUserModal = ref(false)
const editingUser = ref(null)
const userPage = ref(1)
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
const selectedClaim = ref(null)
const selectedKlaim = ref(null)
const newStatus = ref('pending')
const alasan = ref('')
const scoreLoading = ref(false)
const claimPage = ref(1)
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

// Units
const units = ref([])
const unitsTotal = ref(0)
const unitsLoading = ref(false)
const showUnitModal = ref(false)
const editingUnit = ref(null)
const unitPage = ref(1)

// Roles
const roles = ref([])

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
    const params = { page: userPage.value, limit: 10 }
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
  userSearchTimeout = setTimeout(fetchUsers, 300)
}

function onUserPageChange(event) {
  userPage.value = event.page + 1
  fetchUsers()
}

function openUserModal() {
  editingUser.value = null
  showUserModal.value = true
}

function editUser(user) {
  editingUser.value = user
  showUserModal.value = true
}

async function saveUser(payload) {
  try {
    let result
    if (editingUser.value) {
      result = await adminStore.updateUser(editingUser.value.id_user, payload)
    } else {
      result = await adminStore.createUser(payload)
    }
    if (result.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editingUser.value ? 'Pengguna diperbarui' : 'Pengguna ditambahkan', life: 3000 })
      await fetchUsers()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: result.message || 'Gagal menyimpan pengguna', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menyimpan pengguna', life: 3000 })
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
    const params = { page: claimPage.value, limit: 10 }
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
  claimSearchTimeout = setTimeout(fetchClaims, 300)
}

function onClaimPageChange(event) {
  claimPage.value = event.page + 1
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
    const params = { page: auditPage.value, limit: 20 }
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
  fetchUnits()
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
  try {
    await adminStore.fetchRoles()
    roles.value = adminStore.roles
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data role', life: 3000 })
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
  const d = new Date(date)
  return d.toISOString().split('T')[0]
}
</script>