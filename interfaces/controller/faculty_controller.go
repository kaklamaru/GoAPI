package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
	"RESTAPI/utility"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// FacultyController struct สำหรับจัดการคำขอเกี่ยวกับ Faculty
type FacultyController struct {
    usecase usecase.FacultyUsecase  // ใช้งาน usecase เพื่อทำงานกับข้อมูล
}

// NewFacultyController สร้าง instance ของ FacultyController
func NewFacultyController(usecase usecase.FacultyUsecase) *FacultyController {
    return &FacultyController{usecase: usecase}  // ส่งคืน FacultyController พร้อม usecase
}

// AddFaculty ฟังก์ชันสำหรับเพิ่มข้อมูลคณะ
func (c *FacultyController) AddFaculty(ctx *fiber.Ctx) error {
    // สร้างตัวแปรสำหรับเก็บข้อมูลคณะจาก request body
    faculty := new(entities.Faculty)
    
    // ตรวจสอบข้อมูลที่รับมาว่าถูกต้องหรือไม่
    if err := ctx.BodyParser(faculty); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    // เรียกใช้ usecase ในการเพิ่มคณะ
    if err := c.usecase.AddFaculty(faculty); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to create faculty",
        })
    }

    // ส่งคืนสถานะ OK พร้อมข้อมูลคณะที่เพิ่ม
    return ctx.Status(fiber.StatusCreated).JSON(faculty)
}

// UpdateFaculty ฟังก์ชันสำหรับแก้ไขข้อมูลคณะ
func (c *FacultyController) UpdateFaculty(ctx *fiber.Ctx) error {

    faculty := new(entities.Faculty)
    id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

    faculty.FacultyID=id
    // ตรวจสอบข้อมูลที่รับมาว่าถูกต้องหรือไม่
    if err := ctx.BodyParser(faculty); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    if err := c.usecase.UpdateFaculty(faculty); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to update faculty",
        })
    }

    // ส่งคืนสถานะ OK พร้อมข้อมูลคณะที่อัปเดต
    return ctx.Status(fiber.StatusOK).JSON(faculty)
}

// GetAllFaculties ฟังก์ชันสำหรับดึงข้อมูลคณะทั้งหมด
func (c *FacultyController) GetAllFaculties(ctx *fiber.Ctx) error {
    // เรียกใช้ usecase เพื่อดึงข้อมูลคณะทั้งหมด
    faculties, err := c.usecase.GetAllFaculties()
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to retrieve faculties",
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(faculties)
}

func (c *FacultyController) GetFaculty(ctx *fiber.Ctx) error {
    // ดึงค่า ID จากพารามิเตอร์ของ URL
    id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
    // เรียกใช้ usecase เพื่อดึงข้อมูลคณะตาม ID
    faculty, err := c.usecase.GetFaculty(id)
    if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Faculty not found",
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(faculty)
}

func (c *FacultyController) DeleteFacultyByID(ctx *fiber.Ctx) error{
    id, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
    faculty,err := c.usecase.DeleteFacultyByID(id)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":err.Error(),
        })
    }
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Faculty deleted successfully",
        "faculty": faculty, 
    })
}

func (c *FacultyController) AddFacultyStaff(ctx *fiber.Ctx) error{
    facultyID, err := utility.GetUintID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}
	idStr := ctx.Params("userid")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UserID",
		})
	}
	userID := uint(idInt)

    if err:=c.usecase.AddFacultyStaff(facultyID,userID);err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":err,
        })
    }
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "massage":"successfully",
    })
}