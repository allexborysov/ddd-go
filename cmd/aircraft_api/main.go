package main

import (
	"log"
	"log/slog"

	redisv9 "github.com/redis/go-redis/v9"

	"github.com/allexborysov/aircraft/config"
	bookingsvc "github.com/allexborysov/aircraft/internal/application/booking"
	inventorysvc "github.com/allexborysov/aircraft/internal/application/inventory"
	schedulingsvc "github.com/allexborysov/aircraft/internal/application/scheduling"
	"github.com/allexborysov/aircraft/internal/infrastructure/logger"
	redis "github.com/allexborysov/aircraft/internal/infrastructure/storage"
	bookingsync "github.com/allexborysov/aircraft/internal/infrastructure/sync"

	"github.com/allexborysov/aircraft/internal/infrastructure/services/ticketspdf"
	"github.com/allexborysov/aircraft/internal/infrastructure/storage/inmem"
	bookingifc "github.com/allexborysov/aircraft/internal/interface/rest/booking"
	inventoryifc "github.com/allexborysov/aircraft/internal/interface/rest/inventory"
	schedulingifc "github.com/allexborysov/aircraft/internal/interface/rest/scheduling"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.New(cfg.Env)

	redis := redis.MustConnect(&redisv9.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	aircrafts := inmem.NewAircraftRepository()
	flights := inmem.NewFlightRepository()
	tickets := inmem.NewTicketRepository()

	bookingSync := bookingsync.New(redis)

	inventory := inventorysvc.New(aircrafts)
	booking := bookingsvc.New(bookingSync, flights, tickets, ticketspdf.New())
	scheduling := schedulingsvc.New(flights, aircrafts)

	e := echo.New()
	e.Use(middleware.Recover())

	inventoryifc.NewInventoryController(e, *inventory)
	bookingifc.NewBookingController(e, *booking)
	schedulingifc.NewSchedulingController(e, *scheduling)

	err := e.Start(cfg.HttpServer.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	logger.Info("Running..", slog.Any("config", cfg))
}
