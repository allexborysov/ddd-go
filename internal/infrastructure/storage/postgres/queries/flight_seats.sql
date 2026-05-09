-- name: UpsertFlightSeat :exec
INSERT INTO flight_seats (flight_id, seat_number, price)
VALUES ($1, $2, $3)
ON CONFLICT (flight_id, seat_number)
DO UPDATE SET
    price = EXCLUDED.price;

-- name: GetFlightSeatsByFlightID :many
SELECT fs.flight_id, fs.seat_number, fs.price, t.passenger_id
FROM flight_seats fs
LEFT JOIN tickets t ON t.flight_id = fs.flight_id AND t.seat = fs.seat_number
WHERE fs.flight_id = $1
ORDER BY fs.seat_number;
