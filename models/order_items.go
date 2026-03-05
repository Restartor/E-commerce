package models

import "gorm.io/gorm"

type OrderItem struct{
	gorm.Model
	OrderID uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`

}