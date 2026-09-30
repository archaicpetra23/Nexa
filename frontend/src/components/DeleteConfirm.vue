<template>
  <Dialog v-model:visible="visible" modal :header="title" :style="{ width: '28rem' }">
    <p class="text-sm text-secondary">{{ message }}</p>
    <template #footer>
      <Button label="Batal" severity="secondary" text @click="visible = false" />
      <Button label="Hapus" severity="danger" @click="confirm" />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'

const props = defineProps({
  show: { type: Boolean, default: false },
  title: { type: String, default: 'Konfirmasi Hapus' },
  message: { type: String, default: 'Apakah Anda yakin ingin menghapus data ini?' }
})

const emit = defineEmits(['update:show', 'confirm'])

const visible = ref(props.show)

watch(() => props.show, (val) => { visible.value = val })
watch(visible, (val) => { emit('update:show', val) })

function confirm() {
  emit('confirm')
  visible.value = false
}
</script>
