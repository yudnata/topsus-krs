---
name: Slice 7 — Deploy PaaS
overview: Vercel frontend, Azure backend, DB/Redis cloud, GitHub Actions, verifikasi production.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: pending
  - id: github-desktop
    content: Upload repo via GitHub Desktop (GITHUB.md §2.23)
    status: pending
  - id: db-redis-cloud
    content: PostgreSQL Neon Console + Redis Upstash/Railway
    status: pending
  - id: azure-backend
    content: Azure App Service + env + Actions deploy
    status: pending
  - id: vercel-frontend
    content: Vercel import GitHub + VITE_API_URL
    status: pending
  - id: production-test
    content: Login + CRUD + DPA di production
    status: pending
---

# Slice 7 — Deploy PaaS + CI/CD

## Prompt Plan Mode

```text
Buatkan plan deploy Sistem KRS ke production.
- §2.23: GitHub Desktop — Publish/Push ke main (GITHUB.md)
- §2.24: Vercel frontend + Azure backend (DEPLOY.md)
- Neon + Redis env di Azure; VITE_API_URL di Vercel
- Verifikasi production: login, CRUD, DPA
```

## Checklist

- [ ] GitHub Desktop push selesai (§2.23)
- [ ] DB + Redis cloud
- [ ] Azure backend live
- [ ] Vercel frontend live
- [ ] Env production lengkap
- [ ] Production E2E OK

## Referensi

- [GITHUB.md](../../GITHUB.md) — §2.23 GitHub Desktop
- [DEPLOY.md](../../DEPLOY.md) — §2.24 deploy PaaS
