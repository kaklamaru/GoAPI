package middleware

import (
	"RESTAPI/infrastructure/jwt"


	"github.com/gofiber/fiber/v2"
	jwtPkg "github.com/golang-jwt/jwt/v5"
)

// JWTMiddlewareFromCookie ตรวจสอบ JWT จาก Cookie
func JWTMiddlewareFromCookie(jwtService *jwt.JWTService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		tokenString := ctx.Cookies("token")
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

		ctx.Locals("claims", claims)
		return ctx.Next()
	}
}
// Middleware ตรวจสอบหลายๆ role ได้
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        claims, ok := ctx.Locals("claims").(jwtPkg.MapClaims) // ใช้ jwt.MapClaims จากที่ import ใน jwtService
        if !ok {
            return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Unauthorized access",
            })
        }

        role, ok := claims["role"].(string)
        if !ok {
            return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Role not found or invalid type",
            })
        }

        for _, allowedRole := range allowedRoles {
            if role == allowedRole {
                return ctx.Next()
            }
        }

        return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Access denied for role: " + role,
        })
    }
}





