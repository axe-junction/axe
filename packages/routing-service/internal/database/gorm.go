package database

import (
	"log"

	"github.com/axe-junction/axe-server/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	cfg := config.Get()
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Println("Error enabling uuid extension:", err)
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "postgis"`).Error; err != nil {
		log.Println("Error enabling uuid extension:", err)
	}
	log.Println("Connected to PostgreSQL using GORM")
	return db, nil
}
