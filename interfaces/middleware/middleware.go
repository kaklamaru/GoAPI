package middleware

import (
	"RESTAPI/infrastructure/jwt"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddlewareFromCookie(jwtService *jwtimpl.JWTService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Cookies("token") // ดึง JWT จาก Cookie
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		claims, err := jwtService.ValidateJWT(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// เก็บ claims ลงใน Locals สำหรับใช้งานใน Controller
		ctx.Locals("claims", claims)
		return ctx.Next()
	}
}
