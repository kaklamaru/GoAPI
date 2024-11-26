package entities

type User struct {
	UserID   	uint   	`gorm:"primaryKey;autoIncrement" json:"user_id"`
	Email    	string 	`gorm:"unique;not null" json:"email"`
	Password 	string 	`gorm:"not null" json:"password"`
	Role     	string 	`gorm:"default:'user'" json:"role"`
	Student   	*Student `gorm:"foreignKey:UserID"`
	Teacher   	*Teacher `gorm:"foreignKey:UserID"`
}
type Teacher struct {
	UserID    uint   `gorm:"primaryKey" json:"user_id"`
    TitleName  string `gorm:"not null" json:"title_name"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Phone     string `gorm:"unique" json:"phone"`
	Code      string `json:"code"`
}
type Student struct {
	UserID    uint   `gorm:"primaryKey" json:"user_id"`
    TitleName  string `gorm:"not null" json:"title_name"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Phone     string `gorm:"unique" json:"phone"`
	Code      string `gorm:"unique" json:"code"`
	Year      uint   `gorm:"not null" json:"year"`
    BranchID   uint   `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;references:BranchID" json:"branch_id"` // Foreign key ไปยัง Branch
}

