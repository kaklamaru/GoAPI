package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"

	"gorm.io/gorm"
)

// User Repository
type UserRepository interface {
	CreateUser(tx transaction.Transaction, user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetStudentByUserID(userID uint) (*entities.Student, error)
	GetTeacherByUserID(userID uint) (*entities.Teacher, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(tx transaction.Transaction, user *entities.User) error {
	gormTx := tx.(*transaction.GormTransaction)
	if err := gormTx.GetDB().Create(user).Error; err != nil {
		return err 
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
