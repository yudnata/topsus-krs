<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getCurrentKRS, addKelasToKRS, removeKelasFromKRS, submitKRS } from '@/features/mahasiswa/api/krs.api'
import { getAvailableKelas } from '@/features/mahasiswa/api/kelas.api'
import type { KrsHeader } from '@/features/mahasiswa/api/krs.api'
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
  draft: 'Draft', diajukan: 'Menunggu ACC', disetujui: 'Disetujui ✓', ditolak: 'Ditolak'
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
    if (res.success) { krs.value = res.data!; success.value = '🎉 KRS berhasil diajukan! Menunggu persetujuan dosen.' }
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
        <span v-if="submitting" class="spinner" style="width:16px;height:16px;"></span>
        {{ submitting ? 'Mengajukan…' : '🚀 Ajukan KRS' }}
      </button>
    </div>

    <!-- Alerts -->
    <div v-if="error" class="alert alert-error fade-in" style="margin-bottom:1rem;">{{ error }}</div>
    <div v-if="success" class="alert alert-success fade-in" style="margin-bottom:1rem;">{{ success }}</div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex;align-items:center;gap:0.75rem;padding:3rem 0;color:var(--color-muted);">
      <span class="spinner"></span> Memuat data KRS…
    </div>

    <template v-else-if="krs">
      <!-- Status & SKS bar -->
      <div class="card fade-in" style="margin-bottom:1.5rem;">
        <div style="display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:1rem;">
          <div style="display:flex;align-items:center;gap:1rem;">
            <div>
              <div style="font-size:0.8rem;color:var(--color-muted);margin-bottom:0.25rem;">Status KRS</div>
              <span class="badge" :class="`badge-${krs.status}`">{{ statusLabel[krs.status] ?? krs.status }}</span>
            </div>
            <div style="width:1px;height:36px;background:var(--color-border);"></div>
            <div>
              <div style="font-size:0.8rem;color:var(--color-muted);margin-bottom:0.25rem;">Semester</div>
              <div style="font-size:0.9rem;font-weight:600;">{{ krs.kode_ta }} {{ krs.semester }}</div>
            </div>
          </div>
          <div style="min-width:200px;">
            <div style="display:flex;justify-content:space-between;margin-bottom:0.5rem;">
              <span style="font-size:0.8rem;color:var(--color-muted);">Total SKS</span>
              <span style="font-size:0.9rem;font-weight:700;color:var(--color-text);">
                {{ krs.total_sks }} / {{ maxSKS }} SKS
              </span>
            </div>
            <div class="sks-bar">
              <div class="sks-fill" :style="{ width: sksPercent + '%', background: `linear-gradient(90deg, ${sksColor}, #7c3aed)` }"></div>
            </div>
          </div>
        </div>

        <!-- Catatan reviewer -->
        <div v-if="krs.catatan_reviewer" class="alert alert-info" style="margin-top:1rem;">
          <strong>Catatan Reviewer:</strong> {{ krs.catatan_reviewer }}
        </div>
      </div>

      <!-- KRS saat ini -->
      <div class="card fade-in" style="margin-bottom:1.5rem;">
        <div class="card-header">
          <h2 class="card-title">📋 Kelas Dipilih ({{ krs.detail?.length ?? 0 }} kelas)</h2>
        </div>
        <div v-if="!krs.detail?.length" style="padding:2rem;text-align:center;color:var(--color-muted);">
          Belum ada kelas. Pilih dari daftar di bawah.
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
                  <div style="font-weight:500;">{{ d.nama_mk }}</div>
                  <div style="font-size:0.78rem;color:var(--color-muted);">{{ d.kode_mk }}</div>
                </td>
                <td><span class="badge badge-diajukan">{{ d.nama_kelas }}</span></td>
                <td style="color:var(--color-muted);">{{ d.nama_dosen }}</td>
                <td>
                  <div v-for="j in d.jadwal" :key="j.hari" style="font-size:0.8rem;white-space:nowrap;">
                    {{ j.hari }} {{ j.jam_mulai.slice(0,5) }}–{{ j.jam_selesai.slice(0,5) }}
                  </div>
                </td>
                <td style="font-weight:600;color:var(--color-accent);">{{ d.sks }}</td>
                <td v-if="isEditable">
                  <button
                    class="btn btn-danger btn-sm"
                    :disabled="removingKelasId === d.kelas_id"
                    @click="removeKelas(d.kelas_id)"
                  >
                    <span v-if="removingKelasId === d.kelas_id" class="spinner" style="width:12px;height:12px;"></span>
                    {{ removingKelasId === d.kelas_id ? '' : '✕' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pilih kelas -->
      <div class="card fade-in" v-if="isEditable">
        <div class="card-header">
          <h2 class="card-title">🔍 Pilih Kelas Tersedia</h2>
          <input
            v-model="searchQuery"
            class="form-input"
            placeholder="Cari mata kuliah, dosen…"
            style="max-width:260px;"
          />
        </div>
        <div v-if="!filteredKelas.length" style="padding:2rem;text-align:center;color:var(--color-muted);">
          Tidak ada kelas ditemukan.
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
                <td style="font-weight:600;color:var(--color-accent);">{{ 3 }}</td>
                <td>
                  <button
                    v-if="isInKRS(k.id)"
                    class="btn btn-ghost btn-sm"
                    disabled
                  >✓ Dipilih</button>
                  <button
                    v-else
                    class="btn btn-primary btn-sm"
                    :disabled="k.terisi >= k.kapasitas || addingKelasId === k.id"
                    @click="addKelas(k.id)"
                  >
                    <span v-if="addingKelasId === k.id" class="spinner" style="width:12px;height:12px;"></span>
                    {{ k.terisi >= k.kapasitas ? 'Penuh' : addingKelasId === k.id ? '' : '+ Pilih' }}
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
