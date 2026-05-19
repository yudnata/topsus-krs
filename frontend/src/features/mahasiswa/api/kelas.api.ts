import api from '@/core/api/client'
import type { ApiResponse } from '@/core/types/api'

export interface Kelas {
  id: string; mata_kuliah_id: string; kode_mk: string; nama_mk: string
  tahun_akademik_id: string; dosen_id: string; nama_dosen: string
  nama_kelas: string; kapasitas: number; terisi: number
}

export async function getAvailableKelas(taId?: string) {
  const { data } = await api.get<ApiResponse<Kelas[]>>('/api/admin/kelas', {
    params: taId ? { tahun_akademik_id: taId } : {}
  })
  return data
}

export interface TahunAkademik {
  id: string; kode_ta: string; semester: string; is_active: boolean
}
export async function getTahunAkademikAktif() {
  const { data } = await api.get<ApiResponse<TahunAkademik>>('/api/admin/tahun-akademik/aktif')
  return data
}
