package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(tx transaction.Transaction, event *entities.Event) error
	GetAllEvent() ([]entities.Event, error)
	EditEvent(event *entities.Event) error
	GetEventByID(id uint) (*entities.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(tx transaction.Transaction, event *entities.Event) error {
	gormTx := tx.(*transaction.GormTransaction)
	if err := gormTx.GetDB().Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetAllEvent() ([]entities.Event, error) {
	var events []entities.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) EditEvent(event *entities.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) GetEventByID(id uint) (*entities.Event, error) {
	var event entities.Event
	if err := r.db.Find(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}







type PermissionRepository interface {
	CreatePermission(tx transaction.Transaction, permission *entities.Permission) error
	DeletePermission(eventID uint) error
}
type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) CreatePermission(tx transaction.Transaction, permission *entities.Permission) error {
    gormTx := tx.(*transaction.GormTransaction)

    // ตรวจสอบและตั้งค่าของ AllowAllBranch และ AllowAllYear
    if permission.AllowAllBranch {
        permission.BranchIDs = "" // กำหนดเป็นค่าว่างหาก AllowAllBranch = true
    }

    if permission.AllowAllYear {
        permission.YearIDs = "" // กำหนดเป็นค่าว่างหาก AllowAllYear = true
    }

    // แปลงค่า BranchIDs และ YearIDs เป็น JSON string ก่อนที่จะบันทึก
    branchIDsJSON, err := json.Marshal(permission.BranchIDs)
    if err != nil {
        return fmt.Errorf("failed to marshal branch IDs: %w", err)
    }
    permission.BranchIDs = string(branchIDsJSON)

    yearIDsJSON, err := json.Marshal(permission.YearIDs)
    if err != nil {
        return fmt.Errorf("failed to marshal year IDs: %w", err)
    }
    permission.YearIDs = string(yearIDsJSON)

    // สร้าง Permission ใหม่ในฐานข้อมูล
    if err := gormTx.GetDB().Create(permission).Error; err != nil {
        return fmt.Errorf("failed to create permission: %w", err)
    }

    return nil
}





func (r *permissionRepository) DeletePermission(eventID uint) error {
	// ลบ permission ที่เกี่ยวข้องกับ eventID
	return r.db.Where("event_id = ?", eventID).Delete(&entities.Permission{}).Error
}
