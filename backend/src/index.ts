import { createApp } from './app';
import { loadConfig } from './config';
import { CacheClient, Invalidator } from './shared/cache/redis';
import { connect, migrate } from './shared/database/postgres';

async function main(): Promise<void> {
  const cfg = loadConfig();
  const db = connect(cfg.databaseUrl);
  await migrate(db);

  const cache = await CacheClient.create(cfg.redisUrl);
  const inv = new Invalidator(cache);

  const app = createApp({ cfg, db, cache, inv });

  const server = app.listen(cfg.port, () => {
    console.log(`✓ Server listening on :${cfg.port}`);
  });

  const shutdown = async (): Promise<void> => {
    server.close();
    await cache.close();
    await db.end();
    process.exit(0);
  };

  process.on('SIGINT', () => void shutdown());
  process.on('SIGTERM', () => void shutdown());
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
