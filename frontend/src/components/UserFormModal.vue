<template>
  <Dialog v-model:visible="visible" modal :header="title" :style="{ width: '36rem' }">
    <div class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
            Username <span style="color: #B54245">*</span>
          </label>
          <InputText v-model="form.username" class="w-full text-sm" placeholder="Username" />
          <small v-if="errors.username" class="text-xs" style="color: #B54245">{{ errors.username }}</small>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
            Password <span v-if="!isEdit" style="color: #B54245">*</span>
          </label>
          <InputText v-model="form.password" type="password" class="w-full text-sm" :placeholder="isEdit ? 'Kosongkan jika tidak diubah' : 'Password'" />
          <small v-if="errors.password" class="text-xs" style="color: #B54245">{{ errors.password }}</small>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
          Nama <span style="color: #B54245">*</span>
        </label>
        <InputText v-model="form.nama" class="w-full text-sm" placeholder="Nama lengkap" />
        <small v-if="errors.nama" class="text-xs" style="color: #B54245">{{ errors.nama }}</small>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
            Profesi <span style="color: #B54245">*</span>
          </label>
          <Select v-model="form.profesi" :options="profesiOptions" optionLabel="label" optionValue="value" placeholder="Pilih profesi" class="w-full text-sm" />
          <small v-if="errors.profesi" class="text-xs" style="color: #B54245">{{ errors.profesi }}</small>
        </div>
        <div v-if="form.profesi === 'Dokter'">
          <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">Spesialisasi</label>
          <InputText v-model="form.spesialisasi" class="w-full text-sm" placeholder="Spesialisasi" />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
          No. STR <span style="color: #B54245">*</span>
        </label>
        <InputText v-model="form.no_str" class="w-full text-sm font-mono" placeholder="Nomor STR" />
        <small v-if="errors.no_str" class="text-xs" style="color: #B54245">{{ errors.no_str }}</small>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
            Role <span style="color: #B54245">*</span>
          </label>
          <Select v-model="form.id_role" :options="roles" optionLabel="nama_role" optionValue="id_role" placeholder="Pilih role" class="w-full text-sm" />
          <small v-if="errors.id_role" class="text-xs" style="color: #B54245">{{ errors.id_role }}</small>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
            Unit <span style="color: #B54245">*</span>
          </label>
          <Select v-model="form.id_unit" :options="units" optionLabel="nama_unit" optionValue="id_unit" placeholder="Pilih unit" class="w-full text-sm" />
          <small v-if="errors.id_unit" class="text-xs" style="color: #B54245">{{ errors.id_unit }}</small>
        </div>
      </div>
    </div>
    <template #footer>
      <Button label="Batal" severity="secondary" text @click="visible = false" class="text-sm" />
      <Button label="Simpan" @click="submit" :loading="loading" class="text-sm" />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch, reactive, computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'

const props = defineProps({
  show: { type: Boolean, default: false },
  user: { type: Object, default: null },
  roles: { type: Array, default: () => [] },
  units: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:show', 'saved'])

const visible = ref(props.show)
const loading = ref(false)
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
  visible.value = val
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
})

watch(visible, (val) => { emit('update:show', val) })

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

async function submit() {
  if (!validate()) return
  loading.value = true
  const payload = { ...form }
  if (isEdit.value && !payload.password) delete payload.password
  emit('saved', { ...payload, id: props.user?.id_user })
  loading.value = false
  visible.value = false
}
</script>
