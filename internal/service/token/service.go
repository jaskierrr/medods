package service

import (
	"context"
	"log/slog"
	"main/internal/models"
	repo_email "main/internal/repositories/email_mock"
	repo_token "main/internal/repositories/token"
	"time"
)

const day = time.Hour * 24

type service struct {
	logger          *slog.Logger
	repo            repo_token.RepositoryToken
	emailRepo       repo_email.RepositoryEmail
	secret          string
	accessTokenTTL  int64
	refreshTokenTTL time.Time
}

type Service interface {
	Login(ctx context.Context, user models.User) (*models.Response, error)
	Refresh(ctx context.Context, req models.RefreshRequest) (*models.Response, error)
	NewRefreshToken(user models.User) (string, error)
}

func New(repo repo_token.RepositoryToken, emailRepo repo_email.RepositoryEmail, logger *slog.Logger, secret string, accessTokenTTL, refreshTokenTTL int) Service {
	return &service{
		logger:          logger,
		repo:            repo,
		emailRepo:       emailRepo,
		secret:          secret,
		accessTokenTTL:  time.Now().Add(time.Duration(accessTokenTTL) * time.Hour).Unix(),
		refreshTokenTTL: time.Now().Add(time.Duration(refreshTokenTTL) * day),
	}
}
