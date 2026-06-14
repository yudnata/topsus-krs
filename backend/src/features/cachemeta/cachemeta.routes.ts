import { Router } from 'express';
import { CacheClient } from '../../shared/cache/redis';
import { CachemetaController } from './cachemeta.controller';

export function createCachemetaRoutes(cache: CacheClient): Router {
  const router = Router();
  const ctrl = new CachemetaController(cache);

  router.get('/status', ctrl.status);
  router.get('/demo', ctrl.demoGet);
  router.post('/demo', ctrl.demoSet);
  router.delete('/demo', ctrl.demoInvalidate);

  return router;
}
