import { Request, Response } from 'express';
import { sendJson } from '../../shared/utils/response';
import { DosenService } from './dosen.service';

export class DosenController {
  constructor(private svc: DosenService) {}

  listPending = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.listPending(req.userId!, req.role!);
      sendJson(res, 200, true, 'OK', data);
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  approve = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.approve(req.userId!, req.role!, req.params.id, req.body);
      sendJson(res, 200, true, 'KRS disetujui', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };

  reject = async (req: Request, res: Response): Promise<void> => {
    try {
      const data = await this.svc.reject(req.userId!, req.role!, req.params.id, req.body);
      sendJson(res, 200, true, 'KRS ditolak', data);
    } catch (err) {
      sendJson(res, 422, false, err instanceof Error ? err.message : 'Error');
    }
  };
}
