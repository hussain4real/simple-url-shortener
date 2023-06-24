package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Shortlies []Shortly
}
