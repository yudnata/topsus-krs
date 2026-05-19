import api from '@/core/api/client'
import type { ApiResponse } from '@/core/types/api'

export interface JadwalInfo {
  hari: string; jam_mulai: string; jam_selesai: string; ruangan: string
}
export interface KrsDetail {
  id: string; kelas_id: string; nama_kelas: string
  kode_mk: string; nama_mk: string; sks: number
  nama_dosen: string; jadwal: JadwalInfo[]
}
export interface KrsHeader {
  id: string; mahasiswa_id: string; nama_mahasiswa: string
  tahun_akademik_id: string; kode_ta: string; semester: string
  status: 'draft' | 'diajukan' | 'disetujui' | 'ditolak'
  total_sks: number; catatan_mhs?: string; catatan_reviewer?: string
  reviewed_at?: string; created_at: string; updated_at: string
  detail?: KrsDetail[]
}

export async function getCurrentKRS() {
  const { data } = await api.get<ApiResponse<KrsHeader>>('/api/krs/current')
  return data
}

export async function addKelasToKRS(krsId: string, kelasId: string) {
  const { data } = await api.post<ApiResponse<KrsHeader>>(`/api/krs/${krsId}/add-class`, { kelas_id: kelasId })
  return data
}

export async function removeKelasFromKRS(krsId: string, kelasId: string) {
  const { data } = await api.delete<ApiResponse<KrsHeader>>(`/api/krs/${krsId}/remove-class`, {
    data: { kelas_id: kelasId }
  })
  return data
}

export async function submitKRS(krsId: string) {
  const { data } = await api.post<ApiResponse<KrsHeader>>(`/api/krs/${krsId}/submit`)
  return data
}

export async function getKRSHistory() {
  const { data } = await api.get<ApiResponse<KrsHeader[]>>('/api/krs/history')
  return data
}
