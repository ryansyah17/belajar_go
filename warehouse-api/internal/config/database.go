package config

import (
	"log"
	"warehouse-api/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),// log semua query SQL ke terminal
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// AutoMigrate — GORM akan CREATE TABLE jika belum ada,
    // atau ALTER TABLE jika ada kolom baru. Tidak menghapus kolom lama.
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
	)
	if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    log.Println("Database migrated successfully")

    return db
} 