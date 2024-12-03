package entities

type Faculty struct {
    FacultyID    uint      `gorm:"primaryKey" json:"faculty_id"`
    FacultyName  string    `gorm:"not null" json:"faculty_name"`
    Branches     []Branch  `gorm:"foreignKey:FacultyID" json:"branches"` // ความสัมพันธ์แบบ one-to-many
}

type Branch struct {
    BranchID    uint       `gorm:"primaryKey" json:"branch_id"` // Primary key ของ Branch
    BranchName  string     `gorm:"not null" json:"branch_name"`
    Students    []Student  `gorm:"foreignKey:BranchID"`         // ความสัมพันธ์แบบ one-to-many
    FacultyID   uint       `gorm:"not null" json:"faculty_id"`  // Foreign key ที่เชื่อมกับ Faculty
    Permissions  []Permission  `gorm:"foreignKey:BranchID"`
}

