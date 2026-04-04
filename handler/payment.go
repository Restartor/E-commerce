package handler

import (
	"crypto/sha512"
	"ecommerce/config"
	"ecommerce/models"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func CreatePayment(order models.Order) (string, error) {

	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	s := snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("ORDER-%d", order.ID),
			GrossAmt: int64(order.TotalPrice),
		},
	}
	response, err := snap.CreateTransaction(request)
	if err != nil {
		return "", err
	}
	return response.RedirectURL, nil

}

func extractID(orderID string) uint {
	var id uint
	fmt.Sscanf(orderID, "ORDER-%d", &id)
	return id
}

func MidtransNotification(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload!"})
		return
	}

	orderID, err := payload["order_id"].(string)
	statusCode, err := payload["status_code"].(string)
	grossAmount, err := payload["gross_amount"].(string)
	signatureKey, err := payload["signature_key"].(string)
	transactionStatus, err := payload["transaction_status"].(string)
	paymentType, err := payload["payment_type"].(string)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload!"})
		return
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	// signature formula: SHA512(order_id + status_code + gross_amount + server_key)

	hash := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	expectedSignature := hex.EncodeToString(hash[:])

	if signatureKey != expectedSignature {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature!"})
		return
	}

	var order models.Order
	if err := config.DB.Preload("OrderItems").Where("id = ?", extractID(orderID)).
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found!"})
		return
	}
	fmt.Println("Extracted ID:", extractID(orderID))

	if transactionStatus == "settlement" || transactionStatus == "capture" {
		tx := config.DB.Begin()

		order.Status = "paid"
		order.PaymentType = paymentType

		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order!"})
			return
		}

		for _, item := range order.OrderItems {
			var product models.Product
			if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product stock!"})
				return
			}
			product.Stock -= item.Quantity
			if err := tx.Save(&product).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product stock!"})
				return
			}
		}
		var cart models.Cart
		if err := tx.Where("user_id = ?", order.UserID).First(&cart).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart!"})
			return
		}
		if err := tx.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart!"})
			return
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction!"})
			return
		}

	}
	c.JSON(200, gin.H{"message": "notification received"})
	fmt.Println("STATUS:", statusCode)
	fmt.Println("ORDER ID:", orderID)
}
