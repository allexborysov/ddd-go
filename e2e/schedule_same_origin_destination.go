package e2e

import (
	"context"
	"strings"
	"testing"
	"time"
)

type ScheduleSameOriginDestinationFlightParams struct {
	ICAO                  string
	Departure             time.Time
	Arrival               time.Time
	ExpectMessageContains string
}

func ScheduleSameOriginDestinationFlight(ctx context.Context, t *testing.T, env *Env, aircraft CreatedAircraft, params ScheduleSameOriginDestinationFlightParams) {
	t.Helper()

	icao := params.ICAO
	if icao == "" {
		icao = "KJFK"
	}
	departure := params.Departure
	if departure.IsZero() {
		departure = time.Now().Add(24 * time.Hour)
	}
	arrival := params.Arrival
	if arrival.IsZero() {
		arrival = departure.Add(3 * time.Hour)
	}
	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "origin and destination must differ"
	}

	req := map[string]any{
		"aircraft_msn":         aircraft.MSN,
		"origin_icao":          icao,
		"destination_icao":     icao,
		"scheduled_departure":  departure,
		"scheduled_arrival":    arrival,
		"close_booking_buffer": time.Duration(0),
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/scheduling/flight", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("schedule_same_origin_destination_flight: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected same-origin-destination flight as expected: status=%d msg=%q icao=%s",
		status, msg, icao)
}
