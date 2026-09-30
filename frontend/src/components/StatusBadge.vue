<template>
  <span
    :style="{ backgroundColor: bgColor, color: textColor }"
    class="inline-flex items-center gap-1.5 px-2.5 py-0.5 text-xs font-medium"
    style="border-radius: 999px"
  >
    <i :class="icon" class="text-xs"></i>
    <span>{{ label }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: { type: String, required: true },
  type: { type: String, default: 'claim' }
})

// Exact colors from Design_UI.md §2
const config = {
  claim: {
    draft: { icon: 'pi pi-file', bg: '#F3F4F6', text: '#9CA3AF', label: 'Draft' },
    pending: { icon: 'pi pi-clock', bg: '#FEF3C7', text: '#B7791F', label: 'Pending' },
    disetujui: { icon: 'pi pi-check', bg: '#D1FAE5', text: '#2F7A4F', label: 'Disetujui' },
    ditolak: { icon: 'pi pi-times', bg: '#FEE2E2', text: '#B54245', label: 'Ditolak' }
  },
  user: {
    aktif: { icon: 'pi pi-check', bg: '#D1FAE5', text: '#2F7A4F', label: 'Aktif' },
    nonaktif: { icon: 'pi pi-times', bg: '#FEE2E2', text: '#B54245', label: 'Nonaktif' }
  },
  activity: {
    INSERT: { icon: 'pi pi-plus', bg: '#DBEAFE', text: '#1E40AF', label: 'INSERT' },
    UPDATE: { icon: 'pi pi-pencil', bg: '#FEF3C7', text: '#B7791F', label: 'UPDATE' },
    DELETE: { icon: 'pi pi-trash', bg: '#FEE2E2', text: '#B54245', label: 'DELETE' },
    LOGIN: { icon: 'pi pi-sign-in', bg: '#D1FAE5', text: '#2F7A4F', label: 'LOGIN' },
    LOGOUT: { icon: 'pi pi-sign-out', bg: '#F3F4F6', text: '#9CA3AF', label: 'LOGOUT' }
  }
}

const item = computed(() => config[props.type]?.[props.status] || config.claim.draft)
const bgColor = computed(() => item.value.bg)
const textColor = computed(() => item.value.text)
const icon = computed(() => item.value.icon)
const label = computed(() => item.value.label)
</script>
