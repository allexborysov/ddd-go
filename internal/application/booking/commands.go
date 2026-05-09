package bookingsvc

import (
	"github.com/allexborysov/aircraft/internal/domain/flight"
)

type BookFlightCommand struct {
	FlightID    string
	SeatNumber  string
	PassengerID string
}

type BookFlightCommandResult struct {
	Ticket    *flight.Ticket
	TicketPDF string
}
