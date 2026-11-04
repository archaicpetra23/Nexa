# Nexa Deployment Readiness Checklist

## Completed Tasks ✅

### TASK 1: Complete Backend Admin Panel
- [x] Added 3 missing Role CRUD endpoints (POST/PUT/DELETE /master/roles)
- [x] Updated routes.go with all endpoints
- [x] Created Role handlers in master_handler.go
- [x] Updated Role model with Deskripsi field
- [x] Created migration 003_deskripsi_role.sql

### TASK 2: Security Hardening
- [x] Rate limiting middleware (100 req/min per IP) - `rate_limit.go`
- [x] Security headers middleware (CSP, X-Frame-Options, etc.) - `security_headers.go`
- [x] CORS configurable via env (`CORS_ALLOWED_ORIGINS`)
- [x] Cookie Secure & SameSite configurable via env (`COOKIE_SECURE`)
- [x] Config updated with `CORSAllowedOrigins` and `RateLimitEnabled`
- [x] All middleware integrated in routes.go

### TASK 3: Input Validation
- [x] Created sanitization helper - `pkg/validation/sanitize.go`

### TASK 4: Error Handling
- [x] Created global error handler - `middleware/error_handler.go`

### TASK 5: Frontend Performance
- [x] vite.config.js with code-splitting and compression
- [x] Router with dynamic imports for all views
- [x] Production Dockerfile with nginx
- [x] nginx.conf with security headers and API proxy
- [x] package.json with build:prod and analyze scripts

### TASK 6: Docker Production
- [x] Backend multi-stage binary build Dockerfile
- [x] Healthchecks for all services
- [x] Updated docker-compose.yml for production
- [x] Comprehensive .env.example

### TASK 7: Analytics Export
- [x] Analytics handler with CSV/JSON export
- [x] Route with RBAC (admin_ti, petugas_casemix)

## Remaining Tasks ⚠️

### Code Quality
- [ ] Apply sanitization to all handlers (use `validation.SanitizeString()`)
- [ ] Replace `err.Error()` with generic messages in all handlers
- [ ] Remove `fmt.Println` from migrate/main.go
- [ ] Standardize JSON tag casing (snake_case)
- [ ] Extract duplicate validation logic to pkg/validation

### Database
- [ ] Apply migration 003_deskripsi_role.sql
- [ ] Apply 002_performance_indexes.sql
- [ ] Verify trigram index `idx_trgm_pasien_nama`

### Frontend
- [ ] Install `vite-plugin-compression` dependency
- [ ] Run `npm install` in frontend
- [ ] Test production build with `npm run build:prod`

### Docker & Deployment
- [ ] Create .env file from .env.example with production values
- [ ] Build and test: `docker compose up -d --build`
- [ ] Verify all services healthy
- [ ] Test healthcheck endpoints

### Testing
- [ ] Run `go build ./...` and `go vet` clean
- [ ] Run `npm run lint` in frontend
- [ ] Test all 12 admin panel endpoints
- [ ] Test analytics export CSV download
- [ ] Verify RBAC works correctly

## Quick Verification Commands

```bash
# 1. Apply database migrations
docker compose exec -T postgres psql -U postgres -d nexa < backend/scripts/003_deskripsi_role.sql
docker compose exec -T postgres psql -U postgres -d nexa < backend/scripts/002_performance_indexes.sql

# 2. Build and start
docker compose up -d --build

# 3. Check health
curl http://localhost:8080/health
curl http://localhost:5173/

# 4. Verify services
docker compose ps

# 5. Test admin endpoints
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <token>"

# 6. Test analytics export
curl "http://localhost:8080/api/v1/analytics/export-klaim?format=csv" \
  -H "Authorization: Bearer <token>"
```

## Production Environment Variables

```env
DB_DSN=postgres://user:pass@postgres:5432/nexa
POSTGRES_USER=nexa
POSTGRES_PASSWORD=secure_password
POSTGRES_DB=nexa
JWT_SECRET=32_character_minimum_random_string
SERVER_PORT=:8080
COOKIE_DOMAIN=yourdomain.com
COOKIE_SECURE=true
CORS_ALLOWED_ORIGINS=https://yourdomain.com
RATE_LIMIT_ENABLED=true
VITE_API_PROXY=https://yourdomain.com/api
```