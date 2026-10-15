#!/bin/sh
#
# Entrypoint frontend Nexa.
#
# Masalah yang dicegah:
# docker-compose memasang named volume di atas /app/node_modules. Volume itu
# menutupi node_modules hasil build image dan isinya bertahan lintas restart.
# Akibatnya dependency yang ditambahkan SETELAH container start tidak pernah
# terpasang, dan Vite gagal dengan error seperti:
#   Failed to resolve import "chart.js" from "src/views/DashboardView.vue"
#
# Strategi: jalankan npm install hanya ketika manifest dependency berubah.
# Restart biasa tetap cepat, sedangkan perubahan dependency langsung tersinkron
# tanpa perlu menghapus volume.

set -e

STAMP=/app/node_modules/.deps-hash

# Hash package.json DAN package-lock.json. Memantau lockfile saja tidak cukup:
# kalau package.json diedit tanpa regenerasi lock (persis kasus chart.js),
# hash lock tidak berubah dan dependency baru akan terlewat.
current=$(cat /app/package.json /app/package-lock.json 2>/dev/null | md5sum | cut -d' ' -f1)

stored=""
if [ -f "$STAMP" ]; then
  stored=$(cat "$STAMP" 2>/dev/null || true)
fi

need_install=0
if [ ! -d /app/node_modules ]; then
  need_install=1
elif [ "$current" != "$stored" ]; then
  need_install=1
fi

if [ "$need_install" = "1" ]; then
  echo "[entrypoint] manifest dependency berubah (sekarang=$current tersimpan=${stored:-kosong}) - menjalankan npm install"
  npm install
  # Hitung ulang setelah install: npm menulis ulang lockfile, jadi hash
  # sebelum dan sesudah install bisa berbeda.
  cat /app/package.json /app/package-lock.json 2>/dev/null | md5sum | cut -d' ' -f1 > "$STAMP"
  echo "[entrypoint] dependency sudah sinkron (hash=$current)"
else
  echo "[entrypoint] dependency sudah sinkron (hash=$current)"
fi

exec "$@"
