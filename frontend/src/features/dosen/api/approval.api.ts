import api from '@/core/api/client'
import type { ApiResponse } from '@/core/types/api'

export interface KrsDetailItem {
  kelas_id: string; nama_kelas: string; kode_mk: string; nama_mk: string; sks: number
}
export interface PendingKRS {
  id: string; mahasiswa_id: string; nama_mahasiswa: string; nim: string
  tahun_akademik_id: string; kode_ta: string; semester: string
  status: string; total_sks: number
  catatan_mhs?: string; catatan_reviewer?: string; reviewed_at?: string
  created_at: string; detail?: KrsDetailItem[]
}

export async function getPendingKRS() {
  const { data } = await api.get<ApiResponse<PendingKRS[]>>('/api/approval/pending')
  return data
}

export async function approveKRS(krsId: string, catatan?: string) {
  const { data } = await api.post<ApiResponse<PendingKRS>>(`/api/approval/${krsId}/approve`, { catatan: catatan ?? '' })
  return data
}

export async function rejectKRS(krsId: string, catatan: string) {
  const { data } = await api.post<ApiResponse<PendingKRS>>(`/api/approval/${krsId}/reject`, { catatan })
  return data
}
