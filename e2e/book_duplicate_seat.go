package e2e

import (
	"context"
	"strings"
	"testing"
)

type BookDuplicateSeatParams struct {
	PassengerID           string
	ExpectMessageContains string
}

func BookDuplicateSeat(ctx context.Context, t *testing.T, env *Env, booked BookedFlight, params BookDuplicateSeatParams) {
	t.Helper()

	passenger := params.PassengerID
	if passenger == "" {
		passenger = "passenger-" + UniqueSuffix()
	}

	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "seat is already booked"
	}

	req := map[string]any{
		"flight_id":    booked.FlightID,
		"seat_number":  booked.SeatNumber,
		"passenger_id": passenger,
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/booking/flight", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("book_duplicate_seat: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected duplicate booking as expected: status=%d msg=%q flight=%s seat=%s",
		status, msg, booked.FlightID, booked.SeatNumber)
}
