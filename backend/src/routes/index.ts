import { Express, Request, Response } from 'express';
import { Pool } from 'pg';
import { createAdminRoutes } from '../features/admin/admin.routes';
import { AdminRepository } from '../features/admin/admin.repository';
import { AdminService } from '../features/admin/admin.service';
import { createAuthRoutes } from '../features/auth/auth.routes';
import { AuthRepository } from '../features/auth/auth.repository';
import { AuthService } from '../features/auth/auth.service';
import { createCachemetaRoutes } from '../features/cachemeta/cachemeta.routes';
import { createDosenRoutes } from '../features/dosen/dosen.routes';
import { DosenRepository } from '../features/dosen/dosen.repository';
import { DosenService } from '../features/dosen/dosen.service';
import { createMahasiswaRoutes } from '../features/mahasiswa/mahasiswa.routes';
import { MahasiswaRepository } from '../features/mahasiswa/mahasiswa.repository';
import { MahasiswaService } from '../features/mahasiswa/mahasiswa.service';
import { CacheClient, Invalidator } from '../shared/cache/redis';
import { sendJson } from '../shared/utils/response';
import { Config } from '../config';

export interface AppDeps {
  db: Pool;
  cache: CacheClient;
  inv: Invalidator;
  cfg: Config;
}

async function healthHandler(req: Request, res: Response, cache: CacheClient): Promise<void> {
  let redisStatus = 'disabled';
  if (cache.isEnabled()) {
    try {
      await cache.ping();
      redisStatus = 'ok';
    } catch {
      redisStatus = 'error';
    }
  }
  sendJson(res, 200, true, 'OK', { status: 'ok', service: 'krs-api', redis: redisStatus });
}

export function setupRoutes(app: Express, deps: AppDeps): void {
  const { db, cache, inv, cfg } = deps;

  const authRepo = new AuthRepository(db);
  const authSvc = new AuthService(authRepo, cfg);

  const adminSvc = new AdminService(new AdminRepository(db), inv);
  const mhsSvc = new MahasiswaService(new MahasiswaRepository(db), inv);
  const dosenSvc = new DosenService(new DosenRepository(db), inv);

  const api = app;

  api.get('/api/health', (req, res) => void healthHandler(req, res, cache));
  api.use('/api/auth', createAuthRoutes(authSvc));
  api.use('/api/cache', createCachemetaRoutes(cache));
  api.use('/api/admin', createAdminRoutes(adminSvc, authSvc));
  api.use('/api/krs', createMahasiswaRoutes(mhsSvc, authSvc));
  api.use('/api/approval', createDosenRoutes(dosenSvc, authSvc));
}
