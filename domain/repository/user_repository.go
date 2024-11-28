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
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(tx transaction.Transaction, user *entities.User) error {
	gormTx := tx.(*transaction.GormTransaction)
	if err := gormTx.GetDB().Create(user).Error; err != nil {
		return err // ถ้ามีข้อผิดพลาดใน CreateUser ให้ส่งกลับ
	}
	return nil
}

func (r *GormUserRepository) GetUserByEmail(email string) (*entities.User, error) {
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
}

type GormStudentRepository struct {
    db *gorm.DB
}

func NewGormStudentRepository(db *gorm.DB) *GormStudentRepository {
    return &GormStudentRepository{db: db}
}

func (r *GormStudentRepository) CreateStudent(tx transaction.Transaction, student *entities.Student) error {
    gormTx := tx.(*transaction.GormTransaction)
    if err := gormTx.GetDB().Create(student).Error; err != nil {
        return err // ถ้ามีข้อผิดพลาดใน CreateStudent ให้ส่งกลับ
    }
    return nil
}
func (r *GormUserRepository) GetStudentByUserID(userID uint) (*entities.Student, error) {
    var student entities.Student
    result := r.db.Where("user_id = ?", userID).First(&student)
    if result.Error != nil {
        return nil, result.Error
    }
    return &student, nil
}


// Teacher Repository
type TeacherRepository interface {
    CreateTeacher(tx transaction.Transaction, teacher *entities.Teacher) error
}

type GormTeacherRepository struct {
    db *gorm.DB
}

func NewGormTeacherRepository(db *gorm.DB) *GormTeacherRepository {
    return &GormTeacherRepository{db: db}
}

func (r *GormTeacherRepository) CreateTeacher(tx transaction.Transaction, teacher *entities.Teacher) error {
    gormTx := tx.(*transaction.GormTransaction)
    return gormTx.GetDB().Create(teacher).Error
}
