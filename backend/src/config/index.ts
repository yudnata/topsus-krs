import dotenv from 'dotenv';

dotenv.config();

export interface Config {
  port: string;
  databaseUrl: string;
  redisUrl: string;
  jwtSecret: string;
  corsOrigin: string;
  mode: string;
}

function getEnv(key: string, fallback = ''): string {
  return process.env[key]?.trim() || fallback;
}

export function loadConfig(): Config {
  return {
    port: getEnv('PORT', '8080'),
    databaseUrl: getEnv('DATABASE_URL'),
    redisUrl: getEnv('REDIS_URL', 'redis://localhost:6379/0'),
    jwtSecret: getEnv('JWT_SECRET', 'secret'),
    corsOrigin: getEnv('CORS_ORIGIN', 'http://localhost:5173'),
    mode: getEnv('MODE', 'development'),
  };
}

export function splitOrigins(origins: string): string[] {
  const parts = origins
    .split(',')
    .map((o) => o.trim())
    .filter(Boolean);
  return parts.length > 0 ? parts : ['http://localhost:5173'];
}
