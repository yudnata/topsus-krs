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
        {{ showForm ? '✕ Batal' : '+ Tambah Dosen' }}
      </button>
    </div>

    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom:1rem;">{{ error }}</div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom:1rem;">{{ success }}</div>

    <!-- Form -->
    <div v-if="showForm" class="card fade-in" style="margin-bottom:1.5rem;">
      <h2 class="card-title" style="margin-bottom:1.25rem;">Tambah Dosen Baru</h2>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;">
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
          <input v-model="form.password" class="form-input" type="password" />
        </div>
      </div>
      <div style="margin-top:1.25rem;display:flex;gap:0.75rem;">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleCreate">
          <span v-if="formLoading" class="spinner" style="width:15px;height:15px;"></span>
          {{ formLoading ? 'Menyimpan…' : 'Simpan Dosen' }}
        </button>
        <button class="btn btn-ghost" @click="showForm = false">Batal</button>
      </div>
    </div>

    <!-- Table -->
    <div class="card fade-in">
      <div v-if="loading" style="display:flex;align-items:center;gap:0.75rem;padding:2rem;color:var(--color-muted);">
        <span class="spinner"></span> Memuat…
      </div>
      <div v-else class="table-wrap">
        <table class="table">
          <thead>
            <tr><th>NIP</th><th>Nama</th><th>Email</th></tr>
          </thead>
          <tbody>
            <tr v-for="d in dosenList" :key="d.id" class="fade-in">
              <td style="font-family:monospace;font-size:0.85rem;">{{ d.nip }}</td>
              <td style="font-weight:500;">{{ d.nama }}</td>
              <td style="color:var(--color-muted);">{{ d.email }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </AppLayout>
</template>
