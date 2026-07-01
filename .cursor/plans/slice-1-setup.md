---
name: Slice 1 — Setup Project
overview: Scaffold Vue 3 frontend, rapikan Express.js backend, hubungkan ke PostgreSQL+Redis Cloud.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: completed
  - id: scaffold-vue
    content: Folder frontend/ Vite+Vue3+TS+Pinia+Router+Axios
    status: completed
  - id: fix-config
    content: Hubungkan backend ke Cloud Neon Database
    status: completed
  - id: extend-config
    content: config/index.ts REDIS_URL + CORS_ORIGIN
    status: completed
  - id: connect-cloud
    content: Verifikasi koneksi postgres + redis cloud
    status: completed
  - id: run-backend
    content: npm run dev backend tanpa error
    status: completed
  - id: run-frontend
    content: npm run dev frontend tanpa error
    status: completed
---

# Slice 1 — Setup Project (Vue + Express)

**Status:** Selesai (2026-05-19)

## Prompt Plan Mode

```text
Saya membangun Sistem KRS (Mahasiswa, Dosen, pemetaan DPA).
Stack: Vue 3 + Vite + TypeScript + Tailwind CSS v4 (frontend), Express.js (backend), PostgreSQL, Redis.
Buatkan plan untuk Slice 1:
- Scaffold frontend/ dengan Vue 3 + Vite + TS + Pinia + Vue Router + Axios + Tailwind CSS v4
- Rapikan backend: hubungkan ke Cloud Neon, tambah REDIS_URL dan CORS_ORIGIN di config
- .env.example di root, backend/, frontend/
- Health endpoint /api/health sudah jalan
```

## Checklist

- [x] Plan Mode & plan tersimpan
- [x] `frontend/` Vite + Vue 3 + TS + Tailwind v4
- [x] Pinia, Vue Router, Axios
- [x] Scaffold feature-based `src/features/*` + `app/router`
- [x] Backend terhubung dengan Database Cloud Neon
- [x] `config/index.ts`: `REDIS_URL`, `CORS_ORIGIN`
- [x] Backend & frontend build jalan

## File target

- `frontend/**`
- `backend/src/config/index.ts`
- `.env.example`
