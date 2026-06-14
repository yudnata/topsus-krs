import { Router } from 'express';
import { requireAuth, requireRole } from '../auth/auth.middleware';
import { AuthService } from '../auth/auth.service';
import { AdminController } from './admin.controller';
import { AdminService } from './admin.service';

export function createAdminRoutes(svc: AdminService, authSvc: AuthService): Router {
  const router = Router();
  const ctrl = new AdminController(svc);

  router.use(requireAuth(authSvc));

  router.get('/prodi', requireRole('ADMIN'), ctrl.listProdi);
  router.get('/tahun-akademik', requireRole('ADMIN'), ctrl.listTahunAkademik);
  router.get('/tahun-akademik/aktif', requireRole('ADMIN', 'MAHASISWA', 'DOSEN', 'STAFF'), ctrl.getTahunAkademikAktif);
  router.get('/mata-kuliah', requireRole('ADMIN'), ctrl.listMataKuliah);
  router.post('/mata-kuliah', requireRole('ADMIN'), ctrl.createMataKuliah);
  router.get('/kelas', requireRole('ADMIN', 'MAHASISWA'), ctrl.listKelas);
  router.post('/kelas', requireRole('ADMIN'), ctrl.createKelas);
  router.get('/dosen', requireRole('ADMIN'), ctrl.listDosen);
  router.get('/dosen/:id', requireRole('ADMIN'), ctrl.getDosen);
  router.post('/dosen', requireRole('ADMIN'), ctrl.createDosen);
  router.get('/mahasiswa', requireRole('ADMIN'), ctrl.listMahasiswa);
  router.get('/mahasiswa/:id', requireRole('ADMIN'), ctrl.getMahasiswa);
  router.post('/mahasiswa', requireRole('ADMIN'), ctrl.createMahasiswa);
  router.patch('/mahasiswa/:id/dpa', requireRole('ADMIN'), ctrl.patchDPA);

  return router;
}
