package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"RESTAPI/utility"
	// "fmt"
	// "time"

	"github.com/gofiber/fiber/v2"
)

type EventInsideController struct {
	insideUsecase usecase.EventInsideUsecase
	eventUsecase  usecase.EventUsecase
	userUsecase   usecase.UserUsecase
}

func NewEventInsideController(insideUsecase usecase.EventInsideUsecase, 
	eventUsecase usecase.EventUsecase, 
	userUsecase usecase.UserUsecase) *EventInsideController {
	return &EventInsideController{
		insideUsecase: insideUsecase,
		eventUsecase:  eventUsecase,
		userUsecase:   userUsecase,
	}
}

func checkPermission(permission *entities.EventResponse, user *entities.Student) bool {
	if permission == nil || user == nil {
		return false
	}
	permissionBranch := permission.AllowAllBranch
	permissionYear := permission.AllowAllYear

	if !permissionBranch && permission.BranchIDs != nil {
		for _, branch := range permission.BranchIDs {
			if user.BranchId == branch {
				permissionBranch = true
				break
			}
		}
	}

	if !permissionYear && permission.Years != nil {
		for _, year := range permission.Years {
			if user.Year == year {
				permissionYear = true
				break
			}
		}
	}

	return permissionBranch && permissionYear
}

func (c *EventInsideController) JoinEvent(ctx *fiber.Ctx) error {
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to get user claims",
		})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	student, err := c.userUsecase.GetStudentByUserID(userID)
	if err != nil || student == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	event, err := c.eventUsecase.GetEventByID(id)
	if err != nil || event == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	if event.FreeSpace == 0 {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "The event is full. Cannot join.",
		})
	}

	if !checkPermission(event, student) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Not allowed to join this event",
		})
	}

	eventInside := &entities.EventInside{
		EventId:   id,
		User:      userID,
		Status:    false,
		Certifier: event.Creator,
	}

	if err := c.insideUsecase.JoinEventInside(eventInside); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join event",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Joined event successfully",
	})
}

func (c *EventInsideController) UnJoinEventInside(ctx *fiber.Ctx) error {
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to get user claims",
		})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	if err := c.insideUsecase.UnJoinEventInside(id, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Faild to Unjoin event",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "UnJoined event successfully",
	})
}

func (c *EventInsideController) ConfirmAndCheck(ctx *fiber.Ctx) error {
	var req struct {
		EventID uint   `json:"event_id"`
		UserID  uint   `json:"user_id"`
		Status  bool   `json:"status"`
		Comment string `json:"comment"`
	}

	return ctx.Status(fiber.StatusOK).JSON(req)

}

func (c *EventInsideController) CountEventInside(ctx *fiber.Ctx) error {
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	count, err := c.insideUsecase.CountEventInside(id)
	if err != nil {
		return err
	}

	event, err := c.eventUsecase.GetEventByID(id)
	if err != nil {
		return err
	}
	freeSpace := event.Limit - count

	return ctx.Status(fiber.StatusOK).JSON(freeSpace)
}
