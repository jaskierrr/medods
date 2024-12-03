package repoToken

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) SelectToken(ctx context.Context, userID string, token string) (*models.RefreshToken, error) {
	queryBuilder := sq.Select("user_id", "ip", "token", "expiration_time").
		From("refresh_tokens").
		Where(sq.Eq{"user_id": userID, "token": token})
	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error(fmt.Sprintf("repo_token.SelectToken(): %v", err))
		return nil, err
	}

	tokenInfo := new(models.RefreshToken)

	err = r.db.GetConn().QueryRow(ctx, sql, args...).
		Scan(
			&tokenInfo.UserID,
			&tokenInfo.UserIP,
			&tokenInfo.TokenHash,
			&tokenInfo.ExpirationTime,
		)
	if err != nil {
		r.logger.Error(fmt.Sprintf("repo_token.SelectToken(): %v", err))
		return nil, err
	}


	r.logger.Debug("Success SELECT refresh token from storage",
		slog.Any("user_id", tokenInfo.UserID),
		slog.Any("ip", tokenInfo.UserIP),
	)

	return tokenInfo, nil
}
