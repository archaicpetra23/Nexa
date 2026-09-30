FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 5173

# node_modules is bind-mounted as an anonymous/named volume by docker-compose,
# which masks the image's node_modules. Re-sync on every start so newly added
# deps in package.json (e.g. primeicons) are actually present.
CMD ["sh", "-c", "npm install && npm run dev -- --host 0.0.0.0"]
