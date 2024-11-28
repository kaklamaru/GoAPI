package controller

import (
    "RESTAPI/domain/entities"
    "RESTAPI/domain/transaction"
    "RESTAPI/usecase"
    "RESTAPI/pkg" 
    "github.com/gofiber/fiber/v2"
	"time"
	"RESTAPI/infrastructure/jwt"
)

type UserController struct {
    userUsecase    usecase.UserUsecase
    txManager      transaction.TransactionManager
	jwtService     jwt.JWTService
}

func NewUserController(userUsecase usecase.UserUsecase, txManager transaction.TransactionManager,jwtService jwt.JWTService) *UserController {
    return &UserController{
        userUsecase: userUsecase,
        txManager:   txManager,
		jwtService: jwtService,
    }
}

func (c *UserController) RegisterStudent(ctx *fiber.Ctx) error {
    var req struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        TitleName string `json:"title_name"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Phone     string `json:"phone"`
        Code      string `json:"code"`
        Year      uint   `json:"year"`
        BranchID  uint   `json:"branch_id"`
    }

    // รับค่าจาก body
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
    }

    // แฮช password ก่อนเก็บ
    hashedPassword, err := pkg.HashPassword(req.Password)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
    }

    // เริ่มต้นธุรกรรม
    tx := c.txManager.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            if err := tx.Commit(); err != nil {
                tx.Rollback()
            }
        }
    }()

    // สร้าง user
    user := &entities.User{
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "student",
    }

    // สร้าง student โดยตั้ง `UserID` ที่เชื่อมโยงกับ user ที่สร้างขึ้น
    student := &entities.Student{
        TitleName: req.TitleName,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Phone:     req.Phone,
        Code:      req.Code,
        Year:      req.Year,
        BranchID:  req.BranchID,
        UserID:    user.UserID, // ระบุ `UserID` จาก `user`
    }

    // เรียกใช้ UserUsecase เพื่อสร้าง User และ Student พร้อมกัน
    if err := c.userUsecase.RegisterUserAndStudent(tx, user, student); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user and student"})
    }

    // ยืนยันธุรกรรม
    return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Student registered successfully"})
}

func (c *UserController) RegisterTeacher(ctx *fiber.Ctx) error {
    var req struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        TitleName string `json:"title_name"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Phone     string `json:"phone"`
        Code      string `json:"code"`
    }

    // รับค่าจาก body
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
    }

    // แฮช password ก่อนเก็บ โดยใช้ฟังก์ชันของคุณ
    hashedPassword, err := pkg.HashPassword(req.Password)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
    }

    // เริ่มต้นธุรกรรม
    tx := c.txManager.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback() // ใช้ rollback เมื่อเกิดข้อผิดพลาด
        } else {
            if err := tx.Commit(); err != nil {
                tx.Rollback() // ใช้ rollback เมื่อ commit ล้มเหลว
            }
        }
    }()

    // สร้าง user
    user := &entities.User{
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "teacher", // ระบุ role เป็น "teacher"
    }

    // สร้าง teacher
    teacher := &entities.Teacher{
        TitleName: req.TitleName,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Phone:     req.Phone,
        Code:      req.Code,
    }

    // เรียกใช้ UserUsecase เพื่อสร้าง User และ Teacher พร้อมกัน
    if err := c.userUsecase.RegisterUserAndTeacher(tx, user, teacher); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user and teacher"})
    }

    // ส่งผลลัพธ์สำเร็จ
    return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Teacher registered successfully"})
}



func (c *UserController) Login(ctx *fiber.Ctx) error {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var request loginRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	user, err := c.userUsecase.GetUserByEmail(request.Email)
	if err != nil || user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if !pkg.CheckPasswordHash(request.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	token, err := c.jwtService.GenerateJWT(user.UserID, user.Role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set JWT token as a cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",                        // ชื่อคุกกี้
		Value:    token,                          // ค่า JWT
		Expires:  time.Now().Add(24 * time.Hour), // วันหมดอายุ (24 ชั่วโมง)
		HTTPOnly: false,                           // ป้องกันการเข้าถึงผ่าน JavaScript
		Secure:   false,             // ใช้งานเฉพาะ HTTPS (แนะนำสำหรับ Production)
		SameSite: "Lax", // นโยบาย SameSite
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"role":    user.Role,
	})
}

