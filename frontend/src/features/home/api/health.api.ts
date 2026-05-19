import api from '@/core/api/client'
import type { ApiResponse } from '@/core/types/api'

export type HealthPayload = {
  status: string
  service: string
}

export async function fetchHealth() {
  const { data } = await api.get<ApiResponse<HealthPayload>>('/api/health')
  return data
}
