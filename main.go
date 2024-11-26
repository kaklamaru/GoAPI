package main

import (
	"RESTAPI/config"
	"RESTAPI/infrastructure/database"
	"RESTAPI/infrastructure/jwt"
	"RESTAPI/interfaces/server"
	"log"

	// "github.com/golang-jwt/jwt/v5"
)

func main() {
	// โหลดค่าคอนฟิกจากไฟล์ .env
	cfg := config.LoadConfig()

	// ตั้งค่าและเชื่อมต่อฐานข้อมูล
	db := database.SetupDatabase(cfg)
	
	jwtService := jwtimpl.NewJWTService(cfg)
	
	// สร้างและเริ่มต้นเซิร์ฟเวอร์ Fiber โดยใช้ค่าคอนฟิกและฐานข้อมูล
	// สร้าง instance ของ server
	srv, err := server.NewServer(cfg, db ,jwtService)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	
	// เริ่มต้น server
	log.Printf("Starting server on port %d...", cfg.ServerPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
