package repository

import (
	"RESTAPI/domain/entities"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(event *entities.Event) error
	GetAllEvent() ([]entities.Event, error)
	EditEvent(event *entities.Event) error
	GetEventByID(id uint) (*entities.Event, error)
	DeleteEvent(id uint) error
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
    // Preload Permission เพื่อดึงข้อมูล Permission พร้อมกับ Event
    if err := r.db.Preload("Teacher").Find(&events).Error; err != nil {
        return nil, err
    }
    return events, nil
}

func (r *eventRepository) EditEvent(event *entities.Event) error {
    // การแปลง BranchIDs และ Years ให้เป็น JSON string
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

    if event.EventID == 0 {
        return fmt.Errorf("event ID is missing")
    }

    err = r.db.Model(&entities.Event{}).Where("event_id = ?", event.EventID).Updates(map[string]interface{}{
        "event_name":      event.EventName,
        "start_date":      event.StartDate,
        "limit":           event.Limit,
        "working_hour":    event.WorkingHour,
        "detail":          event.Detail,
        "branch_ids":      event.BranchIDs,
        "years":           event.Years,
        "allow_all_branch": event.AllowAllBranch,
        "allow_all_year":  event.AllowAllYear,
    }).Error
    if err != nil {
        return fmt.Errorf("failed to update event: %w", err)
    }

    // ตรวจสอบข้อมูลหลังการอัปเดต
    var updatedEvent entities.Event
    if err := r.db.Where("event_id = ?", event.EventID).First(&updatedEvent).Error; err != nil {
        return fmt.Errorf("failed to fetch updated event: %w", err)
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

	// คืนค่าข้อมูล Event ที่พบ
	return &event, nil
}


func (r *eventRepository) DeleteEvent(id uint) error {
    // สร้างตัวแปรสำหรับ Event
	event := &entities.Event{}

	// ค้นหา Event ด้วย ID
	if err := r.db.First(event, "event_id = ?", id).Error; err != nil {
		// ถ้าไม่พบข้อมูล
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event with ID %d not found", id)
		}
		// ถ้ามีข้อผิดพลาดอื่น ๆ
		return fmt.Errorf("failed to find event: %w", err)
	}

	// ลบข้อมูล Event ที่ค้นพบ
	if err := r.db.Delete(event).Error; err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}





