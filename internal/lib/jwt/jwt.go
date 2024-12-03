package jwt

import (
	"fmt"
	"log/slog"
	"main/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func NewAccessToken(user models.User, secret string, duration int64, logger *slog.Logger) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": user.ID,
		"ip":   user.IP,
		"exp":  duration,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Error(fmt.Sprintf("jwt.NewAccessToken(): %v", err))
		return "", err
	}

	return tokenString, nil
}
