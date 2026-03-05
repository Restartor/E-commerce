package routes

import (
	"net/http"

	"ecommerce/handler"
	"ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "E-Commerce API Running",
		})
	})

	// user register routes
	router.POST("/register", handler.UserDaftar)

	// user login routes
	router.POST("/login", handler.UserLogin)

	// user profile routes
	router.GET("/profile", middleware.AuthMiddleware(), handler.Userprofile)

	// routes untuk produk, bisa diakses oleh semua user
	router.GET("/products", handler.GetProducts)

	// routes untuk admin
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/dashboard", handler.AdminDashboard)
		admin.POST("/products", handler.CreateProduct)
		admin.PUT("/products/:id", handler.UpdateProduct)
		admin.DELETE("/products/:id", handler.DeleteProduct)
	}

	// routes untuk keranjang, hanya bisa diakses oleh user yang sudah login
	user := router.Group("/cart")
	user.Use(middleware.AuthMiddleware())
	{
		user.POST("/", handler.MenambahkanKeKeranjang)
		user.PUT("/items/:id", handler.UpdateKeranjang)
		user.GET("/", handler.GetCart)
		user.DELETE("/items/:id", handler.HapusDariKeranjang)
	}
	// routes untuk checkout, hanya bisa diakses oleh user yang sudah login
	userCheckout := router.Group("/order")
	userCheckout.Use(middleware.AuthMiddleware())
	{
		userCheckout.POST("/checkout", handler.Checkout)
	}
	router.POST("/payment/notification", handler.MidtransNotification) // untuk midtrans
	router.GET("/my-orders", middleware.AuthMiddleware(), handler.GetMyOrders)

}
