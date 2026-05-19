import type { Router } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/auth.store'
import type { UserRole } from '@/core/types/api'

export function setupRouterGuards(router: Router) {
  router.beforeEach(async (to) => {
    const auth = useAuthStore()

    if (auth.token && !auth.user) {
      try {
        await auth.loadProfile()
      } catch {
        auth.logout()
      }
    }

    if (to.meta.guest && auth.isAuthenticated) {
      return { name: 'home' }
    }

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }

    const roles = to.meta.roles as UserRole[] | undefined
    if (roles?.length && auth.role && !roles.includes(auth.role)) {
      return { name: 'home' }
    }

    return true
  })
}
