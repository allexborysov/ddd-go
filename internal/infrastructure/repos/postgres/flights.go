package postgres

import (
	"context"
	"errors"
	"time"

	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	db "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type flightRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewFlightRepository(pool *pgxpool.Pool) *flightRepository {
	return &flightRepository{pool: pool, queries: db.New(pool)}
}

func (r *flightRepository) Store(ctx context.Context, f *flight.Flight) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	if err := qtx.UpsertFlight(ctx, db.UpsertFlightParams{
		ID:                   string(f.ID),
		AircraftID:           string(f.Aircraft),
		Origin:               string(f.Origin),
		Destination:          string(f.Destination),
		ScheduledDeparture:   toTimestamptz(f.ScheduledDeparture),
		ScheduledArrival:     toTimestamptz(f.ScheduledArrival),
		CloseBookingBufferNs: f.CloseBookingBuffer.Nanoseconds(),
	}); err != nil {
		return err
	}

	for sn, sa := range f.Seats {
		if err := qtx.UpsertFlightSeat(ctx, db.UpsertFlightSeatParams{
			FlightID:   string(f.ID),
			SeatNumber: string(sn),
			Price:      float64(sa.Price),
		}); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *flightRepository) StoreTicket(ctx context.Context, t *flight.Ticket) error {
	return r.queries.UpsertTicket(ctx, db.UpsertTicketParams{
		ID:          string(t.ID),
		FlightID:    string(t.FlightID),
		PassengerID: string(t.PassengerID),
		Seat:        string(t.Seat),
		Price:       float64(t.Price),
	})
}

func (r *flightRepository) Find(ctx context.Context, id string) (*flight.Flight, error) {
	row, err := r.queries.GetFlightByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	seatRows, err := r.queries.GetFlightSeatsByFlightID(ctx, id)
	if err != nil {
		return nil, err
	}

	seats := make(map[inventory.SeatNumber]flight.SeatAssignment, len(seatRows))
	for _, s := range seatRows {
		seats[inventory.SeatNumber(s.SeatNumber)] = flight.SeatAssignment{
			PassengerID: flight.PassengerID(fromNullableText(s.PassengerID)),
			Price:       shared.Amount(s.Price),
		}
	}

	return &flight.Flight{
		ID:                 flight.FlightID(row.ID),
		Aircraft:           inventory.MSN(row.AircraftID),
		Seats:              seats,
		Origin:             flight.ICAO(row.Origin),
		Destination:        flight.ICAO(row.Destination),
		ScheduledDeparture: fromTimestamptz(row.ScheduledDeparture),
		ScheduledArrival:   fromTimestamptz(row.ScheduledArrival),
		CloseBookingBuffer: time.Duration(row.CloseBookingBufferNs),
	}, nil
}
