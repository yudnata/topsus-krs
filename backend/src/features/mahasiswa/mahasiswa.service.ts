import { Invalidator } from '../../shared/cache/redis';
import { MahasiswaRepository } from './mahasiswa.repository';
import { KelasDetail, KrsHeader } from './mahasiswa.types';

function timeStr(t: string): string {
  return t.length >= 5 ? t.slice(0, 5) : t;
}

function checkScheduleConflict(existing: KelasDetail[], newKelas: KelasDetail): Error | null {
  for (const ex of existing) {
    for (const exJ of ex.jadwal) {
      for (const newJ of newKelas.jadwal) {
        if (exJ.hari.toLowerCase() !== newJ.hari.toLowerCase()) continue;
        if (
          timeStr(exJ.jam_mulai) < timeStr(newJ.jam_selesai) &&
          timeStr(newJ.jam_mulai) < timeStr(exJ.jam_selesai)
        ) {
          return new Error(`bentrok jadwal hari ${newJ.hari} (${newJ.jam_mulai}–${newJ.jam_selesai})`);
        }
      }
    }
  }
  return null;
}

export class MahasiswaService {
  constructor(
    private repo: MahasiswaRepository,
    private inv: Invalidator,
  ) {}

  async getCurrentKRS(userId: string): Promise<KrsHeader> {
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil mahasiswa tidak ditemukan');
    const ta = await this.repo.getTahunAkademikAktif();
    if (!ta) throw new Error('tidak ada tahun akademik aktif');
    const krs = await this.repo.getOrCreateDraft(profile.id, ta.id);
    krs.kode_ta = ta.kode_ta;
    krs.semester = ta.semester;
    krs.detail = await this.repo.getDetails(krs.id);
    return krs;
  }

  async addClass(userId: string, krsId: string, kelasId: string): Promise<KrsHeader> {
    const krs = await this.repo.getKrsById(krsId);
    if (!krs) throw new Error('KRS tidak ditemukan');
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil mahasiswa tidak ditemukan');
    if (krs.mahasiswa_id !== profile.id) throw new Error('akses ditolak: bukan KRS Anda');
    if (krs.status !== 'draft') throw new Error(`KRS tidak dapat diedit (status: ${krs.status})`);

    const newKelas = await this.repo.getKelasDetail(kelasId);
    if (!newKelas) throw new Error('kelas tidak ditemukan');
    if (newKelas.terisi >= newKelas.kapasitas) throw new Error('kapasitas kelas sudah penuh');

    const existing = await this.repo.getCurrentKelasInKrs(krsId);
    for (const ex of existing) {
      if (ex.id === kelasId) throw new Error('kelas ini sudah ada di KRS Anda');
      if (ex.mata_kuliah_id === newKelas.mata_kuliah_id) {
        throw new Error('mata kuliah ini sudah ada di KRS Anda (kelas berbeda)');
      }
    }
    if (krs.total_sks + newKelas.sks > profile.max_sks) {
      throw new Error(`melebihi batas SKS (${krs.total_sks + newKelas.sks}/${profile.max_sks})`);
    }
    const conflict = checkScheduleConflict(existing, newKelas);
    if (conflict) throw conflict;

    await this.repo.addKelas(krsId, kelasId, newKelas.sks);
    await this.inv.onKrsMutation(profile.id, profile.dosen_pembimbing_id ?? '');
    return this.getKrsWithDetail(krsId);
  }

  async removeClass(userId: string, krsId: string, kelasId: string): Promise<KrsHeader> {
    const krs = await this.repo.getKrsById(krsId);
    if (!krs) throw new Error('KRS tidak ditemukan');
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil mahasiswa tidak ditemukan');
    if (krs.mahasiswa_id !== profile.id) throw new Error('akses ditolak: bukan KRS Anda');
    if (krs.status !== 'draft') throw new Error(`KRS tidak dapat diedit (status: ${krs.status})`);

    const kd = await this.repo.getKelasDetail(kelasId);
    if (!kd) throw new Error('kelas tidak ditemukan');
    const ok = await this.repo.removeKelas(krsId, kelasId, kd.sks);
    if (!ok) throw new Error('kelas tidak ada di KRS Anda');

    await this.inv.onKrsMutation(profile.id, profile.dosen_pembimbing_id ?? '');
    return this.getKrsWithDetail(krsId);
  }

  async submitKRS(userId: string, krsId: string): Promise<KrsHeader> {
    const krs = await this.repo.getKrsById(krsId);
    if (!krs) throw new Error('KRS tidak ditemukan');
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil mahasiswa tidak ditemukan');
    if (krs.mahasiswa_id !== profile.id) throw new Error('akses ditolak: bukan KRS Anda');
    if (krs.total_sks === 0) throw new Error('KRS kosong, tambahkan minimal 1 kelas');
    const ok = await this.repo.submitKrs(krsId);
    if (!ok) throw new Error('KRS tidak dalam status draft atau tidak ditemukan');

    await this.inv.onKrsMutation(profile.id, profile.dosen_pembimbing_id ?? '');
    return this.getKrsWithDetail(krsId);
  }

  async getHistory(userId: string): Promise<KrsHeader[]> {
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil mahasiswa tidak ditemukan');
    return this.repo.getHistory(profile.id);
  }

  private async getKrsWithDetail(krsId: string): Promise<KrsHeader> {
    const krs = await this.repo.getKrsById(krsId);
    if (!krs) throw new Error('KRS tidak ditemukan');
    krs.detail = await this.repo.getDetails(krsId);
    return krs;
  }
}
