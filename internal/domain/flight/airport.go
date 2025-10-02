package flight

import (
	"errors"
	"strings"

	shared "github.com/allexborysov/aircraft"
)

type Airport struct {
	ICAO     ICAO
	TimeZone shared.TimeZone
	Location shared.Location
}

func NewAirport(icao, tz string, loc shared.Location) (*Airport, error) {
	i, err := NewICAO(icao)
	if err != nil {
		return nil, err
	}
	if tz == "" {
		return nil, errors.New("Time Zone cannot be empty")
	}
	if loc.Latitude < -90 || loc.Latitude > 90 {
		return nil, errors.New("Invalid latitude")
	}
	if loc.Longitude < -180 || loc.Longitude > 180 {
		return nil, errors.New("Invalid longitude")
	}

	return &Airport{
		ICAO:     i,
		TimeZone: shared.TimeZone(tz),
		Location: loc,
	}, nil
}

// International Civil Aviation Organization uniquely identifies a particular Airport.
type ICAO string

func NewICAO(s string) (ICAO, error) {
	s = strings.TrimSpace(s)
	max := 70
	if s == "" {
		return "", errors.New("ICAO cannot be empty")
	}
	if len(s) > max {
		return "", errors.New("ICAO must be at most 70 characters")
	}
	return ICAO(s), nil
}
