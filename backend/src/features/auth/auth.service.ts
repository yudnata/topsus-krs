import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import { v4 as uuidv4 } from 'uuid';
import { Config } from '../../config';
import { AuthRepository } from './auth.repository';
import { LoginReq, RegisterReq, toPublicUser, User, UserPublic } from './auth.types';

const ALLOWED_ROLES = new Set(['ADMIN', 'MAHASISWA', 'DOSEN', 'STAFF']);

export class AuthService {
  constructor(
    private repo: AuthRepository,
    private cfg: Config,
  ) {}

  async register(input: RegisterReq): Promise<UserPublic> {
    let role = (input.role ?? 'MAHASISWA').trim().toUpperCase();
    if (!ALLOWED_ROLES.has(role)) throw new Error('invalid role');

    const hash = await bcrypt.hash(input.password, 10);
    const now = new Date();
    const user: User = {
      id: uuidv4(),
      email: input.email.trim(),
      role,
      password_hash: hash,
      is_active: true,
      created_at: now,
      updated_at: now,
    };
    await this.repo.create(user);
    return toPublicUser(user);
  }

  async login(input: LoginReq): Promise<{ token: string; user: UserPublic }> {
    const user = await this.repo.findByEmail(input.email.trim());
    if (!user || !user.is_active) throw new Error('invalid credentials');
    const ok = await bcrypt.compare(input.password, user.password_hash);
    if (!ok) throw new Error('invalid credentials');
    const token = this.issueToken(user);
    return { token, user: toPublicUser(user) };
  }

  validateToken(tokenStr: string): string {
    try {
      const payload = jwt.verify(tokenStr, this.cfg.jwtSecret) as jwt.JwtPayload;
      const sub = payload.sub;
      if (typeof sub !== 'string' || !sub) throw new Error('invalid subject');
      return sub;
    } catch {
      throw new Error('invalid or expired token');
    }
  }

  async findById(id: string): Promise<User | null> {
    return this.repo.findById(id);
  }

  private issueToken(user: User): string {
    return jwt.sign({ sub: user.id, role: user.role }, this.cfg.jwtSecret, { expiresIn: '24h' });
  }
}
