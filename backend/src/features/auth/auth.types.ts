export interface User {
  id: string;
  email: string;
  role: string;
  password_hash: string;
  is_active: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface RegisterReq {
  email: string;
  password: string;
  role?: string;
}

export interface LoginReq {
  email: string;
  password: string;
}

export interface UserPublic {
  id: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: Date;
  updated_at: Date;
}

export function toPublicUser(u: User): UserPublic {
  return {
    id: u.id,
    email: u.email,
    role: u.role,
    is_active: u.is_active,
    created_at: u.created_at,
    updated_at: u.updated_at,
  };
}
