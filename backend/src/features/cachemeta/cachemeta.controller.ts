import { Request, Response } from 'express';
import { TTL_DOSEN_MAHASISWA_LIST } from '../../shared/cache/keys';
import { CacheClient } from '../../shared/cache/redis';
import { sendJson } from '../../shared/utils/response';

export class CachemetaController {
  constructor(private cache: CacheClient) {}

  status = async (_req: Request, res: Response): Promise<void> => {
    let redisStatus = 'disabled';
    if (this.cache.isEnabled()) {
      try {
        await this.cache.ping();
        redisStatus = 'ok';
      } catch {
        redisStatus = 'error';
      }
    }
    sendJson(res, 200, true, 'OK', { redis: redisStatus });
  };

  demoGet = async (req: Request, res: Response): Promise<void> => {
    const key = (req.query.key as string) || 'cache:demo';
    try {
      const { value, hit } = await this.cache.get(key);
      if (hit) {
        res.setHeader('X-Cache', 'HIT');
        sendJson(res, 200, true, 'Cache hit', { key, value });
        return;
      }
      res.setHeader('X-Cache', 'MISS');
      sendJson(res, 200, true, 'Cache miss', { key });
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  demoSet = async (req: Request, res: Response): Promise<void> => {
    const key = (req.query.key as string) || 'cache:demo';
    const value = (req.query.value as string) || 'hello-krs';
    try {
      await this.cache.set(key, value, TTL_DOSEN_MAHASISWA_LIST);
      sendJson(res, 200, true, 'Cached', { key, ttl_seconds: TTL_DOSEN_MAHASISWA_LIST });
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };

  demoInvalidate = async (req: Request, res: Response): Promise<void> => {
    const key = (req.query.key as string) || 'cache:demo';
    try {
      await this.cache.del(key);
      sendJson(res, 200, true, 'Invalidated', { key });
    } catch (err) {
      sendJson(res, 500, false, err instanceof Error ? err.message : 'Error');
    }
  };
}
