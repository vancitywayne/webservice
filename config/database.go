package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func CreateDBConnection() *gorm.DB {
	// Memuat file .env
	LoadEnv()

	// Memuat string connection database dari variabel environment
	dbConfig := os.Getenv("SQLSTRING")

	// Membuat koneksi database
	DB, err := gorm.Open(mysql.Open(dbConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// Mengatur konfigurasi kumpulan koneksi database
	dbO, err := DB.DB()
	if err != nil {
		panic(err)
	}
	dbO.SetConnMaxIdleTime(time.Duration(1) * time.Minute)
	dbO.SetMaxIdleConns(2)
	dbO.SetConnMaxLifetime(time.Duration(1) * time.Minute)

	// Munculkan kesalahan jika record tidak ditemukan
	DB.Statement.RaiseErrorOnNotFound = true

	return DB
}
