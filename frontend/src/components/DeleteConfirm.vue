<template>
  <Dialog v-model:visible="visible" modal :header="title" :style="{ width: '28rem' }">
    <p class="text-sm" style="color: var(--color-text-secondary)">{{ message }}</p>
    <template #footer>
      <Button label="Batal" severity="secondary" text @click="visible = false" class="text-sm" />
      <Button label="Hapus" @click="confirm" class="text-sm" :style="{ backgroundColor: '#B54245' }" />
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
