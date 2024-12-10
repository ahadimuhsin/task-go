package models

import (
	"time"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	Id        int    `gorm:"type:int;primaryKey'autoIncrement" json:"id"`
	Role      string `gorm:"type:varchar(10)" json:"role,omitempty"`
	Name      string `gorm:"type:varchar(255)" json:"name,omitempty"`
	Email     string `gorm:"type:varchar(100)" json:"email,omitempty"`
	Password  string `gorm:"type:varchar(100)" json:"password,omitempty"`
	CreatedAt time.Time`json:"created_at,omitempty"`
	UpdatedAt time.Time`json:"updated_at,omitempty"`
	Tasks []Task `gorm:"constraint:OnDelete:CASCADE" json:"tasks,omitempty"` //has many
}

func (u *User) AfterDelete(c *gorm.DB) (err error){
	c.Clauses(clause.Returning{}).Where("user_id = ?", u.Id).Delete(&Task{})
	return
}