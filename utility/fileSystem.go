package utility

// import (
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"os"
// )

// func SaveFile(file *multipart.FileHeader, eventID uint, userID uint) (string, error) {
//     // เปิดไฟล์ที่อัปโหลด
//     srcFile, err := file.Open()
//     if err != nil {
//         return "", fmt.Errorf("failed to open file: %w", err)
//     }
//     defer srcFile.Close()

//     // กำหนดโฟลเดอร์การอัปโหลด
//     uploadDir := "./uploads"
//     if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
//         err = os.Mkdir(uploadDir, os.ModePerm)
//         if err != nil {
//             return "", fmt.Errorf("failed to create upload directory: %w", err)
//         }
//     }

//     // สร้างโฟลเดอร์สำหรับผู้ใช้
//     userDir := fmt.Sprintf("%s/%d", uploadDir, userID)
//     if _, err := os.Stat(userDir); os.IsNotExist(err) {
//         err = os.Mkdir(userDir, os.ModePerm)
//         if err != nil {
//             return "", fmt.Errorf("failed to create user directory: %w", err)
//         }
//     }

//     // กำหนด path ของไฟล์ที่ต้องการบันทึก
//     filePath := fmt.Sprintf("%s/event_%d.pdf", userDir, eventID)

//     // สร้างไฟล์ปลายทาง
//     dstFile, err := os.Create(filePath)
//     if err != nil {
//         return "", fmt.Errorf("failed to create destination file: %w", err)
//     }
//     defer dstFile.Close()

//     // คัดลอกเนื้อหาไฟล์จากต้นทางไปยังปลายทาง
//     _, err = io.Copy(dstFile, srcFile)
//     if err != nil {
//         return "", fmt.Errorf("failed to copy file content: %w", err)
//     }

//     // ส่งคืน path ของไฟล์ที่บันทึกสำเร็จ
//     return filePath, nil
// }

import (
    "github.com/google/uuid"
    "io"
    "mime/multipart"
    "os"
    "fmt"
)

func SaveFile(file *multipart.FileHeader, userID uint) (string, error) {
    // เปิดไฟล์ที่อัปโหลด
    srcFile, err := file.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open file: %w", err)
    }
    defer srcFile.Close()

    // กำหนดโฟลเดอร์อัปโหลด
    uploadDir := "./uploads"
    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        err = os.Mkdir(uploadDir, os.ModePerm)
        if err != nil {
            return "", fmt.Errorf("failed to create upload directory: %w", err)
        }
    }

    // โฟลเดอร์สำหรับผู้ใช้
    userDir := fmt.Sprintf("%s/%d", uploadDir, userID)
    if _, err := os.Stat(userDir); os.IsNotExist(err) {
        err = os.Mkdir(userDir, os.ModePerm)
        if err != nil {
            return "", fmt.Errorf("failed to create user directory: %w", err)
        }
    }

    // สร้าง UUID สำหรับชื่อไฟล์
    uniqueID := uuid.New().String()
    filePath := fmt.Sprintf("%s/%s.pdf", userDir, uniqueID)

    // สร้างหรือเขียนทับไฟล์ใหม่
    dstFile, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to create destination file: %w", err)
    }
    defer dstFile.Close()

    // คัดลอกเนื้อหาไฟล์จากต้นทางไปยังปลายทาง
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return "", fmt.Errorf("failed to copy file content: %w", err)
    }

    // ส่งคืน path ของไฟล์ที่บันทึกสำเร็จ
    return filePath, nil
}

