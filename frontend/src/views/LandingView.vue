<template>
  <div class="min-h-screen bg-base">
    <header class="bg-primary text-white py-4 px-8 shadow-sm">
      <div class="max-w-screen-xl mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-semibold">Nexa</h1>
        <router-link to="/login" class="px-6 py-2 bg-white text-primary rounded hover:bg-gray-100 transition-colors">
          Login
        </router-link>
      </div>
    </header>

    <main class="max-w-screen-xl mx-auto px-8 py-16">
      <section class="text-center mb-20">
        <h2 class="text-5xl font-semibold text-primary mb-6">
          Sistem Informasi Pengelolaan Data Casemix Terintegrasi
        </h2>
        <p class="text-xl text-secondary max-w-3xl mx-auto leading-relaxed">
          Mempercepat siklus operasional klaim BPJS Kesehatan di Rumah Sakit Umum Daerah (RSUD) dengan integrasi end-to-end, pencarian berkas instan, dan transparansi penuh.
        </p>
      </section>

      <section class="grid md:grid-cols-2 gap-8 mb-16">
        <div class="bg-surface p-8 rounded border border-border">
          <h3 class="text-xl font-semibold text-primary mb-4">Pencarian Cepat</h3>
          <p class="text-secondary leading-relaxed">
            Temu-balik rekam medis dan master kode ICD dalam hitungan detik, menurunkan latensi dari 15-60 menit menjadi kurang dari 5 detik.
          </p>
        </div>

        <div class="bg-surface p-8 rounded border border-border">
          <h3 class="text-xl font-semibold text-primary mb-4">Zero Duplicate Coding</h3>
          <p class="text-secondary leading-relaxed">
            Mencegah insiden duplikasi kode diagnosis atau tindakan pada nomor kunjungan yang sama melalui constraint database.
          </p>
        </div>

        <div class="bg-surface p-8 rounded border border-border">
          <h3 class="text-xl font-semibold text-primary mb-4">Otomasi Rekomendasi Tarif INA-CBGs</h3>
          <p class="text-secondary leading-relaxed">
            Memberikan pemetaan paket tarif klaim otomatis begitu diagnosis primer dan tindakan selesai diinput, memangkas waktu koding minimal 50%.
          </p>
        </div>

        <div class="bg-surface p-8 rounded border border-border">
          <h3 class="text-xl font-semibold text-primary mb-4">Audit Trail & Keamanan Ketat</h3>
          <p class="text-secondary leading-relaxed">
            Mengamankan catatan rekam medis dari modifikasi tidak sah dan mencatat seluruh mutasi data secara otomatis sesuai OWASP Top 10.
          </p>
        </div>
      </section>

      <section class="text-center">
        <router-link to="/login" class="inline-block px-10 py-4 bg-primary text-white text-lg font-medium rounded hover:bg-primary-hover transition-colors">
          Mulai Sekarang
        </router-link>
      </section>
    </main>

    <footer class="bg-surface border-t border-border py-6 mt-16">
      <div class="max-w-screen-xl mx-auto px-8 text-center text-secondary text-sm">
        Nexa Casemix Management System &copy; 2026
      </div>
    </footer>
  </div>
</template>

<script setup>
</script>

<style scoped>
</style>
