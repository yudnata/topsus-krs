export interface MahasiswaProfile {
  id: string;
  nama: string;
  nim: string;
  max_sks: number;
  dosen_pembimbing_id?: string | null;
}

export interface KrsHeader {
  id: string;
  mahasiswa_id: string;
  nama_mahasiswa: string;
  tahun_akademik_id: string;
  kode_ta: string;
  semester: string;
  status: string;
  total_sks: number;
  catatan_mhs?: string;
  catatan_reviewer?: string;
  reviewed_at?: Date | null;
  created_at: Date;
  updated_at: Date;
  detail?: KrsDetail[];
}

export interface KrsDetail {
  id: string;
  krs_id: string;
  kelas_id: string;
  nama_kelas: string;
  mata_kuliah_id: string;
  kode_mk: string;
  nama_mk: string;
  sks: number;
  nama_dosen: string;
  jadwal: JadwalInfo[];
}

export interface JadwalInfo {
  hari: string;
  jam_mulai: string;
  jam_selesai: string;
  ruangan: string;
}

export interface AddClassReq {
  kelas_id: string;
}

export interface RemoveClassReq {
  kelas_id: string;
}

export interface KelasDetail {
  id: string;
  mata_kuliah_id: string;
  sks: number;
  kapasitas: number;
  terisi: number;
  jadwal: JadwalRaw[];
}

export interface JadwalRaw {
  hari: string;
  jam_mulai: string;
  jam_selesai: string;
}
