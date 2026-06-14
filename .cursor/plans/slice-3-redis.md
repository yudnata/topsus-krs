---
name: Slice 3 — Redis Cache
overview: src/shared/cache/redis.ts dengan Get/Set/Del/TTL dan invalidate on write.
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: completed
  - id: redis-client
    content: Connect REDIS_URL dari config
    status: completed
  - id: cache-ops
    content: Get/Set/Del dengan TTL 120s
    status: completed
  - id: invalidate
    content: Pola invalidate on CRUD
    status: completed
  - id: fallback
    content: Graceful fallback jika Redis down
    status: completed
---

# Slice 3 — Redis Cache Layer

**Status:** Selesai (2026-05-19)

## Checklist

- [x] `src/shared/cache/redis.ts` — Get, Set, Del, Ping, getJSON/setJSON
- [x] `src/shared/cache/keys.ts` — konvensi key KRS + master data
- [x] `src/shared/cache/invalidator.ts` — invalidate on write (Slice 4 pakai)
- [x] Graceful fallback jika Redis down
- [x] Endpoint uji: `/api/cache/status`, `/api/cache/demo` (X-Cache: HIT|MISS)

## Uji manual

```bash
# MISS
curl -i "http://localhost:8080/api/cache/demo?key=test-krs"
# SET
curl -X POST "http://localhost:8080/api/cache/demo?key=test-krs&value=ok"
# HIT
curl -i "http://localhost:8080/api/cache/demo?key=test-krs"
```
