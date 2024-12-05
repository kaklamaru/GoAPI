package utility

import (
	"github.com/gofiber/fiber/v2"
	"RESTAPI/domain/transaction"
)

// HandleTransaction ทำการเริ่มต้น transaction และจัดการ commit/rollback
func HandleTransaction(ctx *fiber.Ctx, tx transaction.Transaction, fn func() error) error {
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
