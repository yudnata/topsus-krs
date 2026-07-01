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
  ADMIN: 'Admin', MAHASISWA: 'Mahasiswa', DOSEN: 'Dosen', STAFF: 'Staff'
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
        <div class="sidebar-logo-icon">
          <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 8.25V6a2.25 2.25 0 0 0-2.25-2.25H6A2.25 2.25 0 0 0 3.75 6v12A2.25 2.25 0 0 0 6 20.25h12A2.25 2.25 0 0 0 20.25 18V9.75A2.25 2.25 0 0 0 18 7.5h-1.5ZM13.5 3.75v3c0 .621.504 1.125 1.125 1.125h3m-9 3.75h4.5m-4.5 3h7.5"/>
          </svg>
        </div>
        <div>
          <div class="sidebar-logo-text">KRS System</div>
          <div class="sidebar-logo-sub">{{ roleLabel }}</div>
        </div>
      </div>

      <nav style="flex:1; display:flex; flex-direction:column; gap:0.375rem;">
        <!-- MAHASISWA -->
        <template v-if="role === 'MAHASISWA'">
          <RouterLink to="/mahasiswa/krs" class="nav-item" :class="{ active: route.path.startsWith('/mahasiswa/krs') }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"/>
            </svg>
            Pengajuan KRS
          </RouterLink>
          <RouterLink to="/mahasiswa/history" class="nav-item" :class="{ active: route.path === '/mahasiswa/history' }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>
            Riwayat KRS
          </RouterLink>
        </template>

        <!-- DOSEN -->
        <template v-if="role === 'DOSEN'">
          <RouterLink to="/dosen/persetujuan" class="nav-item" :class="{ active: route.path.startsWith('/dosen') }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>
            Persetujuan KRS
          </RouterLink>
        </template>

        <!-- STAFF -->
        <template v-if="role === 'STAFF'">
          <RouterLink to="/staff/persetujuan" class="nav-item" :class="{ active: route.path.startsWith('/staff') }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>
            Antrian KRS
          </RouterLink>
        </template>

        <!-- ADMIN -->
        <template v-if="role === 'ADMIN'">
          <RouterLink to="/admin/dashboard" class="nav-item" :class="{ active: route.path === '/admin/dashboard' }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 14.25v2.25m3-4.5v4.5m3-6.75v6.75m3-9v9M6 20.25h12A2.25 2.25 0 0 0 20.25 18V6A2.25 2.25 0 0 0 18 3.75H6A2.25 2.25 0 0 0 3.75 6v12A2.25 2.25 0 0 0 6 20.25Z"/>
            </svg>
            Dashboard
          </RouterLink>
          <RouterLink to="/admin/mahasiswa" class="nav-item" :class="{ active: route.path === '/admin/mahasiswa' }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"/>
            </svg>
            Mahasiswa
          </RouterLink>
          <RouterLink to="/admin/dosen" class="nav-item" :class="{ active: route.path === '/admin/dosen' }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4.26 10.147a60.436 60.436 0 0 0-.491 6.347A48.62 48.62 0 0 1 12 20.904a48.62 48.62 0 0 1 8.232-4.41 60.46 60.46 0 0 0-.491-6.347m-15.482 0a50.57 50.57 0 0 0-2.658-.813A59.905 59.905 0 0 1 12 3.493a59.902 59.902 0 0 1 9.918 5.842 50.45 50.45 0 0 0-2.658.814m-15.482 0a50.58 50.58 0 0 0 2.658-.814M21.121 10.147a50.514 50.514 0 0 1-2.658-.814M12 14v.008M12 14a2.25 2.25 0 0 0-2.248-2.354 2.25 2.25 0 0 0-2.252 2.248m4.5 0a2.25 2.25 0 0 1 2.248-2.354M12 14v4m0 0a2.25 2.25 0 0 0-2.25-2.25M12 18a2.25 2.25 0 0 1 2.25-2.25m-4.5 0v.008"/>
            </svg>
            Dosen
          </RouterLink>
          <RouterLink to="/admin/mata-kuliah" class="nav-item" :class="{ active: route.path === '/admin/mata-kuliah' }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-16.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-16.25v16.25"/>
            </svg>
            Mata Kuliah
          </RouterLink>
          <RouterLink to="/admin/kelas" class="nav-item" :class="{ active: route.path === '/admin/kelas' }">
            <svg class="nav-icon" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h18v18H3V3Z"/>
            </svg>
            Kelas
          </RouterLink>
        </template>
      </nav>

      <!-- User footer -->
      <div style="border-top: 1.5px solid var(--color-border); padding-top: 1.5rem; margin-top: auto; display: flex; flex-direction: column; gap: 1rem;">
        <div style="padding: 0 0.5rem;">
          <div style="font-size: 0.875rem; font-weight: 700; color: var(--color-text); text-transform: uppercase; letter-spacing: 0.05em;">{{ userName }}</div>
          <div style="font-size: 0.75rem; color: var(--color-text-muted); font-family: var(--font-mono); word-break: break-all; margin-top: 0.25rem;">{{ auth.user?.email }}</div>
        </div>
        <button class="btn btn-danger btn-sm" style="width: 100%;" @click="logout">
          <svg fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="width: 14px; height: 14px; margin-right: 0.5rem;">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75"/>
          </svg>
          Keluar
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>
