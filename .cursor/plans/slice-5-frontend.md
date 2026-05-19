---
name: Slice 5 — Frontend Vue KRS
overview: UI Role-based untuk KRS Flow, DPA, Master Data, JWT Pinia, Axios interceptor.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: done
  - id: auth-ui
    content: Login page premium + redirect per role + Pinia auth store
    status: done
  - id: app-layout
    content: AppLayout sidebar role-aware dengan logout
    status: done
  - id: ui-mahasiswa
    content: KRS Form (add/remove kelas, SKS bar, submit) + Riwayat KRS
    status: done
  - id: ui-dosen
    content: Persetujuan KRS (approve + reject modal dengan catatan)
    status: done
  - id: ui-staff
    content: Antrian semua KRS (reuse PersetujuanView, role-aware)
    status: done
  - id: ui-admin
    content: Dashboard + Mahasiswa + Dosen + MK + Kelas (CRUD + DPA modal)
    status: done
  - id: axios-interceptor
    content: Bearer token interceptor dari localStorage
    status: done
  - id: design-system
    content: Dark mode CSS system (tokens, btn, card, table, badge, modal, animations)
    status: done
  - id: build-ok
    content: npm run build sukses tanpa error
    status: done
---

# Slice 5 — Frontend Vue KRS Flow & Master Data

## Prompt Plan Mode

```text
Buatkan plan untuk frontend Vue Sistem KRS (Role-Based).
- Halaman login/register (JWT di Pinia + localStorage)
- Halaman Mahasiswa: Form pengisian KRS, tambah kelas (validasi visual kapasitas/SKS), submit KRS.
- Halaman Dosen: Tabel antrian persetujuan KRS mahasiswa bimbingan (Approve/Reject).
- Halaman Admin: Master Data (Mata Kuliah, Kelas, Mahasiswa, Dosen).
- Axios interceptor Authorization Bearer
- UI responsif, konsisten, menggunakan tema Dark Mode Minimalis
```

## Checklist

- [ ] Login/logout
- [ ] UI Mahasiswa (KRS Form & Submit)
- [ ] UI Dosen (KRS Approval)
- [ ] UI Admin (Master Data CRUD)
- [ ] Loading & error states
- [ ] Tema Dark Mode Minimalis terimplementasi

## File target (feature-based)

- `frontend/src/features/auth/` — login, store, api
- `frontend/src/features/mahasiswa/` — form KRS
- `frontend/src/features/dosen/`, `staff/` — persetujuan
- `frontend/src/features/admin/` — users, DPA
- `frontend/src/app/router/guards.ts` — role guard
- `frontend/src/core/api/client.ts` — Axios shared
