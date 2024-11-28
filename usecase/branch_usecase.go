package usecase

import(
	"RESTAPI/domain/entities"
	"RESTAPI/domain/repository"
)

type BranchUsecase interface{
	AddBranch(branch *entities.Branch) error
	GetAllBranches() ([]entities.Branch, error)
	GetBranch(id uint) (*entities.Branch, error)
	UpdateBranch(branch *entities.Branch) error
	GetBranchesByFaculty(id uint) ([]entities.Branch,error)
}

type branchUsecase struct{
	repo repository.BranchRepository
}

func NewBranchUsecase(repo repository.BranchRepository) BranchUsecase{
	return &branchUsecase{repo: repo}
}

func (u *branchUsecase) AddBranch(branch *entities.Branch)error{
	return u.repo.CreateBranch(branch)
}

func (u *branchUsecase) GetAllBranches() ([]entities.Branch,error){
	return u.repo.GetAllBranches()
}

func (u *branchUsecase) GetBranchesByFaculty(id uint) ([]entities.Branch,error){
	return u.repo.GetAllBranchesByFaculty(int(id))
}

func (u *branchUsecase) GetBranch(id uint) (*entities.Branch,error){
	return u.repo.GetBranch(id)
}

func (u *branchUsecase) UpdateBranch(branch *entities.Branch) error{
	return u.repo.UpdateBranch(branch)
}