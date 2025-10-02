package schedulingifc

import (
	"net/http"

	schedulingsvc "github.com/allexborysov/aircraft/internal/application/scheduling"

	"github.com/labstack/echo/v4"
)

type SchedulingController struct {
	service schedulingsvc.Service
}

func NewSchedulingController(router *echo.Echo, service schedulingsvc.Service) *SchedulingController {
	controller := &SchedulingController{service: service}
	router.POST("/api/scheduling/flight", controller.ScheduleFlightController)

	return controller
}

func (c *SchedulingController) ScheduleFlightController(ctx echo.Context) error {
	req := new(ScheduleFlightRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to parse request body",
		})
	}

	if err := req.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	command := req.ToScheduleFlightCommand()
	result, err := c.service.ScheduleFlight(ctx.Request().Context(), command)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	response := ToScheduleFlightResponse(result)
	return ctx.JSON(http.StatusCreated, response)
}
