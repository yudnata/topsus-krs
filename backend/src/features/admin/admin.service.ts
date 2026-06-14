import bcrypt from 'bcrypt';
import { Invalidator } from '../../shared/cache/redis';
import { AdminRepository } from './admin.repository';
import {
  CreateDosenReq,
  CreateKelasReq,
  CreateMahasiswaReq,
  CreateMataKuliahReq,
  Dosen,
  Kelas,
  Mahasiswa,
  MataKuliah,
  PatchDPAReq,
  ProgramStudi,
  TahunAkademik,
} from './admin.types';

export class AdminService {
  constructor(
    private repo: AdminRepository,
    private inv: Invalidator,
  ) {}

  listProdi(): Promise<ProgramStudi[]> {
    return this.repo.listProdi();
  }

  listTahunAkademik(): Promise<TahunAkademik[]> {
    return this.repo.listTahunAkademik();
  }

  getTahunAkademikAktif(): Promise<TahunAkademik | null> {
    return this.repo.getTahunAkademikAktif();
  }

  listMataKuliah(): Promise<MataKuliah[]> {
    return this.repo.listMataKuliah();
  }

  async createMataKuliah(req: CreateMataKuliahReq): Promise<MataKuliah> {
    req.kode_mk = req.kode_mk?.trim() ?? '';
    req.nama_mk = req.nama_mk?.trim() ?? '';
    if (!req.kode_mk || !req.nama_mk || req.sks <= 0) {
      throw new Error('kode_mk, nama_mk, dan sks wajib diisi (sks > 0)');
    }
    return this.repo.createMataKuliah(req);
  }

  listKelas(taId: string): Promise<Kelas[]> {
    return this.repo.listKelas(taId);
  }

  async createKelas(req: CreateKelasReq): Promise<Kelas> {
    if (!req.mata_kuliah_id || !req.tahun_akademik_id || !req.nama_kelas) {
      throw new Error('mata_kuliah_id, tahun_akademik_id, nama_kelas wajib diisi');
    }
    if (!req.kapasitas || req.kapasitas <= 0) req.kapasitas = 40;
    return this.repo.createKelas(req);
  }

  listDosen(): Promise<Dosen[]> {
    return this.repo.listDosen();
  }

  getDosen(id: string): Promise<Dosen | null> {
    return this.repo.getDosen(id);
  }

  async createDosen(req: CreateDosenReq): Promise<Dosen> {
    req.nip = req.nip?.trim() ?? '';
    req.nama = req.nama?.trim() ?? '';
    req.email = req.email?.trim().toLowerCase() ?? '';
    if (!req.nip || !req.nama || !req.email || !req.password) {
      throw new Error('nip, nama, email, password wajib diisi');
    }
    const hash = await bcrypt.hash(req.password, 10);
    const d = await this.repo.createDosenWithUser(req.nip, req.nama, req.email, hash);
    await this.inv.onDosenWrite(d.id);
    return d;
  }

  listMahasiswa(): Promise<Mahasiswa[]> {
    return this.repo.listMahasiswa();
  }

  getMahasiswa(id: string): Promise<Mahasiswa | null> {
    return this.repo.getMahasiswa(id);
  }

  async createMahasiswa(req: CreateMahasiswaReq): Promise<Mahasiswa> {
    req.nim = req.nim?.trim() ?? '';
    req.nama = req.nama?.trim() ?? '';
    req.email = req.email?.trim().toLowerCase() ?? '';
    if (!req.nim || !req.nama || !req.prodi_id || !req.email || !req.password) {
      throw new Error('nim, nama, prodi_id, email, password wajib diisi');
    }
    const hash = await bcrypt.hash(req.password, 10);
    const m = await this.repo.createMahasiswaWithUser(req, hash);
    await this.inv.onMahasiswaWrite(m.id, m.dosen_pembimbing_id ?? '');
    return m;
  }

  async patchDPA(mahasiswaId: string, req: PatchDPAReq): Promise<void> {
    const oldDosenId = await this.repo.getOldDosenPembimbing(mahasiswaId);
    await this.repo.patchDPA(mahasiswaId, req.dosen_pembimbing_id);
    await this.inv.onDpaChange(oldDosenId, req.dosen_pembimbing_id, mahasiswaId);
  }
}
