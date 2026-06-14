import { Invalidator } from '../../shared/cache/redis';
import { DosenRepository } from './dosen.repository';
import { ApproveReq, PendingKRS, RejectReq } from './dosen.types';

export class DosenService {
  constructor(
    private repo: DosenRepository,
    private inv: Invalidator,
  ) {}

  async listPending(userId: string, role: string): Promise<PendingKRS[]> {
    const r = role.toUpperCase();
    if (r === 'STAFF' || r === 'ADMIN') return this.repo.listPendingAll();
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil dosen tidak ditemukan');
    return this.repo.listPendingForDosen(profile.id);
  }

  async approve(userId: string, role: string, krsId: string, req: ApproveReq): Promise<PendingKRS> {
    await this.checkDosenAuthorization(userId, role, krsId);
    const ok = await this.repo.approve(krsId, userId, req.catatan ?? '');
    if (!ok) throw new Error('KRS tidak dalam status diajukan atau tidak ditemukan');
    const { mahasiswaId } = await this.repo.getDosenPembimbingOfKrs(krsId);
    await this.inv.onKrsMutation(mahasiswaId, '');
    const krs = await this.repo.getKrs(krsId);
    if (!krs) throw new Error('KRS tidak ditemukan');
    return krs;
  }

  async reject(userId: string, role: string, krsId: string, req: RejectReq): Promise<PendingKRS> {
    if (!req.catatan?.trim()) throw new Error('catatan penolakan wajib diisi');
    await this.checkDosenAuthorization(userId, role, krsId);
    const ok = await this.repo.reject(krsId, userId, req.catatan);
    if (!ok) throw new Error('KRS tidak dalam status diajukan atau tidak ditemukan');
    const { mahasiswaId } = await this.repo.getDosenPembimbingOfKrs(krsId);
    await this.inv.onKrsMutation(mahasiswaId, '');
    const krs = await this.repo.getKrs(krsId);
    if (!krs) throw new Error('KRS tidak ditemukan');
    return krs;
  }

  private async checkDosenAuthorization(userId: string, role: string, krsId: string): Promise<void> {
    const r = role.toUpperCase();
    if (r === 'STAFF' || r === 'ADMIN') return;
    const { dosenId } = await this.repo.getDosenPembimbingOfKrs(krsId);
    const profile = await this.repo.getProfileByUserId(userId);
    if (!profile) throw new Error('profil dosen tidak ditemukan');
    if (dosenId !== profile.id) {
      throw new Error('otorisasi ditolak: mahasiswa ini bukan bimbingan Anda');
    }
  }
}
