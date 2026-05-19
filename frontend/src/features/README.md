# Features (feature-based / vertical slice)

Setiap folder = satu domain bisnis. Isi standar per feature:

```text
features/<nama>/
├── api/          # Panggilan HTTP ke backend
├── stores/       # Pinia (state feature)
├── views/        # Halaman Vue
├── components/   # Komponen khusus feature (opsional)
└── routes.ts     # Route feature — didaftarkan di app/router
```

| Feature | Route (rencana) | Role |
|---------|-----------------|------|
| `auth` | `/login` | Semua |
| `home` | `/` | Publik |
| `mahasiswa` | `/mahasiswa/krs` | MAHASISWA |
| `dosen` | `/dosen/persetujuan` | DOSEN |
| `staff` | `/staff/persetujuan` | STAFF |
| `admin` | `/admin/users`, `/admin/dpa` | ADMIN |

Slice 5 mengisi UI lengkap; placeholder views ada untuk guard & routing.
