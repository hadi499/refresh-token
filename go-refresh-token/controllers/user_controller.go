package controllers

import (
	"go-refresh-token/database"
	"go-refresh-token/middleware"

	"go-refresh-token/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type GetUserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var validate = validator.New()

// Fungsi untuk mengubah error validator ke format yang lebih jelas
func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " harus diisi"
			case "email":
				errors[field] = "Format email tidak valid"
			case "min":
				errors[field] = field + " minimal " + e.Param() + " karakter"
			case "max":
				errors[field] = field + " maksimal " + e.Param() + " karakter"
			default:
				errors[field] = "Format tidak valid"
			}
		}
	}
	return errors
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔥 Validasi otomatis dengan library validator
	if err := validate.Struct(user); err != nil {
		formattedErrors := formatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": formattedErrors})
		return
	}

	// 🔥 Cek apakah username atau email sudah digunakan
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Name).Or("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau Email sudah digunakan"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", inputUser.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	// Gunakan fungsi GenerateToken dari middleware
	accessToken, err := middleware.GenerateAccessToken(dbUser.Id.String(), dbUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(dbUser.Id.String(), dbUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Save refresh token to database
	dbUser.RefreshToken = refreshToken
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	// Membuat response dengan UserResponse
	response := UserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func LogoutUser(c *gin.Context) {
	var reqBody struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil || reqBody.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token harus disediakan"})
		return
	}

	var user models.User
	if err := database.DB.Where("refresh_token = ?", reqBody.Token).First(&user).Error; err != nil {
		c.Status(http.StatusNoContent) // token sudah tidak ada, anggap logout berhasil
		return
	}

	// Gunakan Update agar lebih pasti
	if err := database.DB.Model(&user).Update("refresh_token", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func MyRefreshToken(c *gin.Context) {
	var reqBody struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil || reqBody.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token diperlukan"})
		return
	}

	var user models.User

	// Cek refresh token di database
	if err := database.DB.Where("refresh_token = ?", reqBody.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token tidak valid"})
		return
	}

	// Parse dan verifikasi token
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(reqBody.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return middleware.JwtKey(), nil // gunakan fungsi helper agar tidak expose variabel langsung
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token tidak valid atau expired"})
		return
	}

	// Generate access token baru
	newAccessToken, err := middleware.GenerateAccessToken(claims.UserId, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  newAccessToken,
		"refreshToken": reqBody.Token,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	// Ambil semua user dari database
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	// Konversi ke response struct
	var response []GetUserResponse
	for _, user := range users {
		response = append(response, GetUserResponse{
			ID:        user.Id.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, response)
}
