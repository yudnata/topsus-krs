import type { RouteRecordRaw } from 'vue-router'

export const dosenRoutes: RouteRecordRaw[] = [
  {
    path: '/dosen/persetujuan',
    name: 'dosen-persetujuan',
    component: () => import('@/features/dosen/views/PersetujuanView.vue'),
    meta: { requiresAuth: true, roles: ['DOSEN'] },
  },
]
