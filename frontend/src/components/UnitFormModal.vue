<template>
  <Dialog v-model:visible="visible" modal :header="title" :style="{ width: '28rem' }">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium mb-1" style="color: var(--color-text-secondary)">
          Nama Unit <span style="color: #B54245">*</span>
        </label>
        <InputText v-model="form.nama_unit" class="w-full text-sm" placeholder="Nama unit" />
        <small v-if="errors.nama_unit" class="text-xs" style="color: #B54245">{{ errors.nama_unit }}</small>
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

const props = defineProps({
  show: { type: Boolean, default: false },
  unit: { type: Object, default: null }
})

const emit = defineEmits(['update:show', 'saved'])

const visible = ref(props.show)
const loading = ref(false)
const form = reactive({ nama_unit: '' })
const errors = reactive({})

const isEdit = computed(() => !!props.unit)
const title = computed(() => isEdit.value ? 'Edit Unit' : 'Tambah Unit')

watch(() => props.show, (val) => {
  visible.value = val
  if (val) {
    if (props.unit) {
      form.nama_unit = props.unit.nama_unit
    } else {
      form.nama_unit = ''
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
  emit('saved', form)
  loading.value = false
  visible.value = false
}
</script>
