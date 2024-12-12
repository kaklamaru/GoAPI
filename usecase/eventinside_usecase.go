package usecase

import(
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"

)

type EventInsideUsecase interface{
	JoinEventInside(eventInside *entities.EventInside) error
	UnJoinEventInside(eventID uint , userID uint) error
	UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error
	CountEventInside(eventID uint) (uint,error)

	
}

type eventInsideUsecase struct{
	repo repository.EventInsideRepository
}

func NewEventInsideUsecase(repo repository.EventInsideRepository) EventInsideUsecase{
	return &eventInsideUsecase{
		repo: repo,
	}
}

func (u *eventInsideUsecase) JoinEventInside(eventInside *entities.EventInside) error{
	return u.repo.JoinEventInside(eventInside)
}

func (u *eventInsideUsecase) UnJoinEventInside(eventID uint , userID uint) error{
	return u.repo.UnJoinEventInside(eventID,userID)
}

func (u *eventInsideUsecase) UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error{
	return u.repo.UpdateEventStatusAndComment(eventID,userID,status,comment)
}

func (u *eventInsideUsecase) CountEventInside(eventID uint) (uint,error){
	return u.repo.CountEventInside(eventID)
}