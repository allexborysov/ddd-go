package e2e

import (
	"context"
	"strings"
	"testing"
)

type BookClosedFlightParams struct {
	SeatNumber            string
	PassengerID           string
	ExpectMessageContains string
}

func BookClosedFlight(ctx context.Context, t *testing.T, env *Env, flight ScheduledFlight, params BookClosedFlightParams) {
	t.Helper()

	seat := params.SeatNumber
	if seat == "" {
		seat = "1A"
	}
	passenger := params.PassengerID
	if passenger == "" {
		passenger = "passenger-" + UniqueSuffix()
	}
	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "booking is closed"
	}

	req := map[string]any{
		"flight_id":    flight.FlightID,
		"seat_number":  seat,
		"passenger_id": passenger,
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/booking/flight", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("book_closed_flight: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected closed-booking attempt as expected: status=%d msg=%q flight=%s seat=%s",
		status, msg, flight.FlightID, seat)
}
