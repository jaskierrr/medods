package bootstrapper

import (
	"context"
	"log"
	"log/slog"
	"main/config"
	"main/internal/controller"
	"main/internal/database"
	"main/internal/handlers"
	logger "main/internal/lib/logger"
	repo "main/internal/repositories"
	"main/internal/service"
	"net/http"
)

var Secret string = ""

type RootBootstrapper struct {
	Infrastructure struct {
		Logger *slog.Logger
		Server *http.Server
		DB     database.DB
	}
	Controller controller.Controller
	Config     *config.Config
	Handlers   handlers.Handlers
	Repository repo.Repository
	Service    service.Service
}

type RootBoot interface {
	registerRepositoriesAndServices(ctx context.Context, db database.DB)
	registerAPIServer(cfg config.Config) error
	RunAPI() error
}

func New() RootBoot {
	return &RootBootstrapper{
		Config: config.NewConfig(),
	}
}

func (r *RootBootstrapper) RunAPI() error {
	Secret = r.Config.Secret
	ctx := context.Background()
	r.Infrastructure.Logger = logger.NewLogger()

	r.registerRepositoriesAndServices(ctx, r.Infrastructure.DB)
	err := r.registerAPIServer(*r.Config)
	if err != nil {
		return err
	}

	return nil
}

func (r *RootBootstrapper) registerRepositoriesAndServices(ctx context.Context, db database.DB) {
	logger := r.Infrastructure.Logger
	r.Infrastructure.DB = database.NewDB().NewConn(ctx, *r.Config, logger)
	r.Repository = repo.NewUserRepo(r.Infrastructure.DB, logger)
	r.Service = service.New(r.Repository, logger, r.Config.Secret, r.Config.AccessTokenTTL, r.Config.RefreshTokenTTL)
}

func (r *RootBootstrapper) registerAPIServer(cfg config.Config) error {
	logger := r.Infrastructure.Logger

	r.Controller = controller.New(r.Service, logger)

	r.Handlers = handlers.New(r.Controller, logger)
	mux := http.NewServeMux()
	r.Handlers.Link(mux)
	if r.Handlers == nil {
		log.Fatal("handlers initialization failed")
	}

	log.Printf("Serving server at http://127.0.0.1:%v", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		return err
	}

	return nil
}
