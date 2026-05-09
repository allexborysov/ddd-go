package e2e

import (
	"context"
	"strings"
	"testing"
)

type BookUnknownSeatParams struct {
	SeatNumber            string
	PassengerID           string
	ExpectMessageContains string
}

func BookUnknownSeat(ctx context.Context, t *testing.T, env *Env, flight ScheduledFlight, params BookUnknownSeatParams) {
	t.Helper()

	seat := params.SeatNumber
	if seat == "" {
		seat = "99Z"
	}
	passenger := params.PassengerID
	if passenger == "" {
		passenger = "passenger-" + UniqueSuffix()
	}
	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "seat does not exist"
	}

	req := map[string]any{
		"flight_id":    flight.FlightID,
		"seat_number":  seat,
		"passenger_id": passenger,
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/booking/flight", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("book_unknown_seat: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected unknown-seat booking as expected: status=%d msg=%q flight=%s seat=%s",
		status, msg, flight.FlightID, seat)
}
