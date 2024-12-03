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
        &entities.Event{},
        &entities.Permission{},
        
    )
}

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



