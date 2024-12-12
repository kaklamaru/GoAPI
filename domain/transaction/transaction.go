package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

type Transaction interface {
    Commit() error
    Rollback() error
    GetDB() *gorm.DB
}

type GormTransaction struct {
    db *gorm.DB
}

func (t *GormTransaction) Commit() error {
    if t.db == nil {
        return fmt.Errorf("no active transaction")
    }
    return t.db.Commit().Error
}

func (t *GormTransaction) Rollback() error {
    if t.db == nil {
        return fmt.Errorf("no active transaction")
    }
    return t.db.Rollback().Error
}


func (t *GormTransaction) GetDB() *gorm.DB {
    return t.db
}



type TransactionManager interface {
    Begin() Transaction
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