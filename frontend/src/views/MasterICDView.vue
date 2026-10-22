<template>
  <DefaultLayout>
    <div class="master-icd-page">
      <header class="page-header">
        <div>
          <span class="page-kicker">REFERENSI KLINIS</span>
          <h1>Master Data ICD</h1>
          <p>Kelola kode diagnosis (ICD-10), tindakan (ICD-9 CM), dan tarif INA-CBGs.</p>
        </div>
      </header>

      <div class="tabs-container">
        <nav class="tabs" role="tablist" aria-label="Jenis master data">
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

      <!-- ICD-10 Tab -->
      <section v-if="activeTab === 'icd10'" class="tab-content">
        <div class="section-bar">
          <span class="count-label">{{ icd10Total }} kode terdaftar</span>
          <button v-if="canMutate" @click="openICD10Modal(null)" class="primary-button" type="button">
            <i class="pi pi-plus" aria-hidden="true"></i>
            Tambah Kode
          </button>
        </div>

        <ModernTable
          :value="icd10List"
          :loading="icd10Loading"
          paginator
          :rows="icd10Rows"
          :first="(icd10Page - 1) * icd10Rows"
          :totalRecords="icd10Total"
          :rowsPerPageOptions="[10, 20, 50]"
          lazy
          dataKey="kode_icd10"
          ariaLabel="Daftar ICD-10"
          @page="onICD10PageChange"
        >
          <template #header>
            <th scope="col">Kode ICD-10</th>
            <th scope="col">Nama Diagnosis</th>
            <th scope="col" v-if="canMutate">Aksi</th>
          </template>
          <template #row="{ item }">
            <td><span class="code-badge font-mono">{{ item.kode_icd10 }}</span></td>
            <td>{{ item.nama_diagnosis }}</td>
            <td v-if="canMutate">
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openICD10Modal(item)" aria-label="Edit diagnosis" />
                <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDeleteICD10(item)" aria-label="Hapus diagnosis" />
              </div>
            </td>
          </template>
        </ModernTable>
      </section>

      <!-- ICD-9 CM Tab -->
      <section v-if="activeTab === 'icd9'" class="tab-content">
        <div class="section-bar">
          <span class="count-label">{{ icd9Total }} kode terdaftar</span>
          <button v-if="canMutate" @click="openICD9Modal(null)" class="primary-button" type="button">
            <i class="pi pi-plus" aria-hidden="true"></i>
            Tambah Kode
          </button>
        </div>

        <ModernTable
          :value="icd9List"
          :loading="icd9Loading"
          paginator
          :rows="icd9Rows"
          :first="(icd9Page - 1) * icd9Rows"
          :totalRecords="icd9Total"
          :rowsPerPageOptions="[10, 20, 50]"
          lazy
          dataKey="kode_tindakan"
          ariaLabel="Daftar ICD-9 CM"
          @page="onICD9PageChange"
        >
          <template #header>
            <th scope="col">Kode ICD-9 CM</th>
            <th scope="col">Nama Tindakan</th>
            <th scope="col">Tarif Standar</th>
            <th scope="col" v-if="canMutate">Aksi</th>
          </template>
          <template #row="{ item }">
            <td><span class="code-badge font-mono">{{ item.kode_tindakan }}</span></td>
            <td>{{ item.nama_tindakan }}</td>
            <td class="font-mono amount-cell">{{ formatRupiah(item.tarif_standar) }}</td>
            <td v-if="canMutate">
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openICD9Modal(item)" aria-label="Edit tindakan" />
                <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDeleteICD9(item)" aria-label="Hapus tindakan" />
              </div>
            </td>
          </template>
        </ModernTable>
      </section>

      <!-- INA-CBGs Tab -->
      <section v-if="activeTab === 'cbgs'" class="tab-content">
        <div class="section-bar">
          <span class="count-label">{{ cbgsTotal }} kode terdaftar</span>
          <button v-if="canMutate" @click="openCBGSModal(null)" class="primary-button" type="button">
            <i class="pi pi-plus" aria-hidden="true"></i>
            Tambah Kode
          </button>
        </div>

        <ModernTable
          :value="cbgsList"
          :loading="cbgsLoading"
          paginator
          :rows="cbgsRows"
          :first="(cbgsPage - 1) * cbgsRows"
          :totalRecords="cbgsTotal"
          :rowsPerPageOptions="[10, 20, 50]"
          lazy
          dataKey="kode_cbgs"
          ariaLabel="Daftar INA-CBGs"
          @page="onCBGSPageChange"
        >
          <template #header>
            <th scope="col">Kode INA-CBGs</th>
            <th scope="col">Deskripsi</th>
            <th scope="col">Tarif</th>
            <th scope="col" v-if="canMutate">Aksi</th>
          </template>
          <template #row="{ item }">
            <td><span class="code-badge font-mono">{{ item.kode_cbgs }}</span></td>
            <td>{{ item.deskripsi }}</td>
            <td class="font-mono amount-cell">{{ formatRupiah(item.tarif) }}</td>
            <td v-if="canMutate">
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openCBGSModal(item)" aria-label="Edit CBGs" />
                <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDeleteCBGS(item)" aria-label="Hapus CBGs" />
              </div>
            </td>
          </template>
        </ModernTable>
      </section>

      <!-- ICD-10 Form Modal -->
      <Dialog v-model:visible="showICD10Modal" modal :header="editingICD10 ? 'Edit Diagnosis' : 'Tambah Diagnosis'" :style="{ width: '32rem' }">
        <div class="form-grid">
          <div class="form-field">
            <label for="icd10-kode">Kode ICD-10 <span class="required">*</span></label>
            <InputText id="icd10-kode" v-model="icd10Form.kode_icd10" :disabled="!!editingICD10" class="w-full" />
          </div>
          <div class="form-field">
            <label for="icd10-nama">Nama Diagnosis <span class="required">*</span></label>
            <InputText id="icd10-nama" v-model="icd10Form.nama_diagnosis" class="w-full" />
          </div>
          <div class="form-field">
            <label for="icd10-kategori">Kategori</label>
            <InputText id="icd10-kategori" v-model="icd10Form.kategori" class="w-full" />
          </div>
        </div>
        <template #footer>
          <Button label="Batal" severity="secondary" @click="showICD10Modal = false" />
          <Button label="Simpan" @click="saveICD10" :loading="saving" />
        </template>
      </Dialog>

      <!-- ICD-9 CM Form Modal -->
      <Dialog v-model:visible="showICD9Modal" modal :header="editingICD9 ? 'Edit Tindakan' : 'Tambah Tindakan'" :style="{ width: '32rem' }">
        <div class="form-grid">
          <div class="form-field">
            <label for="icd9-kode">Kode ICD-9 CM <span class="required">*</span></label>
            <InputText id="icd9-kode" v-model="icd9Form.kode_tindakan" :disabled="!!editingICD9" class="w-full" />
          </div>
          <div class="form-field">
            <label for="icd9-nama">Nama Tindakan <span class="required">*</span></label>
            <InputText id="icd9-nama" v-model="icd9Form.nama_tindakan" class="w-full" />
          </div>
          <div class="form-field">
            <label for="icd9-tarif">Tarif Standar <span class="required">*</span></label>
            <InputNumber id="icd9-tarif" v-model="icd9Form.tarif_standar" :min="0" mode="currency" currency="IDR" locale="id-ID" class="w-full" />
            <small v-if="icd9FormError" class="field-error">{{ icd9FormError }}</small>
          </div>
        </div>
        <template #footer>
          <Button label="Batal" severity="secondary" @click="showICD9Modal = false" />
          <Button label="Simpan" @click="saveICD9" :loading="saving" />
        </template>
      </Dialog>

      <!-- INA-CBGs Form Modal -->
      <Dialog v-model:visible="showCBGSModal" modal :header="editingCBGS ? 'Edit Kode CBGs' : 'Tambah Kode CBGs'" :style="{ width: '32rem' }">
        <div class="form-grid">
          <div class="form-field">
            <label for="cbgs-kode">Kode INA-CBGs <span class="required">*</span></label>
            <InputText id="cbgs-kode" v-model="cbgsForm.kode_cbgs" :disabled="!!editingCBGS" class="w-full" />
          </div>
          <div class="form-field">
            <label for="cbgs-deskripsi">Deskripsi <span class="required">*</span></label>
            <Textarea id="cbgs-deskripsi" v-model="cbgsForm.deskripsi" rows="3" class="w-full" />
          </div>
          <div class="form-field">
            <label for="cbgs-tarif">Tarif <span class="required">*</span></label>
            <InputNumber id="cbgs-tarif" v-model="cbgsForm.tarif" :min="0" mode="currency" currency="IDR" locale="id-ID" class="w-full" />
            <small v-if="cbgsFormError" class="field-error">{{ cbgsFormError }}</small>
          </div>
        </div>
        <template #footer>
          <Button label="Batal" severity="secondary" @click="showCBGSModal = false" />
          <Button label="Simpan" @click="saveCBGS" :loading="saving" />
        </template>
      </Dialog>

      <DeleteConfirm
        :show="showDeleteConfirm"
        @update:show="showDeleteConfirm = $event"
        :title="deleteConfirmTitle"
        :message="deleteConfirmMessage"
        @confirm="executeDelete"
      />

      <Toast />
    </div>
  </DefaultLayout>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '../stores/auth'
import client from '../api/client'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import ModernTable from '../components/ModernTable.vue'
import DeleteConfirm from '../components/DeleteConfirm.vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Toast from 'primevue/toast'
import { formatRupiah } from '../utils/format.js'

const authStore = useAuthStore()
const toast = useToast()

const tabs = [
  { key: 'icd10', label: 'ICD-10 (Diagnosis)' },
  { key: 'icd9', label: 'ICD-9 CM (Tindakan)' },
  { key: 'cbgs', label: 'INA-CBGs' }
]

const activeTab = ref('icd10')
const canMutate = computed(() => authStore.user?.role === 'admin_ti')

// ICD-10
const icd10List = ref([])
const icd10Total = ref(0)
const icd10Loading = ref(false)
const icd10Page = ref(1)
const icd10Rows = ref(10)
const showICD10Modal = ref(false)
const editingICD10 = ref(null)
const icd10Form = ref({ kode_icd10: '', nama_diagnosis: '', kategori: '' })

// ICD-9 CM
const icd9List = ref([])
const icd9Total = ref(0)
const icd9Loading = ref(false)
const icd9Page = ref(1)
const icd9Rows = ref(10)
const showICD9Modal = ref(false)
const editingICD9 = ref(null)
const icd9Form = ref({ kode_tindakan: '', nama_tindakan: '', tarif_standar: 0 })
const icd9FormError = ref('')

// INA-CBGs
const cbgsList = ref([])
const cbgsTotal = ref(0)
const cbgsLoading = ref(false)
const cbgsPage = ref(1)
const cbgsRows = ref(10)
const showCBGSModal = ref(false)
const editingCBGS = ref(null)
const cbgsForm = ref({ kode_cbgs: '', deskripsi: '', tarif: 0 })
const cbgsFormError = ref('')

// Delete
const showDeleteConfirm = ref(false)
const deleteConfirmTitle = ref('')
const deleteConfirmMessage = ref('')
const deleteTarget = ref(null)
const deleteType = ref('')
const saving = ref(false)

watch(activeTab, async (newTab) => {
  if (newTab === 'icd10' && icd10List.value.length === 0) await fetchICD10()
  if (newTab === 'icd9' && icd9List.value.length === 0) await fetchICD9()
  if (newTab === 'cbgs' && cbgsList.value.length === 0) await fetchCBGS()
})

// Fetch ICD-10
async function fetchICD10() {
  icd10Loading.value = true
  try {
    const res = await client.get('/master/icd10', { params: { page: icd10Page.value, limit: icd10Rows.value } })
    if (res.data.success) {
      icd10List.value = res.data.data.items || []
      icd10Total.value = res.data.data.total || 0
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data ICD-10', life: 3000 })
  } finally {
    icd10Loading.value = false
  }
}

function onICD10PageChange(event) {
  icd10Page.value = event.page + 1
  icd10Rows.value = event.rows
  fetchICD10()
}

function openICD10Modal(item) {
  editingICD10.value = item
  if (item) {
    icd10Form.value = { kode_icd10: item.kode_icd10, nama_diagnosis: item.nama_diagnosis, kategori: item.kategori || '' }
  } else {
    icd10Form.value = { kode_icd10: '', nama_diagnosis: '', kategori: '' }
  }
  showICD10Modal.value = true
}

async function saveICD10() {
  if (!icd10Form.value.kode_icd10 || !icd10Form.value.nama_diagnosis) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Kode dan nama diagnosis wajib diisi', life: 3000 })
    return
  }
  saving.value = true
  try {
    const payload = { ...icd10Form.value }
    const res = editingICD10.value
      ? await client.put(`/master/icd10/${editingICD10.value.kode_icd10}`, payload)
      : await client.post('/master/icd10', payload)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editingICD10.value ? 'Diagnosis diperbarui' : 'Diagnosis ditambahkan', life: 3000 })
      showICD10Modal.value = false
      await fetchICD10()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menyimpan', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menyimpan', life: 3000 })
  } finally {
    saving.value = false
  }
}

function confirmDeleteICD10(item) {
  deleteTarget.value = item
  deleteType.value = 'icd10'
  deleteConfirmTitle.value = 'Hapus Diagnosis'
  deleteConfirmMessage.value = `Apakah Anda yakin ingin menghapus diagnosis "${item.kode_icd10} - ${item.nama_diagnosis}"?`
  showDeleteConfirm.value = true
}

// Fetch ICD-9 CM
async function fetchICD9() {
  icd9Loading.value = true
  try {
    const res = await client.get('/master/icd9', { params: { page: icd9Page.value, limit: icd9Rows.value } })
    if (res.data.success) {
      icd9List.value = res.data.data.items || []
      icd9Total.value = res.data.data.total || 0
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data ICD-9 CM', life: 3000 })
  } finally {
    icd9Loading.value = false
  }
}

function onICD9PageChange(event) {
  icd9Page.value = event.page + 1
  icd9Rows.value = event.rows
  fetchICD9()
}

function openICD9Modal(item) {
  editingICD9.value = item
  icd9FormError.value = ''
  if (item) {
    icd9Form.value = { kode_tindakan: item.kode_tindakan, nama_tindakan: item.nama_tindakan, tarif_standar: item.tarif_standar }
  } else {
    icd9Form.value = { kode_tindakan: '', nama_tindakan: '', tarif_standar: 0 }
  }
  showICD9Modal.value = true
}

async function saveICD9() {
  if (!icd9Form.value.kode_tindakan || !icd9Form.value.nama_tindakan) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Kode dan nama tindakan wajib diisi', life: 3000 })
    return
  }
  if (icd9Form.value.tarif_standar < 0) {
    icd9FormError.value = 'Tarif tidak boleh negatif'
    return
  }
  saving.value = true
  try {
    const payload = { ...icd9Form.value }
    const res = editingICD9.value
      ? await client.put(`/master/icd9/${editingICD9.value.kode_tindakan}`, payload)
      : await client.post('/master/icd9', payload)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editingICD9.value ? 'Tindakan diperbarui' : 'Tindakan ditambahkan', life: 3000 })
      showICD9Modal.value = false
      await fetchICD9()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menyimpan', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menyimpan', life: 3000 })
  } finally {
    saving.value = false
  }
}

function confirmDeleteICD9(item) {
  deleteTarget.value = item
  deleteType.value = 'icd9'
  deleteConfirmTitle.value = 'Hapus Tindakan'
  deleteConfirmMessage.value = `Apakah Anda yakin ingin menghapus tindakan "${item.kode_tindakan} - ${item.nama_tindakan}"?`
  showDeleteConfirm.value = true
}

// Fetch INA-CBGs
async function fetchCBGS() {
  cbgsLoading.value = true
  try {
    const res = await client.get('/master/cbgs', { params: { page: cbgsPage.value, limit: cbgsRows.value } })
    if (res.data.success) {
      cbgsList.value = res.data.data.items || []
      cbgsTotal.value = res.data.data.total || 0
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data INA-CBGs', life: 3000 })
  } finally {
    cbgsLoading.value = false
  }
}

function onCBGSPageChange(event) {
  cbgsPage.value = event.page + 1
  cbgsRows.value = event.rows
  fetchCBGS()
}

function openCBGSModal(item) {
  editingCBGS.value = item
  cbgsFormError.value = ''
  if (item) {
    cbgsForm.value = { kode_cbgs: item.kode_cbgs, deskripsi: item.deskripsi, tarif: item.tarif }
  } else {
    cbgsForm.value = { kode_cbgs: '', deskripsi: '', tarif: 0 }
  }
  showCBGSModal.value = true
}

async function saveCBGS() {
  if (!cbgsForm.value.kode_cbgs || !cbgsForm.value.deskripsi) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Kode dan deskripsi wajib diisi', life: 3000 })
    return
  }
  if (cbgsForm.value.tarif < 0) {
    cbgsFormError.value = 'Tarif tidak boleh negatif'
    return
  }
  saving.value = true
  try {
    const payload = { ...cbgsForm.value }
    const res = editingCBGS.value
      ? await client.put(`/master/cbgs/${editingCBGS.value.kode_cbgs}`, payload)
      : await client.post('/master/cbgs', payload)
    if (res.data.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: editingCBGS.value ? 'Kode CBGs diperbarui' : 'Kode CBGs ditambahkan', life: 3000 })
      showCBGSModal.value = false
      await fetchCBGS()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menyimpan', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menyimpan', life: 3000 })
  } finally {
    saving.value = false
  }
}

function confirmDeleteCBGS(item) {
  deleteTarget.value = item
  deleteType.value = 'cbgs'
  deleteConfirmTitle.value = 'Hapus Kode CBGs'
  deleteConfirmMessage.value = `Apakah Anda yakin ingin menghapus kode "${item.kode_cbgs} - ${item.deskripsi}"?`
  showDeleteConfirm.value = true
}

async function executeDelete() {
  try {
    let res
    if (deleteType.value === 'icd10') {
      res = await client.delete(`/master/icd10/${deleteTarget.value.kode_icd10}`)
      if (res.data.success) {
        toast.add({ severity: 'success', summary: 'Sukses', detail: 'Diagnosis dihapus', life: 3000 })
        await fetchICD10()
      }
    } else if (deleteType.value === 'icd9') {
      res = await client.delete(`/master/icd9/${deleteTarget.value.kode_tindakan}`)
      if (res.data.success) {
        toast.add({ severity: 'success', summary: 'Sukses', detail: 'Tindakan dihapus', life: 3000 })
        await fetchICD9()
      }
    } else if (deleteType.value === 'cbgs') {
      res = await client.delete(`/master/cbgs/${deleteTarget.value.kode_cbgs}`)
      if (res.data.success) {
        toast.add({ severity: 'success', summary: 'Sukses', detail: 'Kode CBGs dihapus', life: 3000 })
        await fetchCBGS()
      }
    }
    if (res && !res.data.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: res.data.message || 'Gagal menghapus', life: 3000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.message || 'Gagal menghapus', life: 3000 })
  }
  showDeleteConfirm.value = false
}

// Initial load
fetchICD10()
</script>

<style scoped>
.master-icd-page {
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

.tabs-container {
  overflow-x: auto;
  padding: 5px;
  border: 1px solid var(--color-border);
  border-radius: 11px;
  background: #eeede9;
}

.tabs {
  display: flex;
  min-width: max-content;
  gap: 4px;
}

.tabs button {
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

.tabs button:hover:not(.is-active) {
  background: rgb(255 255 255 / 55%);
  color: var(--color-text-primary);
}

.tabs button.is-active {
  background: white;
  box-shadow: 0 1px 3px rgb(31 41 55 / 10%);
  color: var(--color-primary);
  font-weight: 600;
}

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
  animation: tab-fade-in 200ms ease-out both;
}

.section-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.count-label {
  color: var(--color-text-secondary);
  font-size: 13px;
  font-weight: 500;
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
  gap: 16px;
  margin-top: 10px;
}

.form-field {
  display: grid;
  gap: 6px;
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

@keyframes tab-fade-in {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .master-icd-page,
  .tab-content {
    animation: none;
  }
}
</style>
