package service

import (
	"fmt"
	"main/internal/models"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const RefreshTokenSeparator = "_"

func (s *service) NewRefreshToken(user models.User) (string, error) {
	// сторим тело refresh токена в формате: "userID_userIP"
	tokenBodyBuilder := new(strings.Builder)
	_, err := tokenBodyBuilder.WriteString(user.ID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.newRefreshToken(): %v", err))
		return "", err
	}
	_, err = tokenBodyBuilder.WriteString(RefreshTokenSeparator)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.newRefreshToken(): %v", err))
		return "", err
	}
	_, err = tokenBodyBuilder.WriteString(user.IP)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.newRefreshToken(): %v", err))
		return "", err
	}
	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(tokenBodyBuilder.String()), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.newRefreshToken(): %v", err))
		return "", err
	}

	return string(refreshTokenHash), nil
}
