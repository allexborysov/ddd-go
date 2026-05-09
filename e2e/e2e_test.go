package e2e

import (
	"context"
	"testing"
	"time"
)

func TestE2E(t *testing.T) {
	ctx := context.Background()
	env := NewEnv()

	seatPrices := map[string]float64{
		"1A": 250, "1B": 250, "2A": 200, "2B": 200,
	}
	departure := time.Now().Add(24 * time.Hour)
	arrival := departure.Add(6 * time.Hour)
	aircraft := CreateAircraft(ctx, t, env, CreateAircraftParams{
		MSN:   "MSN-" + UniqueSuffix(),
		Seats: []string{"1A", "1B", "2A", "2B"},
	})

	t.Run("open_flight", func(t *testing.T) {
		t.Parallel()
		flight := ScheduleFlight(ctx, t, env, aircraft, ScheduleFlightParams{
			OriginICAO:         "KJFK",
			DestinationICAO:    "KLAX",
			ScheduledDeparture: departure,
			ScheduledArrival:   arrival,
			CloseBookingBuffer: 1 * time.Hour,
			SeatPrices:         seatPrices,
		})

		t.Run("book_then_duplicate", func(t *testing.T) {
			t.Parallel()
			booked := BookFlight(ctx, t, env, flight, BookFlightParams{
				SeatNumber:  "1A",
				PassengerID: "passenger-" + UniqueSuffix(),
			})
			BookDuplicateSeat(ctx, t, env, booked, BookDuplicateSeatParams{})
		})

		t.Run("book_unknown_seat", func(t *testing.T) {
			t.Parallel()
			BookUnknownSeat(ctx, t, env, flight, BookUnknownSeatParams{})
		})

		t.Run("concurrent_book", func(t *testing.T) {
			t.Parallel()
			ConcurrentBook(ctx, t, env, flight, ConcurrentBookParams{SeatNumber: "2A"})
		})

		t.Run("book_available_seat_not_locked", func(t *testing.T) {
			t.Parallel()
			BookFlight(ctx, t, env, flight, BookFlightParams{
				SeatNumber:  "2B",
				PassengerID: "passenger-" + UniqueSuffix(),
			})
		})
	})

	t.Run("closed_flight", func(t *testing.T) {
		t.Parallel()
		// Booking window already closed: time-to-departure (30m) < buffer (2h).
		closedDeparture := time.Now().Add(30 * time.Minute)
		closedArrival := closedDeparture.Add(2 * time.Hour)
		closedBuffer := 2 * time.Hour
		flight := ScheduleFlight(ctx, t, env, aircraft, ScheduleFlightParams{
			OriginICAO:         "KJFK",
			DestinationICAO:    "KLAX",
			ScheduledDeparture: closedDeparture,
			ScheduledArrival:   closedArrival,
			CloseBookingBuffer: closedBuffer,
			SeatPrices:         seatPrices,
		})
		BookClosedFlight(ctx, t, env, flight, BookClosedFlightParams{SeatNumber: "2A"})
	})

	t.Run("invalid_timing_flight", func(t *testing.T) {
		t.Parallel()
		ScheduleInvalidTimingFlight(ctx, t, env, aircraft, ScheduleInvalidTimingFlightParams{
			Departure: departure,
			Arrival:   departure.Add(-1 * time.Hour),
		})
	})

	t.Run("same_origin_destination_flight", func(t *testing.T) {
		t.Parallel()
		ScheduleSameOriginDestinationFlight(ctx, t, env, aircraft, ScheduleSameOriginDestinationFlightParams{
			ICAO:      "KJFK",
			Departure: departure,
			Arrival:   arrival,
		})
	})

	t.Run("duplicate_aircraft", func(t *testing.T) {
		t.Parallel()
		RegisterDuplicateAircraft(ctx, t, env, aircraft, RegisterDuplicateAircraftParams{})
	})

	t.Run("unknown_aircraft_flight", func(t *testing.T) {
		t.Parallel()
		ScheduleUnknownAircraftFlight(ctx, t, env, ScheduleUnknownAircraftFlightParams{
			Departure: departure,
			Arrival:   arrival,
		})
	})
}
