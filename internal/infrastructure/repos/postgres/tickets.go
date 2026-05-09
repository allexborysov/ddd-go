package postgres

import (
	"context"
	"errors"

	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent"
	ticketent "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent/ticket"
)

type ticketRepository struct {
	client *ent.Client
}

func NewTicketRepository(c *ent.Client) *ticketRepository {
	return &ticketRepository{client: c}
}

func (r *ticketRepository) Store(ctx context.Context, t *flight.Ticket) error {
	return r.client.Ticket.
		Create().
		SetID(string(t.ID)).
		SetFlightID(string(t.FlightID)).
		SetPassengerID(string(t.PassengerID)).
		SetSeat(string(t.Seat)).
		SetPrice(float64(t.Price)).
		OnConflictColumns("id").
		UpdateNewValues().
		Exec(ctx)
}

func (r *ticketRepository) Find(ctx context.Context, id string) (*flight.Ticket, error) {
	row, err := r.client.Ticket.
		Query().
		Where(ticketent.IDEQ(id)).
		WithFlight().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainTicket(row)
}

func toDomainTicket(row *ent.Ticket) (*flight.Ticket, error) {
	if row == nil {
		return nil, errors.New("nil ticket row")
	}
	if row.Edges.Flight == nil {
		return nil, errors.New("ticket flight edge not loaded")
	}
	return &flight.Ticket{
		ID:          flight.TicketID(row.ID),
		FlightID:    flight.FlightID(row.Edges.Flight.ID),
		PassengerID: flight.PassengerID(row.PassengerID),
		Seat:        inventory.SeatNumber(row.Seat),
		Price:       shared.Amount(row.Price),
	}, nil
}
