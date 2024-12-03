package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"encoding/json"

	"RESTAPI/usecase"
	"strconv"
	"time"

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

func (c *EventController) handleTransaction(ctx *fiber.Ctx, tx transaction.Transaction, fn func() error) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err := tx.Commit(); err != nil {
			tx.Rollback()
		}
	}()

	if err := fn(); err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return nil
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
	claims, ok := ctx.Locals("claims").(jwtpkg.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
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
	// รับค่าจาก body
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	location, err := time.LoadLocation("Asia/Bangkok")
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load location"})
    }

    startDate, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartDate, location)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, use 'YYYY-MM-DD HH:MM:SS'"})
    }

	tx := c.txManager.Begin()

	return c.handleTransaction(ctx, tx, func() error {
		// สร้าง Event
		event := &entities.Event{
			EventName:   req.EventName,
			StartDate:   startDate,
			WorkingHour: req.WorkingHour,
			Limit:       req.Limit,
			Detail:      req.Detail,
			Creator:     userID,
		}

		// สร้าง Permission หากมี
		var permission *entities.Permission
		if len(req.Branches) > 0 || len(req.Years) > 0 {
			branchData, err := json.Marshal(req.Branches)
			if err != nil {
				return err
			}

			yearData, err := json.Marshal(req.Years)
			if err != nil {
				return err
			}

			permission = &entities.Permission{
				BranchIDs: string(branchData),
				YearIDs:   string(yearData),
			}
		}

		// สร้าง Event และ Permission (ถ้ามี)
		if err := c.usecase.CreateEventWithPermission(tx, event, permission); err != nil {
			return err
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Create successfully",
			"event":   event,
		})
	})
}

// func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
//     var req struct {
//         EventName   string   `json:"event_name"`
//         StartDate   string   `json:"start_date"`
//         WorkingHour uint     `json:"working_hour"`
//         Limit       uint     `json:"limit"`
//         Detail      string   `json:"detail"`
//         Branches    []uint   `json:"branches"`
//         Years       []uint   `json:"years"`
//     }

//     claims, ok := ctx.Locals("claims").(jwtpkg.MapClaims)
//     if !ok {
//         return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
//     }

//     role := claims["role"].(string)
//     if role != "teacher" && role != "admin" {
//         return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to create an event"})
//     }

//     userIDFloat, ok := claims["user_id"].(float64)
//     if !ok {
//         return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id in claims"})
//     }
//     userID := uint(userIDFloat)

//     if err := ctx.BodyParser(&req); err != nil {
//         return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
//     }

//     location, err := time.LoadLocation("Asia/Bangkok")
//     if err != nil {
//         return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load location"})
//     }

//     startDate, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartDate, location)
//     if err != nil {
//         return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, use 'YYYY-MM-DD HH:MM:SS'"})
//     }

//     tx := c.txManager.Begin()

//     return c.handleTransaction(ctx, tx, func() error {
//         // สร้าง Event ก่อน
//         event := &entities.Event{
//             EventName:   req.EventName,
//             StartDate:   startDate,
//             WorkingHour: req.WorkingHour,
//             Limit:       req.Limit,
//             Detail:      req.Detail,
//             Creator:     userID,
//         }

//         if err := c.usecase.CreateEvent(tx, event); err != nil {
//             return err
//         }

//         // ตรวจสอบว่า Event ถูกสร้างเสร็จและมี EventID
//         if event.EventID == 0 {
//             return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Event creation failed"})
//         }

//         // สร้าง Permission เมื่อ EventID ถูกกำหนดแล้ว
//         if len(req.Branches) > 0 || len(req.Years) > 0 {
//             branchData, err := json.Marshal(req.Branches)
//             if err != nil {
//                 return err
//             }

//             yearData, err := json.Marshal(req.Years)
//             if err != nil {
//                 return err
//             }

//             permission := &entities.Permission{
//                 EventID:   event.EventID,  // ใช้ EventID ที่ได้จากการสร้าง Event
//                 BranchIDs: string(branchData),
//                 YearIDs:   string(yearData),
//             }

//             if err := c.usecase.CreatePermission(tx, permission); err != nil {
//                 return err
//             }
//         }

//         return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
//             "message": "Create successfully",
//             "event":   event,
//         })
//     })
// }

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
	event.Limit = req.Limit
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
