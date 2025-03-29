package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"roulette/internal/config"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBConfig.Host, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Name, cfg.DBConfig.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	_, err = db.DB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
