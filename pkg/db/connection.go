package db

import (
	"github.com/caarlos0/env/v6"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type config struct {
	DSN string `env:"DATASOURCENAME" envDefault:"file::memory:?cache=shared"`
}

func NewGormClient() (*gorm.DB, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
