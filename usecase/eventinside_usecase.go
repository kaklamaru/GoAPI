package usecase

import(
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"

)

type EventInsideUsecase interface{
	CreateEventInside(eventInside *entities.EventInside) error

}

type eventInsideUsecase struct{
	repo repository.EventInsideRepository
}

func NewEventInsideUsecase(repo repository.EventInsideRepository) EventInsideUsecase{
	return &eventInsideUsecase{
		repo: repo,
	}
}

func (u *eventInsideUsecase) CreateEventInside(eventInside *entities.EventInside) error{
	return u.repo.CreateEventInside(eventInside)
}