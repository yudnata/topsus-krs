import { Request, Response } from 'express';
import { sendJson } from '../../shared/utils/response';
import { AdminService } from './admin.service';

export class AdminController {
  constructor(private svc: AdminService) {}

  listProdi = async (_req: Request, res: Response): Promise<void> => {
    try {
      sendJson(res, 200, true, 'OK', await this.svc.listProdi());
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  listTahunAkademik = async (_req: Request, res: Response): Promise<void> => {
    try {
      sendJson(res, 200, true, 'OK', await this.svc.listTahunAkademik());
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  getTahunAkademikAktif = async (_req: Request, res: Response): Promise<void> => {
    const data = await this.svc.getTahunAkademikAktif();
    if (!data) {
      sendJson(res, 404, false, 'Tidak ada tahun akademik aktif');
      return;
    }
    sendJson(res, 200, true, 'OK', data);
  };

  listMataKuliah = async (_req: Request, res: Response): Promise<void> => {
    try {
      sendJson(res, 200, true, 'OK', await this.svc.listMataKuliah());
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  createMataKuliah = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.createMataKuliah(req.body);
      sendJson(res, 201, true, 'Mata kuliah dibuat', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };

  listKelas = async (req: Request, res: Response): Promise<void> => {
    try {
      const taId = (req.query.tahun_akademik_id as string) ?? '';
      sendJson(res, 200, true, 'OK', await this.svc.listKelas(taId));
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  createKelas = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.createKelas(req.body);
      sendJson(res, 201, true, 'Kelas dibuat', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };

  listDosen = async (_req: Request, res: Response): Promise<void> => {
    try {
      sendJson(res, 200, true, 'OK', await this.svc.listDosen());
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  getDosen = async (req: Request, res: Response): Promise<void> => {
    const data = await this.svc.getDosen(req.params.id);
    if (!data) {
      sendJson(res, 404, false, 'Dosen tidak ditemukan');
      return;
    }
    sendJson(res, 200, true, 'OK', data);
  };

  createDosen = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.createDosen(req.body);
      sendJson(res, 201, true, 'Dosen dibuat', data);
    } catch (err) {
      sendJson(res, 409, false, err instanceof Error ? err.message : 'Error');
    }
  };

  listMahasiswa = async (_req: Request, res: Response): Promise<void> => {
    try {
      sendJson(res, 200, true, 'OK', await this.svc.listMahasiswa());
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  getMahasiswa = async (req: Request, res: Response): Promise<void> => {
    const data = await this.svc.getMahasiswa(req.params.id);
    if (!data) {
      sendJson(res, 404, false, 'Mahasiswa tidak ditemukan');
      return;
    }
    sendJson(res, 200, true, 'OK', data);
  };

  createMahasiswa = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.createMahasiswa(req.body);
      sendJson(res, 201, true, 'Mahasiswa dibuat', data);
    } catch (err) {
      sendJson(res, 409, false, err instanceof Error ? err.message : 'Error');
    }
  };

  patchDPA = async (req: Request, res: Response): Promise<void> => {
    if (!req.body?.dosen_pembimbing_id) {
      sendJson(res, 400, false, 'dosen_pembimbing_id wajib diisi');
      return;
    }
    try {
      await this.svc.patchDPA(req.params.id, req.body);
      sendJson(res, 200, true, 'DPA berhasil diperbarui');
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };
}
