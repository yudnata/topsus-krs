import { NextFunction, Request, Response } from 'express';
import { sendJson } from '../../shared/utils/response';
import { AuthService } from './auth.service';

export function requireAuth(svc: AuthService) {
  return async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      sendJson(res, 401, false, 'Missing authorization header');
      return;
    }
    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0].toLowerCase() !== 'bearer') {
      sendJson(res, 401, false, 'Invalid authorization format');
      return;
    }
    try {
      const userId = svc.validateToken(parts[1]);
      const user = await svc.findById(userId);
      if (!user) {
        sendJson(res, 401, false, 'User not found');
        return;
      }
      req.user = user;
      req.userId = userId;
      req.role = user.role;
      next();
    } catch {
      sendJson(res, 401, false, 'Invalid or expired token');
    }
  };
}

export function requireRole(...roles: string[]) {
  const allowed = new Set(roles.map((r) => r.toUpperCase()));
  return (req: Request, res: Response, next: NextFunction): void => {
    const role = (req.role ?? '').toUpperCase();
    if (!allowed.has(role)) {
      sendJson(res, 403, false, 'Forbidden');
      return;
    }
    next();
  };
}
