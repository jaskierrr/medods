package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"main/internal/lib/jwt"
	"main/internal/models"
)

func (s *service) Login(ctx context.Context, user models.User) (*models.Response, error) {
	accessToken, err := jwt.NewAccessToken(user, s.secret, s.accessTokenTTL, s.logger)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.Login(): failed to generate access token: %v", err))
		return nil, err
	}

	refreshToken, err := s.NewRefreshToken(user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.Login(): %v", err))
		return nil, err
	}

	err = s.repo.Login(ctx, user, refreshToken, s.refreshTokenTTL)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.Login(): %v", err))
		return nil, err
	}

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	return &models.Response{
		Access:  accessToken,
		Refresh: refreshTokenBase64,
	}, nil
}
