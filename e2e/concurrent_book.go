package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"
)

type ConcurrentBookParams struct {
	SeatNumber string
	Count      int // defaults to 10
}

// ConcurrentBook fires Count concurrent booking requests for the same seat and
// asserts that exactly one succeeds — the rest must be rejected because the seat
// is either held by the in-flight lock or already taken.
func ConcurrentBook(ctx context.Context, t *testing.T, env *Env, flight ScheduledFlight, params ConcurrentBookParams) {
	t.Helper()

	count := params.Count
	if count == 0 {
		count = 100
	}

	type result struct {
		status int
	}

	results := make([]result, count)
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		i := i
		go func() {
			defer wg.Done()

			body, _ := json.Marshal(map[string]any{
				"flight_id":    flight.FlightID,
				"seat_number":  params.SeatNumber,
				"passenger_id": "passenger-" + UniqueSuffix(),
			})

			req, err := http.NewRequestWithContext(ctx, http.MethodPost,
				env.BaseURL+"/api/booking/flight", bytes.NewReader(body))
			if err != nil {
				results[i] = result{status: -1}
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := env.Client.Do(req)
			if err != nil {
				results[i] = result{status: -1}
				return
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			results[i] = result{status: resp.StatusCode}
		}()
	}
	wg.Wait()

	successes := 0
	for _, r := range results {
		if r.status >= 200 && r.status < 300 {
			successes++
		}
	}

	if successes != 1 {
		t.Fatalf("concurrent_book: expected exactly 1 success out of %d requests for seat %s, got %d",
			count, params.SeatNumber, successes)
	}
	t.Logf("concurrent_book: exactly 1 of %d requests succeeded for seat=%s", count, params.SeatNumber)
}
