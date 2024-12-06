package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"

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
    // Preload Permission เพื่อดึงข้อมูล Permission พร้อมกับ Event
    if err := r.db.Preload("Teacher").Preload("Permission").Find(&events).Error; err != nil {
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






