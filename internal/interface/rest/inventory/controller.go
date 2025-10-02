package inventoryifc

import (
	"net/http"

	inventorysvc "github.com/allexborysov/aircraft/internal/application/inventory"

	"github.com/labstack/echo/v4"
)

type InventoryController struct {
	service inventorysvc.Service
}

func NewInventoryController(router *echo.Echo, service inventorysvc.Service) *InventoryController {
	controller := &InventoryController{service}
	router.POST("/api/inventory/aircraft", controller.RegisterAircraftController)

	return controller
}

func (controller *InventoryController) RegisterAircraftController(ctx echo.Context) error {
	req := new(RegisterAircraftRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to parse request body",
		})
	}

	err := req.Validate()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	command := req.ToRegisterAircraftCommand()
	result, err := controller.service.RegisterAircraft(ctx.Request().Context(), command)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	response := ToRegisterAircraftResponse(*result)
	return ctx.JSON(http.StatusCreated, response)
}
