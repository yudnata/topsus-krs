---
name: Slice 2 — Database
overview: Skema 11 tabel KRS, migrasi SQL, seed, update migrate().
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: completed
  - id: migration-sql
    content: migrations/ SQL 11 tabel lengkap
    status: completed
  - id: update-migrate
    content: migrate() baca file SQL berurutan
    status: completed
  - id: seed-dosen
    content: Seed minimal master + jadwal
    status: completed
  - id: verify-fk
    content: Query FK mahasiswa → dosen valid
    status: completed
---

# Slice 2 — Database (11 Tabel KRS)

**Status:** Selesai (2026-05-19)

## File migrasi

| File | Isi |
| ---- | --- |
| `001_extensions.sql` | `pgcrypto` untuk `gen_random_uuid()` |
| `002_krs_schema.sql` | 11 tabel + index |
| `003_seed.sql` | 2 prodi, 1 TA, 3 dosen, 2 mhs, 5 MK, 3 kelas, jadwal |

## Checklist

- [x] migrations/ SQL 11 tabel lengkap
- [x] migrate() mengeksekusi file SQL
- [x] Seed master data berjalan
- [x] Query FK KRS dan Jadwal Kelas valid
