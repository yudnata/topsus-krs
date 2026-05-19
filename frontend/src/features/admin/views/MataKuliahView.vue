<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { listMataKuliah, createMataKuliah, listProdi } from '@/features/admin/api/admin.api'
import type { MataKuliah, ProgramStudi } from '@/features/admin/api/admin.api'

const mkList = ref<MataKuliah[]>([])
const prodiList = ref<ProgramStudi[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')
const showForm = ref(false)
const form = ref({ kode_mk:'', nama_mk:'', sks: 3, prodi_id:'' })
const formLoading = ref(false)

async function load() {
  loading.value = true
  try {
    const [mk, p] = await Promise.all([listMataKuliah(), listProdi()])
    mkList.value = mk.data ?? []
    prodiList.value = p.data ?? []
  } catch { error.value = 'Gagal memuat data' }
  finally { loading.value = false }
}

async function handleCreate() {
  if (!form.value.kode_mk || !form.value.nama_mk || !form.value.prodi_id) {
    error.value = 'Semua field wajib diisi'; return
  }
  formLoading.value = true; error.value = ''
  try {
    const res = await createMataKuliah(form.value)
    if (res.success) {
      success.value = 'Mata kuliah berhasil dibuat'
      showForm.value = false
      form.value = { kode_mk:'', nama_mk:'', sks: 3, prodi_id:'' }
      await load()
    } else error.value = res.message
  } catch (e: any) { error.value = e.response?.data?.message ?? 'Gagal membuat mata kuliah' }
  finally { formLoading.value = false }
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Mata Kuliah</h1>
        <p class="page-sub">Katalog mata kuliah aktif</p>
      </div>
      <button class="btn btn-primary" @click="showForm = !showForm">
        {{ showForm ? '✕ Batal' : '+ Tambah MK' }}
      </button>
    </div>

    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom:1rem;">{{ error }}</div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom:1rem;">{{ success }}</div>

    <div v-if="showForm" class="card fade-in" style="margin-bottom:1.5rem;">
      <h2 class="card-title" style="margin-bottom:1.25rem;">Tambah Mata Kuliah</h2>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;">
        <div class="form-field">
          <label class="form-label">Kode MK</label>
          <input v-model="form.kode_mk" class="form-input" placeholder="IF101" />
        </div>
        <div class="form-field">
          <label class="form-label">Nama Mata Kuliah</label>
          <input v-model="form.nama_mk" class="form-input" placeholder="Algoritma & Pemrograman" />
        </div>
        <div class="form-field">
          <label class="form-label">SKS</label>
          <input v-model.number="form.sks" class="form-input" type="number" min="1" max="6" />
        </div>
        <div class="form-field">
          <label class="form-label">Program Studi</label>
          <select v-model="form.prodi_id" class="form-input">
            <option value="">— Pilih Prodi —</option>
            <option v-for="p in prodiList" :key="p.id" :value="p.id">{{ p.nama_prodi }}</option>
          </select>
        </div>
      </div>
      <div style="margin-top:1.25rem;display:flex;gap:0.75rem;">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleCreate">
          <span v-if="formLoading" class="spinner" style="width:15px;height:15px;"></span>
          {{ formLoading ? 'Menyimpan…' : 'Simpan' }}
        </button>
        <button class="btn btn-ghost" @click="showForm = false">Batal</button>
      </div>
    </div>

    <div class="card fade-in">
      <div v-if="loading" style="display:flex;align-items:center;gap:0.75rem;padding:2rem;color:var(--color-muted);">
        <span class="spinner"></span> Memuat…
      </div>
      <div v-else class="table-wrap">
        <table class="table">
          <thead>
            <tr><th>Kode</th><th>Nama Mata Kuliah</th><th>SKS</th></tr>
          </thead>
          <tbody>
            <tr v-for="mk in mkList" :key="mk.id" class="fade-in">
              <td><span class="badge badge-diajukan">{{ mk.kode_mk }}</span></td>
              <td style="font-weight:500;">{{ mk.nama_mk }}</td>
              <td style="font-weight:600;color:var(--color-accent);">{{ mk.sks }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </AppLayout>
</template>
