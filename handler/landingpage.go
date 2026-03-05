package handler

import (
	"ecommerce/config"
	"ecommerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// get all produk untuk landing page
func GetProducts(c *gin.Context) {
	var products []models.Product

	// ambil query parameter untuk pagination, jika tidak ada set default page ke 1 dan limit ke 5
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")

	page, err := strconv.Atoi(pageStr) // jika terjadi error atau page kurang dari 1, set page ke 1
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr) // jika terjadi error atau limit kurang dari 1, set limit ke 5
	if err != nil || limit < 1 {
		limit = 5
	}

	if limit > 50 {
		limit = 50
	} // batasi limit maksimal 100 untuk mencegah beban server yang terlalu berat

	offset := (page - 1) * limit
	if err := config.DB.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Gagal mengambil produk!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"page":    page,
		"limit":   limit,
		"data":    products,
	})
}
