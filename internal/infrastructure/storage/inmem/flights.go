package inmem

import (
	"context"
	"sync"

	"github.com/allexborysov/aircraft/internal/domain/flight"
)

type flightRepository struct {
	mtx     sync.RWMutex
	flights map[string]*flight.Flight
}

func NewFlightRepository() *flightRepository {
	return &flightRepository{
		flights: make(map[string]*flight.Flight),
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
