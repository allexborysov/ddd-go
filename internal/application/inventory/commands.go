package inventorysvc

import "github.com/allexborysov/aircraft/internal/domain/inventory"

type RegisterAircraftCommand struct {
	MSN   string
	Seats []string
}

type RegisterAircraftCommandResult struct {
	Aircraft *inventory.Aircraft
}
