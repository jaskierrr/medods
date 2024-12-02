package jwt

import (
	"fmt"
	"log"
	"main/internal/models"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type data struct {
	Secret string `envconfig:"secret" required:"true"`
}

func newSecret() *data {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := &data{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal("Failed load envconfig " + err.Error())
	}

	return cfg
}

var secret = newSecret().Secret

func NewAccessToken(user models.User, secret string, duration int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": user.ID,
		"ip": user.IP,
		"exp":  duration,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string) (interface{}, error) {
	bearerToken := strings.Split(tokenStr, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return "token valid", nil
}
