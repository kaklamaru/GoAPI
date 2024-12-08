package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"RESTAPI/utility"
)

type EventUsecase interface {
	CreateEvent(event *entities.Event) error
	// GetAllEvent() ([]entities.Event, error)
	GetAllEvent() ([]entities.EventResponse, error)
	EditEvent(event *entities.Event) error
	// GetEventByID(id uint) (*entities.Event, error)
	GetEventByID(id uint) (*entities.EventResponse, error)
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

func (u *eventUsecase) GetAllEvent() ([]entities.EventResponse, error) {
	events,err := u.eventRepo.GetAllEvent()
	if err != nil {
		return nil, err
	}
	var res []entities.EventResponse
	for _, event := range events {
		// แปลงข้อมูลแต่ละรายการ
		mappedEvent, err := mapEventResponse(event)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}


func (u *eventUsecase) GetEventByID(id uint) (*entities.EventResponse, error) {
	event, err := u.eventRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	
	return mapEventResponse(*event)
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


func mapEventResponse(event entities.Event) (*entities.EventResponse, error) {
	branches, err := utility.DecodeIDs(event.BranchIDs)
	if err != nil {
		return &entities.EventResponse{}, err
	}

	years, err := utility.DecodeIDs(event.Years)
	if err != nil {
		return &entities.EventResponse{}, err
	}

	return &entities.EventResponse{
		EventID:        event.EventID,
		EventName:      event.EventName,
		Creator:        event.Creator,
		StartDate:      event.StartDate,
		WorkingHour:    event.WorkingHour,
		Limit:          event.Limit,
		Detail:         event.Detail,
		BranchIDs:      branches,
		Years:          years,
		AllowAllBranch: event.AllowAllBranch,
		AllowAllYear:   event.AllowAllYear,
		Teacher: struct {
			UserID    uint   `json:"user_id"`
			TitleName string `json:"title_name"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
			Code      string `json:"code"`
		}{
			UserID:    event.Teacher.UserID,
			TitleName: event.Teacher.TitleName,
			FirstName: event.Teacher.FirstName,
			LastName:  event.Teacher.LastName,
			Phone:     event.Teacher.Phone,
			Code:      event.Teacher.Code,
		},
	}, nil
}