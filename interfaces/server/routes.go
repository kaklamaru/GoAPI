package server

import (
	"RESTAPI/domain/repository"
	"RESTAPI/domain/transaction"
	"RESTAPI/infrastructure/database"
	"RESTAPI/infrastructure/jwt"
	"RESTAPI/interfaces/controller"
	"RESTAPI/interfaces/middleware"
	"RESTAPI/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// SetupRoutes ฟังก์ชันสำหรับกำหนดเส้นทางทั้งหมด
func SetupRoutes(app *fiber.App, db database.Database, jwtService *jwtimpl.JWTService) {
	// เส้นทางสำหรับ health check
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	txManager := transaction.NewGormTransactionManager(db.GetDb())

	// สร้าง Repository
	userRepo := repository.NewGormUserRepository(db.GetDb())
	studentRepo := repository.NewGormStudentRepository(db.GetDb())
	teacherRepo := repository.NewGormTeacherRepository(db.GetDb())

	// สร้าง Usecase
	userUsecase := usecase.NewUserUsecase(userRepo, studentRepo,teacherRepo)
	userController := controller.NewUserController(userUsecase,txManager,*jwtService)

	app.Post("/register/student", userController.RegisterStudent)
	app.Post("/register/teacher", userController.RegisterTeacher)
	app.Post("/login", userController.Login)

	// Protected Routes Group
	protected := app.Group("/protected", middleware.JWTMiddlewareFromCookie(jwtService))

	protected.Get("/profile", func(ctx *fiber.Ctx) error {
	// แปลง claims เป็น jwt.MapClaims
	claims, ok := ctx.Locals("claims").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// ใช้ claims เพื่อนำค่า user_id และ role ออกมา
	userID := claims["user_id"]
	role := claims["role"]

	return ctx.JSON(fiber.Map{
		"user_id": userID,
		"role":    role,
	})
})


	facultyRepo := repository.NewFacultyRepository(db.GetDb())
    facultyUsecase := usecase.NewFacultyUsecase(facultyRepo)
    facultyController := controller.NewFacultyController(facultyUsecase)

	// เส้นทางสำหรับการเพิ่มคณะ
	protected.Post("/faculties", facultyController.AddFaculty)
	// เส้นทางสำหรับการแก้ไขคณะ
	protected.Put("/faculties/:id", facultyController.UpdateFaculty)
	// เส้นทางสำหรับการดึงข้อมูลคณะทั้งหมด
	protected.Get("/faculties", facultyController.GetAllFaculties)
	// เส้นทางสำหรับการดึงข้อมูลคณะตาม ID
	protected.Get("/faculties/:id", facultyController.GetFaculty)

	branchRepo := repository.NewBranchRepository(db.GetDb())
	branchUsecase := usecase.NewBranchUsecase(branchRepo)
	branchController := controller.NewBranchController(branchUsecase)
	// เพิ่ม
	protected.Post("/branch",branchController.AddBranch)
	// ดึงสาขาทั้งหมด
	protected.Get("/branch",branchController.GetAllBranches)
	// ดึงบางสาขา
	protected.Get("/branch/:id",branchController.GetBranch)
	
	protected.Put("/branch/:id", branchController.UpdateBranch)

	
	// protected.Post("/student",stdController.AddStudent)
	// protected.Get("/student/:id",stdController.GetStudentByID)
}