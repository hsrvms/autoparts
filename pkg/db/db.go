package db

import (
	"context"
	"log"
	"time"

	"github.com/hsrvms/autoparts/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database represents a connection to the database
type Database struct {
	Pool *pgxpool.Pool
	cfg  *config.Config
}

// New creates a new database connection
func New(cfg *config.Config) (*Database, error) {
	// Create a connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDBConnString())
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create the connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return &Database{Pool: pool, cfg: cfg}, nil
}

// Close closes the database connection
func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("Database connection closed")
	}
}
