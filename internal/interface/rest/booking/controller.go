package bookingifc

import (
	"net/http"

	bookingsvc "github.com/allexborysov/aircraft/internal/application/booking"

	"github.com/labstack/echo/v4"
)

type BookingController struct {
	service bookingsvc.Service
}

func NewBookingController(router *echo.Echo, service bookingsvc.Service) *BookingController {
	controller := &BookingController{service}
	router.POST("/api/booking/flight", controller.BookFlightController)

	return controller
}

func (controller *BookingController) BookFlightController(ctx echo.Context) error {
	req := new(BookFlightRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to parse request body",
		})
	}

	err := req.Validate()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	command := req.ToBookFlightCommand()
	result, err := controller.service.BookFlight(ctx.Request().Context(), command)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	response := ToBookFlightResposnse(result)
	return ctx.JSON(http.StatusCreated, response)
}
