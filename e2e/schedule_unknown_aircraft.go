package e2e

import (
	"context"
	"strings"
	"testing"
	"time"
)

type ScheduleUnknownAircraftFlightParams struct {
	AircraftMSN           string
	OriginICAO            string
	DestinationICAO       string
	Departure             time.Time
	Arrival               time.Time
	ExpectMessageContains string
}

func ScheduleUnknownAircraftFlight(ctx context.Context, t *testing.T, env *Env, params ScheduleUnknownAircraftFlightParams) {
	t.Helper()

	msn := params.AircraftMSN
	if msn == "" {
		msn = "GHOST-" + UniqueSuffix()
	}
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
		arrival = departure.Add(3 * time.Hour)
	}
	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "aircraft not found"
	}

	req := map[string]any{
		"aircraft_msn":         msn,
		"origin_icao":          origin,
		"destination_icao":     dest,
		"scheduled_departure":  departure,
		"scheduled_arrival":    arrival,
		"close_booking_buffer": time.Duration(0),
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/scheduling/flight", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("schedule_unknown_aircraft_flight: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected unknown-aircraft flight as expected: status=%d msg=%q msn=%s",
		status, msg, msn)
}
