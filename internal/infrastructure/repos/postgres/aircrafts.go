package postgres

import (
	"context"
	"errors"

	"github.com/allexborysov/aircraft/internal/domain/inventory"
	db "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type aircraftRepository struct {
	queries *db.Queries
}

func NewAircraftRepository(pool *pgxpool.Pool) *aircraftRepository {
	return &aircraftRepository{queries: db.New(pool)}
}

func (r *aircraftRepository) Store(ctx context.Context, a *inventory.Aircraft) error {
	seats := make([]string, len(a.Seats))
	for i, s := range a.Seats {
		seats[i] = string(s)
	}
	return r.queries.UpsertAircraft(ctx, db.UpsertAircraftParams{
		ID:    string(a.MSN),
		Seats: seats,
	})
}

func (r *aircraftRepository) Find(ctx context.Context, msn string) (*inventory.Aircraft, error) {
	row, err := r.queries.GetAircraftByID(ctx, msn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return inventory.NewAircraft(row.ID, row.Seats)
}
