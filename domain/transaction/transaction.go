package transaction

import "gorm.io/gorm"

// Transaction interface สำหรับจัดการธุรกรรม
type Transaction interface {
    Commit() error
    Rollback() error
}

// TransactionManager interface สำหรับเริ่มต้นธุรกรรม
type TransactionManager interface {
    Begin() Transaction
}

type GormTransaction struct {
    db *gorm.DB
}

func (t *GormTransaction) Commit() error {
    return t.db.Commit().Error
}

func (t *GormTransaction) Rollback() error {
    return t.db.Rollback().Error
}

// เพิ่มฟังก์ชัน Getter เพื่อเข้าถึง DB
func (t *GormTransaction) GetDB() *gorm.DB {
    return t.db
}

type GormTransactionManager struct {
    db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) *GormTransactionManager {
    return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) Begin() Transaction {
    return &GormTransaction{db: tm.db.Begin()}
}
