<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/auth.store'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const role = computed(() => auth.role)
const userName = computed(() => auth.user?.email?.split('@')[0] ?? 'User')
const roleBadgeMap: Record<string, string> = {
  ADMIN: '🛡️ Admin', MAHASISWA: '🎓 Mahasiswa', DOSEN: '📚 Dosen', STAFF: '🏢 Staff'
}
const roleLabel = computed(() => roleBadgeMap[role.value ?? ''] ?? role.value)

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="app-shell">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-logo">
        <div class="sidebar-logo-icon">📋</div>
        <div>
          <div class="sidebar-logo-text">KRS System</div>
          <div class="sidebar-logo-sub">{{ roleLabel }}</div>
        </div>
      </div>

      <nav style="flex:1; display:flex; flex-direction:column; gap:0.25rem;">
        <!-- MAHASISWA -->
        <template v-if="role === 'MAHASISWA'">
          <RouterLink to="/mahasiswa/krs" class="nav-item" :class="{ active: route.path.startsWith('/mahasiswa') }">
            <span class="nav-icon">📄</span> Pengajuan KRS
          </RouterLink>
          <RouterLink to="/mahasiswa/history" class="nav-item" :class="{ active: route.path === '/mahasiswa/history' }">
            <span class="nav-icon">📂</span> Riwayat KRS
          </RouterLink>
        </template>

        <!-- DOSEN -->
        <template v-if="role === 'DOSEN'">
          <RouterLink to="/dosen/persetujuan" class="nav-item" :class="{ active: route.path.startsWith('/dosen') }">
            <span class="nav-icon">✅</span> Persetujuan KRS
          </RouterLink>
        </template>

        <!-- STAFF -->
        <template v-if="role === 'STAFF'">
          <RouterLink to="/staff/persetujuan" class="nav-item" :class="{ active: route.path.startsWith('/staff') }">
            <span class="nav-icon">📋</span> Antrian KRS
          </RouterLink>
        </template>

        <!-- ADMIN -->
        <template v-if="role === 'ADMIN'">
          <RouterLink to="/admin/dashboard" class="nav-item" :class="{ active: route.path === '/admin/dashboard' }">
            <span class="nav-icon">📊</span> Dashboard
          </RouterLink>
          <RouterLink to="/admin/mahasiswa" class="nav-item" :class="{ active: route.path === '/admin/mahasiswa' }">
            <span class="nav-icon">🎓</span> Mahasiswa
          </RouterLink>
          <RouterLink to="/admin/dosen" class="nav-item" :class="{ active: route.path === '/admin/dosen' }">
            <span class="nav-icon">📚</span> Dosen
          </RouterLink>
          <RouterLink to="/admin/mata-kuliah" class="nav-item" :class="{ active: route.path === '/admin/mata-kuliah' }">
            <span class="nav-icon">📖</span> Mata Kuliah
          </RouterLink>
          <RouterLink to="/admin/kelas" class="nav-item" :class="{ active: route.path === '/admin/kelas' }">
            <span class="nav-icon">🏫</span> Kelas
          </RouterLink>
        </template>
      </nav>

      <!-- User footer -->
      <div style="border-top: 1px solid var(--color-border); padding-top: 1rem; margin-top: auto;">
        <div style="padding: 0.5rem 0.75rem; margin-bottom: 0.5rem;">
          <div style="font-size: 0.8rem; font-weight: 600; color: var(--color-text);">{{ userName }}</div>
          <div style="font-size: 0.72rem; color: var(--color-muted);">{{ auth.user?.email }}</div>
        </div>
        <button class="nav-item btn-danger" style="width: 100%; border: none;" @click="logout">
          <span class="nav-icon">🚪</span> Keluar
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>
