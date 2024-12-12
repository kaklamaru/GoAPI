package entities

import "time"

type Event struct {
	EventID        uint      `gorm:"primaryKey;autoIncrement" json:"event_id"`
	EventName      string    `gorm:"not null" json:"event_name"`
	Creator        uint      `gorm:"not null" json:"creator"`
	StartDate      time.Time `gorm:"not null" json:"start_date"`
	WorkingHour    uint      `gorm:"not null" json:"working_hour"`
	FreeSpace      uint      `gorm:"not null" json:"free_space"`
	Detail         string    `json:"detail"`
	BranchIDs      string    `gorm:"type:json" json:"branches"`
	Years          string    `gorm:"type:json" json:"years"`
	AllowAllBranch bool      `json:"allow_all_branch"`
	AllowAllYear   bool      `json:"allow_all_year"`
	Status         bool      `gorm:"default:true" json:"status"` 
	Teacher        Teacher   `gorm:"foreignKey:Creator;references:UserID" json:"teacher"`
}


type EventInside struct {
	EventId   uint    `gorm:"primaryKey" json:"event_id"`
	User      uint    `gorm:"primaryKey" json:"user_id"`
	Event     Event   `gorm:"foreignKey:EventId;references:EventID" json:"event"`
	Student   Student `gorm:"foreignKey:User;references:UserID" json:"student"`
	Certifier uint    `gorm:"default:null" json:"certifier"`
	Teacher   Teacher `gorm:"foreignKey:Certifier;references:UserID" json:"teacher"`
	Status    bool    `json:"status"`
	Comment   string  `json:"comment"`
	FilePDF   string  `gorm:"size:255" json:"file_pdf"`
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

//
type EventResponse struct {
	EventID        uint      `json:"event_id"`
	EventName      string    `json:"event_name"`
	Creator        uint      `json:"creator"`
	StartDate      time.Time `json:"start_date"`
	WorkingHour    uint      `json:"working_hour"`
	Limit          uint      `json:"limit"`
	FreeSpace      uint      `json:"free_space"`
	Detail         string    `json:"detail"`
	BranchIDs      []uint    `json:"branches"`
	Years          []uint    `json:"years"`
	Status         bool      `json:"status"` 
	AllowAllBranch bool      `json:"allow_all_branch"`
	AllowAllYear   bool      `json:"allow_all_year"`
	Teacher        struct {
		UserID    uint   `json:"user_id"`
		TitleName string `json:"title_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Code      string `json:"code"`
	} `json:"teacher"`
}
