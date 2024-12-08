package repository

import(
	"gorm.io/gorm"
	"RESTAPI/domain/entities"
)

type EventInsideRepository interface{
	CreateEventInside(eventInside *entities.EventInside) error
	
}

type eventInsideRepository struct{
	db *gorm.DB
}

func NewEventInsideRepository(db *gorm.DB) EventInsideRepository{
	return &eventInsideRepository{db: db}
}

func (r *eventInsideRepository) CreateEventInside(eventInside *entities.EventInside) error{
	return r.db.Create(eventInside).Error
}