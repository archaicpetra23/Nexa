<template>
  <div class="landing-page">
    <header class="landing-header">
      <router-link to="/" class="brand-lockup" aria-label="Nexa, beranda">
        <span class="brand-mark"><i class="pi pi-plus" aria-hidden="true"></i></span>
        <span class="brand-name">nexa<span>.</span></span>
      </router-link>
      <div class="header-actions">
        <span class="header-caption"><i class="pi pi-shield" aria-hidden="true"></i> Sistem informasi rumah sakit</span>
        <router-link to="/login" class="button button-secondary">
          Masuk <i class="pi pi-arrow-up-right" aria-hidden="true"></i>
        </router-link>
      </div>
    </header>

    <main>
      <section class="landing-hero">
        <div class="hero-copy">
          <div class="eyebrow"><span></span> CASEMIX MANAGEMENT SYSTEM</div>
          <h1>Alur klinis dan klaim, <span>dalam satu ruang kerja.</span></h1>
          <p class="hero-description">
            Nexa menghubungkan proses rekam medis, koding, dan klaim BPJS agar tim rumah sakit dapat bekerja lebih terarah, cepat, dan akuntabel.
          </p>
          <div class="hero-actions">
            <router-link to="/login" class="button button-primary">
              Masuk ke Nexa <i class="pi pi-arrow-right" aria-hidden="true"></i>
            </router-link>
            <span class="access-note"><i class="pi pi-lock" aria-hidden="true"></i> Akses khusus staf rumah sakit</span>
          </div>
          <div class="hero-footnote">
            <span class="footnote-rule"></span>
            <span>Terintegrasi untuk setiap tahap perjalanan klaim</span>
          </div>
        </div>

        <div class="workflow-visual" aria-label="Ilustrasi alur kerja Casemix">
          <div class="visual-topline">
            <div>
              <span class="visual-kicker">ALUR KERJA</span>
              <h2>Perjalanan klaim</h2>
            </div>
            <span class="visual-icon"><i class="pi pi-sitemap" aria-hidden="true"></i></span>
          </div>
          <div class="workflow-list">
            <div class="workflow-item">
              <span class="workflow-index">01</span>
              <span class="workflow-icon"><i class="pi pi-id-card" aria-hidden="true"></i></span>
              <span class="workflow-copy"><strong>Registrasi & rekam medis</strong><small>Data pasien dan kunjungan</small></span>
              <i class="pi pi-check-circle workflow-check" aria-hidden="true"></i>
            </div>
            <div class="workflow-connector"></div>
            <div class="workflow-item">
              <span class="workflow-index">02</span>
              <span class="workflow-icon"><i class="pi pi-file-edit" aria-hidden="true"></i></span>
              <span class="workflow-copy"><strong>Diagnosis & tindakan</strong><small>Koding klinis terstruktur</small></span>
              <i class="pi pi-check-circle workflow-check" aria-hidden="true"></i>
            </div>
            <div class="workflow-connector"></div>
            <div class="workflow-item workflow-item-current">
              <span class="workflow-index">03</span>
              <span class="workflow-icon"><i class="pi pi-receipt" aria-hidden="true"></i></span>
              <span class="workflow-copy"><strong>Verifikasi Casemix</strong><small>Grouping dan validasi klaim</small></span>
              <span class="current-label"><span></span> TERHUBUNG</span>
            </div>
            <div class="workflow-connector"></div>
            <div class="workflow-item">
              <span class="workflow-index">04</span>
              <span class="workflow-icon"><i class="pi pi-chart-bar" aria-hidden="true"></i></span>
              <span class="workflow-copy"><strong>Monitoring & audit</strong><small>Transparansi status klaim</small></span>
              <i class="pi pi-arrow-right workflow-arrow" aria-hidden="true"></i>
            </div>
          </div>
          <div class="visual-footer">
            <span><i class="pi pi-lock" aria-hidden="true"></i> Akses berbasis peran</span>
            <span class="visual-footer-divider"></span>
            <span><i class="pi pi-history" aria-hidden="true"></i> Jejak aktivitas tercatat</span>
          </div>
        </div>
      </section>

      <section class="capabilities">
        <div class="section-heading">
          <div>
            <span class="section-kicker">DIBUAT UNTUK ALUR KERJA RUMAH SAKIT</span>
            <h2>Data lebih rapi. Koordinasi lebih jelas.</h2>
          </div>
          <p>Satu sistem untuk menjaga kesinambungan informasi dari pelayanan hingga administrasi klaim.</p>
        </div>
        <div class="capability-grid">
          <article class="capability-card">
            <span class="capability-icon"><i class="pi pi-search" aria-hidden="true"></i></span>
            <span class="capability-number">01</span>
            <h3>Pencarian klinis cepat</h3>
            <p>Temukan rekam medis serta referensi kode ICD dengan lebih mudah dalam alur kerja sehari-hari.</p>
          </article>
          <article class="capability-card">
            <span class="capability-icon"><i class="pi pi-shield" aria-hidden="true"></i></span>
            <span class="capability-number">02</span>
            <h3>Validasi koding</h3>
            <p>Validasi membantu menjaga diagnosis dan tindakan tetap konsisten pada kunjungan pasien.</p>
          </article>
          <article class="capability-card">
            <span class="capability-icon"><i class="pi pi-share-alt" aria-hidden="true"></i></span>
            <span class="capability-number">03</span>
            <h3>Alur klaim terhubung</h3>
            <p>Informasi klaim mengikuti proses klinis, dari pengelolaan rekam medis hingga verifikasi.</p>
          </article>
          <article class="capability-card">
            <span class="capability-icon"><i class="pi pi-history" aria-hidden="true"></i></span>
            <span class="capability-number">04</span>
            <h3>Akuntabilitas aktivitas</h3>
            <p>Perubahan penting tercatat untuk mendukung transparansi dan penelusuran aktivitas.</p>
          </article>
        </div>
      </section>

      <section class="landing-cta">
        <div>
          <span class="section-kicker">MULAI BEKERJA DENGAN NEXA</span>
          <h2>Siap melanjutkan alur kerja?</h2>
          <p>Masuk menggunakan akun yang diberikan oleh administrator rumah sakit.</p>
        </div>
        <router-link to="/login" class="button button-light">
          Masuk ke akun <i class="pi pi-arrow-right" aria-hidden="true"></i>
        </router-link>
      </section>
    </main>

    <footer class="landing-footer">
      <router-link to="/" class="brand-lockup footer-brand" aria-label="Nexa, beranda">
        <span class="brand-mark"><i class="pi pi-plus" aria-hidden="true"></i></span>
        <span class="brand-name">nexa<span>.</span></span>
      </router-link>
      <span>Nexa Casemix Management System <span class="footer-separator">/</span> © 2026</span>
      <span class="footer-secure"><i class="pi pi-lock" aria-hidden="true"></i> Akses aman & berbasis peran</span>
    </footer>
  </div>
</template>

<style scoped>
.landing-page {
  min-height: 100vh;
  background: var(--color-bg-base);
}

.landing-header,
.landing-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px max(5.5vw, 24px);
  border-bottom: 1px solid var(--color-border);
  background: var(--color-surface);
}

.brand-lockup {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--color-primary);
  text-decoration: none;
}

.brand-mark {
  display: grid;
  width: 32px;
  height: 32px;
  place-items: center;
  border-radius: 6px;
  background: var(--color-primary);
  color: white;
  font-size: 14px;
}

.brand-name {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.8px;
}

.brand-name > span {
  color: var(--color-status-pending);
}

.header-actions,
.hero-actions {
  display: flex;
  align-items: center;
  gap: 24px;
}

.header-caption,
.access-note {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-secondary);
  font-size: 12px;
}

.header-caption i,
.access-note i {
  color: var(--color-primary);
}

.button {
  display: inline-flex;
  min-height: 44px;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 17px;
  border: 1px solid transparent;
  border-radius: var(--radius-base);
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  transition: background-color var(--motion-fast), border-color var(--motion-fast), color var(--motion-fast), transform var(--motion-fast);
}

.button:hover {
  transform: translateY(-1px);
}

.button-primary {
  background: var(--color-primary);
  color: white;
}

.button-primary:hover {
  background: var(--color-primary-hover);
}

.button-secondary {
  border-color: var(--color-border);
  background: white;
  color: var(--color-primary);
}

.button-secondary:hover {
  border-color: var(--color-primary);
  background: var(--color-bg-base);
}

.button-light {
  background: white;
  color: var(--color-primary);
}

.button-light:hover {
  background: var(--color-bg-base);
}

.landing-hero {
  display: grid;
  max-width: 1440px;
  min-height: 590px;
  grid-template-columns: 0.92fr 1.08fr;
  align-items: center;
  gap: clamp(48px, 8vw, 120px);
  margin: 0 auto;
  padding: 68px max(8vw, 32px) 76px;
}

.hero-copy {
  animation: hero-in 380ms ease-out both;
}

.eyebrow,
.section-kicker,
.visual-kicker {
  color: var(--color-text-secondary);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1.25px;
}

.eyebrow {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 20px;
  color: var(--color-primary);
}

.eyebrow > span,
.current-label > span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--color-status-approved);
}

.hero-copy h1 {
  max-width: 610px;
  margin: 0;
  color: var(--color-primary-dark);
  font-size: clamp(38px, 4.5vw, 60px);
  font-weight: 600;
  letter-spacing: -2.6px;
  line-height: 1.1;
}

.hero-copy h1 span {
  display: block;
  color: var(--color-primary);
}

.hero-description {
  max-width: 495px;
  margin: 23px 0 28px;
  color: var(--color-text-secondary);
  font-size: 15px;
  line-height: 1.75;
}

.hero-actions {
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 17px;
}

.hero-footnote {
  display: flex;
  align-items: center;
  gap: 11px;
  margin-top: 47px;
  color: var(--color-text-secondary);
  font-size: 11px;
}

.footnote-rule {
  width: 27px;
  height: 1px;
  background: var(--color-primary);
}

.workflow-visual {
  position: relative;
  padding: 28px 30px 21px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-surface);
  box-shadow: 0 18px 50px rgb(31 41 55 / 7%);
  animation: panel-in 500ms 90ms ease-out both;
}

.visual-topline {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1px 20px;
  border-bottom: 1px solid var(--color-border);
}

.visual-kicker {
  color: var(--color-text-disabled);
  font-size: 9px;
}

.visual-topline h2 {
  margin: 4px 0 0;
  color: var(--color-primary-dark);
  font-size: 17px;
  font-weight: 600;
}

.visual-icon,
.workflow-icon,
.capability-icon {
  display: grid;
  flex: 0 0 auto;
  place-items: center;
  color: var(--color-primary);
  background: #EEF2F5;
}

.visual-icon {
  width: 38px;
  height: 38px;
  border-radius: 6px;
}

.workflow-list {
  padding-top: 10px;
}

.workflow-item {
  display: grid;
  min-height: 66px;
  grid-template-columns: 26px 36px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  border-radius: 6px;
  transition: background-color var(--motion-fast);
}

.workflow-item-current {
  border: 1px solid #D9E1E7;
  background: #F7F9FA;
}

.workflow-index {
  color: var(--color-text-disabled);
  font-family: var(--font-mono);
  font-size: 10px;
}

.workflow-icon {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  font-size: 14px;
}

.workflow-copy {
  display: grid;
  gap: 2px;
}

.workflow-copy strong {
  color: var(--color-text-primary);
  font-size: 12px;
  font-weight: 600;
}

.workflow-copy small {
  color: var(--color-text-secondary);
  font-size: 10px;
}

.workflow-check {
  color: var(--color-status-approved);
  font-size: 15px;
}

.workflow-arrow {
  color: var(--color-text-disabled);
  font-size: 12px;
}

.workflow-connector {
  height: 10px;
  margin-left: 27px;
  border-left: 1px dashed var(--color-border);
}

.current-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--color-primary);
  font-size: 9px;
  font-weight: 600;
  letter-spacing: 0.4px;
}

.current-label > span {
  width: 6px;
  height: 6px;
}

.visual-footer {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-top: 19px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: 10px;
}

.visual-footer span:not(.visual-footer-divider) {
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.visual-footer i {
  color: var(--color-primary);
  font-size: 11px;
}

.visual-footer-divider {
  width: 1px;
  height: 15px;
  background: var(--color-border);
}

.capabilities {
  padding: 63px max(8vw, 32px) 76px;
  border-top: 1px solid var(--color-border);
  background: var(--color-surface);
}

.section-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 30px;
  max-width: 1200px;
  margin: 0 auto 30px;
}

.section-heading h2,
.landing-cta h2 {
  margin: 8px 0 0;
  color: var(--color-primary-dark);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.6px;
}

.section-heading > p {
  max-width: 390px;
  margin: 0 0 2px;
  color: var(--color-text-secondary);
  font-size: 12px;
  line-height: 1.7;
}

.capability-grid {
  display: grid;
  max-width: 1200px;
  grid-template-columns: repeat(4, 1fr);
  margin: 0 auto;
  border: 1px solid var(--color-border);
  border-radius: 6px;
}

.capability-card {
  position: relative;
  min-height: 204px;
  padding: 24px 22px;
}

.capability-card + .capability-card {
  border-left: 1px solid var(--color-border);
}

.capability-icon {
  width: 36px;
  height: 36px;
  margin-bottom: 20px;
  border-radius: 6px;
}

.capability-number {
  position: absolute;
  top: 31px;
  right: 21px;
  color: var(--color-text-disabled);
  font-family: var(--font-mono);
  font-size: 10px;
}

.capability-card h3 {
  margin: 0 0 8px;
  color: var(--color-primary-dark);
  font-size: 14px;
  font-weight: 600;
}

.capability-card p {
  max-width: 250px;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 11px;
  line-height: 1.7;
}

.landing-cta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 36px max(8vw, 32px);
  background: var(--color-primary);
}

.landing-cta > div {
  max-width: 1200px;
}

.landing-cta .section-kicker {
  color: #CAD4DC;
}

.landing-cta h2 {
  margin-top: 6px;
  color: white;
  font-size: 21px;
}

.landing-cta p {
  margin: 6px 0 0;
  color: #E1E8ED;
  font-size: 12px;
}

.landing-footer {
  min-height: 68px;
  gap: 20px;
  border-top: 1px solid var(--color-border);
  border-bottom: 0;
  color: var(--color-text-secondary);
  font-size: 10px;
}

.footer-brand .brand-mark {
  width: 25px;
  height: 25px;
  font-size: 11px;
}

.footer-brand .brand-name {
  font-size: 16px;
}

.footer-secure {
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.footer-secure i {
  color: var(--color-primary);
}

.footer-separator {
  margin: 0 5px;
  color: var(--color-text-disabled);
}

@keyframes hero-in {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes panel-in {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 960px) {
  .landing-hero {
    min-height: auto;
    grid-template-columns: 1fr;
    gap: 42px;
    padding-top: 64px;
    padding-bottom: 64px;
  }

  .hero-copy h1 {
    max-width: 680px;
  }

  .workflow-visual {
    max-width: 650px;
  }

  .capability-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .capability-card:nth-child(3) {
    border-left: 0;
    border-top: 1px solid var(--color-border);
  }

  .capability-card:nth-child(4) {
    border-top: 1px solid var(--color-border);
  }
}

@media (max-width: 600px) {
  .landing-header {
    padding: 14px 20px;
  }

  .header-caption {
    display: none;
  }

  .landing-hero {
    padding: 48px 20px 52px;
  }

  .hero-copy h1 {
    font-size: 40px;
    letter-spacing: -1.8px;
  }

  .hero-description {
    font-size: 14px;
  }

  .hero-actions {
    align-items: flex-start;
    flex-direction: column;
  }

  .hero-footnote {
    margin-top: 34px;
  }

  .workflow-visual {
    padding: 21px 14px 17px;
  }

  .workflow-item {
    grid-template-columns: 20px 32px minmax(0, 1fr) auto;
    gap: 8px;
    padding-right: 5px;
    padding-left: 5px;
  }

  .workflow-icon {
    width: 32px;
    height: 32px;
  }

  .workflow-copy strong {
    font-size: 11px;
  }

  .workflow-copy small {
    font-size: 9px;
  }

  .current-label {
    font-size: 7px;
  }

  .visual-footer {
    gap: 9px;
    font-size: 9px;
  }

  .capabilities {
    padding: 45px 20px;
  }

  .section-heading {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
  }

  .section-heading h2 {
    font-size: 21px;
  }

  .capability-card {
    min-height: 195px;
    padding: 18px 14px;
  }

  .capability-card h3 {
    font-size: 12px;
  }

  .capability-card p {
    font-size: 10px;
  }

  .capability-number {
    top: 25px;
    right: 12px;
  }

  .landing-cta {
    align-items: flex-start;
    flex-direction: column;
    padding: 30px 20px;
  }

  .landing-footer {
    align-items: flex-start;
    flex-direction: column;
    gap: 10px;
    padding: 18px 20px;
  }
}
</style>
