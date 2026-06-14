import fs from 'fs';
import path from 'path';
import { Pool } from 'pg';

let pool: Pool | null = null;

export function connect(databaseUrl: string): Pool {
  if (!databaseUrl) {
    console.error('DATABASE_URL is required');
    process.exit(1);
  }
  pool = new Pool({ connectionString: databaseUrl });
  return pool;
}

export function getPool(): Pool {
  if (!pool) throw new Error('Database not initialized');
  return pool;
}

function migrationDir(): string {
  const candidates = [
    path.join(process.cwd(), 'migrations'),
    path.join(process.cwd(), 'backend', 'migrations'),
  ];
  for (const dir of candidates) {
    if (fs.existsSync(dir) && fs.statSync(dir).isDirectory()) return dir;
  }
  console.error('migrations directory not found (run from backend/ or project root)');
  process.exit(1);
}

async function ensureMigrationTable(client: Pool): Promise<void> {
  await client.query(`
    CREATE TABLE IF NOT EXISTS schema_migrations (
      filename   VARCHAR(255) PRIMARY KEY,
      applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    )`);
}

async function isMigrationApplied(client: Pool, name: string): Promise<boolean> {
  const res = await client.query(
    'SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)',
    [name],
  );
  return res.rows[0]?.exists === true;
}

async function markMigrationApplied(client: Pool, name: string): Promise<void> {
  await client.query(
    'INSERT INTO schema_migrations (filename) VALUES ($1) ON CONFLICT DO NOTHING',
    [name],
  );
}

export async function migrate(client: Pool): Promise<void> {
  const dir = migrationDir();
  const names = fs
    .readdirSync(dir)
    .filter((f) => f.endsWith('.sql'))
    .sort();

  await ensureMigrationTable(client);

  for (const name of names) {
    if (await isMigrationApplied(client, name)) {
      console.log(`⊘ Migration skipped (already applied): ${name}`);
      continue;
    }
    const sql = fs.readFileSync(path.join(dir, name), 'utf8');
    await client.query(sql);
    await markMigrationApplied(client, name);
    console.log(`✓ Migration applied: ${name}`);
  }
  console.log('✓ Database migrate check complete');
}
