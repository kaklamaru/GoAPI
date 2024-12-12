package controller

import (
	// "RESTAPI/domain/entities"
	// "RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"RESTAPI/usecase"

	"RESTAPI/utility"

	"github.com/gofiber/fiber/v2"
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

func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var req usecase.EventRequest

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
	}

	role, ok := claims["role"].(string)
	if !ok || (role != "teacher" && role != "admin") {
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

	if err := c.usecase.CreateEvent(&req, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

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
	var req usecase.EventRequest

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

    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request payload",
        })
    }

    if err := c.usecase.EditEvent(ctx.Context(),id, &req ,userID); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Update event success",
    })
}

func (c *EventController) DeleteEvent(ctx *fiber.Ctx) error {
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

	if err := c.usecase.DeleteEvent(id, userID); err != nil {

        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),  // ส่งข้อความ error ที่ได้จาก usecase
        })
    }

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"massage": "Event deleted successfully",
	})

}
