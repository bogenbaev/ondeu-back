package repository

import (
	"fmt"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
	SSLMode  string
}

func NewPostgresRepository(cfg Config) *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(connStr), nil)
	if err != nil {
		log.Fatalf("an error is occured while connecting: %s", err.Error())
	}

	if err = AutoMigrate(db); err != nil {
		log.Fatalf("an error is occurred while migrating: %s", err.Error())
	}

	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&dto.Tree{},
		&dto.Document{},
	)
}
