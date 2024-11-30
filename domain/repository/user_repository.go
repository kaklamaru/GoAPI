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

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(tx transaction.Transaction, user *entities.User) error {
	gormTx := tx.(*transaction.GormTransaction)
	if err := gormTx.GetDB().Create(user).Error; err != nil {
		return err // ถ้ามีข้อผิดพลาดใน CreateUser ให้ส่งกลับ
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

// Student Repository
type StudentRepository interface {
	CreateStudent(tx transaction.Transaction, student *entities.Student) error
    EditStudentByID(student *entities.Student) error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *studentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) CreateStudent(tx transaction.Transaction, student *entities.Student) error {
	gormTx := tx.(*transaction.GormTransaction)
	if err := gormTx.GetDB().Create(student).Error; err != nil {
		return err // ถ้ามีข้อผิดพลาดใน CreateStudent ให้ส่งกลับ
	}
	return nil
}
func (r *userRepository) GetStudentByUserID(id uint) (*entities.Student, error) {
	var student entities.Student
	result := r.db.Where("user_id = ?", id).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

func (r *studentRepository) EditStudentByID(student *entities.Student) error{
    return r.db.Save(student).Error
}

// Teacher Repository
type TeacherRepository interface {
	CreateTeacher(tx transaction.Transaction, teacher *entities.Teacher) error
    EditTeacherByID(teacher *entities.Teacher) error
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
func (r *teacherRepository) EditTeacherByID(teacher *entities.Teacher) error{
    return r.db.Save(teacher).Error
}