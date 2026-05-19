import type { RouteRecordRaw } from 'vue-router'

export const staffRoutes: RouteRecordRaw[] = [
  {
    path: '/staff/persetujuan',
    name: 'staff-persetujuan',
    // Reuse view yang sama — komponen sudah role-aware (STAFF lihat semua)
    component: () => import('@/features/dosen/views/PersetujuanView.vue'),
    meta: { requiresAuth: true, roles: ['STAFF', 'ADMIN'] },
  },
]
