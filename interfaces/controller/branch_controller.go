package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"net/http"
	"strconv"
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
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
	}
	if err := c.usecase.AddBranch(branch) ;err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to create branch",
        })
	}
	return ctx.Status(http.StatusCreated).JSON(branch)
}

func (c *BranchController) GetAllBranches(ctx *fiber.Ctx)error{
	branches,err := c.usecase.GetAllBranches()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to retrieve branches",
        })
	}
	return ctx.Status(http.StatusOK).JSON(branches)
}

func (c *BranchController) GetBranch(ctx *fiber.Ctx) error{
	idstr := ctx.Params("id")
	idint, err := strconv.Atoi(idstr)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid id format",
        })
    }
    // แปลง facultyID จาก int เป็น uint
    id := uint(idint)
	branch,err:= c.usecase.GetBranch(id)
	if err != nil {
        return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
            "error": "Branch not found",
        })
    }
	// ส่งคืนข้อมูลคณะที่พบในรูปแบบ JSON
    return ctx.Status(http.StatusOK).JSON(branch)
}

func (c *BranchController) UpdateBranch(ctx *fiber.Ctx) error{
	branch := new(entities.Branch)
	if err:=ctx.BodyParser(branch); err != nil {
		return ctx.Status(http.StatusBadGateway).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if err:=c.usecase.UpdateBranch(branch);err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to update branch",
        })
	}
	return ctx.Status(http.StatusOK).JSON(branch)
}
