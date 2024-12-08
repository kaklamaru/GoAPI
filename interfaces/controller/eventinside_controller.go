package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"RESTAPI/utility"

	"github.com/gofiber/fiber/v2"
)

type EventInsideController struct{
	usecase usecase.EventInsideUsecase
	eventusecase usecase.EventUsecase

}

func NewEventInsideController(usecase usecase.EventInsideUsecase, eventusecase usecase.EventUsecase) *EventInsideController{
	return &EventInsideController{
		usecase: usecase,
		eventusecase: eventusecase,
	}
}

func (c *EventInsideController) JoinEvent(ctx *fiber.Ctx) error{
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
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	event, err := c.eventusecase.GetEventByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	eventInside := &entities.EventInside{
		EventId: id,
		User:    userID,
		Status:  false,
		Certifier: event.Creator,

	}

	if err := c.usecase.CreateEventInside(eventInside); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Joined event successfully",
	})
}