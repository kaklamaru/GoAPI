package database

import (
	"RESTAPI/config"
	"RESTAPI/domain/entities"
	"fmt"
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
    // Migrate Faculty table first
    if err := m.Db.AutoMigrate(&entities.Faculty{}); err != nil {
        return fmt.Errorf("failed to migrate Faculty: %w", err)
    }

    // Migrate Branch table after Faculty
    if err := m.Db.AutoMigrate(&entities.Branch{}); err != nil {
        return fmt.Errorf("failed to migrate Branch: %w", err)
    }

    // Migrate other tables as needed
    if err := m.Db.AutoMigrate(&entities.User{}); err != nil {
        return fmt.Errorf("failed to migrate User: %w", err)
    }

    if err := m.Db.AutoMigrate(&entities.Teacher{}); err != nil {
        return fmt.Errorf("failed to migrate Teacher: %w", err)
    }
    if err := m.Db.AutoMigrate(&entities.Student{}); err != nil {
        return fmt.Errorf("failed to migrate Student: %w", err)
    }
    if err := m.Db.AutoMigrate(&entities.Event{}); err != nil {
        return fmt.Errorf("failed to migrate Event: %w", err)
    }
    if err := m.Db.AutoMigrate(&entities.Permission{}); err != nil {
        return fmt.Errorf("failed to migrate Permission: %w", err)
    }

    return nil
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

    addTriggerIfNotExists(db, "before_insert_students", `
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
    addTriggerIfNotExists(db, "before_insert_teachers", `
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

// ฟังก์ชันเพื่อเพิ่ม Trigger ถ้ามันยังไม่มี
func addTriggerIfNotExists(db Database, triggerName, triggerSQL string) {
    var count int
    // ตรวจสอบว่า trigger มีอยู่แล้วหรือไม่
    query := fmt.Sprintf("SHOW TRIGGERS LIKE '%s'", triggerName)
    err := db.GetDb().Raw(query).Scan(&count).Error
    if err != nil {
        log.Printf("Error checking trigger existence: %v", err)
        return
    }

    if count == 0 {
        // ถ้าไม่มี trigger นี้ก็สร้างใหม่
        if err := db.GetDb().Exec(triggerSQL).Error; err != nil {
            log.Printf("Failed to create trigger %s: %v", triggerName, err)
        } else {
            log.Printf("Trigger %s created successfully!", triggerName)
        }
    } else {
        log.Printf("Trigger %s already exists, skipping creation.", triggerName)
    }
}