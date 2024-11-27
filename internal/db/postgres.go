package db

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func LoadDB(connString string, log *logrus.Logger) (*pgxpool.Pool, error) {
	log.Debug("Attempting to connect to the database")
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Errorf("Failed to connect to the database: %v", err)
		return nil, err
	}

	log.Info("Successfully connected to the database")
	return pool, nil
}

func RunMigrations(connString string, log *logrus.Logger) error {
	log.Debug("Starting database migrations")
	m, err := migrate.New(
		fmt.Sprintf("file://%s", "./migrations"),
		connString)
	if err != nil {
		log.Errorf("Failed to initialize migrations: %v", err)
		return err
	}
	log.Debug("Migrations initialized successfully")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Errorf("Failed to apply migrations: %v", err)
		return fmt.Errorf("postgres.go: run migrations: up: %w", err)
	}
	log.Info("Database migrations applied successfully")
	return nil
}
