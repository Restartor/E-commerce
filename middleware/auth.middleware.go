package middleware

import (
	"ecommerce/config"
	"ecommerce/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	// penjelasan: fungsi ini mengembalikan sebuah middleware // yang akan digunakan untuk memeriksa token JWT pada header Authorization

	return func(c *gin.Context) {
		// ambil token dari header authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization Header tidak ada!",
			})
			c.Abort()
			return
		}
		// token biasanya dikirim dengan format// "Bearer <token>", jadi kita perlu memisahkan "Bearer"
		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Format token salah/error",
			})
			c.Abort()
			return
		}

		// cek apakah token sudah di-blacklist (logout)
		var blacklisted models.BlacklistedToken
		if err := config.DB.Where("token = ?", tokenString[1]).First(&blacklisted).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token sudah tidak valid, silakan login kembali!",
			})
			c.Abort()
			return
		}

		// parsing token dengan menggunakan jwt.Parse, // dan memberikan fungsi untuk mengambil secret key dari environment variable
		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// jika terjadi error saat parsing token atau token tidak valid, // maka kembalikan response unauthorized
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid",
			})
			c.Abort()
			return
		}

		// jika token valid, maka ambil klaim dari token dan simpan di context
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", userID)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
		}

		c.Next() // lanjutkan ke handler berikutnya
	}
}
