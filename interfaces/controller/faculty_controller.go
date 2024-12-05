package controller

import (
	"RESTAPI/domain/entities"
	"RESTAPI/usecase"
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
    // สร้างตัวแปรสำหรับเก็บข้อมูลคณะจาก request body

    faculty := new(entities.Faculty)
    idstr := ctx.Params("id")
    idint, err := strconv.Atoi(idstr)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid id format",
        })
    }
    id := uint(idint)

    faculty.FacultyID=id
    // ตรวจสอบข้อมูลที่รับมาว่าถูกต้องหรือไม่
    if err := ctx.BodyParser(faculty); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    // เรียกใช้ usecase ในการอัปเดตคณะ
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

    // ส่งคืนข้อมูลคณะทั้งหมดในรูปแบบ JSON
    return ctx.Status(fiber.StatusOK).JSON(faculties)
}

// GetFaculty ฟังก์ชันสำหรับดึงข้อมูลคณะตาม ID
func (c *FacultyController) GetFaculty(ctx *fiber.Ctx) error {
    // ดึงค่า ID จากพารามิเตอร์ของ URL
    idstr := ctx.Params("id")
    idint, err := strconv.Atoi(idstr)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid id format",
        })
    }
    id := uint(idint)
    // เรียกใช้ usecase เพื่อดึงข้อมูลคณะตาม ID
    faculty, err := c.usecase.GetFaculty(id)
    if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Faculty not found",
        })
    }

    // ส่งคืนข้อมูลคณะที่พบในรูปแบบ JSON
    return ctx.Status(fiber.StatusOK).JSON(faculty)
}

func (c *FacultyController) DeleteFacultyByID(ctx *fiber.Ctx) error{
    idstr := ctx.Params("id")
    idint,err := strconv.Atoi(idstr)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid id format",
        })
    }
    id := uint(idint)
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