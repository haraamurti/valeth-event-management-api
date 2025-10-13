package handler

import (
	"context"
	"strconv"
	"valeth-twice-management-api/internal/config"
	"valeth-twice-management-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReservationHandler struct {
    Svc *service.ReservationService
}

func NewReservationHandler(s *service.ReservationService) *ReservationHandler {
    return &ReservationHandler{Svc: s}
}

type ReserveRequest struct {
    UserID string `json:"user_id"`
}

// route: POST /events/:id/reserve
func (h *ReservationHandler) Reserve(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id64, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil { return c.Status(400).JSON(fiber.Map{"error":"invalid event id"}) }

    var req ReserveRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error":"invalid body"})
    }
    // Choose approach: env var USE_REDIS=true to use Redis approach
    useRedis := config.Get("USE_REDIS","false")
    if useRedis == "true" {
        key := "event:" + strconv.FormatUint(id64,10) + ":stock"
        if err := h.Svc.ReserveWithRedisAtomic(context.Background(), key, uint(id64), req.UserID); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"status":"reserved (redis)"})
    } else {
        if err := h.Svc.ReserveWithDB(uint(id64), req.UserID); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"status":"reserved (db)"})
    }
}
