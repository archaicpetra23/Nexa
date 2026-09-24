import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const isAuthenticated = ref(false)

  async function login(username, password) {
    try {
      const response = await client.post('/auth/login', { username, password })
      if (response.data.success) {
        user.value = response.data.data.user
        isAuthenticated.value = true
        return { success: true }
      }
      return { success: false, message: response.data.message }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Login failed' }
    }
  }

  async function fetchMe() {
    try {
      const response = await client.get('/auth/me')
      if (response.data.success) {
        user.value = response.data.data
        isAuthenticated.value = true
      }
    } catch (error) {
      user.value = null
      isAuthenticated.value = false
    }
  }

  async function logout() {
    try {
      await client.post('/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    }
    user.value = null
    isAuthenticated.value = false
  }

  return { user, isAuthenticated, login, logout, fetchMe }
})