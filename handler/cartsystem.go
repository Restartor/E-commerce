package handler

import (
	"ecommerce/config"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MenambahkanKeKeranjang(c *gin.Context) {

	userID := c.GetUint("user_id")

	var produkInput struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	// bind json dengan struct produk
	if err := c.ShouldBindJSON(&produkInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// cek apakah produk ada atau tidak agar tidak menambahkan produk yang tidak ada ke keranjang
	var product models.Product
	if err := config.DB.First(&product, produkInput.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false, "error": "product not found",
		})
		return
	}
	// jika stock produk kurang dari quantity yang diminta, kembalikan error
	if produkInput.Quantity > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false, "error": "stock not enough",
		})
		return
	}

	// cek keranjang pengguna
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart = models.Cart{UserID: userID}
		config.DB.Create(&cart)
	}

	// cek apakah produk sudah ada dikeranjang
	var cartItem models.CartItem
	err := config.DB.Where("cart_id = ? AND product_id = ?",
		cart.ID, produkInput.ProductID).First(&cartItem).Error

	if err == nil {
		// kalau produk sudah ada, update quantity
		cartItem.Quantity += produkInput.Quantity
		config.DB.Save(&cartItem)
	} else {
		// kalau produk belum ada, buat item baru
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: produkInput.ProductID,
			Quantity:  produkInput.Quantity,
		}
		config.DB.Create(&cartItem)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "item telah ditambahkan ke keranjang"})

}

func UpdateKeranjang(c *gin.Context) {
	itemID := c.Param("id")

	var keranjangInput struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}
	// bind json dengan struct keranjang
	if err := c.ShouldBindJSON(&keranjangInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// cek apakah item ada di keranjang
	var cartItem models.CartItem
	if err := config.DB.First(&cartItem, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Item tidak ditemukan"})
		return
	}
	// jika input quantity lebih besar dari stock produk, kembalikan error
	var product models.Product
	config.DB.First(&product, cartItem.ProductID)

	if keranjangInput.Quantity > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "stock tidak mencukupi",
		})
		return
	}
	// jika 0 maka hapus item dari keranjang
	if keranjangInput.Quantity == 0 {
		config.DB.Delete(&cartItem)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "item removed from cart!"})
		return
	}

	// update quantity item di keranjang
	cartItem.Quantity = keranjangInput.Quantity
	if err := config.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// respon sukses

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item keranjang berhasil diperbarui",
	})

}

func HapusDariKeranjang(c *gin.Context) {
	itemID := c.Param("id")

	// cek apakah item ada di keranjang
	var cartItem models.CartItem
	if err := config.DB.First(&cartItem, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Item tidak ditemukan",
		})
		return
	}
	// hapus item dari keranjang
	if err := config.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "produk berhasil dihapus",
	})

}
