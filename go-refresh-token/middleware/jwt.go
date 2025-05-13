package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key untuk JWT
var jwtKey = []byte("my_secret_key")

// Claims adalah struktur untuk menyimpan payload token
type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken membuat token JWT untuk user yang berhasil login
func GenerateAccessToken(userId string, username string) (string, error) {
	expirationTime := time.Now().Add(3 * time.Minute) // Token berlaku selama 3 menit

	claims := &Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Gunakan HS256 sebagai metode signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Buat token string
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userId string, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku selama 24 jam

	claims := &Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Gunakan HS256 sebagai metode signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Buat token string
	tokenString, err := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JwtKey() []byte {
	return jwtKey
}
