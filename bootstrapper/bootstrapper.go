package bootstrapper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"main/config"
	"main/internal/controller"
	"main/internal/database"
	"main/internal/handlers"
	logger "main/internal/lib/logger"
	repoEmail "main/internal/repositories/email_mock"
	repoToken "main/internal/repositories/token"
	service "main/internal/service/token"
	"net/http"
)

type RootBootstrapper struct {
	Infrastructure struct {
		Logger *slog.Logger
		Server *http.Server
		DB     database.DB
	}
	Controller controller.Controller
	Config     *config.Config
	Handlers   handlers.Handlers
	RepoToken  repoToken.RepositoryToken
	RepoEmail  repoEmail.RepositoryEmail
	Service    service.Service
}

type RootBoot interface {
	registerRepositoriesAndServices(ctx context.Context) error
	registerAPIServer(cfg config.Config) error
	RunAPI() error
}

func New() RootBoot {
	return &RootBootstrapper{
		Config: config.NewConfig(),
	}
}

func (r *RootBootstrapper) RunAPI() error {
	ctx := context.Background()
	r.Infrastructure.Logger = logger.NewLogger()

	r.registerRepositoriesAndServices(ctx)
	err := r.registerAPIServer(*r.Config)
	if err != nil {
		return err
	}

	return nil
}

func (r *RootBootstrapper) registerRepositoriesAndServices(ctx context.Context) error {
	db, err := database.NewDB().NewConn(ctx, *r.Config, r.Infrastructure.Logger)
	if err != nil {
		return err
	}
	r.Infrastructure.DB = db

	r.RepoToken = repoToken.NewUserRepo(r.Infrastructure.DB, r.Infrastructure.Logger)
	r.RepoEmail = repoEmail.NewEmailRepo(r.Infrastructure.Logger)
	r.Service = service.New(
		r.RepoToken,
		r.RepoEmail,
		r.Infrastructure.Logger,
		r.Config.Secret,
		r.Config.AccessTokenTTL,
		r.Config.RefreshTokenTTL,
	)

	return nil
}

func (r *RootBootstrapper) registerAPIServer(cfg config.Config) error {
	r.Controller = controller.New(r.Service, r.Infrastructure.Logger)

	r.Handlers = handlers.New(r.Controller, r.Infrastructure.Logger)
	mux := http.NewServeMux()
	r.Handlers.Link(mux)
	if r.Handlers == nil {
		return errors.New("handlers initialization failed")
	}

	log.Printf("Serving server at http://127.0.0.1:%v", cfg.ServerPort)

	addr := fmt.Sprintf(":%v", cfg.ServerPort)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return err
	}

	return nil
}
