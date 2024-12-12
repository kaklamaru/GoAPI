package repository

import (
	"RESTAPI/domain/entities"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository interface {
	CreateEvent(event *entities.Event) error
	GetAllEvent() ([]entities.Event, error)
	EditEvent(event *entities.Event) error
	GetEventByID(id uint) (*entities.Event, error)
	DeleteEvent(id uint) error
	UpdateEventStatusIfNoSpace() error
	CanJoinEvent(eventID uint) (bool, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(event *entities.Event) error {
	branchIDsJSON, err := json.Marshal(event.BranchIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal branch IDs: %w", err)
	}

	yearIDsJSON, err := json.Marshal(event.Years)
	if err != nil {
		return fmt.Errorf("failed to marshal year IDs: %w", err)
	}

	event.BranchIDs = string(branchIDsJSON)
	event.Years = string(yearIDsJSON)

	if err := r.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetAllEvent() ([]entities.Event, error) {
	var events []entities.Event
	if err := r.db.Preload("Teacher").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) EditEvent(event *entities.Event) error {
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    var currentEvent entities.Event
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("event_id = ?", event.EventID).
        First(&currentEvent).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to lock event: %w", err)
    }

    var joinedCount int64
    if err := tx.Model(&entities.EventInside{}).
        Where("event_id = ?", event.EventID).
        Count(&joinedCount).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to count joined participants: %w", err)
    }

    if event.FreeSpace < uint(joinedCount)  {
        tx.Rollback()
        return fmt.Errorf("free space (%d) cannot be less than joined participants (%d)", event.FreeSpace, joinedCount)
    }

    branchIDsJSON, err := json.Marshal(event.BranchIDs)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to marshal branch IDs: %w", err)
    }
    yearIDsJSON, err := json.Marshal(event.Years)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to marshal year IDs: %w", err)
    }

    if err := tx.Model(&entities.Event{}).Where("event_id = ?", event.EventID).Updates(map[string]interface{}{
        "event_name":       event.EventName,
        "start_date":       event.StartDate,
        "free_space":       event.FreeSpace,
        "working_hour":     event.WorkingHour,
        "detail":           event.Detail,
        "branch_ids":       string(branchIDsJSON),
        "years":            string(yearIDsJSON),
        "allow_all_branch": event.AllowAllBranch,
        "allow_all_year":   event.AllowAllYear,
    }).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update event: %w", err)
    }

    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}



func (r *eventRepository) GetEventByID(id uint) (*entities.Event, error) {
	var event entities.Event

	// ค้นหา Event โดยใช้ ID
	if err := r.db.Preload("Teacher").First(&event, "event_id = ?", id).Error; err != nil {
		// ตรวจสอบว่าไม่พบข้อมูล (Record Not Found)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with ID %d not found", id)
		}
		// ถ้ามีข้อผิดพลาดอื่น ๆ
		return nil, fmt.Errorf("failed to retrieve event: %w", err)
	}

	return &event, nil
}

func (r *eventRepository) DeleteEvent(id uint) error {
	event := &entities.Event{}

	if err := r.db.First(event, "event_id = ?", id).Error; err != nil {
		// ถ้าไม่พบข้อมูล
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event with ID %d not found", id)
		}
		return fmt.Errorf("failed to find event: %w", err)
	}

	var eventInsideCount int64
	if err := r.db.Model(&entities.EventInside{}).Where("event_id = ?", id).Count(&eventInsideCount).Error; err != nil {
		return fmt.Errorf("failed to check event references: %w", err)
	}

	if eventInsideCount > 0 {
		return fmt.Errorf("cannot delete event with ID %d because it is referenced in event_insides", id)
	}

	if err := r.db.Delete(event).Error; err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (r *eventRepository) UpdateEventStatusIfNoSpace() error {
	err := r.db.Model(&entities.Event{}).
		Where("free_space = ?", 0).
		Update("status", false).Error

	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}
	return nil
}

func (r *eventRepository) CanJoinEvent(eventID uint) (bool, error) {
	var event entities.Event
	if err := r.db.Where("event_id = ? AND status = ?", eventID, true).First(&event).Error; err != nil {
		return false, fmt.Errorf("event not available: %w", err)
	}
	return true, nil
}
