-- name: UpsertFlight :exec
INSERT INTO flights (id, aircraft_id, origin, destination, scheduled_departure, scheduled_arrival, close_booking_buffer_ns)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id)
DO UPDATE SET
    aircraft_id = EXCLUDED.aircraft_id,
    origin = EXCLUDED.origin,
    destination = EXCLUDED.destination,
    scheduled_departure = EXCLUDED.scheduled_departure,
    scheduled_arrival = EXCLUDED.scheduled_arrival,
    close_booking_buffer_ns = EXCLUDED.close_booking_buffer_ns;

-- name: GetFlightByID :one
SELECT id, aircraft_id, origin, destination, scheduled_departure, scheduled_arrival, close_booking_buffer_ns
FROM flights
WHERE id = $1;
