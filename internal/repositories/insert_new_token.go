package repo

import (
	"context"
	"log/slog"
	"main/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repository) InsertNewToken(ctx context.Context, user models.User, tx pgx.Tx, refreshToken string, refreshTokenDuration time.Time) error {
	queryBuilder := sq.Insert("refresh_tokens").Columns("user_id", "ip", "token", "expiration_time").Values(user.ID, user.IP, refreshToken, refreshTokenDuration)

	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	r.logger.Debug("Success INSERT refresh token in storage",
		slog.Any("user_id", user.ID),
		slog.Any("ip", user.IP),
	)

	return nil
}
