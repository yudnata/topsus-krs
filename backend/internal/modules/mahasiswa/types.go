package mahasiswa

import "time"

// ─── Profil Mahasiswa dari JWT ────────────────────────────────────────────────

type MahasiswaProfile struct {
	ID                string  `json:"id"`
	Nama              string  `json:"nama"`
	NIM               string  `json:"nim"`
	MaxSKS            int     `json:"max_sks"`
	DosenPembimbingID *string `json:"dosen_pembimbing_id,omitempty"`
}

// ─── KRS Header ──────────────────────────────────────────────────────────────

type KrsHeader struct {
	ID              string     `json:"id"`
	MahasiswaID     string     `json:"mahasiswa_id"`
	NamaMahasiswa   string     `json:"nama_mahasiswa"`
	TahunAkademikID string     `json:"tahun_akademik_id"`
	KodeTa          string     `json:"kode_ta"`
	Semester        string     `json:"semester"`
	Status          string     `json:"status"`
	TotalSKS        int        `json:"total_sks"`
	CatatanMhs      string     `json:"catatan_mhs,omitempty"`
	CatatanReviewer string     `json:"catatan_reviewer,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Detail          []KrsDetail `json:"detail,omitempty"`
}

// ─── KRS Detail (baris kelas) ─────────────────────────────────────────────────

type KrsDetail struct {
	ID          string       `json:"id"`
	KrsID       string       `json:"krs_id"`
	KelasID     string       `json:"kelas_id"`
	NamaKelas   string       `json:"nama_kelas"`
	MataKuliahID string      `json:"mata_kuliah_id"`
	KodeMK      string       `json:"kode_mk"`
	NamaMK      string       `json:"nama_mk"`
	SKS         int          `json:"sks"`
	NamaDosen   string       `json:"nama_dosen"`
	Jadwal      []JadwalInfo `json:"jadwal"`
}

type JadwalInfo struct {
	Hari       string `json:"hari"`
	JamMulai   string `json:"jam_mulai"`
	JamSelesai string `json:"jam_selesai"`
	Ruangan    string `json:"ruangan"`
}

// ─── Request DTO ─────────────────────────────────────────────────────────────

type AddClassReq struct {
	KelasID string `json:"kelas_id"`
}

type RemoveClassReq struct {
	KelasID string `json:"kelas_id"`
}

// ─── Kelas dengan Jadwal (untuk validasi bentrok) ────────────────────────────

type KelasDetail struct {
	ID            string
	MataKuliahID  string
	SKS           int
	Kapasitas     int
	Terisi        int
	Jadwal        []JadwalRaw
}

type JadwalRaw struct {
	Hari       string
	JamMulai   string // "HH:MM:SS"
	JamSelesai string
}
