package entities

import "time"

type EventResponse struct {
	EventID        uint      `json:"event_id"`
	EventName      string    `json:"event_name"`
	StartDate      time.Time `json:"start_date"`
	WorkingHour    uint      `json:"working_hour"`
	Limit          uint      `json:"limit"`
	FreeSpace      uint      `json:"free_space"`
	Location       string    `json:"location"`
	Detail         string    `json:"detail"`
	Status         bool      `json:"status"`
	BranchIDs      []uint    `json:"branches"`
	Years          []uint    `json:"years"`
	AllowAllBranch bool      `json:"allow_all_branch"`
	AllowAllYear   bool      `json:"allow_all_year"`
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
	BranchName  string `json:"branch_name"`
	FacultyName string `json:"faculty_name"`
}
type OutsideResponse struct {
	EventID     uint            `json:"event_id"`
	EventName   string          `json:"event_name"`
	Location    string          `json:"location"`
	StartDate   time.Time       `json:"start_date"`
	WorkingHour uint            `json:"working_hour"`
	Intendant   string          `json:"intendent"`
	Student     StudentResponse `json:"student"`
}