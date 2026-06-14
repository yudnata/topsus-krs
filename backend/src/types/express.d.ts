declare global {
  namespace Express {
    interface Request {
      user?: {
        id: string;
        email: string;
        role: string;
        password_hash: string;
        is_active: boolean;
        created_at: Date;
        updated_at: Date;
      };
      userId?: string;
      role?: string;
    }
  }
}

export {};
