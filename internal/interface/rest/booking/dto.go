package bookingifc

import (
	bookingsvc "github.com/allexborysov/aircraft/internal/application/booking"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/infrastructure/validation"
)

type BookFlightRequest struct {
	FlightID            string `json:"flight_id" validate:"required"`
	SeatNumber          string `json:"seat_number" validate:"required"`
	PassengerPassportID string `json:"passenger_passport_id" validate:"required"`
}

func (req *BookFlightRequest) Validate() error {
	return validation.Validate.Struct(req)
}

func (req *BookFlightRequest) ToBookFlightCommand() *bookingsvc.BookFlightCommand {
	return &bookingsvc.BookFlightCommand{
		FlightID:            req.FlightID,
		SeatNumber:          req.SeatNumber,
		PassengerPassportID: req.PassengerPassportID,
	}
}

type BookFlightResponse struct {
	Ticket    *flight.Ticket `json:"ticket"`
	TicketPDF string         `json:"ticket_pdf"`
}

func ToBookFlightResposnse(result *bookingsvc.BookFlightCommandResult) BookFlightResponse {
	return BookFlightResponse{
		Ticket:    result.Ticket,
		TicketPDF: result.TicketPDF,
	}
}
