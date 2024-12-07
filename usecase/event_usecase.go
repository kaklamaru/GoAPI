package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
)

type EventUsecase interface {
	CreateEvent(event *entities.Event) error
	GetAllEvent() ([]entities.Event, error)
	EditEvent(event *entities.Event) error
	GetEventByID(id uint) (*entities.Event, error)
	CheckBranch(branchID uint) (bool, error)
	DeleteEvent(id uint) error 

}

type eventUsecase struct {
	eventRepo      repository.EventRepository
	branchRepo     repository.BranchRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, branchRepo repository.BranchRepository) EventUsecase {
	return &eventUsecase{
		eventRepo:      eventRepo,
		branchRepo:     branchRepo,
	}
}
func (u *eventUsecase) CreateEvent(event *entities.Event) error {
	return u.eventRepo.CreateEvent(event)
}

func (u *eventUsecase) GetAllEvent() ([]entities.Event, error) {
	return u.eventRepo.GetAllEvent()
}

func (u *eventUsecase) GetEventByID(id uint) (*entities.Event, error) {
	return u.eventRepo.GetEventByID(id)
}

func (u *eventUsecase) EditEvent(event *entities.Event) error {
	return u.eventRepo.EditEvent(event)
}

func (u *eventUsecase) DeleteEvent(id uint) error {
	return u.eventRepo.DeleteEvent(id)
}

func (u *eventUsecase) CheckBranch(branchID uint) (bool, error) {
    return u.branchRepo.BranchExists(branchID)
}
