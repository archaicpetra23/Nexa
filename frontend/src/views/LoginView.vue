<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-gray-800 mb-6">Nexa</h1>
      <h2 class="text-lg text-center text-gray-600 mb-8">Sistem Casemix Terintegrasi</h2>

      <form @submit.prevent="handleLogin" novalidate>
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2">Username</label>
          <input
            v-model="username"
            type="text"
            :class="['w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500',
              usernameError ? 'border-red-400' : 'border-gray-300']"
            placeholder="Masukkan username"
            autocomplete="username"
          />
          <p v-if="usernameError" class="text-sm text-red-600 mt-1">{{ usernameError }}</p>
        </div>

        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-bold mb-2">Password</label>
          <input
            v-model="password"
            type="password"
            :class="['w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500',
              passwordError ? 'border-red-400' : 'border-gray-300']"
            placeholder="Masukkan password"
            autocomplete="current-password"
          />
          <p v-if="passwordError" class="text-sm text-red-600 mt-1">{{ passwordError }}</p>
        </div>

        <div v-if="error" class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {{ error }}
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50"
        >
          {{ loading ? 'Loading...' : 'Login' }}
        </button>
      </form>

      <div class="mt-6 text-center text-sm text-gray-600">
        <p class="mt-2">Hubungi admin TI untuk akun pengguna.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const usernameError = ref('')
const passwordError = ref('')

watch(username, () => { usernameError.value = ''; error.value = '' })
watch(password, () => { passwordError.value = ''; error.value = '' })

function validate() {
  let valid = true
  usernameError.value = ''
  passwordError.value = ''

  if (!username.value.trim()) {
    usernameError.value = 'Username wajib diisi'
    valid = false
  } else if (username.value.trim().length < 3) {
    usernameError.value = 'Username minimal 3 karakter'
    valid = false
  }
  if (!password.value) {
    passwordError.value = 'Password wajib diisi'
    valid = false
  } else if (password.value.length < 6) {
    passwordError.value = 'Password minimal 6 karakter'
    valid = false
  }
  return valid
}

async function handleLogin() {
  if (!validate()) return

  loading.value = true
  error.value = ''

  const result = await authStore.login(username.value.trim(), password.value)

  if (result.success) {
    router.push('/dashboard')
  } else {
    error.value = result.message || 'Login gagal. Periksa username dan password.'
  }

  loading.value = false
}
</script>
