package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `json:"username" gorm:"uniqueindex;not null"`
	Email    string `json:"email" gorm:"uniqueindex;not null"`
	Password string `json:"password" gorm:"not null"`
	Role	 string `json:"role" gorm:"not null;default:'customer'"`
}