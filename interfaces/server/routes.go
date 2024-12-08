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

	eventRepo := repository.NewEventRepository(db.GetDb())
	eventUsecase := usecase.NewEventUsecase(eventRepo,branchRepo)
	eventController := controller.NewEventController(eventUsecase,txManager)

	insideRepo := repository.NewEventInsideRepository(db.GetDb())
	insideUsecase := usecase.NewEventInsideUsecase(insideRepo)
	insideController := controller.NewEventInsideController(insideUsecase,eventUsecase)

	app.Post("/register/student", userController.RegisterStudent)
	app.Post("/register/teacher", userController.RegisterTeacher)
	app.Post("/login", userController.Login)

	protected := app.Group("/protected", middleware.JWTMiddlewareFromCookie(jwtService))
	admin := protected.Group("/admin", middleware.RoleMiddleware("admin"))
	teacher := protected.Group("/teacher", middleware.RoleMiddleware("teacher"))
	student := protected.Group("/student", middleware.RoleMiddleware("student"))



	protected.Get("/userbyclaim", userController.GetUserByClaims)
	student.Put("/personalinfo", userController.EditStudent)
	teacher.Put("/personalinfo", userController.EditTeacher)
	admin.Put("/personalinfo", userController.EditTeacher)
	admin.Put("/studentinfo", userController.EditStudentByID)
	admin.Put("/teacherinfo", userController.EditTeacherByID)
	admin.Post("/faculty", facultyController.AddFaculty)
	admin.Put("/faculty/:id", facultyController.UpdateFaculty)
	admin.Delete("/faculty/:id", facultyController.DeleteFacultyByID)
	app.Get("/faculties", facultyController.GetAllFaculties)
	app.Get("/faculty/:id", facultyController.GetFaculty)
	admin.Post("/branch", branchController.AddBranch)
	app.Get("/branches", branchController.GetAllBranches)
	app.Get("/branch/:id", branchController.GetBranch)
	app.Get("/branchbyfaculty/:id", branchController.GetBranchesByFaculty)
	admin.Put("/branch/:id", branchController.UpdateBranch)
	admin.Delete("/branch/:id", branchController.DeleteBranchByID)
	admin.Get("/students", userController.GetAllStudent)
	admin.Get("/teachers", userController.GetAllTeacher)

	protected.Post("/event", eventController.CreateEvent)
	protected.Get("/events", eventController.GetAllEvent)
	protected.Get("/event/:id", eventController.GetEventByID)
	protected.Put("/event/:id", eventController.EditEvent)
	protected.Delete("/event/:id",eventController.DeleteEvent)
	
	student.Post("/join/:id",insideController.JoinEvent)

}
