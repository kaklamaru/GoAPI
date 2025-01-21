package repository

import (
	"RESTAPI/domain/entities"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type OutsideRepository interface {
	CreateOutside(outside *entities.EventOutside) (uint,error)
	GetOutsideByID(id uint) (*entities.EventOutside,error)
	AllOutsideThisYears(userID uint, year uint) ([]entities.EventOutside, error) 
}

type outsideRepository struct {
	db *gorm.DB
}

func NewOutsideRepository(db *gorm.DB) OutsideRepository {
	return &outsideRepository{
		db: db,
	}
}

func (r *outsideRepository) CreateOutside(outside *entities.EventOutside) (uint,error) {
	if err := r.db.Create(outside).Error; err != nil {
		return 0,err
	}
	return outside.EventID, nil
}

func (r *outsideRepository) GetOutsideByID(id uint) (*entities.EventOutside,error){
	var outside entities.EventOutside
	if err := r.db.Preload("Student.Branch.Faculty").First(&outside, "event_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve event: %w", err)
	}
	return &outside, nil
}

func (r *outsideRepository) AllOutsideThisYears(userID uint, year uint) ([]entities.EventOutside, error) {
	var eventOutside []entities.EventOutside
	if err := r.db.Where("user = ? AND school_year = ?", userID, year).Find(&eventOutside).Error; err != nil {
		return nil, err
	}
	return eventOutside, nil
}
