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

    <div v-if="error" class="alert alert-error" style="margin-bottom:1rem;">{{ error }}</div>

    <div v-if="loading" style="display:flex;align-items:center;gap:0.75rem;padding:3rem 0;color:var(--color-muted);">
      <span class="spinner"></span> Memuat riwayat…
    </div>

    <div v-else-if="!history.length" style="text-align:center;padding:4rem;color:var(--color-muted);">
      <div style="font-size:2.5rem;margin-bottom:1rem;">📂</div>
      <div>Belum ada riwayat KRS.</div>
    </div>

    <div v-else style="display:flex;flex-direction:column;gap:1rem;">
      <div v-for="krs in history" :key="krs.id" class="card fade-in">
        <div style="display:flex;align-items:flex-start;justify-content:space-between;flex-wrap:wrap;gap:1rem;">
          <div>
            <div style="font-size:1rem;font-weight:600;margin-bottom:0.5rem;">
              {{ krs.kode_ta }} — Semester {{ krs.semester }}
            </div>
            <div style="display:flex;gap:0.75rem;align-items:center;flex-wrap:wrap;">
              <span class="badge" :class="`badge-${krs.status}`">{{ statusLabel[krs.status] }}</span>
              <span style="font-size:0.8rem;color:var(--color-muted);">{{ krs.total_sks }} SKS</span>
            </div>
          </div>
          <div style="text-align:right;">
            <div style="font-size:0.78rem;color:var(--color-muted);">
              Diajukan: {{ krs.created_at ? new Date(krs.created_at).toLocaleDateString('id-ID') : '-' }}
            </div>
            <div v-if="krs.reviewed_at" style="font-size:0.78rem;color:var(--color-muted);">
              Direview: {{ new Date(krs.reviewed_at).toLocaleDateString('id-ID') }}
            </div>
          </div>
        </div>
        <div v-if="krs.catatan_reviewer" class="alert alert-info" style="margin-top:1rem;font-size:0.85rem;">
          <strong>Catatan:</strong> {{ krs.catatan_reviewer }}
        </div>
      </div>
    </div>
  </AppLayout>
</template>
