-- Skema Sistem KRS (11 tabel) — sumber: ARSITEKTUR.md §5

DROP TABLE IF EXISTS pengajuan_krs_detail CASCADE;
DROP TABLE IF EXISTS pengajuan_krs CASCADE;
DROP TABLE IF EXISTS jadwal_kelas CASCADE;
DROP TABLE IF EXISTS kelas CASCADE;
DROP TABLE IF EXISTS prasyarat_mk CASCADE;
DROP TABLE IF EXISTS mata_kuliah CASCADE;
DROP TABLE IF EXISTS mahasiswa CASCADE;
DROP TABLE IF EXISTS dosen CASCADE;
DROP TABLE IF EXISTS tahun_akademik CASCADE;
DROP TABLE IF EXISTS program_studi CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email         VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role          VARCHAR(20) NOT NULL CHECK (role IN ('ADMIN','MAHASISWA','DOSEN','STAFF')),
  is_active     BOOLEAN NOT NULL DEFAULT true,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE program_studi (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  kode_prodi VARCHAR(10) UNIQUE NOT NULL,
  nama_prodi VARCHAR(255) NOT NULL,
  fakultas   VARCHAR(255) NOT NULL
);

CREATE TABLE tahun_akademik (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  kode_ta    VARCHAR(10) UNIQUE NOT NULL,
  semester   VARCHAR(10) NOT NULL CHECK (semester IN ('Ganjil','Genap','Pendek')),
  is_active  BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE dosen (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID UNIQUE REFERENCES users(id) ON DELETE SET NULL,
  nip        VARCHAR(32) UNIQUE NOT NULL,
  nama       VARCHAR(255) NOT NULL,
  email      VARCHAR(255) UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mahasiswa (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id             UUID UNIQUE REFERENCES users(id) ON DELETE SET NULL,
  prodi_id            UUID REFERENCES program_studi(id) ON DELETE RESTRICT,
  nim                 VARCHAR(32) UNIQUE NOT NULL,
  nama                VARCHAR(255) NOT NULL,
  dosen_pembimbing_id UUID REFERENCES dosen(id) ON DELETE SET NULL,
  max_sks             SMALLINT NOT NULL DEFAULT 24,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mata_kuliah (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  prodi_id   UUID REFERENCES program_studi(id),
  kode_mk    VARCHAR(20) UNIQUE NOT NULL,
  nama_mk    VARCHAR(255) NOT NULL,
  sks        SMALLINT NOT NULL CHECK (sks > 0)
);

CREATE TABLE prasyarat_mk (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mata_kuliah_id UUID REFERENCES mata_kuliah(id) ON DELETE CASCADE,
  prasyarat_id   UUID REFERENCES mata_kuliah(id) ON DELETE CASCADE,
  UNIQUE (mata_kuliah_id, prasyarat_id)
);

CREATE TABLE kelas (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mata_kuliah_id    UUID REFERENCES mata_kuliah(id) ON DELETE CASCADE,
  tahun_akademik_id UUID REFERENCES tahun_akademik(id) ON DELETE CASCADE,
  dosen_id          UUID REFERENCES dosen(id) ON DELETE SET NULL,
  nama_kelas        VARCHAR(10) NOT NULL,
  kapasitas         SMALLINT NOT NULL DEFAULT 40,
  terisi            SMALLINT NOT NULL DEFAULT 0,
  UNIQUE (mata_kuliah_id, tahun_akademik_id, nama_kelas)
);

CREATE TABLE jadwal_kelas (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  kelas_id    UUID REFERENCES kelas(id) ON DELETE CASCADE,
  hari        VARCHAR(10) NOT NULL,
  jam_mulai   TIME NOT NULL,
  jam_selesai TIME NOT NULL,
  ruangan     VARCHAR(50) NOT NULL
);

CREATE TABLE pengajuan_krs (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mahasiswa_id      UUID NOT NULL REFERENCES mahasiswa(id) ON DELETE CASCADE,
  tahun_akademik_id UUID NOT NULL REFERENCES tahun_akademik(id) ON DELETE RESTRICT,
  status            VARCHAR(20) NOT NULL DEFAULT 'draft'
    CHECK (status IN ('draft','diajukan','disetujui','ditolak')),
  total_sks         SMALLINT NOT NULL DEFAULT 0,
  catatan_mhs       TEXT,
  catatan_reviewer  TEXT,
  reviewed_by       UUID REFERENCES users(id) ON DELETE SET NULL,
  reviewed_at       TIMESTAMPTZ,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (mahasiswa_id, tahun_akademik_id)
);

CREATE TABLE pengajuan_krs_detail (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pengajuan_krs_id UUID NOT NULL REFERENCES pengajuan_krs(id) ON DELETE CASCADE,
  kelas_id         UUID NOT NULL REFERENCES kelas(id) ON DELETE RESTRICT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (pengajuan_krs_id, kelas_id)
);

CREATE INDEX idx_mahasiswa_dpa ON mahasiswa(dosen_pembimbing_id);
CREATE INDEX idx_pengajuan_krs_status ON pengajuan_krs(status);
CREATE INDEX idx_kelas_ta ON kelas(tahun_akademik_id);
