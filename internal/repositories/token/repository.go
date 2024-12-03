//go:generate mockgen -source=./repository.go -destination=../../../test/mock/repo_token_mock.go -package=mock

package repoToken

import (
	"context"
	"log/slog"
	"main/internal/database"
	"main/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	db     database.DB
	logger *slog.Logger
}

type RepositoryToken interface {
	Login(ctx context.Context, user models.User, refreshToken string, refreshTokenDuration time.Time) error
	SelectToken(ctx context.Context, userID string, token string) (*models.RefreshToken, error)
	DeleteToken(ctx context.Context, tx pgx.Tx, userID string, token string) error
	InsertNewToken(ctx context.Context, user models.User, tx pgx.Tx, refreshToken string, refreshTokenDuration time.Time) error

	StartTx(ctx context.Context) (pgx.Tx, error)
}

func NewUserRepo(db database.DB, logger *slog.Logger) RepositoryToken {
	return &repository{
		db:     db,
		logger: logger,
	}
}
