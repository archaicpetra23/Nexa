<template>
  <div class="modern-table" :aria-busy="loading">
    <div class="table-viewport">
      <table :aria-label="ariaLabel">
        <thead>
          <tr><slot name="header" /></tr>
        </thead>
        <tbody v-if="displayRows.length">
          <tr v-for="(item, index) in displayRows" :key="rowKey(item, index)">
            <slot name="row" :item="item" :index="first + index" />
          </tr>
        </tbody>
      </table>

      <div v-if="loading" class="table-state table-loading">
        <i class="pi pi-spin pi-spinner" aria-hidden="true"></i>
        <span>Memuat data...</span>
      </div>
      <div v-else-if="!displayRows.length" class="table-state">
        <i class="pi pi-inbox" aria-hidden="true"></i>
        <strong>Belum ada data</strong>
        <span>Data akan muncul di sini setelah tersedia.</span>
      </div>
    </div>

    <footer v-if="paginator" class="table-pagination">
      <span class="pagination-summary">
        {{ totalRecords ? `Menampilkan ${first + 1}–${Math.min(first + rows, totalRecords)} dari ${totalRecords}` : 'Tidak ada data' }}
      </span>
      <div class="pagination-controls">
        <label class="page-size">
          <span>Baris</span>
          <select :value="rows" aria-label="Jumlah baris per halaman" @change="changeRows">
            <option v-for="option in rowsPerPageOptions" :key="option" :value="option">{{ option }}</option>
          </select>
        </label>
        <button
          type="button"
          class="page-button"
          aria-label="Halaman sebelumnya"
          :disabled="currentPage <= 0 || loading"
          @click="goToPage(currentPage - 1)"
        >
          <i class="pi pi-angle-left" aria-hidden="true"></i>
        </button>
        <span class="page-indicator">{{ totalRecords ? currentPage + 1 : 0 }} / {{ pageCount }}</span>
        <button
          type="button"
          class="page-button"
          aria-label="Halaman berikutnya"
          :disabled="currentPage >= pageCount - 1 || loading"
          @click="goToPage(currentPage + 1)"
        >
          <i class="pi pi-angle-right" aria-hidden="true"></i>
        </button>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  paginator: { type: Boolean, default: false },
  rows: { type: Number, default: 10 },
  first: { type: Number, default: 0 },
  totalRecords: { type: Number, default: 0 },
  rowsPerPageOptions: { type: Array, default: () => [10, 20, 50] },
  lazy: { type: Boolean, default: false },
  dataKey: { type: String, default: '' },
  ariaLabel: { type: String, default: 'Data tabel' }
})

const emit = defineEmits(['page'])

const currentPage = computed(() => Math.floor(props.first / props.rows))
const pageCount = computed(() => Math.max(1, Math.ceil(props.totalRecords / props.rows)))
const displayRows = computed(() => {
  if (props.lazy || !props.paginator) return props.value
  return props.value.slice(props.first, props.first + props.rows)
})

function rowKey(item, index) {
  return props.dataKey ? item[props.dataKey] : `${props.first + index}`
}

function goToPage(page) {
  if (page < 0 || page >= pageCount.value) return
  emit('page', { page, first: page * props.rows, rows: props.rows })
}

function changeRows(event) {
  const nextRows = Number(event.target.value)
  emit('page', { page: 0, first: 0, rows: nextRows })
}
</script>

<style scoped>
.modern-table {
  display: flex;
  width: 100%;
  height: 100%;
  flex: 1 1 auto;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid #e3e8ea;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 4px 14px rgb(31 41 55 / 5%);
}

.table-viewport {
  position: relative;
  min-width: 0;
  min-height: 0;
  flex: 1 1 auto;
  overflow: auto;
  overscroll-behavior: contain;
  scrollbar-color: #bdc9d0 #f4f6f7;
  scrollbar-width: thin;
  -webkit-overflow-scrolling: touch;
}

.table-viewport::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.table-viewport::-webkit-scrollbar-track {
  background: #f4f6f7;
}

.table-viewport::-webkit-scrollbar-thumb {
  border: 2px solid #f4f6f7;
  border-radius: 999px;
  background: #bdc9d0;
}

table {
  width: max-content;
  min-width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  text-align: left;
  white-space: nowrap;
}

:slotted(th) {
  position: sticky;
  z-index: 2;
  top: 0;
  padding: 12px 16px;
  border-right: 1px solid #e7ecee;
  border-bottom: 1px solid #e4e9eb;
  background: #f7f9fa;
  color: #62717e;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.075em;
  text-transform: uppercase;
}

:slotted(th:last-child) {
  border-right: 0;
}

:slotted(td) {
  padding: 12px 16px;
  border-right: 1px solid #edf0f1;
  border-bottom: 1px solid #edf0f1;
  color: #35424c;
  font-size: 12px;
  vertical-align: middle;
}

:slotted(td:last-child) {
  border-right: 0;
}

:slotted(.description-cell) {
  min-width: 280px;
  max-width: 460px;
  white-space: normal;
}

:slotted(.user-name-cell) {
  color: #253746;
  font-weight: 600;
}

:slotted(.amount-cell) {
  color: #253746;
  font-weight: 600;
  text-align: right;
}

:slotted(.role-cell) {
  min-width: 150px;
}

tbody tr {
  transition: background-color 140ms ease;
}

tbody tr:nth-child(even) {
  background: #fbfcfc;
}

tbody tr:hover {
  background: #f1f6f8;
}

tbody tr:hover :slotted(td) {
  background: #f1f6f8;
}

tbody tr:last-child td {
  border-bottom: 0;
}

.table-state {
  position: absolute;
  inset: 44px 0 0;
  display: flex;
  min-height: 150px;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8px;
  background: #fff;
  color: #798691;
  font-size: 12px;
  text-align: center;
}

.table-state > i {
  color: #91a3af;
  font-size: 22px;
}

.table-state strong {
  color: #34424d;
  font-size: 13px;
}

.table-loading {
  background: rgb(255 255 255 / 82%);
}

.table-loading > i {
  color: #1c3d5a;
}

.table-pagination {
  display: flex;
  min-height: 52px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 14px;
  border-top: 1px solid #e9edef;
  background: #fff;
}

.pagination-summary {
  color: #75818a;
  font-size: 11px;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 7px;
}

.page-size {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-right: 4px;
  color: #75818a;
  font-size: 11px;
}

.page-size select {
  height: 30px;
  padding: 0 7px;
  border: 1px solid #e1e6e8;
  border-radius: 6px;
  background: #fff;
  color: #34424d;
  font-size: 11px;
}

.page-button {
  display: inline-grid;
  width: 30px;
  height: 30px;
  place-items: center;
  border: 1px solid #e1e6e8;
  border-radius: 7px;
  background: white;
  color: #435360;
  cursor: pointer;
}

.page-button:hover:not(:disabled) {
  border-color: #bdcbd3;
  background: #f4f8fa;
  color: #1c3d5a;
}

.page-button:disabled {
  color: #b8c0c5;
  cursor: not-allowed;
}

.page-indicator {
  min-width: 44px;
  color: #52616c;
  font-size: 11px;
  font-variant-numeric: tabular-nums;
  text-align: center;
}

@media (max-width: 640px) {
  .table-pagination {
    min-height: 48px;
    padding: 7px 9px;
  }

  .pagination-summary {
    font-size: 10px;
  }

  .page-size {
    gap: 4px;
    margin-right: 0;
  }
}
</style>
