package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"encoding/json"

	"RESTAPI/usecase"
	"strconv"
	"time"
	"RESTAPI/utility"
	"github.com/gofiber/fiber/v2"
	jwtpkg "github.com/golang-jwt/jwt/v5"
)

type EventController struct {
	usecase   usecase.EventUsecase
	txManager transaction.TransactionManager
}

func NewEventController(usecase usecase.EventUsecase, txManager transaction.TransactionManager) *EventController {
	return &EventController{
		usecase:   usecase,
		txManager: txManager,
	}
}

func (c *EventController) buildPermission(branches []uint, years []uint) (*entities.Permission, error) {
	var permission *entities.Permission

	if len(branches) > 0 && len(years) > 0 {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		permission = &entities.Permission{
			BranchIDs:      string(branchData),
			YearIDs:        string(yearData),
			AllowAllBranch: false,
			AllowAllYear:   false,
		}
	} else if len(branches) > 0 {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		permission = &entities.Permission{
			BranchIDs:      string(branchData),
			YearIDs:        "",
			AllowAllBranch: false,
			AllowAllYear:   true,
		}
	} else if len(years) > 0 {
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		permission = &entities.Permission{
			BranchIDs:      "",
			YearIDs:        string(yearData),
			AllowAllBranch: true,
			AllowAllYear:   false,
		}
	} else {
		permission = &entities.Permission{
			BranchIDs:      "",
			YearIDs:        "",
			AllowAllBranch: true,
			AllowAllYear:   true,
		}
	}

	return permission, nil
}

func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var req struct {
		EventName   string `json:"event_name"`
		StartDate   string `json:"start_date"` // This will be the date string
		WorkingHour uint   `json:"working_hour"`
		Limit       uint   `json:"limit"`
		Detail      string `json:"detail"`
		Branches    []uint `json:"branches"`
		Years       []uint `json:"years"`
	}
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	role := claims["role"].(string)
	if role != "teacher" && role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to create an event"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id in claims"})
	}
	userID := uint(userIDFloat)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return err
	}

	tx := c.txManager.Begin()

	return utility.HandleTransaction(ctx, tx, func() error {
		event := &entities.Event{
			EventName:   req.EventName,
			StartDate:   startDate,
			Limit: req.Limit,
			WorkingHour: req.WorkingHour,
			Detail:      req.Detail,
			Creator:     userID,
		}

		permission, err := c.buildPermission(req.Branches, req.Years)
		if err != nil {
			return err
		}
		
		if err := c.usecase.CreateEventWithPermission(tx, event, permission); err != nil {
			return err
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Create successfully",
			"event":   event,
		})
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

func (c *EventController) GetEventByID(ctx *fiber.Ctx) error {
	idstr := ctx.Params("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid id format",
		})
	}
	id := uint(idint)
	event, err := c.usecase.GetEventByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(event)
}

func (c *EventController) EditEvent(ctx *fiber.Ctx) error {
	eventIDStr := ctx.Params("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID format",
		})
	}

	// ดึง claims จาก JWT (เพื่อตรวจสอบ user_id)
	claims, ok := ctx.Locals("claims").(jwtpkg.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid claims",
		})
	}
	role := claims["role"]
	if role != "teacher" && role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to create an event",
		})
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	event, err := c.usecase.GetEventByID(uint(eventID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	// ตรวจสอบว่า Creator ของ Event ตรงกับ userID จาก JWT หรือไม่
	if event.Creator != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to edit this event",
		})
	}

	// Parse JSON body
	var req struct {
		EventName   string `json:"event_name"`
		StartDate   string `json:"start_date"` // ใช้ string แล้วค่อยแปลง
		WorkingHour uint   `json:"working_hour"`
		Limit       uint   `json:"limit"`
		Detail      string `json:"detail"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// แปลง StartDate จาก string เป็น time.Time
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

	// อัปเดตข้อมูล Event
	event.EventName = req.EventName
	event.StartDate = startDate
	event.WorkingHour = req.WorkingHour
	event.Detail = req.Detail

	// ส่ง pointer ของ event ไปยัง usecase
	if err := c.usecase.EditEvent(event); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to update event",
		})
	}

	// ส่งข้อมูล Event ที่แก้ไขกลับไปยังผู้ใช้
	return ctx.Status(fiber.StatusOK).JSON(event)
}
