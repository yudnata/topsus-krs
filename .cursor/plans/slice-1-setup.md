---
name: Slice 1 — Setup Project
overview: Scaffold Vue 3 frontend, rapikan Go backend, docker-compose PostgreSQL+Redis, perbaiki Dockerfile dan config.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: completed
  - id: scaffold-vue
    content: Folder frontend/ Vite+Vue3+TS+Pinia+Router+Axios
    status: completed
  - id: fix-dockerfile
    content: Dockerfile build ./cmd/main.go
    status: completed
  - id: extend-config
    content: config.go REDIS_URL + CORS_ORIGIN
    status: completed
  - id: docker-compose
    content: docker compose up -d postgres + redis
    status: completed
  - id: run-backend
    content: go run cmd/main.go tanpa error
    status: completed
  - id: run-frontend
    content: npm run dev frontend tanpa error
    status: completed
---

# Slice 1 — Setup Project (Vue + Go + Docker)

**Status:** Selesai (2026-05-19)

## Prompt Plan Mode

```text
Saya membangun Sistem KRS (Mahasiswa, Dosen, pemetaan DPA).
Stack: Vue 3 + Vite + TypeScript + Tailwind CSS v4 (frontend), Go Fiber v3 (backend), PostgreSQL, Redis.
Buatkan plan untuk Slice 1:
- Scaffold frontend/ dengan Vue 3 + Vite + TS + Pinia + Vue Router + Axios + Tailwind CSS v4
- Rapikan backend: perbaiki Dockerfile (build ./cmd/main.go), tambah REDIS_URL dan CORS_ORIGIN di config
- docker-compose up untuk PostgreSQL + Redis
- .env.example di root, backend/, frontend/
- Health endpoint /api/health sudah jalan
```

## Checklist

- [x] Plan Mode & plan tersimpan
- [x] `frontend/` Vite + Vue 3 + TS + Tailwind v4
- [x] Pinia, Vue Router, Axios
- [x] Scaffold feature-based `src/features/*` + `app/router`
- [x] `backend/Dockerfile` → `go build -o main ./cmd/main.go`
- [x] `config.go`: `REDIS_URL`, `CORS_ORIGIN`
- [x] `docker compose up -d` OK (postgres; redis opsional jika port 6379 bebas)
- [x] Backend & frontend build jalan

## File target

- `frontend/**`
- `backend/Dockerfile`, `backend/internal/config/config.go`
- `docker-compose.yml`, `.env.example`
