package platform

import (
	"context"
	"fmt"
	"time"

	"pos-system/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func ConnectDatabase(ctx context.Context, cfg config.Config, logger zerolog.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	var pool *pgxpool.Pool
	for attempt := 1; attempt <= cfg.DBConnectRetries; attempt++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = pool.Ping(pingCtx)
			cancel()
			if err == nil {
				return pool, nil
			}
			pool.Close()
		}

		logger.Warn().Err(err).Int("attempt", attempt).Msg("database connection attempt failed")
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("connect database after retries: %w", err)
}
