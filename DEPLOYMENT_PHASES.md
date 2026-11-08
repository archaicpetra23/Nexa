# Deployment Phases — Nexa

## PHASE 1: PRE‑DEPLOYMENT (Fix & Prepare Before Deploy)

### 1.1 Environment & Secrets
- [ ] Create valid `.env` from `.env.example`
  - `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
  - `JWT_SECRET` (32+ bytes: `openssl rand -base64 32`)
  - `COOKIE_SECURE=true` (for HTTPS)
  - `COOKIE_DOMAIN` = remote hostname (e.g., `nexa.cfargotunnel.com`)
  - `CORS_ALLOWED_ORIGINS` = `https://nexa.cfargotunnel.com`
  - `RATE_LIMIT_ENABLED=true`
  - `SERVER_PORT=:8080`
- [ ] Verify `.env` not committed (check `.gitignore`)

### 1.2 Docker Images & Build Validation
- [ ] `docker compose build` — all images build without error
- [ ] `docker compose up -d` → all services **Up (healthy)**
- [ ] Health check: `curl http://localhost:8080/health` → `{"message":"Nexa API running","success":true}`
- [ ] Frontend loads at `http://localhost:5173` (or nginx port)

### 1.3 Unit Test Suite (Mandatory per brief)
- [ ] **Backend Go tests** — create files:
  - `backend/pkg/jwt_test.go`
  - `backend/pkg/bcrypt_test.go`
  - `backend/internal/delivery/http/date_helper_test.go`
  - `backend/internal/usecase/klaim_status_test.go`
  - `backend/internal/repository/mock_pasien_repository_test.go`
  - `backend/internal/repository/mock_klaim_repository_test.go`
  - `backend/internal/repository/mock_rekam_medis_repository_test.go`
  - `backend/internal/usecase/klaim_usecase_test.go`
  - `backend/internal/usecase/pasien_usecase_test.go`
  - `backend/internal/usecase/rekam_medis_usecase_test.go`
- [ ] Run: `docker compose exec backend go test ./... -cover` → all `ok` with coverage
- [ ] **Frontend Vitest** — add deps & config:
  - `vitest`, `@vue/test-utils`, `jsdom` in `frontend/package.json`
  - `frontend/vitest.config.js`
  - Test files:
    - `src/views/__tests__/pasien-form.test.js`
    - `src/views/__tests__/klaim-status.test.js`
    - `src/api/__tests__/client.test.js`
    - `src/utils/__tests__/format.test.js`
- [ ] Run: `cd frontend && npm run test` → all pass
- [ ] Run: `cd frontend && npm run build` → succeeds

### 1.4 Code Consolidation (Behaviour-preserving)
- [ ] Extract `validatePasienForm` → `frontend/src/utils/pasienValidation.js`, update `PasienView.vue`
- [ ] Create `frontend/src/utils/klaimValidation.js` (transition rules), wire into `KlaimView.vue` modal
- [ ] Consolidate 4× `formatRupiah` + `formatTanggalWaktu` → `frontend/src/utils/format.js`, update 5 Vue files
- [ ] Verify `npm run build` output size unchanged (±1%)

### 1.5 Documentation Updates
- [ ] Add **"Menjalankan Unit Test"** section in README (after "Verifikasi Semua Service Jalan")
- [ ] Add DevLog entry in `docs/DevLog.md` (newest-first, exact template)
  - Agent: SSERAPHIM
  - Tipe: Test & Quality
  - Status: Selesai
  - Modul: Backend (unit test), Frontend (Vitest), Dokumentasi
  - File terdampak: [list all created/modified]
  - Deskripsi: before→after counts, what added
  - Verifikasi: literal commands + last output lines
  - Error/Kendala: – atau detail error

### 1.6 Cloudflare Tunnel Setup (Named Tunnel = Fixed URL)
- [ ] Create Named Tunnel in Cloudflare Dashboard → get Tunnel ID + credentials JSON
- [ ] Save credentials to `~/.cloudflared/nexa-tunnel.json`
- [ ] Create `~/.cloudflared/config.yml`:
  ```yaml
  credentials-file: /home/razan/.cloudflared/nexa-tunnel.json
  ingress:
    - hostname: nexa.cfargotunnel.com
      service: http://localhost:8080
  ```
- [ ] Test tunnel: `cloudflared tunnel --config ~/.cloudflared/config.yml run nexa-tunnel`
- [ ] Verify accessible at `https://nexa.cfargotunnel.com`
- [ ] Create systemd service for auto-start:
  ```ini
  # /etc/systemd/system/cloudflared-nexa.service
  [Unit]
  Description=Cloudflare Tunnel for Nexa
  After=network.target docker.service
  Requires=docker.service

  [Service]
  Type=simple
  User=razan
  ExecStart=/usr/bin/cloudflared tunnel --config /home/razan/.cloudflared/config.yml run nexa-tunnel
  Restart=always
  RestartSec=5

  [Install]
  WantedBy=multi-user.target
  ```
- [ ] Enable: `sudo systemctl enable --now cloudflared-nexa`

### 1.7 Backup Strategy
- [ ] Create backup script `/home/razan/scripts/backup-nexa.sh`:
  ```bash
  #!/bin/bash
  DATE=$(date +%F)
  docker exec nexa-postgres-1 pg_dump -U nexa nexa > /home/razan/backups/nexa_${DATE}.sql
  find /home/razan/backups -name "nexa_*.sql" -mtime +7 -delete
  ```
- [ ] Make executable + add cron: `0 2 * * * /home/razan/scripts/backup-nexa.sh`

---

## PHASE 2: DEPLOYMENT (Run on STB Mini Server)

### 2.1 Server Preparation
- [ ] Install Docker & Docker Compose v2
- [ ] Install `cloudflared` (Arch: `pacman -S cloudflared` or download binary)
- [ ] Ensure ports 8080, 5432 available (or adjust in docker-compose.yml)

### 2.2 Deploy Application
- [ ] Clone repo: `git clone https://github.com/archaicpetra23/Nexa.git ~/nexa`
- [ ] Copy pre-prepared `.env` to `~/nexa/.env`
- [ ] Build & start:
  ```bash
  cd ~/nexa
  docker compose up --build -d
  ```
- [ ] Wait 30s, verify:
  ```bash
  docker compose ps
  curl http://localhost:8080/health
  ```

### 2.3 Start Cloudflare Tunnel
- [ ] Copy tunnel credentials & config to STB:
  - `~/.cloudflared/nexa-tunnel.json`
  - `~/.cloudflared/config.yml`
- [ ] Start/enable systemd service:
  ```bash
  sudo systemctl enable --now cloudflared-nexa
  ```
- [ ] Verify tunnel status: `systemctl status cloudflared-nexa`
- [ ] Test remote access: open `https://nexa.cfargotunnel.com` from external device

### 2.4 Post‑Deploy Verification
- [ ] Login with demo accounts (admin / password123)
- [ ] Test each role RBAC (7 roles)
- [ ] Create pasien, rekam medis, klaim → verify flow
- [ ] Check audit trail logs
- [ ] Run unit tests again on server:
  ```bash
  docker compose exec backend go test ./... -cover
  cd frontend && npm run test
  ```

### 2.5 Monitoring & Maintenance
- [ ] Verify backup cron runs (check `/home/razan/backups/`)
- [ ] Monitor logs: `docker compose logs -f backend`
- [ ] Set up log rotation (optional): `docker compose logs --tail=1000 > nexa_$(date +%F).log`

---

## QUICK REFERENCE

| Phase | Command | Expected Result |
|-------|---------|-----------------|
| Pre‑1 | `docker compose build` | All images built |
| Pre‑2 | `docker compose up -d` | All Up (healthy) |
| Pre‑3 | `go test ./... -cover` | All `ok` + coverage |
| Pre‑4 | `npm run test` | All Vitest pass |
| Pre‑5 | `cloudflared tunnel run` | Tunnel connected |
| Deploy‑1 | `docker compose up --build -d` | Services running |
| Deploy‑2 | `systemctl start cloudflared-nexa` | Tunnel active |
| Deploy‑3 | Open `https://nexa.cfargotunnel.com` | App loads, login works |

---

## ROLLBACK PLAN
- **App only**: `docker compose down && docker compose up -d` (keeps DB)
- **Full reset**: `docker compose down -v && docker compose up --build -d` (re-seeds DB)
- **Tunnel issues**: `systemctl restart cloudflared-nexa` or re-create Named Tunnel in dashboard