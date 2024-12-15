package fileSystem

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(file *multipart.FileHeader, eventID uint, userID uint) error {
	// เปิดไฟล์เพื่ออ่านข้อมูล
	srcFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer srcFile.Close()

	// กำหนดโฟลเดอร์การอัปโหลด
	uploadDir := "./uploads"
	// ตรวจสอบว่าโฟลเดอร์ 'uploads' มีหรือไม่ ถ้าไม่มีก็สร้างขึ้น
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	// กำหนดโฟลเดอร์สำหรับผู้ใช้
	userDir := fmt.Sprintf("%s/%d", uploadDir, userID)
	// ตรวจสอบว่าโฟลเดอร์ผู้ใช้มีหรือไม่ ถ้าไม่มีก็สร้างขึ้น
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		err = os.Mkdir(userDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create user directory: %w", err)
		}
	}

	// กำหนด path ของไฟล์ที่ต้องการบันทึก
	filePath := fmt.Sprintf("%s/event_%d.pdf", userDir, eventID)

	// เปิดไฟล์ปลายทางสำหรับเขียนข้อมูล
	dstFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// คัดลอกข้อมูลจากไฟล์ต้นทางไปยังไฟล์ปลายทาง
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}