package database

import (
    "RESTAPI/config"
    "RESTAPI/domain/entities"
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// Interface สำหรับการเชื่อมต่อฐานข้อมูล
type Database interface {
    GetDb() *gorm.DB
    AutoMigrate() error
}

type mysqlDatabase struct {
    Db *gorm.DB
}

// ฟังก์ชันสำหรับสร้าง instance ของฐานข้อมูล
func NewMySQLDatabase(cfg *config.Config) (Database, error) {
    db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &mysqlDatabase{Db: db}, nil
}

// ฟังก์ชันสำหรับดึง instance ของฐานข้อมูล
func (m *mysqlDatabase) GetDb() *gorm.DB {
    return m.Db
}

// ฟังก์ชันสำหรับทำ AutoMigrate
func (m *mysqlDatabase) AutoMigrate() error {
    return m.Db.AutoMigrate(
        &entities.User{},
        &entities.Faculty{},
        &entities.Branch{},
        &entities.Student{},
        &entities.Teacher{},
        
    )
}

// // ฟังก์ชันสำหรับตั้งค่าและเชื่อมต่อฐานข้อมูล
// func SetupDatabase(cfg *config.Config) Database {
//     db, err := NewMySQLDatabase(cfg)
//     if err != nil {
//         log.Fatalf("failed to connect to database: %v", err)
//     }

//     err = db.AutoMigrate()
//     if err != nil {
//         log.Fatalf("failed to migrate database: %v", err)
//     }

//     log.Println("Database connected and migrated successfully!")
//     return db
// }

// MYSQL
func SetupDatabase(cfg *config.Config) Database {
    db, err := NewMySQLDatabase(cfg)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    err = db.AutoMigrate()
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    // เพิ่ม Trigger สำหรับ Students
    db.GetDb().Exec(`
        CREATE TRIGGER before_insert_students
        BEFORE INSERT ON students
        FOR EACH ROW
        BEGIN
            IF EXISTS (SELECT 1 FROM teachers WHERE user_id = NEW.user_id) THEN
                SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'User ID already exists in teachers';
            END IF;
        END;
    `)

    // เพิ่ม Trigger สำหรับ Teachers
    db.GetDb().Exec(`
        CREATE TRIGGER before_insert_teachers
        BEFORE INSERT ON teachers
        FOR EACH ROW
        BEGIN
            IF EXISTS (SELECT 1 FROM students WHERE user_id = NEW.user_id) THEN
                SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'User ID already exists in students';
            END IF;
        END;
    `)

    log.Println("Database connected, migrated, and triggers added successfully!")
    return db
}

// POSTGRESQL
// func SetupDatabase(cfg *config.Config) Database {
//     db, err := NewMySQLDatabase(cfg)
//     if err != nil {
//         log.Fatalf("failed to connect to database: %v", err)
//     }

//     err = db.AutoMigrate()
//     if err != nil {
//         log.Fatalf("failed to migrate database: %v", err)
//     }

//     // เพิ่ม Constraint เพื่อป้องกัน UserID ซ้ำระหว่าง Students และ Teachers
//     if err := db.GetDb().Exec(`
//         ALTER TABLE students
//         ADD CONSTRAINT check_user_id_not_in_teacher
//         CHECK (user_id NOT IN (SELECT user_id FROM teachers))
//     `).Error; err != nil {
//         log.Printf("failed to add constraint to students: %v", err)
//     }

//     if err := db.GetDb().Exec(`
//         ALTER TABLE teachers
//         ADD CONSTRAINT check_user_id_not_in_student
//         CHECK (user_id NOT IN (SELECT user_id FROM students))
//     `).Error; err != nil {
//         log.Printf("failed to add constraint to teachers: %v", err)
//     }

//     log.Println("Database connected and migrated successfully!")
//     return db
// }

