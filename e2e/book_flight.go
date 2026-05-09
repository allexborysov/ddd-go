package e2e

import (
	"context"
	"testing"
)

type BookFlightParams struct {
	SeatNumber  string
	PassengerID string
}

type BookedFlight struct {
	TicketID    string
	FlightID    string
	PassengerID string
	SeatNumber  string
}

func BookFlight(ctx context.Context, t *testing.T, env *Env, flight ScheduledFlight, params BookFlightParams) BookedFlight {
	t.Helper()

	req := map[string]any{
		"flight_id":    flight.FlightID,
		"seat_number":  params.SeatNumber,
		"passenger_id": params.PassengerID,
	}

	var resp struct {
		Ticket struct {
			ID          string `json:"ID"`
			FlightID    string `json:"FlightID"`
			PassengerID string `json:"PassengerID"`
			Seat        string `json:"Seat"`
		} `json:"ticket"`
		TicketPDF string `json:"ticket_pdf"`
	}

	postJSON(ctx, t, env, "/api/booking/flight", req, &resp)

	if resp.Ticket.ID == "" {
		t.Fatalf("book_flight: server returned empty ticket ID")
	}
	if resp.Ticket.FlightID != flight.FlightID {
		t.Fatalf("book_flight: got FlightID %q, want %q", resp.Ticket.FlightID, flight.FlightID)
	}
	if resp.Ticket.Seat != params.SeatNumber {
		t.Fatalf("book_flight: got Seat %q, want %q", resp.Ticket.Seat, params.SeatNumber)
	}
	if resp.Ticket.PassengerID != params.PassengerID {
		t.Fatalf("book_flight: got PassengerID %q, want %q", resp.Ticket.PassengerID, params.PassengerID)
	}

	out := BookedFlight{
		TicketID:    resp.Ticket.ID,
		FlightID:    resp.Ticket.FlightID,
		PassengerID: resp.Ticket.PassengerID,
		SeatNumber:  resp.Ticket.Seat,
	}
	t.Logf("booked flight: ticket=%s flight=%s passenger=%s seat=%s",
		out.TicketID, out.FlightID, out.PassengerID, out.SeatNumber)
	return out
}
