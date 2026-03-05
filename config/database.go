package config

import (
	"ecommerce/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", // sslmode adalah
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	// buka koneksi ke database dengan gorm dan postgres driver
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("GAGAL KONEKSI KE DATABASE!", err)
	}

	// migrasi otomatis untuk semua model
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.BlacklistedToken{},
	)
	if err != nil {
		log.Fatal("migrate gagal!:", err)
	}

	// set global variable DB dengan koneksi database yang sudah berhasil
	DB = db
	fmt.Println("BERHASIL KONEKSI KE DATABASE!")
}
