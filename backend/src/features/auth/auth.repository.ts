import { Pool } from 'pg';
import { User } from './auth.types';

export class AuthRepository {
  constructor(private db: Pool) {}

  async create(u: User): Promise<void> {
    await this.db.query(
      `INSERT INTO users (id, email, password_hash, role, is_active, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6, $7)`,
      [u.id, u.email, u.password_hash, u.role, u.is_active, u.created_at, u.updated_at],
    );
  }

  async findByEmail(email: string): Promise<User | null> {
    const res = await this.db.query(
      `SELECT id, email, role, password_hash, is_active, created_at, updated_at
       FROM users WHERE email = $1`,
      [email],
    );
    return res.rows[0] ?? null;
  }

  async findById(id: string): Promise<User | null> {
    const res = await this.db.query(
      `SELECT id, email, role, password_hash, is_active, created_at, updated_at
       FROM users WHERE id = $1`,
      [id],
    );
    return res.rows[0] ?? null;
  }
}
