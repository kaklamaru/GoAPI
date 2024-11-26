package repository

import (
	"RESTAPI/domain/entities"
	"gorm.io/gorm"
)

type FacultyRepository interface {
    CreateFaculty(faculty *entities.Faculty) error
    UpdateFaculty(faculty *entities.Faculty) error
	GetAllFaculties() ([]entities.Faculty, error)
    GetFacultyByID(id uint) (*entities.Faculty, error)
}

type facultyRepository struct{
	db *gorm.DB
}

// NewFacultyRepository สร้าง instance ของ FacultyRepository
func NewFacultyRepository(db *gorm.DB) FacultyRepository {
    return &facultyRepository{db: db}
}

// CreateFaculty เพิ่มข้อมูลคณะ
func (r *facultyRepository) CreateFaculty(faculty *entities.Faculty) error {
    return r.db.Create(faculty).Error
}

// UpdateFaculty แก้ไขข้อมูลคณะ
func (r *facultyRepository) UpdateFaculty(faculty *entities.Faculty) error {
    return r.db.Save(faculty).Error
}

// GetAllFaculties แสดงข้อมูลคณะทั้งหมด
func (r *facultyRepository) GetAllFaculties() ([]entities.Faculty, error) {
    var faculties []entities.Faculty
    if err := r.db.Find(&faculties).Error; err != nil {
        return nil, err
    }
    return faculties, nil
}

// GetFacultyByID ค้นหาข้อมูลคณะตามรหัส
func (r *facultyRepository) GetFacultyByID(id uint) (*entities.Faculty, error) {
    var faculty entities.Faculty
    if err := r.db.First(&faculty, id).Error; err != nil {
        return nil, err
    }
    return &faculty, nil
}
