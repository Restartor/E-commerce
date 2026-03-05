package handler

import (
	"net/http"
	"strings"

	"ecommerce/config"
	"ecommerce/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// routes untuk user yang sudah login, menggunakan middleware authentication
func Userprofile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID tidak valid",
		})
		return
	}
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Role tidak valid",
		})
		return
	}
	username, _ := c.Get("username")
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Username tidak valid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userIDStr,
		"role":     roleStr,
		"username": usernameStr,
		"message":  "Profile user berhasil diakses!",
	})
}

func UserDaftar(c *gin.Context) {
	// deklarasi variabel input dengan tipe data RegisterInput
	var input models.RegisterInput

	// bind json request body ke variabel input,
	// jika error maka kembalikan response bad request dengan pesan error
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// cek email sudah terdaftar atau belum dengan query ke database,
	// jika sudah terdaftar maka kembalikan response bad request dengan pesan error
	var existingEmail models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingEmail).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email telah terdaftar!",
		})
		return
	}

	// cek username sudah terdaftar atau belum
	var existingUsername models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUsername).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username telah terdaftar!",
		})
		return
	}
	// hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal Hash Password!",
		})
		return
	}
	// buat user baru dengan data dari input dan password yang sudah di hash
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	// simpan user baru ke database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal menyimpan user!",
		})
		return
	}
	// kembalikan response sukses dengan data user yang sudah dibuat
	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"message":  "Berhasil terdaftar!",
	})

}

func UserLogin(c *gin.Context) {
	var input models.LoginInput

	// bind json dengan variable input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// mencari user dengan email input
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email atau password salah!",
		})
		return
	}

	// cek password jika user ada, bandingkan input dgn hash password di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email atau password salah!",
		})
		return
	}

	// generate jwt token jika login berhasil
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 5).Unix(), // token expired dalam 5jam
	}
	// buat token dengan klaim yang sudah dibuat dan signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// cek revoked / logged out token

	// sign token dengan secret key
	tokenString, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
	// jika terjadi error saat signing token,
	// kembalikan response internal server error dengan pesan error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token gagal ter-generate!",
		})
		return
	}
	// response sukses dengan token yang sudah di generate
	c.JSON(http.StatusOK, gin.H{
		"message": "login success!",
		"token":   tokenString,
	})
}

func UserLogout(c *gin.Context) {
	// ambil token dari header authorization
	authHeader := c.GetHeader("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format token salah!",
		})
		return
	}
	tokenString := tokenParts[1]

	// parse token untuk mendapatkan expiry time
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token tidak valid!",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membaca token!",
		})
		return
	}

	// ambil waktu kadaluarsa dari token
	exp, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token tidak memiliki waktu kadaluarsa!",
		})
		return
	}
	expiresAt := time.Unix(int64(exp), 0)

	// simpan token ke blacklist
	blacklistedToken := models.BlacklistedToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	if err := config.DB.Create(&blacklistedToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal logout!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil logout!",
	})
}

func GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cart models.Cart
	var total float64

	// preload adalah fitur gorm untuk memuat relasi data secara otomatis, dalam kasus ini kita memuat relasi CartItems dan Product dalam satu query untuk mengurangi jumlah query ke database

	if err := config.DB.
		Preload("CartItems.Product").Where("user_id = ?", userID).
		First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "cart empty",
		})
		return
	}
	// hitung total harga keranjang menggunakan for loop untuk menjumlahkan harga produk dikalikan dengan quantity untuk setiap item di keranjang
	for _, item := range cart.CartItems {
		total += float64(item.Quantity) * item.Product.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"total":   total,
		"data":    cart,
	})

}
