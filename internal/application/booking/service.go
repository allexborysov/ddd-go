package bookingsvc

import (
	"context"
	"errors"

	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
)

type BookingSync interface {
	Lock(ctx context.Context, key string) bool
	Unlock(ctx context.Context, key string)
}

type Service struct {
	mutex              BookingSync
	flights            flight.FlightRepository
	ticketPDFPresenter flight.TicketPDFPresenter
}

func New(
	mutex BookingSync,
	flights flight.FlightRepository,
	ticketPDFPresenter flight.TicketPDFPresenter) *Service {
	return &Service{
		mutex:              mutex,
		flights:            flights,
		ticketPDFPresenter: ticketPDFPresenter,
	}
}

func (s *Service) BookFlight(ctx context.Context, command *BookFlightCommand) (*BookFlightCommandResult, error) {
	lockKey := command.FlightID + ":" + command.SeatNumber
	acquired := s.mutex.Lock(ctx, lockKey)
	if !acquired {
		return nil, errors.New("Seat is being held by another customer.")
	}
	defer s.mutex.Unlock(ctx, lockKey)

	fl, err := s.flights.Find(ctx, command.FlightID)
	if err != nil {
		return nil, err
	}
	if fl == nil {
		return nil, errors.New("Flight not found")
	}

	seatNumber, err := inventory.NewSeatNumber(command.SeatNumber)
	if err != nil {
		return nil, err
	}

	ticket, err := fl.IssueTicket(
		flight.PassengerID(command.PassengerID),
		seatNumber,
	)
	if err != nil {
		return nil, err
	}

	pdf, err := s.ticketPDFPresenter.GeneratePDF(ticket)
	if err != nil {
		return nil, err
	}

	err = s.flights.StoreTicket(ctx, ticket)
	if err != nil {
		return nil, err
	}

	return &BookFlightCommandResult{
		Ticket:    ticket,
		TicketPDF: pdf,
	}, nil
}
