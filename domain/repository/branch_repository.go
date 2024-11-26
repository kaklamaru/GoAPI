package repository

import (
	"RESTAPI/domain/entities"
	"gorm.io/gorm"
)

type BranchRepository interface {
	CreateBranch(branch *entities.Branch) error
	GetAllBranches() ([]entities.Branch, error)
	GetBranch(id uint) (*entities.Branch, error)
	UpdateBranch(branch *entities.Branch) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) CreateBranch(branch *entities.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) GetAllBranches() ([]entities.Branch, error) {
	var branches []entities.Branch
	if err := r.db.Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

func (r *branchRepository) GetBranch(id uint) (*entities.Branch, error) {
	var branch entities.Branch
	if err := r.db.First(&branch, id).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) UpdateBranch(branch *entities.Branch) error {
	return r.db.Save(branch).Error
}
