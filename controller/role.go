package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ulbithebest/todolist-be/model"
	repo "github.com/ulbithebest/todolist-be/repository"
	"gorm.io/gorm"
	"net/http"
)

func GetAllRole(c *fiber.Ctx) error {
	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk mendapatkan semua role
	roles, err := repo.GetAllRole(db)
	if err != nil {
		// Jika terjadi kesalahan saat mengambil data role, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Jika tidak ada role yang ditemukan, mengembalikan pesan kesalahan
	if len(roles) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"code": http.StatusNotFound, "success": false, "status": "error", "message": "Data role tidak ditemukan", "data": nil})
	}

	// Jika tidak ada kesalahan, mengembalikan data role sebagai respons JSON
	response := fiber.Map{
		"code":    http.StatusOK,
		"success": true,
		"status":  "success",
		"data":    roles,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func GetRoleById(c *fiber.Ctx) error {
	// Mendapatkan parameter ID role dari URL
	id := c.Params("id_role")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak ditemukan"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk mendapatkan role berdasarkan ID
	role, err := repo.GetRoleById(db, id)
	if err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Jika tidak ada kesalahan, mengembalikan data role sebagai respons JSON
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "data": role})
}

func InsertRole(c *fiber.Ctx) error {
	// Mendeklarasikan variabel untuk menyimpan data role dari body request
	var role model.Roles

	// Mem-parsing body request ke dalam variabel role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gagal memproses request"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk menyisipkan data role ke dalam database
	if err := repo.InsertRole(db, role); err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan role"})
	}

	// Mengembalikan respons sukses dengan pesan
	return c.Status(http.StatusCreated).JSON(fiber.Map{"code": http.StatusCreated, "success": true, "status": "success", "message": "Buku berhasil disimpan", "data": role})
}

func UpdateRole(c *fiber.Ctx) error {
	// Mendapatkan parameter ID
	id := c.Params("id_role")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak ditemukan"})
	}

	// Mendeklarasikan variabel untuk menyimpan data role yang diperbarui dari body request
	var updatedRole model.Roles

	// Mem-parsing body request ke dalam variabel updatedRole
	if err := c.BodyParser(&updatedRole); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gagal memproses request"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk memperbarui data role di dalam database
	if err := repo.UpdateRole(db, id, updatedRole); err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui role"})
	}

	// Mengembalikan respons sukses dengan pesan
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "message": "Role berhasil diperbarui"})
}

func DeleteRole(c *fiber.Ctx) error {
	// Mendapatkan parameter ID role dari URL
	id := c.Params("id_role")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak ditemukan"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk menghapus data role dari database berdasarkan ID
	if err := repo.DeleteRole(db, id); err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus role"})
	}

	// Mengembalikan respons sukses
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "message": "Role berhasil dihapus", "deleted_id": id})
}
