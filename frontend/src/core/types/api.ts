export type ApiResponse<T = unknown> = {
  success: boolean
  message: string
  data?: T
}

export type UserRole = 'ADMIN' | 'MAHASISWA' | 'DOSEN' | 'STAFF'
