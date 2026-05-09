package flight

import (
	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/pborman/uuid"
)

type TicketID string
type Ticket struct {
	ID          TicketID
	FlightID    FlightID
	PassengerID PassengerID

	Seat  inventory.SeatNumber
	Price shared.Amount
}

func NewTicket(flightID FlightID, passengerID PassengerID, seatNumber inventory.SeatNumber, price shared.Amount) Ticket {
	id := TicketID(uuid.New())
	return Ticket{
		ID:          id,
		FlightID:    flightID,
		PassengerID: passengerID,
		Seat:        seatNumber,
		Price:       price,
	}
}

type TicketPDFPresenter interface {
	GeneratePDF(t *Ticket) (string, error)
}
