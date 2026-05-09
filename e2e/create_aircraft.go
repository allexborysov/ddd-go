package e2e

import (
	"context"
	"testing"
)

type CreateAircraftParams struct {
	MSN   string
	Seats []string
}

type CreatedAircraft struct {
	MSN   string
	Seats []string
}

func CreateAircraft(ctx context.Context, t *testing.T, env *Env, params CreateAircraftParams) CreatedAircraft {
	t.Helper()

	req := map[string]any{
		"msn":   params.MSN,
		"seats": params.Seats,
	}

	var resp struct {
		Aircraft struct {
			MSN   string   `json:"MSN"`
			Seats []string `json:"Seats"`
		} `json:"aircraft"`
	}

	postJSON(ctx, t, env, "/api/inventory/aircraft", req, &resp)

	if resp.Aircraft.MSN == "" {
		t.Fatalf("create_aircraft: server returned empty MSN")
	}

	out := CreatedAircraft{
		MSN:   resp.Aircraft.MSN,
		Seats: resp.Aircraft.Seats,
	}
	t.Logf("created aircraft: MSN=%s seats=%v", out.MSN, out.Seats)
	return out
}
