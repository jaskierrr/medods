package database

import (
	"context"
	"fmt"
	"log/slog"
	"main/config"

	"github.com/jackc/pgx/v5"
)

const dbConfigString = "postgres://%s:%s@%s:%s/%s"

type db struct {
	logger *slog.Logger
	conn   *pgx.Conn
}

type DB interface {
	NewConn(ctx context.Context, config config.Config, logger *slog.Logger) (DB, error)
	GetConn() *pgx.Conn
}

func NewDB() DB {
	return new(db)
}

func (d *db) NewConn(ctx context.Context, config config.Config, logger *slog.Logger) (DB, error) {
	connString := fmt.Sprintf(dbConfigString,
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	logger.Debug("New Postgres connection opened")

	return &db{
		logger: logger,
		conn:   conn,
	}, nil
}

func (d *db) GetConn() *pgx.Conn {
	return d.conn
}
