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
	userController := controller.NewUserController(userUsecase, *jwtService, txManager)

	facultyRepo := repository.NewFacultyRepository(db.GetDb())
	facultyUsecase := usecase.NewFacultyUsecase(facultyRepo)
	facultyController := controller.NewFacultyController(facultyUsecase)

	branchRepo := repository.NewBranchRepository(db.GetDb())
	branchUsecase := usecase.NewBranchUsecase(branchRepo)
	branchController := controller.NewBranchController(branchUsecase)

	eventRepo := repository.NewEventRepository(db.GetDb())

	insideRepo := repository.NewEventInsideRepository(db.GetDb())
	eventUsecase := usecase.NewEventUsecase(eventRepo, branchRepo, insideRepo)
	insideUsecase := usecase.NewEventInsideUsecase(insideRepo, userRepo, eventUsecase, txManager)
	eventController := controller.NewEventController(eventUsecase, txManager)
	insideController := controller.NewEventInsideController(insideUsecase, eventUsecase, userUsecase)

	outsideRepo := repository.NewOutsideRepository(db.GetDb())
	outsideUsecase := usecase.NewOutsideUsecase(outsideRepo)
	outsideController := controller.NewOutsideController(outsideUsecase)

	app.Post("/register/student", userController.RegisterStudent)
	app.Post("/register/teacher", userController.RegisterTeacher)
	app.Post("/login", userController.Login)

	protected := app.Group("/protected", middleware.JWTMiddlewareFromCookie(jwtService))
	admin := protected.Group("/admin", middleware.RoleMiddleware("admin"))
	teacher := protected.Group("/teacher", middleware.RoleMiddleware("teacher"))
	student := protected.Group("/student", middleware.RoleMiddleware("student"))
	// adminTeacher := protected.Group("/admin-teacher", middleware.RoleMiddleware("admin", "teacher"))
	admin.Post("/event", eventController.CreateEvent)
	teacher.Post("/event", eventController.CreateEvent)

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

	admin.Put("/status/:id", eventController.StatusEvent)
	teacher.Put("/status/:id", eventController.StatusEvent)
	
	app.Get("/events", eventController.GetAllEvent)
	app.Get("/event/:id", eventController.GetEventByID)

	admin.Put("/event/:id", eventController.EditEvent)
	teacher.Put("/event/:id", eventController.EditEvent)
	admin.Delete("/event/:id", eventController.DeleteEvent)
	teacher.Delete("/event/:id", eventController.DeleteEvent)

	protected.Get("/count/:id", insideController.CountEventInside)

	student.Post("/join/:id", insideController.JoinEvent)
	student.Delete("/unjoin/:id", insideController.UnJoinEventInside)

	student.Post("upload/:id", insideController.UploadFile)

	student.Get("/file/:id", insideController.GetFile)

	student.Post("/outside",outsideController.CreateOutside)
	student.Get("/outside/:id",outsideController.GetOutsideByID)
	student.Get("/download/:id",outsideController.DownloadPDF)
}
