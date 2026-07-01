<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getPendingKRS, approveKRS, rejectKRS } from '@/features/dosen/api/approval.api'
import type { PendingKRS } from '@/features/dosen/api/approval.api'
import { useAuthStore } from '@/features/auth/stores/auth.store'

const auth = useAuthStore()
const pending = ref<PendingKRS[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')

// Modal state
const showRejectModal = ref(false)
const rejectKrsId = ref('')
const rejectNote = ref('')
const processingId = ref<string | null>(null)

// Expand detail
const expandedId = ref<string | null>(null)

const pageTitle = auth.role === 'DOSEN' ? 'Persetujuan KRS Bimbingan' : 'Antrian KRS'
const pageSub = auth.role === 'DOSEN'
  ? 'KRS mahasiswa yang Anda bimbing (DPA)'
  : 'Semua pengajuan KRS yang menunggu persetujuan'

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await getPendingKRS()
    if (res.success) pending.value = res.data ?? []
    else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal memuat data'
  } finally {
    loading.value = false
  }
}

async function doApprove(krsId: string) {
  if (!confirm('Setujui KRS ini?')) return
  processingId.value = krsId
  error.value = ''; success.value = ''
  try {
    const res = await approveKRS(krsId)
    if (res.success) {
      success.value = `KRS ${res.data?.nama_mahasiswa ?? ''} berhasil disetujui`
      pending.value = pending.value.filter(p => p.id !== krsId)
    } else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal menyetujui KRS'
  } finally {
    processingId.value = null
  }
}

function openRejectModal(krsId: string) {
  rejectKrsId.value = krsId
  rejectNote.value = ''
  showRejectModal.value = true
}

async function doReject() {
  if (!rejectNote.value.trim()) { error.value = 'Catatan penolakan wajib diisi'; return }
  processingId.value = rejectKrsId.value
  error.value = ''; success.value = ''
  showRejectModal.value = false
  try {
    const res = await rejectKRS(rejectKrsId.value, rejectNote.value)
    if (res.success) {
      success.value = `KRS ${res.data?.nama_mahasiswa ?? ''} ditolak`
      pending.value = pending.value.filter(p => p.id !== rejectKrsId.value)
    } else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal menolak KRS'
  } finally {
    processingId.value = null
  }
}

function toggleDetail(id: string) {
  expandedId.value = expandedId.value === id ? null : id
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">{{ pageTitle }}</h1>
        <p class="page-sub">{{ pageSub }}</p>
      </div>
      <div style="display:flex; align-items:center; gap:0.75rem;">
        <div class="stat-card" style="padding:0.75rem 1.5rem; display:flex; flex-direction:row; align-items:center; gap: 0.75rem;">
          <div class="stat-value" style="font-size:1.5rem; color:var(--color-accent);">{{ pending.length }}</div>
          <div class="stat-label" style="margin-top: 0;">Antrian</div>
        </div>
      </div>
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

    <!-- Loading -->
    <div v-if="loading" style="display:flex; align-items:center; gap:0.75rem; padding:3rem 0; color:var(--color-text-muted); font-weight: 500;">
      <span class="spinner"></span> Memuat antrian…
    </div>

    <!-- Empty State -->
    <div v-else-if="!pending.length" style="text-align:center; padding:5rem 2rem; border: 1.5px solid var(--color-border); background: var(--color-surface);">
      <div style="color: var(--color-success); margin-bottom: 1.25rem; display: flex; justify-content: center;">
        <svg width="48" height="48" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
        </svg>
      </div>
      <div style="font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.25rem;">Tidak ada antrian KRS</div>
      <div style="font-size: 0.875rem; color: var(--color-text-muted);">Semua pengajuan KRS bimbingan Anda sudah diproses dan disetujui.</div>
    </div>

    <!-- Antrian List -->
    <div v-else style="display:flex; flex-direction:column; gap:1.5rem;">
      <div v-for="krs in pending" :key="krs.id" class="card fade-in">
        <!-- Header baris -->
        <div style="display:flex; align-items:flex-start; justify-content:space-between; gap:1.5rem; flex-wrap:wrap;">
          <div>
            <div style="font-size:0.75rem; color:var(--color-text-muted); font-family: var(--font-mono); margin-bottom:0.25rem;">NIM: {{ krs.nim }}</div>
            <div style="font-size:1.125rem; font-weight:700;">{{ krs.nama_mahasiswa }}</div>
            <div style="font-size:0.8125rem; color:var(--color-text-muted); font-family: var(--font-mono); margin-top:0.375rem;">
              {{ krs.kode_ta }} {{ krs.semester }} · {{ krs.total_sks }} SKS
            </div>
          </div>
          <div style="display:flex; gap:0.75rem; align-items:center; flex-shrink:0; flex-wrap: wrap;">
            <button
              class="btn btn-ghost btn-sm"
              @click="toggleDetail(krs.id)"
            >
              {{ expandedId === krs.id ? 'Sembunyikan' : 'Lihat Detail' }}
            </button>
            <button
              class="btn btn-danger btn-sm"
              :disabled="processingId === krs.id"
              @click="openRejectModal(krs.id)"
            >
              <span v-if="processingId === krs.id" class="spinner" style="margin-right: 0.375rem;"></span>
              Tolak
            </button>
            <button
              class="btn btn-primary btn-sm"
              :disabled="processingId === krs.id"
              @click="doApprove(krs.id)"
            >
              <span v-if="processingId === krs.id" class="spinner" style="margin-right: 0.375rem;"></span>
              Setujui
            </button>
          </div>
        </div>

        <!-- Catatan mahasiswa -->
        <div v-if="krs.catatan_mhs" class="alert alert-info" style="margin-top:1.25rem;">
          <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 111.083.87l-.467 1.327a.75.75 0 00.177.834l.211.211m-.172-4.902h.008V9.75h-.008v.007zm-2.247 11.02a9 9 0 1118 0 9 9 0 01-18 0z"/>
          </svg>
          <div>
            <strong>Catatan Mahasiswa:</strong> {{ krs.catatan_mhs }}
          </div>
        </div>

        <!-- Detail kelas -->
        <div v-if="expandedId === krs.id && krs.detail?.length" style="margin-top:1.5rem;">
          <div class="table-wrap">
            <table class="table">
              <thead>
                <tr>
                  <th>Mata Kuliah</th>
                  <th>Kelas</th>
                  <th>SKS</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="d in krs.detail" :key="d.kelas_id">
                  <td>
                    <div style="font-weight:700; color: var(--color-text);">{{ d.nama_mk }}</div>
                    <div style="font-size:0.75rem; color:var(--color-text-muted); font-family: var(--font-mono); margin-top: 0.125rem;">{{ d.kode_mk }}</div>
                  </td>
                  <td><span class="badge badge-draft">{{ d.nama_kelas }}</span></td>
                  <td style="font-weight:700; color:var(--color-accent); font-family: var(--font-mono);">{{ d.sks }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Reject Modal -->
    <div v-if="showRejectModal" class="modal-backdrop" @click.self="showRejectModal = false">
      <div class="modal-box">
        <h3 class="modal-title">Tolak KRS</h3>
        <div class="form-field" style="margin-bottom:1.5rem;">
          <label class="form-label">Alasan Penolakan <span style="color:var(--color-danger);">*</span></label>
          <textarea
            v-model="rejectNote"
            class="form-input"
            rows="4"
            placeholder="Jelaskan alasan penolakan KRS ini…"
            style="resize:vertical;"
          ></textarea>
        </div>
        <div style="display:flex; gap:0.75rem; justify-content:flex-end;">
          <button class="btn btn-ghost" @click="showRejectModal = false">Batal</button>
          <button class="btn btn-danger" @click="doReject">Tolak</button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
