package db

import (
	"context"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Config() *pgxpool.Config {
	dbUrl := os.Getenv("DB_CONNECTION_STRING")

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = int32(max(runtime.NumCPU(), 4))
	dbConfig.MinConns = 0
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = time.Minute * 30
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	return dbConfig
}

func NewConnectionPool() (pool *pgxpool.Pool) {
	pool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal(err)
		return
	}

	return pool
}
