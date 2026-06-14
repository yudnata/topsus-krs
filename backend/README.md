# Backend API (Express.js)

REST API untuk Sistem KRS dengan arsitektur **feature-based** (`controller` → `service` → `repository`).

## Struktur

```text
backend/
├── src/
│   ├── features/
│   │   ├── auth/
│   │   ├── admin/
│   │   ├── mahasiswa/
│   │   ├── dosen/
│   │   └── cachemeta/
│   ├── shared/
│   │   ├── cache/
│   │   ├── database/
│   │   └── utils/
│   ├── config/
│   ├── routes/
│   ├── app.ts
│   └── index.ts
├── migrations/
└── package.json
```

## Menjalankan lokal

```bash
cp .env.example .env
npm install
npm run dev
```

Health check: `GET http://localhost:8080/api/health`

## Build production

```bash
npm run build
npm start
```

## Endpoint utama

- `POST /api/auth/login` — JWT + profil user
- `GET /api/admin/*` — master data (ADMIN)
- `GET|POST /api/krs/*` — alur KRS mahasiswa
- `GET|POST /api/approval/*` — persetujuan dosen/staff
