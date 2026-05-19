# Frontend — Sistem KRS

Vue 3 + Vite + TypeScript + Pinia + Vue Router + Axios + Tailwind CSS v4.

## Arsitektur feature-based

```text
src/
├── app/router/          # Gabung semua routes + role guards
├── core/api/            # Axios client (shared)
├── core/types/          # Tipe API bersama
└── features/            # Satu folder per domain bisnis
    ├── auth/            # api, stores, views, routes.ts
    ├── home/
    ├── mahasiswa/
    ├── dosen/
    ├── staff/
    └── admin/
```

Detail: [ARSITEKTUR.md §8](../ARSITEKTUR.md) dan [src/features/README.md](src/features/README.md).

## Development

```bash
cp .env.example .env
npm install
npm run dev
```

`VITE_API_URL` → backend (default `http://localhost:8080`).

## Build

```bash
npm run build
```

## Route (scaffold)

| Path | Feature |
|------|---------|
| `/` | home |
| `/login` | auth |
| `/mahasiswa/krs` | mahasiswa (Slice 5) |
| `/dosen/persetujuan` | dosen |
| `/staff/persetujuan` | staff |
| `/admin/users` | admin |
