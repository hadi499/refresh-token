package main

import (
	"fmt"
	"go-refresh-token/controllers"
	"go-refresh-token/database"
	"go-refresh-token/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	refresh := os.Getenv("REFRESH_TOKEN_SECRET")
	if secret == "" {
		log.Fatal("ACCESS_TOKEN_SECRET is not set")
	}

	fmt.Println("Secret token loaded:", secret)
	fmt.Println("Refresh token loaded:", refresh)

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
