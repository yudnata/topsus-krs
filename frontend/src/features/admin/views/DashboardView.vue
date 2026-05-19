<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { listDosen, listMahasiswa, listMataKuliah, listKelas } from '@/features/admin/api/admin.api'

const dosenCount = ref(0)
const mahasiswaCount = ref(0)
const mkCount = ref(0)
const kelasCount = ref(0)
const loading = ref(true)

onMounted(async () => {
  try {
    const [d, m, mk, k] = await Promise.all([listDosen(), listMahasiswa(), listMataKuliah(), listKelas()])
    dosenCount.value   = d.data?.length ?? 0
    mahasiswaCount.value = m.data?.length ?? 0
    mkCount.value      = mk.data?.length ?? 0
    kelasCount.value   = k.data?.length ?? 0
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Dashboard Admin</h1>
        <p class="page-sub">Ringkasan data master Sistem KRS</p>
      </div>
    </div>

    <div v-if="loading" style="display:flex;align-items:center;gap:0.75rem;padding:3rem 0;color:var(--color-muted);">
      <span class="spinner"></span> Memuat data…
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:1rem;" class="fade-in">
      <div class="stat-card">
        <div style="font-size:1.75rem;margin-bottom:0.25rem;">📚</div>
        <div class="stat-value">{{ dosenCount }}</div>
        <div class="stat-label">Dosen Terdaftar</div>
      </div>
      <div class="stat-card">
        <div style="font-size:1.75rem;margin-bottom:0.25rem;">🎓</div>
        <div class="stat-value">{{ mahasiswaCount }}</div>
        <div class="stat-label">Mahasiswa Terdaftar</div>
      </div>
      <div class="stat-card">
        <div style="font-size:1.75rem;margin-bottom:0.25rem;">📖</div>
        <div class="stat-value">{{ mkCount }}</div>
        <div class="stat-label">Mata Kuliah</div>
      </div>
      <div class="stat-card">
        <div style="font-size:1.75rem;margin-bottom:0.25rem;">🏫</div>
        <div class="stat-value">{{ kelasCount }}</div>
        <div class="stat-label">Kelas Tersedia</div>
      </div>
    </div>

    <div class="card fade-in" style="margin-top:2rem;">
      <h2 class="card-title" style="margin-bottom:1rem;">🔗 Menu Cepat</h2>
      <div style="display:flex;gap:0.75rem;flex-wrap:wrap;">
        <RouterLink to="/admin/mahasiswa" class="btn btn-primary">+ Tambah Mahasiswa</RouterLink>
        <RouterLink to="/admin/dosen" class="btn btn-primary">+ Tambah Dosen</RouterLink>
        <RouterLink to="/admin/mata-kuliah" class="btn btn-ghost">📖 Mata Kuliah</RouterLink>
        <RouterLink to="/admin/kelas" class="btn btn-ghost">🏫 Kelas</RouterLink>
      </div>
    </div>
  </AppLayout>
</template>
