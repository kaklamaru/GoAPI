package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"RESTAPI/utility"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type EventInsideController struct {
	usecase      usecase.EventInsideUsecase
	eventusecase usecase.EventUsecase
	userUsecase  usecase.UserUsecase
}

func NewEventInsideController(usecase usecase.EventInsideUsecase, eventusecase usecase.EventUsecase, userUsecase usecase.UserUsecase) *EventInsideController {
	return &EventInsideController{
		usecase:      usecase,
		eventusecase: eventusecase,
		userUsecase: userUsecase,
	}
}

func checkPermission(permission *entities.EventResponse, user *entities.Student) bool {
	if permission == nil || user == nil {
		return false
	}
	fmt.Print(permission.AllowAllBranch)
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

	event, err := c.eventusecase.GetEventByID(id)
	if err != nil || event == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
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

	if err := c.usecase.CreateEventInside(eventInside); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join event",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Joined event successfully",
	})
}


