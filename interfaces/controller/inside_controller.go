package controller

import (
	"RESTAPI/usecase"
	"RESTAPI/utility"
	"strconv"

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
	userID, ok := utility.GetUserIDFromClaims(claims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}

	if err := c.insideUsecase.JoinEventInside(id, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
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

func (c *EventInsideController) UploadFile(ctx *fiber.Ctx) error {
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

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get file",
		})
	}

	if err := c.insideUsecase.UploadFile(file, id, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file uploaded successfully",
	})
}

func (c *EventInsideController) GetFileForMe(ctx *fiber.Ctx) error {
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}
    // ดึง claims ของผู้ใช้
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to get user claims",
		})
	}

    // ตรวจสอบ user_id ใน claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)
    // ดึง URL ของไฟล์ที่ต้องการ
	filePath, err := c.insideUsecase.GetFile(id, userID)
    if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.SendFile(filePath)
}

func (c *EventInsideController) GetFile(ctx *fiber.Ctx) error{
	eventID, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}
	idStr := ctx.Params("userid")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UserID",
		})
	}
	userID := uint(idInt)
	filePath, err := c.insideUsecase.GetFile(eventID, userID)
    if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.SendFile(filePath)

}

func (c *EventInsideController) MyChecklist(ctx *fiber.Ctx) error{
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
	checklist,err:= c.insideUsecase.MyChecklist(userID,id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(checklist)
}


















func (c *EventInsideController) ConfirmAndCheck(ctx *fiber.Ctx) error {
	var req struct {
		EventID uint   `json:"event_id"`
		UserID  uint   `json:"user_id"`
		Status  bool   `json:"status"`
		Comment string `json:"comment"`
	}

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.insideUsecase.UpdateEventStatusAndComment(req.EventID, req.UserID, req.Status, req.Comment); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Checking successfully",
	})

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






