package jwt

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"RESTAPI/config"
)

type JWTService struct {
	SecretKey string
}

// NewJWTService สร้าง Service สำหรับจัดการ JWT
func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		SecretKey: cfg.JWTSecret,
	}
}

// GenerateJWT สร้าง JWT Token
func (j *JWTService) GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // กำหนดหมดอายุ 24 ชั่วโมง
	}

	// สร้าง token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เซ็นต์ token และคืนค่า
	return token.SignedString([]byte(j.SecretKey))
}

// ValidateJWT ตรวจสอบ JWT Token
func (j *JWTService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า token เป็น valid และ cast claims ไปยัง jwt.MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
