-- name: UpsertFlightSeat :exec
INSERT INTO flight_seats (flight_id, seat_number, price, passenger_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (flight_id, seat_number)
DO UPDATE SET
    price = EXCLUDED.price,
    passenger_id = EXCLUDED.passenger_id;

-- name: GetFlightSeatsByFlightID :many
SELECT flight_id, seat_number, price, passenger_id
FROM flight_seats
WHERE flight_id = $1
ORDER BY seat_number;
