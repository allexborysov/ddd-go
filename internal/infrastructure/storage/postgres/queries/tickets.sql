-- name: UpsertTicket :exec
INSERT INTO tickets (id, flight_id, passenger_id, seat, price)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id)
DO UPDATE SET
    flight_id = EXCLUDED.flight_id,
    passenger_id = EXCLUDED.passenger_id,
    seat = EXCLUDED.seat,
    price = EXCLUDED.price;

-- name: GetTicketByID :one
SELECT id, flight_id, passenger_id, seat, price
FROM tickets
WHERE id = $1;
