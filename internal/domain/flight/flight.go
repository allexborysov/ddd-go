package flight

import (
	"context"
	"errors"
	"time"

	shared "github.com/allexborysov/aircraft"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/pborman/uuid"
)

type FlightID string

type Flight struct {
	ID FlightID

	Aircraft inventory.MSN
	Seats    map[inventory.SeatNumber]SeatAssignment

	Origin      ICAO
	Destination ICAO

	ScheduledDeparture time.Time
	ScheduledArrival   time.Time

	CloseBookingBuffer time.Duration
}

type PassengerPassportID string
type SeatAssignment struct {
	PassengerID PassengerPassportID
	Price       shared.Amount
}

type FlightRepository interface {
	Store(ctx context.Context, fl *Flight) error
	Find(ctx context.Context, id string) (*Flight, error)
}

var (
	ErrSeatBooked    = errors.New("Seat is already booked")
	ErrBookingClosed = errors.New("Booking is closed")
)

func (f *Flight) AssignSeat(passengerId PassengerPassportID, seatNumber inventory.SeatNumber) (*Ticket, error) {
	untilDeparture := time.Until(f.ScheduledDeparture)
	if untilDeparture < f.CloseBookingBuffer {
		return nil, ErrBookingClosed
	}

	existingSeat := f.Seats[seatNumber]
	if existingSeat.PassengerID != "" {
		return nil, ErrSeatBooked
	}

	f.Seats[seatNumber] = SeatAssignment{
		PassengerID: passengerId,
		Price:       existingSeat.Price,
	}

	ticket := NewTicket(f.ID, passengerId, seatNumber, existingSeat.Price)

	return &ticket, nil
}

var (
	ErrInvalidTimes      = errors.New("invalid schedule times")
	ErrNegativeBuffer    = errors.New("booking buffer must be positive duration")
	ErrOriginDestination = errors.New("origin and destination must differ")
)

func NewFlight(
	aircraft inventory.Aircraft,
	originICAO string,
	destinationICAO string,
	departure time.Time,
	arrival time.Time,
	closeBookingBuffer time.Duration,
	seatPrices map[inventory.SeatNumber]shared.Amount,
) (*Flight, error) {
	id := FlightID(uuid.New())

	if arrival.IsZero() || departure.IsZero() || !arrival.After(departure) {
		return nil, ErrInvalidTimes
	}
	if closeBookingBuffer < 0 {
		return nil, ErrNegativeBuffer
	}
	if originICAO == destinationICAO {
		return nil, ErrOriginDestination
	}

	orig, err := NewICAO(originICAO)
	if err != nil {
		return nil, err
	}
	dest, err := NewICAO(destinationICAO)
	if err != nil {
		return nil, err
	}

	seats := aircraft.Seats
	pricedSeats := make(map[inventory.SeatNumber]SeatAssignment, len(seats))
	for _, sn := range seats {
		pricedSeats[sn] = SeatAssignment{Price: seatPrices[sn]}
	}

	return &Flight{
		ID:                 id,
		Aircraft:           aircraft.MSN,
		Seats:              pricedSeats,
		Origin:             orig,
		Destination:        dest,
		ScheduledDeparture: departure,
		ScheduledArrival:   arrival,
		CloseBookingBuffer: closeBookingBuffer,
	}, nil
}
