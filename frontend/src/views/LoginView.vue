<template>
  <AuthLayout>
    <section class="login-card">
      <div class="login-card-heading">
        <span class="login-card-icon"><i class="pi pi-sign-in" aria-hidden="true"></i></span>
        <span class="login-card-kicker">SELAMAT DATANG KEMBALI</span>
      </div>
      <h2>Masuk ke Nexa</h2>
      <p class="login-intro">Gunakan akun staf rumah sakit untuk melanjutkan.</p>

      <form class="login-form" @submit.prevent="handleLogin" novalidate>
        <div class="form-field">
          <label for="username">Username</label>
          <div :class="['input-wrap', { 'input-wrap-error': usernameError }]">
            <i class="pi pi-user input-icon" aria-hidden="true"></i>
            <input
              id="username"
              v-model="username"
              type="text"
              placeholder="Masukkan username"
              autocomplete="username"
              :aria-invalid="Boolean(usernameError)"
              :aria-describedby="usernameError ? 'username-error' : undefined"
              autofocus
            />
          </div>
          <p v-if="usernameError" id="username-error" class="field-error">{{ usernameError }}</p>
        </div>

        <div class="form-field">
          <label for="password">Password</label>
          <div :class="['input-wrap', { 'input-wrap-error': passwordError }]">
            <i class="pi pi-lock input-icon" aria-hidden="true"></i>
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Masukkan password"
              autocomplete="current-password"
              :aria-invalid="Boolean(passwordError)"
              :aria-describedby="passwordError ? 'password-error' : undefined"
            />
            <button
              class="password-toggle"
              type="button"
              :aria-label="showPassword ? 'Sembunyikan password' : 'Tampilkan password'"
              @click="showPassword = !showPassword"
            >
              <i :class="showPassword ? 'pi pi-eye-slash' : 'pi pi-eye'" aria-hidden="true"></i>
            </button>
          </div>
          <p v-if="passwordError" id="password-error" class="field-error">{{ passwordError }}</p>
        </div>

        <div v-if="error" class="login-error" role="alert">
          <i class="pi pi-exclamation-circle" aria-hidden="true"></i>
          <span>{{ error }}</span>
        </div>

        <button class="login-submit" type="submit" :disabled="loading">
          <i v-if="loading" class="pi pi-spinner pi-spin" aria-hidden="true"></i>
          <span>{{ loading ? 'Memeriksa akun...' : 'Masuk' }}</span>
          <i v-if="!loading" class="pi pi-arrow-right" aria-hidden="true"></i>
        </button>
      </form>

      <div class="login-divider"><span></span><small>AKSES TERBATAS</small><span></span></div>
      <p class="login-security"><i class="pi pi-shield" aria-hidden="true"></i> Aktivitas akun tercatat untuk keamanan sistem.</p>
    </section>
  </AuthLayout>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import AuthLayout from '../layouts/AuthLayout.vue'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
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

  try {
    const result = await authStore.login(username.value.trim(), password.value)
    if (result.success) {
      router.push('/dashboard')
    } else {
      error.value = result.message || 'Login gagal. Periksa username dan password.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-card {
  padding: 30px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: var(--color-surface);
  box-shadow: var(--shadow-panel);
}

.login-card-heading {
  display: flex;
  align-items: center;
  gap: 11px;
  margin-bottom: 16px;
}

.login-card-icon {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border-radius: 6px;
  background: #EEF2F5;
  color: var(--color-primary);
}

.login-card-kicker {
  color: var(--color-text-secondary);
  font-size: 9px;
  font-weight: 600;
  letter-spacing: 1px;
}

.login-card h2 {
  margin: 0;
  color: var(--color-primary-dark);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.5px;
}

.login-intro {
  margin: 6px 0 23px;
  color: var(--color-text-secondary);
  font-size: 12px;
}

.login-form {
  display: grid;
  gap: 16px;
}

.form-field label {
  display: block;
  margin-bottom: 6px;
  color: var(--color-text-primary);
  font-size: 12px;
  font-weight: 500;
}

.input-wrap {
  display: flex;
  min-height: 44px;
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-base);
  background: white;
  transition: border-color var(--motion-fast), box-shadow var(--motion-fast);
}

.input-wrap:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgb(28 61 90 / 9%);
}

.input-wrap-error {
  border-color: var(--color-status-rejected);
}

.input-icon {
  color: var(--color-text-disabled);
  font-size: 13px;
}

.input-wrap input {
  width: 100%;
  min-width: 0;
  height: 42px;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--color-text-primary);
  font-size: 12px;
}

.input-wrap input::placeholder {
  color: var(--color-text-disabled);
}

.input-wrap input:focus-visible {
  outline: 0;
}

.password-toggle {
  display: grid;
  width: 34px;
  height: 34px;
  flex: 0 0 auto;
  place-items: center;
  border: 0;
  border-radius: 4px;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
}

.password-toggle:hover {
  background: var(--color-bg-base);
  color: var(--color-primary);
}

.field-error {
  margin: 5px 0 0;
  color: var(--color-status-rejected);
  font-size: 11px;
}

.login-error {
  display: flex;
  align-items: flex-start;
  gap: 9px;
  padding: 10px 11px;
  border: 1px solid #E9C7C8;
  border-radius: var(--radius-base);
  background: #FCF4F4;
  color: #8F3033;
  font-size: 11px;
}

.login-error i {
  margin-top: 1px;
}

.login-submit {
  display: flex;
  min-height: 44px;
  align-items: center;
  justify-content: center;
  gap: 9px;
  margin-top: 2px;
  border: 0;
  border-radius: var(--radius-base);
  background: var(--color-primary);
  color: white;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color var(--motion-fast);
}

.login-submit:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.login-submit:disabled {
  cursor: wait;
  opacity: 0.7;
}

.login-divider {
  display: flex;
  align-items: center;
  gap: 11px;
  margin: 22px 0 13px;
}

.login-divider span {
  height: 1px;
  flex: 1;
  background: var(--color-border);
}

.login-divider small {
  color: var(--color-text-disabled);
  font-size: 8px;
  font-weight: 600;
  letter-spacing: 0.8px;
}

.login-security {
  display: flex;
  justify-content: center;
  gap: 7px;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 10px;
}

.login-security i {
  color: var(--color-primary);
}

@media (max-width: 760px) {
  .login-card {
    padding: 24px 20px;
  }
}
</style>
