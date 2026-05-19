import api from '@/core/api/client'
import type { ApiResponse } from '@/core/types/api'
import type { UserRole } from '@/core/types/api'

export type AuthUser = {
  id: string
  email: string
  role: UserRole
  is_active: boolean
}

export type LoginPayload = { email: string; password: string }
export type LoginResult = { token: string; user: AuthUser }

export async function login(payload: LoginPayload) {
  const { data } = await api.post<ApiResponse<LoginResult>>('/api/auth/login', payload)
  return data
}

export async function fetchProfile() {
  const { data } = await api.get<ApiResponse<AuthUser>>('/api/auth/profile')
  return data
}
