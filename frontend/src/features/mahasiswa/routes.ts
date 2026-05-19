import type { RouteRecordRaw } from 'vue-router'

export const mahasiswaRoutes: RouteRecordRaw[] = [
  {
    path: '/mahasiswa/krs',
    name: 'mahasiswa-krs',
    component: () => import('@/features/mahasiswa/views/KrsView.vue'),
    meta: { requiresAuth: true, roles: ['MAHASISWA'] },
  },
  {
    path: '/mahasiswa/history',
    name: 'mahasiswa-history',
    component: () => import('@/features/mahasiswa/views/KrsHistoryView.vue'),
    meta: { requiresAuth: true, roles: ['MAHASISWA'] },
  },
]
