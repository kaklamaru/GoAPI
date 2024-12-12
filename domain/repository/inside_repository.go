package repository

import (
	"RESTAPI/domain/entities"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventInsideRepository interface {
	JoinEventInside(eventInside *entities.EventInside) error
	UnJoinEventInside(eventID uint, userID uint) error
	UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error 
	CountEventInside(eventID uint) (uint, error)

}

type eventInsideRepository struct {
	db *gorm.DB
}

func NewEventInsideRepository(db *gorm.DB) EventInsideRepository {
	return &eventInsideRepository{db: db}
}

func (r *eventInsideRepository) JoinEventInside(eventInside *entities.EventInside) error {
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // กำหนด Timeout ให้ Lock
    if err := tx.Exec("SET innodb_lock_wait_timeout = 5").Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to set lock timeout: %w", err)
    }

    var event entities.Event
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
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
    if err := tx.Save(&event).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update event free space: %w", err)
    }

    if err := tx.Create(eventInside).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create event inside record: %w", err)
    }

    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}


func (r *eventInsideRepository) UnJoinEventInside(eventID uint, userID uint) error {
	// ลบ Event ที่มี event_id และ user_id ตรงกัน
	if err := r.db.Where("event_id = ? AND user = ?", eventID, userID).Delete(&entities.EventInside{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event with ID %d and user ID %d not found", eventID, userID)
		}
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (r *eventInsideRepository) UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error {
	// สร้าง map สำหรับฟิลด์ที่ต้องการอัปเดต
	updates := map[string]interface{}{
		"status":  status,
		"comment": comment,
	}

	// อัปเดตเรคคอร์ดที่ตรงกับเงื่อนไข
	if err := r.db.Model(&entities.EventInside{}).
		Where("event_id = ? AND user = ?", eventID, userID).
		Updates(updates).Error; err != nil {

		return fmt.Errorf("failed to update event: %w", err)

	}
	return nil
}

func (r *eventInsideRepository) CountEventInside(eventID uint) (uint, error) {
	var count int64

	// ใช้ Count เพื่อคำนวณจำนวนเรคคอร์ด
	if err := r.db.Model(&entities.EventInside{}).Where("event_id = ?", eventID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}

	return uint(count), nil
}

