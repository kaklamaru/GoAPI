package entities

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


