//go:generate mockgen -source=./repository.go -destination=../../../test/mock/repo_email_mock.go -package=mock

package repoEmail_mock

import (
	"log/slog"
)

type emailRepo struct {
	logger *slog.Logger
}

type RepositoryEmail interface {
	Send(ip string) error
}

func NewEmailRepo(logger *slog.Logger) RepositoryEmail {
	return &emailRepo{
		logger: logger,
	}
}

func (e *emailRepo) Send(ip string) error {
	e.logger.Info("When trying to login to your account your ip has changed, if it is not you write to the support team",
		slog.Any("suspicious ip: ", ip),
	)

	return nil
}
