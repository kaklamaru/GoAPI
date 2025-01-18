package controller

import (
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
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid JWT claims",
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

	if err := c.usecase.CreateEvent(&req, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
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


func (c *EventController) MyEvent(ctx *fiber.Ctx) error {
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid JWT claims",
		})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)
	events, err := c.usecase.MyEvent(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}
func (c *EventController) AllAllowedEvent(ctx *fiber.Ctx) error {
	events, err := c.usecase.AllAllowedEvent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}
func (c *EventController) AllCurrentEvent(ctx *fiber.Ctx) error {
	events, err := c.usecase.AllCurrentEvent()
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

    if err := c.usecase.EditEvent(id, &req ,userID); err != nil {
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
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	if err := c.usecase.DeleteEvent(id, userID); err != nil {

        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(), 
        })
    }

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"massage": "Event deleted successfully",
	})
}

func (c *EventController) StatusEvent(ctx *fiber.Ctx) error {
    id, err := utility.GetUintID(ctx)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ID",
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
    if err := c.usecase.ToggleEventStatus(id,userID); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to toggle event status",
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Event status toggled successfully",
    })
}
