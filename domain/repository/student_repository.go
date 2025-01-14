package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"fmt"

	"gorm.io/gorm"
)

// Student Repository
type StudentRepository interface {
	CreateStudent(tx transaction.Transaction,student *entities.Student) error
	EditStudentByID(student *entities.Student) error
	GetAllStudent() ([]entities.Student, error)
	// GetStudentByUserID(id uint) (*entities.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *studentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) CreateStudent(tx transaction.Transaction,student *entities.Student) error {
	gormTx := tx.(*transaction.GormTransaction)
	if err := gormTx.GetDB().Create(student).Error; err != nil {
		return err 
	}
	return nil
}

func (r *userRepository) GetStudentByUserID(id uint) (*entities.Student, error) {
	var student entities.Student
	result := r.db.
		Preload("Branch.Faculty").
		Where("user_id = ?", id).
		First(&student)

	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}


func (r *studentRepository) EditStudentByID(student *entities.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) GetAllStudent() ([]entities.Student, error) {
	var student []entities.Student
	if err := r.db.Preload("Branch.Faculty").Find(&student).Error; err != nil {
		return nil, err
	}
	fmt.Println(student)
	return student, nil
}
