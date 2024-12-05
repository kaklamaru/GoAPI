package utility

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtservice "github.com/golang-jwt/jwt/v5"
)

func ParseStartDate(dateStr string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to load location")
	}
	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, location)
	if err != nil {
		return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid date format, use 'YYYY-MM-DD HH:MM:SS'")
	}
	return startDate, nil
}

func GetClaimsFromContext(ctx *fiber.Ctx) (jwtservice.MapClaims, error) {
	claims, ok := ctx.Locals("claims").(jwtservice.MapClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid claims")
	}
	return claims, nil
}