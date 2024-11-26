package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
)

// FacultyUsecase interface สำหรับการทำงานที่เกี่ยวข้องกับ Faculty
// ใช้ interface เพื่อความยืดหยุ่นในการทดสอบและใช้งาน
type FacultyUsecase interface {
    AddFaculty(faculty *entities.Faculty) error       // เพิ่มข้อมูลคณะ
    UpdateFaculty(faculty *entities.Faculty) error    // แก้ไขข้อมูลคณะ
    GetAllFaculties() ([]entities.Faculty, error)    // แสดงคณะทั้งหมด
    GetFaculty(id uint) (*entities.Faculty, error)    // ค้นหาคณะตาม ID
}

// facultyUsecase struct ซึ่งจะใช้งาน repository ในการดึงข้อมูลจากฐานข้อมูล
type facultyUsecase struct {
    repo repository.FacultyRepository  // ใช้ repository เพื่อทำงานกับฐานข้อมูล
}

// NewFacultyUsecase สร้าง instance ของ FacultyUsecase
func NewFacultyUsecase(repo repository.FacultyRepository) FacultyUsecase {
    return &facultyUsecase{repo: repo}  // ส่งคืน facultyUsecase พร้อม repository
}

// AddFaculty เพิ่มคณะใหม่
func (u *facultyUsecase) AddFaculty(faculty *entities.Faculty) error {
    // เรียกใช้ repository เพื่อเพิ่มคณะในฐานข้อมูล
    return u.repo.CreateFaculty(faculty)
}

// UpdateFaculty แก้ไขข้อมูลคณะ
func (u *facultyUsecase) UpdateFaculty(faculty *entities.Faculty) error {
    // เรียกใช้ repository เพื่ออัปเดตคณะในฐานข้อมูล
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