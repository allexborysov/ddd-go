package inventoryifc

import (
	inventorysvc "github.com/allexborysov/aircraft/internal/application/inventory"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/allexborysov/aircraft/internal/infrastructure/validation"
)

type RegisterAircraftRequest struct {
	MSN   string   `json:"msn" validate:"required"`
	Seats []string `json:"seats" validate:"required,dive,required"`
}

func (req *RegisterAircraftRequest) Validate() error {
	return validation.Validate.Struct(req)
}

func (req *RegisterAircraftRequest) ToRegisterAircraftCommand() *inventorysvc.RegisterAircraftCommand {
	return &inventorysvc.RegisterAircraftCommand{
		MSN:   req.MSN,
		Seats: req.Seats,
	}
}

type RegisterAircraftResponse struct {
	Aircraft inventory.Aircraft `json:"aircraft"`
}

func ToRegisterAircraftResponse(result inventorysvc.RegisterAircraftCommandResult) RegisterAircraftResponse {
	return RegisterAircraftResponse{
		Aircraft: *result.Aircraft,
	}
}
