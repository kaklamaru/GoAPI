package utility

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

func GetUserIDFromClaims(claims map[string]interface{}) (uint, bool) {
	userIDFloat, ok := claims["user_id"].(float64) // JWT claims มักจะใช้ float64
	if !ok {
		return 0, false
	}
	return uint(userIDFloat), true
}

func GetUintID(ctx *fiber.Ctx) (uint, error) {
	idStr := ctx.Params("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id format")
	}
	return uint(idInt), nil
}

func DecodeIDs(dataStr string) ([]uint, error) {
	var ids []uint

	// เช็คว่าเป็น string ว่างหรือไม่ (หลังจาก Trim เครื่องหมาย escape ออก)
	if strings.Trim(dataStr, "\"") == "" {
		// // ถ้าเป็น string ว่าง ให้ return slice ว่าง
		// fmt.Println("IDs: []") // กรณีที่เป็น string ว่าง
		return ids, nil
	}

	// ถ้า string ไม่ว่างให้ Trim เครื่องหมายคำพูดที่ escape ออก
	dataStr = strings.Trim(dataStr, "\"")

	// แปลง string ที่เหลือเป็น []uint
	if err := json.Unmarshal([]byte(dataStr), &ids); err != nil {
		// ถ้าเกิดข้อผิดพลาดในการ decode
		fmt.Println("Error decoding IDs:", err)
		return nil, err
	}

	return ids, nil
}
