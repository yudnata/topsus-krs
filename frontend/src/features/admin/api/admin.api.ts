import api from '@/core/api/client'
import type { ApiResponse } from '@/core/types/api'

// ─── Types ───────────────────────────────────────────────────────────────────

export interface Dosen {
  id: string; nip: string; nama: string; email: string; user_id?: string
}
export interface Mahasiswa {
  id: string; nim: string; nama: string; prodi_id: string; nama_prodi: string
  dosen_pembimbing_id?: string; nama_dosen?: string; max_sks: number
}
export interface MataKuliah {
  id: string; prodi_id: string; kode_mk: string; nama_mk: string; sks: number
}
export interface Kelas {
  id: string; mata_kuliah_id: string; kode_mk: string; nama_mk: string
  tahun_akademik_id: string; dosen_id: string; nama_dosen: string
  nama_kelas: string; kapasitas: number; terisi: number
}
export interface ProgramStudi { id: string; kode_prodi: string; nama_prodi: string; fakultas: string }
export interface TahunAkademik { id: string; kode_ta: string; semester: string; is_active: boolean }

// ─── API calls ───────────────────────────────────────────────────────────────

export const listDosen         = () => api.get<ApiResponse<Dosen[]>>('/api/admin/dosen').then(r => r.data)
export const createDosen       = (d: { nip:string;nama:string;email:string;password:string }) =>
  api.post<ApiResponse<Dosen>>('/api/admin/dosen', d).then(r => r.data)

export const listMahasiswa     = () => api.get<ApiResponse<Mahasiswa[]>>('/api/admin/mahasiswa').then(r => r.data)
export const createMahasiswa   = (m: { nim:string;nama:string;prodi_id:string;email:string;password:string;dosen_pembimbing_id?:string }) =>
  api.post<ApiResponse<Mahasiswa>>('/api/admin/mahasiswa', m).then(r => r.data)
export const patchDPA          = (mhsId: string, dosenId: string) =>
  api.patch<ApiResponse<null>>(`/api/admin/mahasiswa/${mhsId}/dpa`, { dosen_pembimbing_id: dosenId }).then(r => r.data)

export const listMataKuliah    = () => api.get<ApiResponse<MataKuliah[]>>('/api/admin/mata-kuliah').then(r => r.data)
export const createMataKuliah  = (m: { prodi_id:string;kode_mk:string;nama_mk:string;sks:number }) =>
  api.post<ApiResponse<MataKuliah>>('/api/admin/mata-kuliah', m).then(r => r.data)

export const listKelas         = (taId?: string) =>
  api.get<ApiResponse<Kelas[]>>('/api/admin/kelas', { params: taId ? { tahun_akademik_id: taId } : {} }).then(r => r.data)
export const createKelas       = (k: { mata_kuliah_id:string;tahun_akademik_id:string;dosen_id:string;nama_kelas:string;kapasitas:number }) =>
  api.post<ApiResponse<Kelas>>('/api/admin/kelas', k).then(r => r.data)

export const listProdi         = () => api.get<ApiResponse<ProgramStudi[]>>('/api/admin/prodi').then(r => r.data)
export const listTahunAkademik = () => api.get<ApiResponse<TahunAkademik[]>>('/api/admin/tahun-akademik').then(r => r.data)
