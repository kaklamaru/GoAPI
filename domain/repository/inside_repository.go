package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventInsideRepository interface {
	JoinEventInside(eventInside *entities.EventInside, txManager transaction.TransactionManager) error
	UnJoinEventInside(eventID uint, userID uint, txManager transaction.TransactionManager) error
	UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error
	CountEventInside(eventID uint) (uint, error)
	IsUserJoinedEvent(eventID uint, userID uint) (bool, error) 
}

type eventInsideRepository struct {
	db *gorm.DB
}

func NewEventInsideRepository(db *gorm.DB) EventInsideRepository {
	return &eventInsideRepository{db: db}
}

func (r *eventInsideRepository) JoinEventInside(eventInside *entities.EventInside, txManager transaction.TransactionManager) error {
    tx := txManager.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.GetDB().Exec("SET innodb_lock_wait_timeout = 5").Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to set lock timeout: %w", err)
    }

    var event entities.Event
    if err := tx.GetDB().Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("event_id = ?", eventInside.EventId).
        First(&event).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to fetch event: %w", err)
    }

    if event.FreeSpace <= 0 {
        tx.Rollback()
        return fmt.Errorf("no free space available for event")
    }

    event.FreeSpace -= 1
    if err := tx.GetDB().Save(&event).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update event free space: %w", err)
    }

    if err := tx.GetDB().Create(eventInside).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create event inside record: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (r *eventInsideRepository) UnJoinEventInside(eventID uint, userID uint, txManager transaction.TransactionManager) error {
    tx := txManager.Begin()

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.GetDB().Exec("SET innodb_lock_wait_timeout = 5").Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to set lock timeout: %w", err)
    }

    var event entities.Event
    if err := tx.GetDB().Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("event_id = ?", eventID).
        First(&event).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("event with ID %d not found", eventID)
        }
        return fmt.Errorf("failed to fetch event: %w", err)
    }

    if err := tx.GetDB().
        Where("event_id = ? AND user = ?", eventID, userID).
        Delete(&entities.EventInside{}).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("event with ID %d and user ID %d not found", eventID, userID)
        }
        return fmt.Errorf("failed to delete event: %w", err)
    }

    event.FreeSpace += 1
    if err := tx.GetDB().Save(&event).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update event free space: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (r *eventInsideRepository) UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error {
	updates := map[string]interface{}{
		"status":  status,
		"comment": comment,
	}

	if err := r.db.Model(&entities.EventInside{}).
		Where("event_id = ? AND user = ?", eventID, userID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

func (r *eventInsideRepository) CountEventInside(eventID uint) (uint, error) {
	var count int64

	if err := r.db.Model(&entities.EventInside{}).Where("event_id = ?", eventID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}
	return uint(count), nil
}

func (r *eventInsideRepository) IsUserJoinedEvent(eventID uint, userID uint) (bool, error) {
    var count int64
    if err := r.db.Model(&entities.EventInside{}).
        Where("event_id = ? AND user = ?", eventID, userID).
        Count(&count).Error; err != nil {
        return false, fmt.Errorf("failed to query event inside: %w", err)
    }

    return count > 0, nil
}


