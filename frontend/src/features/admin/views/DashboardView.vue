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

    <!-- Loading -->
    <div v-if="loading" style="display:flex; align-items:center; gap:0.75rem; padding:3rem 0; color:var(--color-text-muted); font-weight: 500;">
      <span class="spinner"></span> Memuat data…
    </div>

    <!-- Stats Grid -->
    <div v-else class="stats-grid fade-in">
      <div class="stat-card">
        <div class="stat-icon-wrapper">
          <svg width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4.26 10.147a60.436 60.436 0 0 0-.491 6.347A48.62 48.62 0 0 1 12 20.904a48.62 48.62 0 0 1 8.232-4.41 60.46 60.46 0 0 0-.491-6.347m-15.482 0a50.57 50.57 0 0 0-2.658-.813A59.905 59.905 0 0 1 12 3.493a59.902 59.902 0 0 1 9.918 5.842 50.45 50.45 0 0 0-2.658.814m-15.482 0a50.58 50.58 0 0 0 2.658-.814M21.121 10.147a50.514 50.514 0 0 1-2.658-.814M12 14v.008M12 14a2.25 2.25 0 0 0-2.248-2.354 2.25 2.25 0 0 0-2.252 2.248m4.5 0a2.25 2.25 0 0 1 2.248-2.354M12 14v4m0 0a2.25 2.25 0 0 0-2.25-2.25M12 18a2.25 2.25 0 0 1 2.25-2.25m-4.5 0v.008"/>
          </svg>
        </div>
        <div class="stat-value">{{ dosenCount }}</div>
        <div class="stat-label">Dosen Terdaftar</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon-wrapper">
          <svg width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"/>
          </svg>
        </div>
        <div class="stat-value">{{ mahasiswaCount }}</div>
        <div class="stat-label">Mahasiswa Terdaftar</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon-wrapper">
          <svg width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-16.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-16.25v16.25"/>
          </svg>
        </div>
        <div class="stat-value">{{ mkCount }}</div>
        <div class="stat-label">Mata Kuliah</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon-wrapper">
          <svg width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h18v18H3V3Z"/>
          </svg>
        </div>
        <div class="stat-value">{{ kelasCount }}</div>
        <div class="stat-label">Kelas Tersedia</div>
      </div>
    </div>

    <!-- Quick Links -->
    <div class="card fade-in" style="margin-top:2.5rem;">
      <h2 class="card-title">Menu Cepat</h2>
      <div style="display:flex; gap:1rem; flex-wrap:wrap; margin-top:1.5rem;">
        <RouterLink to="/admin/mahasiswa" class="btn btn-primary">Tambah Mahasiswa</RouterLink>
        <RouterLink to="/admin/dosen" class="btn btn-primary">Tambah Dosen</RouterLink>
        <RouterLink to="/admin/mata-kuliah" class="btn btn-ghost">Mata Kuliah</RouterLink>
        <RouterLink to="/admin/kelas" class="btn btn-ghost">Kelas</RouterLink>
      </div>
    </div>
  </AppLayout>
</template>
