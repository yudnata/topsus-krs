---
name: Slice 4 — Backend CRUD + DPA
overview: Feature admin (master data), mahasiswa (KRS flow + 5 validasi bisnis), dosen (approval DPA), JWT RBAC, cache di service.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: done
  - id: module-admin
    content: Feature admin controller/service/repository/routes (CRUD MK, Kelas, Dosen, Mahasiswa, DPA)
    status: done
  - id: module-mahasiswa-krs
    content: Feature mahasiswa KRS flow (add-class, remove-class, submit, history) + 5 validasi bisnis
    status: done
  - id: module-dosen-approval
    content: Feature dosen approval (pending DPA, approve, reject) + otorisasi DPA
    status: done
  - id: auth-roles
    content: ADMIN / MAHASISWA / DOSEN / STAFF RBAC via requireRole middleware
    status: done
  - id: cache-integration
    content: Cache invalidation di setiap service layer
    status: done
  - id: seed-users
    content: Seed user accounts (password123) + link user_id ke dosen & mahasiswa
    status: done
  - id: build-ok
    content: npm run build backend sukses tanpa error
    status: done
---

# Slice 4 — Backend CRUD + DPA + Auth

## Prompt Plan Mode

```text
Buatkan plan untuk API backend Sistem KRS beserta validasinya.
- Modul master: CRUD mata kuliah, kelas (tersedia sesuai tahun akademik aktif).
- Modul krs: POST /api/krs/:id/add-class dengan validasi bisnis:
  1. Kapasitas kelas tidak penuh
  2. Bentrok jadwal (waktu tumpang tindih)
  3. Max SKS per semester (limit dari profil mahasiswa)
  4. Duplikat kelas/MK dicegah
- Modul approval: DOSEN hanya melihat mhs bimbingan; mengunci KRS jika status disetujui.
- Integrasi cache Redis di service layer
- JWT middleware dengan Role (ADMIN, MAHASISWA, DOSEN)
```

## Checklist

- [ ] API Master Data (kelas, MK, prodi) lengkap
- [ ] API KRS Flow (add-class, remove-class, submit) berfungsi
- [ ] Validasi bentrok jadwal & max SKS teruji
- [ ] Validasi kapasitas kelas berjalan (kapasitas penuh ditolak)
- [ ] Approval dosen mengunci draft mahasiswa
- [ ] Auth JWT + RBAC beroperasi sempurna

## File target

- `backend/src/features/dosen/`
- `backend/src/features/mahasiswa/`
- `backend/src/routes/index.ts`
