package repository

import (
	"fmt"
	"log"

	"online-quiz/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto Migrate
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Quiz{},
		&domain.Question{},
		&domain.Option{},
		&domain.QuizSession{},
		&domain.SessionQuestion{},
		&domain.SessionOption{},
		&domain.SessionAnswer{},
	)
	if err != nil {
		log.Printf("failed to auto migrate database: %v", err)
		return nil, err
	}

	return db, nil
}
