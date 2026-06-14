import { Pool } from 'pg';
import { DosenProfile, KrsDetailItem, PendingKRS } from './dosen.types';

export class DosenRepository {
  constructor(private db: Pool) {}

  async getProfileByUserId(userId: string): Promise<DosenProfile | null> {
    const res = await this.db.query('SELECT id, nama FROM dosen WHERE user_id = $1', [userId]);
    return res.rows[0] ?? null;
  }

  async listPendingForDosen(dosenId: string): Promise<PendingKRS[]> {
    return this.queryPending(
      `SELECT pk.id, pk.mahasiswa_id, m.nama AS nama_mahasiswa, m.nim,
              pk.tahun_akademik_id, ta.kode_ta, ta.semester,
              pk.status, pk.total_sks,
              COALESCE(pk.catatan_mhs,'') AS catatan_mhs,
              COALESCE(pk.catatan_reviewer,'') AS catatan_reviewer,
              pk.reviewed_at, pk.created_at
       FROM pengajuan_krs pk
       JOIN mahasiswa m ON m.id = pk.mahasiswa_id
       JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
       WHERE pk.status = 'diajukan' AND m.dosen_pembimbing_id = $1
       ORDER BY pk.created_at ASC`,
      [dosenId],
    );
  }

  async listPendingAll(): Promise<PendingKRS[]> {
    return this.queryPending(
      `SELECT pk.id, pk.mahasiswa_id, m.nama AS nama_mahasiswa, m.nim,
              pk.tahun_akademik_id, ta.kode_ta, ta.semester,
              pk.status, pk.total_sks,
              COALESCE(pk.catatan_mhs,'') AS catatan_mhs,
              COALESCE(pk.catatan_reviewer,'') AS catatan_reviewer,
              pk.reviewed_at, pk.created_at
       FROM pengajuan_krs pk
       JOIN mahasiswa m ON m.id = pk.mahasiswa_id
       JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
       WHERE pk.status = 'diajukan'
       ORDER BY pk.created_at ASC`,
    );
  }

  private async queryPending(query: string, args: string[] = []): Promise<PendingKRS[]> {
    const res = await this.db.query(query, args);
    const result: PendingKRS[] = [];
    for (const row of res.rows) {
      row.detail = await this.getDetailRingkas(row.id);
      result.push(row);
    }
    return result;
  }

  private async getDetailRingkas(krsId: string): Promise<KrsDetailItem[]> {
    const res = await this.db.query(
      `SELECT pkd.kelas_id, k.nama_kelas, mk.kode_mk, mk.nama_mk, mk.sks
       FROM pengajuan_krs_detail pkd
       JOIN kelas k ON k.id = pkd.kelas_id
       JOIN mata_kuliah mk ON mk.id = k.mata_kuliah_id
       WHERE pkd.pengajuan_krs_id = $1
       ORDER BY mk.kode_mk`,
      [krsId],
    );
    return res.rows;
  }

  async getKrs(krsId: string): Promise<PendingKRS | null> {
    const res = await this.db.query(
      `SELECT pk.id, pk.mahasiswa_id, m.nama AS nama_mahasiswa, m.nim,
              pk.tahun_akademik_id, ta.kode_ta, ta.semester,
              pk.status, pk.total_sks,
              COALESCE(pk.catatan_mhs,'') AS catatan_mhs,
              COALESCE(pk.catatan_reviewer,'') AS catatan_reviewer,
              pk.reviewed_at, pk.created_at
       FROM pengajuan_krs pk
       JOIN mahasiswa m ON m.id = pk.mahasiswa_id
       JOIN tahun_akademik ta ON ta.id = pk.tahun_akademik_id
       WHERE pk.id = $1`,
      [krsId],
    );
    if (!res.rows[0]) return null;
    res.rows[0].detail = await this.getDetailRingkas(krsId);
    return res.rows[0];
  }

  async getDosenPembimbingOfKrs(krsId: string): Promise<{ dosenId: string; mahasiswaId: string }> {
    const res = await this.db.query(
      `SELECT COALESCE(m.dosen_pembimbing_id::text,'') AS dosen_id, m.id AS mahasiswa_id
       FROM pengajuan_krs pk JOIN mahasiswa m ON m.id = pk.mahasiswa_id WHERE pk.id = $1`,
      [krsId],
    );
    return { dosenId: res.rows[0]?.dosen_id ?? '', mahasiswaId: res.rows[0]?.mahasiswa_id ?? '' };
  }

  async approve(krsId: string, reviewerId: string, catatan: string): Promise<boolean> {
    const res = await this.db.query(
      `UPDATE pengajuan_krs
       SET status = 'disetujui', reviewed_by = $1, reviewed_at = NOW(),
           catatan_reviewer = $2, updated_at = NOW()
       WHERE id = $3 AND status = 'diajukan'`,
      [reviewerId, catatan, krsId],
    );
    return (res.rowCount ?? 0) > 0;
  }

  async reject(krsId: string, reviewerId: string, catatan: string): Promise<boolean> {
    const res = await this.db.query(
      `UPDATE pengajuan_krs
       SET status = 'ditolak', reviewed_by = $1, reviewed_at = NOW(),
           catatan_reviewer = $2, updated_at = NOW()
       WHERE id = $3 AND status = 'diajukan'`,
      [reviewerId, catatan, krsId],
    );
    return (res.rowCount ?? 0) > 0;
  }
}
