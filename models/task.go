package models

import "time"

type Task struct {
	Id           int       `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserId       int       `gorm:"type:int" json:"user_id"`
	Title        string    `gorm:"type:varchar(255)" json:"title"`
	Description  string    `gorm:"type:text" json:"description"`
	Status       string    `gorm:"type:varchar(50)" json:"status"`
	Reason       string    `gorm:"type:text;default:" json:"reason"`
	Revision     int8      `gorm:"type:int;default:0" json:"revision"`
	DueDate      string    `gorm:"type:varchar(100)" json:"due_date"`
	SubmitDate   string    `gorm:"type:varchar(100)" json:"submit_date"`
	RejectedDate string    `gorm:"type:varchar(100)" json:"rejected_date"`
	ApprovedDate string    `gorm:"type:varchar(100)" json:"approved_date"`
	Attachment   string    `gorm:"type:varchar(255)" json:"attachment"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	User         User      `gorm:"foreignKey:user_id" json:"user,omitempty"` //belongs to
}
