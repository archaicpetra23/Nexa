<template>
  <span :class="classes" class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium">
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

const config = {
  claim: {
    draft: { icon: 'pi pi-file', color: 'bg-gray-100 text-gray-700', label: 'Draft' },
    pending: { icon: 'pi pi-clock', color: 'bg-amber-100 text-amber-800', label: 'Pending' },
    disetujui: { icon: 'pi pi-check', color: 'bg-green-100 text-green-700', label: 'Disetujui' },
    ditolak: { icon: 'pi pi-times', color: 'bg-red-100 text-red-700', label: 'Ditolak' }
  },
  user: {
    aktif: { icon: 'pi pi-check', color: 'bg-green-100 text-green-700', label: 'Aktif' },
    nonaktif: { icon: 'pi pi-times', color: 'bg-red-100 text-red-700', label: 'Nonaktif' }
  },
  activity: {
    INSERT: { icon: 'pi pi-plus', color: 'bg-blue-100 text-blue-700', label: 'INSERT' },
    UPDATE: { icon: 'pi pi-pencil', color: 'bg-amber-100 text-amber-800', label: 'UPDATE' },
    DELETE: { icon: 'pi pi-trash', color: 'bg-red-100 text-red-700', label: 'DELETE' },
    LOGIN: { icon: 'pi pi-sign-in', color: 'bg-green-100 text-green-700', label: 'LOGIN' },
    LOGOUT: { icon: 'pi pi-sign-out', color: 'bg-gray-100 text-gray-700', label: 'LOGOUT' }
  }
}

const item = computed(() => config[props.type]?.[props.status] || config.claim.draft)
const classes = computed(() => item.value.color)
const icon = computed(() => item.value.icon)
const label = computed(() => item.value.label)
</script>
