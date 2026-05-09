-- name: UpsertAircraft :exec
INSERT INTO aircrafts (id, seats)
VALUES ($1, $2)
ON CONFLICT (id)
DO UPDATE SET seats = EXCLUDED.seats;

-- name: GetAircraftByID :one
SELECT id, seats FROM aircrafts WHERE id = $1;
