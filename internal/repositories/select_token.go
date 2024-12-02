package repo

import (
	"context"
	"log/slog"
	"main/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) SelectToken(ctx context.Context, req models.RefreshRequest) (models.RefreshToken, error) {
	queryBuilder := sq.Select("user_id", "ip", "token", "expiration_time").From("refresh_tokens").Where(sq.Eq{"user_id": req.User.ID, "token": req.RefreshTokenHash})
	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return models.RefreshToken{}, err
	}

	r.logger.Debug("user_id from select token",
		slog.Any("user_id", req.User.ID),
	)

	tokenInfo := models.RefreshToken{}

	err = r.db.GetConn().QueryRow(ctx, sql, args...).Scan(&tokenInfo.UserID, &tokenInfo.UserIP, &tokenInfo.TokenHash, &tokenInfo.ExpirationTime)
	if err != nil {
		return models.RefreshToken{}, err
	}

	r.logger.Debug("token from select token",
		slog.Any("token info", tokenInfo),
	)

	r.logger.Debug("Success SELECT refresh token from storage",
		slog.Any("user_id", tokenInfo.UserID),
		slog.Any("ip", tokenInfo.UserIP),
	)

	return tokenInfo, nil
}
