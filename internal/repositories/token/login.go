package repoToken

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) Login(
	ctx context.Context,
	user models.User,
	refreshToken string,
	refreshTokenDuration time.Time,
) error {
	queryBuilder := sq.Insert("refresh_tokens").
		Columns("user_id", "ip", "token", "expiration_time").
		Values(user.ID, user.IP, refreshToken, refreshTokenDuration)

	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error(fmt.Sprintf("repo_token.Login(): %v", err))
		return err
	}

	_, err = r.db.GetConn().Exec(ctx, sql, args...)
	if err != nil {
		if pgxErr, ok := err.(*pgconn.PgError); ok {
			if pgxErr.Code == "23505" {
				err = fmt.Errorf("access token for this user:%q from this ip:%q already exists", user.ID, user.IP)
			}
		}
		r.logger.Error(fmt.Sprintf("repo_token.Login(): %v", err))
		return err
	}

	r.logger.Debug("Success INSERT refresh token in storage",
		slog.Any("user_id", user.ID),
		slog.Any("ip", user.IP),
	)

	return nil
}
