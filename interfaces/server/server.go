package server

import (
	"RESTAPI/config"
	"RESTAPI/infrastructure/database"
	"RESTAPI/infrastructure/jwt"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server interface สำหรับการเริ่มต้นเซิร์ฟเวอร์
type Server interface {
	Start() error
}

// โครงสร้างสำหรับเซิร์ฟเวอร์ Fiber
type fiberServer struct {
	app  *fiber.App
	port int
}

// NewServer ฟังก์ชันสำหรับสร้าง instance ของเซิร์ฟเวอร์ Fiber
func NewServer(cfg *config.Config, db database.Database ,jwtService *jwt.JWTService) (Server, error) {
	// ตรวจสอบค่าพอร์ต
	if cfg.ServerPort == 0 {
		return nil, fmt.Errorf("Server port not specified in config")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000/, https://your-frontend-domain.com/", // อนุญาตเฉพาะ origin ที่ระบุ
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",            // Headers ที่อนุญาต
        AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",                        // Methods ที่อนุญาต
        AllowCredentials: true,                                                // รองรับ cookies
    }))

	// กำหนด middleware สำหรับการกู้คืนจาก panic
	app.Use(recover.New())

	// กำหนด middleware สำหรับการบันทึก log ของการร้องขอ
	app.Use(logger.New())

	// กำหนดเส้นทางทั้งหมดและส่งผ่านฐานข้อมูล
	SetupRoutes(app, db,jwtService)

	return &fiberServer{
		app:  app,
		port: cfg.ServerPort,
	}, nil
}

// Start method สำหรับเริ่มต้นเซิร์ฟเวอร์
func (s *fiberServer) Start() error {
	// สร้าง URL ของเซิร์ฟเวอร์จากพอร์ตที่กำหนด
	serverUrl := fmt.Sprintf(":%d", s.port)

	// เริ่มต้นเซิร์ฟเวอร์และบันทึกข้อผิดพลาดหากมี
	return s.app.Listen(serverUrl)
}
