package repository

import (
	"RESTAPI/domain/entities"

	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(event *entities.Event) error
	GetAllEvent() ([]entities.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(event *entities.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetAllEvent() ([]entities.Event, error) {
	var events []entities.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
