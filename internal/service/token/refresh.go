package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"main/internal/lib/jwt"
	"main/internal/models"
	"time"
)

var (
	ErrWrongRefreshToken = errors.New("wrong refresh token")
	ErrTokenExpired      = errors.New("token expired")
	ErrSendEmail         = errors.New("failed send email")
)

func (s *service) Refresh(ctx context.Context, req models.RefreshRequest) (*models.Response, error) {
	refreshTokenHashByte, err := base64.StdEncoding.DecodeString(req.RefreshToken)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.Refresh(): %v", err))
		return nil, err
	}

	refreshTokenHash := string(refreshTokenHashByte)

	tokenFromDB, err := s.repo.SelectToken(ctx, req.User.ID, refreshTokenHash)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.Refresh(): %v", err))
		return nil, err
	}

	if tokenFromDB.TokenHash != refreshTokenHash {
		s.logger.Error(fmt.Sprintf("service.Refresh(): %v", ErrWrongRefreshToken))
		return nil, ErrWrongRefreshToken
	}

	if tokenFromDB.ExpirationTime.Before(time.Now()) {
		s.logger.Error(fmt.Sprintf("service.Refresh(): %v", ErrTokenExpired))
		return nil, ErrTokenExpired
	}

	// не возвращаю здесь ошибку,а иду дальше т.к. у пользователя просто мог смениться ip
	if tokenFromDB.UserIP != req.User.IP {
		if err := s.emailRepo.Send(req.User.IP); err != nil {
			s.logger.Error(fmt.Sprintf("service.Refresh(): %v", ErrSendEmail))
			return nil, ErrSendEmail
		}
	}

	tx, err := s.repo.StartTx(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service.Refresh(): failed to begin transaction: %v", err))
		return nil, err
	}

	defer func() {
		if err != nil {
		s.logger.Error(fmt.Sprintf("service.Refresh(): failed to write transaction: %v", err))
		_ = tx.Rollback(ctx)
		}
	}()

	err = s.repo.DeleteToken(ctx, tx, req.User.ID, refreshTokenHash)
	if err != nil {
			s.logger.Error(fmt.Sprintf("service.Refresh(): %v", err))
			return nil, err
	}

	refreshToken, err := s.NewRefreshToken(req.User)
	if err != nil {
			s.logger.Error(fmt.Sprintf("service.Refresh(): %v", err))
			return nil, err
	}

	err = s.repo.InsertNewToken(ctx, req.User, tx, refreshToken, s.refreshTokenTTL)
	if err != nil {
			s.logger.Error(fmt.Sprintf("service.Refresh(): %v", err))
			return nil, err
	}

	tx.Commit(ctx)

	accessToken, err := jwt.NewAccessToken(req.User, s.secret, s.accessTokenTTL, s.logger)
	if err != nil {
			s.logger.Error(fmt.Sprintf("service.Refresh(): failed to generate access token: %v", err))
			return nil, err
	}

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	response := models.Response{
		Access:  accessToken,
		Refresh: refreshTokenBase64,
	}

	s.logger.Debug("Success refresh tokens",
		slog.Any("user_id", req.User.ID),
		slog.Any("ip", req.User.IP),
	)

	return &response, nil
}
