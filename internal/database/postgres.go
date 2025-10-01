package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/quantfidential/trading-ecosystem/exchange-data-adapter-go/internal/config"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	DB     *sql.DB
	config *config.Config
	logger *logrus.Logger
}

func NewPostgresDB(cfg *config.Config, logger *logrus.Logger) (*PostgresDB, error) {
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL is required")
	}

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnectionMaxIdleTime)

	return &PostgresDB{
		DB:     db,
		config: cfg,
		logger: logger,
	}, nil
}

func (p *PostgresDB) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := p.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.logger.Info("PostgreSQL connection established")
	return nil
}

func (p *PostgresDB) Disconnect(ctx context.Context) error {
	if p.DB != nil {
		if err := p.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		p.logger.Info("PostgreSQL connection closed")
	}
	return nil
}

func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}
