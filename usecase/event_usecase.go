package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"RESTAPI/domain/transaction"
	"errors"
)

type EventUsecase interface {
	// CreateEvent(tx transaction.Transaction,event *entities.Event) error
	CreateEventWithPermission(tx transaction.Transaction, event *entities.Event, permission *entities.Permission) error
	GetAllEvent() ([]entities.Event, error)
	EditEvent(event *entities.Event) error
	GetEventByID(id uint) (*entities.Event, error)
	// CreatePermission(tx transaction.Transaction,permission *entities.Permission) error
	DeletePermission(eventID uint) error
}

type eventUsecase struct {
	eventRepo      repository.EventRepository
	permissionRepo repository.PermissionRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, permissionRepo repository.PermissionRepository) EventUsecase {
	return &eventUsecase{
		eventRepo:      eventRepo,
		permissionRepo: permissionRepo,
	}
}
func (u *eventUsecase) CreateEventWithPermission(tx transaction.Transaction, event *entities.Event, permission *entities.Permission) error {
	// สร้าง Event
	if err := u.eventRepo.CreateEvent(tx, event); err != nil {
		tx.Rollback()
		return err
	}

	// ตรวจสอบว่า event.EventID ได้รับการอัปเดต
	if event.EventID == 0 {
		tx.Rollback()
		return errors.New("failed to get valid EventID")
	}

	// เชื่อมโยง Permission กับ Event
	if permission != nil {
		permission.EventID = event.EventID
		// สร้าง Permission
		if err := u.permissionRepo.CreatePermission(tx, permission); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit ธุรกรรม
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
// func (u *eventUsecase) CreateEvent(tx transaction.Transaction, event *entities.Event) error {
// 	return u.eventRepo.CreateEvent(tx, event)
// }

func (u *eventUsecase) GetAllEvent() ([]entities.Event, error) {
	return u.eventRepo.GetAllEvent()
}

func (u *eventUsecase) EditEvent(event *entities.Event) error {
	return u.eventRepo.EditEvent(event)
}

func (u *eventUsecase) GetEventByID(id uint) (*entities.Event, error) {
	return u.eventRepo.GetEventByID(id)
}

// // CreatePermission สำหรับสร้าง permission ใหม่
// func (u *eventUsecase) CreatePermission(tx transaction.Transaction, permission *entities.Permission) error {
// 	return u.permissionRepo.CreatePermission(tx, permission)
// }

// DeletePermission สำหรับลบ permission ตาม eventID
func (u *eventUsecase) DeletePermission(eventID uint) error {
	return u.permissionRepo.DeletePermission(eventID)
}
