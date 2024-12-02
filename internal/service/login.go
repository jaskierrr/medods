package service

import (
	"context"
	"main/internal/lib/jwt"
	"main/internal/models"
)

func (s *service) Login(ctx context.Context, user models.User) (models.Response, error) {

	accessToken, err := jwt.NewAccessToken(user, s.secret, s.AccesstokenTTL)
	if err != nil {
		s.logger.Error("failed to generate access token" + err.Error())
		return models.Response{}, err
	}

	refreshToken, err := newRefreshToken(user)
	if err != nil {
		return models.Response{}, err
	}

	err = s.repo.Login(ctx, user, refreshToken, s.RefreshtokenTTL)
	if err != nil {
		return models.Response{}, err
	}

	response := models.Response{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	return response, nil
}
