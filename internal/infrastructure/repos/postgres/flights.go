package postgres

import (
	"context"
	"errors"
	"time"

	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent"
	flightent "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent/flight"
	flightseatent "github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres/ent/flightseat"
)

type flightRepository struct {
	client *ent.Client
}

func NewFlightRepository(c *ent.Client) *flightRepository {
	return &flightRepository{client: c}
}

func (r *flightRepository) Store(ctx context.Context, f *flight.Flight) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.Flight.
		Create().
		SetID(string(f.ID)).
		SetAircraftID(string(f.Aircraft)).
		SetOrigin(string(f.Origin)).
		SetDestination(string(f.Destination)).
		SetScheduledDeparture(f.ScheduledDeparture).
		SetScheduledArrival(f.ScheduledArrival).
		SetCloseBookingBufferNs(f.CloseBookingBuffer.Nanoseconds()).
		OnConflictColumns("id").
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return err
	}

	for sn, sa := range f.Seats {
		b := tx.FlightSeat.
			Create().
			SetFlightID(string(f.ID)).
			SetSeatNumber(string(sn)).
			SetPrice(float64(sa.Price))
		if sa.PassengerID != "" {
			b.SetPassengerID(string(sa.PassengerID))
		}
		if err := b.
			OnConflictColumns("flight_id", "seat_number").
			UpdateNewValues().
			Exec(ctx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *flightRepository) Find(ctx context.Context, id string) (*flight.Flight, error) {
	row, err := r.client.Flight.
		Query().
		Where(flightent.IDEQ(id)).
		WithAircraft().
		WithSeats(func(q *ent.FlightSeatQuery) {
			q.Order(ent.Asc(flightseatent.FieldSeatNumber))
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainFlight(row)
}

func toDomainFlight(row *ent.Flight) (*flight.Flight, error) {
	if row == nil {
		return nil, errors.New("nil flight row")
	}
	if row.Edges.Aircraft == nil {
		return nil, errors.New("flight aircraft edge not loaded")
	}

	seats := make(map[inventory.SeatNumber]flight.SeatAssignment, len(row.Edges.Seats))
	for _, s := range row.Edges.Seats {
		seats[inventory.SeatNumber(s.SeatNumber)] = flight.SeatAssignment{
			PassengerID: flight.PassengerID(s.PassengerID),
			Price:       shared.Amount(s.Price),
		}
	}

	return &flight.Flight{
		ID:                 flight.FlightID(row.ID),
		Aircraft:           inventory.MSN(row.Edges.Aircraft.ID),
		Seats:              seats,
		Origin:             flight.ICAO(row.Origin),
		Destination:        flight.ICAO(row.Destination),
		ScheduledDeparture: row.ScheduledDeparture,
		ScheduledArrival:   row.ScheduledArrival,
		CloseBookingBuffer: time.Duration(row.CloseBookingBufferNs),
	}, nil
}
