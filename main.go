package main

import (
	"RESTAPI/config"
	"RESTAPI/infrastructure/database"
	"RESTAPI/infrastructure/jwt"
	"RESTAPI/interfaces/server"
	"log"
	"time"
)

func deleteOldNews(db database.Database) {
	for {
		// ตั้งค่าเขตเวลาให้ตรงกับข้อมูลที่จัดเก็บ
		db.GetDb().Exec("SET time_zone = '+07:00'")
		
		result := db.GetDb().Exec("DELETE FROM `news` WHERE created_at < NOW() - INTERVAL 1 DAY")
		if result.Error != nil {
			log.Printf("Error deleting old news: %v", result.Error)
		} else {
			log.Printf("Deleted %d old news data successfully.", result.RowsAffected)
		}

		time.Sleep(8 * time.Hour)
	}
}


func main() {
	// โหลดค่าคอนฟิกจากไฟล์ .env
	cfg := config.LoadConfig()

	// ตั้งค่าและเชื่อมต่อฐานข้อมูล
	db := database.SetupDatabase(cfg)

	// เริ่มรัน goroutine สำหรับการลบข้อมูลเก่า
	go deleteOldNews(db)

	// สร้าง instance ของ JWT service
	jwtService := jwt.NewJWTService(cfg)

	// สร้าง instance ของ server
	srv, err := server.NewServer(cfg, db, jwtService)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// เริ่มต้น server
	log.Printf("Starting server on port %d...", cfg.ServerPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
