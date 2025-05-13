package main

import (
	"go-refresh-token/controllers"
	"go-refresh-token/database"
	"go-refresh-token/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Hanya izinkan origin tertentu
		if origin == "http://localhost:3000" { // ganti sesuai kebutuhan
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	database.ConnectDatabase()
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/api/users/register", controllers.Register)
	r.POST("/api/users/auth", controllers.Login)
	r.POST("/api/users/logout", controllers.LogoutUser)
	r.POST("/api/users/refresh", controllers.MyRefreshToken)

	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthenticateToken())
	{

		authRoutes.GET("/api/users", controllers.GetUsers)

	}

	r.Run(":5000")
}
