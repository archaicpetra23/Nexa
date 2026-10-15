<template>
  <DefaultLayout>
    <div class="klaim-page">
      <header class="page-header">
        <div>
          <span class="page-kicker">PENAGIHAN</span>
          <h1>Klaim BPJS</h1>
          <p>Kelola pengajuan klaim dan pantau status persetujuan dari BPJS.</p>
        </div>
      </header>

      <div class="stat-strip">
        <div class="stat-card">
          <span class="stat-label">Pending</span>
          <span class="stat-value" :style="{ color: 'var(--color-status-pending)' }">{{ stats.pending }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-label">Disetujui</span>
          <span class="stat-value" :style="{ color: 'var(--color-status-approved)' }">{{ stats.disetujui }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-label">Ditolak</span>
          <span class="stat-value" :style="{ color: 'var(--color-status-rejected)' }">{{ stats.ditolak }}</span>
        </div>
      </div>

      <div class="filter-bar">
        <label class="search-input">
          <i class="pi pi-search" aria-hidden="true"></i>
          <InputText v-model="search" placeholder="Cari nomor RM atau pasien..." @input="debouncedSearch" />
        </label>
        <div class="filter-group">
          <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="Semua Status" @change="resetPage" class="w-full" />
        </div>
        <button @click="openCreateModal" class="primary-button" type="button">
          <i class="pi pi-plus" aria-hidden="true"></i>
          Ajukan Klaim
        </button>
      </div>

      <ModernTable
        :value="klaimList"
        :loading="loading"
        paginator
        lazy
        :rows="limit"
        :first="(page - 1) * limit"
        :totalRecords="total"
        :rowsPerPageOptions="[10, 20, 50]"
        dataKey="id_klaim"
        ariaLabel="Daftar klaim"
        @page="onPageChange"
      >
        <template #header>
          <th scope="col">ID Klaim</th>
          <th scope="col">No. RM</th>
          <th scope="col">Kode INA-CBGs</th>
          <th scope="col">Nominal</th>
          <th scope="col">Tgl. Klaim</th>
          <th scope="col">Status</th>
          <th scope="col">Aksi</th>
        </template>
        <template #row="{ item }">
          <td class="font-mono">{{ item.id_klaim }}</td>
          <td class="font-mono">{{ item.id_rekam }}</td>
          <td><span class="code-badge">{{ item.kode_cbgs }}</span></td>
          <td class="font-mono amount-cell">{{ formatRupiah(item.nominal_klaim) }}</td>
          <td>{{ formatDate(item.tanggal_klaim) }}</td>
          <td><StatusBadge :status="item.status_klaim" type="claim" /></td>
          <td>
            <div class="flex gap-1">
              <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openStatusModal(item)" aria-label="Ubah status klaim" />
              <Button v-if="canDelete" icon="pi pi-trash" severity="danger" text rounded @click="confirmDelete(item)" aria-label="Hapus klaim" />
            </div>
          </td>
        </template>
      </ModernTable>

      <!-- Create Modal -->
      <Dialog v-model:visible="showCreateModal" modal header="Ajukan Klaim Baru" :style="{ width: '36rem' }">
        <div class="form-grid">
          <div class="form-field">
            <label for="id_rekam">No. RM <span class="required">*</span></label>
            <InputNumber id="id_rekam" v-model="createForm.id_rekam" :min="1" class="w-full" />
            <small v-if="createErrors.id_rekam" class="field-error">{{ createErrors.id_rekam }}</small>
          </div>
          <div class="form-field">
            <label for="kode_cbgs">Kode INA-CBGs <span class="required">*</span></label>
            <Select id="kode_cbgs" v-model="createForm.kode_cbgs" :options="cbgsOptions" optionLabel="label" optionValue="value" placeholder="Pilih kode CBGs" class="w-full" @change="onCBGSChange" />
            <small v-if="createErrors.kode_cbgs" class="field-error">{{ createErrors.kode_cbgs }}</small>
            <small v-if="selectedCBGSTarif" class="field-hint">Tarif standar: {{ formatRupiah(selectedCBGSTarif) }}</small>
          </div>
          <div class="form-field full-width">
            <label for="nominal_klaim">Nominal <span class="required">*</span></label>
            <InputNumber id="nominal_klaim" v-model="createForm.nominal_klaim" :min="0" mode="currency" currency="IDR" locale="id-ID" class="w-full" />
            <small v-if="createErrors.nominal_klaim" class="field-error">{{ createErrors.nominal_klaim }}</small>
          </div>
        </div>
        <template #footer>
          <Button label="Batal" severity="secondary" @click="showCreateModal = false" />
          <Button label="Ajukan" @click="saveCreate" :loading="saving" :disabled="saving" />
        </template>
      </Dialog>

      <!-- Status Change Modal -->
      <Dialog v-model:visible="showStatusModal" modal header="Ubah Status Klaim" :style="{ width: '36rem' }">
        <div v-if="editingKlaim" class="form-grid">
          <div class="form-field full-width">
            <label>ID Klaim</label>
            <span class="static-value font-mono">{{ editingKlaim.id_klaim }}</span>
          </div>
          <div class="form-field">
            <label for="new-status">Status Baru <span class="required">*</span></label>
            <Select id="new-status" v-model="statusForm.status_klaim" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>
          <div v-if="showReasonField" class="form-field full-width">
            <label for="alasan">Alasan <span class="required">*</span></label>
            <Textarea id="alasan" v-model="statusForm.alasan_pending_tolak" rows="3" placeholder="Masukkan alasan..." class="w-full" />
            <small v-if="statusErrors.alasan" class="field-error">{{ statusErrors.alasan }}</small>
          </div>
        </div>
        <template #footer>
          <Button label="Batal" severity="secondary" @click="showStatusModal = false" />
          <Button label="Simpan" @click="saveStatus" :loading="saving" :disabled="saving" />
        </template>
      </Dialog>

      <DeleteConfirm
        :show="showDeleteConfirm"
        @update:show="showDeleteConfirm = $event"
        :message="`Apakah Anda yakin ingin menghapus klaim ID: ${deleteTarget?.id_klaim}?`"
        @confirm="executeDelete"
      />

      <Toast />
    </div>
  </DefaultLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '../stores/auth'
import client from '../api/client'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import ModernTable from '../components/ModernTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import DeleteConfirm from '../components/DeleteConfirm.vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import Toast from 'primevue/toast'

const authStore = useAuthStore()
const toast = useToast()

const klaimList = ref([])
const loading = ref(false)
const page = ref(1)
const limit = ref(10)
const total = ref(0)
const search = ref('')
const statusFilter = ref('')
let searchTimeout = null

const canDelete = computed(() => authStore.user?.role === 'admin_ti')

const stats = ref({ pending: 0, disetujui: 0, ditolak: 0 })

const statusOptions = [
  { label: 'Semua status', value: '' },
  { label: 'Draft', value: 'draft' },
  { label: 'Pending', value: 'pending' },
  { label: 'Disetujui', value: 'disetujui' },
  { label: 'Ditolak', value: 'ditolak' }
]

const cbgsOptions = ref([])

const showCreateModal = ref(false)
const createForm = ref({ id_rekam: null, kode_cbgs: '', nominal_klaim: 0 })
const createErrors = ref({})
const selectedCBGSTarif = ref(0)

const showStatusModal = ref(false)
const editingKlaim = ref(null)
const statusForm = ref({ status_klaim: '', alasan_pending_tolak: '' })
const statusErrors = ref({})
const showReasonField = computed(() => statusForm.value.status_klaim === 'pending' || statusForm.value.status_klaim === 'ditolak')

const showDeleteConfirm = ref(false)
const deleteTarget = ref(null)
const saving = ref(false)

async function fetchKlaim() {
  loading.value = true
  try {
    const params = { page: page.value, limit: limit.value }
    if (search.value) params.search = search.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await client.get('/klaim', { params })
    if (res.data.success) {
      klaimList.value = res.data.data.items || []
      total.value = res.data.data.total || 0
      updateStats()
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data klaim', life: 3000 })
  } finally {
    loading.value = false
  }
}

function updateStats() {
  stats.value = { pending: 0, disetujui: 0, ditolak: 0 }
  klaimList.value.forEach(k => {
    if (k.status_klaim === 'pending') stats.value.pending++
    else if (k.status_klaim === 'disetujui') stats.value.disetujui++
    else if (k.status_klaim === 'ditolak') stats.value.ditolak++
  })
}

async function fetchCBGS() {
  try {
    const res = await client.get('/master/cbgs', { params: { page: 1, limit: 1000 } })
    if (res.data.success) {
      cbgsOptions.value = (res.data.data.items || []).map(c => ({
        label: `${c.kode_cbgs} - ${c.deskripsi}`,
        value: c.kode_cbgs,
        tarif: c.tarif
      }))
    }
  } catch (e) {
    console.error('Gagal memuat CBGs:', e)
  }
}

function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchKlaim()
  }, 300)
}

function onPageChange(event) {
  page.value = event.page + 1
  limit.value = event.rows
  fetchKlaim()
}

function resetPage() {
  page.value = 1
  fetchKlaim()
}

function onCBGSChange() {
  const selected = cbgsOptions.value.find(c => c.value === createForm.value.kode_cbgs)
  selectedCBGSTarif.value = selected?.tarif || 0
  createForm.value.nominal_klaim = selectedCBGSTarif.value
}

function openCreateModal() {
  createForm.value = { id_rekam: null, kode_cbgs: '', nominal_klaim: 0 }
  createErrors.value = {}
  selectedCBGSTarif.value = 0
  showCreateModal.value = true
}

function validateCreateForm() {
  createErrors.value = {}
  if (!createForm.value.id_rekam) createErrors.value.id_rekam = 'No. RM wajib diisi'
  if (!createForm.value.kode_cbgs) createErrors.value.kode_cbgs = 'Kode CBGs wajib dipilih'
  if (!createForm.value.nominal_klaim || createForm.value.nominal_klaim < 0) createErrors.value.nominal_klaim = 'Nominal wajib diisi dan tidak negatif'
  return Object.keys(createErrors.value).length === 0
}

async function saveCreate() {
  if (!validateCreateForm()) return
  saving.value = true
  try {
    const res = await client.post('/klaim', createForm.value)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Klaim berhasil diajukan', life: 3000 })
      showCreateModal.value = false
      await fetchKlaim()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal mengajukan klaim', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal mengajukan klaim', life: 3000 })
  } finally {
    saving.value = false
  }
}

function openStatusModal(item) {
  editingKlaim.value = item
  statusForm.value = { status_klaim: item.status_klaim, alasan_pending_tolak: item.alasan_pending_tolak || '' }
  statusErrors.value = {}
  showStatusModal.value = true
}

function validateStatusForm() {
  statusErrors.value = {}
  if (!statusForm.value.status_klaim) statusErrors.value.status = 'Status wajib dipilih'
  if (showReasonField.value && !statusForm.value.alasan_pending_tolak?.trim()) {
    statusErrors.value.alasan = 'Alasan wajib diisi untuk status pending atau ditolak'
  }
  return Object.keys(statusErrors.value).length === 0
}

async function saveStatus() {
  if (!validateStatusForm()) return
  saving.value = true
  try {
    const payload = {
      status_klaim: statusForm.value.status_klaim,
      alasan_pending_tolak: statusForm.value.alasan_pending_tolak || null
    }
    const res = await client.patch(`/klaim/${editingKlaim.value.id_klaim}/status`, payload)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Status klaim diperbarui', life: 3000 })
      showStatusModal.value = false
      await fetchKlaim()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal memperbarui status', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal memperbarui status', life: 3000 })
  } finally {
    saving.value = false
  }
}

function confirmDelete(item) {
  deleteTarget.value = item
  showDeleteConfirm.value = true
}

async function executeDelete() {
  try {
    const res = await client.delete(`/klaim/${deleteTarget.value.id_klaim}`)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Klaim dihapus', life: 3000 })
      await fetchKlaim()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menghapus', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menghapus', life: 3000 })
  }
  showDeleteConfirm.value = false
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: '2-digit' })
}

function formatRupiah(value) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value)
}

fetchKlaim()
fetchCBGS()
</script>

<style scoped>
.klaim-page {
  display: grid;
  gap: 18px;
  animation: page-enter 240ms ease-out both;
}

.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
}

.page-kicker {
  display: block;
  color: var(--color-text-secondary);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1.25px;
}

.page-header h1 {
  margin: 8px 0 5px;
  color: var(--color-primary-dark);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.7px;
}

.page-header p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.stat-strip {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
}

.stat-label {
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 500;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  font-family: var(--font-mono);
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  background: white;
}

.search-input {
  display: flex;
  min-width: min(100%, 280px);
  min-height: 40px;
  flex: 1 1 280px;
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  border: 1px solid #dfe3e6;
  border-radius: 7px;
  background: #fbfcfc;
  color: #7b8790;
}

.search-input:focus-within {
  border-color: #8ba6ba;
  box-shadow: 0 0 0 3px rgb(28 61 90 / 8%);
}

.search-input :deep(.p-inputtext) {
  width: 100%;
  min-width: 0;
  padding: 9px 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

.filter-group {
  display: flex;
  min-width: 180px;
  align-items: center;
}

.primary-button {
  display: inline-flex;
  min-height: 40px;
  align-items: center;
  gap: 8px;
  padding: 0 16px;
  border: 0;
  border-radius: var(--radius-base);
  background: var(--color-primary);
  color: white;
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  transition: background-color 140ms ease;
}

.primary-button:hover {
  background: var(--color-primary-hover);
}

.code-badge {
  display: inline-block;
  padding: 4px 9px;
  border: 1px solid #dce4ea;
  border-radius: 5px;
  background: #f5f8fa;
  color: var(--color-primary-dark);
  font-size: 11px;
  font-weight: 600;
}

.amount-cell {
  text-align: right;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 10px;
}

.form-field {
  display: grid;
  gap: 6px;
}

.form-field.full-width {
  grid-column: 1 / -1;
}

.form-field label {
  color: var(--color-text-primary);
  font-size: 13px;
  font-weight: 500;
}

.static-value {
  padding: 8px 0;
  color: var(--color-text-primary);
  font-size: 13px;
}

.required {
  color: #B54245;
}

.field-error {
  color: #B54245;
  font-size: 11px;
}

.field-hint {
  color: var(--color-text-secondary);
  font-size: 11px;
}

@media (max-width: 900px) {
  .stat-strip {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 640px) {
  .stat-strip {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input {
    min-width: 100%;
  }

  .filter-group {
    width: 100%;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {
  .klaim-page {
    animation: none;
  }
}
</style>
