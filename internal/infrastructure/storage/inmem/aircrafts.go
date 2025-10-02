package inmem

import (
	"context"
	"sync"

	"github.com/allexborysov/aircraft/internal/domain/inventory"
)

type aircraftRepository struct {
	mtx       sync.RWMutex
	aircrafts map[string]*inventory.Aircraft
}

func NewAircraftRepository() *aircraftRepository {
	return &aircraftRepository{
		aircrafts: make(map[string]*inventory.Aircraft),
	}
}

func (r *aircraftRepository) Store(ctx context.Context, a *inventory.Aircraft) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.aircrafts[string(a.MSN)] = a
	return nil
}

func (r *aircraftRepository) Find(ctx context.Context, msn string) (*inventory.Aircraft, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.aircrafts[string(msn)], nil
}
