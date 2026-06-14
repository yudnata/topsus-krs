export interface ProgramStudi {
  id: string;
  kode_prodi: string;
  nama_prodi: string;
  fakultas: string;
}

export interface TahunAkademik {
  id: string;
  kode_ta: string;
  semester: string;
  is_active: boolean;
}

export interface MataKuliah {
  id: string;
  prodi_id: string;
  kode_mk: string;
  nama_mk: string;
  sks: number;
}

export interface CreateMataKuliahReq {
  prodi_id: string;
  kode_mk: string;
  nama_mk: string;
  sks: number;
}

export interface Kelas {
  id: string;
  mata_kuliah_id: string;
  nama_mk: string;
  kode_mk: string;
  tahun_akademik_id: string;
  dosen_id: string;
  nama_dosen: string;
  nama_kelas: string;
  kapasitas: number;
  terisi: number;
}

export interface CreateKelasReq {
  mata_kuliah_id: string;
  tahun_akademik_id: string;
  dosen_id: string;
  nama_kelas: string;
  kapasitas: number;
}

export interface Dosen {
  id: string;
  user_id?: string | null;
  nip: string;
  nama: string;
  email: string;
  created_at: Date;
}

export interface CreateDosenReq {
  nip: string;
  nama: string;
  email: string;
  password: string;
}

export interface Mahasiswa {
  id: string;
  user_id?: string | null;
  prodi_id: string;
  nama_prodi: string;
  nim: string;
  nama: string;
  dosen_pembimbing_id?: string | null;
  nama_dosen?: string;
  max_sks: number;
  created_at: Date;
}

export interface CreateMahasiswaReq {
  prodi_id: string;
  nim: string;
  nama: string;
  dosen_pembimbing_id?: string | null;
  max_sks?: number;
  password: string;
  email: string;
}

export interface PatchDPAReq {
  dosen_pembimbing_id: string;
}
