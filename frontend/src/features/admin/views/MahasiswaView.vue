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
        {{ showForm ? '✕ Batal' : '+ Tambah Mahasiswa' }}
      </button>
    </div>

    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom:1rem;">{{ error }}</div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom:1rem;">{{ success }}</div>

    <!-- Form tambah -->
    <div v-if="showForm" class="card fade-in" style="margin-bottom:1.5rem;">
      <h2 class="card-title" style="margin-bottom:1.25rem;">Tambah Mahasiswa Baru</h2>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;">
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
          <input v-model="form.password" class="form-input" type="password" placeholder="Password login" />
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
      <div style="margin-top:1.25rem;display:flex;gap:0.75rem;">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleCreate">
          <span v-if="formLoading" class="spinner" style="width:15px;height:15px;"></span>
          {{ formLoading ? 'Menyimpan…' : 'Simpan Mahasiswa' }}
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
            <tr><th>NIM</th><th>Nama</th><th>Prodi</th><th>DPA</th><th>Max SKS</th><th>Aksi</th></tr>
          </thead>
          <tbody>
            <tr v-for="m in mahasiswaList" :key="m.id" class="fade-in">
              <td style="font-family:monospace;font-size:0.85rem;">{{ m.nim }}</td>
              <td style="font-weight:500;">{{ m.nama }}</td>
              <td style="color:var(--color-muted);">{{ m.nama_prodi }}</td>
              <td>
                <span v-if="m.nama_dosen" style="font-size:0.85rem;">{{ m.nama_dosen }}</span>
                <span v-else style="color:var(--color-danger);font-size:0.8rem;">Belum dipetakan</span>
              </td>
              <td style="text-align:center;">{{ m.max_sks }}</td>
              <td>
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
        <h3 class="modal-title">📌 Set Dosen Pembimbing (DPA)</h3>
        <p style="font-size:0.875rem;color:var(--color-muted);margin-bottom:1.25rem;">
          Mahasiswa: <strong style="color:var(--color-text);">{{ dpaTarget?.nama }}</strong>
        </p>
        <div class="form-field" style="margin-bottom:1.25rem;">
          <label class="form-label">Dosen Pembimbing Akademik</label>
          <select v-model="dpaDosenId" class="form-input">
            <option value="">— Pilih Dosen —</option>
            <option v-for="d in dosenList" :key="d.id" :value="d.id">{{ d.nama }}</option>
          </select>
        </div>
        <div style="display:flex;gap:0.75rem;justify-content:flex-end;">
          <button class="btn btn-ghost" @click="dpaModal = false">Batal</button>
          <button class="btn btn-primary" :disabled="dpaLoading || !dpaDosenId" @click="saveDPA">
            <span v-if="dpaLoading" class="spinner" style="width:14px;height:14px;"></span>
            {{ dpaLoading ? 'Menyimpan…' : 'Simpan DPA' }}
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
