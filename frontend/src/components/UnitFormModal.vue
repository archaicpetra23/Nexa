<template>
  <Dialog v-model:visible="visible" modal :header="title" :style="{ width: '28rem' }">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-secondary mb-1">
          Nama Unit <span class="text-red-600">*</span>
        </label>
        <InputText v-model="form.nama_unit" class="w-full" placeholder="Nama unit" />
        <small v-if="errors.nama_unit" class="text-red-600">{{ errors.nama_unit }}</small>
      </div>
    </div>
    <template #footer>
      <Button label="Batal" severity="secondary" text @click="visible = false" />
      <Button label="Simpan" @click="submit" :loading="loading" />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch, reactive } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'

const props = defineProps({
  show: { type: Boolean, default: false },
  unit: { type: Object, default: null }
})

const emit = defineEmits(['update:show', 'saved'])

const visible = ref(props.show)
const loading = ref(false)
const form = reactive({ nama_unit: '' })
const errors = reactive({})

const title = ref('Tambah Unit')

watch(() => props.show, (val) => {
  visible.value = val
  if (val) {
    if (props.unit) {
      form.nama_unit = props.unit.nama_unit
      title.value = 'Edit Unit'
    } else {
      form.nama_unit = ''
      title.value = 'Tambah Unit'
    }
    Object.keys(errors).forEach(k => delete errors[k])
  }
})

watch(visible, (val) => { emit('update:show', val) })

function validate() {
  Object.keys(errors).forEach(k => delete errors[k])
  if (!form.nama_unit.trim()) errors.nama_unit = 'Nama unit wajib diisi'
  return Object.keys(errors).length === 0
}

async function submit() {
  if (!validate()) return
  loading.value = true
  emit('saved', { ...form, id: props.unit?.id_unit })
  loading.value = false
  visible.value = false
}
</script>
