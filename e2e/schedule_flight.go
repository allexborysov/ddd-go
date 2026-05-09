package e2e

import (
	"context"
	"testing"
	"time"
)

type ScheduleFlightParams struct {
	OriginICAO         string
	DestinationICAO    string
	ScheduledDeparture time.Time
	ScheduledArrival   time.Time
	CloseBookingBuffer time.Duration
	SeatPrices         map[string]float64
}

type ScheduledFlight struct {
	FlightID    string
	AircraftMSN string
}

func ScheduleFlight(ctx context.Context, t *testing.T, env *Env, aircraft CreatedAircraft, params ScheduleFlightParams) ScheduledFlight {
	t.Helper()

	req := map[string]any{
		"aircraft_msn":         aircraft.MSN,
		"origin_icao":          params.OriginICAO,
		"destination_icao":     params.DestinationICAO,
		"scheduled_departure":  params.ScheduledDeparture,
		"scheduled_arrival":    params.ScheduledArrival,
		"close_booking_buffer": params.CloseBookingBuffer,
		"seat_prices":          params.SeatPrices,
	}

	var resp struct {
		Flight struct {
			ID       string `json:"ID"`
			Aircraft string `json:"Aircraft"`
		} `json:"flight"`
	}

	postJSON(ctx, t, env, "/api/scheduling/flight", req, &resp)

	if resp.Flight.ID == "" {
		t.Fatalf("schedule_flight: server returned empty flight ID")
	}
	if resp.Flight.Aircraft != aircraft.MSN {
		t.Fatalf("schedule_flight: got Aircraft %q, want %q", resp.Flight.Aircraft, aircraft.MSN)
	}

	out := ScheduledFlight{
		FlightID:    resp.Flight.ID,
		AircraftMSN: resp.Flight.Aircraft,
	}
	t.Logf("scheduled flight: id=%s aircraft=%s %s->%s",
		out.FlightID, out.AircraftMSN, params.OriginICAO, params.DestinationICAO)
	return out
}
