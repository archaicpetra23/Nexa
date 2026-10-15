<template>
  <DefaultLayout>
    <div class="rekam-medis-page">
      <header class="page-header">
        <div>
          <span class="page-kicker">PELAYANAN</span>
          <h1>Rekam Medis</h1>
          <p>Kelola catatan kunjungan pasien, diagnosis, dan tindakan medis.</p>
        </div>
      </header>

      <div class="filter-bar">
        <label class="search-input">
          <i class="pi pi-search" aria-hidden="true"></i>
          <InputText v-model="search" placeholder="Cari nomor RM atau pasien..." @input="debouncedSearch" />
        </label>
        <button @click="openModal(null)" class="primary-button" type="button">
          <i class="pi pi-plus" aria-hidden="true"></i>
          Tambah Kunjungan
        </button>
      </div>

      <ModernTable
        :value="rekamList"
        :loading="loading"
        paginator
        lazy
        :rows="limit"
        :first="(page - 1) * limit"
        :totalRecords="total"
        :rowsPerPageOptions="[10, 20, 50]"
        dataKey="id_rekam"
        ariaLabel="Daftar rekam medis"
        @page="onPageChange"
      >
        <template #header>
          <th scope="col">No. RM</th>
          <th scope="col">ID Pasien</th>
          <th scope="col">Tgl. Kunjungan</th>
          <th scope="col">Jenis Perawatan</th>
          <th scope="col">Keluhan</th>
          <th scope="col">Dokter</th>
          <th scope="col">Aksi</th>
        </template>
        <template #row="{ item }">
          <td class="font-mono">{{ item.id_rekam }}</td>
          <td class="font-mono">{{ item.id_pasien }}</td>
          <td>{{ formatDate(item.tanggal_kunjungan) }}</td>
          <td>{{ item.jenis_perawatan === 'rawat_jalan' ? 'Rawat Jalan' : 'Rawat Inap' }}</td>
          <td class="complaint-cell">{{ item.keluhan }}</td>
          <td>{{ item.id_dokter }}</td>
          <td>
            <div class="flex gap-1">
              <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openModal(item)" aria-label="Edit rekam medis" />
              <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDelete(item)" aria-label="Hapus rekam medis" />
            </div>
          </td>
        </template>
      </ModernTable>

      <!-- Form Modal with Tabs -->
      <Dialog v-model:visible="showModal" modal :header="editing ? 'Edit Kunjungan' : 'Tambah Kunjungan'" :style="{ width: '40rem' }">
        <div v-if="showModal" class="tabs-container">
          <nav class="form-tabs" role="tablist">
            <button
              v-for="tab in formTabs"
              :key="tab.key"
              @click="activeFormTab = tab.key"
              :class="{ 'is-active': activeFormTab === tab.key }"
              role="tab"
              :aria-selected="activeFormTab === tab.key"
              type="button"
            >
              {{ tab.label }}
            </button>
          </nav>

          <!-- Tab: Kunjungan -->
          <div v-if="activeFormTab === 'kunjungan'" class="form-grid">
            <div class="form-field">
              <label for="id_pasien">ID Pasien <span class="required">*</span></label>
              <InputNumber id="id_pasien" v-model="form.id_pasien" :min="1" :disabled="!!editing" class="w-full" />
              <small v-if="formErrors.id_pasien" class="field-error">{{ formErrors.id_pasien }}</small>
            </div>
            <div class="form-field">
              <label for="id_dokter">ID Dokter <span class="required">*</span></label>
              <InputNumber id="id_dokter" v-model="form.id_dokter" :min="1" class="w-full" />
              <small v-if="formErrors.id_dokter" class="field-error">{{ formErrors.id_dokter }}</small>
            </div>
            <div class="form-field">
              <label for="jenis_perawatan">Jenis Perawatan <span class="required">*</span></label>
              <Select id="jenis_perawatan" v-model="form.jenis_perawatan" :options="jenisPerawatanOptions" optionLabel="label" optionValue="value" placeholder="Pilih jenis perawatan" class="w-full" />
              <small v-if="formErrors.jenis_perawatan" class="field-error">{{ formErrors.jenis_perawatan }}</small>
            </div>
            <div class="form-field">
              <label for="tanggal_kunjungan">Tanggal Kunjungan <span class="required">*</span></label>
              <DatePicker id="tanggal_kunjungan" v-model="form.tanggal_kunjungan" dateFormat="dd/mm/yy" showIcon class="w-full" />
              <small v-if="formErrors.tanggal_kunjungan" class="field-error">{{ formErrors.tanggal_kunjungan }}</small>
            </div>
            <div class="form-field">
              <label for="tanggal_pulang">Tanggal Pulang</label>
              <DatePicker id="tanggal_pulang" v-model="form.tanggal_pulang" dateFormat="dd/mm/yy" showIcon class="w-full" />
            </div>
            <div class="form-field full-width">
              <label for="keluhan">Keluhan <span class="required">*</span></label>
              <Textarea id="keluhan" v-model="form.keluhan" rows="3" class="w-full" />
              <small v-if="formErrors.keluhan" class="field-error">{{ formErrors.keluhan }}</small>
            </div>
            <div class="form-field full-width">
              <label for="catatan">Catatan</label>
              <Textarea id="catatan" v-model="form.catatan" rows="3" class="w-full" />
            </div>
          </div>

          <!-- Tab: Diagnosis -->
          <div v-if="activeFormTab === 'diagnosis'" class="diagnosis-tab">
            <div v-if="!editing" class="empty-state">
              <i class="pi pi-info-circle" aria-hidden="true"></i>
              <p>Simpan kunjungan dulu untuk menambah diagnosis</p>
            </div>
            <div v-else class="form-grid">
              <div class="form-field full-width">
                <label for="icd10-search">ICD-10 <span class="required">*</span></label>
                <AutoComplete
                  id="icd10-search"
                  v-model="diagnosisForm.kode_icd10"
                  :suggestions="icd10Suggestions"
                  :loading="icd10Loading"
                  placeholder="Cari kode atau diagnosis..."
                  @complete="searchICD10"
                  field="label"
                  class="w-full"
                >
                  <template #item="{ item }">
                    <div class="icd-item">
                      <span class="code font-mono">{{ item.kode }}</span>
                      <span class="name">{{ item.nama }}</span>
                    </div>
                  </template>
                </AutoComplete>
              </div>
              <div class="form-field">
                <label for="jenis">Jenis <span class="required">*</span></label>
                <Select id="jenis" v-model="diagnosisForm.jenis" :options="jenisOptions" optionLabel="label" optionValue="value" placeholder="Pilih jenis" class="w-full" />
              </div>
              <div class="form-field full-width">
                <button @click="addDiagnosis" class="secondary-button" type="button">
                  <i class="pi pi-plus" aria-hidden="true"></i>
                  Tambah Diagnosis
                </button>
              </div>
              <div v-if="form.diagnoses.length" class="form-field full-width">
                <span class="label">Diagnosis yang ditambahkan</span>
                <div class="diagnosis-list">
                  <div v-for="(diag, idx) in form.diagnoses" :key="idx" class="diagnosis-item">
                    <span class="font-mono">{{ diag.kode_icd10 }} ({{ diag.jenis }})</span>
                    <Button icon="pi pi-trash" severity="danger" text rounded @click="removeDiagnosis(idx)" aria-label="Hapus diagnosis" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Tab: Tindakan -->
          <div v-if="activeFormTab === 'tindakan'" class="tindakan-tab">
            <div v-if="!editing" class="empty-state">
              <i class="pi pi-info-circle" aria-hidden="true"></i>
              <p>Simpan kunjungan dulu untuk menambah tindakan</p>
            </div>
            <div v-else class="form-grid">
              <div class="form-field full-width">
                <label for="icd9-search">ICD-9 CM <span class="required">*</span></label>
                <AutoComplete
                  id="icd9-search"
                  v-model="tindakanForm.kode_tindakan"
                  :suggestions="icd9Suggestions"
                  :loading="icd9Loading"
                  placeholder="Cari kode atau tindakan..."
                  @complete="searchICD9"
                  field="label"
                  class="w-full"
                >
                  <template #item="{ item }">
                    <div class="icd-item">
                      <span class="code font-mono">{{ item.kode }}</span>
                      <span class="name">{{ item.nama }}</span>
                      <span class="tarif">{{ formatRupiah(item.tarif) }}</span>
                    </div>
                  </template>
                </AutoComplete>
              </div>
              <div class="form-field">
                <label for="jumlah">Jumlah <span class="required">*</span></label>
                <InputNumber id="jumlah" v-model="tindakanForm.jumlah" :min="1" class="w-full" />
              </div>
              <div class="form-field full-width">
                <button @click="addTindakan" class="secondary-button" type="button">
                  <i class="pi pi-plus" aria-hidden="true"></i>
                  Tambah Tindakan
                </button>
              </div>
              <div v-if="form.tindakans.length" class="form-field full-width">
                <span class="label">Tindakan yang ditambahkan</span>
                <div class="tindakan-list">
                  <div v-for="(tind, idx) in form.tindakans" :key="idx" class="tindakan-item">
                    <span class="font-mono">{{ tind.kode_tindakan }} x{{ tind.jumlah }}</span>
                    <Button icon="pi pi-trash" severity="danger" text rounded @click="removeTindakan(idx)" aria-label="Hapus tindakan" />
                  </div>
                </div>
              </div>
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
        :message="`Apakah Anda yakin ingin menghapus rekam medis ID: ${deleteTarget?.id_rekam}?`"
        @confirm="executeDelete"
      />

      <Toast />
    </div>
  </DefaultLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import client from '../api/client'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import ModernTable from '../components/ModernTable.vue'
import DeleteConfirm from '../components/DeleteConfirm.vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import DatePicker from 'primevue/datepicker'
import Select from 'primevue/select'
import AutoComplete from 'primevue/autocomplete'
import Toast from 'primevue/toast'

const toast = useToast()

const rekamList = ref([])
const loading = ref(false)
const page = ref(1)
const limit = ref(10)
const total = ref(0)
const search = ref('')
let searchTimeout = null

const formTabs = [
  { key: 'kunjungan', label: 'Kunjungan' },
  { key: 'diagnosis', label: 'Diagnosis' },
  { key: 'tindakan', label: 'Tindakan' }
]
const activeFormTab = ref('kunjungan')

const showModal = ref(false)
const editing = ref(null)
const saving = ref(false)

const form = ref({
  id_pasien: null,
  id_dokter: null,
  jenis_perawatan: '',
  tanggal_kunjungan: null,
  tanggal_pulang: null,
  keluhan: '',
  catatan: '',
  diagnoses: [],
  tindakans: []
})
const formErrors = ref({})

const jenisPerawatanOptions = [
  { label: 'Rawat Jalan', value: 'rawat_jalan' },
  { label: 'Rawat Inap', value: 'rawat_inap' }
]

const diagnosisForm = ref({ kode_icd10: '', jenis: '' })
const jenisOptions = [
  { label: 'Primer', value: 'primer' },
  { label: 'Sekunder', value: 'sekunder' }
]
const icd10Suggestions = ref([])
const icd10Loading = ref(false)
let icd10SearchTimeout = null

const tindakanForm = ref({ kode_tindakan: '', jumlah: 1 })
const icd9Suggestions = ref([])
const icd9Loading = ref(false)
let icd9SearchTimeout = null

const showDeleteConfirm = ref(false)
const deleteTarget = ref(null)

async function fetchRekam() {
  loading.value = true
  try {
    const params = { page: page.value, limit: limit.value }
    if (search.value) params.search = search.value
    const res = await client.get('/rekam-medis', { params })
    if (res.data.success) {
      rekamList.value = res.data.data.items || []
      total.value = res.data.data.total || 0
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data rekam medis', life: 3000 })
  } finally {
    loading.value = false
  }
}

function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchRekam()
  }, 300)
}

function onPageChange(event) {
  page.value = event.page + 1
  limit.value = event.rows
  fetchRekam()
}

function openModal(item) {
  editing.value = item
  activeFormTab.value = 'kunjungan'
  formErrors.value = {}
  if (item) {
    form.value = {
      id_pasien: item.id_pasien,
      id_dokter: item.id_dokter,
      jenis_perawatan: item.jenis_perawatan,
      tanggal_kunjungan: new Date(item.tanggal_kunjungan),
      tanggal_pulang: item.tanggal_pulang ? new Date(item.tanggal_pulang) : null,
      keluhan: item.keluhan,
      catatan: item.catatan || '',
      diagnoses: [],
      tindakans: []
    }
  } else {
    form.value = {
      id_pasien: null,
      id_dokter: null,
      jenis_perawatan: '',
      tanggal_kunjungan: null,
      tanggal_pulang: null,
      keluhan: '',
      catatan: '',
      diagnoses: [],
      tindakans: []
    }
  }
  showModal.value = true
}

function searchICD10(event) {
  clearTimeout(icd10SearchTimeout)
  icd10SearchTimeout = setTimeout(async () => {
    if (!event.query || event.query.length < 1) {
      icd10Suggestions.value = []
      return
    }
    icd10Loading.value = true
    try {
      const res = await client.get('/master/icd10/search', { params: { q: event.query, limit: 10 } })
      if (res.data.success) {
        icd10Suggestions.value = (res.data.data || []).map(d => ({
          label: `${d.kode_icd10} - ${d.nama_diagnosis}`,
          kode: d.kode_icd10,
          nama: d.nama_diagnosis
        }))
      }
    } catch (e) {
      console.error('ICD-10 search error:', e)
    } finally {
      icd10Loading.value = false
    }
  }, 300)
}

function addDiagnosis() {
  if (!diagnosisForm.value.kode_icd10 || !diagnosisForm.value.jenis) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Kode ICD dan jenis wajib diisi', life: 3000 })
    return
  }
  form.value.diagnoses.push({ ...diagnosisForm.value })
  diagnosisForm.value = { kode_icd10: '', jenis: '' }
  icd10Suggestions.value = []
}

function removeDiagnosis(index) {
  form.value.diagnoses.splice(index, 1)
}

function searchICD9(event) {
  clearTimeout(icd9SearchTimeout)
  icd9SearchTimeout = setTimeout(async () => {
    if (!event.query || event.query.length < 1) {
      icd9Suggestions.value = []
      return
    }
    icd9Loading.value = true
    try {
      const res = await client.get('/master/icd9/search', { params: { q: event.query, limit: 10 } })
      if (res.data.success) {
        icd9Suggestions.value = (res.data.data || []).map(t => ({
          label: `${t.kode_tindakan} - ${t.nama_tindakan}`,
          kode: t.kode_tindakan,
          nama: t.nama_tindakan,
          tarif: t.tarif_standar
        }))
      }
    } catch (e) {
      console.error('ICD-9 search error:', e)
    } finally {
      icd9Loading.value = false
    }
  }, 300)
}

function addTindakan() {
  if (!tindakanForm.value.kode_tindakan || !tindakanForm.value.jumlah) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Kode tindakan dan jumlah wajib diisi', life: 3000 })
    return
  }
  form.value.tindakans.push({ ...tindakanForm.value })
  tindakanForm.value = { kode_tindakan: '', jumlah: 1 }
  icd9Suggestions.value = []
}

function removeTindakan(index) {
  form.value.tindakans.splice(index, 1)
}

function validate() {
  formErrors.value = {}
  if (!form.value.id_pasien) formErrors.value.id_pasien = 'ID Pasien wajib diisi'
  if (!form.value.id_dokter) formErrors.value.id_dokter = 'ID Dokter wajib diisi'
  if (!form.value.jenis_perawatan) formErrors.value.jenis_perawatan = 'Jenis Perawatan wajib dipilih'
  if (!form.value.tanggal_kunjungan) formErrors.value.tanggal_kunjungan = 'Tanggal Kunjungan wajib diisi'
  if (!form.value.keluhan) formErrors.value.keluhan = 'Keluhan wajib diisi'
  return Object.keys(formErrors.value).length === 0
}

async function save() {
  if (!validate()) return
  saving.value = true
  try {
    const payload = {
      id_pasien: form.value.id_pasien,
      id_dokter: form.value.id_dokter,
      jenis_perawatan: form.value.jenis_perawatan,
      tanggal_kunjungan: formatDateForAPI(form.value.tanggal_kunjungan),
      tanggal_pulang: form.value.tanggal_pulang ? formatDateForAPI(form.value.tanggal_pulang) : null,
      keluhan: form.value.keluhan,
      catatan: form.value.catatan || ''
    }
    const res = editing.value
      ? await client.put(`/rekam-medis/${editing.value.id_rekam}`, payload)
      : await client.post('/rekam-medis', payload)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editing.value ? 'Rekam medis diperbarui' : 'Rekam medis ditambahkan', life: 3000 })
      showModal.value = false
      await fetchRekam()
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
    const res = await client.delete(`/rekam-medis/${deleteTarget.value.id_rekam}`)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Rekam medis dihapus', life: 3000 })
      await fetchRekam()
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

function formatDateForAPI(date) {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatRupiah(value) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value)
}

fetchRekam()
</script>

<style scoped>
.rekam-medis-page {
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

.complaint-cell {
  max-width: 240px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tabs-container {
  margin-top: 10px;
}

.form-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--color-border);
}

.form-tabs button {
  padding: 6px 12px;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: background-color 140ms ease, color 140ms ease;
}

.form-tabs button:hover:not(.is-active) {
  background: #f0f3f5;
}

.form-tabs button.is-active {
  color: var(--color-primary);
  font-weight: 600;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 32px;
  color: var(--color-text-secondary);
  text-align: center;
}

.empty-state i {
  font-size: 24px;
  color: var(--color-text-disabled);
}

.diagnosis-tab,
.tindakan-tab {
  margin-bottom: 12px;
}

.icd-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
}

.icd-item .code {
  flex: 0 0 auto;
  font-weight: 600;
  color: var(--color-primary);
}

.icd-item .name {
  flex: 1;
  color: var(--color-text-primary);
}

.icd-item .tarif {
  flex: 0 0 auto;
  color: var(--color-text-secondary);
  font-size: 11px;
}

.secondary-button {
  display: inline-flex;
  min-height: 36px;
  align-items: center;
  gap: 6px;
  padding: 0 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: white;
  color: var(--color-primary);
  cursor: pointer;
  font: inherit;
  font-size: 12px;
  font-weight: 600;
  transition: background-color 140ms ease;
}

.secondary-button:hover {
  background: #f5f8fa;
}

.diagnosis-list,
.tindakan-list {
  display: grid;
  gap: 8px;
}

.diagnosis-item,
.tindakan-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid #dce6eb;
  border-radius: 6px;
  background: #f8fafc;
}

.label {
  display: block;
  color: var(--color-text-primary);
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 10px;
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
  .rekam-medis-page {
    animation: none;
  }
}
</style>
