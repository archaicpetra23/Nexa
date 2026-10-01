<template>
  <Dialog
    :visible="show"
    modal
    :header="title"
    :closable="!saving"
    :closeOnEscape="!saving"
    :dismissableMask="!saving"
    :style="{ width: 'min(48rem, calc(100vw - 2rem))' }"
    :pt="{
      root: { style: { backgroundColor: '#FFFFFF', border: '1px solid var(--color-border)', borderRadius: '12px' } },
      mask: { style: { backgroundColor: 'rgba(15, 23, 42, 0.48)' } },
      content: { style: { backgroundColor: '#FFFFFF' } }
    }"
    @update:visible="emit('update:show', $event)"
  >
    <template #header>
      <div>
        <h3 class="text-lg font-semibold" style="color: var(--color-primary)">{{ title }}</h3>
        <p class="mt-1 text-sm text-secondary">
          {{ isEdit ? 'Perbarui informasi akun dan penempatan pengguna.' : 'Lengkapi informasi untuk membuat akun pengguna baru.' }}
        </p>
      </div>
    </template>

    <form class="space-y-6 px-1 pb-1 pt-2 sm:px-2" @submit.prevent="submit">
      <div>
        <h4 class="mb-3 text-sm font-semibold" style="color: var(--color-text-primary)">Informasi akun</h4>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label for="user-form-username" class="mb-1 block text-sm font-medium text-secondary">
              Username <span class="text-red-600">*</span>
            </label>
            <InputText id="user-form-username" v-model="form.username" class="w-full text-sm" placeholder="Contoh: nama.pengguna" autocomplete="username" />
            <small v-if="errors.username" class="mt-1 block text-xs text-red-700">{{ errors.username }}</small>
          </div>
          <div>
            <label for="user-form-password" class="mb-1 block text-sm font-medium text-secondary">
              Password <span v-if="!isEdit" class="text-red-600">*</span>
            </label>
            <InputText
              id="user-form-password"
              v-model="form.password"
              type="password"
              class="w-full text-sm"
              :placeholder="isEdit ? 'Kosongkan jika tidak diubah' : 'Minimal 8 karakter'"
              autocomplete="new-password"
            />
            <small v-if="errors.password" class="mt-1 block text-xs text-red-700">{{ errors.password }}</small>
            <small v-else-if="isEdit" class="mt-1 block text-xs text-secondary">Biarkan kosong untuk mempertahankan password saat ini.</small>
          </div>
          <div class="sm:col-span-2">
            <label for="user-form-name" class="mb-1 block text-sm font-medium text-secondary">
              Nama lengkap <span class="text-red-600">*</span>
            </label>
            <InputText id="user-form-name" v-model="form.nama" class="w-full text-sm" placeholder="Masukkan nama lengkap" autocomplete="name" />
            <small v-if="errors.nama" class="mt-1 block text-xs text-red-700">{{ errors.nama }}</small>
          </div>
        </div>
      </div>

      <div>
        <h4 class="mb-3 text-sm font-semibold" style="color: var(--color-text-primary)">Informasi profesi</h4>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label for="user-form-profession" class="mb-1 block text-sm font-medium text-secondary">
              Profesi <span class="text-red-600">*</span>
            </label>
            <Select
              inputId="user-form-profession"
              v-model="form.profesi"
              :options="profesiOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Pilih profesi"
              class="w-full text-sm"
            />
            <small v-if="errors.profesi" class="mt-1 block text-xs text-red-700">{{ errors.profesi }}</small>
          </div>
          <div v-if="form.profesi === 'Dokter'">
            <label for="user-form-specialization" class="mb-1 block text-sm font-medium text-secondary">Spesialisasi</label>
            <InputText id="user-form-specialization" v-model="form.spesialisasi" class="w-full text-sm" placeholder="Contoh: Anak" />
          </div>
          <div>
            <label for="user-form-str" class="mb-1 block text-sm font-medium text-secondary">
              No. STR <span class="text-red-600">*</span>
            </label>
            <InputText id="user-form-str" v-model="form.no_str" class="w-full text-sm font-mono" placeholder="Masukkan nomor STR" />
            <small v-if="errors.no_str" class="mt-1 block text-xs text-red-700">{{ errors.no_str }}</small>
          </div>
        </div>
      </div>

      <div>
        <h4 class="mb-3 text-sm font-semibold" style="color: var(--color-text-primary)">Akses dan penempatan</h4>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label for="user-form-role" class="mb-1 block text-sm font-medium text-secondary">
              Role <span class="text-red-600">*</span>
            </label>
            <Select
              inputId="user-form-role"
              v-model="form.id_role"
              :options="roles"
              optionLabel="nama_role"
              optionValue="id_role"
              placeholder="Pilih role"
              class="w-full text-sm"
            />
            <small v-if="errors.id_role" class="mt-1 block text-xs text-red-700">{{ errors.id_role }}</small>
          </div>
          <div>
            <label for="user-form-unit" class="mb-1 block text-sm font-medium text-secondary">
              Unit <span class="text-red-600">*</span>
            </label>
            <Select
              inputId="user-form-unit"
              v-model="form.id_unit"
              :options="units"
              optionLabel="nama_unit"
              optionValue="id_unit"
              placeholder="Pilih unit"
              class="w-full text-sm"
            />
            <small v-if="errors.id_unit" class="mt-1 block text-xs text-red-700">{{ errors.id_unit }}</small>
          </div>
        </div>
      </div>

      <div class="flex flex-col-reverse gap-2 border-t pt-4 sm:flex-row sm:justify-end" style="border-color: var(--color-border)">
        <Button label="Batal" severity="secondary" outlined class="text-sm" :disabled="saving" @click="close" />
        <Button
          :label="isEdit ? 'Simpan Perubahan' : 'Tambah Pengguna'"
          icon="pi pi-check"
          type="submit"
          class="text-sm"
          :loading="saving"
          :disabled="saving"
        />
      </div>
    </form>
  </Dialog>
</template>

<script setup>
import { watch, reactive, computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'

const props = defineProps({
  show: { type: Boolean, default: false },
  user: { type: Object, default: null },
  roles: { type: Array, default: () => [] },
  units: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})

const emit = defineEmits(['update:show', 'saved'])

const form = reactive({
  username: '',
  password: '',
  nama: '',
  profesi: '',
  spesialisasi: '',
  no_str: '',
  id_role: null,
  id_unit: null
})
const errors = reactive({})

const profesiOptions = [
  { label: 'Dokter', value: 'Dokter' },
  { label: 'Perawat', value: 'Perawat' },
  { label: 'Petugas RM', value: 'Petugas RM' },
  { label: 'Koder', value: 'Koder' },
  { label: 'Staf Keuangan', value: 'Staf Keuangan' },
  { label: 'Admin TI', value: 'Admin TI' }
]

const isEdit = computed(() => !!props.user)
const title = computed(() => isEdit.value ? 'Edit Pengguna' : 'Tambah Pengguna')

watch(() => props.show, (val) => {
  if (val) {
    if (props.user) {
      Object.assign(form, {
        username: props.user.username,
        password: '',
        nama: props.user.nama,
        profesi: props.user.profesi,
        spesialisasi: props.user.spesialisasi || '',
        no_str: props.user.no_str,
        id_role: props.user.id_role,
        id_unit: props.user.id_unit
      })
    } else {
      Object.assign(form, {
        username: '',
        password: '',
        nama: '',
        profesi: '',
        spesialisasi: '',
        no_str: '',
        id_role: null,
        id_unit: null
      })
    }
    Object.keys(errors).forEach(k => delete errors[k])
  }
}, { immediate: true })

function validate() {
  Object.keys(errors).forEach(k => delete errors[k])
  if (!form.username.trim()) errors.username = 'Username wajib diisi'
  if (!isEdit.value && !form.password) errors.password = 'Password wajib diisi'
  if (!isEdit.value && form.password && form.password.length < 8) errors.password = 'Password minimal 8 karakter'
  if (!form.nama.trim()) errors.nama = 'Nama wajib diisi'
  if (!form.profesi) errors.profesi = 'Profesi wajib dipilih'
  if (!form.no_str.trim()) errors.no_str = 'No. STR wajib diisi'
  if (!form.id_role) errors.id_role = 'Role wajib dipilih'
  if (!form.id_unit) errors.id_unit = 'Unit wajib dipilih'
  return Object.keys(errors).length === 0
}

function submit() {
  if (!validate()) return
  const payload = { ...form }
  if (isEdit.value && !payload.password) delete payload.password
  emit('saved', { ...payload, id: props.user?.id_user })
}

</script>
