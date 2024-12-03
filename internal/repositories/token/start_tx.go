package repoToken

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) StartTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.GetConn().
		BeginTx(
			ctx,
			pgx.TxOptions{
				IsoLevel:   pgx.Serializable,
				AccessMode: pgx.ReadWrite,
			})

	if err != nil {
		r.logger.Error(fmt.Sprintf("repo_token.StartTx(): %v", err))
		return nil, err
	}

	return tx, nil
}
