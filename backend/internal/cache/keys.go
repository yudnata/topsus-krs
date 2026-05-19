package cache

import "fmt"

// TTL default (detik) — selaras ARSITEKTUR.md §9
const (
	TTLDosenMahasiswaList = 120
	TTLKrsPending         = 60
	TTLKrsMahasiswaList   = 120
)

// KRS / antrian persetujuan
func KeyKrsPendingDosen(dosenID string) string {
	return fmt.Sprintf("krs:pending:dosen:%s", dosenID)
}

func KeyKrsPendingStaff() string {
	return "krs:pending:staff"
}

func KeyKrsMahasiswaList(mahasiswaID string) string {
	return fmt.Sprintf("krs:mhs:%s:list", mahasiswaID)
}

// Master data (invalidate on CRUD / PATCH DPA — Slice 4)
func KeyDosenList() string {
	return "dosen:list"
}

func KeyDosenDetail(dosenID string) string {
	return fmt.Sprintf("dosen:%s", dosenID)
}

func KeyMahasiswaList() string {
	return "mahasiswa:list"
}

func KeyMahasiswaDetail(mahasiswaID string) string {
	return fmt.Sprintf("mahasiswa:%s", mahasiswaID)
}
