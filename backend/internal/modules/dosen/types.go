package dosen

import "time"

// ─── Profil Dosen dari JWT ────────────────────────────────────────────────────

type DosenProfile struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// ─── Pending KRS (antrian persetujuan) ───────────────────────────────────────

type PendingKRS struct {
	ID              string     `json:"id"`
	MahasiswaID     string     `json:"mahasiswa_id"`
	NamaMahasiswa   string     `json:"nama_mahasiswa"`
	NIM             string     `json:"nim"`
	TahunAkademikID string     `json:"tahun_akademik_id"`
	KodeTa          string     `json:"kode_ta"`
	Semester        string     `json:"semester"`
	Status          string     `json:"status"`
	TotalSKS        int        `json:"total_sks"`
	CatatanMhs      string     `json:"catatan_mhs,omitempty"`
	CatatanReviewer string     `json:"catatan_reviewer,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	Detail          []KrsDetailItem `json:"detail,omitempty"`
}

type KrsDetailItem struct {
	KelasID   string `json:"kelas_id"`
	NamaKelas string `json:"nama_kelas"`
	KodeMK    string `json:"kode_mk"`
	NamaMK    string `json:"nama_mk"`
	SKS       int    `json:"sks"`
}

// ─── Request DTO ─────────────────────────────────────────────────────────────

type ApproveReq struct {
	Catatan string `json:"catatan,omitempty"`
}

type RejectReq struct {
	Catatan string `json:"catatan"`
}
