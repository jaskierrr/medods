package repoToken

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repository) DeleteToken(ctx context.Context, tx pgx.Tx, userID string, token string) error {
	queryBuilder := sq.Delete("refresh_tokens").Where(sq.Eq{"user_id": userID, "token": token})
	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error(fmt.Sprintf("repo_token.DeleteToken(): %v", err))
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error(fmt.Sprintf("repo_token.DeleteToken(): %v", err))
		return err
	}

	return nil
}
