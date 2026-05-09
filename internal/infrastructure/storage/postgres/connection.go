package postgres

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent"
)

func MustConnectPostgres(connString string) *ent.Client {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("Failed to open Postgres: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping Postgres: %v", err)
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed to apply Postgres schema: %v", err)
	}

	return client
}
