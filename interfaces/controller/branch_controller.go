package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"RESTAPI/utility"

	"github.com/gofiber/fiber/v2"
)
type BranchController struct{
	usecase usecase.BranchUsecase
}

func NewBranchController(usecase usecase.BranchUsecase) *BranchController{
	return &BranchController{usecase: usecase}
}

func (c *BranchController) AddBranch(ctx *fiber.Ctx)error{
	branch := new(entities.Branch)
	if err :=ctx.BodyParser(branch) ;err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
	}
	if err := c.usecase.AddBranch(branch) ;err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to create branch",
        })
	}
	return ctx.Status(fiber.StatusCreated).JSON(branch)
}

func (c *BranchController) GetAllBranches(ctx *fiber.Ctx)error{
	branches,err := c.usecase.GetAllBranches()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to retrieve branches",
        })
	}
	return ctx.Status(fiber.StatusOK).JSON(branches)
}

func (c *BranchController) GetBranch(ctx *fiber.Ctx) error{
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	branch,err:= c.usecase.GetBranch(id)
	if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Branch not found",
        })
    }
	// ส่งคืนข้อมูลคณะที่พบในรูปแบบ JSON
    return ctx.Status(fiber.StatusOK).JSON(branch)
}
func (c *BranchController) GetBranchesByFaculty(ctx *fiber.Ctx) error{
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	branch,err:= c.usecase.GetBranchesByFaculty(id)
	if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Branch not found",
        })
    }
	// ส่งคืนข้อมูลคณะที่พบในรูปแบบ JSON
    return ctx.Status(fiber.StatusOK).JSON(branch)
}

func (c *BranchController) UpdateBranch(ctx *fiber.Ctx) error{
	branch := new(entities.Branch)
    id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

    branch.BranchID=id

	if err:=ctx.BodyParser(branch); err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if err:=c.usecase.UpdateBranch(branch);err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to update branch",
        })
	}
	return ctx.Status(fiber.StatusOK).JSON(branch)
}

func (c *BranchController) DeleteBranchByID(ctx *fiber.Ctx) error {
    id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

    // เรียกใช้ usecase เพื่อลบ branch ตาม ID
    branch, err := c.usecase.DeleteBranchByID(id)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Branch deleted successfully",
        "branch": branch, 
    })
}




