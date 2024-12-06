package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"

	"gorm.io/gorm"
)

// Teacher Repository
type TeacherRepository interface {
	CreateTeacher(tx transaction.Transaction, teacher *entities.Teacher) error
	EditTeacherByID(teacher *entities.Teacher) error
	GetAllTeacher() ([]entities.Teacher, error)
}

type teacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *teacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) CreateTeacher(tx transaction.Transaction, teacher *entities.Teacher) error {
	gormTx := tx.(*transaction.GormTransaction)
	return gormTx.GetDB().Create(teacher).Error
}
func (r *userRepository) GetTeacherByUserID(id uint) (*entities.Teacher, error) {
	var teacher entities.Teacher
	result := r.db.Where("user_id = ?", id).First(&teacher)
	if result.Error != nil {
		return nil, result.Error
	}
	return &teacher, nil
}
func (r *teacherRepository) EditTeacherByID(teacher *entities.Teacher) error {
	return r.db.Save(teacher).Error
}

func (r *teacherRepository) GetAllTeacher() ([]entities.Teacher, error) {
	var teacher []entities.Teacher
	if err := r.db.Find(&teacher).Error; err != nil {
		return nil, err
	}
	return teacher, nil
}
