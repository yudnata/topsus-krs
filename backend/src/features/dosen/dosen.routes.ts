import { Router } from 'express';
import { requireAuth, requireRole } from '../auth/auth.middleware';
import { AuthService } from '../auth/auth.service';
import { DosenController } from './dosen.controller';
import { DosenService } from './dosen.service';

export function createDosenRoutes(svc: DosenService, authSvc: AuthService): Router {
  const router = Router();
  const ctrl = new DosenController(svc);

  router.use(requireAuth(authSvc), requireRole('DOSEN', 'STAFF', 'ADMIN'));

  router.get('/pending', ctrl.listPending);
  router.post('/:id/approve', ctrl.approve);
  router.post('/:id/reject', ctrl.reject);

  return router;
}
