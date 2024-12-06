package repository

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"encoding/json"

	// "encoding/json"
	"fmt"

	"gorm.io/gorm"
)



type PermissionRepository interface {
	CreatePermission(tx transaction.Transaction, permission *entities.Permission) error
	DeletePermission(eventID uint) error
}
type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) CreatePermission(tx transaction.Transaction, permission *entities.Permission) error {
    gormTx := tx.(*transaction.GormTransaction)

    // แปลง BranchIDs และ YearIDs เป็น JSON string ก่อนบันทึก
    branchIDsJSON, err := json.Marshal(permission.BranchIDs)
    if err != nil {
        return fmt.Errorf("failed to marshal branch IDs: %w", err)
    }

    yearIDsJSON, err := json.Marshal(permission.Years)
    if err != nil {
        return fmt.Errorf("failed to marshal year IDs: %w", err)
    }

    // เก็บข้อมูล JSON ลงในฐานข้อมูล
    permission.BranchIDs = string(branchIDsJSON) // เก็บเป็น JSON string
    permission.Years = string(yearIDsJSON)       // เก็บเป็น JSON string

    // สร้าง Permission ใหม่ในฐานข้อมูล
    if err := gormTx.GetDB().Create(permission).Error; err != nil {
        return fmt.Errorf("failed to create permission: %w", err)
    }

    return nil
}


func (r *permissionRepository) DeletePermission(eventID uint) error {
	return r.db.Where("event_id = ?", eventID).Delete(&entities.Permission{}).Error
}