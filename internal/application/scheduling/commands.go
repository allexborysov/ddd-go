package schedulingsvc

import (
	"time"

	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
)

type ScheduleFlightCommand struct {
	AircraftMSN        string
	OriginICAO         string
	DestinationICAO    string
	ScheduledDeparture time.Time
	ScheduledArrival   time.Time
	CloseBookingBuffer time.Duration
	SeatPrices         map[inventory.SeatNumber]shared.Amount
}

type ScheduleFlightCommandResult struct {
	Flight *flight.Flight
}
