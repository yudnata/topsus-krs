import { Request, Response } from 'express';
import { sendJson } from '../../shared/utils/response';
import { MahasiswaService } from './mahasiswa.service';

export class MahasiswaController {
  constructor(private svc: MahasiswaService) {}

  getCurrentKRS = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.getCurrentKRS(req.userId!);
      sendJson(res, 200, true, 'OK', data);
    } catch (err) {
      sendJson(res, 404, false, err instanceof Error ? err.message : 'Error');
    }
  };

  addClass = async (req: Request, res: Response): Promise<void> => {
    if (!req.body?.kelas_id) {
      sendJson(res, 400, false, 'kelas_id wajib diisi');
      return;
    }
    try {
      const data = await this.svc.addClass(req.userId!, req.params.id, req.body.kelas_id);
      sendJson(res, 200, true, 'Kelas ditambahkan', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };

  removeClass = async (req: Request, res: Response): Promise<void> => {
    if (!req.body?.kelas_id) {
      sendJson(res, 400, false, 'kelas_id wajib diisi');
      return;
    }
    try {
      const data = await this.svc.removeClass(req.userId!, req.params.id, req.body.kelas_id);
      sendJson(res, 200, true, 'Kelas dihapus', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };

  submitKRS = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.submitKRS(req.userId!, req.params.id);
      sendJson(res, 200, true, 'KRS berhasil diajukan', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };

  getHistory = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.getHistory(req.userId!);
      sendJson(res, 200, true, 'OK', data);
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };
}
