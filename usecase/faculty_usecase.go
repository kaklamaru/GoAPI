package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"fmt"
)

// FacultyUsecase interface สำหรับการทำงานที่เกี่ยวข้องกับ Faculty
// ใช้ interface เพื่อความยืดหยุ่นในการทดสอบและใช้งาน
type FacultyUsecase interface {
	AddFaculty(faculty *entities.Faculty) error    // เพิ่มข้อมูลคณะ
	UpdateFaculty(faculty *entities.Faculty) error // แก้ไขข้อมูลคณะ
	GetAllFaculties() ([]entities.Faculty, error)  // แสดงคณะทั้งหมด
	GetFaculty(id uint) (*entities.Faculty, error) // ค้นหาคณะตาม ID
    DeleteFacultyByID(id uint) (*entities.Faculty, error)
	AddFacultyStaff(facultyID uint,userID uint) error
}

// facultyUsecase struct ซึ่งจะใช้งาน repository ในการดึงข้อมูลจากฐานข้อมูล
type facultyUsecase struct {
	repo repository.FacultyRepository // ใช้ repository เพื่อทำงานกับฐานข้อมูล
	userRepo repository.UserRepository
}

// NewFacultyUsecase สร้าง instance ของ FacultyUsecase
func NewFacultyUsecase(repo repository.FacultyRepository,userRepo repository.UserRepository) FacultyUsecase {
	return &facultyUsecase{
		repo: repo,
		userRepo: userRepo,
		} // ส่งคืน facultyUsecase พร้อม repository
}

// AddFaculty เพิ่มคณะใหม่
func (u *facultyUsecase) AddFaculty(faculty *entities.Faculty) error {
	// เรียกใช้ repository เพื่อเพิ่มคณะในฐานข้อมูล
	return u.repo.CreateFaculty(faculty)
}

// UpdateFaculty แก้ไขข้อมูลคณะ
func (u *facultyUsecase) UpdateFaculty(faculty *entities.Faculty) error {
	return u.repo.UpdateFaculty(faculty)
}

// GetAllFaculties ดึงข้อมูลคณะทั้งหมด
func (u *facultyUsecase) GetAllFaculties() ([]entities.Faculty, error) {
	// เรียกใช้ repository เพื่อดึงข้อมูลคณะทั้งหมดจากฐานข้อมูล
	return u.repo.GetAllFaculties()
}

// GetFaculty ค้นหาคณะตาม ID
func (u *facultyUsecase) GetFaculty(id uint) (*entities.Faculty, error) {
	// เรียกใช้ repository เพื่อตรวจสอบข้อมูลคณะตาม ID
	return u.repo.GetFacultyByID(id)
}

func (u *facultyUsecase) DeleteFacultyByID(id uint) (*entities.Faculty, error) {
	faculty, err := u.repo.DeleteFacultyByID(id)
	if err != nil {
		return nil, err
	}
	return faculty, nil
}

func (u *facultyUsecase) AddFacultyStaff(facultyID uint,userID uint) error{
	faculty,err := u.repo.GetFacultyByID(facultyID)
	if err != nil {
		return fmt.Errorf("faculty not found")
	}
	teacher ,err:= u.userRepo.GetTeacherByUserID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	faculty.SuperUser=&teacher.UserID
	return u.repo.AddFacultyStaff(faculty)
}
