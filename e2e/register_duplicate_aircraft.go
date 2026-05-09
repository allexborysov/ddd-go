package e2e

import (
	"context"
	"strings"
	"testing"
)

type RegisterDuplicateAircraftParams struct {
	Seats                 []string
	ExpectMessageContains string
}

func RegisterDuplicateAircraft(ctx context.Context, t *testing.T, env *Env, aircraft CreatedAircraft, params RegisterDuplicateAircraftParams) {
	t.Helper()

	seats := params.Seats
	if len(seats) == 0 {
		seats = []string{"9A"}
	}
	expect := params.ExpectMessageContains
	if expect == "" {
		expect = "already registered"
	}

	req := map[string]any{
		"msn":   aircraft.MSN,
		"seats": seats,
	}

	status, msg := postJSONExpectError(ctx, t, env, "/api/inventory/aircraft", req)

	if !strings.Contains(strings.ToLower(msg), strings.ToLower(expect)) {
		t.Fatalf("register_duplicate_aircraft: expected message to contain %q, got %q (status %d)",
			expect, msg, status)
	}
	t.Logf("rejected duplicate aircraft as expected: status=%d msg=%q msn=%s",
		status, msg, aircraft.MSN)
}
