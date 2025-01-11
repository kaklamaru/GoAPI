package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"RESTAPI/utility"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type OutsideController struct {
	usecase usecase.OutsideUsecase
}

func NewOutsideController(usecase usecase.OutsideUsecase) *OutsideController {
	return &OutsideController{
		usecase: usecase,
	}
}

func (c *OutsideController) CreateOutside(ctx *fiber.Ctx) error {
	var req entities.OutsideRequest

	// ดึง JWT claims จาก context
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid JWT claims",
		})
	}

	// แปลง user_id จาก claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	// แปลง request body เป็นโครงสร้าง OutsideRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// เรียกใช้ usecase เพื่อสร้าง EventOutside และรับ eventID
	eventID, err := c.usecase.CreateOutside(req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// ส่งคืน eventID ในการตอบสนอง
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create successfully",
		"event_id": eventID,
	})
}

func (c *OutsideController) GetOutsideByID(ctx *fiber.Ctx) error{
	id ,err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid event ID",
		})
	}
	outside,err:= c.usecase.GetOutsideByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Failed to retrieve event",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(outside)
}

func (c *OutsideController) DownloadPDF(ctx *fiber.Ctx) error{
	id ,err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid event ID",
		})
	}

	// เรียกใช้ usecase เพื่อสร้าง PDF ในหน่วยความจำ
	data, fileName, err := c.usecase.CreateFile(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error: %v", err))
	}

	// กำหนดค่า HTTP headers เพื่อให้ผู้ใช้ดาวน์โหลดไฟล์
	ctx.Set("Content-Type", "application/pdf")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	return ctx.Send(data)
}