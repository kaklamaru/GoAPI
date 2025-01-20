package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"RESTAPI/domain/transaction"
	// "RESTAPI/utility"
	"RESTAPI/utility/fileSystem"
	"fmt"
	"mime/multipart"
	"os"
)


type EventInsideUsecase interface{
	JoinEventInside(eventID uint, userID uint) error
	UnJoinEventInside(eventID uint , userID uint) error
	UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error
	CountEventInside(eventID uint) (uint,error)
	UploadFile(file *multipart.FileHeader, eventID uint, userID uint) error 	
	GetFile(eventID uint,userID uint) (string,error)
    MyChecklist(userID uint,eventID uint) ([]entities.MyChecklist,error)

}

type eventInsideUsecase struct{
	insideRepo repository.EventInsideRepository
	userRepo repository.UserRepository
	eventUsecase EventUsecase
	txManager transaction.TransactionManager
	
}

func NewEventInsideUsecase(insideRepo repository.EventInsideRepository,userRepo repository.UserRepository,eventUsecase EventUsecase,txManager transaction.TransactionManager) EventInsideUsecase{
	return &eventInsideUsecase{
		insideRepo: insideRepo,
		userRepo: userRepo,
		eventUsecase: eventUsecase,
		txManager: txManager,
	}
}

func checkPermission(permission *entities.EventResponse, user *entities.Student) bool {
	if permission == nil || user == nil {
		return false
	}
	permissionBranch := permission.AllowAllBranch
	permissionYear := permission.AllowAllYear

	if !permissionBranch && permission.BranchIDs != nil {
		for _, branch := range permission.BranchIDs {
			if user.BranchId == branch {
				permissionBranch = true
				break
			}
		}
	}

	if !permissionYear && permission.Years != nil {
		for _, year := range permission.Years {
			if user.Year == year {
				permissionYear = true
				break
			}
		}
	}

	return permissionBranch && permissionYear
}

func (u *eventInsideUsecase) JoinEventInside(eventID uint, userID uint) error {
    student, err := u.userRepo.GetStudentByUserID(userID)
    if err != nil || student == nil {
        return fmt.Errorf("student not found")
    }

    event, err := u.eventUsecase.GetEventByID(eventID)
    if err != nil || event == nil {
        return fmt.Errorf("event not found")
    }

    if !event.Status {
        return fmt.Errorf("event not allowed")
    }

    if event.FreeSpace == 0 {
        return fmt.Errorf("the event is full")
    }

    if !checkPermission(event, student) {
        return fmt.Errorf("user is not allowed to join this event")
    }

    eventInside := &entities.EventInside{
        EventId:   eventID,
        User:      userID,
        Status:    false,
        Certifier: event.Creator.UserID,
    }

    err = u.insideRepo.JoinEventInside(eventInside, u.txManager)
    if err != nil {
        return fmt.Errorf("failed to join event inside: %w", err)
    }
    return nil
}

func (u *eventInsideUsecase) UnJoinEventInside(eventID uint, userID uint) error {
    student, err := u.userRepo.GetStudentByUserID(userID)
    if err != nil || student == nil {
        return fmt.Errorf("user not found or not a student")
    }

    event, err := u.eventUsecase.GetEventByID(eventID)
    if err != nil || event == nil {
        return fmt.Errorf("event not found")
    }

    isMember, err := u.insideRepo.IsUserJoinedEvent(eventID, userID)
    if err != nil {
        return fmt.Errorf("failed to verify user participation: %w", err)
    }
    if !isMember {
        return fmt.Errorf("user is not a member of this event")
    }

    if err := u.insideRepo.UnJoinEventInside(eventID, userID, u.txManager); err != nil {
        return fmt.Errorf("failed to unjoin event: %w", err)
    }

    return nil
}

func (u *eventInsideUsecase) UploadFile(file *multipart.FileHeader, eventID uint, userID uint) error {
    if file.Header.Get("Content-Type") != "application/pdf" {
        return fmt.Errorf("only PDF files are allowed")
    }
    // ค้นหาว่า event นี้มีการบันทึกไฟล์ไว้หรือยัง
    currentFilePath, err := u.insideRepo.GetFilePathByEvent(eventID, userID)
    if err != nil {
        return fmt.Errorf("failed to fetch current file path: %w", err)
    }

    // ถ้ามีไฟล์เก่าอยู่ ให้ลบไฟล์เก่าก่อน 
    if currentFilePath != "" {
        removeErr := os.Remove(currentFilePath)
        if removeErr != nil {
            return fmt.Errorf("failed to remove old file: %v", removeErr)
        }
    }
    // บันทึกไฟล์ใหม่
    path, err := filesystem.SaveFile(file, userID)
    if err != nil {
        return fmt.Errorf("failed to save file: %w", err)
    }
    // อัปเดตฐานข้อมูลด้วย path ใหม่
    err = u.insideRepo.UpdateFile(eventID, userID, path)
    if err != nil {
        removeErr := os.Remove(path)
        if removeErr != nil {
            return fmt.Errorf("failed to update database and remove file: %v, cleanup error: %w", err, removeErr)
        }
        return fmt.Errorf("failed to update database: %w", err)
    }
    return nil
}

func (u *eventInsideUsecase) GetFile(eventID uint, userID uint) (string, error) {
    filePath, err := u.insideRepo.GetFilePath(eventID, userID)
    if err != nil {
        return "", err
    }
    return filePath, nil
}

func (u *eventInsideUsecase) MyChecklist(userID uint,eventID uint) ([]entities.MyChecklist,error){
    checklist,err:=u.insideRepo.MyChecklist(userID,eventID)
    if err != nil {
		return nil, err
	}
    var res []entities.MyChecklist
    for _, inside := range checklist {
		mappedEvent := entities.MyChecklist{
            EventID: inside.EventId,
            UserID: inside.User,
            TitleName: inside.Student.TitleName,
            FirstName: inside.Student.FirstName,
            LastName: inside.Student.LastName,
            Code: inside.Student.Code,
            Certifier: inside.Certifier,
            Status: inside.Status,
            Comment: inside.Comment,
            FilePDF: inside.FilePDF,
        }
		res = append(res,mappedEvent)
	}
    return res,nil
}

func (u *eventInsideUsecase) UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error{
	return u.insideRepo.UpdateEventStatusAndComment(eventID,userID,status,comment)
}

func (u *eventInsideUsecase) CountEventInside(eventID uint) (uint,error){
	return u.insideRepo.CountEventInside(eventID)
}
