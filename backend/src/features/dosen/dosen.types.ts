export interface DosenProfile {
  id: string;
  nama: string;
}

export interface PendingKRS {
  id: string;
  mahasiswa_id: string;
  nama_mahasiswa: string;
  nim: string;
  tahun_akademik_id: string;
  kode_ta: string;
  semester: string;
  status: string;
  total_sks: number;
  catatan_mhs?: string;
  catatan_reviewer?: string;
  reviewed_at?: Date | null;
  created_at: Date;
  detail?: KrsDetailItem[];
}

export interface KrsDetailItem {
  kelas_id: string;
  nama_kelas: string;
  kode_mk: string;
  nama_mk: string;
  sks: number;
}

export interface ApproveReq {
  catatan?: string;
}

export interface RejectReq {
  catatan: string;
}
