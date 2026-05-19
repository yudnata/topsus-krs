---
name: Slice 6 — Integrasi E2E
overview: CORS, alur login→CRUD→DPA, verifikasi cache invalidation.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: done
  - id: cors-check
    content: Verifikasi konfigurasi CORS untuk localhost dan production
    status: done
  - id: e2e-flow
    content: Uji alur E2E lengkap lulus manual (Master -> DPA -> KRS -> Approval)
    status: done
  - id: cache-test
    content: Cache invalidation terverifikasi (GET hit -> UPDATE -> GET fresh)
    status: done
  - id: auth-check
    content: Verifikasi tidak ada regresi auth (guard & interceptor berfungsi)
    status: done
---

# Slice 6 — Integrasi & Testing E2E

## Prompt Plan Mode

```text
Buatkan plan integrasi E2E Sistem KRS.
- CORS_ORIGIN mengizinkan localhost:5173 dan URL Vercel production
- Uji alur: login → CRUD dosen → CRUD mahasiswa → assign DPA → verifikasi list per dosen
- Uji cache: GET dua kali (hit), lalu UPDATE dan GET lagi (data fresh)
- Dokumentasi hasil uji singkat di PROMPT_LOG
```

## Checklist

- [ ] CORS OK dev
- [ ] E2E flow lulus
- [ ] Cache invalidation OK
- [ ] PROMPT_LOG updated
