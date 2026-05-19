<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { listKelas, createKelas, listMataKuliah, listDosen, listTahunAkademik } from '@/features/admin/api/admin.api'
import type { Kelas, MataKuliah, Dosen, TahunAkademik } from '@/features/admin/api/admin.api'

const kelasList = ref<Kelas[]>([])
const mkList = ref<MataKuliah[]>([])
const dosenList = ref<Dosen[]>([])
const tahunList = ref<TahunAkademik[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')
const showForm = ref(false)
const form = ref({ mata_kuliah_id:'', tahun_akademik_id:'', dosen_id:'', nama_kelas:'A', kapasitas: 40 })
const formLoading = ref(false)

async function load() {
  loading.value = true
  try {
    const [k, mk, d, ta] = await Promise.all([listKelas(), listMataKuliah(), listDosen(), listTahunAkademik()])
    kelasList.value = k.data ?? []
    mkList.value = mk.data ?? []
    dosenList.value = d.data ?? []
    tahunList.value = ta.data ?? []
    // Default pilih TA aktif
    const aktif = tahunList.value.find(t => t.is_active)
    if (aktif) form.value.tahun_akademik_id = aktif.id
  } catch { error.value = 'Gagal memuat data' }
  finally { loading.value = false }
}

async function handleCreate() {
  if (!form.value.mata_kuliah_id || !form.value.tahun_akademik_id || !form.value.nama_kelas) {
    error.value = 'Field wajib: Mata Kuliah, Tahun Akademik, Nama Kelas'; return
  }
  formLoading.value = true; error.value = ''
  try {
    const res = await createKelas(form.value)
    if (res.success) {
      success.value = 'Kelas berhasil dibuat'
      showForm.value = false
      await load()
    } else error.value = res.message
  } catch (e: any) { error.value = e.response?.data?.message ?? 'Gagal membuat kelas' }
  finally { formLoading.value = false }
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Manajemen Kelas</h1>
        <p class="page-sub">Penawaran kelas per semester</p>
      </div>
      <button class="btn btn-primary" @click="showForm = !showForm">
        {{ showForm ? '✕ Batal' : '+ Tambah Kelas' }}
      </button>
    </div>

    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom:1rem;">{{ error }}</div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom:1rem;">{{ success }}</div>

    <div v-if="showForm" class="card fade-in" style="margin-bottom:1.5rem;">
      <h2 class="card-title" style="margin-bottom:1.25rem;">Tambah Kelas Baru</h2>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;">
        <div class="form-field">
          <label class="form-label">Mata Kuliah</label>
          <select v-model="form.mata_kuliah_id" class="form-input">
            <option value="">— Pilih MK —</option>
            <option v-for="mk in mkList" :key="mk.id" :value="mk.id">{{ mk.kode_mk }} — {{ mk.nama_mk }}</option>
          </select>
        </div>
        <div class="form-field">
          <label class="form-label">Tahun Akademik</label>
          <select v-model="form.tahun_akademik_id" class="form-input">
            <option v-for="ta in tahunList" :key="ta.id" :value="ta.id">
              {{ ta.kode_ta }} {{ ta.semester }} {{ ta.is_active ? '(Aktif)' : '' }}
            </option>
          </select>
        </div>
        <div class="form-field">
          <label class="form-label">Dosen Pengampu</label>
          <select v-model="form.dosen_id" class="form-input">
            <option value="">— Opsional —</option>
            <option v-for="d in dosenList" :key="d.id" :value="d.id">{{ d.nama }}</option>
          </select>
        </div>
        <div class="form-field">
          <label class="form-label">Nama Kelas</label>
          <input v-model="form.nama_kelas" class="form-input" placeholder="A" maxlength="5" />
        </div>
        <div class="form-field">
          <label class="form-label">Kapasitas</label>
          <input v-model.number="form.kapasitas" class="form-input" type="number" min="1" />
        </div>
      </div>
      <div style="margin-top:1.25rem;display:flex;gap:0.75rem;">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleCreate">
          <span v-if="formLoading" class="spinner" style="width:15px;height:15px;"></span>
          {{ formLoading ? 'Menyimpan…' : 'Simpan Kelas' }}
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
            <tr><th>Mata Kuliah</th><th>Kelas</th><th>Dosen</th><th>Kuota</th></tr>
          </thead>
          <tbody>
            <tr v-for="k in kelasList" :key="k.id" class="fade-in">
              <td>
                <div style="font-weight:500;">{{ k.nama_mk }}</div>
                <div style="font-size:0.78rem;color:var(--color-muted);">{{ k.kode_mk }}</div>
              </td>
              <td><span class="badge badge-draft">{{ k.nama_kelas }}</span></td>
              <td style="color:var(--color-muted);">{{ k.nama_dosen || '-' }}</td>
              <td>
                <span :style="{ color: k.terisi >= k.kapasitas ? 'var(--color-danger)' : 'var(--color-success)' }">
                  {{ k.terisi }}/{{ k.kapasitas }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </AppLayout>
</template>
