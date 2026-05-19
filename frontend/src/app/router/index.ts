import { createRouter, createWebHistory } from 'vue-router'
import { adminRoutes } from '@/features/admin/routes'
import { authRoutes } from '@/features/auth/routes'
import { dosenRoutes } from '@/features/dosen/routes'
import { homeRoutes } from '@/features/home/routes'
import { mahasiswaRoutes } from '@/features/mahasiswa/routes'
import { staffRoutes } from '@/features/staff/routes'
import { setupRouterGuards } from '@/app/router/guards'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    ...homeRoutes,
    ...authRoutes,
    ...mahasiswaRoutes,
    ...dosenRoutes,
    ...staffRoutes,
    ...adminRoutes,
  ],
})

setupRouterGuards(router)

export default router
