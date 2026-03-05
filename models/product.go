package models

import "gorm.io/gorm"

type Product struct{
	gorm.Model
	Productname string `json:"productname" gorm:"not null" binding:"required"`
	Description string `json:"description" gorm:"not null" binding:"required"`
	Price float64		`json:"price" gorm:"not null" binding:"required"`
	Stock int			`json:"stock" gorm:"not null" binding:"required"`
}
