package handler

import (
	"ecommerce/config"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {

	userID := c.GetUint("user_id")

	var cart models.Cart
	if err := config.DB.Preload("CartItems.Product").
		Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cart empty!",
		})
		return
	}

	var totalPrice float64
	order := models.Order{
		UserID: userID,
		Status: "pending",
	}

	config.DB.Create(&order)

	for _, item := range cart.CartItems {
		totalPrice += float64(item.Quantity) * item.Product.Price

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}
		config.DB.Create(&orderItem)
	}
	order.TotalPrice = totalPrice
	config.DB.Save(&order)

	token, err := CreatePayment(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "payment creation failed!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "Checkout berhasil, silakan lanjutkan pembayaran!",
		"snap_token":  token,
		"payment_url": token,
	})

}

func GetMyOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	var orders []models.Order
	config.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}
