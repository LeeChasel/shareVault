package utils

import (
	"time"

	"github.com/LeeChasel/shareVault/configs"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

func getJWTKey() []byte {
	return []byte(configs.Env.JwtSecret)
}

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTKey())
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTKey(), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
