//go:generate mockgen -source=./service.go -destination=../mocks/service_mock.go -package=mock
package service

import (
	"context"
	"log/slog"
	"main/internal/models"
	repo "main/internal/repositories"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger          *slog.Logger
	repo            repo.Repository
	secret          string
	AccesstokenTTL  int64
	RefreshtokenTTL time.Time
}

type Service interface {
	Login(ctx context.Context, user models.User) (models.Response, error)
	Refresh(ctx context.Context, req models.RefreshRequest) (models.Response, error)
}

func New(repo repo.Repository, logger *slog.Logger, secret string, AccesstokenTTL int, RefreshtokenTTL int) Service {
	return &service{
		logger:          logger,
		repo:            repo,
		secret:          secret,
		AccesstokenTTL:  time.Now().Add(time.Duration(AccesstokenTTL) * time.Hour).Unix(),
		RefreshtokenTTL: time.Now().Add(time.Duration(RefreshtokenTTL) * time.Minute * 24),
	}
}

func newRefreshToken(user models.User) (string, error) {
	// сторим тело refresh токена в формате: "userID_userIP"
	tokenBodyBuilder := strings.Builder{}
	_, err := tokenBodyBuilder.WriteString(user.ID)
	if err != nil {
		return "", err
	}
	_, err = tokenBodyBuilder.WriteString("_")
	if err != nil {
		return "", err
	}
	_, err = tokenBodyBuilder.WriteString(user.IP)
	if err != nil {
		return "", err
	}
	tokenBody := tokenBodyBuilder.String()

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(tokenBody), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	refreshToken := string(refreshTokenHash)

	return refreshToken, nil
}
