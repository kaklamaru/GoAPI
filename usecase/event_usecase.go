package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"RESTAPI/utility"
	"context"
	"encoding/json"
	"fmt"
)

type EventUsecase interface {
	CreateEvent(req *EventRequest, userID uint) error 
	GetAllEvent() ([]entities.EventResponse, error)
	EditEvent(ctx context.Context,eventID uint,req *EventRequest, userID uint) error
	GetEventByID(id uint) (*entities.EventResponse, error)
	CheckBranch(branchID uint) (bool, error)
	DeleteEvent(eventID uint, userID uint) error
	buildPermission(branches []uint, years []uint) (*Permission, error)
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
	var permission *Permission

	if len(branches) > 0 {
		for _, branchID := range branches {
			exists, err := u.branchRepo.BranchExists(branchID)
			if err != nil {
				return nil, fmt.Errorf("error checking branch: %v", err)
			}
			if !exists {
				return nil, fmt.Errorf("branch with ID %d does not exist", branchID)
			}
		}
	}

	if len(branches) > 0 && len(years) > 0 {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		permission = &Permission{
			BranchIDs:      string(branchData),
			Years:          string(yearData),
			AllowAllBranch: false,
			AllowAllYear:   false,
		}
	} else if len(branches) > 0 {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		permission = &Permission{
			BranchIDs:      string(branchData),
			Years:          "",
			AllowAllBranch: false,
			AllowAllYear:   true,
		}
	} else if len(years) > 0 {
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		permission = &Permission{
			BranchIDs:      "",
			Years:          string(yearData),
			AllowAllBranch: true,
			AllowAllYear:   false,
		}
	} else {
		permission = &Permission{
			BranchIDs:      "",
			Years:          "",
			AllowAllBranch: true,
			AllowAllYear:   true,
		}
	}

	return permission, nil
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
        FreeSpace:      req.FreeSpace,
        WorkingHour:    req.WorkingHour,
        Detail:         req.Detail,
        Creator:        userID,
        BranchIDs:      permission.BranchIDs,
        Years:          permission.Years,
        AllowAllBranch: permission.AllowAllBranch,
        AllowAllYear:   permission.AllowAllYear,
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
		// แปลงข้อมูลแต่ละรายการ
		count, err := u.insideRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event,count)
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

func (u *eventUsecase) EditEvent(ctx context.Context,eventID uint,req *EventRequest, userID uint) error {
	event, err := u.eventRepo.GetEventByID(eventID)
    if err != nil {
        return fmt.Errorf("event not found")
    }

    // ตรวจสอบว่า User มีสิทธิ์แก้ไข Event หรือไม่
    if event.Creator != userID {
        return fmt.Errorf("you do not have permission to edit this event")
    }

    // แปลง StartDate
    startDate, err := utility.ParseStartDate(req.StartDate)
    if err != nil {
        return fmt.Errorf("invalid start date format")
    }

    // สร้าง permission
    permission, err := u.buildPermission(req.Branches, req.Years)
    if err != nil {
        return fmt.Errorf("failed to build permissions: %w", err)
    }

    event.EventName = req.EventName
    event.StartDate = startDate
    event.FreeSpace = req.FreeSpace
    event.WorkingHour = req.WorkingHour
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
		Creator:        event.Creator,
		StartDate:      event.StartDate,
		WorkingHour:    event.WorkingHour,
		Limit:          limit,
		FreeSpace:      event.FreeSpace,
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
