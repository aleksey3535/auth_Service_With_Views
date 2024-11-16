package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Name     string
	Sslmode  string
	DBName   string
	Port     string
	Password string
}

func New(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=%s ", 
	cfg.Port, cfg.Name, cfg.Password, cfg.DBName, cfg.Sslmode ))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}