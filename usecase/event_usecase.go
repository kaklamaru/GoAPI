package usecase

import(
	"RESTAPI/domain/repository"
	"RESTAPI/domain/entities"
)

type EventUsecase interface{
	CreateEvent(event *entities.Event) error
	GetAllEvent() ([]entities.Event,error)

}

type eventUsecase struct{
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) EventUsecase{
	return &eventUsecase{repo: repo}
}

func (u *eventUsecase) CreateEvent(event *entities.Event) error{
	return u.repo.CreateEvent(event)
}

func (u *eventUsecase) GetAllEvent() ([]entities.Event,error){
	return u.repo.GetAllEvent()
}

