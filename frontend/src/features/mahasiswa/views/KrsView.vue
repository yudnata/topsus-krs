<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getCurrentKRS, addKelasToKRS, removeKelasFromKRS, submitKRS } from '@/features/mahasiswa/api/krs.api'
import type { KrsHeader } from '@/features/mahasiswa/api/krs.api'
import { getAvailableKelas } from '@/features/mahasiswa/api/kelas.api'
import type { Kelas } from '@/features/mahasiswa/api/kelas.api'

const krs = ref<KrsHeader | null>(null)
const kelasList = ref<Kelas[]>([])
const loading = ref(true)
const submitting = ref(false)
const error = ref('')
const success = ref('')
const addingKelasId = ref<string | null>(null)
const removingKelasId = ref<string | null>(null)

// Search kelas
const searchQuery = ref('')
const filteredKelas = computed(() => {
  const q = searchQuery.value.toLowerCase()
  return kelasList.value.filter(k =>
    k.nama_mk.toLowerCase().includes(q) ||
    k.kode_mk.toLowerCase().includes(q) ||
    k.nama_dosen.toLowerCase().includes(q)
  )
})

// SKS info
const maxSKS = 24
const sksPercent = computed(() => Math.min(100, ((krs.value?.total_sks ?? 0) / maxSKS) * 100))
const sksColor = computed(() => sksPercent.value > 80 ? 'var(--color-warning)' : 'var(--color-accent)')

// Status info
const isEditable = computed(() => krs.value?.status === 'draft')
const statusLabel: Record<string, string> = {
  draft: 'Draft', diajukan: 'Menunggu ACC', disetujui: 'Disetujui', ditolak: 'Ditolak'
}

// Cek kelas sudah di KRS
function isInKRS(kelasId: string) {
  return krs.value?.detail?.some(d => d.kelas_id === kelasId) ?? false
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [krsRes, kelasRes] = await Promise.all([
      getCurrentKRS(),
      getAvailableKelas()
    ])
    if (krsRes.success) krs.value = krsRes.data ?? null
    if (kelasRes.success) kelasList.value = kelasRes.data ?? []
  } catch (e: any) {
    error.value = e.response?.data?.message ?? e.message ?? 'Gagal memuat data'
  } finally {
    loading.value = false
  }
}

async function addKelas(kelasId: string) {
  if (!krs.value || !isEditable.value) return
  addingKelasId.value = kelasId
  error.value = ''; success.value = ''
  try {
    const res = await addKelasToKRS(krs.value.id, kelasId)
    if (res.success) { krs.value = res.data!; success.value = 'Kelas berhasil ditambahkan' }
    else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal menambahkan kelas'
  } finally {
    addingKelasId.value = null
  }
}

async function removeKelas(kelasId: string) {
  if (!krs.value || !isEditable.value) return
  removingKelasId.value = kelasId
  error.value = ''; success.value = ''
  try {
    const res = await removeKelasFromKRS(krs.value.id, kelasId)
    if (res.success) { krs.value = res.data!; success.value = 'Kelas dihapus dari KRS' }
    else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal menghapus kelas'
  } finally {
    removingKelasId.value = null
  }
}

async function doSubmit() {
  if (!krs.value || !isEditable.value) return
  if (!confirm('Ajukan KRS? Setelah diajukan, KRS tidak bisa diedit.')) return
  submitting.value = true
  error.value = ''; success.value = ''
  try {
    const res = await submitKRS(krs.value.id)
    if (res.success) { krs.value = res.data!; success.value = 'KRS berhasil diajukan! Menunggu persetujuan dosen.' }
    else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal mengajukan KRS'
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Pengajuan KRS</h1>
        <p class="page-sub">Kartu Rencana Studi — Semester aktif</p>
      </div>
      <button
        v-if="isEditable && krs"
        class="btn btn-primary"
        :disabled="submitting || (krs?.total_sks ?? 0) === 0"
        @click="doSubmit"
      >
        <span v-if="submitting" class="spinner" style="margin-right: 0.5rem;"></span>
        {{ submitting ? 'Mengajukan…' : 'Ajukan KRS' }}
      </button>
    </div>

    <!-- Alerts -->
    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom: 1.5rem;">
      <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"/>
      </svg>
      <span>{{ error }}</span>
    </div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom: 1.5rem;">
      <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
      </svg>
      <span>{{ success }}</span>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex; align-items:center; gap:0.75rem; padding:3rem 0; color:var(--color-text-muted); font-weight: 500;">
      <span class="spinner"></span> Memuat data KRS…
    </div>

    <template v-else-if="krs">
      <!-- Status & SKS bar -->
      <div class="card fade-in" style="margin-bottom:2rem;">
        <div style="display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:1.5rem;">
          <div style="display:flex; align-items:center; gap:2rem; flex-wrap:wrap;">
            <div>
              <div style="font-size:0.75rem; font-weight:700; color:var(--color-text-muted); text-transform:uppercase; letter-spacing:0.05em; margin-bottom:0.375rem;">Status KRS</div>
              <span class="badge" :class="`badge-${krs.status}`">{{ statusLabel[krs.status] ?? krs.status }}</span>
            </div>
            <div style="width:1px; height:40px; background:var(--color-border-subtle);" class="hidden md:block"></div>
            <div>
              <div style="font-size:0.75rem; font-weight:700; color:var(--color-text-muted); text-transform:uppercase; letter-spacing:0.05em; margin-bottom:0.375rem;">Semester</div>
              <div style="font-size:1rem; font-weight:700; font-family: var(--font-mono);">{{ krs.kode_ta }} {{ krs.semester }}</div>
            </div>
          </div>
          <div style="min-width:260px; flex: 1; max-width: 400px;">
            <div style="display:flex; justify-content:space-between; margin-bottom:0.5rem; align-items: baseline;">
              <span style="font-size:0.75rem; font-weight:700; color:var(--color-text-muted); text-transform:uppercase; letter-spacing:0.05em;">Total SKS</span>
              <span style="font-size:1rem; font-weight:700; color:var(--color-text); font-family: var(--font-mono);">
                {{ krs.total_sks }} / {{ maxSKS }} SKS
              </span>
            </div>
            <div class="sks-bar">
              <div class="sks-fill" :style="{ width: sksPercent + '%', backgroundColor: sksColor }"></div>
            </div>
          </div>
        </div>

        <!-- Catatan reviewer -->
        <div v-if="krs.catatan_reviewer" class="alert alert-info" style="margin-top:1.5rem;">
          <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 111.083.87l-.467 1.327a.75.75 0 00.177.834l.211.211m-.172-4.902h.008V9.75h-.008v.007zm-2.247 11.02a9 9 0 1118 0 9 9 0 01-18 0z"/>
          </svg>
          <div>
            <strong>Catatan Reviewer:</strong> {{ krs.catatan_reviewer }}
          </div>
        </div>
      </div>

      <!-- KRS saat ini -->
      <div class="card fade-in" style="margin-bottom:2rem;">
        <div class="card-header">
          <h2 class="card-title">Kelas Dipilih ({{ krs.detail?.length ?? 0 }} kelas)</h2>
        </div>
        <div v-if="!krs.detail?.length" style="padding:3rem 2rem; text-align:center; color:var(--color-text-muted); font-weight: 500;">
          Belum ada kelas yang dipilih. Silakan pilih kelas yang tersedia di bawah.
        </div>
        <div class="table-wrap" v-else>
          <table class="table">
            <thead>
              <tr>
                <th>Mata Kuliah</th>
                <th>Kelas</th>
                <th>Dosen</th>
                <th>Jadwal</th>
                <th>SKS</th>
                <th v-if="isEditable"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="d in krs.detail" :key="d.id" class="fade-in">
                <td>
                  <div style="font-weight:700; color: var(--color-text);">{{ d.nama_mk }}</div>
                  <div style="font-size:0.75rem; color:var(--color-text-muted); font-family: var(--font-mono); margin-top: 0.125rem;">{{ d.kode_mk }}</div>
                </td>
                <td><span class="badge badge-diajukan">{{ d.nama_kelas }}</span></td>
                <td style="color:var(--color-text-muted); font-weight: 500;">{{ d.nama_dosen }}</td>
                <td>
                  <div v-for="j in d.jadwal" :key="j.hari" style="font-size:0.8125rem; font-family: var(--font-mono); color: var(--color-text);">
                    {{ j.hari }}, {{ j.jam_mulai.slice(0,5) }}–{{ j.jam_selesai.slice(0,5) }}
                  </div>
                </td>
                <td style="font-weight:700; color:var(--color-accent); font-family: var(--font-mono);">{{ d.sks }}</td>
                <td v-if="isEditable" style="text-align: right;">
                  <button
                    class="btn btn-danger btn-sm"
                    :disabled="removingKelasId === d.kelas_id"
                    @click="removeKelas(d.kelas_id)"
                    style="padding: 0.375rem 0.625rem;"
                  >
                    <span v-if="removingKelasId === d.kelas_id" class="spinner"></span>
                    <span v-else>Hapus</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pilih kelas -->
      <div class="card fade-in" v-if="isEditable">
        <div class="card-header" style="flex-wrap: wrap; gap: 1rem;">
          <h2 class="card-title">Pilih Kelas Tersedia</h2>
          <div style="position:relative; width: 100%; max-width: 300px;">
            <input
              v-model="searchQuery"
              class="form-input"
              placeholder="Cari mata kuliah, dosen…"
              style="padding-left: 2.25rem;"
            />
            <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="position:absolute; left:0.875rem; top:50%; transform:translateY(-50%); color:var(--color-text-muted);">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.637 10.637z"/>
            </svg>
          </div>
        </div>
        <div v-if="!filteredKelas.length" style="padding:3rem 2rem; text-align:center; color:var(--color-text-muted); font-weight: 500;">
          Tidak ada kelas yang cocok dengan pencarian.
        </div>
        <div class="table-wrap" v-else>
          <table class="table">
            <thead>
              <tr>
                <th>Mata Kuliah</th>
                <th>Kelas</th>
                <th>Dosen</th>
                <th>Kuota</th>
                <th>SKS</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="k in filteredKelas" :key="k.id" class="fade-in">
                <td>
                  <div style="font-weight:700; color: var(--color-text);">{{ k.nama_mk }}</div>
                  <div style="font-size:0.75rem; color:var(--color-text-muted); font-family: var(--font-mono); margin-top: 0.125rem;">{{ k.kode_mk }}</div>
                </td>
                <td><span class="badge badge-draft">{{ k.nama_kelas }}</span></td>
                <td style="color:var(--color-text-muted); font-weight: 500;">{{ k.nama_dosen || '-' }}</td>
                <td>
                  <span :style="{ color: k.terisi >= k.kapasitas ? 'var(--color-danger)' : 'var(--color-success)', fontWeight: '700', fontFamily: 'var(--font-mono)' }">
                    {{ k.terisi }}/{{ k.kapasitas }}
                  </span>
                </td>
                <td style="font-weight:700; color:var(--color-accent); font-family: var(--font-mono);">3</td>
                <td style="text-align: right;">
                  <button
                    v-if="isInKRS(k.id)"
                    class="btn btn-ghost btn-sm"
                    disabled
                    style="border-color: var(--color-border-subtle); color: var(--color-text-muted);"
                  >
                    Terpilih
                  </button>
                  <button
                    v-else
                    class="btn btn-primary btn-sm"
                    :disabled="k.terisi >= k.kapasitas || addingKelasId === k.id"
                    @click="addKelas(k.id)"
                  >
                    <span v-if="addingKelasId === k.id" class="spinner"></span>
                    <span v-else>{{ k.terisi >= k.kapasitas ? 'Penuh' : 'Pilih' }}</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </AppLayout>
</template>
