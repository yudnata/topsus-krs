import { Router } from 'express';
import { AuthController } from './auth.controller';
import { requireAuth } from './auth.middleware';
import { AuthService } from './auth.service';

export function createAuthRoutes(svc: AuthService): Router {
  const router = Router();
  const ctrl = new AuthController(svc);

  router.post('/register', ctrl.register);
  router.post('/login', ctrl.login);
  router.get('/profile', requireAuth(svc), ctrl.profile);

  return router;
}
