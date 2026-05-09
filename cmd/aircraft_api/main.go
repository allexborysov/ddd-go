package main

import (
	"fmt"
	"log"
	"log/slog"

	redisv9 "github.com/redis/go-redis/v9"

	"github.com/allexborysov/aircraft/config"
	bookingsvc "github.com/allexborysov/aircraft/internal/application/booking"
	inventorysvc "github.com/allexborysov/aircraft/internal/application/inventory"
	schedulingsvc "github.com/allexborysov/aircraft/internal/application/scheduling"
	"github.com/allexborysov/aircraft/internal/domain/flight"
	"github.com/allexborysov/aircraft/internal/domain/inventory"
	"github.com/allexborysov/aircraft/internal/infrastructure/logger"
	"github.com/allexborysov/aircraft/internal/infrastructure/repos/inmem"
	pgrepos "github.com/allexborysov/aircraft/internal/infrastructure/repos/postgres"
	redis "github.com/allexborysov/aircraft/internal/infrastructure/storage"
	"github.com/allexborysov/aircraft/internal/infrastructure/storage/postgres"
	bookingsync "github.com/allexborysov/aircraft/internal/infrastructure/sync"

	"github.com/allexborysov/aircraft/internal/infrastructure/services/ticketspdf"
	bookingifc "github.com/allexborysov/aircraft/internal/interface/rest/booking"
	inventoryifc "github.com/allexborysov/aircraft/internal/interface/rest/inventory"
	schedulingifc "github.com/allexborysov/aircraft/internal/interface/rest/scheduling"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.New(cfg.Env)

	redis := redis.MustConnectRedis(&redisv9.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Repositories
	var (
		aircrafts inventory.AircraftRepository
		flights   flight.FlightRepository
		tickets   flight.TicketRepository
	)
	if cfg.InMemoryStorage {
		aircrafts = inmem.NewAircraftRepository()
		flights = inmem.NewFlightRepository()
		tickets = inmem.NewTicketRepository()
	} else {
		pg := postgres.MustConnectPostgres(fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName,
		))
		defer pg.Close()

		aircrafts = pgrepos.NewAircraftRepository(pg)
		flights = pgrepos.NewFlightRepository(pg)
		tickets = pgrepos.NewTicketRepository(pg)
	}

	// Sync
	bookingSync := bookingsync.New(redis)

	// Services
	inventory := inventorysvc.New(aircrafts)
	booking := bookingsvc.New(bookingSync, flights, tickets, ticketspdf.New())
	scheduling := schedulingsvc.New(flights, aircrafts)

	// Interface
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
