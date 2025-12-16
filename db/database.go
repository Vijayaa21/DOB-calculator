package db

import (
	"context"
	"fmt"

	"github.com/dob-calculator/config"
	"github.com/dob-calculator/db/sqlc"
	"github.com/dob-calculator/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Database holds the database connection pool and queries
type Database struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Failed to parse database config", zap.Error(err))
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	queries := sqlc.New(&pgxPoolAdapter{pool: pool})

	return &Database{
		Pool:    pool,
		Queries: queries,
	}, nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		logger.Info("Database connection closed")
	}
}

// pgxPoolAdapter adapts pgxpool.Pool to sqlc.DBTX interface
type pgxPoolAdapter struct {
	pool *pgxpool.Pool
}

func (a *pgxPoolAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	result, err := a.pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (a *pgxPoolAdapter) Query(ctx context.Context, sql string, args ...interface{}) (sqlc.Rows, error) {
	rows, err := a.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRowsAdapter{rows: rows}, nil
}

func (a *pgxPoolAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) sqlc.Row {
	return a.pool.QueryRow(ctx, sql, args...)
}

// pgxRowsAdapter adapts pgx.Rows to sqlc.Rows interface
type pgxRowsAdapter struct {
	rows interface {
		Close()
		Err() error
		Next() bool
		Scan(dest ...interface{}) error
	}
}

func (r *pgxRowsAdapter) Close() {
	r.rows.Close()
}

func (r *pgxRowsAdapter) Err() error {
	return r.rows.Err()
}

func (r *pgxRowsAdapter) Next() bool {
	return r.rows.Next()
}

func (r *pgxRowsAdapter) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}
