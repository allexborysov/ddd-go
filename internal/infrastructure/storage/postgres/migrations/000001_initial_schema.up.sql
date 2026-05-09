CREATE TABLE IF NOT EXISTS aircrafts (
    id TEXT PRIMARY KEY,
    seats TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS flights (
    id TEXT PRIMARY KEY,
    aircraft_id TEXT NOT NULL REFERENCES aircrafts(id),
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    scheduled_departure TIMESTAMPTZ NOT NULL,
    scheduled_arrival TIMESTAMPTZ NOT NULL,
    close_booking_buffer_ns BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS flight_seats (
    flight_id TEXT NOT NULL REFERENCES flights(id),
    seat_number TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (flight_id, seat_number)
);

CREATE TABLE IF NOT EXISTS tickets (
    id TEXT PRIMARY KEY,
    flight_id TEXT NOT NULL REFERENCES flights(id),
    passenger_id TEXT NOT NULL,
    seat TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL
);
