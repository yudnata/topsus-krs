import type { RouteRecordRaw } from 'vue-router'

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin/dashboard',
    name: 'admin-dashboard',
    component: () => import('@/features/admin/views/DashboardView.vue'),
    meta: { requiresAuth: true, roles: ['ADMIN'] },
  },
  {
    path: '/admin/mahasiswa',
    name: 'admin-mahasiswa',
    component: () => import('@/features/admin/views/MahasiswaView.vue'),
    meta: { requiresAuth: true, roles: ['ADMIN'] },
  },
  {
    path: '/admin/dosen',
    name: 'admin-dosen',
    component: () => import('@/features/admin/views/DosenView.vue'),
    meta: { requiresAuth: true, roles: ['ADMIN'] },
  },
  {
    path: '/admin/mata-kuliah',
    name: 'admin-mata-kuliah',
    component: () => import('@/features/admin/views/MataKuliahView.vue'),
    meta: { requiresAuth: true, roles: ['ADMIN'] },
  },
  {
    path: '/admin/kelas',
    name: 'admin-kelas',
    component: () => import('@/features/admin/views/KelasView.vue'),
    meta: { requiresAuth: true, roles: ['ADMIN'] },
  },
]
