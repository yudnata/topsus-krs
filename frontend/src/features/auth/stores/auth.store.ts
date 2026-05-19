import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { AuthUser, LoginPayload } from '@/features/auth/api/auth.api'
import { fetchProfile, login as loginApi } from '@/features/auth/api/auth.api'
import type { UserRole } from '@/core/types/api'

const TOKEN_KEY = 'krs_token'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(null)
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const role = computed(() => user.value?.role ?? null)

  async function login(payload: LoginPayload) {
    loading.value = true
    try {
      const res = await loginApi(payload)
      if (!res.success || !res.data) {
        throw new Error(res.message || 'Login gagal')
      }
      token.value = res.data.token
      user.value = res.data.user
      localStorage.setItem(TOKEN_KEY, res.data.token)
    } finally {
      loading.value = false
    }
  }

  async function loadProfile() {
    if (!token.value) return
    const res = await fetchProfile()
    if (res.success && res.data) {
      user.value = res.data
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  function hasRole(...roles: UserRole[]) {
    return !!role.value && roles.includes(role.value)
  }

  return {
    user,
    token,
    loading,
    isAuthenticated,
    role,
    login,
    loadProfile,
    logout,
    hasRole,
  }
})
