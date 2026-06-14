import { createClient, RedisClientType } from 'redis';
import {
  keyDosenDetail,
  keyDosenList,
  keyKrsMahasiswaList,
  keyKrsPendingDosen,
  keyKrsPendingStaff,
  keyMahasiswaDetail,
  keyMahasiswaList,
} from './keys';

export class CacheClient {
  private client: RedisClientType | null = null;
  private enabled = false;

  static async create(redisUrl: string): Promise<CacheClient> {
    const instance = new CacheClient();
    if (!redisUrl) {
      console.log('⚠ REDIS_URL kosong — cache dinonaktifkan (fallback DB only)');
      return instance;
    }
    try {
      const client = createClient({
        url: redisUrl,
        socket: { connectTimeout: 3000, reconnectStrategy: () => false },
      });
      client.on('error', () => {});
      await Promise.race([
        client.connect(),
        new Promise((_, reject) => setTimeout(() => reject(new Error('timeout')), 3000)),
      ]);
      await client.ping();
      instance.client = client as RedisClientType;
      instance.enabled = true;
      console.log('✓ Redis connected');
    } catch (err) {
      console.log(`⚠ Redis tidak tersedia: ${err} — fallback DB only`);
    }
    return instance;
  }

  isEnabled(): boolean {
    return this.enabled;
  }

  async ping(): Promise<void> {
    if (!this.client) throw new Error('redis disabled');
    await this.client.ping();
  }

  async get(key: string): Promise<{ value: string; hit: boolean }> {
    if (!this.enabled || !this.client) return { value: '', hit: false };
    const val = await this.client.get(key);
    if (val === null) return { value: '', hit: false };
    return { value: val, hit: true };
  }

  async set(key: string, value: string, ttlSeconds: number): Promise<void> {
    if (!this.enabled || !this.client) return;
    await this.client.set(key, value, { EX: ttlSeconds });
  }

  async del(...keys: string[]): Promise<void> {
    if (!this.enabled || !this.client || keys.length === 0) return;
    await this.client.del(keys);
  }

  async close(): Promise<void> {
    if (this.client) await this.client.quit();
  }
}

export class Invalidator {
  constructor(private cache: CacheClient) {}

  async onDosenWrite(dosenId: string): Promise<void> {
    await this.cache.del(keyDosenList(), keyDosenDetail(dosenId), keyKrsPendingDosen(dosenId));
  }

  async onMahasiswaWrite(mahasiswaId: string, dosenPembimbingId: string): Promise<void> {
    const keys = [keyMahasiswaList(), keyMahasiswaDetail(mahasiswaId), keyKrsMahasiswaList(mahasiswaId)];
    if (dosenPembimbingId) keys.push(keyKrsPendingDosen(dosenPembimbingId));
    await this.cache.del(...keys);
  }

  async onKrsMutation(mahasiswaId: string, dosenPembimbingId: string): Promise<void> {
    const keys = [keyKrsPendingStaff(), keyKrsMahasiswaList(mahasiswaId)];
    if (dosenPembimbingId) keys.push(keyKrsPendingDosen(dosenPembimbingId));
    await this.cache.del(...keys);
  }

  async onDpaChange(oldDosenId: string, newDosenId: string, mahasiswaId: string): Promise<void> {
    await this.onMahasiswaWrite(mahasiswaId, newDosenId);
    if (oldDosenId && oldDosenId !== newDosenId) {
      await this.cache.del(keyKrsPendingDosen(oldDosenId));
    }
  }
}
