<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/features/auth/stores/auth.store'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

// Redirect berdasarkan role
const roleDest = computed(() => {
  switch (auth.role) {
    case 'MAHASISWA': return '/mahasiswa/krs'
    case 'DOSEN': return '/dosen/persetujuan'
    case 'STAFF': return '/staff/persetujuan'
    case 'ADMIN': return '/admin/dashboard'
    default: return null
  }
})

if (auth.isAuthenticated && roleDest.value) {
  router.replace(roleDest.value)
} else if (!auth.isAuthenticated) {
  router.replace('/login')
}
</script>

<template>
  <div style="display:flex;align-items:center;justify-content:center;min-height:100vh;">
    <span class="spinner"></span>
  </div>
</template>
