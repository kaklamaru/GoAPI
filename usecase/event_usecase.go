package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"RESTAPI/utility"

	"encoding/json"
	"fmt"
)

type EventUsecase interface {
	CreateEvent(req *EventRequest, userID uint) error
	GetAllEvent() ([]entities.EventResponse, error)
	EditEvent(eventID uint, req *EventRequest, userID uint) error
	GetEventByID(id uint) (*entities.EventResponse, error)
	CheckBranch(branchID uint) (bool, error)
	DeleteEvent(eventID uint, userID uint) error
	buildPermission(branches []uint, years []uint) (*Permission, error)
	ToggleEventStatus(eventID uint, userID uint) error
	AllAllowedEvent() ([]entities.EventResponse, error)
	AllCurrentEvent() ([]entities.EventResponse, error)
	MyEvent(userID uint) ([]entities.EventResponse, error)
}
type Permission struct {
	BranchIDs      string `json:"branches"`
	Years          string `json:"years"`
	AllowAllBranch bool   `json:"allow_all_branch"`
	AllowAllYear   bool   `json:"allow_all_year"`
}

type EventRequest struct {
	EventName   string `json:"event_name"`
	StartDate   string `json:"start_date"`
	WorkingHour uint   `json:"working_hour"`
	SchoolYear  uint   `json:"school_year"`
	Location    string `json:"location"`
	FreeSpace   uint   `json:"free_space"`
	Detail      string `json:"detail"`
	Branches    []uint `json:"branches"`
	Years       []uint `json:"years"`
}

type eventUsecase struct {
	eventRepo  repository.EventRepository
	branchRepo repository.BranchRepository
	insideRepo repository.EventInsideRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, branchRepo repository.BranchRepository, insideRepo repository.EventInsideRepository) EventUsecase {
	return &eventUsecase{
		eventRepo:  eventRepo,
		branchRepo: branchRepo,
		insideRepo: insideRepo,
	}
}

func (u *eventUsecase) buildPermission(branches []uint, years []uint) (*Permission, error) {
	if len(branches) > 0 {
		if err := u.validateBranches(branches); err != nil {
			return nil, err
		}
	}

	createPermission := func(branches []uint, years []uint, allowAllBranch, allowAllYear bool) (*Permission, error) {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		return &Permission{
			BranchIDs:      string(branchData),
			Years:          string(yearData),
			AllowAllBranch: allowAllBranch,
			AllowAllYear:   allowAllYear,
		}, nil
	}

	switch {
	case len(branches) > 0 && len(years) > 0:
		return createPermission(branches, years, false, false)
	case len(branches) > 0:
		return createPermission(branches, nil, false, true)
	case len(years) > 0:
		return createPermission(nil, years, true, false)
	default:
		return createPermission(nil, nil, true, true)
	}
}

func (u *eventUsecase) validateBranches(branches []uint) error {
	for _, branchID := range branches {
		exists, err := u.branchRepo.BranchExists(branchID)
		if err != nil {
			return fmt.Errorf("error checking branch: %v", err)
		}
		if !exists {
			return fmt.Errorf("branch with ID %d does not exist", branchID)
		}
	}
	return nil
}

func (u *eventUsecase) CreateEvent(req *EventRequest, userID uint) error {
	permission, err := u.buildPermission(req.Branches, req.Years)
	if err != nil {
		return err
	}

	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return err
	}

	event := &entities.Event{
		EventName:      req.EventName,
		StartDate:      startDate,
		SchoolYear:     req.SchoolYear,
		FreeSpace:      req.FreeSpace,
		WorkingHour:    req.WorkingHour,
		Detail:         req.Detail,
		Location:       req.Location,
		Creator:        userID,
		AllowAllBranch: permission.AllowAllBranch,
		AllowAllYear:   permission.AllowAllYear,
		BranchIDs:      permission.BranchIDs,
		Years:          permission.Years,
	}

	return u.eventRepo.CreateEvent(event)
}

func (u *eventUsecase) GetAllEvent() ([]entities.EventResponse, error) {
	events, err := u.eventRepo.GetAllEvent()
	if err != nil {
		return nil, err
	}
	var res []entities.EventResponse
	for _, event := range events {
		count, err := u.insideRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}
func (u *eventUsecase) AllAllowedEvent() ([]entities.EventResponse, error) {
	events, err := u.eventRepo.AllAllowedEvent()
	if err != nil {
		return nil, err
	}
	var res []entities.EventResponse
	for _, event := range events {
		count, err := u.insideRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) MyEvent(userID uint) ([]entities.EventResponse, error) {
	events, err := u.eventRepo.MyEvent(userID)
	if err != nil {
		return nil, err
	}
	var res []entities.EventResponse
	for _, event := range events {
		count, err := u.insideRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) AllCurrentEvent() ([]entities.EventResponse, error) {
	events, err := u.eventRepo.AllCurrentEvent()
	if err != nil {
		return nil, err
	}
	var res []entities.EventResponse
	for _, event := range events {
		count, err := u.insideRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
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
	count, err := u.insideRepo.CountEventInside(id)
	if err != nil {
		return nil, err
	}
	return mapEventResponse(*event, count)
}

func (u *eventUsecase) EditEvent(eventID uint, req *EventRequest, userID uint) error {
	event, err := u.eventRepo.GetEventByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	if event.Creator != userID {
		return fmt.Errorf("you do not have permission to edit this event")
	}

	if event.FreeSpace == 0 {
		event.Status = false
	}

	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date format")
	}

	permission, err := u.buildPermission(req.Branches, req.Years)
	if err != nil {
		return fmt.Errorf("failed to build permissions: %w", err)
	}

	event.EventName = req.EventName
	event.StartDate = startDate
	event.FreeSpace = req.FreeSpace
	event.SchoolYear = req.SchoolYear
	event.WorkingHour = req.WorkingHour
	event.Location = req.Location
	event.Detail = req.Detail
	event.BranchIDs = permission.BranchIDs
	event.Years = permission.Years
	event.AllowAllBranch = permission.AllowAllBranch
	event.AllowAllYear = permission.AllowAllYear

	return u.eventRepo.EditEvent(event)
}

func (u *eventUsecase) DeleteEvent(eventID uint, userID uint) error {
	event, err := u.eventRepo.GetEventByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.Creator != userID {
		return fmt.Errorf("you do not have permission to edit this event")
	}

	return u.eventRepo.DeleteEvent(event.EventID)
}

func (u *eventUsecase) CheckBranch(branchID uint) (bool, error) {
	return u.branchRepo.BranchExists(branchID)
}

func mapEventResponse(event entities.Event, count uint) (*entities.EventResponse, error) {
	branches, err := utility.DecodeIDs(event.BranchIDs)
	if err != nil {
		return &entities.EventResponse{}, err
	}

	years, err := utility.DecodeIDs(event.Years)
	if err != nil {
		return &entities.EventResponse{}, err
	}
	limit := event.FreeSpace + count

	return &entities.EventResponse{
		EventID:        event.EventID,
		EventName:      event.EventName,
		StartDate:      utility.FormatToThaiDate(event.StartDate),
		StartTime:      utility.FormatToThaiTime(event.StartDate),
		SchoolYear: event.SchoolYear,
		WorkingHour:    event.WorkingHour,
		Limit:          limit,
		FreeSpace:      event.FreeSpace,
		Detail:         event.Detail,
		Location:       event.Location,
		BranchIDs:      branches,
		Status:         event.Status,
		Years:          years,
		AllowAllBranch: event.AllowAllBranch,
		AllowAllYear:   event.AllowAllYear,
		Creator: struct {
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

func (u *eventUsecase) ToggleEventStatus(eventID uint, userID uint) error {

	event, err := u.eventRepo.GetEventByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.Creator != userID {
		return fmt.Errorf("you do not have permission to edit this event")
	}
	return u.eventRepo.ToggleEventStatus(event.EventID)
}
