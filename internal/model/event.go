package model

import "time"

type Event struct {
    ID               uint `gorm:"primaryKey"`
    Name             string
    TotalTickets     int
    AvailableTickets int
    CreatedAt        time.Time
}
