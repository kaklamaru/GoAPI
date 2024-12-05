package entities

// type Faculty struct {
//     FacultyID   uint      `gorm:"primaryKey;autoIncrement" json:"faculty_id"`
//     FacultyCode string    `gorm:"unique;not null" json:"faculty_code"`
//     FacultyName string    `gorm:"unique;not null" json:"faculty_name"`
//     Branches    []Branch  `gorm:"foreignKey:FacultyID" json:"branches"` // One-to-Many relationship
// }

// type Branch struct {
//     BranchID   uint      `gorm:"primaryKey;autoIncrement" json:"branch_id"`
//     BranchCode string    `gorm:"unique;not null" json:"branch_code"`
//     BranchName string    `gorm:"unique;not null" json:"branch_name"`
//     FacultyID  uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"faculty_id"`
//     Faculty    Faculty   `gorm:"foreignKey:FacultyID" json:"faculties"` // Optional: Load Faculty data
//     Students   []Student `gorm:"foreignKey:BranchID"`                      // ความสัมพันธ์แบบ one-to-many

// }
// type Faculty struct {
//     FacultyID   uint      `gorm:"primaryKey;autoIncrement" json:"faculty_id"`
//     FacultyCode string    `gorm:"unique;not null" json:"faculty_code"`
//     FacultyName string    `gorm:"unique;not null" json:"faculty_name"`
//     Branches    []Branch  `gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"branches"`
// }

// type Branch struct {
//     BranchID   uint    `gorm:"primaryKey;autoIncrement" json:"branch_id"`
//     BranchCode string  `gorm:"unique;not null" json:"branch_code"`
//     BranchName string  `gorm:"unique;not null" json:"branch_name"`
//     FacultyID  uint    `gorm:"not null" json:"faculty_id"`  // ใช้ FacultyID เป็น FK
// }
type Branch struct {
    BranchID   uint    `gorm:"primaryKey;autoIncrement" json:"branch_id"`
    BranchCode string  `gorm:"unique;not null" json:"branch_code"`
    BranchName string  `gorm:"unique;not null" json:"branch_name"`
    FacultyID  uint    `gorm:"not null" json:"faculty_id"`
    Faculty    Faculty `gorm:"foreignKey:FacultyID" json:"faculty"` // จะถูก preload โดยอัตโนมัติ
}

type Faculty struct {
    FacultyID   uint    `gorm:"primaryKey;autoIncrement" json:"faculty_id"`
    FacultyCode string  `gorm:"unique;not null" json:"faculty_code"`
    FacultyName string  `gorm:"unique;not null" json:"faculty_name"`
}


