package main

import (
	"fmt"
	"os"
	"valeth-twice-management-api/internal/config"
	"valeth-twice-management-api/internal/db"
	"valeth-twice-management-api/internal/handler"
	redis "valeth-twice-management-api/internal/redisclient"
	"valeth-twice-management-api/internal/repo"
	"valeth-twice-management-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
    config.Load()
dsn := config.Get("PG_URL", "")
if dsn == "" {
    panic("PG_URL not set in .env")
}

    dbConn, err := db.Connect(dsn)
    if err != nil {
        panic(err)
    }

    fmt.Println("✅ Database connected!")
    fmt.Println("🚀 Running AutoMigrate...")
    if err := db.AutoMigrate(dbConn); err != nil {
    panic(err)
    }
    fmt.Println("✅ AutoMigrate done!")
    

    // init redis
    redis.Init()

    // Optional: preload available_tickets into Redis if USE_REDIS=true
    if config.Get("USE_REDIS","false") == "true" {
        // get event row and set redis key
        var e repo.EventRepo
        e = *repo.NewEventRepo(dbConn)
        ev, err := e.GetByID(1)
        if err != nil {
            panic(err)
        }
        redis.Rdb.Set(redis.Ctx, "event:1:stock", ev.AvailableTickets, 0)
    }

    eventRepo := repo.NewEventRepo(dbConn)
    resSvc := service.NewReservationService(eventRepo, dbConn)
    resHandler := handler.NewReservationHandler(resSvc)

    app := fiber.New()
    app.Post("/events/:id/reserve", resHandler.Reserve)

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    app.Listen(":" + port)
}
