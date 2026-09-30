<template>
  <div class="auth-page">
    <router-link to="/" class="auth-brand" aria-label="Kembali ke beranda Nexa">
      <span class="auth-brand-mark"><i class="pi pi-plus" aria-hidden="true"></i></span>
      <span>nexa<span>.</span></span>
    </router-link>
    <main class="auth-main">
      <div class="auth-aside">
        <span class="auth-aside-kicker">CASEMIX MANAGEMENT SYSTEM</span>
        <h1>Satu alur yang menghubungkan setiap layanan.</h1>
        <p>Ruang kerja terintegrasi untuk membantu tim rumah sakit mengelola data klinis dan proses klaim dengan lebih jelas.</p>
        <div class="auth-aside-rule"></div>
        <div class="auth-assurance">
          <i class="pi pi-shield" aria-hidden="true"></i>
          <span><strong>Akses terjaga</strong><small>Hak akses menyesuaikan peran Anda</small></span>
        </div>
      </div>
      <div class="auth-form-wrap">
        <slot />
        <p class="auth-help">Belum memiliki akun? <span>Hubungi administrator TI rumah sakit.</span></p>
      </div>
    </main>
    <footer class="auth-footer">
      <span>Nexa Casemix Management System</span>
      <span><i class="pi pi-lock" aria-hidden="true"></i> Sesi terlindungi</span>
    </footer>
  </div>
</template>

<script setup>
</script>

<style scoped>
.auth-page {
  display: flex;
  min-height: 100vh;
  flex-direction: column;
  padding: 0 max(7vw, 24px);
  background: var(--color-bg-base);
}

.auth-brand {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  gap: 10px;
  margin-top: 30px;
  color: var(--color-primary);
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.8px;
  text-decoration: none;
}

.auth-brand > span:last-child > span {
  color: var(--color-status-pending);
}

.auth-brand-mark {
  display: grid;
  width: 32px;
  height: 32px;
  place-items: center;
  border-radius: 6px;
  background: var(--color-primary);
  color: white;
  font-size: 14px;
}

.auth-main {
  display: grid;
  width: 100%;
  max-width: 1080px;
  flex: 1;
  grid-template-columns: 1fr 0.78fr;
  align-items: center;
  gap: clamp(50px, 10vw, 130px);
  margin: 0 auto;
  padding: 52px 0;
}

.auth-aside {
  max-width: 480px;
  animation: page-enter 300ms ease-out both;
}

.auth-aside-kicker {
  color: var(--color-text-secondary);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1.2px;
}

.auth-aside h1 {
  margin: 18px 0 14px;
  color: var(--color-primary-dark);
  font-size: clamp(36px, 4vw, 51px);
  font-weight: 600;
  letter-spacing: -2px;
  line-height: 1.13;
}

.auth-aside > p {
  max-width: 420px;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 14px;
  line-height: 1.75;
}

.auth-aside-rule {
  width: 100%;
  height: 1px;
  margin: 31px 0 21px;
  background: var(--color-border);
}

.auth-assurance {
  display: flex;
  align-items: center;
  gap: 13px;
  color: var(--color-primary);
}

.auth-assurance > i {
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: var(--color-surface);
}

.auth-assurance span {
  display: grid;
  gap: 2px;
}

.auth-assurance strong {
  color: var(--color-text-primary);
  font-size: 12px;
  font-weight: 600;
}

.auth-assurance small {
  color: var(--color-text-secondary);
  font-size: 11px;
}

.auth-form-wrap {
  width: 100%;
  max-width: 400px;
  justify-self: end;
  animation: page-enter 340ms 60ms ease-out both;
}

.auth-help {
  margin: 16px 0 0;
  color: var(--color-text-secondary);
  font-size: 11px;
  text-align: center;
}

.auth-help span {
  color: var(--color-primary);
  font-weight: 500;
}

.auth-footer {
  display: flex;
  justify-content: space-between;
  padding: 16px 0 20px;
  border-top: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: 10px;
}

.auth-footer > span:last-child {
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.auth-footer i {
  color: var(--color-primary);
}

@media (max-width: 760px) {
  .auth-page {
    padding: 0 20px;
  }

  .auth-brand {
    margin-top: 20px;
  }

  .auth-main {
    max-width: 440px;
    grid-template-columns: 1fr;
    gap: 30px;
    padding: 35px 0;
  }

  .auth-aside {
    text-align: center;
  }

  .auth-aside h1 {
    margin: 10px 0;
    font-size: 31px;
    letter-spacing: -1.2px;
  }

  .auth-aside > p,
  .auth-assurance,
  .auth-aside-rule {
    display: none;
  }

  .auth-form-wrap {
    max-width: none;
  }
}
</style>
