package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Fungsi middleware untuk autentikasi JWT
func AuthenticateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Ambil token dari header (format: "Bearer <token>")
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenStr := splitToken[1]

		// Parse dan verifikasi token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return JwtKey(), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Simpan data user ke context untuk digunakan di handler berikutnya
		c.Set("user", claims)

		// Lanjut ke handler berikutnya
		c.Next()
	}
}
