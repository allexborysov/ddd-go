package schedulingsvc

import (
	"context"
	"errors"

	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
)

type Service struct {
	flights   flight.FlightRepository
	aircrafts inventory.AircraftRepository
}

func New(flights flight.FlightRepository, aircrafts inventory.AircraftRepository) *Service {
	return &Service{flights: flights, aircrafts: aircrafts}
}

func (s *Service) ScheduleFlight(ctx context.Context, command *ScheduleFlightCommand) (*ScheduleFlightCommandResult, error) {
	aircraft, err := s.aircrafts.Find(ctx, command.AircraftMSN)
	if err != nil {
		return nil, err
	}
	if aircraft == nil {
		return nil, errors.New("Aircraft not found")
	}

	fl, err := flight.NewFlight(
		*aircraft,
		command.OriginICAO,
		command.DestinationICAO,
		command.ScheduledDeparture,
		command.ScheduledArrival,
		command.CloseBookingBuffer,
		command.SeatPrices,
	)
	if err != nil {
		return nil, err
	}

	err = s.flights.Store(ctx, fl)
	if err != nil {
		return nil, err
	}

	return &ScheduleFlightCommandResult{Flight: fl}, nil
}
