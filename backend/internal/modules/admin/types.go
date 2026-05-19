package admin

import "time"

// ─── Program Studi ──────────────────────────────────────────────────────────

type ProgramStudi struct {
	ID        string `json:"id"`
	KodeProdi string `json:"kode_prodi"`
	NamaProdi string `json:"nama_prodi"`
	Fakultas  string `json:"fakultas"`
}

// ─── Tahun Akademik ──────────────────────────────────────────────────────────

type TahunAkademik struct {
	ID       string `json:"id"`
	KodeTa   string `json:"kode_ta"`
	Semester string `json:"semester"`
	IsActive bool   `json:"is_active"`
}

// ─── Mata Kuliah ─────────────────────────────────────────────────────────────

type MataKuliah struct {
	ID       string `json:"id"`
	ProdiID  string `json:"prodi_id"`
	KodeMK   string `json:"kode_mk"`
	NamaMK   string `json:"nama_mk"`
	SKS      int    `json:"sks"`
}

type CreateMataKuliahReq struct {
	ProdiID string `json:"prodi_id"`
	KodeMK  string `json:"kode_mk"`
	NamaMK  string `json:"nama_mk"`
	SKS     int    `json:"sks"`
}

// ─── Kelas ────────────────────────────────────────────────────────────────────

type Kelas struct {
	ID              string `json:"id"`
	MataKuliahID    string `json:"mata_kuliah_id"`
	NamaMK          string `json:"nama_mk"`
	KodeMK          string `json:"kode_mk"`
	TahunAkademikID string `json:"tahun_akademik_id"`
	DosenID         string `json:"dosen_id"`
	NamaDosen       string `json:"nama_dosen"`
	NamaKelas       string `json:"nama_kelas"`
	Kapasitas       int    `json:"kapasitas"`
	Terisi          int    `json:"terisi"`
	Jadwal          []JadwalKelas `json:"jadwal,omitempty"`
}

type JadwalKelas struct {
	ID         string `json:"id"`
	KelasID    string `json:"kelas_id"`
	Hari       string `json:"hari"`
	JamMulai   string `json:"jam_mulai"`
	JamSelesai string `json:"jam_selesai"`
	Ruangan    string `json:"ruangan"`
}

type CreateKelasReq struct {
	MataKuliahID    string `json:"mata_kuliah_id"`
	TahunAkademikID string `json:"tahun_akademik_id"`
	DosenID         string `json:"dosen_id"`
	NamaKelas       string `json:"nama_kelas"`
	Kapasitas       int    `json:"kapasitas"`
}

// ─── Dosen ───────────────────────────────────────────────────────────────────

type Dosen struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"user_id,omitempty"`
	NIP       string    `json:"nip"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDosenReq struct {
	NIP      string `json:"nip"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"` // buat akun user sekaligus
}

// ─── Mahasiswa ────────────────────────────────────────────────────────────────

type Mahasiswa struct {
	ID                 string    `json:"id"`
	UserID             *string   `json:"user_id,omitempty"`
	ProdiID            string    `json:"prodi_id"`
	NamProdi           string    `json:"nama_prodi"`
	NIM                string    `json:"nim"`
	Nama               string    `json:"nama"`
	DosenPembimbingID  *string   `json:"dosen_pembimbing_id,omitempty"`
	NamaDosen          string    `json:"nama_dosen,omitempty"`
	MaxSKS             int       `json:"max_sks"`
	CreatedAt          time.Time `json:"created_at"`
}

type CreateMahasiswaReq struct {
	ProdiID           string  `json:"prodi_id"`
	NIM               string  `json:"nim"`
	Nama              string  `json:"nama"`
	DosenPembimbingID *string `json:"dosen_pembimbing_id,omitempty"`
	MaxSKS            int     `json:"max_sks"`
	Password          string  `json:"password"` // buat akun user sekaligus
	Email             string  `json:"email"`
}

type PatchDPAReq struct {
	DosenPembimbingID string `json:"dosen_pembimbing_id"`
}
