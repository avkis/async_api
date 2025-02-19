package store

import (
	"async_api/config"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewPostgresDB(conf *config.Config) (*sql.DB, error) {
	dsn := conf.DataSourceName()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database connection: %w", err)
	}

	return db, nil
}
