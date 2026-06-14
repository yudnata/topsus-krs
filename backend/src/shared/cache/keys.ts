export const TTL_DOSEN_MAHASISWA_LIST = 120;
export const TTL_KRS_PENDING = 60;
export const TTL_KRS_MAHASISWA_LIST = 120;

export function keyKrsPendingDosen(dosenId: string): string {
  return `krs:pending:dosen:${dosenId}`;
}

export function keyKrsPendingStaff(): string {
  return 'krs:pending:staff';
}

export function keyKrsMahasiswaList(mahasiswaId: string): string {
  return `krs:mhs:${mahasiswaId}:list`;
}

export function keyDosenList(): string {
  return 'dosen:list';
}

export function keyDosenDetail(dosenId: string): string {
  return `dosen:${dosenId}`;
}

export function keyMahasiswaList(): string {
  return 'mahasiswa:list';
}

export function keyMahasiswaDetail(mahasiswaId: string): string {
  return `mahasiswa:${mahasiswaId}`;
}
