package service

import (
	"context"
	"errors"
	"fmt"
	"valeth-twice-management-api/internal/redisclient"
	"valeth-twice-management-api/internal/repo"

	redisv9 "github.com/redis/go-redis/v9" // ✅ alias untuk hindari bentrok
	"gorm.io/gorm"
)

type ReservationService struct {
    Repo *repo.EventRepo
    DB   *gorm.DB
}

func NewReservationService(r *repo.EventRepo, db *gorm.DB) *ReservationService {
    return &ReservationService{Repo: r, DB: db}
}

// 1) Pure DB transaction
func (s *ReservationService) ReserveWithDB(eventID uint, userID string) error {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        return s.Repo.ReserveTicketTx(tx, eventID, userID)
    })
}

// 2) Redis atomic decrement + persist to DB
var decScript = redisv9.NewScript(`
local k = KEYS[1]
local val = tonumber(redis.call("GET", k) or "-1")
if val <= 0 then
return -1
else
return redis.call("DECR", k)
end`)

func (s *ReservationService) ReserveWithRedisAtomic(ctx context.Context, eventKey string, eventID uint, userID string) error {
    // Step 1: Redis atomic decrement
    result, err := decScript.Run(ctx, redisclient.Rdb, []string{eventKey}).Result()
    if err != nil {
        return err
    }

    v, ok := result.(int64)
    if !ok {
        return errors.New("unexpected redis result")
    }
    if v < 0 {
        return errors.New("sold out (redis)")
    }

    // Step 2: DB transaction
    txErr := s.DB.Transaction(func(tx *gorm.DB) error {
        return s.Repo.ReserveTicketTx(tx, eventID, userID)
    })

    // Step 3: Compensation if DB fails
    if txErr != nil {
        if incErr := redisclient.Rdb.Incr(ctx, eventKey).Err(); incErr != nil {
            return fmt.Errorf("db error: %v; redis compensation error: %v", txErr, incErr)
        }
        return txErr
    }

    return nil
}
