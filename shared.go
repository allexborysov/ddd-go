package shared

import "time"

type Location struct {
	Latitude  float64
	Longitude float64
}

// Store IANA ID like "Europe/Paris", "UTC", "America/New_York"
type TimeZone string

func (tz TimeZone) ToLocation() (*time.Location, error) {
	return time.LoadLocation(string(tz))
}

// Represent amount in any currency
type Amount float64
