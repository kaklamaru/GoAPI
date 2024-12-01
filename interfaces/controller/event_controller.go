package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtpkg "github.com/golang-jwt/jwt/v5"
)

type EventController struct {
	usecase usecase.EventUsecase
}

func NewEventController(usecase usecase.EventUsecase) *EventController {
	return &EventController{usecase: usecase}
}

func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var req struct {
		EventName   string `json:"event_name"`
		StartDate   string `json:"start_date"` // ใช้ string แล้วค่อยแปลง
		WorkingHour uint   `json:"working_hour"`
		Limit       uint   `json:"limit"`
		Detail      string `json:"detail"`
	}

	// ดึง claims จาก context
	claims, ok := ctx.Locals("claims").(jwtpkg.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid claims",
		})
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	// Parse JSON body
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// แปลง start_date จาก string เป็น time.Time และกำหนด TimeZone
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load location",
		})
	}

	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartDate, location)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format, use 'YYYY-MM-DD HH:MM:SS'",
		})
	}

	// สร้าง event
	event := &entities.Event{
		EventName:   req.EventName,
		StartDate:   startDate,
		WorkingHour: req.WorkingHour,
		Limit:       req.Limit,
		Detail:      req.Detail,
		Creator:     userID,
	}

	// เรียก usecase เพื่อสร้าง event
	if err := c.usecase.CreateEvent(event); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to create event",
		})
	}

	// ส่ง response
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create successfully",
	})
}


func (c *EventController) GetAllEvent(ctx *fiber.Ctx) error {
	events, err := c.usecase.GetAllEvent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}
