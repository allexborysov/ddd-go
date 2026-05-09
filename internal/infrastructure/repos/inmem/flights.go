package inmem

import (
	"context"
	"errors"
	"sync"

	"github.com/allexborysov/aircraft/internal/domain/flight"
)

type flightRepository struct {
	mtx     sync.RWMutex
	flights map[string]*flight.Flight
	tickets map[string]*flight.Ticket
}

func NewFlightRepository() *flightRepository {
	return &flightRepository{
		flights: make(map[string]*flight.Flight),
		tickets: make(map[string]*flight.Ticket),
	}
}

func (r *flightRepository) Store(ctx context.Context, f *flight.Flight) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.flights[string(f.ID)] = f
	return nil
}

func (r *flightRepository) Find(ctx context.Context, id string) (*flight.Flight, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.flights[id], nil
}

func (r *flightRepository) StoreTicket(ctx context.Context, t *flight.Ticket) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	fl, ok := r.flights[string(t.FlightID)]
	if !ok {
		return errors.New("flight not found")
	}
	seat := fl.Seats[t.Seat]
	seat.PassengerID = t.PassengerID
	fl.Seats[t.Seat] = seat
	r.storeTicket(t)
	return nil
}

func (r *flightRepository) storeTicket(t *flight.Ticket) {
	r.tickets[string(t.ID)] = t
}
