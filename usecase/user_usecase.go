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
func (u *userUsecase) GetStudentByUserID(userID uint) (*entities.Student, error){
    return u.userRepo.GetStudentByUserID(userID)
}

// รับ tx เป็นพารามิเตอร์ และส่ง tx ไปให้กับ repository ในแต่ละฟังก์ชัน
func (u *userUsecase) RegisterUserAndStudent(tx transaction.Transaction, user *entities.User, student *entities.Student) error {
    // สร้าง User ก่อน
    if err := u.userRepo.CreateUser(tx, user); err != nil {
        tx.Rollback()
        return err
    }

    // ตรวจสอบว่า user.UserID ได้รับการอัปเดตแล้ว
    if user.UserID == 0 {
        tx.Rollback()
        return errors.New("failed to get valid UserID")
    }

    // ตอนนี้สามารถใช้ user.UserID เพื่อสร้าง Student ได้
    student.UserID = user.UserID

    // สร้าง Student
    if err := u.studentRepo.CreateStudent(tx, student); err != nil {
        tx.Rollback()
        return err
    }

    // Commit ธุรกรรม
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

func (u *userUsecase) RegisterUserAndTeacher(tx transaction.Transaction, user *entities.User, teacher *entities.Teacher) error {
    // สร้าง user
    if err := u.userRepo.CreateUser(tx, user); err != nil {
        tx.Rollback()
        return err
    }
	  // ตรวจสอบว่า user.UserID ได้รับการอัปเดตแล้ว
	  if user.UserID == 0 {
        tx.Rollback()
        return errors.New("failed to get valid UserID")
    }
    // เชื่อมโยง teacher กับ user
    teacher.UserID = user.UserID

    // สร้าง teacher
    if err := u.teacherRepo.CreateTeacher(tx, teacher); err != nil {
        tx.Rollback()
        return err
    }

    // ยืนยันธุรกรรม
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

func (u *userUsecase) GetUserByEmail(email string) (*entities.User, error){
    return u.userRepo.GetUserByEmail(email)
}


