package postgres

import (
	"context"
	"errors"

	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	db "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ticketRepository struct {
	queries *db.Queries
}

func NewTicketRepository(pool *pgxpool.Pool) *ticketRepository {
	return &ticketRepository{queries: db.New(pool)}
}

func (r *ticketRepository) Store(ctx context.Context, t *flight.Ticket) error {
	return r.queries.UpsertTicket(ctx, db.UpsertTicketParams{
		ID:          string(t.ID),
		FlightID:    string(t.FlightID),
		PassengerID: string(t.PassengerID),
		Seat:        string(t.Seat),
		Price:       float64(t.Price),
	})
}

func (r *ticketRepository) Find(ctx context.Context, id string) (*flight.Ticket, error) {
	row, err := r.queries.GetTicketByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &flight.Ticket{
		ID:          flight.TicketID(row.ID),
		FlightID:    flight.FlightID(row.FlightID),
		PassengerID: flight.PassengerID(row.PassengerID),
		Seat:        inventory.SeatNumber(row.Seat),
		Price:       shared.Amount(row.Price),
	}, nil
}
