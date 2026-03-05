package handler

import (
	"ecommerce/config"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome King Admin!!",
	})
}

// Membuat produk baru
func CreateProduct(c *gin.Context) {
	var product models.Product

	// bind json request body ke variabel product,
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// simpan data product ke database
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal membuat produk!",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product berhasil dibuat!",
		"data":    product,
		"id":      product.ID,
	})
}

// Update produk berdasarkan ID
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Product tidak ditemukan!",
		})
		return
	}
	c.ShouldBindJSON(&product)
	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated!",
		"data":    product,
	})

}

// Hapus produk berdasarkan ID
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var products models.Product
	if err := config.DB.Where("id = ?", id).First(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Product tidak ditemukan!",
		})
		return
	}
	if err := config.DB.Delete(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal menghapus product!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data produk berhasil dihapus!",
	})
}
