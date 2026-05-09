package postgres

import (
	"context"
	"errors"

	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent"
)

type aircraftRepository struct {
	client *ent.Client
}

func NewAircraftRepository(c *ent.Client) *aircraftRepository {
	return &aircraftRepository{client: c}
}

func (r *aircraftRepository) Store(ctx context.Context, a *inventory.Aircraft) error {
	seats := make([]string, len(a.Seats))
	for i, s := range a.Seats {
		seats[i] = string(s)
	}

	return r.client.Aircraft.
		Create().
		SetID(string(a.MSN)).
		SetSeats(seats).
		OnConflictColumns("id").
		UpdateNewValues().
		Exec(ctx)
}

func (r *aircraftRepository) Find(ctx context.Context, msn string) (*inventory.Aircraft, error) {
	row, err := r.client.Aircraft.Get(ctx, msn)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainAircraft(row)
}

func toDomainAircraft(row *ent.Aircraft) (*inventory.Aircraft, error) {
	if row == nil {
		return nil, errors.New("nil aircraft row")
	}
	return inventory.NewAircraft(row.ID, row.Seats)
}
