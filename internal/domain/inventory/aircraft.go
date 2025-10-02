package inventory

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Aircraft struct {
	MSN   MSN
	Seats Seats
}

type AircraftRepository interface {
	Store(ctx context.Context, a *Aircraft) error
	Find(ctx context.Context, msn string) (*Aircraft, error)
}

func NewAircraft(msn string, seats []string) (*Aircraft, error) {
	m, err := NewMSN(msn)
	if err != nil {
		return nil, err
	}
	ss, err := NewSeats(seats)
	if err != nil {
		return nil, err
	}

	return &Aircraft{
		MSN:   m,
		Seats: ss,
	}, nil
}

// Manufacturer Serial Number assigned at factory. Uniquely identifies a particular Aircraft.
type MSN string
type SeatNumber string
type Seats []SeatNumber

func NewMSN(s string) (MSN, error) {
	min := 5
	max := 20
	if len(s) < min {
		return "", errors.New("MSN must be at least 5 characters")
	}
	if len(s) > max {
		return "", errors.New("MSN must be at most 20 characters")
	}
	return MSN(s), nil
}

func NewSeatNumber(s string) (SeatNumber, error) {
	s = strings.TrimSpace(s)
	max := 3
	if s == "" {
		return "", errors.New("seat number cannot be empty")
	}
	if len(s) > max {
		return "", errors.New("seat number must be at most 3 characters")
	}
	return SeatNumber(s), nil
}

func NewSeats(input []string) (Seats, error) {
	if len(input) < 1 {
		return nil, errors.New("at least one seat is required")
	}
	if len(input) > 200 {
		return nil, errors.New("too many seats defined")
	}
	seen := make(map[string]struct{}, len(input))
	result := make(Seats, 0, len(input))
	for _, raw := range input {
		if _, ok := seen[raw]; ok {
			return nil, fmt.Errorf("duplicate seat: %s", raw)
		}
		seen[raw] = struct{}{}
		seatNumber, err := NewSeatNumber(raw)
		if err != nil {
			return nil, err
		}
		result = append(result, seatNumber)
	}
	return result, nil
}
