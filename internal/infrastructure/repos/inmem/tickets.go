package inmem

import (
	"context"
	"sync"

	"github.com/allexborysov/aircraft/internal/domain/flight"
)

type ticketRepository struct {
	mtx     sync.RWMutex
	tickets map[string]*flight.Ticket
}

func NewTicketRepository() *ticketRepository {
	return &ticketRepository{
		tickets: make(map[string]*flight.Ticket),
	}
}

func (r *ticketRepository) Store(ctx context.Context, t *flight.Ticket) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.tickets[string(t.ID)] = t
	return nil
}

func (r *ticketRepository) Find(ctx context.Context, id string) (*flight.Ticket, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.tickets[id], nil
}
