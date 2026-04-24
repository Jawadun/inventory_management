package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")

type DB struct {
	*sql.DB
}

func ConnectDB(cfg *Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		if os.Getenv("DB_CREATE") == "true" {
			return createDBAndConnect(cfg)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func createDBAndConnect(cfg *Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	db.Close()

	return ConnectDB(cfg)
}
