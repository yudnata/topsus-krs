import { Request, Response } from 'express';
import { sendJson } from '../../shared/utils/response';
import { toPublicUser } from './auth.types';
import { AuthService } from './auth.service';

export class AuthController {
  constructor(private svc: AuthService) {}

  register = async (req: Request, res: Response): Promise<void> => {
    try {
      const user = await this.svc.register(req.body);
      sendJson(res, 201, true, 'Registered', user);
    } catch (err) {
      sendJson(res, 409, false, err instanceof Error ? err.message : 'Error');
    }
  };

  login = async (req: Request, res: Response): Promise<void> => {
    try {
      const result = await this.svc.login(req.body);
      sendJson(res, 200, true, 'Success', result);
    } catch {
      sendJson(res, 401, false, 'invalid credentials');
    }
  };

  profile = (req: Request, res: Response): void => {
    if (!req.user) {
      sendJson(res, 401, false, 'Unauthorized');
      return;
    }
    sendJson(res, 200, true, 'Success', toPublicUser(req.user));
  };
}
