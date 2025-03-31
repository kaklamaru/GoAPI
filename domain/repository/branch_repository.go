package repository

import (
	"RESTAPI/domain/entities"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type BranchRepository interface {
	CreateBranch(branch *entities.Branch) error
	GetAllBranches() ([]entities.Branch, error)
	GetBranch(id uint) (*entities.Branch, error)
	UpdateBranch(branch *entities.Branch) error
	GetAllBranchesByFaculty(facultyId int) ([]entities.Branch, error)
	DeleteBranchByID(id uint) (*entities.Branch, error)
	BranchExists(id uint) (bool, error)
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
	if err := r.db.Preload("Faculty").Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}


func (r *branchRepository) GetBranch(id uint) (*entities.Branch, error) {
    var branch entities.Branch
    // ดึงข้อมูล Branch พร้อมกับ Faculty ที่เกี่ยวข้อง
    if err := r.db.Preload("Faculty").First(&branch, id).Error; err != nil {
        return nil, err
    }

    return &branch, nil
}

func (r *branchRepository) BranchExists(branchID uint) (bool, error) {
    var branch entities.Branch
    if err := r.db.Where("branch_id = ?", branchID).First(&branch).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return false, nil 
        }
        return false, err 
    }
    return true, nil
}




func (r *branchRepository) UpdateBranch(branch *entities.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) GetAllBranchesByFaculty(facultyId int) ([]entities.Branch, error) {
	var branches []entities.Branch
	result := r.db.Where("faculty_id = ?", facultyId).Find(&branches)
	if result.Error != nil {
		return nil, result.Error
	}
	return branches, nil
}

func (r *branchRepository) DeleteBranchByID(id uint) (*entities.Branch, error) {
	branch := &entities.Branch{}

	if err := r.db.First(branch, "branch_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("branch with ID %d not found", id)
		}
		return nil, err
	}

	// ลบ branch
	if err := r.db.Delete(branch).Error; err != nil {
		return nil, err
	}

	return branch, nil
}
