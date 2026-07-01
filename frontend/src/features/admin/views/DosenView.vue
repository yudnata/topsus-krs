<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { listDosen, createDosen } from '@/features/admin/api/admin.api'
import type { Dosen } from '@/features/admin/api/admin.api'

const dosenList = ref<Dosen[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')
const showForm = ref(false)
const form = ref({ nip:'', nama:'', email:'', password:'' })
const formLoading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await listDosen()
    dosenList.value = res.data ?? []
  } catch { error.value = 'Gagal memuat data dosen' }
  finally { loading.value = false }
}

async function handleCreate() {
  if (!form.value.nip || !form.value.nama || !form.value.email || !form.value.password) {
    error.value = 'Semua field wajib diisi'; return
  }
  formLoading.value = true; error.value = ''
  try {
    const res = await createDosen(form.value)
    if (res.success) {
      success.value = 'Dosen berhasil dibuat'
      showForm.value = false
      form.value = { nip:'', nama:'', email:'', password:'' }
      await load()
    } else error.value = res.message
  } catch (e: any) { error.value = e.response?.data?.message ?? 'Gagal membuat dosen' }
  finally { formLoading.value = false }
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Manajemen Dosen</h1>
        <p class="page-sub">Tambah & lihat daftar dosen</p>
      </div>
      <button class="btn btn-primary" @click="showForm = !showForm">
        {{ showForm ? 'Batal' : 'Tambah Dosen' }}
      </button>
    </div>

    <!-- Alerts -->
    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom:1.5rem;">
      <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"/>
      </svg>
      <span>{{ error }}</span>
    </div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom:1.5rem;">
      <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
      </svg>
      <span>{{ success }}</span>
    </div>

    <!-- Form -->
    <div v-if="showForm" class="card fade-in" style="margin-bottom:2rem;">
      <h2 class="card-title" style="margin-bottom:1.5rem;">Tambah Dosen Baru</h2>
      <div style="display:grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap:1.25rem;">
        <div class="form-field">
          <label class="form-label">NIP</label>
          <input v-model="form.nip" class="form-input" placeholder="19800101001" />
        </div>
        <div class="form-field">
          <label class="form-label">Nama Lengkap</label>
          <input v-model="form.nama" class="form-input" placeholder="Dr. Nama Dosen" />
        </div>
        <div class="form-field">
          <label class="form-label">Email</label>
          <input v-model="form.email" class="form-input" type="email" placeholder="dosen@kampus.ac.id" />
        </div>
        <div class="form-field">
          <label class="form-label">Password</label>
          <input v-model="form.password" class="form-input" type="password" placeholder="••••••••" />
        </div>
      </div>
      <div style="margin-top:1.5rem; display:flex; gap:0.75rem;">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleCreate">
          <span v-if="formLoading" class="spinner" style="margin-right:0.5rem;"></span>
          Simpan Dosen
        </button>
        <button class="btn btn-ghost" @click="showForm = false">Batal</button>
      </div>
    </div>

    <!-- Table -->
    <div class="card fade-in">
      <div v-if="loading" style="display:flex; align-items:center; gap:0.75rem; padding:2rem; color:var(--color-text-muted); font-weight: 500;">
        <span class="spinner"></span> Memuat…
      </div>
      <div v-else class="table-wrap">
        <table class="table">
          <thead>
            <tr><th>NIP</th><th>Nama</th><th>Email</th></tr>
          </thead>
          <tbody>
            <tr v-for="d in dosenList" :key="d.id" class="fade-in">
              <td style="font-family:var(--font-mono); font-size:0.8125rem;">{{ d.nip }}</td>
              <td style="font-weight:700; color: var(--color-text);">{{ d.nama }}</td>
              <td style="color:var(--color-text-muted); font-weight: 500;">{{ d.email }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </AppLayout>
</template>
