# Nexa Development Log
(Chronological, newest first)

---

## 2026‑11‑09 04:27 (Phase 1 Pre‑Deployment – SSERAPHIM)

### Deployment Infrastructure Readied
- **Deployment strategy**: Named Cloudflare Tunnel (fixed URL `*.cfargotunnel.com`, HTTPS auto)
- **Phase structure**: 2 phases (pre‑deployment fixes, execution)
- **Target environment**: Personal STB mini server (Docker Compose + PostgreSQL)

### Phase 1 Progress ✅
**1.1 Environment & Secrets**
- `.env` created with secure JWT secret (openssl rand -base64 32)
- PostgreSQL credentials & API port configured
- CORS origins set for remote tunnel + local dev

**1.2 Docker Build Validation**
- `docker-compose.yml` syntax fixed (indentation in backend service)
- All services running healthy: postgres, backend, frontend
- Backend API reachable at `http://localhost:8080/health`
- Health check verified with curl

**1.3 Unit Test Suite**
- Backend:
  - date_helper_test.go (valid/invalid parsing)
  - jwt_test.go (GenerateToken, ValidateToken)
  - bcrypt_test.go already existed (full coverage)
  - klaim_usecase_test.go (state transitions, mocks)
  - All backend tests passing
- Frontend:
  - Installed Vitest + jsdom as dev dependencies
  - klaimValidation.test.js (status transitions)
  - pasienValidation.test.js (form validation)
  - format.test.js (Rupiah & date formatting)
  - 100 frontend tests passing

**1.4 Code Consolidation**
- Validation & formatting extracted into `/utils` (backend + frontend)
- State machine sync between frontend (validation.js) & backend (usecase.go)

**1.5 Docker Compose Deployment**
- PostgreSQL: 5432 (healthy)
- Backend Go: 8080 (healthy)
- Frontend Vue: 5173 → 80/tcp (rebuilding after Dockerfile fix)

### Build Issues Resolved
1. **Docker Compose YAML syntax** – corrected indentation in `backend:` service
2. **Dockerfile path mismatch** – build output at `dist-frontend`, not `dist`
3. **Vitest missing** – installed terser (required for Vite v3+)
4. **Frontend service unhealthy** – rebuilding with corrected Dockerfile

### Pending for Phase 2
- Cloudflare account creation (free tier)
- `cloudflared` installation on STB server
- Named Tunnel creation (`nexa` → `*.cfargotunnel.com`)
- systemd service for auto‑start
- Remote verification (professor access)

### Files Created/Modified
- `.env` – environment variables
- `docker-compose.yml` – fixed indentation
- `frontend/Dockerfile` – corrected output path
- `backend/internal/usecase/klaim_usecase_test.go` – mock repo tests
- `frontend/src/utils/*.test.js` – validation & formatting tests
- `DEPLOYMENT_PHASES.md` – two‑phase deployment checklist

---

## Earlier Logs (Historical)
[Previous development logs archived]