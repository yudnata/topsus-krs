<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getKRSHistory } from '@/features/mahasiswa/api/krs.api'
import type { KrsHeader } from '@/features/mahasiswa/api/krs.api'

const history = ref<KrsHeader[]>([])
const loading = ref(true)
const error = ref('')

const statusLabel: Record<string, string> = {
  draft: 'Draft', diajukan: 'Menunggu ACC', disetujui: 'Disetujui', ditolak: 'Ditolak'
}

async function load() {
  loading.value = true
  try {
    const res = await getKRSHistory()
    if (res.success) history.value = res.data ?? []
    else error.value = res.message
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal memuat riwayat'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <AppLayout>
    <div class="topbar">
      <div>
        <h1 class="page-title">Riwayat KRS</h1>
        <p class="page-sub">Semua pengajuan KRS Anda</p>
      </div>
    </div>

    <!-- Alerts -->
    <div v-if="error" class="alert alert-error" style="margin-bottom:1.5rem;">
      <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"/>
      </svg>
      <span>{{ error }}</span>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="display:flex; align-items:center; gap:0.75rem; padding:3rem 0; color:var(--color-text-muted); font-weight: 500;">
      <span class="spinner"></span> Memuat riwayat…
    </div>

    <!-- Empty State -->
    <div v-else-if="!history.length" style="text-align:center; padding:5rem 2rem; border: 1.5px solid var(--color-border); background: var(--color-surface);">
      <div style="color: var(--color-text-muted); margin-bottom: 1.25rem; display: flex; justify-content: center;">
        <svg width="48" height="48" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-19.5 0A2.25 2.25 0 0 0 4.5 15h15a2.25 2.25 0 0 0 2.25-2.25m-19.5 0v.15A2.25 2.25 0 0 0 4.5 15h15a2.25 2.25 0 0 0 2.25-2.25V14a2.25 2.25 0 0 0-2.25-2.25H4.5A2.25 2.25 0 0 0 2.25 14v.75m3-3V7.5A2.25 2.25 0 0 1 7.5 5.25h9a2.25 2.25 0 0 1 2.25 2.25v4.25"/>
        </svg>
      </div>
      <div style="font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.25rem;">Belum ada riwayat</div>
      <div style="font-size: 0.875rem; color: var(--color-text-muted);">Riwayat KRS Anda akan muncul setelah Anda mengirim pengajuan KRS pertama.</div>
    </div>

    <!-- History List -->
    <div v-else style="display:flex; flex-direction:column; gap:1.5rem;">
      <div v-for="krs in history" :key="krs.id" class="card fade-in">
        <div style="display:flex; align-items:flex-start; justify-content:space-between; flex-wrap:wrap; gap:1.5rem;">
          <div>
            <div style="font-size:1.125rem; font-weight:700; margin-bottom:0.75rem; font-family: var(--font-mono);">
              {{ krs.kode_ta }} — Semester {{ krs.semester }}
            </div>
            <div style="display:flex; gap:1rem; align-items:center; flex-wrap:wrap;">
              <span class="badge" :class="`badge-${krs.status}`">{{ statusLabel[krs.status] }}</span>
              <span style="font-size:0.8125rem; font-weight: 600; color:var(--color-text-muted); font-family: var(--font-mono);">{{ krs.total_sks }} SKS</span>
            </div>
          </div>
          <div style="text-align:right; font-family: var(--font-mono); font-size:0.8125rem; color:var(--color-text-muted); display:flex; flex-direction:column; gap:0.25rem;">
            <div>
              Diajukan: {{ krs.created_at ? new Date(krs.created_at).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) : '-' }}
            </div>
            <div v-if="krs.reviewed_at">
              Direview: {{ new Date(krs.reviewed_at).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) }}
            </div>
          </div>
        </div>
        
        <!-- Catatan reviewer -->
        <div v-if="krs.catatan_reviewer" class="alert alert-info" style="margin-top:1.5rem;">
          <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 111.083.87l-.467 1.327a.75.75 0 00.177.834l.211.211m-.172-4.902h.008V9.75h-.008v.007zm-2.247 11.02a9 9 0 1118 0 9 9 0 01-18 0z"/>
          </svg>
          <div>
            <strong>Catatan:</strong> {{ krs.catatan_reviewer }}
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
