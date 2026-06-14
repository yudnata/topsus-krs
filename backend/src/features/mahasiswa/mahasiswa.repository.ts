import { Pool } from 'pg';
import { JadwalInfo, JadwalRaw, KelasDetail, KrsDetail, KrsHeader, MahasiswaProfile } from './mahasiswa.types';

export class MahasiswaRepository {
  constructor(private db: Pool) {}

  async getProfileByUserId(userId: string): Promise<MahasiswaProfile | null> {
    const res = await this.db.query(
      `SELECT id, nama, nim, max_sks, dosen_pembimbing_id
       FROM mahasiswa WHERE user_id = $1`,
      [userId],
    );
    return res.rows[0] ?? null;
  }

  async getTahunAkademikAktif(): Promise<{ id: string; kode_ta: string; semester: string } | null> {
    const res = await this.db.query(
      'SELECT id, kode_ta, semester FROM tahun_akademik WHERE is_active = true LIMIT 1',
    );
    return res.rows[0] ?? null;
  }

  async getOrCreateDraft(mahasiswaId: string, taId: string): Promise<KrsHeader> {
    const existing = await this.db.query(
      `SELECT pk.id, pk.mahasiswa_id, m.nama AS nama_mahasiswa,
              pk.tahun_akademik_id, ta.kode_ta, ta.semester,
              pk.status, pk.total_sks,
              COALESCE(pk.catatan_mhs,'') AS catatan_mhs,
              COALESCE(pk.catatan_reviewer,'') AS catatan_reviewer,
              pk.reviewed_at, pk.created_at, pk.updated_at
       FROM pengajuan_krs pk
       JOIN mahasiswa m ON m.id = pk.mahasiswa_id
       JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
       WHERE pk.mahasiswa_id = $1 AND pk.tahun_akademik_id = $2`,
      [mahasiswaId, taId],
    );
    if (existing.rows[0]) return existing.rows[0];

    const now = new Date();
    const created = await this.db.query(
      `INSERT INTO pengajuan_krs (mahasiswa_id, tahun_akademik_id, status, total_sks, created_at, updated_at)
       VALUES ($1, $2, 'draft', 0, $3, $3)
       RETURNING id, mahasiswa_id, tahun_akademik_id, status, total_sks, created_at, updated_at`,
      [mahasiswaId, taId, now],
    );
    return {
      ...created.rows[0],
      nama_mahasiswa: '',
      kode_ta: '',
      semester: '',
      catatan_mhs: '',
      catatan_reviewer: '',
    };
  }

  async getKrsById(krsId: string): Promise<KrsHeader | null> {
    const res = await this.db.query(
      `SELECT pk.id, pk.mahasiswa_id, m.nama AS nama_mahasiswa,
              pk.tahun_akademik_id, ta.kode_ta, ta.semester,
              pk.status, pk.total_sks,
              COALESCE(pk.catatan_mhs,'') AS catatan_mhs,
              COALESCE(pk.catatan_reviewer,'') AS catatan_reviewer,
              pk.reviewed_at, pk.created_at, pk.updated_at
       FROM pengajuan_krs pk
       JOIN mahasiswa m ON m.id = pk.mahasiswa_id
       JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
       WHERE pk.id = $1`,
      [krsId],
    );
    return res.rows[0] ?? null;
  }

  async getDetails(krsId: string): Promise<KrsDetail[]> {
    const res = await this.db.query(
      `SELECT pkd.id, pkd.pengajuan_krs_id AS krs_id, pkd.kelas_id,
              k.nama_kelas, mk.id AS mata_kuliah_id, mk.kode_mk, mk.nama_mk, mk.sks,
              COALESCE(d.nama,'') AS nama_dosen
       FROM pengajuan_krs_detail pkd
       JOIN kelas k ON k.id = pkd.kelas_id
       JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
       LEFT JOIN dosen d ON d.id = k.dosen_id
       WHERE pkd.pengajuan_krs_id = $1
       ORDER BY mk.kode_mk`,
      [krsId],
    );
    const result: KrsDetail[] = [];
    for (const row of res.rows) {
      const jadwal = await this.getJadwalForKelas(row.kelas_id);
      result.push({ ...row, jadwal });
    }
    return result;
  }

  private async getJadwalForKelas(kelasId: string): Promise<JadwalInfo[]> {
    const res = await this.db.query(
      `SELECT hari, jam_mulai::text, jam_selesai::text, ruangan
       FROM jadwal_kelas WHERE kelas_id = $1`,
      [kelasId],
    );
    return res.rows;
  }

  async getKelasDetail(kelasId: string): Promise<KelasDetail | null> {
    const res = await this.db.query(
      `SELECT k.id, k.mata_kuliah_id, mk.sks, k.kapasitas, k.terisi
       FROM kelas k JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id WHERE k.id = $1`,
      [kelasId],
    );
    if (!res.rows[0]) return null;
    const k: KelasDetail = { ...res.rows[0], jadwal: [] };
    const jRes = await this.db.query(
      `SELECT hari, jam_mulai::text, jam_selesai::text FROM jadwal_kelas WHERE kelas_id = $1`,
      [kelasId],
    );
    k.jadwal = jRes.rows as JadwalRaw[];
    return k;
  }

  async getCurrentKelasInKrs(krsId: string): Promise<KelasDetail[]> {
    const res = await this.db.query(
      `SELECT k.id, k.mata_kuliah_id, mk.sks, k.kapasitas, k.terisi
       FROM pengajuan_krs_detail pkd
       JOIN kelas k ON k.id = pkd.kelas_id
       JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
       WHERE pkd.pengajuan_krs_id = $1`,
      [krsId],
    );
    const result: KelasDetail[] = [];
    for (const row of res.rows) {
      const jRes = await this.db.query(
        `SELECT hari, jam_mulai::text, jam_selesai::text FROM jadwal_kelas WHERE kelas_id = $1`,
        [row.id],
      );
      result.push({ ...row, jadwal: jRes.rows });
    }
    return result;
  }

  async addKelas(krsId: string, kelasId: string, sks: number): Promise<void> {
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');
      await client.query(
        'INSERT INTO pengajuan_krs_detail (pengajuan_krs_id, kelas_id) VALUES ($1, $2)',
        [krsId, kelasId],
      );
      await client.query('UPDATE kelas SET terisi = terisi + 1 WHERE id = $1', [kelasId]);
      await client.query(
        'UPDATE pengajuan_krs SET total_sks = total_sks + $1, updated_at = NOW() WHERE id = $2',
        [sks, krsId],
      );
      await client.query('COMMIT');
    } catch (err) {
      await client.query('ROLLBACK');
      throw err;
    } finally {
      client.release();
    }
  }

  async removeKelas(krsId: string, kelasId: string, sks: number): Promise<boolean> {
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');
      const del = await client.query(
        'DELETE FROM pengajuan_krs_detail WHERE pengajuan_krs_id = $1 AND kelas_id = $2',
        [krsId, kelasId],
      );
      if (del.rowCount === 0) {
        await client.query('ROLLBACK');
        return false;
      }
      await client.query('UPDATE kelas SET terisi = GREATEST(terisi - 1, 0) WHERE id = $1', [kelasId]);
      await client.query(
        'UPDATE pengajuan_krs SET total_sks = GREATEST(total_sks - $1, 0), updated_at = NOW() WHERE id = $2',
        [sks, krsId],
      );
      await client.query('COMMIT');
      return true;
    } catch (err) {
      await client.query('ROLLBACK');
      throw err;
    } finally {
      client.release();
    }
  }

  async submitKrs(krsId: string): Promise<boolean> {
    const res = await this.db.query(
      `UPDATE pengajuan_krs SET status = 'diajukan', updated_at = NOW()
       WHERE id = $1 AND status = 'draft'`,
      [krsId],
    );
    return (res.rowCount ?? 0) > 0;
  }

  async getHistory(mahasiswaId: string): Promise<KrsHeader[]> {
    const res = await this.db.query(
      `SELECT pk.id, pk.mahasiswa_id, m.nama AS nama_mahasiswa,
              pk.tahun_akademik_id, ta.kode_ta, ta.semester,
              pk.status, pk.total_sks,
              COALESCE(pk.catatan_mhs,'') AS catatan_mhs,
              COALESCE(pk.catatan_reviewer,'') AS catatan_reviewer,
              pk.reviewed_at, pk.created_at, pk.updated_at
       FROM pengajuan_krs pk
       JOIN mahasiswa m ON m.id = pk.mahasiswa_id
       JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
       WHERE pk.mahasiswa_id = $1
       ORDER BY pk.created_at DESC`,
      [mahasiswaId],
    );
    return res.rows;
  }
}
