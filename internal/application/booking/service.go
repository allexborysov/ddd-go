package bookingsvc

import (
	"context"
	"errors"

	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
)

type BookingSync interface {
	Lock(ctx context.Context, seatNumber string) bool
	Unlock(ctx context.Context, seatNumber string)
}

type Service struct {
	mutex              BookingSync
	flights            flight.FlightRepository
	tickets            flight.TicketRepository
	ticketPDFPresenter flight.TicketPDFPresenter
}

func New(
	mutex BookingSync,
	flights flight.FlightRepository,
	tickets flight.TicketRepository,
	ticketPDFPresenter flight.TicketPDFPresenter) *Service {
	return &Service{
		mutex:              mutex,
		flights:            flights,
		tickets:            tickets,
		ticketPDFPresenter: ticketPDFPresenter,
	}
}

func (s *Service) BookFlight(ctx context.Context, command *BookFlightCommand) (*BookFlightCommandResult, error) {
	locked := s.mutex.Lock(ctx, command.SeatNumber)
	if !locked {
		return nil, errors.New("Seat is being held by another customer.")
	}
	defer s.mutex.Unlock(ctx, command.SeatNumber)

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

	ticket, err := fl.AssignSeat(
		flight.PassengerPassportID(command.PassengerPassportID),
		seatNumber,
	)
	if err != nil {
		return nil, err
	}

	pdf, err := s.ticketPDFPresenter.GeneratePDF(ticket)
	if err != nil {
		return nil, err
	}

	err = s.tickets.Store(ctx, ticket)
	if err != nil {
		return nil, err
	}
	err = s.flights.Store(ctx, fl)
	if err != nil {
		return nil, err
	}

	return &BookFlightCommandResult{
		Ticket:    ticket,
		TicketPDF: pdf,
	}, nil
}
