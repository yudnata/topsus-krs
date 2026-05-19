package admin

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

// ─── Program Studi ───────────────────────────────────────────────────────────

func (r *Repository) ListProdi(ctx context.Context) ([]ProgramStudi, error) {
	rows, err := r.db.Query(ctx, `SELECT id, kode_prodi, nama_prodi, fakultas FROM program_studi ORDER BY nama_prodi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []ProgramStudi
	for rows.Next() {
		var p ProgramStudi
		if err := rows.Scan(&p.ID, &p.KodeProdi, &p.NamaProdi, &p.Fakultas); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// ─── Tahun Akademik ──────────────────────────────────────────────────────────

func (r *Repository) GetTahunAkademikAktif(ctx context.Context) (*TahunAkademik, error) {
	var ta TahunAkademik
	err := r.db.QueryRow(ctx, `SELECT id, kode_ta, semester, is_active FROM tahun_akademik WHERE is_active = true LIMIT 1`).
		Scan(&ta.ID, &ta.KodeTa, &ta.Semester, &ta.IsActive)
	if err != nil {
		return nil, err
	}
	return &ta, nil
}

func (r *Repository) ListTahunAkademik(ctx context.Context) ([]TahunAkademik, error) {
	rows, err := r.db.Query(ctx, `SELECT id, kode_ta, semester, is_active FROM tahun_akademik ORDER BY kode_ta DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []TahunAkademik
	for rows.Next() {
		var ta TahunAkademik
		if err := rows.Scan(&ta.ID, &ta.KodeTa, &ta.Semester, &ta.IsActive); err != nil {
			return nil, err
		}
		result = append(result, ta)
	}
	return result, nil
}

// ─── Mata Kuliah ─────────────────────────────────────────────────────────────

func (r *Repository) ListMataKuliah(ctx context.Context) ([]MataKuliah, error) {
	rows, err := r.db.Query(ctx, `SELECT id, COALESCE(prodi_id::text,''), kode_mk, nama_mk, sks FROM mata_kuliah ORDER BY kode_mk`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []MataKuliah
	for rows.Next() {
		var m MataKuliah
		if err := rows.Scan(&m.ID, &m.ProdiID, &m.KodeMK, &m.NamaMK, &m.SKS); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *Repository) CreateMataKuliah(ctx context.Context, req CreateMataKuliahReq) (*MataKuliah, error) {
	var m MataKuliah
	err := r.db.QueryRow(ctx, `
		INSERT INTO mata_kuliah (prodi_id, kode_mk, nama_mk, sks)
		VALUES ($1, $2, $3, $4)
		RETURNING id, COALESCE(prodi_id::text,''), kode_mk, nama_mk, sks`,
		req.ProdiID, req.KodeMK, req.NamaMK, req.SKS,
	).Scan(&m.ID, &m.ProdiID, &m.KodeMK, &m.NamaMK, &m.SKS)
	return &m, err
}

// ─── Kelas ───────────────────────────────────────────────────────────────────

func (r *Repository) ListKelas(ctx context.Context, taID string) ([]Kelas, error) {
	query := `
		SELECT k.id, k.mata_kuliah_id, mk.nama_mk, mk.kode_mk,
		       k.tahun_akademik_id, COALESCE(k.dosen_id::text,''), COALESCE(d.nama,''),
		       k.nama_kelas, k.kapasitas, k.terisi
		FROM kelas k
		JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
		LEFT JOIN dosen d ON d.id = k.dosen_id`
	args := []any{}
	if taID != "" {
		query += " WHERE k.tahun_akademik_id = $1"
		args = append(args, taID)
	}
	query += " ORDER BY mk.kode_mk, k.nama_kelas"
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Kelas
	for rows.Next() {
		var k Kelas
		if err := rows.Scan(&k.ID, &k.MataKuliahID, &k.NamaMK, &k.KodeMK,
			&k.TahunAkademikID, &k.DosenID, &k.NamaDosen,
			&k.NamaKelas, &k.Kapasitas, &k.Terisi); err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	return result, nil
}

func (r *Repository) CreateKelas(ctx context.Context, req CreateKelasReq) (*Kelas, error) {
	var k Kelas
	err := r.db.QueryRow(ctx, `
		INSERT INTO kelas (mata_kuliah_id, tahun_akademik_id, dosen_id, nama_kelas, kapasitas)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, mata_kuliah_id, tahun_akademik_id, COALESCE(dosen_id::text,''), nama_kelas, kapasitas, terisi`,
		req.MataKuliahID, req.TahunAkademikID, req.DosenID, req.NamaKelas, req.Kapasitas,
	).Scan(&k.ID, &k.MataKuliahID, &k.TahunAkademikID, &k.DosenID, &k.NamaKelas, &k.Kapasitas, &k.Terisi)
	return &k, err
}

// ─── Dosen ───────────────────────────────────────────────────────────────────

func (r *Repository) ListDosen(ctx context.Context) ([]Dosen, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, nip, nama, COALESCE(email,''), created_at
		FROM dosen ORDER BY nama`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Dosen
	for rows.Next() {
		var d Dosen
		var uid *string
		if err := rows.Scan(&d.ID, &uid, &d.NIP, &d.Nama, &d.Email, &d.CreatedAt); err != nil {
			return nil, err
		}
		d.UserID = uid
		result = append(result, d)
	}
	return result, nil
}

func (r *Repository) GetDosen(ctx context.Context, id string) (*Dosen, error) {
	var d Dosen
	var uid *string
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, nip, nama, COALESCE(email,''), created_at
		FROM dosen WHERE id = $1`, id,
	).Scan(&d.ID, &uid, &d.NIP, &d.Nama, &d.Email, &d.CreatedAt)
	d.UserID = uid
	return &d, err
}

func (r *Repository) CreateDosenWithUser(ctx context.Context, nip, nama, email, passwordHash string) (*Dosen, error) {
	now := time.Now()
	// Buat user dulu
	var userID string
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, 'DOSEN', true, $3, $3)
		RETURNING id`, email, passwordHash, now,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}
	var d Dosen
	err = r.db.QueryRow(ctx, `
		INSERT INTO dosen (user_id, nip, nama, email, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, nip, nama, email, created_at`,
		userID, nip, nama, email, now,
	).Scan(&d.ID, &d.UserID, &d.NIP, &d.Nama, &d.Email, &d.CreatedAt)
	return &d, err
}

// ─── Mahasiswa ───────────────────────────────────────────────────────────────

func (r *Repository) ListMahasiswa(ctx context.Context) ([]Mahasiswa, error) {
	rows, err := r.db.Query(ctx, `
		SELECT m.id, m.user_id, m.prodi_id, COALESCE(ps.nama_prodi,''),
		       m.nim, m.nama,
		       m.dosen_pembimbing_id, COALESCE(d.nama,''),
		       m.max_sks, m.created_at
		FROM mahasiswa m
		LEFT JOIN program_studi ps ON ps.id = m.prodi_id
		LEFT JOIN dosen d ON d.id = m.dosen_pembimbing_id
		ORDER BY m.nim`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		var uid *string
		var dpID *string
		if err := rows.Scan(&m.ID, &uid, &m.ProdiID, &m.NamProdi,
			&m.NIM, &m.Nama,
			&dpID, &m.NamaDosen,
			&m.MaxSKS, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.UserID = uid
		m.DosenPembimbingID = dpID
		result = append(result, m)
	}
	return result, nil
}

func (r *Repository) GetMahasiswa(ctx context.Context, id string) (*Mahasiswa, error) {
	var m Mahasiswa
	var uid *string
	var dpID *string
	err := r.db.QueryRow(ctx, `
		SELECT m.id, m.user_id, m.prodi_id, COALESCE(ps.nama_prodi,''),
		       m.nim, m.nama,
		       m.dosen_pembimbing_id, COALESCE(d.nama,''),
		       m.max_sks, m.created_at
		FROM mahasiswa m
		LEFT JOIN program_studi ps ON ps.id = m.prodi_id
		LEFT JOIN dosen d ON d.id = m.dosen_pembimbing_id
		WHERE m.id = $1`, id,
	).Scan(&m.ID, &uid, &m.ProdiID, &m.NamProdi,
		&m.NIM, &m.Nama,
		&dpID, &m.NamaDosen,
		&m.MaxSKS, &m.CreatedAt)
	m.UserID = uid
	m.DosenPembimbingID = dpID
	return &m, err
}

func (r *Repository) CreateMahasiswaWithUser(ctx context.Context, req CreateMahasiswaReq, passwordHash string) (*Mahasiswa, error) {
	now := time.Now()
	var userID string
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, 'MAHASISWA', true, $3, $3)
		RETURNING id`, req.Email, passwordHash, now,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	maxSKS := req.MaxSKS
	if maxSKS == 0 {
		maxSKS = 24
	}

	var m Mahasiswa
	var dpID *string
	err = r.db.QueryRow(ctx, `
		INSERT INTO mahasiswa (user_id, prodi_id, nim, nama, dosen_pembimbing_id, max_sks, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, prodi_id, nim, nama, dosen_pembimbing_id, max_sks, created_at`,
		userID, req.ProdiID, req.NIM, req.Nama, req.DosenPembimbingID, maxSKS, now,
	).Scan(&m.ID, &m.UserID, &m.ProdiID, &m.NIM, &m.Nama, &dpID, &m.MaxSKS, &m.CreatedAt)
	m.DosenPembimbingID = dpID
	return &m, err
}

func (r *Repository) PatchDPA(ctx context.Context, mahasiswaID, dosenID string) (oldDosenID string, err error) {
	err = r.db.QueryRow(ctx, `
		UPDATE mahasiswa SET dosen_pembimbing_id = $1
		WHERE id = $2
		RETURNING COALESCE((SELECT dosen_pembimbing_id::text FROM mahasiswa WHERE id = $2), '')`,
		dosenID, mahasiswaID,
	).Scan(&oldDosenID)
	return oldDosenID, err
}

func (r *Repository) GetOldDosenPembimbing(ctx context.Context, mahasiswaID string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(dosen_pembimbing_id::text,'') FROM mahasiswa WHERE id = $1`, mahasiswaID,
	).Scan(&id)
	return id, err
}
