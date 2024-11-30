package server

import (
	"RESTAPI/domain/repository"
	"RESTAPI/domain/transaction"
	"RESTAPI/infrastructure/database"
	"RESTAPI/infrastructure/jwt"
	"RESTAPI/infrastructure/middleware"
	"RESTAPI/interfaces/controller"
	"RESTAPI/usecase"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5"
)

// SetupRoutes ฟังก์ชันสำหรับกำหนดเส้นทางทั้งหมด
func SetupRoutes(app *fiber.App, db database.Database, jwtService *jwt.JWTService) {

	txManager := transaction.NewGormTransactionManager(db.GetDb())
	userRepo := repository.NewUserRepository(db.GetDb())
	studentRepo := repository.NewStudentRepository(db.GetDb())
	teacherRepo := repository.NewTeacherRepository(db.GetDb())
	userUsecase := usecase.NewUserUsecase(userRepo, studentRepo, teacherRepo)
	userController := controller.NewUserController(userUsecase, txManager, *jwtService)

	facultyRepo := repository.NewFacultyRepository(db.GetDb())
	facultyUsecase := usecase.NewFacultyUsecase(facultyRepo)
	facultyController := controller.NewFacultyController(facultyUsecase)

	branchRepo := repository.NewBranchRepository(db.GetDb())
	branchUsecase := usecase.NewBranchUsecase(branchRepo)
	branchController := controller.NewBranchController(branchUsecase)

	app.Post("/register/student", userController.RegisterStudent)
	app.Post("/register/teacher", userController.RegisterTeacher)
	app.Post("/login", userController.Login)

	protected := app.Group("/protected", middleware.JWTMiddlewareFromCookie(jwtService))
	admin := protected.Group("/admin", middleware.RoleMiddleware("admin"))
	teacher := protected.Group("/teacher", middleware.RoleMiddleware("teacher"))
	student := protected.Group("/student", middleware.RoleMiddleware("student"))

	// ดึงข้อมูลจาก claim  /protected
	protected.Get("/userbyclaim", userController.GetUserByClaims)
	student.Put("/editstudent", userController.EditStudentByID)
	teacher.Put("/editteacher", userController.EditTeacherByID)
	// เส้นทางสำหรับการเพิ่มคณะ(admin) /protected/admin
	admin.Post("/faculty", facultyController.AddFaculty)
	// เส้นทางสำหรับการแก้ไขคณะ(admin) /protected/admin
	admin.Put("/faculty/:id", facultyController.UpdateFaculty)
	admin.Delete("/faculty/:id", facultyController.DeleteFacultyByID)
	// เส้นทางสำหรับการดึงข้อมูลคณะทั้งหมด
	app.Get("/faculties", facultyController.GetAllFaculties)
	// เส้นทางสำหรับการดึงข้อมูลคณะตาม ID
	app.Get("/faculty/:id", facultyController.GetFaculty)

	// เพิ่ม (admin) /protected/admin
	admin.Post("/branch", branchController.AddBranch)
	// ดึงสาขาทั้งหมด
	app.Get("/branches", branchController.GetAllBranches)
	// ดึงบางสาขา
	app.Get("/branch/:id", branchController.GetBranch)
	// ดึงตามรหัสคณะ
	app.Get("/branchbyfaculty/:id", branchController.GetBranchesByFaculty)
	// แก้ไข (admin) /protected/admin
	admin.Put("/branch/:id", branchController.UpdateBranch)
	admin.Delete("/branch/:id", branchController.DeleteBranchByID)

}
