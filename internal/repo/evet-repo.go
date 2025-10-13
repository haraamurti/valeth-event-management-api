package repo

import (
	"errors"
	"valeth-twice-management-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepo struct {
    DB *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo { return &EventRepo{DB: db} }

func (r *EventRepo) GetByID(id uint) (*model.Event, error) {
    var e model.Event
    if err := r.DB.First(&e, id).Error; err != nil {
        return nil, err
    }
    return &e, nil
}

// Pessimistic lock approach: SELECT ... FOR UPDATE inside a transaction
func (r *EventRepo) ReserveTicketTx(tx *gorm.DB, eventID uint, userID string) error {
    var e model.Event
    // lock row for update
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&e, eventID).Error; err != nil {
        return err
    }
    if e.AvailableTickets <= 0 {
        return errors.New("sold out")
    }
    e.AvailableTickets -= 1
    if err := tx.Save(&e).Error; err != nil {
        return err
    }
    res := model.Reservation{EventID: eventID, UserID: userID}
    if err := tx.Create(&res).Error; err != nil {
        return err
    }
    return nil
}
