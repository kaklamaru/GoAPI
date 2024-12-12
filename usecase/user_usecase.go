package usecase

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
	"RESTAPI/domain/transaction"
	"errors"
)

type UserUsecase interface {
	RegisterUserAndStudent(tx transaction.Transaction, user *entities.User, student *entities.Student) error
	RegisterUserAndTeacher(tx transaction.Transaction, user *entities.User, teacher *entities.Teacher) error
	GetUserByEmail(email string) (*entities.User, error)
	GetStudentByUserID(userID uint) (*entities.Student, error)
	GetTeacherByUserID(userID uint) (*entities.Teacher, error)
	EditStudentByID(student *entities.Student) error
	EditTeacherByID(teacher *entities.Teacher) error
	GetAllStudent() ([]entities.Student, error)
	GetAllTeacher() ([]entities.Teacher, error)
}

type userUsecase struct {
	userRepo    repository.UserRepository
	studentRepo repository.StudentRepository
	teacherRepo repository.TeacherRepository
}

func NewUserUsecase(userRepo repository.UserRepository, studentRepo repository.StudentRepository, teacherRepo repository.TeacherRepository) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		studentRepo: studentRepo,
		teacherRepo: teacherRepo,
	}
}

func (u *userUsecase) registerUserAndEntity(tx transaction.Transaction,user *entities.User, entity interface{}, createEntityFunc func(tx transaction.Transaction, entity interface{}) error) error {
	// สร้าง User
	if err := u.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return err
	}

	// ตรวจสอบว่า user.UserID ได้รับการอัปเดต
	if user.UserID == 0 {
		tx.Rollback()
		return errors.New("failed to get valid UserID")
	}

	// เชื่อมโยง Entity กับ User
	switch e := entity.(type) {
	case *entities.Student:
		e.UserID = user.UserID
	case *entities.Teacher:
		e.UserID = user.UserID
	default:
		tx.Rollback()
		return errors.New("unsupported entity type")
	}

	// สร้าง Entity
	if err := createEntityFunc(tx, entity); err != nil {
		tx.Rollback()
		return err
	}

	// Commit ธุรกรรม
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) RegisterUserAndStudent(tx transaction.Transaction,user *entities.User, student *entities.Student) error {
	return u.registerUserAndEntity(tx, user, student, func(tx transaction.Transaction, entity interface{}) error {
		return u.studentRepo.CreateStudent(tx, entity.(*entities.Student))
	})
}

func (u *userUsecase) RegisterUserAndTeacher(tx transaction.Transaction, user *entities.User, teacher *entities.Teacher) error {
	return u.registerUserAndEntity(tx, user, teacher, func(tx transaction.Transaction, entity interface{}) error {
		return u.teacherRepo.CreateTeacher(tx, entity.(*entities.Teacher))
	})
}

func (u *userUsecase) GetUserByEmail(email string) (*entities.User, error) {
	return u.userRepo.GetUserByEmail(email)
}

func (u *userUsecase) GetStudentByUserID(userID uint) (*entities.Student, error) {
	return u.userRepo.GetStudentByUserID(userID)
}

func (u *userUsecase) GetTeacherByUserID(userID uint) (*entities.Teacher, error) {
	return u.userRepo.GetTeacherByUserID(userID)
}

func (u *userUsecase) EditStudentByID(student *entities.Student) error {
	return u.studentRepo.EditStudentByID(student)
}

func (u *userUsecase) EditTeacherByID(teacher *entities.Teacher) error {
	return u.teacherRepo.EditTeacherByID(teacher)
}

func (u *userUsecase) GetAllStudent() ([]entities.Student, error) {
	return u.studentRepo.GetAllStudent()
}

func (u *userUsecase) GetAllTeacher() ([]entities.Teacher, error) {
	return u.teacherRepo.GetAllTeacher()
}
