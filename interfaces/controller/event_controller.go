package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"RESTAPI/usecase"
	"encoding/json"
	"fmt"
	// "time"

	"RESTAPI/utility"
	// "time"

	"github.com/gofiber/fiber/v2"
)

type EventController struct {
	usecase   usecase.EventUsecase
	txManager transaction.TransactionManager
}

type Permission struct {
	BranchIDs      string `json:"branches"`
	Years          string `json:"years"`
	AllowAllBranch bool   `json:"allow_all_branch"`
	AllowAllYear   bool   `json:"allow_all_year"`
}

func NewEventController(usecase usecase.EventUsecase, txManager transaction.TransactionManager) *EventController {
	return &EventController{
		usecase:   usecase,
		txManager: txManager,
	}
}

func (c *EventController) buildPermission(branches []uint, years []uint) (*Permission, error) {
	var permission *Permission

	if len(branches) > 0 {
		for _, branchID := range branches {
			exists, err := c.usecase.CheckBranch(branchID)
			if err != nil {
				return nil, fmt.Errorf("error checking branch: %v", err)
			}
			if !exists {
				return nil, fmt.Errorf("branch with ID %d does not exist", branchID)
			}
		}
	}

	if len(branches) > 0 && len(years) > 0 {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		permission = &Permission{
			BranchIDs:      string(branchData),
			Years:          string(yearData),
			AllowAllBranch: false,
			AllowAllYear:   false,
		}
	} else if len(branches) > 0 {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		permission = &Permission{
			BranchIDs:      string(branchData),
			Years:          "",
			AllowAllBranch: false,
			AllowAllYear:   true,
		}
	} else if len(years) > 0 {
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		permission = &Permission{
			BranchIDs:      "",
			Years:          string(yearData),
			AllowAllBranch: true,
			AllowAllYear:   false,
		}
	} else {
		permission = &Permission{
			BranchIDs:      "",
			Years:          "",
			AllowAllBranch: true,
			AllowAllYear:   true,
		}
	}

	return permission, nil
}


func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var req struct {
		EventName   string `json:"event_name"`
		StartDate   string `json:"start_date"`
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

	permission, err := c.buildPermission(req.Branches, req.Years)
	if err != nil {
		return err
	}

	event := &entities.Event{
		EventName:      req.EventName,
		StartDate:      startDate,
		Limit:          req.Limit,
		WorkingHour:    req.WorkingHour,
		Detail:         req.Detail,
		Creator:        userID,
		BranchIDs:      permission.BranchIDs,
		Years:          permission.Years,
		AllowAllBranch: permission.AllowAllBranch,
		AllowAllYear:   permission.AllowAllYear,
	}

	if err := c.usecase.CreateEvent(event); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create successfully ",
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
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}
	event, err := c.usecase.GetEventByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve event",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(event)
}

func (c *EventController) EditEvent(ctx *fiber.Ctx) error {
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
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

	event, err := c.usecase.GetEventByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	// ตรวจสอบว่าผู้ใช้มีสิทธิ์ในการแก้ไข event นี้
	if event.Creator != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to edit this event",
		})
	}

	var req struct {
		EventName   string `json:"event_name"`
		StartDate   string `json:"start_date"`
		WorkingHour uint   `json:"working_hour"`
		Limit       uint   `json:"limit"`
		Detail      string `json:"detail"`
		Branches    []uint `json:"branches"`
		Years       []uint `json:"years"`
	}

	// ตรวจสอบ payload ที่ส่งมา
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return err
	}

	permission, err := c.buildPermission(req.Branches, req.Years)
	if err != nil {
		return err
	}
	var eventEdit entities.Event

	eventEdit.EventID = event.EventID
	eventEdit.EventName = req.EventName
	eventEdit.StartDate = startDate
	eventEdit.Limit = req.Limit
	eventEdit.WorkingHour = req.WorkingHour
	eventEdit.Detail = req.Detail
	eventEdit.BranchIDs = permission.BranchIDs
	eventEdit.Years = permission.Years
	eventEdit.AllowAllBranch = permission.AllowAllBranch
	eventEdit.AllowAllYear = permission.AllowAllYear

	// เรียกใช้ usecase เพื่ออัปเดต event
	if err := c.usecase.EditEvent(&eventEdit); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to update event",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update event success",
	})
}

func (c *EventController) DeleteEvent(ctx *fiber.Ctx) error{
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
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

	event, err := c.usecase.GetEventByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	if event.Creator != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to edit this event",
		})
	}

	if err:= c.usecase.DeleteEvent(id);err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"massage":"Event deleted successfully",
	})

}