package repo

import (
	"context"

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
		return nil, err
	}

	return tx, nil
}
