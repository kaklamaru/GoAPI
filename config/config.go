package config

import (
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

// Config struct เก็บค่าคอนฟิกต่างๆ
type Config struct {
    DSN        string // Data Source Name สำหรับการเชื่อมต่อฐานข้อมูล
    JWTSecret  string // Secret key สำหรับ JWT
    ServerPort int    // พอร์ตของเซิร์ฟเวอร์
}

// LoadConfig ฟังก์ชันสำหรับโหลดค่าคอนฟิกจากไฟล์ .env
func LoadConfig() *Config {
    // โหลดค่าคอนฟิกจากไฟล์ .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // สร้าง DSN สำหรับการเชื่อมต่อฐานข้อมูล MySQL
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_USER"),      // ดึงค่า DB_USER จากไฟล์ .env
        os.Getenv("DB_PASSWORD"),  // ดึงค่า DB_PASSWORD จากไฟล์ .env
        os.Getenv("DB_HOST"),      // ดึงค่า DB_HOST จากไฟล์ .env
        os.Getenv("DB_PORT"),      // ดึงค่า DB_PORT จากไฟล์ .env
        os.Getenv("DB_NAME"),      // ดึงค่า DB_NAME จากไฟล์ .env
    )

    // ดึงค่า JWT_SECRET จากไฟล์ .env
    jwtSecret := os.Getenv("JWT_SECRET")

    // ดึงค่า SERVER_PORT จากไฟล์ .env และแปลงเป็น int
    serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
    if err != nil {
        log.Fatalf("Invalid SERVER_PORT value")
    }

    // ตรวจสอบว่าค่าที่จำเป็นถูกตั้งค่าแล้ว
    if dsn == "" || jwtSecret == "" {
        log.Fatalf("Required environment variables are missing")
    }

    // คืนค่า Config struct ที่มีค่าคอนฟิกทั้งหมด
    return &Config{
        DSN:        dsn,
        JWTSecret:  jwtSecret,
        ServerPort: serverPort,
    }
}
