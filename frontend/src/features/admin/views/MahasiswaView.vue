<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { listMahasiswa, listDosen, listProdi, createMahasiswa, patchDPA } from '@/features/admin/api/admin.api'
import type { Mahasiswa, Dosen, ProgramStudi } from '@/features/admin/api/admin.api'

const mahasiswaList = ref<Mahasiswa[]>([])
const dosenList = ref<Dosen[]>([])
const prodiList = ref<ProgramStudi[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')

// Form tambah mahasiswa
const showForm = ref(false)
const form = ref({ nim:'', nama:'', email:'', password:'', prodi_id:'', dosen_pembimbing_id:'', max_sks: 24 })
const formLoading = ref(false)

// DPA patch
const dpaModal = ref(false)
const dpaTarget = ref<Mahasiswa | null>(null)
const dpaDosenId = ref('')
const dpaLoading = ref(false)

async function load() {
  loading.value = true
  try {
    const [m, d, p] = await Promise.all([listMahasiswa(), listDosen(), listProdi()])
    mahasiswaList.value = m.data ?? []
    dosenList.value = d.data ?? []
    prodiList.value = p.data ?? []
  } catch { error.value = 'Gagal memuat data' }
  finally { loading.value = false }
}

async function handleCreate() {
  if (!form.value.nim || !form.value.nama || !form.value.email || !form.value.password || !form.value.prodi_id) {
    error.value = 'Semua field wajib diisi'; return
  }
  formLoading.value = true; error.value = ''
  try {
    const res = await createMahasiswa({
      nim: form.value.nim, nama: form.value.nama,
      email: form.value.email, password: form.value.password,
      prodi_id: form.value.prodi_id,
      dosen_pembimbing_id: form.value.dosen_pembimbing_id || undefined,
    })
    if (res.success) {
      success.value = 'Mahasiswa berhasil dibuat'
      showForm.value = false
      form.value = { nim:'', nama:'', email:'', password:'', prodi_id:'', dosen_pembimbing_id:'', max_sks: 24 }
      await load()
    } else error.value = res.message
  } catch (e: any) { error.value = e.response?.data?.message ?? 'Gagal membuat mahasiswa' }
  finally { formLoading.value = false }
}

function openDPA(m: Mahasiswa) {
  dpaTarget.value = m
  dpaDosenId.value = m.dosen_pembimbing_id ?? ''
  dpaModal.value = true
}
async function saveDPA() {
  if (!dpaTarget.value || !dpaDosenId.value) return
  dpaLoading.value = true; error.value = ''
  try {
    const res = await patchDPA(dpaTarget.value.id, dpaDosenId.value)
    if (res.success) { success.value = 'DPA berhasil diperbarui'; dpaModal.value = false; await load() }
    else error.value = res.message
  } catch (e: any) { error.value = e.response?.data?.message ?? 'Gagal update DPA' }
  finally { dpaLoading.value = false }
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Manajemen Mahasiswa</h1>
        <p class="page-sub">Daftar mahasiswa & pemetaan DPA</p>
      </div>
      <button class="btn btn-primary" @click="showForm = !showForm">
        {{ showForm ? 'Batal' : 'Tambah Mahasiswa' }}
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

    <!-- Form tambah -->
    <div v-if="showForm" class="card fade-in" style="margin-bottom:2rem;">
      <h2 class="card-title" style="margin-bottom:1.5rem;">Tambah Mahasiswa Baru</h2>
      <div style="display:grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap:1.25rem;">
        <div class="form-field">
          <label class="form-label">NIM</label>
          <input v-model="form.nim" class="form-input" placeholder="2201001" />
        </div>
        <div class="form-field">
          <label class="form-label">Nama Lengkap</label>
          <input v-model="form.nama" class="form-input" placeholder="Nama mahasiswa" />
        </div>
        <div class="form-field">
          <label class="form-label">Email</label>
          <input v-model="form.email" class="form-input" type="email" placeholder="nim@student.ac.id" />
        </div>
        <div class="form-field">
          <label class="form-label">Password</label>
          <input v-model="form.password" class="form-input" type="password" placeholder="••••••••" />
        </div>
        <div class="form-field">
          <label class="form-label">Program Studi</label>
          <select v-model="form.prodi_id" class="form-input">
            <option value="">— Pilih Prodi —</option>
            <option v-for="p in prodiList" :key="p.id" :value="p.id">{{ p.nama_prodi }}</option>
          </select>
        </div>
        <div class="form-field">
          <label class="form-label">Dosen Pembimbing (DPA)</label>
          <select v-model="form.dosen_pembimbing_id" class="form-input">
            <option value="">— Opsional —</option>
            <option v-for="d in dosenList" :key="d.id" :value="d.id">{{ d.nama }}</option>
          </select>
        </div>
      </div>
      <div style="margin-top:1.5rem; display:flex; gap:0.75rem;">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleCreate">
          <span v-if="formLoading" class="spinner" style="margin-right:0.5rem;"></span>
          Simpan Mahasiswa
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
            <tr><th>NIM</th><th>Nama</th><th>Prodi</th><th>DPA</th><th>Max SKS</th><th></th></tr>
          </thead>
          <tbody>
            <tr v-for="m in mahasiswaList" :key="m.id" class="fade-in">
              <td style="font-family:var(--font-mono); font-size:0.8125rem;">{{ m.nim }}</td>
              <td style="font-weight:700; color: var(--color-text);">{{ m.nama }}</td>
              <td style="color:var(--color-text-muted); font-weight:500;">{{ m.nama_prodi }}</td>
              <td>
                <span v-if="m.nama_dosen" style="font-weight: 500;">{{ m.nama_dosen }}</span>
                <span v-else style="color:var(--color-danger); font-size:0.8125rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.02em;">Belum dipetakan</span>
              </td>
              <td style="font-family:var(--font-mono); font-weight: 700;">{{ m.max_sks }}</td>
              <td style="text-align: right;">
                <button class="btn btn-ghost btn-sm" @click="openDPA(m)">Set DPA</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- DPA Modal -->
    <div v-if="dpaModal" class="modal-backdrop" @click.self="dpaModal = false">
      <div class="modal-box">
        <h3 class="modal-title">Set Dosen Pembimbing (DPA)</h3>
        <p style="font-size:0.875rem; color:var(--color-text-muted); margin-bottom:1.5rem; margin-top: -0.5rem;">
          Mahasiswa: <strong style="color:var(--color-text);">{{ dpaTarget?.nama }}</strong>
        </p>
        <div class="form-field" style="margin-bottom:1.5rem;">
          <label class="form-label">Dosen Pembimbing Akademik</label>
          <select v-model="dpaDosenId" class="form-input">
            <option value="">— Pilih Dosen —</option>
            <option v-for="d in dosenList" :key="d.id" :value="d.id">{{ d.nama }}</option>
          </select>
        </div>
        <div style="display:flex; gap:0.75rem; justify-content:flex-end;">
          <button class="btn btn-ghost" @click="dpaModal = false">Batal</button>
          <button class="btn btn-primary" :disabled="dpaLoading || !dpaDosenId" @click="saveDPA">
            <span v-if="dpaLoading" class="spinner" style="margin-right:0.5rem;"></span>
            Simpan DPA
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
