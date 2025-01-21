package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/domain/transaction"
	"RESTAPI/infrastructure/jwt"
	"RESTAPI/pkg"
	"RESTAPI/usecase"
	"RESTAPI/utility"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userUsecase usecase.UserUsecase
	jwtService  jwt.JWTService
	txManager   transaction.TransactionManager
}

func NewUserController(userUsecase usecase.UserUsecase, jwtService jwt.JWTService, txManager transaction.TransactionManager) *UserController {
	return &UserController{
		userUsecase: userUsecase,
		jwtService:  jwtService,
		txManager:   txManager,
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

	return utility.HandleTransaction(ctx, tx, func() error {
		// สร้าง user และ student
		user := &entities.User{
			Email:    req.Email,
			Password: hashedPassword,
			Role:     "student",
		}
		student := &entities.Student{
			TitleName: req.TitleName,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
			Code:      req.Code,
			Year:      req.Year,
			BranchId:  req.BranchID,
		}

		// เรียกใช้ UserUsecase เพื่อสร้าง User และ Student พร้อมกัน
		if err := c.userUsecase.RegisterUserAndStudent(tx, user, student); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user and student"})
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "registered successfully"})
	})
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
		Role      string `json:"role"`
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

	return utility.HandleTransaction(ctx, tx, func() error {
		req.Role = "teacher"

		// สร้าง user และ teacher
		user := &entities.User{
			Email:    req.Email,
			Password: hashedPassword,
			Role:     req.Role,
		}
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

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "registered successfully"})
	})
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

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",                        // ชื่อคุกกี้
		Value:    token,                          // ค่า JWT
		Expires:  time.Now().Add(24 * time.Hour), // วันหมดอายุ (24 ชั่วโมง)
		HTTPOnly: false,                          // ป้องกันการเข้าถึงผ่าน JavaScript
		Secure:   false,                          // ใช้งานเฉพาะ HTTPS (แนะนำสำหรับ Production)
		SameSite: "Lax",                          // นโยบาย SameSite
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"role":    user.Role,
	})
}

func (c *UserController) GetUserByClaims(ctx *fiber.Ctx) error {
	// ดึง claims จาก context
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	// ดึง user_id จาก claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	// ตรวจสอบ role จาก claims
	role := claims["role"]

	// ใช้ if-else เพื่อตรวจสอบ role
	if role == "student" {
		student, err := c.userUsecase.GetStudentByUserID(userID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve student data",
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(student)

	} else if role == "teacher" || role == "admin" {
		// ถ้า role เป็น "teacher"
		teacher, err := c.userUsecase.GetTeacherByUserID(userID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve teacher data",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(teacher)

	} else if role == "superadmin" {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"massage": "welcome superadmin",
			"role":"superadmin",
		})
	} else {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}
}

func (c *UserController) EditStudent(ctx *fiber.Ctx) error {
	var req struct {
		TitleName string `json:"title_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Code      string `json:"code"`
		Year      uint   `json:"year"`
		BranchID  uint   `json:"branch_id"`
	}

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid userID in claims",
		})
	}
	userID := uint(userIDFloat)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	student := &entities.Student{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
		Year:      req.Year,
		BranchId:  req.BranchID,
		UserID:    userID,
	}
	// เรียกใช้ UserUsecase เพื่อสร้าง User และ Student พร้อมกัน
	if err := c.userUsecase.EditStudentByID(student); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to edit student"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student edited successfully",
	})
}

func (c *UserController) EditStudentByID(ctx *fiber.Ctx) error {
	var req struct {
		TitleName string `json:"title_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Code      string `json:"code"`
		Year      uint   `json:"year"`
		BranchID  uint   `json:"branch_id"`
	}
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	student := &entities.Student{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
		Year:      req.Year,
		BranchId:  req.BranchID,
		UserID:    id,
	}
	// เรียกใช้ UserUsecase เพื่อสร้าง User และ Student พร้อมกัน
	if err := c.userUsecase.EditStudentByID(student); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to edit student"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student edited successfully",
	})
}

func (c *UserController) EditTeacherByID(ctx *fiber.Ctx) error {
	var req struct {
		TitleName string `json:"title_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Code      string `json:"code"`
	}
	id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	teacher := &entities.Teacher{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
		UserID:    id,
	}
	// เรียกใช้ UserUsecase เพื่อสร้าง User และ Student พร้อมกัน
	if err := c.userUsecase.EditTeacherByID(teacher); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to edit teacher"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Teacher edited successfully",
	})
}

func (c *UserController) EditTeacher(ctx *fiber.Ctx) error {
	var req struct {
		TitleName string `json:"title_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Code      string `json:"code"`
	}
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in claims",
		})
	}
	userID := uint(userIDFloat)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	teacher := &entities.Teacher{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
		UserID:    userID,
	}
	// เรียกใช้ UserUsecase เพื่อสร้าง User และ Student พร้อมกัน
	if err := c.userUsecase.EditTeacherByID(teacher); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to edit teacher"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Teacher edited successfully",
	})
}

func (c *UserController) GetAllStudent(ctx *fiber.Ctx) error {
	student, err := c.userUsecase.GetAllStudent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve student",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(student)
}

func (c *UserController) GetAllTeacher(ctx *fiber.Ctx) error {
	teacher, err := c.userUsecase.GetAllTeacher()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve teacher",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(teacher)
}

func (c *UserController) EditRole(ctx *fiber.Ctx) error {
	var req struct {
		Role string `json:"role"`
	}
	userID, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	err = c.userUsecase.EditRole(userID, req.Role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Edit role success",
	})
}
