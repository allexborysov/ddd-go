package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"github.com/allexborysov/aircraft/config"
	pg "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres"
)

func setupTestContainers(cfg *config.Config) {
	ctx := context.Background()
	slog.Info("Setup test containers")

	pgCtr, err := tcpostgres.Run(ctx, "postgres:18-alpine",
		tcpostgres.WithDatabase("aircraft"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
	)
	if err != nil {
		log.Fatalf("postgres container: %v", err)
	}

	redisCtr, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		log.Fatalf("redis container: %v", err)
	}

	pgHost, _ := pgCtr.Host(ctx)
	pgPort, _ := pgCtr.MappedPort(ctx, "5432")
	redisAddr, _ := redisCtr.Endpoint(ctx, "")

	cfg.Postgres.Host = pgHost
	cfg.Postgres.Port = int(pgPort.Num())
	cfg.Postgres.User = "postgres"
	cfg.Postgres.Password = "postgres"
	cfg.Postgres.DBName = "aircraft"
	cfg.Redis.Addr = redisAddr

	pgConnStr, _ := pgCtr.ConnectionString(ctx, "sslmode=disable")
	runMigrations(pgConnStr)
}

func runMigrations(connStr string) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("open migrate db: %v", err)
	}
	defer db.Close()

	deadline := time.Now().Add(30 * time.Second)
	for {
		if err = db.Ping(); err == nil {
			break
		}
		if time.Now().After(deadline) {
			log.Fatalf("postgres ping timeout: %v", err)
		}
		time.Sleep(250 * time.Millisecond)
	}

	driver, err := pgmigrate.WithInstance(db, &pgmigrate.Config{})
	if err != nil {
		log.Fatalf("migrate driver: %v", err)
	}

	src, err := iofs.New(pg.Migrations, "migrations")
	if err != nil {
		log.Fatalf("migrate source: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		log.Fatalf("migrate instance: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
}
