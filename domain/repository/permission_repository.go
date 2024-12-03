package repository

import (
	"RESTAPI/domain/entities"
	"fmt"
	"gorm.io/gorm"
)

type PermissionRepository interface{

}
type permissionRepository struct{
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository{
	return &permissionRepository{db: db}
}

func (r *permissionRepository) CreatePermission(permission *entities.Permission) error{
	// ตรวจสอบก่อนว่า EventID และ BranchID นี้มีในฐานข้อมูลหรือไม่
	var existingPermission entities.Permission
	if err := r.db.Where("event_id = ? AND branch_id = ?", permission.EventID, permission.BranchID).First(&existingPermission).Error; err == nil {
		// ถ้ามีอยู่แล้ว จะ return error
		return fmt.Errorf("permission for this event and branch already exists")
	}
	// ถ้าไม่มี ให้ทำการสร้างใหม่
	return r.db.Create(permission).Error
}
