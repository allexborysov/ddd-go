package inventorysvc

import (
	"context"
	"errors"

	"github.com/allexborysov/aircraft/internal/domain/inventory"
)

type Service struct {
	aircrafts inventory.AircraftRepository
}

func New(aircrafts inventory.AircraftRepository) *Service {
	return &Service{
		aircrafts,
	}
}

func (s *Service) RegisterAircraft(ctx context.Context, command *RegisterAircraftCommand) (*RegisterAircraftCommandResult, error) {
	found, err := s.aircrafts.Find(ctx, command.MSN)
	if err != nil {
		return nil, err
	}
	if found != nil {
		return nil, errors.New("Aircraft with same MSN is already registered")
	}

	aircraft, err := inventory.NewAircraft(command.MSN, command.Seats)
	if err != nil {
		return nil, err
	}

	err = s.aircrafts.Store(ctx, aircraft)
	if err != nil {
		return nil, err
	}

	return &RegisterAircraftCommandResult{aircraft}, nil
}
