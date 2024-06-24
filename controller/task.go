package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/ulbithebest/todolist-be/model"
	repo "github.com/ulbithebest/todolist-be/repository"
	"gorm.io/gorm"
	"net/http"
)

func GetAllTask(c *fiber.Ctx) error {
	// Cek token header autentikasi
	token := c.Get("login")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk mendapatkan semua task
	tasks, err := repo.GetAllTask(db)
	if err != nil {
		// Jika terjadi kesalahan saat mengambil data task, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Terjadi Kesalahan dalam Get Data"})
	}

	// Jika tidak ada task yang ditemukan, mengembalikan pesan kesalahan
	if len(tasks) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"code": http.StatusNotFound, "success": false, "status": "error", "message": "Data task tidak ditemukan", "data": nil})
	}

	// Jika tidak ada kesalahan, mengembalikan data task sebagai respons JSON
	response := fiber.Map{
		"code":    http.StatusOK,
		"success": true,
		"status":  "success",
		"data":    tasks,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func GetTaskById(c *fiber.Ctx) error {
	// Cek token header autentikasi
	token := c.Get("login")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	// Mendapatkan parameter ID task dari URL
	id := c.Query("id_task")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID task tidak boleh kosong"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk mendapatkan task berdasarkan ID
	task, err := repo.GetTaskById(db, id)
	if err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Data Tidak Ditemukan"})
	}

	// Jika tidak ada kesalahan, mengembalikan data task sebagai respons JSON
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "data": task})
}

func GetTaskByIdUser(c *fiber.Ctx) error {
	// Cek token header autentikasi
	tokenStr := c.Get("login")
	if tokenStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	// Parse token untuk mendapatkan id_user
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil // Ganti "secret_key" dengan kunci rahasia Anda
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	idUser := claims.IdUser // Dapatkan id_user dari klaim token

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk mendapatkan task berdasarkan ID
	task, err := repo.GetTaskByIdUser(db, int(idUser))
	if err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Data Tidak Ditemukan"})
	}

	// Jika tidak ada kesalahan, mengembalikan data task sebagai respons JSON
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "data": task})
}

func InsertTask(c *fiber.Ctx) error {
	// Cek token header autentikasi
	tokenStr := c.Get("login")
	if tokenStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	// Parse token untuk mendapatkan id_user
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil // Ganti "secret_key" dengan kunci rahasia Anda
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	idUser := claims.IdUser

	// Mendeklarasikan variabel untuk menyimpan data task dari body request
	var task model.Task

	// Mem-parsing body request ke dalam variabel task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gagal memproses request"})
	}

	task.IdUser = int(idUser)
	task.Completed = "false"

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk insert data task ke dalam database
	if err := repo.InsertTask(db, &task); err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan task"})
	}

	// Mengembalikan respons sukses dengan pesan
	return c.Status(http.StatusCreated).JSON(fiber.Map{"code": http.StatusCreated, "success": true, "status": "success", "message": "Task berhasil disimpan", "data": task})
}

func UpdateTask(c *fiber.Ctx) error {
	// Cek token header autentikasi
	token := c.Get("login")
	if token == "" {
		return fiber.NewError(http.StatusBadRequest, "Header tidak ada")
	}

	// Mendapatkan parameter ID
	id := c.Query("id_task")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID task tidak ditemukan"})
	}

	// Mendeklarasikan variabel untuk menyimpan data task yang diperbarui dari body request
	var updatedTask model.Task

	// Mem-parsing body request ke dalam variabel updatedTask
	if err := c.BodyParser(&updatedTask); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gagal memproses request"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	// Memanggil fungsi repo untuk cek data dengan id apakah ada atau tidak
	_, err := repo.GetTaskById(db, id)
	if err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "ID task tidak ditemukan"})
	}

	// Memanggil fungsi repo untuk memperbarui data task di dalam database
	if err := repo.UpdateTask(db, id, updatedTask); err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui task"})
	}

	// Mengembalikan respons sukses dengan pesan
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "message": "Task berhasil diperbarui"})
}

func DeleteTask(c *fiber.Ctx) error {
	// Cek token header autentikasi
	token := c.Get("login")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	// Mendapatkan parameter ID task dari URL
	id := c.Query("id_task")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID task tidak boleh kosong"})
	}

	// Mendapatkan koneksi database dari context Fiber
	db := c.Locals("db").(*gorm.DB)

	//  Memanggil fungsi repo untuk cek data dengan id apakah ada atau tidak
	_, err := repo.GetTaskById(db, id)
	if err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "ID task tidak ditemukan"})
	}

	// Memanggil fungsi repo untuk menghapus data task dari database berdasarkan ID
	if err := repo.DeleteTask(db, id); err != nil {
		// Jika terjadi kesalahan, mengembalikan respons error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus task"})
	}

	// Mengembalikan respons sukses
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "message": "Task berhasil dihapus", "deleted_id": id})
}
