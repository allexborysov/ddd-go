package schedulingifc

import (
	"time"

	shared "github.com/allexborysov/aircraft"
	schedulingsvc "github.com/allexborysov/aircraft/internal/application/scheduling"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/allexborysov/aircraft/internal/infrastructure/validation"
)

type ScheduleFlightRequest struct {
	AircraftMSN        string             `json:"aircraft_msn" validate:"required"`
	OriginICAO         string             `json:"origin_icao" validate:"required"`
	DestinationICAO    string             `json:"destination_icao" validate:"required"`
	ScheduledDeparture time.Time          `json:"scheduled_departure" validate:"required"`
	ScheduledArrival   time.Time          `json:"scheduled_arrival" validate:"required"`
	CloseBookingBuffer time.Duration      `json:"close_booking_buffer" validate:"gte=0"`
	SeatPrices         map[string]float64 `json:"seat_prices" validate:"omitempty"`
}

func (r *ScheduleFlightRequest) Validate() error {
	return validation.Validate.Struct(r)
}

func (r *ScheduleFlightRequest) ToScheduleFlightCommand() *schedulingsvc.ScheduleFlightCommand {
	prices := make(map[inventory.SeatNumber]shared.Amount, len(r.SeatPrices))
	for k, v := range r.SeatPrices {
		prices[inventory.SeatNumber(k)] = shared.Amount(v)
	}
	return &schedulingsvc.ScheduleFlightCommand{
		AircraftMSN:        r.AircraftMSN,
		OriginICAO:         r.OriginICAO,
		DestinationICAO:    r.DestinationICAO,
		ScheduledDeparture: r.ScheduledDeparture,
		ScheduledArrival:   r.ScheduledArrival,
		CloseBookingBuffer: r.CloseBookingBuffer,
		SeatPrices:         prices,
	}
}

type ScheduleFlightResponse struct {
	Flight *flight.Flight `json:"flight"`
}

func ToScheduleFlightResponse(res *schedulingsvc.ScheduleFlightCommandResult) ScheduleFlightResponse {
	return ScheduleFlightResponse{Flight: res.Flight}
}
