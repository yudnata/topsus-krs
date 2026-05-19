-- Seed minimal — development (idempotent dengan ON CONFLICT)
-- Password untuk semua akun dev: password123
-- Hash bcrypt (cost 10): $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

-- ─── Akun Users ───────────────────────────────────────────────────────────────
INSERT INTO users (id, email, password_hash, role, is_active, created_at, updated_at) VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'admin@kampus.ac.id',    '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'ADMIN',     true, NOW(), NOW()),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01', 'ahmad@kampus.ac.id',    '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'DOSEN',     true, NOW(), NOW()),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02', 'siti@kampus.ac.id',     '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'DOSEN',     true, NOW(), NOW()),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb03', 'budi@kampus.ac.id',     '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'DOSEN',     true, NOW(), NOW()),
  ('cccccccc-cccc-cccc-cccc-cccccccccc01', 'andi@student.ac.id',    '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'MAHASISWA', true, NOW(), NOW()),
  ('cccccccc-cccc-cccc-cccc-cccccccccc02', 'bella@student.ac.id',   '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'MAHASISWA', true, NOW(), NOW()),
  ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'staff@kampus.ac.id',    '$2a$10$bn3faSzRRDLifYqLe30y/eeWLfLU4TT9ufYRz7bP0HuTd81xMox0S', 'STAFF',     true, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;


INSERT INTO program_studi (id, kode_prodi, nama_prodi, fakultas) VALUES
  ('11111111-1111-1111-1111-111111111101', 'TI', 'Teknik Informatika', 'FTI'),
  ('11111111-1111-1111-1111-111111111102', 'SI', 'Sistem Informasi', 'FTI')
ON CONFLICT (kode_prodi) DO NOTHING;

INSERT INTO tahun_akademik (id, kode_ta, semester, is_active) VALUES
  ('22222222-2222-2222-2222-222222222201', '2025/2026', 'Ganjil', true)
ON CONFLICT (kode_ta) DO NOTHING;

INSERT INTO dosen (id, user_id, nip, nama, email) VALUES
  ('33333333-3333-3333-3333-333333333301', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01', '19800101001', 'Dr. Ahmad Wijaya', 'ahmad@kampus.ac.id'),
  ('33333333-3333-3333-3333-333333333302', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02', '19800202002', 'Prof. Siti Rahayu', 'siti@kampus.ac.id'),
  ('33333333-3333-3333-3333-333333333303', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb03', '19800303003', 'Dr. Budi Santoso', 'budi@kampus.ac.id')
ON CONFLICT (nip) DO NOTHING;

INSERT INTO mahasiswa (id, user_id, prodi_id, nim, nama, dosen_pembimbing_id, max_sks) VALUES
  ('44444444-4444-4444-4444-444444444401', 'cccccccc-cccc-cccc-cccc-cccccccccc01', '11111111-1111-1111-1111-111111111101', '2201001', 'Andi Pratama', '33333333-3333-3333-3333-333333333301', 24),
  ('44444444-4444-4444-4444-444444444402', 'cccccccc-cccc-cccc-cccc-cccccccccc02', '11111111-1111-1111-1111-111111111101', '2201002', 'Bella Sari', '33333333-3333-3333-3333-333333333302', 24)
ON CONFLICT (nim) DO NOTHING;

INSERT INTO mata_kuliah (id, prodi_id, kode_mk, nama_mk, sks) VALUES
  ('55555555-5555-5555-5555-555555555501', '11111111-1111-1111-1111-111111111101', 'IF101', 'Algoritma & Pemrograman', 3),
  ('55555555-5555-5555-5555-555555555502', '11111111-1111-1111-1111-111111111101', 'IF102', 'Struktur Data', 3),
  ('55555555-5555-5555-5555-555555555503', '11111111-1111-1111-1111-111111111101', 'IF201', 'Basis Data', 3),
  ('55555555-5555-5555-5555-555555555504', '11111111-1111-1111-1111-111111111102', 'SI101', 'Pengantar Sistem Informasi', 3),
  ('55555555-5555-5555-5555-555555555505', '11111111-1111-1111-1111-111111111101', 'IF301', 'Pemrograman Web', 3)
ON CONFLICT (kode_mk) DO NOTHING;

INSERT INTO prasyarat_mk (mata_kuliah_id, prasyarat_id)
SELECT '55555555-5555-5555-5555-555555555502', '55555555-5555-5555-5555-555555555501'
WHERE NOT EXISTS (
  SELECT 1 FROM prasyarat_mk
  WHERE mata_kuliah_id = '55555555-5555-5555-5555-555555555502'
    AND prasyarat_id = '55555555-5555-5555-5555-555555555501'
);

INSERT INTO kelas (id, mata_kuliah_id, tahun_akademik_id, dosen_id, nama_kelas, kapasitas, terisi) VALUES
  ('66666666-6666-6666-6666-666666666601', '55555555-5555-5555-5555-555555555501', '22222222-2222-2222-2222-222222222201', '33333333-3333-3333-3333-333333333301', 'A', 40, 0),
  ('66666666-6666-6666-6666-666666666602', '55555555-5555-5555-5555-555555555503', '22222222-2222-2222-2222-222222222201', '33333333-3333-3333-3333-333333333302', 'A', 35, 0),
  ('66666666-6666-6666-6666-666666666603', '55555555-5555-5555-5555-555555555505', '22222222-2222-2222-2222-222222222201', '33333333-3333-3333-3333-333333333303', 'B', 30, 0)
ON CONFLICT (mata_kuliah_id, tahun_akademik_id, nama_kelas) DO NOTHING;

INSERT INTO jadwal_kelas (kelas_id, hari, jam_mulai, jam_selesai, ruangan)
SELECT v.kelas_id, v.hari, v.jam_mulai::time, v.jam_selesai::time, v.ruangan
FROM (VALUES
  ('66666666-6666-6666-6666-666666666601'::uuid, 'Senin', '08:00', '10:00', 'Lab-1'),
  ('66666666-6666-6666-6666-666666666602'::uuid, 'Selasa', '10:00', '12:00', 'R-201'),
  ('66666666-6666-6666-6666-666666666603'::uuid, 'Rabu', '13:00', '15:00', 'R-305')
) AS v(kelas_id, hari, jam_mulai, jam_selesai, ruangan)
WHERE NOT EXISTS (
  SELECT 1 FROM jadwal_kelas j WHERE j.kelas_id = v.kelas_id AND j.hari = v.hari
);
