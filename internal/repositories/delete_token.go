package repo

import (
	"context"
	"main/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repository) DeleteToken(ctx context.Context,tx pgx.Tx, req models.RefreshRequest) error {
	queryBuilder := sq.Delete("refresh_tokens").Where(sq.Eq{"user_id": req.User.ID, "token": req.RefreshTokenHash})
	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
