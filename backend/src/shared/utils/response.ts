import { Response } from 'express';

export interface ApiResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
}

export function sendJson<T>(
  res: Response,
  status: number,
  success: boolean,
  message: string,
  data?: T,
): Response {
  const body: ApiResponse<T> = { success, message };
  if (data !== undefined) body.data = data;
  return res.status(status).json(body);
}
