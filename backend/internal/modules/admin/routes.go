package admin

import (
	"backend/internal/modules/auth"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes mendaftarkan semua endpoint admin di bawah /api/admin.
// Semua dilindungi JWT + role ADMIN.
func RegisterRoutes(api fiber.Router, h *Handler, authSvc *auth.Service) {
	admin := api.Group("/admin", auth.RequireAuth(authSvc))

	// Referensi
	admin.Get("/prodi", auth.RequireRole("ADMIN"), h.ListProdi)
	admin.Get("/tahun-akademik", auth.RequireRole("ADMIN"), h.ListTahunAkademik)
	admin.Get("/tahun-akademik/aktif", auth.RequireRole("ADMIN", "MAHASISWA", "DOSEN", "STAFF"), h.GetTahunAkademikAktif)

	// Mata Kuliah
	admin.Get("/mata-kuliah", auth.RequireRole("ADMIN"), h.ListMataKuliah)
	admin.Post("/mata-kuliah", auth.RequireRole("ADMIN"), h.CreateMataKuliah)

	// Kelas (filter by ?tahun_akademik_id=)
	admin.Get("/kelas", auth.RequireRole("ADMIN", "MAHASISWA"), h.ListKelas)
	admin.Post("/kelas", auth.RequireRole("ADMIN"), h.CreateKelas)

	// Dosen
	admin.Get("/dosen", auth.RequireRole("ADMIN"), h.ListDosen)
	admin.Get("/dosen/:id", auth.RequireRole("ADMIN"), h.GetDosen)
	admin.Post("/dosen", auth.RequireRole("ADMIN"), h.CreateDosen)

	// Mahasiswa
	admin.Get("/mahasiswa", auth.RequireRole("ADMIN"), h.ListMahasiswa)
	admin.Get("/mahasiswa/:id", auth.RequireRole("ADMIN"), h.GetMahasiswa)
	admin.Post("/mahasiswa", auth.RequireRole("ADMIN"), h.CreateMahasiswa)
	admin.Patch("/mahasiswa/:id/dpa", auth.RequireRole("ADMIN"), h.PatchDPA)
}
