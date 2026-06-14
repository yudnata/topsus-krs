import { Pool } from 'pg';
import {
  CreateDosenReq,
  CreateKelasReq,
  CreateMahasiswaReq,
  CreateMataKuliahReq,
  Dosen,
  Kelas,
  Mahasiswa,
  MataKuliah,
  ProgramStudi,
  TahunAkademik,
} from './admin.types';

export class AdminRepository {
  constructor(private db: Pool) {}

  async listProdi(): Promise<ProgramStudi[]> {
    const res = await this.db.query(
      'SELECT id, kode_prodi, nama_prodi, fakultas FROM program_studi ORDER BY nama_prodi',
    );
    return res.rows;
  }

  async getTahunAkademikAktif(): Promise<TahunAkademik | null> {
    const res = await this.db.query(
      'SELECT id, kode_ta, semester, is_active FROM tahun_akademik WHERE is_active = true LIMIT 1',
    );
    return res.rows[0] ?? null;
  }

  async listTahunAkademik(): Promise<TahunAkademik[]> {
    const res = await this.db.query(
      'SELECT id, kode_ta, semester, is_active FROM tahun_akademik ORDER BY kode_ta DESC',
    );
    return res.rows;
  }

  async listMataKuliah(): Promise<MataKuliah[]> {
    const res = await this.db.query(
      `SELECT id, COALESCE(prodi_id::text,'') AS prodi_id, kode_mk, nama_mk, sks
       FROM mata_kuliah ORDER BY kode_mk`,
    );
    return res.rows;
  }

  async createMataKuliah(req: CreateMataKuliahReq): Promise<MataKuliah> {
    const res = await this.db.query(
      `INSERT INTO mata_kuliah (prodi_id, kode_mk, nama_mk, sks)
       VALUES ($1, $2, $3, $4)
       RETURNING id, COALESCE(prodi_id::text,'') AS prodi_id, kode_mk, nama_mk, sks`,
      [req.prodi_id, req.kode_mk, req.nama_mk, req.sks],
    );
    return res.rows[0];
  }

  async listKelas(taId: string): Promise<Kelas[]> {
    let query = `
      SELECT k.id, k.mata_kuliah_id, mk.nama_mk, mk.kode_mk,
             k.tahun_akademik_id, COALESCE(k.dosen_id::text,'') AS dosen_id, COALESCE(d.nama,'') AS nama_dosen,
             k.nama_kelas, k.kapasitas, k.terisi
      FROM kelas k
      JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
      LEFT JOIN dosen d ON d.id = k.dosen_id`;
    const params: string[] = [];
    if (taId) {
      query += ' WHERE k.tahun_akademik_id = $1';
      params.push(taId);
    }
    query += ' ORDER BY mk.kode_mk, k.nama_kelas';
    const res = await this.db.query(query, params);
    return res.rows;
  }

  async createKelas(req: CreateKelasReq): Promise<Kelas> {
    const res = await this.db.query(
      `INSERT INTO kelas (mata_kuliah_id, tahun_akademik_id, dosen_id, nama_kelas, kapasitas)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING id, mata_kuliah_id, tahun_akademik_id,
                 COALESCE(dosen_id::text,'') AS dosen_id, nama_kelas, kapasitas, terisi`,
      [req.mata_kuliah_id, req.tahun_akademik_id, req.dosen_id || null, req.nama_kelas, req.kapasitas],
    );
    const row = res.rows[0];
    return { ...row, nama_mk: '', kode_mk: '', nama_dosen: '' };
  }

  async listDosen(): Promise<Dosen[]> {
    const res = await this.db.query(
      `SELECT id, user_id, nip, nama, COALESCE(email,'') AS email, created_at
       FROM dosen ORDER BY nama`,
    );
    return res.rows;
  }

  async getDosen(id: string): Promise<Dosen | null> {
    const res = await this.db.query(
      `SELECT id, user_id, nip, nama, COALESCE(email,'') AS email, created_at
       FROM dosen WHERE id = $1`,
      [id],
    );
    return res.rows[0] ?? null;
  }

  async createDosenWithUser(nip: string, nama: string, email: string, passwordHash: string): Promise<Dosen> {
    const now = new Date();
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');
      const userRes = await client.query(
        `INSERT INTO users (email, password_hash, role, is_active, created_at, updated_at)
         VALUES ($1, $2, 'DOSEN', true, $3, $3) RETURNING id`,
        [email, passwordHash, now],
      );
      const userId = userRes.rows[0].id;
      const dosenRes = await client.query(
        `INSERT INTO dosen (user_id, nip, nama, email, created_at)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, user_id, nip, nama, email, created_at`,
        [userId, nip, nama, email, now],
      );
      await client.query('COMMIT');
      return dosenRes.rows[0];
    } catch (err) {
      await client.query('ROLLBACK');
      throw err;
    } finally {
      client.release();
    }
  }

  async listMahasiswa(): Promise<Mahasiswa[]> {
    const res = await this.db.query(`
      SELECT m.id, m.user_id, m.prodi_id, COALESCE(ps.nama_prodi,'') AS nama_prodi,
             m.nim, m.nama, m.dosen_pembimbing_id, COALESCE(d.nama,'') AS nama_dosen,
             m.max_sks, m.created_at
      FROM mahasiswa m
      LEFT JOIN program_studi ps ON ps.id = m.prodi_id
      LEFT JOIN dosen d ON d.id = m.dosen_pembimbing_id
      ORDER BY m.nim`);
    return res.rows;
  }

  async getMahasiswa(id: string): Promise<Mahasiswa | null> {
    const res = await this.db.query(
      `SELECT m.id, m.user_id, m.prodi_id, COALESCE(ps.nama_prodi,'') AS nama_prodi,
              m.nim, m.nama, m.dosen_pembimbing_id, COALESCE(d.nama,'') AS nama_dosen,
              m.max_sks, m.created_at
       FROM mahasiswa m
       LEFT JOIN program_studi ps ON ps.id = m.prodi_id
       LEFT JOIN dosen d ON d.id = m.dosen_pembimbing_id
       WHERE m.id = $1`,
      [id],
    );
    return res.rows[0] ?? null;
  }

  async createMahasiswaWithUser(req: CreateMahasiswaReq, passwordHash: string): Promise<Mahasiswa> {
    const now = new Date();
    const maxSks = req.max_sks || 24;
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');
      const userRes = await client.query(
        `INSERT INTO users (email, password_hash, role, is_active, created_at, updated_at)
         VALUES ($1, $2, 'MAHASISWA', true, $3, $3) RETURNING id`,
        [req.email, passwordHash, now],
      );
      const userId = userRes.rows[0].id;
      const mhsRes = await client.query(
        `INSERT INTO mahasiswa (user_id, prodi_id, nim, nama, dosen_pembimbing_id, max_sks, created_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         RETURNING id, user_id, prodi_id, nim, nama, dosen_pembimbing_id, max_sks, created_at`,
        [userId, req.prodi_id, req.nim, req.nama, req.dosen_pembimbing_id ?? null, maxSks, now],
      );
      await client.query('COMMIT');
      return { ...mhsRes.rows[0], nama_prodi: '', nama_dosen: '' };
    } catch (err) {
      await client.query('ROLLBACK');
      throw err;
    } finally {
      client.release();
    }
  }

  async getOldDosenPembimbing(mahasiswaId: string): Promise<string> {
    const res = await this.db.query(
      `SELECT COALESCE(dosen_pembimbing_id::text,'') AS id FROM mahasiswa WHERE id = $1`,
      [mahasiswaId],
    );
    return res.rows[0]?.id ?? '';
  }

  async patchDPA(mahasiswaId: string, dosenId: string): Promise<void> {
    await this.db.query(
      'UPDATE mahasiswa SET dosen_pembimbing_id = $1 WHERE id = $2',
      [dosenId, mahasiswaId],
    );
  }
}
