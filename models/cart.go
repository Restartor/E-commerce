package models

import "gorm.io/gorm"

// struct untuk cart
type Cart struct{
	gorm.Model
	UserID uint `json:"user_id"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
}
// struct untuk item dalam cart
type CartItem struct{
	gorm.Model
	CartID uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity int `json:"quantity"`
}