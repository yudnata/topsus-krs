---
name: Sistem KRS — Master Plan
overview: Rencana end-to-end Sistem KRS (Vue 3 + Express.js + PostgreSQL + Redis) dengan pemetaan DPA, deploy Vercel + Azure.
todos:
  - id: slice-0-docs
    content: Dokumentasi ARSITEKTUR, CHECKPOINT, CONTEXT, PROMPT_LOG, DEPLOY
    status: completed
  - id: slice-1-setup
    content: Setup Vue + Express + koneksi Cloud Database + config
    status: pending
  - id: slice-2-database
    content: Skema 11 Tabel Master Data & KRS + migrasi + seed
    status: pending
  - id: slice-3-redis
    content: Redis cache layer + invalidation
    status: pending
  - id: slice-4-backend
    content: API KRS Flow + Master Data + Validasi Bisnis + auth JWT
    status: pending
  - id: slice-5-frontend
    content: Vue CRUD + pemetaan DPA UI
    status: pending
  - id: slice-6-integration
    content: E2E testing + CORS
    status: pending
  - id: slice-7-deploy
    content: Deploy Vercel + Azure + CI/CD
    status: pending
---

# Master Plan — Sistem KRS

## Tujuan

Aplikasi CRUD **Mahasiswa** dan **Dosen** dengan relasi **Dosen Pembimbing Akademik** (FK), cache Redis, auth JWT, deploy PaaS tanpa terminal.

## Stack

- Frontend: Vue 3 + Vite + TypeScript → **Vercel**
- Backend: Express.js → **Azure App Service**
- DB: PostgreSQL (Neon) | Cache: Redis (Azure Cache atau Upstash)

## Urutan slice

| # | Slice | Plan file |
|---|-------|-----------|
| 0 | Dokumentasi | (selesai) |
| 1 | Setup | [slice-1-setup.md](slice-1-setup.md) |
| 2 | Database | [slice-2-database.md](slice-2-database.md) |
| 3 | Redis | [slice-3-redis.md](slice-3-redis.md) |
| 4 | Backend API | [slice-4-backend-crud.md](slice-4-backend-crud.md) |
| 5 | Frontend | [slice-5-frontend.md](slice-5-frontend.md) |
| 6 | Integrasi | [slice-6-integration.md](slice-6-integration.md) |
| 7 | Deploy | [slice-7-deploy.md](slice-7-deploy.md) |

## Dokumen kontrol

- [CHECKPOINT.md](../../CHECKPOINT.md) — progress & checklist
- [CONTEXT.md](../../CONTEXT.md) — state aktif
- [ARSITEKTUR.md](../../ARSITEKTUR.md) — desain teknis
- [DEPLOY.md](../../DEPLOY.md) — hosting PaaS

## Mulai dari sini

Slice aktif: **4** — buka [slice-4-backend-crud.md](slice-4-backend-crud.md), aktifkan Plan Mode, eksekusi checklist.
