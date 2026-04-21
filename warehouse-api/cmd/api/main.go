package main

import (
	"log"

	"warehouse-api/internal/config"
	"warehouse-api/internal/router"
)

func main() {
    // 1. Load config
    cfg := config.Load()

    // 2. Connect & migrate database
    db := config.NewDatabase(cfg)

    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal(err)
    }
    defer sqlDB.Close()

    // 3. Setup & jalankan router
    r := router.New(db)
    r.Setup()

    log.Printf("Server running → http://localhost:%s/api/v1/health", cfg.AppPort)
    log.Fatal(r.Run(cfg.AppPort))
}