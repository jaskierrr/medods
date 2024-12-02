//go:generate mockgen -source=./repository.go -destination=../../mocks/repo_mock.go -package=mock
package repo

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

type Repository interface {
	Login(ctx context.Context, user models.User, refreshToken string, refreshTokenDuration time.Time) error
	SelectToken(ctx context.Context, req models.RefreshRequest) (models.RefreshToken, error)
	DeleteToken(ctx context.Context, tx pgx.Tx, req models.RefreshRequest) error
	InsertNewToken(ctx context.Context, user models.User, tx pgx.Tx, refreshToken string, refreshTokenDuration time.Time) error

	StartTx(ctx context.Context) (pgx.Tx, error)
}

func NewUserRepo(db database.DB, logger *slog.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}
