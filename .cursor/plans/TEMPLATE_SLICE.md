---
name: "Slice X — [Nama Slice]"
overview: "[Ringkasan satu kalimat]"
todos:
  - id: plan-mode
    content: Plan Mode dijalankan & plan tersimpan
    status: pending
  - id: todo-1
    content: "[Checklist item 1 dari CHECKPOINT.md]"
    status: pending
  - id: todo-2
    content: "[Checklist item 2]"
    status: pending
---

# Slice X — [Nama Slice]

**Status CHECKPOINT:** Belum mulai / Sedang / Selesai  
**Plan file:** `.cursor/plans/slice-X-nama.md`

## Tujuan

[Deskripsi deliverable]

## Prompt Plan Mode

```text
[Salin dari CHECKPOINT.md section slice ini]
```

## Checklist (mirror CHECKPOINT.md)

- [ ] Item 1
- [ ] Item 2

## File yang diharapkan berubah

- `path/to/file`

## Setelah selesai

1. Update tabel di [CHECKPOINT.md](../../CHECKPOINT.md)
2. Update [CONTEXT.md](../../CONTEXT.md) — slice aktif + plan file berikutnya
3. Isi [PROMPT_LOG.md](../../PROMPT_LOG.md) — PLAN MODE + prompt eksekusi
