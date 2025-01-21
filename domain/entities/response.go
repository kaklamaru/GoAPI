package entities

import "time"

type EventResponse struct {
	EventID        uint   `json:"event_id"`
	EventName      string `json:"event_name"`
	StartDate      string `json:"start_date"`
	StartTime      string `json:"start_time"`
	WorkingHour    uint   `json:"working_hour"`
	SchoolYear     uint   `json:"school_year"`
	Limit          uint   `json:"limit"`
	FreeSpace      uint   `json:"free_space"`
	Location       string `json:"location"`
	Detail         string `json:"detail"`
	Status         bool   `json:"status"`
	BranchIDs      []uint `json:"branches"`
	Years          []uint `json:"years"`
	AllowAllBranch bool   `json:"allow_all_branch"`
	AllowAllYear   bool   `json:"allow_all_year"`
	Creator        struct {
		UserID    uint   `json:"user_id"`
		TitleName string `json:"title_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Code      string `json:"code"`
	} `json:"creator"`
}

type OutsideRequest struct {
	EventName   string `json:"event_name"`
	StartDate   string `json:"start_date"`
	Location    string `json:"location"`
	SchoolYear  uint   `json:"school_year"`
	WorkingHour uint   `json:"working_hour"`
	Intendant   string `json:"intendent"`
}
type StudentResponse struct {
	UserID      uint   `json:"user_id"`
	TitleName   string `json:"title_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Code        string `json:"code"`
	BranchID    uint   `json:"branch_id"`
	BranchName  string `json:"branch_name"`
	FacultyID   uint   `json:"faculty_id"`
	FacultyName string `json:"faculty_name"`
}

type OutsideResponse struct {
	EventID     uint            `json:"event_id"`
	EventName   string          `json:"event_name"`
	Location    string          `json:"location"`
	StartDate   time.Time       `json:"start_date"`
	SchoolYear  uint            `json:"school_year"`
	WorkingHour uint            `json:"working_hour"`
	Intendant   string          `json:"intendent"`
	Student     StudentResponse `json:"student"`
}

type MyChecklist struct {
	EventID   uint   `json:"event_id"`
	UserID    uint   `json:"user_id"`
	TitleName string `json:"title_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Code      string `json:"code"`
	Certifier uint   `json:"certifier"`
	Status    bool   `json:"status"`
	Comment   string `json:"comment"`
	FilePDF   string `json:"file_pdf"`
}

type MyOutside struct {
	EventID     uint   `json:"event_id"`
	EventName   string `json:"event_name"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	StartTime   string `json:"start_time"`
	WorkingHour uint   `json:"working_hour"`
	SchoolYear  uint   `json:"school_year"`
	Intendant   string `json:"intendent"`
	// Status      bool   `json:"status"`
	// Comment     string `json:"comment"`
}

type MyInside struct {
	EventID     uint   `json:"event_id"`
	EventName   string `json:"event_name"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	StartTime   string `json:"start_time"`
	WorkingHour uint   `json:"working_hour"`
	SchoolYear  uint   `json:"school_year"`
	Status      bool   `json:"status"`
	Comment     string `json:"comment"`
}
