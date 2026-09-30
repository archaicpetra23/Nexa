import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'

export const useAdminStore = defineStore('admin', () => {
  const users = ref([])
  const roles = ref([])
  const units = ref([])
  const auditTrail = ref([])
  const claims = ref([])

  const usersTotal = ref(0)
  const auditTotal = ref(0)
  const claimsTotal = ref(0)

  async function fetchUsers(params = {}) {
    const res = await client.get('/admin/users', { params })
    if (res.data.success) {
      const data = res.data.data
      users.value = data.items || data
      usersTotal.value = data.total || users.value.length
    }
  }

  async function createUser(payload) {
    const res = await client.post('/admin/users', payload)
    return res.data
  }

  async function updateUser(id, payload) {
    const res = await client.put(`/admin/users/${id}`, payload)
    return res.data
  }

  async function deleteUser(id) {
    const res = await client.delete(`/admin/users/${id}`)
    return res.data
  }

  async function fetchRoles() {
    const res = await client.get('/master/roles')
    if (res.data.success) roles.value = res.data.data
  }

  async function fetchUnits() {
    const res = await client.get('/master/units')
    if (res.data.success) units.value = res.data.data
  }

  async function createUnit(payload) {
    const res = await client.post('/master/units', payload)
    return res.data
  }

  async function updateUnit(id, payload) {
    const res = await client.put(`/master/units/${id}`, payload)
    return res.data
  }

  async function deleteUnit(id) {
    const res = await client.delete(`/master/units/${id}`)
    return res.data
  }

  async function fetchAuditTrail(params = {}) {
    const res = await client.get('/admin/audit-trail', { params })
    if (res.data.success) {
      const data = res.data.data
      auditTrail.value = data.items || data
      auditTotal.value = data.total || auditTrail.value.length
    }
  }

  async function fetchClaims(params = {}) {
    const res = await client.get('/admin/submissions', { params })
    if (res.data.success) {
      const data = res.data.data
      claims.value = data.items || data
      claimsTotal.value = data.total || claims.value.length
    }
  }

  async function updateClaimStatus(id, payload) {
    const res = await client.patch(`/klaim/${id}/status`, payload)
    return res.data
  }

  return {
    users, roles, units, auditTrail, claims,
    usersTotal, auditTotal, claimsTotal,
    fetchUsers, createUser, updateUser, deleteUser,
    fetchRoles, fetchUnits, createUnit, updateUnit, deleteUnit,
    fetchAuditTrail, fetchClaims, updateClaimStatus
  }
})
