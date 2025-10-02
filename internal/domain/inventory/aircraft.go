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

func NewMSN(msn string) (MSN, error) {
	min := 5
	max := 20
	if len(msn) < min {
		return "", errors.New("MSN must be at least 5 characters")
	}
	if len(msn) > max {
		return "", errors.New("MSN must be at most 20 characters")
	}
	return MSN(msn), nil
}

func NewSeatNumber(seatNum string) (SeatNumber, error) {
	seatNum = strings.TrimSpace(seatNum)
	max := 3
	if seatNum == "" {
		return "", errors.New("seat number cannot be empty")
	}
	if len(seatNum) > max {
		return "", errors.New("seat number must be at most 3 characters")
	}
	return SeatNumber(seatNum), nil
}

func NewSeats(seatNumbers []string) (Seats, error) {
	if len(seatNumbers) < 1 {
		return nil, errors.New("at least one seat is required")
	}
	if len(seatNumbers) > 200 {
		return nil, errors.New("too many seats defined")
	}
	seen := make(map[string]struct{}, len(seatNumbers))
	result := make(Seats, 0, len(seatNumbers))
	for _, raw := range seatNumbers {
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
