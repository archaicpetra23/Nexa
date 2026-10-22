<template>
  <span
    :class="['status-badge', statusClass]"
    :aria-label="`Status: ${label}`"
  >
    <i :class="icon" aria-hidden="true"></i>
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
    draft: { icon: 'pi pi-file', class: 'status-draft', label: 'Draft' },
    pending: { icon: 'pi pi-clock', class: 'status-pending', label: 'Pending' },
    disetujui: { icon: 'pi pi-check', class: 'status-approved', label: 'Disetujui' },
    ditolak: { icon: 'pi pi-times', class: 'status-rejected', label: 'Ditolak' }
  },
  user: {
    aktif: { icon: 'pi pi-check', class: 'status-approved', label: 'Aktif' },
    nonaktif: { icon: 'pi pi-times', class: 'status-rejected', label: 'Nonaktif' }
  },
  activity: {
    INSERT: { icon: 'pi pi-plus', class: 'status-approved', label: 'INSERT' },
    UPDATE: { icon: 'pi pi-pencil', class: 'status-pending', label: 'UPDATE' },
    DELETE: { icon: 'pi pi-trash', class: 'status-rejected', label: 'DELETE' },
    LOGIN: { icon: 'pi pi-sign-in', class: 'status-approved', label: 'LOGIN' },
    LOGOUT: { icon: 'pi pi-sign-out', class: 'status-draft', label: 'LOGOUT' }
  }
}

const item = computed(() => config[props.type]?.[props.status] || config.claim.draft)
const statusClass = computed(() => item.value.class)
const icon = computed(() => item.value.icon)
const label = computed(() => item.value.label)
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.25rem 0.625rem;
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 500;
}

.status-badge i {
  font-size: var(--text-xs);
}

.status-draft {
  background-color: var(--color-status-draft-bg);
  color: var(--color-status-draft-fg);
}

.status-pending {
  background-color: var(--color-status-pending-bg);
  color: var(--color-status-pending-fg);
}

.status-approved {
  background-color: var(--color-status-approved-bg);
  color: var(--color-status-approved-fg);
}

.status-rejected {
  background-color: var(--color-status-rejected-bg);
  color: var(--color-status-rejected-fg);
}
</style>

