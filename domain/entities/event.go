package entities

import "time"

type Event struct {
	EventID        uint          `gorm:"primaryKey;autoIncrement" json:"event_id"`
	EventName      string        `gorm:"not null" json:"event_name"`
	Creator        uint          `gorm:"not null" json:"creator"`    // Foreign Key เชื่อมโยงกับ Teacher.UserID
	StartDate      time.Time     `gorm:"not null" json:"start_date"` // ใช้ time.Time สำหรับวัน/เดือน/ปี
	WorkingHour    uint          `gorm:"not null" json:"working_hour"`
	Limit          uint          `gorm:"not null" json:"limit"`
	Detail         string        `gorm:"not null" json:"detail"`
	BranchIDs      string        `gorm:"type:json" json:"branches"`
	Years          string        `gorm:"type:json" json:"years"`
	AllowAllBranch bool          `json:"allow_all_branch"`
	AllowAllYear   bool          `json:"allow_all_year"`
	EventInsides   []EventInside `gorm:"foreignKey:EventID"`
	Teacher        Teacher       `gorm:"foreignKey:Creator;references:UserID" json:"teacher"` // ความสัมพันธ์กับ Teacher

}

type EventInside struct {
	EventID   uint
	UserID    uint
	Certifier uint
	Status    bool
	Comment   string
	Image1    string
}
type EventOutside struct {
	EventID   uint `gorm:"primaryKey;autoIncrement" json:"event_id"`
	UserID    uint
	EventName string
	StartDate string
	Location  string
	Intendant string
	Certifier uint
	Status    bool
	Comment   string
	Image1    string
	Image2    string
}
