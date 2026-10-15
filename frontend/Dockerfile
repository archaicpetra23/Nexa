FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# node_modules dipasang sebagai named volume oleh docker-compose, sehingga
# menutupi node_modules hasil build image ini. Entrypoint menyinkronkan ulang
# dependency ketika package-lock.json berubah, supaya dependency baru
# (mis. chart.js) benar-benar terpasang tanpa perlu hapus volume.
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 5173

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
