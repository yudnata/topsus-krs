package mahasiswa

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

// ─── Profil Mahasiswa ─────────────────────────────────────────────────────────

// GetProfileByUserID — cari data mahasiswa berdasarkan user_id dari JWT.
func (r *Repository) GetProfileByUserID(ctx context.Context, userID string) (*MahasiswaProfile, error) {
	var m MahasiswaProfile
	var dpID *string
	err := r.db.QueryRow(ctx, `
		SELECT id, nama, nim, max_sks, dosen_pembimbing_id
		FROM mahasiswa WHERE user_id = $1`, userID,
	).Scan(&m.ID, &m.Nama, &m.NIM, &m.MaxSKS, &dpID)
	m.DosenPembimbingID = dpID
	return &m, err
}

// ─── Tahun Akademik Aktif ─────────────────────────────────────────────────────

func (r *Repository) GetTahunAkademikAktif(ctx context.Context) (id, kodeTa, semester string, err error) {
	err = r.db.QueryRow(ctx, `SELECT id, kode_ta, semester FROM tahun_akademik WHERE is_active = true LIMIT 1`).
		Scan(&id, &kodeTa, &semester)
	return
}

// ─── KRS Header ──────────────────────────────────────────────────────────────

// GetOrCreateDraft — ambil atau buat draft KRS mahasiswa untuk TA aktif.
func (r *Repository) GetOrCreateDraft(ctx context.Context, mahasiswaID, taID string) (*KrsHeader, error) {
	var h KrsHeader
	var rat *time.Time
	err := r.db.QueryRow(ctx, `
		SELECT pk.id, pk.mahasiswa_id, m.nama,
		       pk.tahun_akademik_id, ta.kode_ta, ta.semester,
		       pk.status, pk.total_sks,
		       COALESCE(pk.catatan_mhs,''), COALESCE(pk.catatan_reviewer,''),
		       pk.reviewed_at, pk.created_at, pk.updated_at
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
		WHERE pk.mahasiswa_id = $1 AND pk.tahun_akademik_id = $2`,
		mahasiswaID, taID,
	).Scan(&h.ID, &h.MahasiswaID, &h.NamaMahasiswa,
		&h.TahunAkademikID, &h.KodeTa, &h.Semester,
		&h.Status, &h.TotalSKS,
		&h.CatatanMhs, &h.CatatanReviewer,
		&rat, &h.CreatedAt, &h.UpdatedAt)

	if err == nil {
		h.ReviewedAt = rat
		return &h, nil
	}

	// Belum ada — buat draft baru
	now := time.Now()
	err = r.db.QueryRow(ctx, `
		INSERT INTO pengajuan_krs (mahasiswa_id, tahun_akademik_id, status, total_sks, created_at, updated_at)
		VALUES ($1, $2, 'draft', 0, $3, $3)
		RETURNING id, mahasiswa_id, tahun_akademik_id, status, total_sks, created_at, updated_at`,
		mahasiswaID, taID, now,
	).Scan(&h.ID, &h.MahasiswaID, &h.TahunAkademikID, &h.Status, &h.TotalSKS, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, err
	}
	h.KodeTa = ""
	h.Semester = ""
	return &h, nil
}

func (r *Repository) GetKrsByID(ctx context.Context, krsID string) (*KrsHeader, error) {
	var h KrsHeader
	var rat *time.Time
	err := r.db.QueryRow(ctx, `
		SELECT pk.id, pk.mahasiswa_id, m.nama,
		       pk.tahun_akademik_id, ta.kode_ta, ta.semester,
		       pk.status, pk.total_sks,
		       COALESCE(pk.catatan_mhs,''), COALESCE(pk.catatan_reviewer,''),
		       pk.reviewed_at, pk.created_at, pk.updated_at
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
		WHERE pk.id = $1`, krsID,
	).Scan(&h.ID, &h.MahasiswaID, &h.NamaMahasiswa,
		&h.TahunAkademikID, &h.KodeTa, &h.Semester,
		&h.Status, &h.TotalSKS,
		&h.CatatanMhs, &h.CatatanReviewer,
		&rat, &h.CreatedAt, &h.UpdatedAt)
	h.ReviewedAt = rat
	return &h, err
}

// ─── KRS Detail ──────────────────────────────────────────────────────────────

func (r *Repository) GetDetails(ctx context.Context, krsID string) ([]KrsDetail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pkd.id, pkd.pengajuan_krs_id, pkd.kelas_id,
		       k.nama_kelas, mk.id, mk.kode_mk, mk.nama_mk, mk.sks,
		       COALESCE(d.nama,'')
		FROM pengajuan_krs_detail pkd
		JOIN kelas k ON k.id = pkd.kelas_id
		JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
		LEFT JOIN dosen d ON d.id = k.dosen_id
		WHERE pkd.pengajuan_krs_id = $1
		ORDER BY mk.kode_mk`, krsID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []KrsDetail
	for rows.Next() {
		var dt KrsDetail
		if err := rows.Scan(&dt.ID, &dt.KrsID, &dt.KelasID,
			&dt.NamaKelas, &dt.MataKuliahID, &dt.KodeMK, &dt.NamaMK, &dt.SKS,
			&dt.NamaDosen); err != nil {
			return nil, err
		}
		// Ambil jadwal untuk kelas ini
		jadwal, _ := r.getJadwalForKelas(ctx, dt.KelasID)
		dt.Jadwal = jadwal
		result = append(result, dt)
	}
	return result, nil
}

func (r *Repository) getJadwalForKelas(ctx context.Context, kelasID string) ([]JadwalInfo, error) {
	rows, err := r.db.Query(ctx, `
		SELECT hari, jam_mulai::text, jam_selesai::text, ruangan
		FROM jadwal_kelas WHERE kelas_id = $1`, kelasID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jadwal []JadwalInfo
	for rows.Next() {
		var j JadwalInfo
		if err := rows.Scan(&j.Hari, &j.JamMulai, &j.JamSelesai, &j.Ruangan); err != nil {
			return nil, err
		}
		jadwal = append(jadwal, j)
	}
	return jadwal, nil
}

// ─── Kelas Detail (untuk validasi) ───────────────────────────────────────────

func (r *Repository) GetKelasDetail(ctx context.Context, kelasID string) (*KelasDetail, error) {
	var k KelasDetail
	err := r.db.QueryRow(ctx, `
		SELECT k.id, k.mata_kuliah_id, mk.sks, k.kapasitas, k.terisi
		FROM kelas k
		JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
		WHERE k.id = $1`, kelasID,
	).Scan(&k.ID, &k.MataKuliahID, &k.SKS, &k.Kapasitas, &k.Terisi)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `
		SELECT hari, jam_mulai::text, jam_selesai::text FROM jadwal_kelas WHERE kelas_id = $1`, kelasID)
	if err != nil {
		return &k, nil
	}
	defer rows.Close()
	for rows.Next() {
		var j JadwalRaw
		_ = rows.Scan(&j.Hari, &j.JamMulai, &j.JamSelesai)
		k.Jadwal = append(k.Jadwal, j)
	}
	return &k, nil
}

// GetCurrentKelasIDs — ID kelas yang sudah ada di KRS ini (cegah duplikat MK).
func (r *Repository) GetCurrentKelasInKrs(ctx context.Context, krsID string) ([]KelasDetail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT k.id, k.mata_kuliah_id, mk.sks, k.kapasitas, k.terisi
		FROM pengajuan_krs_detail pkd
		JOIN kelas k ON k.id = pkd.kelas_id
		JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
		WHERE pkd.pengajuan_krs_id = $1`, krsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []KelasDetail
	for rows.Next() {
		var k KelasDetail
		if err := rows.Scan(&k.ID, &k.MataKuliahID, &k.SKS, &k.Kapasitas, &k.Terisi); err != nil {
			return nil, err
		}
		// Ambil jadwal
		jrows, _ := r.db.Query(ctx, `
			SELECT hari, jam_mulai::text, jam_selesai::text FROM jadwal_kelas WHERE kelas_id = $1`, k.ID)
		if jrows != nil {
			for jrows.Next() {
				var j JadwalRaw
				_ = jrows.Scan(&j.Hari, &j.JamMulai, &j.JamSelesai)
				k.Jadwal = append(k.Jadwal, j)
			}
			jrows.Close()
		}
		result = append(result, k)
	}
	return result, nil
}

// ─── Mutasi KRS ───────────────────────────────────────────────────────────────

func (r *Repository) AddKelas(ctx context.Context, krsID, kelasID string, sks int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert detail
	_, err = tx.Exec(ctx, `
		INSERT INTO pengajuan_krs_detail (pengajuan_krs_id, kelas_id)
		VALUES ($1, $2)`, krsID, kelasID)
	if err != nil {
		return err
	}
	// Increment terisi
	_, err = tx.Exec(ctx, `UPDATE kelas SET terisi = terisi + 1 WHERE id = $1`, kelasID)
	if err != nil {
		return err
	}
	// Update total_sks & updated_at
	_, err = tx.Exec(ctx, `
		UPDATE pengajuan_krs SET total_sks = total_sks + $1, updated_at = NOW() WHERE id = $2`,
		sks, krsID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) RemoveKelas(ctx context.Context, krsID, kelasID string, sks int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, `
		DELETE FROM pengajuan_krs_detail
		WHERE pengajuan_krs_id = $1 AND kelas_id = $2`, krsID, kelasID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return tx.Rollback(ctx)
	}
	_, err = tx.Exec(ctx, `UPDATE kelas SET terisi = GREATEST(terisi - 1, 0) WHERE id = $1`, kelasID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		UPDATE pengajuan_krs SET total_sks = GREATEST(total_sks - $1, 0), updated_at = NOW() WHERE id = $2`,
		sks, krsID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) SubmitKrs(ctx context.Context, krsID string) error {
	res, err := r.db.Exec(ctx, `
		UPDATE pengajuan_krs SET status = 'diajukan', updated_at = NOW()
		WHERE id = $1 AND status = 'draft'`, krsID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errKrsNotDraft
	}
	return nil
}

var errKrsNotDraft = &businessError{"KRS tidak dalam status draft atau tidak ditemukan"}

// ─── History ─────────────────────────────────────────────────────────────────

func (r *Repository) GetHistory(ctx context.Context, mahasiswaID string) ([]KrsHeader, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pk.id, pk.mahasiswa_id, m.nama,
		       pk.tahun_akademik_id, ta.kode_ta, ta.semester,
		       pk.status, pk.total_sks,
		       COALESCE(pk.catatan_mhs,''), COALESCE(pk.catatan_reviewer,''),
		       pk.reviewed_at, pk.created_at, pk.updated_at
		FROM pengajuan_krs pk
		JOIN mahasiswa m ON m.id = pk.mahasiswa_id
		JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
		WHERE pk.mahasiswa_id = $1
		ORDER BY pk.created_at DESC`, mahasiswaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []KrsHeader
	for rows.Next() {
		var h KrsHeader
		var rat *time.Time
		if err := rows.Scan(&h.ID, &h.MahasiswaID, &h.NamaMahasiswa,
			&h.TahunAkademikID, &h.KodeTa, &h.Semester,
			&h.Status, &h.TotalSKS,
			&h.CatatanMhs, &h.CatatanReviewer,
			&rat, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		h.ReviewedAt = rat
		result = append(result, h)
	}
	return result, nil
}

// ─── Util ─────────────────────────────────────────────────────────────────────

type businessError struct{ msg string }

func (e *businessError) Error() string { return e.msg }
