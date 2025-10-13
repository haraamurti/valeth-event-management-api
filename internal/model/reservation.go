package model

import "time"

type Reservation struct {
    ID        uint `gorm:"primaryKey"`
    EventID   uint
    UserID    string
    CreatedAt time.Time
}
