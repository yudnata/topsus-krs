package dosen

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ─── Profil Dosen ─────────────────────────────────────────────────────────────

func (r *Repository) GetProfileByUserID(ctx context.Context, userID string) (*DosenProfile, error) {
	var d DosenProfile
	err := r.db.QueryRow(ctx, `SELECT id, nama FROM dosen WHERE user_id = $1`, userID).
		Scan(&d.ID, &d.Nama)
	return &d, err
}

// ─── List Pending ─────────────────────────────────────────────────────────────

// ListPendingForDosen — hanya mahasiswa bimbingan DPA dosen ini.
func (r *Repository) ListPendingForDosen(ctx context.Context, dosenID string) ([]PendingKRS, error) {
	return r.queryPending(ctx, `
		SELECT pk.id, pk.mahasiswa_id, m.nama, m.nim,
		       pk.tahun_akademik_id, ta.kode_ta, ta.semester,
		       pk.status, pk.total_sks,
		       COALESCE(pk.catatan_mhs,''), COALESCE(pk.catatan_reviewer,''),
		       pk.reviewed_at, pk.created_at
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
		WHERE pk.status = 'diajukan' AND m.dosen_pembimbing_id = $1
		ORDER BY pk.created_at ASC`, dosenID)
}

// ListPendingAll — semua KRS status diajukan (untuk STAFF).
func (r *Repository) ListPendingAll(ctx context.Context) ([]PendingKRS, error) {
	return r.queryPending(ctx, `
		SELECT pk.id, pk.mahasiswa_id, m.nama, m.nim,
		       pk.tahun_akademik_id, ta.kode_ta, ta.semester,
		       pk.status, pk.total_sks,
		       COALESCE(pk.catatan_mhs,''), COALESCE(pk.catatan_reviewer,''),
		       pk.reviewed_at, pk.created_at
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
		WHERE pk.status = 'diajukan'
		ORDER BY pk.created_at ASC`)
}

func (r *Repository) queryPending(ctx context.Context, q string, args ...any) ([]PendingKRS, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []PendingKRS
	for rows.Next() {
		var p PendingKRS
		var rat *time.Time
		if err := rows.Scan(&p.ID, &p.MahasiswaID, &p.NamaMahasiswa, &p.NIM,
			&p.TahunAkademikID, &p.KodeTa, &p.Semester,
			&p.Status, &p.TotalSKS,
			&p.CatatanMhs, &p.CatatanReviewer,
			&rat, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.ReviewedAt = rat
		// Lampirkan detail ringkas
		p.Detail, _ = r.getDetailRingkas(ctx, p.ID)
		result = append(result, p)
	}
	return result, nil
}

func (r *Repository) getDetailRingkas(ctx context.Context, krsID string) ([]KrsDetailItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pkd.kelas_id, k.nama_kelas, mk.kode_mk, mk.nama_mk, mk.sks
		FROM pengajuan_krs_detail pkd
		JOIN kelas k ON k.id = pkd.kelas_id
		JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
		WHERE pkd.pengajuan_krs_id = $1
		ORDER BY mk.kode_mk`, krsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []KrsDetailItem
	for rows.Next() {
		var d KrsDetailItem
		if err := rows.Scan(&d.KelasID, &d.NamaKelas, &d.KodeMK, &d.NamaMK, &d.SKS); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

// ─── Ambil KRS ───────────────────────────────────────────────────────────────

func (r *Repository) GetKrs(ctx context.Context, krsID string) (*PendingKRS, error) {
	var p PendingKRS
	var rat *time.Time
	err := r.db.QueryRow(ctx, `
		SELECT pk.id, pk.mahasiswa_id, m.nama, m.nim,
		       pk.tahun_akademik_id, ta.kode_ta, ta.semester,
		       pk.status, pk.total_sks,
		       COALESCE(pk.catatan_mhs,''), COALESCE(pk.catatan_reviewer,''),
		       pk.reviewed_at, pk.created_at
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
		WHERE pk.id = $1`, krsID,
	).Scan(&p.ID, &p.MahasiswaID, &p.NamaMahasiswa, &p.NIM,
		&p.TahunAkademikID, &p.KodeTa, &p.Semester,
		&p.Status, &p.TotalSKS,
		&p.CatatanMhs, &p.CatatanReviewer,
		&rat, &p.CreatedAt)
	p.ReviewedAt = rat
	return &p, err
}

// GetDosenPembimbingOfKrs — ambil dosen_pembimbing_id dari mahasiswa pemilik KRS.
func (r *Repository) GetDosenPembimbingOfKrs(ctx context.Context, krsID string) (string, string, error) {
	var dosenID, mahasiswaID string
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(m.dosen_pembimbing_id::text,''), m.id
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		WHERE pk.id = $1`, krsID,
	).Scan(&dosenID, &mahasiswaID)
	return dosenID, mahasiswaID, err
}

// ─── Approve / Reject ─────────────────────────────────────────────────────────

func (r *Repository) Approve(ctx context.Context, krsID, reviewerID, catatan string) error {
	res, err := r.db.Exec(ctx, `
		UPDATE pengajuan_krs
		SET status = 'disetujui', reviewed_by = $1, reviewed_at = NOW(),
		    catatan_reviewer = $2, updated_at = NOW()
		WHERE id = $3 AND status = 'diajukan'`,
		reviewerID, catatan, krsID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errNotPending
	}
	return nil
}

func (r *Repository) Reject(ctx context.Context, krsID, reviewerID, catatan string) error {
	res, err := r.db.Exec(ctx, `
		UPDATE pengajuan_krs
		SET status = 'ditolak', reviewed_by = $1, reviewed_at = NOW(),
		    catatan_reviewer = $2, updated_at = NOW()
		WHERE id = $3 AND status = 'diajukan'`,
		reviewerID, catatan, krsID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errNotPending
	}
	return nil
}

var errNotPending = &repoError{"KRS tidak dalam status diajukan atau tidak ditemukan"}

type repoError struct{ msg string }

func (e *repoError) Error() string { return e.msg }
