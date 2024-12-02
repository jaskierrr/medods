package service

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/lib/jwt"
	"main/internal/models"
	"time"
)

var (
	ErrWrongRefreshToken = errors.New("wrong refrash token")
	ErrTokenExpired      = errors.New("token expired")
	ErrSendEmail         = errors.New("failed send email")
)

func (s *service) Refresh(ctx context.Context, req models.RefreshRequest) (models.Response, error) {
	tokenFromDB, err := s.repo.SelectToken(ctx, req)
	if err != nil {
		return models.Response{}, err
	}

	s.logger.Info("Token from DB",
		slog.Any("token", tokenFromDB),
	)

	if tokenFromDB.TokenHash != req.RefreshTokenHash {
		return models.Response{}, ErrWrongRefreshToken
	}

	if tokenFromDB.ExpirationTime.Before(time.Now()) {
		s.logger.Info("err: "+ErrTokenExpired.Error())
		// !!! почему здесь пустая ошибка возвращается? хотя сверху в лог пишется
		return models.Response{}, ErrTokenExpired
	}

	if tokenFromDB.UserIP != req.User.IP {
		err := sendEmail(s, req.User.IP)
		if err != nil {
			s.logger.Info("err:"+ErrTokenExpired.Error())
			return models.Response{}, ErrSendEmail
		}
	}

	tx, err := s.repo.StartTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction: " + err.Error())
		return models.Response{}, err
	}

	err = s.repo.DeleteToken(ctx, tx, req)
	if err != nil {
		return models.Response{}, err
	}

	refreshToken, err := newRefreshToken(req.User)
	if err != nil {
		return models.Response{}, err
	}

	err = s.repo.InsertNewToken(ctx, req.User, tx, refreshToken, s.RefreshtokenTTL)
	if err != nil {
		return models.Response{}, err
	}

	accessToken, err := jwt.NewAccessToken(req.User, s.secret, s.AccesstokenTTL)
	if err != nil {
		s.logger.Error("failed to generate access token" + err.Error())
		return models.Response{}, err
	}

	response := models.Response{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	s.logger.Debug("Success refresh tokens",
		slog.Any("user_id", req.User.ID),
		slog.Any("ip", req.User.IP),
	)

	return response, nil
}

func sendEmail(s *service, ip string) error {
	s.logger.Info("When trying to login to your account your ip has changed, if it is not you write to the support team",
		slog.Any("suspicious ip: ", ip),
	)

	return nil
}
