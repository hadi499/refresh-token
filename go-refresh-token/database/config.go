package database

import (
	"fmt"
	"log"

	"go-refresh-token/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Ganti dengan kredensial PostgreSQL yang sesuai
	dsn := "host=localhost user=hadi password=admin123 dbname=gotoken_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	// Membuka koneksi ke database dengan konfigurasi logger aktif
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ Gagal terhubung ke database:", err)
	}

	// Migrasi otomatis (pastikan model sudah benar)
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Gagal melakukan migrasi database:", err)
	}

	fmt.Println("✅ Database connected successfully!")
}

// GetDB mengembalikan instance database
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("❌ Database belum diinisialisasi. Panggil ConnectDatabase() terlebih dahulu.")
	}
	return DB
}
