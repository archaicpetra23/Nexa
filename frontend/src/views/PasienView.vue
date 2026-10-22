<template>
  <DefaultLayout>
    <div class="pasien-page">
      <header class="page-header">
        <div>
          <span class="page-kicker">DATA MASTER</span>
          <h1>Data Pasien</h1>
          <p>Kelola informasi pasien untuk keperluan rekam medis dan klaim BPJS.</p>
        </div>
      </header>

      <div class="filter-bar">
        <label class="search-input">
          <i class="pi pi-search" aria-hidden="true"></i>
          <InputText v-model="search" placeholder="Cari nama atau NIK..." @input="debouncedSearch" />
        </label>
        <button @click="openModal(null)" class="primary-button" type="button">
          <i class="pi pi-plus" aria-hidden="true"></i>
          Tambah Pasien
        </button>
      </div>

      <ModernTable
        :value="pasienList"
        :loading="loading"
        paginator
        lazy
        :rows="limit"
        :first="(page - 1) * limit"
        :totalRecords="total"
        :rowsPerPageOptions="[10, 20, 50]"
        dataKey="id_pasien"
        ariaLabel="Daftar pasien"
        @page="onPageChange"
      >
        <template #header>
          <th scope="col">NIK</th>
          <th scope="col">Nama Pasien</th>
          <th scope="col">Tgl. Lahir</th>
          <th scope="col">Jenis Kelamin</th>
          <th scope="col">No. BPJS</th>
          <th scope="col">Alamat</th>
          <th scope="col">Aksi</th>
        </template>
        <template #row="{ item }">
          <td class="font-mono">{{ item.nik }}</td>
          <td class="user-name-cell">{{ item.nama }}</td>
          <td>{{ formatTanggalSingkat(item.tanggal_lahir) }}</td>
          <td>{{ item.jenis_kelamin === 'L' ? 'Laki-laki' : 'Perempuan' }}</td>
          <td class="font-mono">{{ item.no_bpjs || '—' }}</td>
          <td class="address-cell">{{ item.alamat }}</td>
          <td>
            <div class="flex gap-1">
              <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openModal(item)" aria-label="Edit pasien" />
              <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDelete(item)" aria-label="Hapus pasien" />
            </div>
          </td>
        </template>
      </ModernTable>

      <Dialog v-model:visible="showModal" modal :header="editing ? 'Edit Pasien' : 'Tambah Pasien'" :style="{ width: '36rem' }">
        <div class="form-section">
          <h3 class="section-title">Identitas Pasien</h3>
          <div class="form-grid">
            <div class="form-field">
              <label for="nik">NIK <span class="required">*</span></label>
              <InputText id="nik" v-model="form.nik" maxlength="16" class="w-full" />
              <small v-if="errors.nik" class="field-error">{{ errors.nik }}</small>
            </div>
            <div class="form-field">
              <label for="nama">Nama <span class="required">*</span></label>
              <InputText id="nama" v-model="form.nama" class="w-full" />
              <small v-if="errors.nama" class="field-error">{{ errors.nama }}</small>
            </div>
            <div class="form-field">
              <label for="tanggal_lahir">Tanggal Lahir <span class="required">*</span></label>
              <DatePicker id="tanggal_lahir" v-model="form.tanggal_lahir" dateFormat="dd/mm/yy" showIcon class="w-full" />
              <small v-if="errors.tanggal_lahir" class="field-error">{{ errors.tanggal_lahir }}</small>
            </div>
            <div class="form-field">
              <label for="jenis_kelamin">Jenis Kelamin <span class="required">*</span></label>
              <Select id="jenis_kelamin" v-model="form.jenis_kelamin" :options="jenisKelaminOptions" optionLabel="label" optionValue="value" placeholder="Pilih jenis kelamin" class="w-full" />
              <small v-if="errors.jenis_kelamin" class="field-error">{{ errors.jenis_kelamin }}</small>
            </div>
            <div class="form-field">
              <label for="no_bpjs">No. BPJS</label>
              <InputText id="no_bpjs" v-model="form.no_bpjs" maxlength="13" class="w-full" />
              <small v-if="errors.no_bpjs" class="field-error">{{ errors.no_bpjs }}</small>
            </div>
            <div class="form-field">
              <label for="no_hp">No. HP</label>
              <InputText id="no_hp" v-model="form.no_hp" class="w-full" />
            </div>
            <div class="form-field full-width">
              <label for="alamat">Alamat <span class="required">*</span></label>
              <Textarea id="alamat" v-model="form.alamat" rows="3" class="w-full" />
              <small v-if="errors.alamat" class="field-error">{{ errors.alamat }}</small>
            </div>
          </div>
        </div>
        <template #footer>
          <Button label="Batal" severity="secondary" @click="showModal = false" />
          <Button label="Simpan" @click="save" :loading="saving" :disabled="saving" />
        </template>
      </Dialog>

      <DeleteConfirm
        :show="showDeleteConfirm"
        @update:show="showDeleteConfirm = $event"
        :message="deleteConfirmMessage"
        @confirm="executeDelete"
      />

      <Toast />
    </div>
  </DefaultLayout>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import client from '../api/client'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import ModernTable from '../components/ModernTable.vue'
import DeleteConfirm from '../components/DeleteConfirm.vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import DatePicker from 'primevue/datepicker'
import Select from 'primevue/select'
import { validatePasienForm } from '../utils/pasienValidation'
import Toast from 'primevue/toast'
import { formatTanggalSingkat } from '../utils/format.js'

const toast = useToast()

const pasienList = ref([])
const loading = ref(false)
const page = ref(1)
const limit = ref(10)
const total = ref(0)
const search = ref('')
let searchTimeout = null

const showModal = ref(false)
const editing = ref(null)
const saving = ref(false)
const form = ref({
  nik: '',
  nama: '',
  tanggal_lahir: null,
  jenis_kelamin: '',
  alamat: '',
  no_bpjs: '',
  no_hp: ''
})
const errors = ref({})

const jenisKelaminOptions = [
  { label: 'Laki-laki', value: 'L' },
  { label: 'Perempuan', value: 'P' }
]

const showDeleteConfirm = ref(false)
const deleteTarget = ref(null)
const deleteConfirmMessage = computed(() =>
  `Apakah Anda yakin ingin menghapus pasien "${deleteTarget.value?.nama}" (NIK: ${deleteTarget.value?.nik})?`,
)

async function fetchPasien() {
  loading.value = true
  try {
    const params = { page: page.value, limit: limit.value }
    if (search.value) params.search = search.value
    const res = await client.get('/pasien', { params })
    if (res.data.success) {
      pasienList.value = res.data.data.items || []
      total.value = res.data.data.total || 0
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data pasien', life: 3000 })
  } finally {
    loading.value = false
  }
}

function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchPasien()
  }, 300)
}

function onPageChange(event) {
  page.value = event.page + 1
  limit.value = event.rows
  fetchPasien()
}

function openModal(item) {
  editing.value = item
  errors.value = {}
  if (item) {
    form.value = {
      nik: item.nik,
      nama: item.nama,
      tanggal_lahir: new Date(item.tanggal_lahir),
      jenis_kelamin: item.jenis_kelamin,
      alamat: item.alamat,
      no_bpjs: item.no_bpjs || '',
      no_hp: item.no_hp || ''
    }
  } else {
    form.value = {
      nik: '',
      nama: '',
      tanggal_lahir: null,
      jenis_kelamin: '',
      alamat: '',
      no_bpjs: '',
      no_hp: ''
    }
  }
  showModal.value = true
}

function validate() {
  const result = validatePasienForm(form.value)
  errors.value = result.errors
  return result.valid
}

async function save() {
  if (!validate()) return
  saving.value = true
  try {
    const payload = {
      ...form.value,
      tanggal_lahir: formatDateForAPI(form.value.tanggal_lahir)
    }
    const res = editing.value
      ? await client.put(`/pasien/${editing.value.id_pasien}`, payload)
      : await client.post('/pasien', payload)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editing.value ? 'Pasien diperbarui' : 'Pasien ditambahkan', life: 3000 })
      showModal.value = false
      await fetchPasien()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menyimpan', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menyimpan', life: 3000 })
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
    const res = await client.delete(`/pasien/${deleteTarget.value.id_pasien}`)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Pasien dihapus', life: 3000 })
      await fetchPasien()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menghapus', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menghapus', life: 3000 })
  }
  showDeleteConfirm.value = false
}

function formatDateForAPI(date) {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

fetchPasien()
</script>

<style scoped>
.pasien-page {
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
  min-width: min(100%, 320px);
  min-height: 40px;
  flex: 1 1 320px;
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
  font-size: 13px;
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

.address-cell {
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.form-section {
  margin-top: 10px;
}

.section-title {
  margin: 0 0 16px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-primary);
  font-size: 14px;
  font-weight: 600;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
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

.required {
  color: #B54245;
}

.field-error {
  color: #B54245;
  font-size: 11px;
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input {
    min-width: 100%;
  }
}

@media (prefers-reduced-motion: reduce) {
  .pasien-page {
    animation: none;
  }
}
</style>
