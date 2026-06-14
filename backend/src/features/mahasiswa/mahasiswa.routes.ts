import { Router } from 'express';
import { requireAuth, requireRole } from '../auth/auth.middleware';
import { AuthService } from '../auth/auth.service';
import { MahasiswaController } from './mahasiswa.controller';
import { MahasiswaService } from './mahasiswa.service';

export function createMahasiswaRoutes(svc: MahasiswaService, authSvc: AuthService): Router {
  const router = Router();
  const ctrl = new MahasiswaController(svc);

  router.use(requireAuth(authSvc), requireRole('MAHASISWA'));

  router.get('/current', ctrl.getCurrentKRS);
  router.get('/history', ctrl.getHistory);
  router.post('/:id/add-class', ctrl.addClass);
  router.delete('/:id/remove-class', ctrl.removeClass);
  router.post('/:id/submit', ctrl.submitKRS);

  return router;
}
