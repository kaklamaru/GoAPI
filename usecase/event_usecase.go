package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
)

type EventUsecase interface {
	CreateEvent(event *entities.Event) error
	GetAllEvent() ([]entities.Event, error)
	EditEvent(event *entities.Event) error
	GetEventByID(id uint) (*entities.Event,error)
}

type eventUsecase struct {
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) EventUsecase {
	return &eventUsecase{repo: repo}
}

func (u *eventUsecase) CreateEvent(event *entities.Event) error {
	return u.repo.CreateEvent(event)
}

func (u *eventUsecase) GetAllEvent() ([]entities.Event, error) {
	return u.repo.GetAllEvent()
}

func (u *eventUsecase) EditEvent(event *entities.Event) error {
	return u.repo.EditEvent(event)
}

func (u *eventUsecase) GetEventByID(id uint) (*entities.Event, error) {
	return u.repo.GetEventByID(id)
}
