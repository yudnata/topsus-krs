import cors from 'cors';
import express from 'express';
import morgan from 'morgan';
import { splitOrigins } from './config';
import { setupRoutes } from './routes';
import { Config } from './config';

export function createApp(deps: {
  cfg: Config;
  db: import('pg').Pool;
  cache: import('./shared/cache/redis').CacheClient;
  inv: import('./shared/cache/redis').Invalidator;
}): express.Application {
  const app = express();

  app.use(morgan('dev'));
  app.use(
    cors({
      origin: splitOrigins(deps.cfg.corsOrigin),
      allowedHeaders: ['Origin', 'Content-Type', 'Accept', 'Authorization'],
    }),
  );
  app.use(express.json());

  setupRoutes(app, deps);

  return app;
}
