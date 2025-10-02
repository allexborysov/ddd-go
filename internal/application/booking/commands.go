package bookingsvc

import (
	"github.com/allexborysov/aircraft/internal/domain/flight"
)

type BookFlightCommand struct {
	FlightID            string
	SeatNumber          string
	PassengerPassportID string
}

type BookFlightCommandResult struct {
	Ticket    *flight.Ticket
	TicketPDF string
}
