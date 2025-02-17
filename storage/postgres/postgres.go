package postgres

import (
	"context"
	"fmt"
	"weather/config"
	"weather/pkg/logger"
	"weather/storage"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool   *pgxpool.Pool
	logger logger.ILogger
}

func New(ctx context.Context, cfg config.Config, logger logger.ILogger) (storage.IStorage, error) {
	url := fmt.Sprintf(`host=%s port=%v user=%s password=%s database=%s sslmode=disable`,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase)

	pgPoolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgPoolConfig.MaxConns = 100
	pgPoolConfig.MaxConnLifetime = time.Hour

	newPool, err := pgxpool.NewWithConfig(context.Background(), pgPoolConfig)

	if err != nil {
		fmt.Println("error while connecting to db ", err.Error())
		return nil, err
	}
	return Store{
		Pool: newPool,
	}, nil
}
func (s Store) CloseDB() {
	s.Pool.Close()
}

func (s Store) Weather() storage.IWeatherStorage {
	Newcontact := NewWeatherRepo(s.Pool, s.logger)

	return &Newcontact
}
