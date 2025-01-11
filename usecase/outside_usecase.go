package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	filesystem "RESTAPI/utility/fileSystem"
	"fmt"

	// "RESTAPI/usecase"
	"RESTAPI/utility"
	// "time"
)

type OutsideUsecase interface {
	CreateOutside(req entities.OutsideRequest, userID uint) (uint, error) 
	GetOutsideByID(id uint) (*entities.OutsideResponse, error)
	CreateFile(id uint) ([]byte, string, error)
}

type outsideUsecase struct {
	repo repository.OutsideRepository
}

func NewOutsideUsecase(repo repository.OutsideRepository) OutsideUsecase {
	return &outsideUsecase{
		repo: repo,
	}
}

func (u *outsideUsecase) CreateOutside(req entities.OutsideRequest, userID uint) (uint, error) {

	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return 0, err
	}
	outside := &entities.EventOutside{
		User:        userID,
		EventName:   req.EventName,
		StartDate:   startDate,
		Location:    req.Location,
		WorkingHour: req.WorkingHour,
	}
	id, err := u.repo.CreateOutside(outside)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (u *outsideUsecase) GetOutsideByID(id uint) (*entities.OutsideResponse, error) {
	outside, err := u.repo.GetOutsideByID(id)
	if err != nil {
		return nil, err
	}
	outsideRes := entities.OutsideResponse{
		EventID:     outside.EventID,
		EventName:   outside.EventName,
		Location:    outside.Location,
		StartDate:   outside.StartDate,
		WorkingHour: outside.WorkingHour,
		Intendant:   outside.Intendant,
		Student: entities.StudentResponse{
			UserID:      outside.Student.UserID,
			TitleName:   outside.Student.TitleName,
			FirstName:   outside.Student.FirstName,
			LastName:    outside.Student.LastName,
			Phone:       outside.Student.Phone,
			Code:        outside.Student.Code,
			BranchName:  outside.Student.Branch.BranchName,
			FacultyName: outside.Student.Branch.Faculty.FacultyName,
		},
	}
	return &outsideRes, nil

}

func (u *outsideUsecase) CreateFile(id uint) ([]byte, string, error){
	data, err := u.GetOutsideByID(id)
	if err != nil {
		return nil, "", fmt.Errorf("data not found: %v", err)
	}
	pdfBytes, fileName, err := filesystem.CreatePDF(*data)
	if err != nil {
		return nil, " " , fmt.Errorf("error creating PDF: %v", err)
	}

	// คืนค่าข้อมูลไฟล์ PDF ที่สร้างขึ้นมา (ข้อมูลแบบ byte array และชื่อไฟล์)
	return pdfBytes, fileName, nil
	
}
