package e2e

import (
	"context"
	"strings"
	"testing"
	"time"
)

type ScheduleInvalidTimingFlightParams struct {
	OriginICAO            string
	DestinationICAO       string
	Departure             time.Time
	Arrival               time.Time
	ExpectMessageContains string
}

func ScheduleInvalidTimingFlight(ctx context.Context, t *testing.T, env *Env, aircraft CreatedAircraft, params ScheduleInvalidTimingFlightParams) {
	t.Helper()

	origin := params.OriginICAO
	if origin == "" {
		origin = "KJFK"
	}
	dest := params.DestinationICAO
	if dest == "" {
		dest = "KLAX"
	}

	departure := params.Departure
	if departure.IsZero() {
		departure = time.Now().Add(24 * time.Hour)
	}
	arrival := params.Arrival
	if arrival.IsZero() {
		arrival = departure
	}

	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "invalid schedule times"
	}

	req := map[string]any{
		"aircraft_msn":         aircraft.MSN,
		"origin_icao":          origin,
		"destination_icao":     dest,
		"scheduled_departure":  departure,
		"scheduled_arrival":    arrival,
		"close_booking_buffer": time.Duration(0),
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/scheduling/flight", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("schedule_invalid_timing_flight: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected invalid-timing flight as expected: status=%d msg=%q (departure=%s arrival=%s)",
		status, msg, departure.Format(time.RFC3339), arrival.Format(time.RFC3339))
}
